package handler

/*
================================================================================
FILE: internal/gateway/handler/rest.go
================================================================================

PURPOSE:
Implements REST HTTP handlers for client requests coming through the API Gateway.
Translates HTTP JSON requests into gRPC calls to backend services.

LEARNING GO CONCEPTS:
- Writing net/http or Chi handler functions: `func(w http.ResponseWriter, r *http.Request)`.
- Decoding JSON request bodies (`json.NewDecoder(r.Body).Decode(...)`).
- Encoding JSON response bodies (`json.NewEncoder(w).Encode(...)`).
- Error handling and HTTP status codes.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `RequestRide(w http.ResponseWriter, r *http.Request)`
   - Decode pickup/dropoff coordinates from request body.
   - Extract Rider ID from JWT context (set by auth middleware).
   - Call Dispatch gRPC client `RequestTrip(...)`.
   - Return 201 Created with Trip JSON payload.

2. `GetRideStatus(w http.ResponseWriter, r *http.Request)`
   - Extract trip_id from URL path params.
   - Call Dispatch gRPC client `GetTrip(...)`.
   - Return 200 OK with Trip JSON.

3. `RegisterDriver(w http.ResponseWriter, r *http.Request)`
   - Forward driver registration to Driver State REST Service.

4. `UpdateDriverStatus(w http.ResponseWriter, r *http.Request)`
   - Forward status change (ONLINE/OFFLINE) to Driver State REST Service.
================================================================================
*/

import (
	"encoding/json"
	"net/http"

	"github.com/vikasrathod4299/microservice/proto/dispatch"
	"github.com/vikasrathod4299/microservice/proto/driver"
)

type GatewayRESTHandler struct {
	dispatchClient dispatch.DispatchServiceClient
	driverClient   driver.DriverServiceClient
}

func NewGatewayRESTHandler(dispatchClient dispatch.DispatchServiceClient, driverClient driver.DriverServiceClient) *GatewayRESTHandler {
	return &GatewayRESTHandler{
		dispatchClient: dispatchClient,
		driverClient:   driverClient,
	}
}

type CreateRideReq struct {
	PickupLat  float64 `json:"pickup_lat"`
	PickupLng  float64 `json:"pickup_lng"`
	DropoffLat float64 `json:"dropoff_lat"`
	DropoffLng float64 `json:"dropoff_lng"`
}

// RequestRide handles POST /api/rides
func (h *GatewayRESTHandler) RequestRide(w http.ResponseWriter, r *http.Request) {
	var req CreateRideReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `"{"error": invalid request body}`, http.StatusBadRequest)
		return
	}

	rider_id, ok := r.Context().Value("user_id").(string)

	if ok || rider_id == "" {
		rider_id = "demo_id_123"
	}

	resp, err := h.dispatchClient.RequestTrip(r.Context(), &dispatch.TripRequest{
		RiderId:    rider_id,
		PickupLat:  req.PickupLat,
		PickupLng:  req.PickupLng,
		DropoffLat: req.DropoffLat,
		DropoffLng: req.DropoffLng,
	})
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type RequestDriverReq struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	VehicalMake  string `json:"vehical_make"`
	VehicalModel string `json:"vehical_model"`
	LisencePlate string `json:"lisence_plate"`
}
