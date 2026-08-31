package service

/*
================================================================================
FILE: internal/location/service/location.go
================================================================================

PURPOSE:
Business logic layer for the Location Service. Coordinates input validation
and invokes the Redis Repository layer.

LEARNING GO CONCEPTS:
- Go Interface-driven design (Dependency Inversion Principle).
- Defining interfaces for dependencies (`GeoRepository`).

WHAT YOU NEED TO IMPLEMENT HERE:
1. Define `GeoRepository` interface:
   - `SaveLocation(ctx context.Context, driverID string, lat, lng float64) error`
   - `FindNearby(ctx context.Context, lat, lng, radiusKM float64, limit int) ([]NearbyDriver, error)`
   - `Remove(ctx context.Context, driverID string) error`

2. `LocationService` struct holding `GeoRepository`.

3. Methods:
   - `UpdateLocation(...)`: Validates latitude (-90 to +90) and longitude (-180 to +180), then calls repository.
   - `FindNearbyDrivers(...)`: Validates radius > 0, calls repository.
================================================================================
*/

import (
	"context"
	"errors"
	"math"
	"strings"
)

var (
	ErrDriverIDRequired   = errors.New("driver ID is required")
	ErrInvalidCoordinates = errors.New("invalid geographic coordinates")
	ErrInvalidRadius      = errors.New("radius must be greater than zero")
	ErrInvalidLimit       = errors.New("limit must be greater than zero")
)

type NearbyDriver struct {
	DriverID   string
	Latitude   float64
	Longitude  float64
	DistanceKM float64
}

type GeoRepository interface {
	SaveLocation(ctx context.Context, driverID string, lat, lng float64) error
	FindNearby(ctx context.Context, lat, lng, radiusKM float64, limit int) ([]NearbyDriver, error)
	Remove(ctx context.Context, driverID string) error
}

type LocationService struct {
	repo GeoRepository
}

func NewLocationService(repo GeoRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error {
	if strings.TrimSpace(driverID) == "" {
		return ErrDriverIDRequired
	}
	if !validCoordinates(lat, lng) {
		return ErrInvalidCoordinates
	}
	return s.repo.SaveLocation(ctx, driverID, lat, lng)
}

func (s *LocationService) FindNearbyDrivers(ctx context.Context, lat, lng, radiusKM float64, limit int) ([]NearbyDriver, error) {
	if !validCoordinates(lat, lng) {
		return nil, ErrInvalidCoordinates
	}
	if math.IsNaN(radiusKM) || math.IsInf(radiusKM, 0) || radiusKM <= 0 {
		return nil, ErrInvalidRadius
	}
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	return s.repo.FindNearby(ctx, lat, lng, radiusKM, limit)
}

func (s *LocationService) RemoveDriver(ctx context.Context, driverID string) error {
	if strings.TrimSpace(driverID) == "" {
		return ErrDriverIDRequired
	}
	return s.repo.Remove(ctx, driverID)
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
