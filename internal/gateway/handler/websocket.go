package handler

/*
================================================================================
FILE: internal/gateway/handler/websocket.go
================================================================================

PURPOSE:
Handles persistent WebSocket connections for Drivers and Riders.
- Drivers stream GPS coordinates every 3 seconds -> Gateway forwards to Location Service via gRPC stream.
- Riders & Drivers receive real-time updates (e.g., ride request, driver assigned, driver arrived).

LEARNING GO CONCEPTS:
- Gorilla WebSocket upgrader (`websocket.Upgrader`).
- Concurrency with Goroutines (`go client.readPump()`, `go client.writePump()`).
- Channels for receiving incoming messages.
- Integrating WebSockets with gRPC Client Streams.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `HandleDriverWebSocket(w http.ResponseWriter, r *http.Request)`
   - Upgrade HTTP connection to WebSocket.
   - Extract driver_id from query params or JWT token (`/ws/driver?token=...`).
   - Register connection in WebSocket Hub.
   - Start a gRPC stream `LocationService.StreamDriverLocation()`.
   - Read incoming WebSocket messages in a loop and send to gRPC stream.

2. `HandleRiderWebSocket(w http.ResponseWriter, r *http.Request)`
   - Upgrade connection.
   - Register rider connection in Hub to receive trip updates.
================================================================================
*/

import (
	"net/http"
	// TODO: Uncomment when adding gorilla/websocket:
	// "github.com/gorilla/websocket"
)

type GatewayWSHandler struct {
	// TODO: Add fields for Hub and Location gRPC Client
	// hub *hub.Hub
	// locationClient locationPb.LocationServiceClient
}

func NewGatewayWSHandler() *GatewayWSHandler {
	return &GatewayWSHandler{}
}

// HandleDriverWS handles GET /ws/driver
func (h *GatewayWSHandler) HandleDriverWS(w http.ResponseWriter, r *http.Request) {
	// TODO Step 1: Upgrade HTTP to WebSocket using upgrader.Upgrade(w, r, nil)
	// TODO Step 2: Open bidirectional gRPC stream to Location Service
	// TODO Step 3: Loop reading GPS json pings from websocket, forward to gRPC stream
}

// HandleRiderWS handles GET /ws/rider
func (h *GatewayWSHandler) HandleRiderWS(w http.ResponseWriter, r *http.Request) {
	// TODO Step 1: Upgrade connection to WebSocket
	// TODO Step 2: Register connection in Hub for live trip status pushes
}
