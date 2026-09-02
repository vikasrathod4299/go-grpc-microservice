package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/vikasrathod4299/microservice/internal/location/handler"
	"github.com/vikasrathod4299/microservice/internal/location/repository"
	"github.com/vikasrathod4299/microservice/internal/location/service"
	"github.com/vikasrathod4299/microservice/pkg/config"
	locationPb "github.com/vikasrathod4299/microservice/proto/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const locationAddr = ":50051"

func main() {
	cfg := config.LoadConfig()

	ctx, stopSignles := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignles()

	redisClient := repository.NewRedisClient(cfg.RedisAddr)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Close redis client %v", err)
		}
	}()

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("connect to Redis %s: %v", cfg.RedisAddr, err)
	}

	repository := repository.NewRedisGeoRepository(redisClient)
	service := service.NewLocationService(repository)
	handler := handler.NewLocationGrpcHandler(service)

	listner, err := net.Listen("tcp", locationAddr)
	if err != nil {
		log.Fatalf("listen on %s: %v", locationAddr, err)
	}
	defer func() {
		if err := listner.Close(); err != nil {
			log.Printf("Closed location listner %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	locationPb.RegisterLocationServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Location gRPC service is listening on %s", locationAddr)
		serverErr <- grpcServer.Serve(listner)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Sutting down location service")
	case err := <-serverErr:
		if err != nil && errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("servce Location gRPC service: %v", err)
		}
		return
	}

	gracefullStopDone := make(chan struct{})
	select {
	case <-gracefullStopDone:
		log.Printf("Location service stopped cleanly")
	case <-time.After(5 * time.Second):
		log.Print("graceful shutdown timed out; forcing stop")
		grpcServer.Stop()
		<-gracefullStopDone
	}
	if err := <-serverErr; err != nil {
		if !errors.Is(err, grpc.ErrServerStopped) {
			log.Printf("Location gRPC server stopped with error: %v", err)
		}
	}
}

// func main() {
// 	port := ":50051"
// 	log.Printf("🚀 Starting Location Service (gRPC Spatial Engine) on port %s...\n", port)
//
// 	redisClient := repository.NewRedisClient("localhost:6379")
// 	defer redisClient.Close()
//
// 	repo := repository.NewRedisGeoRepository(redisClient)
// 	locService := service.NewLocationService(repo)
//
// 	listener, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to listen on port %s: %v", port, err)
// 	}
//
// 	grpcServer := grpc.NewServer()
// 	locationPb.RegisterLocationServiceServer(grpcServer, handler.NewLocationGrpcHandler(locService))
//
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//
// 	go func() {
// 		log.Printf("📍 Location gRPC Service active on tcp://localhost%s\n", port)
// 		if err := grpcServer.Serve(listener); err != nil {
// 			log.Fatalf("❌ gRPC Server failed: %v", err)
// 		}
// 	}()
//
// 	<-stop
// 	log.Println("🛑 Shutting down Location Service gracefully...")
// 	grpcServer.GracefulStop()
// 	_ = listener.Close()
// 	log.Println("✅ Location Service stopped cleanly.")
// }
