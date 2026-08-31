package handler

/*
================================================================================
FILE: internal/driverstate/handler/rest.go
================================================================================

PURPOSE:
REST HTTP handlers for managing driver profiles and updating driver availability status.

LEARNING GO CONCEPTS:
- Standard REST API handler implementation.
- Reading path parameters (`chi.URLParam(r, "id")` or `r.PathValue("id")`).

WHAT YOU NEED TO IMPLEMENT HERE:
1. `CreateDriver(w http.ResponseWriter, r *http.Request)`
2. `GetDriver(w http.ResponseWriter, r *http.Request)`
3. `UpdateStatus(w http.ResponseWriter, r *http.Request)`
   - Parses status string (`OFFLINE`, `ONLINE`, `ON_TRIP`).
   - Calls `service.UpdateDriverStatus(...)`.
================================================================================
*/

import (
	"encoding/json"
	"net/http"
)

type DriverHandler struct {
	// service DriverService
}

func NewDriverHandler() *DriverHandler {
	return &DriverHandler{}
}

func (h *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "TODO: CreateDriver"})
}

func (h *DriverHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "TODO: GetDriver"})
}

func (h *DriverHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "TODO: UpdateStatus"})
}
