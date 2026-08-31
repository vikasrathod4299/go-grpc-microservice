package service

/*
================================================================================
FILE: internal/dispatch/service/statemachine.go
================================================================================

PURPOSE:
Enforces valid State Transitions for a Trip lifecycle.
Prevents illegal status jumps (e.g. from SEARCHING directly to COMPLETED).

TRIP STATES:
1. SEARCHING
2. DRIVER_ASSIGNED
3. DRIVER_ARRIVED
4. IN_TRANSIT
5. COMPLETED
6. CANCELLED

LEARNING GO CONCEPTS:
- Using Go maps for state transition validation matrix (`map[State][]State`).
- Custom domain errors.

WHAT YOU NEED TO IMPLEMENT HERE:
1. Define custom `TripStatus` string or enum type.
2. `ValidateTransition(currentStatus, newStatus TripStatus) error`:
   - Checks if `newStatus` is allowed from `currentStatus`.
================================================================================
*/

import "errors"

type TripStatus string

const (
	StatusSearching      TripStatus = "SEARCHING"
	StatusDriverAssigned TripStatus = "DRIVER_ASSIGNED"
	StatusDriverArrived  TripStatus = "DRIVER_ARRIVED"
	StatusInTransit      TripStatus = "IN_TRANSIT"
	StatusCompleted      TripStatus = "COMPLETED"
	StatusCancelled      TripStatus = "CANCELLED"
)

var ErrInvalidStateTransition = errors.New("invalid trip state transition")

var validTransitions = map[TripStatus][]TripStatus{
	StatusSearching:      {StatusDriverAssigned, StatusCancelled},
	StatusDriverAssigned: {StatusDriverArrived, StatusCancelled},
	StatusDriverArrived:  {StatusInTransit, StatusCancelled},
	StatusInTransit:      {StatusCompleted, StatusCancelled},
	StatusCompleted:      {},
	StatusCancelled:      {},
}

func ValidateTransition(current, next TripStatus) error {
	allowed, exists := validTransitions[current]
	if !exists {
		return ErrInvalidStateTransition
	}
	for _, s := range allowed {
		if s == next {
			return nil
		}
	}
	return ErrInvalidStateTransition
}
