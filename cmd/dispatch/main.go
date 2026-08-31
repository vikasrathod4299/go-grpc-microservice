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
	"github.com/vikasrathod4299/microservice/internal/dispatch/grpcclient"
	"github.com/vikasrathod4299/microservice/internal/dispatch/handler"
	"github.com/vikasrathod4299/microservice/internal/dispatch/repository"
	"github.com/vikasrathod4299/microservice/internal/dispatch/service"
	"github.com/vikasrathod4299/microservice/pkg/config"
	dispatchPb "github.com/vikasrathod4299/microservice/proto/dispatch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/*
================================================================================
SERVICE: Dispatch (Trip) Service - State Machine Brain (Entry Point)
================================================================================

RESPONSIBILITIES:
1. Connect to PostgreSQL Database (storing trips, state history).
2. Connect to Kafka Broker (for publishing event streams like TripCreated, TripCompleted).
3. Connect to Location gRPC Service Client (to search for available nearby drivers).
4. Instantiate Trip Repository, State Machine, & Service layers.
5. Register Dispatch gRPC Server on port :50052.
6. Handle Graceful Shutdown.

HOW TO RUN THIS SERVICE:
   go run cmd/dispatch/main.go
================================================================================
*/

const dispatchAddress = ":50052"

func main() {
	cfg := config.LoadConfig()
	ctx, stopSingle := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stopSingle()

	databaseCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	database, err := pgxpool.New(
		databaseCtx,
		cfg.PostgresDSN,
	)
	if err != nil {
		log.Fatalf("create PostgreSQL Connection pool: %v", err)
	}
	defer database.Close()

	if err := database.Ping(databaseCtx); err != nil {
		log.Fatalf("connection to PostgreSQL: %v", err)
	}

	locationClient, err := grpcclient.NewLocationClient(cfg.LocationServiceURL)
	if err != nil {
		log.Fatalf("create location gRPC client for %s: %v", cfg.LocationServiceURL, err)
	}
	defer func() {
		if err := locationClient.Close(); err != nil {
			log.Fatalf("close gRPC connection: %v", err)
		}
	}()

	tripRepo := repository.NewPostgresTripRepository(database)
	service := service.NewDispatchService(tripRepo, locationClient)
	handler := handler.NewDispatchGrpcHandler(service)

	listner, err := net.Listen("tcp", dispatchAddress)
	if err != nil {
		log.Fatalf("liston on %s: %v", dispatchAddress, err)
	}
	defer listner.Close()
	grpcServer := grpc.NewServer()

	dispatchPb.RegisterDispatchServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("Dispatch gRPC service listening on: %s", dispatchAddress)
		serveErr <- grpcServer.Serve(listner)
	}()

	select {
	case <-ctx.Done():
		log.Print("shutting down Dispatch service")

	case err := <-serveErr:
		if err != nil &&
			!errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf(
				"serve Dispatch gRPC service: %v",
				err,
			)
		}

		return
	}

	gracefulStopDone := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(gracefulStopDone)
	}()

	select {
	case <-gracefulStopDone:
		log.Print("Dispatch service stopped cleanly")

	case <-time.After(5 * time.Second):
		log.Print(
			"graceful shutdown timed out; forcing stop",
		)

		grpcServer.Stop()
		<-gracefulStopDone
	}

	if err := <-serveErr; err != nil &&
		!errors.Is(err, grpc.ErrServerStopped) {
		log.Printf(
			"Dispatch gRPC server stopped with error: %v",
			err,
		)
	}
}
