.PHONY: proto infra gateway location dispatch driverstate help

help:
	@echo "Available commands:"
	@echo "  make infra        - Start Redis, Postgres, and Kafka containers"
	@echo "  make proto        - Generate Go gRPC code from proto files"
	@echo "  make gateway      - Run API Gateway Service (:8080)"
	@echo "  make location     - Run Location Service (:50051)"
	@echo "  make dispatch     - Run Dispatch Service (:50052)"
	@echo "  make driverstate  - Run Driver State Service (:8081)"

infra:
	docker compose up -d

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/location/location.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/dispatch/dispatch.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/driver/driver.proto

gateway:
	go run cmd/gateway/main.go

location:
	go run cmd/location/main.go

dispatch:
	go run cmd/dispatch/main.go

driverstate:
	go run cmd/driverstate/main.go
