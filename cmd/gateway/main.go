package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vikasrathod4299/microservice/internal/gateway/handler"
	"github.com/vikasrathod4299/microservice/internal/gateway/hub"
	customMiddleware "github.com/vikasrathod4299/microservice/internal/gateway/middleware"
	"github.com/vikasrathod4299/microservice/pkg/config"
	dispatchPb "github.com/vikasrathod4299/microservice/proto/dispatch"
	driverPb "github.com/vikasrathod4299/microservice/proto/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("🚀 Starting API Gateway & WebSocket Hub Service on port :8080...")

	cfg := config.LoadConfig()
	log.Printf("⚙️ Loaded configuration (Port: %s, LocationURL: %s)", cfg.Port, cfg.LocationServiceURL)

	dispatchConn, err := grpc.NewClient(cfg.DispatchServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial dispatch service", cfg.DispatchServiceURL, err)
	}
	defer dispatchConn.Close()

	dispatchPb := dispatchPb.NewDispatchServiceClient(dispatchConn)

	driverConn, err := grpc.NewClient(cfg.DispatchServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial driver service", cfg.DispatchServiceURL, err)
	}
	defer driverConn.Close()

	driverPb := driverPb.NewDriverServiceClient(driverConn)

	restHandler := handler.NewGatewayRESTHandler(dispatchPb, driverPb)

	hub := hub.NewHub()
	go hub.Run()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(customMiddleware.CORS)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok", "service":"api-gateway"}`)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/ride", restHandler.RequestRide)
		r.Post("/ride", restHandler.RequestRide)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("API Gateway listening on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down API Gateway gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("API Gateway stopped cleanly.")
}
