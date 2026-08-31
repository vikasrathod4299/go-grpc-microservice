package service

/*
================================================================================
FILE: internal/driverstate/service/driver.go
================================================================================

PURPOSE:
Business logic layer for managing Driver profiles and status transitions.
Crucial behavior: When a driver switches status to OFFLINE or ON_TRIP,
this service MUST call Location Service's `RemoveDriver` gRPC method to purge
the driver from the Redis geo index!

LEARNING GO CONCEPTS:
- Inter-service synchronization (HTTP Service -> gRPC Service).

WHAT YOU NEED TO IMPLEMENT HERE:
1. `Driver` struct (ID, Name, Phone, VehicleMake, VehicleModel, LicensePlate, Availability).
2. `DriverRepository` interface.
3. `DriverService` struct.
4. `UpdateDriverStatus(ctx context.Context, driverID string, status string) error`:
   - Save updated status in Postgres DB.
   - If status == "OFFLINE" || status == "ON_TRIP":
     Call `locationClient.RemoveDriver(ctx, driverID)`.
================================================================================
*/

import "context"

type Driver struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	VehicleMake  string `json:"vehicle_make"`
	VehicleModel string `json:"vehicle_model"`
	LicensePlate string `json:"license_plate"`
	Availability string `json:"availability"` // "OFFLINE", "ONLINE", "ON_TRIP"
}

type DriverRepository interface {
	CreateDriver(ctx context.Context, d *Driver) error
	GetDriverByID(ctx context.Context, id string) (*Driver, error)
	UpdateAvailability(ctx context.Context, id string, availability string) error
}

type DriverService struct {
	repo DriverRepository
}

func NewDriverService(repo DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}
