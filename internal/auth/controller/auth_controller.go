package controller

import (
	"context"
	"discord/gen/proto/service/auth"
	authService "discord/internal/auth/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	auth.UnimplementedAuthServiceServer
	authService *authService.AuthService
}

func NewAuthController(authService *authService.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
func (c *AuthController) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email, and password are required")
	}

	err := c.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}

// Login handles user authentication
func (c *AuthController) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	accessToken, refreshToken, user, err := c.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Logout handles user logout
func (c *AuthController) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	err := c.authService.Logout(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LogoutResponse{
		Message: "Logged out successfully",
		Success: true,
	}, nil
}

// RefreshToken generates new access token
func (c *AuthController) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	accessToken, refreshToken, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Success:      true,
	}, nil
}

// VerifyEmail verifies user email
func (c *AuthController) VerifyEmail(ctx context.Context, req *auth.VerifyEmailRequest) (*auth.VerifyEmailResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "verification token is required")
	}

	err := c.authService.VerifyEmail(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.VerifyEmailResponse{
		Message: "Email verified successfully",
		Success: true,
	}, nil
}

// ForgotPassword initiates password reset
func (c *AuthController) ForgotPassword(ctx context.Context, req *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	err := c.authService.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.ForgotPasswordResponse{
		Message: "Password reset email sent",
	}, nil
}

// ResetPassword resets user password
func (c *AuthController) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	if req.Token == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "token and new password are required")
	}

	err := c.authService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.ResetPasswordResponse{
		Message: "Password reset successfully",
	}, nil
}

// ChangePassword changes user password
func (c *AuthController) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	if req.OldPassword == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old and new password are required")
	}

	err := c.authService.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.ChangePasswordResponse{
		Message: "Password changed successfully",
		Success: true,
	}, nil
}

// Enable2FA enables two-factor authentication
func (c *AuthController) Enable2FA(ctx context.Context, req *auth.Enable2FARequest) (*auth.Enable2FAResponse, error) {
	secret, qrCode, backupCodes, err := c.authService.Enable2FA(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Enable2FAResponse{
		Secret:      secret,
		QrCodeUrl:   qrCode,
		BackupCodes: backupCodes,
		Success:     true,
	}, nil
}

// Verify2FA verifies 2FA code
func (c *AuthController) Verify2FA(ctx context.Context, req *auth.Verify2FARequest) (*auth.Verify2FAResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "2FA code is required")
	}

	valid, err := c.authService.Verify2FA(ctx, req.UserId, req.Code)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.Verify2FAResponse{
		Valid:   valid,
		Message: "2FA verified",
	}, nil
}
