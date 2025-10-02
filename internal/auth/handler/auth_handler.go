package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth "discord/gen/proto/service/auth"
	"discord/internal/auth/service"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: service.NewAuthService()}
}

func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	response, err := h.authService.Register(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return response, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := h.authService.Login(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return user, nil
}

func (h *AuthHandler) ForgotPassword(ctx context.Context, req *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error) {
	return h.authService.ForgotPassword(ctx, req)
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	return h.authService.ResetPassword(ctx, req)
}
