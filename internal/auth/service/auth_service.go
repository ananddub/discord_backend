package service

import (
	"context"
	"discord/gen/repo"

	auth "discord/gen/proto/service/auth"
	"discord/internal/auth/repository"
)

type AuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	ForgotPassword(ctx context.Context, req *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService() AuthService {
	authRepo, _ := repository.NewAuthRepository()
	return &authService{authRepo: *authRepo}
}

func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.authRepo.CreateUser(ctx, repo.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	return &auth.RegisterResponse{Message: ""}, nil
}

func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	return s.authRepo.Login(ctx, req)
}

func (s *authService) ForgotPassword(ctx context.Context, req *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error) {
	return s.authRepo.ForgotPassword(ctx, req)
}

func (s *authService) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	return s.authRepo.ResetPassword(ctx, req)
}
