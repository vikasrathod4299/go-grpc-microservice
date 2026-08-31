package grpcclient

/*
================================================================================
FILE: internal/driverstate/grpcclient/location_client.go
================================================================================

PURPOSE:
gRPC client wrapper used by Driver State Service to notify Location Service
to purge a driver when they go OFFLINE or ON_TRIP.
================================================================================
*/

import (
	"context"
)

type LocationClient struct {
	// client locationPb.LocationServiceClient
}

func NewLocationClient(target string) (*LocationClient, error) {
	return &LocationClient{}, nil
}

func (c *LocationClient) RemoveDriver(ctx context.Context, driverID string) error {
	// TODO: Call locationPb.RemoveDriver(ctx, &locationPb.DriverID{DriverId: driverID})
	return nil
}
