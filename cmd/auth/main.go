package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vikasrathod4299/microservice/internal/auth/handler"
	authRepo "github.com/vikasrathod4299/microservice/internal/auth/repository/db"
	"github.com/vikasrathod4299/microservice/internal/auth/service"
	"github.com/vikasrathod4299/microservice/pkg/config"
	authPb "github.com/vikasrathod4299/microservice/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const authAddr = ":50053"

func main() {
	cfg := config.LoadConfig()

	if len([]byte(cfg.JWTSecret)) < 32 {
		log.Fatal("JWTSecret must contain at least 32 bytes")
	}

	ctx, stopSignles := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignles()

	databaseCtx, cancelDatabase := context.WithTimeout(ctx, 5*time.Second)
	defer cancelDatabase()

	database, err := pgxpool.New(databaseCtx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal("create Auth Postgres pool: %v", err)
	}
	defer database.Close()

	if err := database.Ping(databaseCtx); err != nil {
		log.Fatalf(
			"connect to Auth PostgreSQL database: %v",
			err,
		)
	}

	repo := authRepo.New(database)
	service, err := service.NewAuthService(repo, cfg.JWTSecret)
	if err != nil {
		log.Fatal("create Auth service: %v", err)
	}
	handler := handler.NewAuthGrpcService(service)

	listner, err := net.Listen("tcp", authAddr)
	if err != nil {
		log.Fatal("listen on %s: %v", authAddr, err)
	}
	defer listner.Close()

	grpcServer := grpc.NewServer()
	authPb.RegisterAuthServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("Auth gRPC service listening on: %v", authAddr)
		serveErr <- grpcServer.Serve(listner)
	}()

	select {
	case <-ctx.Done():
		log.Print("shutting down Auth service")
	case err := <-serveErr:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("serve Auth gRPC service: %v", err)
		}
		return
	}

	gracefullStopDone := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(gracefullStopDone)
	}()

	select {
	case <-gracefullStopDone:
		log.Print("Auth service stopped cleanly")
	case <-time.After(5 * time.Second):
		log.Print("Auth graceful shutdown timeout; forcing stop")
		grpcServer.Stop()
		<-gracefullStopDone
	}
	if err := <-serveErr; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Printf("Auth gRPC server stopped with error: %v", err)
	}
}
