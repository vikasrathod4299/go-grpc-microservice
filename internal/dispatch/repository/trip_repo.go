package repository

/*
================================================================================
FILE: internal/dispatch/repository/trip_repo.go
================================================================================

PURPOSE:
PostgreSQL repository implementation for SQL persistence of trips using `jackc/pgx/v5`.

LEARNING GO CONCEPTS:
- Using SQL driver `pgxpool.Pool` or standard `database/sql`.
- Prepared SQL statements & parameter binding (`$1, $2`).
- Scanning query rows into Go structs.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `PostgresTripRepository` struct holding `*pgxpool.Pool`.

2. `CreateTrip(ctx context.Context, t *service.Trip) error`
   - `INSERT INTO trips (id, rider_id, driver_id, status, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, created_at, updated_at) VALUES ($1, ...)`

3. `GetTripByID(ctx context.Context, id string) (*service.Trip, error)`
   - `SELECT ... FROM trips WHERE id = $1`

4. `UpdateTripStatus(ctx context.Context, id string, status service.TripStatus, driverID string) error`
   - `UPDATE trips SET status = $1, driver_id = COALESCE(NULLIF($2, ''), driver_id), updated_at = NOW() WHERE id = $3`
================================================================================
*/

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	dispatchRepo "github.com/vikasrathod4299/microservice/internal/dispatch/repository/db"
	"github.com/vikasrathod4299/microservice/internal/dispatch/service"
)

type PostgresTripRepository struct {
	queries *dispatchRepo.Queries
}

func NewPostgresTripRepository(db *pgxpool.Pool) *PostgresTripRepository {
	return &PostgresTripRepository{
		queries: dispatchRepo.New(db),
	}
}

func (r *PostgresTripRepository) CreateTrip(ctx context.Context, trip *service.Trip) error {
	createTrip, err := r.queries.CreateTrip(ctx,
		dispatchRepo.CreateTripParams{
			ID:      trip.ID,
			RiderID: trip.RiderID,
			DriverID: pgtype.Text{
				String: trip.DriverID,
				Valid:  true,
			},
			Status:     string(trip.Status),
			PickupLat:  trip.PickupLat,
			PickupLng:  trip.PickupLng,
			DropoffLat: trip.DropoffLat,
			DropoffLng: trip.DropoffLng,
		},
	)
	if err != nil {
		return err
	}

	*trip = *tripFromDatabase(createTrip)

	return nil
}

func (r *PostgresTripRepository) GetTripByID(ctx context.Context, id string) (*service.Trip, error) {
	result, err := r.queries.GetTripByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return tripFromDatabase(result), nil
}

func (r *PostgresTripRepository) UpdateTripStatus(ctx context.Context, id string, status service.TripStatus, driverID string) error {
	_, err := r.queries.UpdateTripStatus(ctx, dispatchRepo.UpdateTripStatusParams{
		ID:       id,
		Status:   string(status),
		DriverID: driverID,
	})

	return err
}

func tripFromDatabase(
	trip dispatchRepo.Trip,
) *service.Trip {
	return &service.Trip{
		ID:         trip.ID,
		RiderID:    trip.RiderID,
		DriverID:   trip.DriverID.String,
		Status:     service.TripStatus(trip.Status),
		PickupLat:  trip.PickupLat,
		PickupLng:  trip.PickupLng,
		DropoffLat: trip.DropoffLat,
		DropoffLng: trip.DropoffLng,
		CreatedAt:  trip.CreatedAt.Time,
		UpdatedAt:  trip.UpdatedAt.Time,
	}
}
