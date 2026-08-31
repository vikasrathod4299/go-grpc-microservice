package grpcclient

/*
================================================================================
FILE: internal/gateway/grpcclient/location_client.go
================================================================================

PURPOSE:
gRPC client wrapper to communicate with the Location Service.
Allows API Gateway handlers (e.g. WebSocket handler) to invoke `StreamDriverLocation`.

LEARNING GO CONCEPTS:
- Creating a gRPC client connection (`grpc.NewClient` or `grpc.Dial`).
- Reusing client instances across HTTP request handlers.
================================================================================
*/

type LocationClient struct {
	// client locationPb.LocationServiceClient
}

func NewLocationClient(target string) (*LocationClient, error) {
	// TODO Step 1: Dial target address (e.g., "localhost:50051") with grpc.WithInsecure()
	// TODO Step 2: Initialize locationPb.NewLocationServiceClient(conn)
	return &LocationClient{}, nil
}
