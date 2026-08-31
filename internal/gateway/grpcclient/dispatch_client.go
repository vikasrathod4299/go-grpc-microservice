package grpcclient

/*
================================================================================
FILE: internal/gateway/grpcclient/dispatch_client.go
================================================================================

PURPOSE:
gRPC client wrapper to communicate with the Dispatch Service.
Used by REST handlers (`POST /api/rides`, `GET /api/rides/:id`) to issue ride commands.

LEARNING GO CONCEPTS:
- Invoking unary gRPC methods (`RequestTrip`, `GetTrip`, `UpdateTripStatus`).
================================================================================
*/

type DispatchClient struct {
	// client dispatchPb.DispatchServiceClient
}

func NewDispatchClient(target string) (*DispatchClient, error) {
	// TODO: Create gRPC connection to target (e.g., "localhost:50052")
	return &DispatchClient{}, nil
}
