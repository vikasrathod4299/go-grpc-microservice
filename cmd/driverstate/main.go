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
)

/*
================================================================================
SERVICE: Driver State Service (Entry Point)
================================================================================

RESPONSIBILITIES:
1. Connect to PostgreSQL Database (driver profiles, vehicles, status).
2. Connect to Location gRPC Service Client (to inform Location Service when a driver goes OFFLINE/ON_TRIP).
3. Instantiate Repository (`internal/driverstate/repository`).
4. Instantiate Service (`internal/driverstate/service`).
5. Instantiate REST Handler & HTTP Router (`internal/driverstate/handler`).
6. Start HTTP Server on port :8081.
7. Handle Graceful Shutdown.

HOW TO RUN THIS SERVICE:
   go run cmd/driverstate/main.go
================================================================================
*/

func main() {
	port := ":8081"
	log.Printf("🚀 Starting Driver State Service on port %s...\n", port)

	// TODO Step 1: Connect to PostgreSQL Database
	// db := repository.NewPostgresDB("postgres://user:secret@localhost:5432/uberclone")
	// defer db.Close()

	// TODO Step 2: Initialize gRPC Client to Location Service
	// locClient := grpcclient.NewLocationClient("localhost:50051")

	// TODO Step 3: Initialize Repository & Service
	// repo := repository.NewDriverRepository(db)
	// driverService := service.NewDriverService(repo, locClient)

	// TODO Step 4: Register HTTP REST routes
	// r := chi.NewRouter()
	// handler := handler.NewDriverHandler(driverService)
	// r.Post("/drivers", handler.CreateDriver)
	// r.Get("/drivers/{id}", handler.GetDriver)
	// r.Put("/drivers/{id}/status", handler.UpdateStatus)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok", "service":"driverstate"}`)
	})

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚘 Driver State HTTP API listening on http://localhost%s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Shutting down Driver State Service gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Graceful shutdown failed: %v", err)
	}

	log.Println("✅ Driver State Service stopped cleanly.")
}
