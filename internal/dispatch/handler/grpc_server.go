package handler

/*
================================================================================
FILE: internal/dispatch/handler/grpc_server.go
================================================================================

PURPOSE:
Implements the `DispatchServiceServer` gRPC interface generated from `proto/dispatch/dispatch.proto`.
Receives RPC requests to create trips, fetch trip details, or update trip status.

LEARNING GO CONCEPTS:
- Mapping gRPC proto messages to internal service domains.
- Standard gRPC status codes (`codes.InvalidArgument`, `codes.NotFound`, `codes.Internal`).

WHAT YOU NEED TO IMPLEMENT HERE:
1. `DispatchGrpcHandler` struct embedding `dispatchPb.UnimplementedDispatchServiceServer`.

2. `RequestTrip(ctx context.Context, req *dispatchPb.TripRequest) (*dispatchPb.Trip, error)`
   - Call `service.CreateTrip(ctx, req.RiderId, req.PickupLat, req.PickupLng, req.DropoffLat, req.DropoffLng)`.
   - Convert resulting Trip struct into proto `*dispatchPb.Trip`.

3. `GetTrip(ctx context.Context, req *dispatchPb.TripID) (*dispatchPb.Trip, error)`
   - Call `service.GetTripByID(ctx, req.TripId)`.

4. `UpdateTripStatus(ctx context.Context, req *dispatchPb.UpdateTripStatusRequest) (*dispatchPb.Trip, error)`
   - Call `service.UpdateTripState(ctx, req.TripId, req.Status)`.
================================================================================
*/

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vikasrathod4299/microservice/internal/dispatch/service"
	dispatchPb "github.com/vikasrathod4299/microservice/proto/dispatch"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DispatchGrpcHandler struct {
	dispatchPb.UnimplementedDispatchServiceServer
	service *service.DispatchService
}

func NewDispatchGrpcHandler(svc *service.DispatchService) *DispatchGrpcHandler {
	return &DispatchGrpcHandler{
		service: svc,
	}
}

func (h *DispatchGrpcHandler) RequestTrip(ctx context.Context, req *dispatchPb.TripRequest) (*dispatchPb.Trip, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	trip, err := h.service.CreateTrip(ctx, req.GetPickupLat(), req.GetPickupLng(), req.GetDropoffLat(), req.GetDropoffLng(), req.GetRiderId())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return toProtoTrip(trip), nil
}

func (h *DispatchGrpcHandler) GetTripByID(ctx context.Context, req *dispatchPb.TripId) (*dispatchPb.Trip, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	trip, err := h.service.GetTripByID(ctx, req.GetTripId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoTrip(trip), nil
}

func (h *DispatchGrpcHandler) UpdateTripStatus(ctx context.Context, req *dispatchPb.UpdateTripStatusRequest) (*dispatchPb.Trip, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	nextStatus, err := fromProtoTripStatus(req.GetStatus())
	if err != nil {
		return nil, err
	}

	trip, err := h.service.UpdateTripStatus(ctx, req.GetTripId(), nextStatus, req.GetDriverId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoTrip(trip), nil
}

func toProtoTrip(t *service.Trip) *dispatchPb.Trip {
	if t == nil {
		return nil
	}

	return &dispatchPb.Trip{
		TripId:     t.ID,
		RiderId:    t.RiderID,
		DriverId:   t.DriverID,
		PickupLat:  t.PickupLat,
		PickupLng:  t.PickupLng,
		DropoffLat: t.DropoffLat,
		DropoffLng: t.DropoffLng,
		Status:     toProtoTripStatus(t.Status),
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
	}
}

func toProtoTripStatus(status service.TripStatus) dispatchPb.TripStatus {
	switch status {
	case service.StatusSearching:
		return dispatchPb.TripStatus_SEARCHING
	case service.StatusDriverAssigned:
		return dispatchPb.TripStatus_DRIVER_ASSIGNED
	case service.StatusDriverArrived:
		return dispatchPb.TripStatus_DRIVER_ARRIVED
	case service.StatusInTransit:
		return dispatchPb.TripStatus_IN_TRANSIT
	case service.StatusCompleted:
		return dispatchPb.TripStatus_COMPLETED
	case service.StatusCancelled:
		return dispatchPb.TripStatus_CANCELED
	default:
		return dispatchPb.TripStatus_TRIP_STATUS_UNSPECIFIED

	}
}

func fromProtoTripStatus(tripStatus dispatchPb.TripStatus) (service.TripStatus, error) {
	switch tripStatus {
	case dispatchPb.TripStatus_SEARCHING:
		return service.StatusSearching, nil
	case dispatchPb.TripStatus_DRIVER_ASSIGNED:
		return service.StatusDriverAssigned, nil
	case dispatchPb.TripStatus_DRIVER_ARRIVED:
		return service.StatusDriverArrived, nil
	case dispatchPb.TripStatus_COMPLETED:
		return service.StatusCompleted, nil
	case dispatchPb.TripStatus_CANCELED:
		return service.StatusCancelled, nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "invalid trip staus")
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, service.ErrDriverIDRequired),
		errors.Is(err, service.ErrTripIDRequired),
		errors.Is(err, service.ErrRiderIDRequired),
		errors.Is(err, service.ErrInvalidTripCoordinates):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrInvalidStateTransition):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, pgx.ErrNoRows):
		return status.Error(
			codes.NotFound,
			"trip not found",
		)

	case errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()

	default:
		return status.Error(
			codes.Internal,
			"dispatch service operation failed",
		)
	}
}
