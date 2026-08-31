package grpcclient

import (
	"context"
	"strings"

	locationPb "github.com/vikasrathod4299/microservice/proto/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LocationClient struct {
	connection *grpc.ClientConn
	client     locationPb.LocationServiceClient
}

func NewLocationClient(target string) (*LocationClient, error) {
	target = strings.TrimSpace(target)

	connection, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &LocationClient{
		connection: connection,
		client:     locationPb.NewLocationServiceClient(connection),
	}, nil
}

func (c *LocationClient) FindNearbyDrivers(ctx context.Context, lat, lng, radiusKM float64, limit int32) ([]string, error) {
	response, err := c.client.FindNearbyDrivers(ctx, &locationPb.NearbyDriversRequest{
		Latitude:  lat,
		Longitude: lng,
		RadiusKm:  radiusKM,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	driverIDs := make([]string, len(response.GetDrivers()))

	for i, driver := range response.GetDrivers() {
		driverIDs[i] = driver.GetDriverId()
	}
	return driverIDs, nil
}

func (c *LocationClient) Close() error {
	return c.connection.Close()
}
