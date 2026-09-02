package handler

import (
	"context"
	"errors"

	"github.com/vikasrathod4299/microservice/internal/auth/service"
	authPb "github.com/vikasrathod4299/microservice/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcHandler struct {
	service *service.AuthSerivce
	authPb.UnimplementedAuthServiceServer
}

func NewAuthGrpcService(service *service.AuthSerivce) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		service: service,
	}
}

func (h *AuthGrpcHandler) Register(ctx context.Context, req *authPb.RegisterRequest) (*authPb.AuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	role, err := roleFromProto(req.GetRole())
	if err != nil {
		return nil, err
	}

	result, err := h.service.Register(ctx, service.RegisterInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
		Role:     role,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return authResultToProto(result), nil
}

func (h *AuthGrpcHandler) Login(ctx context.Context, req *authPb.LoginRequest) (*authPb.AuthResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	result, err := h.service.Login(ctx, service.Login{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return authResultToProto(result), err
}

func (h *AuthGrpcHandler) GetUser(ctx context.Context, req *authPb.GetUserRequest) (*authPb.User, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	result, err := h.service.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return userToProto(result), nil
}

func authResultToProto(result *service.AuthResult) *authPb.AuthResponse {
	return &authPb.AuthResponse{
		User:        userToProto(result.User),
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   result.ExpiresAt.Unix(),
	}
}

func userToProto(user *service.User) *authPb.User {
	if user == nil {
		return nil
	}

	return &authPb.User{
		Id:        user.ID,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		Role:      roleToProto(user.Role),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}

func roleFromProto(role authPb.UserRole) (service.Role, error) {
	switch role {
	case authPb.UserRole_DRIVER:
		return service.RoleDriver, nil
	case authPb.UserRole_RIDER:
		return service.RoleRider, nil
	default:
		return "", status.Error(codes.InvalidArgument, "invalid user role")
	}
}

func roleToProto(
	role service.Role,
) authPb.UserRole {
	switch role {
	case service.RoleDriver:
		return authPb.UserRole_DRIVER
	case service.RoleRider:
		return authPb.UserRole_RIDER
	default:
		return authPb.UserRole_USER_ROLE_UNSPECIFIED
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, service.ErrNameRequired),
		errors.Is(err, service.ErrPhoneRequired),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidPassword),
		errors.Is(err, service.ErrInvalidRole),
		errors.Is(err, service.ErrUserIDRequired):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, service.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()

	default:
		return status.Error(
			codes.Internal,
			"Auth service operation failed",
		)
	}
}
