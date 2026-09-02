// Package config provides centralized configuration management for the Uber Clone application.
package config

/*
================================================================================
PACKAGE: pkg/config
================================================================================

PURPOSE:
Centralized configuration management. Reads environment variables (or .env file)
and parses them into a strongly typed Go struct so services don't hardcode URLs, ports, or credentials.

LEARNING GO CONCEPTS:
- Using `os.Getenv(key)` with default fallback values.
- Struct field tags (`env:"PORT"`).
- Single source of truth for config across all microservices.

WHAT YOU NEED TO IMPLEMENT HERE:
1. Define a `Config` struct holding all settings:
   - Port (e.g. "8080")
   - Database URLs (Postgres, Redis)
   - Kafka Brokers (e.g. "localhost:9092")
   - gRPC service endpoints ("localhost:50051", "localhost:50052")
   - JWT Secret key
2. Write a `LoadConfig() *Config` function that populates and returns the struct.
================================================================================
*/

import "os"

type Config struct {
	Port               string
	LocationServiceURL string
	DispatchServiceURL string
	AuthServiceURL     string
	DriverStateURL     string
	RedisAddr          string
	PostgresDSN        string
	KafkaBrokers       string
	JWTSecret          string
}

func LoadConfig() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		LocationServiceURL: getEnv("LOCATION_SERVICE_URL", "localhost:50051"),
		DispatchServiceURL: getEnv("DISPATCH_SERVICE_URL", "localhost:50052"),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "localhost:50053"),
		DriverStateURL:     getEnv("DRIVER_STATE_URL", "localhost:8081"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		PostgresDSN:        getEnv("POSTGRES_DSN", "postgres://uber:secret@localhost:5432/uberclone?sslmode=disable"),
		KafkaBrokers:       getEnv("KAFKA_BROKERS", "localhost:9092"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-uber-clone-jwt-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultValue
}
