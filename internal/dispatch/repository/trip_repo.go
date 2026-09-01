package repository

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
