package handler

/*
================================================================================
FILE: internal/location/handler/grpc_server.go
================================================================================

PURPOSE:
Implements the generated `LocationServiceServer` gRPC interface from `proto/location/location_grpc.pb.go`.
Receives incoming gRPC calls from the API Gateway and Dispatch Service.

LEARNING GO CONCEPTS:
- Implementing gRPC server interfaces in Go.
- Handling bidirectional gRPC streaming (`StreamDriverLocation`).
- Receiving stream messages in a `for` loop with `stream.Recv()`.
- Returning EOF or errors when the stream closes.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `LocationGrpcHandler` struct embedding `locationPb.UnimplementedLocationServiceServer`.
2. `StreamDriverLocation(stream locationPb.LocationService_StreamDriverLocationServer) error`
   - Loop continuously:
     `msg, err := stream.Recv()`
     if err == io.EOF { return stream.SendAndClose(&locationPb.LocationResponse{Success: true}) }
   - Pass driver coordinates to `service.UpdateDriverLocation(msg.DriverId, msg.Latitude, msg.Longitude)`.

3. `FindNearbyDrivers(ctx context.Context, req *locationPb.NearbyDriversRequest) (*locationPb.NearbyDriversResponse, error)`
   - Call `service.FindNearbyDrivers(req.Latitude, req.Longitude, req.RadiusKm, req.Limit)`.
   - Convert service output to gRPC response struct and return.

4. `RemoveDriver(ctx context.Context, req *locationPb.DriverID) (*locationPb.LocationResponse, error)`
   - Call `service.RemoveDriver(req.DriverId)`.
================================================================================
*/

import (
	"context"
	"errors"
	"io"

	locationService "github.com/vikasrathod4299/microservice/internal/location/service"
	locationPb "github.com/vikasrathod4299/microservice/proto/location"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LocationGrpcHandler struct {
	locationPb.UnimplementedLocationServiceServer
	service *locationService.LocationService
}

func NewLocationGrpcHandler(service *locationService.LocationService) *LocationGrpcHandler {
	return &LocationGrpcHandler{
		service: service,
	}
}

func (h *LocationGrpcHandler) StreamDriverLocation(stream locationPb.LocationService_StreamDriverLocationServer) error {
	for {
		loc, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&locationPb.LocationResponse{
				Success: true,
				Message: "driver location updated",
			})
		}

		if err != nil {
			return err
		}

		err = h.service.UpdateLocation(stream.Context(), loc.GetDriverId(), loc.GetLatitude(), loc.GetLongitude())
		if err != nil {
			return toGRPCError(err)
		}
	}
}

func (h *LocationGrpcHandler) FindNearbyDrivers(ctx context.Context, req *locationPb.NearbyDriversRequest) (*locationPb.NearbyDriversResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	drivers, err := h.service.FindNearbyDrivers(ctx, req.GetLatitude(), req.GetLongitude(), req.RadiusKm, int(req.GetLimit()))
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoDrivers := make([]*locationPb.NearbyDriver, len(drivers))
	for i, driver := range drivers {
		protoDrivers[i] = &locationPb.NearbyDriver{
			DriverId:   driver.DriverID,
			Latitude:   driver.Latitude,
			Longitude:  driver.Longitude,
			DistanceKm: driver.DistanceKM,
		}
	}

	return &locationPb.NearbyDriversResponse{
		Drivers: protoDrivers,
	}, nil
}

func (h *LocationGrpcHandler) RemoveDriver(ctx context.Context, req *locationPb.DriverID) (*locationPb.LocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "req is required")
	}

	err := h.service.RemoveDriver(ctx, req.GetDriverId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &locationPb.LocationResponse{
		Success: true,
		Message: "Driver Removed",
	}, nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, locationService.ErrDriverIDRequired),
		errors.Is(err, locationService.ErrInvalidCoordinates),
		errors.Is(err, locationService.ErrInvalidLimit),
		errors.Is(err, locationService.ErrInvalidRadius):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()
	default:
		return status.Error(codes.Internal, "location service operation failed")
	}
}
