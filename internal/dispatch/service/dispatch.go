package service

/*
================================================================================
FILE: internal/dispatch/service/dispatch.go
================================================================================

PURPOSE:
Core orchestration logic for trips.
1. Creates trip in PostgreSQL.
2. Calls Location Service via gRPC to find nearby available drivers.
3. Assigns closest driver to trip.
4. Updates trip state via State Machine validator.
5. Publishes event to Kafka topic (`TripCreated`, `TripAssigned`, `TripCompleted`).

LEARNING GO CONCEPTS:
- Microservice orchestration pattern.
- Combining DB transactions, gRPC client calls, and Kafka publishing.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `DispatchService` struct:
   - repo TripRepository
   - locationClient LocationClient
   - publisher EventPublisher

2. `CreateTrip(ctx context.Context, riderID string, pLat, pLng, dLat, dLng float64) (*Trip, error)`
   - Create trip record with status = SEARCHING.
   - Query Location Service for 5 nearby drivers within 3 km.
   - If driver found -> Assign driver, update status to DRIVER_ASSIGNED.
   - Save to Postgres repo.
   - Publish `TripAssignedEvent` to Kafka.
   - Return Trip.
================================================================================
*/

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	nearbyDriverRadiusKM = 3.0
	nearbyDriverLimit    = int32(5)
	locationCallTimeout  = 3 * time.Second
)

var (
	ErrRiderIDRequired        = errors.New("rider ID is required")
	ErrTripIDRequired         = errors.New("trip ID is required")
	ErrInvalidTripCoordinates = errors.New("invalid trip coordinates")
	ErrDriverIDRequired       = errors.New("driver ID is required for assignment")
)

type Trip struct {
	ID         string     `json:"id"`
	RiderID    string     `json:"rider_id"`
	DriverID   string     `json:"driver_id"`
	Status     TripStatus `json:"status"`
	PickupLat  float64    `json:"pickup_lat"`
	PickupLng  float64    `json:"pickup_lng"`
	DropoffLat float64    `json:"dropoff_lat"`
	DropoffLng float64    `json:"dropoff_lng"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) error
	GetTripByID(ctx context.Context, id string) (*Trip, error)
	UpdateTripStatus(ctx context.Context, id string, status TripStatus, driverID string) error
}

type LocationFinder interface {
	FindNearbyDrivers(
		ctx context.Context,
		lat float64,
		lng float64,
		radiusKM float64,
		limit int32,
	) ([]string, error)
}

type DispatchService struct {
	repo           TripRepository
	locationFinder LocationFinder
}

func NewDispatchService(repo TripRepository, locationFinder LocationFinder) *DispatchService {
	return &DispatchService{repo: repo, locationFinder: locationFinder}
}

func (svc *DispatchService) CreateTrip(ctx context.Context, pickupLat, pickupLng, dropoffLat, dropoffLng float64, riderID string) (*Trip, error) {
	if strings.TrimSpace(riderID) == "" {
		return nil, ErrRiderIDRequired
	}

	if !validCoordinates(pickupLat, pickupLng) ||
		!validCoordinates(dropoffLat, dropoffLng) {
		return nil, ErrInvalidTripCoordinates
	}

	tripID, err := generateTripID()
	if err != nil {
		return nil, fmt.Errorf("generate trip ID: %w", err)
	}

	now := time.Now().UTC()

	trip := &Trip{
		ID:         tripID,
		RiderID:    riderID,
		Status:     StatusSearching,
		PickupLat:  pickupLat,
		PickupLng:  pickupLng,
		DropoffLat: dropoffLat,
		DropoffLng: dropoffLng,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := svc.repo.CreateTrip(ctx, trip); err != nil {
		return nil, fmt.Errorf("create trp %w", err)
	}

	locationCtx, cancel := context.WithTimeout(ctx, locationCallTimeout)
	defer cancel()

	driverIDs, err := svc.locationFinder.FindNearbyDrivers(locationCtx, pickupLat, pickupLng, nearbyDriverRadiusKM, nearbyDriverLimit)
	if err != nil {
		return nil, fmt.Errorf("find nearby drivers: %w", err)
	}

	if len(driverIDs) == 0 {
		return trip, nil
	}

	driverID := driverIDs[0]

	if err := ValidateTransition(trip.Status, StatusDriverAssigned); err != nil {
		return nil, err
	}

	trip.DriverID = driverID
	trip.Status = StatusDriverAssigned
	trip.UpdatedAt = time.Now().UTC()

	return trip, nil
}

func (svc *DispatchService) GetTripByID(ctx context.Context, tripID string) (*Trip, error) {
	if strings.TrimSpace(tripID) == "" {
		return nil, ErrTripIDRequired
	}

	trip, err := svc.repo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("get trip: %w", err)
	}

	return trip, nil
}

func (svc *DispatchService) UpdateTripStatus(ctx context.Context, tripID string, nextStatus TripStatus, driverID string) (*Trip, error) {
	if strings.TrimSpace(tripID) == "" {
		return nil, ErrTripIDRequired
	}

	trip, err := svc.repo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("update trip state: %w", err)
	}

	if err := ValidateTransition(trip.Status, nextStatus); err != nil {
		return nil, err
	}

	updateDriverID := ""

	if nextStatus != StatusDriverAssigned {
		if strings.TrimSpace(driverID) == "" {
			return nil, ErrDriverIDRequired
		}

		updateDriverID = driverID
	}

	if err := svc.repo.UpdateTripStatus(ctx, tripID, nextStatus, driverID); err != nil {
		return nil, fmt.Errorf("updating trip status: %w", err)
	}

	trip.Status = nextStatus
	trip.UpdatedAt = time.Now().UTC()

	if updateDriverID != "" {
		trip.DriverID = updateDriverID
	}

	return trip, nil
}

func validCoordinates(lat float64, lng float64) bool {
	if math.IsNaN(lat) ||
		math.IsNaN(lng) ||
		math.IsInf(lat, 0) ||
		math.IsInf(lng, 0) {
		return false
	}

	return lat >= -90 &&
		lat <= 90 &&
		lng >= -180 &&
		lng <= 180
}

func generateTripID() (string, error) {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(value[:]), nil
}
