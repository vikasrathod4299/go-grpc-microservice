package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	authPb "github.com/vikasrathod4299/microservice/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type AuthRESTHandler struct {
	client authPb.AuthServiceClient
}

func NewAuthRESTHandler(authClient authPb.AuthServiceClient) *AuthRESTHandler {
	return &AuthRESTHandler{
		client: authClient,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *AuthRESTHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAuthError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	role, ok := parseUserRole(request.Role)
	if !ok {
		writeAuthError(w, http.StatusBadRequest, "role must be rider or driver")
	}

	response, err := h.client.Register(
		r.Context(),
		&authPb.RegisterRequest{
			Name:     request.Name,
			Phone:    request.Phone,
			Email:    request.Email,
			Password: request.Password,
			Role:     role,
		},
	)
	if err != nil {
		writeAuthGRPCError(w, err)
		return
	}

	writeAuthProtoJSON(w, http.StatusCreated, response)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthRESTHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAuthError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.client.Login(r.Context(), &authPb.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeAuthGRPCError(w, err)
		return
	}

	writeAuthProtoJSON(w, http.StatusOK, response)
}

func (h *AuthRESTHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	response, err := h.client.GetUser(r.Context(), &authPb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		writeAuthGRPCError(w, err)
		return
	}

	writeAuthProtoJSON(w, http.StatusOK, response)
}

func writeAuthProtoJSON(w http.ResponseWriter, statusCode int, message proto.Message) {
	body, err := protojson.MarshalOptions{
		UseProtoNames: true,
	}.Marshal(message)
	if err != nil {
		writeAuthError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(append(body, '\n'))
	return
}

func writeAuthGRPCError(w http.ResponseWriter, err error) {
	grpcStatus := status.Convert(err)
	writeAuthError(w, authHTTPStatus(grpcStatus.Code()), grpcStatus.Message())
}

func parseUserRole(
	value string,
) (authPb.UserRole, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "rider":
		return authPb.UserRole_RIDER, true

	case "driver":
		return authPb.UserRole_DRIVER, true

	default:
		return authPb.UserRole_USER_ROLE_UNSPECIFIED, false
	}
}

func authHTTPStatus(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest

	case codes.Unauthenticated:
		return http.StatusUnauthorized

	case codes.PermissionDenied:
		return http.StatusForbidden

	case codes.NotFound:
		return http.StatusNotFound

	case codes.AlreadyExists:
		return http.StatusConflict

	case codes.ResourceExhausted:
		return http.StatusTooManyRequests

	case codes.Unavailable:
		return http.StatusServiceUnavailable

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout

	default:
		return http.StatusInternalServerError
	}
}

func writeAuthError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
