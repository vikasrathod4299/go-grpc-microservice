package models

/*
================================================================================
PACKAGE: pkg/models
================================================================================

PURPOSE:
Defines shared event schemas published to Kafka (e.g. by Dispatch Service)
and consumed by downstream analytics/notifications or API Gateway.

LEARNING GO CONCEPTS:
- Event-driven message payload structures.
- JSON struct tags for serialization.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `TripCreatedEvent` - Payload when a rider requests a new trip.
2. `TripAssignedEvent` - Payload when a driver is assigned.
3. `TripCompletedEvent` - Payload when a trip is finished.
================================================================================
*/

import "time"

type TripEventHeader struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"` // e.g. "TRIP_CREATED", "TRIP_COMPLETED"
	Timestamp time.Time `json:"timestamp"`
}

type TripCreatedEvent struct {
	TripEventHeader
	TripID    string  `json:"trip_id"`
	RiderID   string  `json:"rider_id"`
	PickupLat float64 `json:"pickup_lat"`
	PickupLng float64 `json:"pickup_lng"`
}

type TripAssignedEvent struct {
	TripEventHeader
	TripID   string `json:"trip_id"`
	DriverID string `json:"driver_id"`
	RiderID  string `json:"rider_id"`
}

type TripCompletedEvent struct {
	TripEventHeader
	TripID     string  `json:"trip_id"`
	RiderID    string  `json:"rider_id"`
	DriverID   string  `json:"driver_id"`
	FinalFare  float64 `json:"final_fare,omitempty"`
	DurationMs int64   `json:"duration_ms"`
}
