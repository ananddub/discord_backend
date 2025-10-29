package service

import (
	"context"
	"errors"
	"time"

	"discord/gen/proto/schema"
	authRepo "discord/internal/auth/repository"
	"discord/internal/auth/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *authRepo.AuthRepository
}

func NewAuthService(authRepo *authRepo.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	// Validate input
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Check if user already exists
	existingUser, _ := s.authRepo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	existingUser, _ = s.authRepo.GetUserByUsername(ctx, username)
	if existingUser != nil {
		return errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	return s.authRepo.CreateUser(ctx, username, email, string(hashedPassword))
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, *schema.User, error) {
	// Get user by username or email
	user, err := s.authRepo.GetUserByUsername(ctx, username)
	if err != nil {
		user, err = s.authRepo.GetUserByEmail(ctx, username)
		if err != nil {
			return "", "", nil, errors.New("invalid credentials")
		}
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := util.GenerateJWT(user.ID, 24*time.Hour)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := util.GenerateJWT(user.ID, 30*24*time.Hour)
	if err != nil {
		return "", "", nil, err
	}

	// Update user status to online
	s.authRepo.UpdateUserStatus(ctx, user.ID, "online")

	// Convert to proto user
	protoUser := &schema.User{
		Id:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		FullName:        user.FullName.String,
		ProfilePic:      user.ProfilePic.String,
		Bio:             user.Bio.String,
		ColorCode:       user.ColorCode.String,
		BackgroundColor: user.BackgroundColor.String,
		BackgroundPic:   user.BackgroundPic.String,
		Status:          user.Status,
		CustomStatus:    user.CustomStatus.String,
		IsBot:           user.IsBot.Bool,
		IsVerified:      user.IsVerified.Bool,
		IsDeleted:       user.IsDeleted.Bool,
		CreatedAt:       user.CreatedAt.Time.Unix(),
		UpdatedAt:       user.UpdatedAt.Time.Unix(),
	}

	return accessToken, refreshToken, protoUser, nil
}

// Logout invalidates user tokens
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Extract user ID from token
	userID, err := util.ValidateJWT(accessToken)
	if err != nil {
		return err
	}

	// Update user status to offline
	return s.authRepo.UpdateUserStatus(ctx, userID, "offline")
}

// RefreshToken generates new access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate refresh token
	userID, err := util.ValidateJWT(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Generate new tokens
	newAccessToken, err := util.GenerateJWT(userID, 24*time.Hour)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := util.GenerateJWT(userID, 30*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// VerifyEmail verifies user email address
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// TODO: Implement email verification logic
	// 1. Decode token
	// 2. Extract user ID
	// 3. Update user verified status
	return nil
}

// ForgotPassword sends password reset email
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	resetToken, err := util.GenerateJWT(user.ID, 1*time.Hour)
	if err != nil {
		return err
	}

	// TODO: Send email with reset token
	_ = resetToken

	return nil
}

// ResetPassword resets user password using token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate token
	userID, err := util.ValidateJWT(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.authRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID int32, oldPassword, newPassword string) error {
	// Get user
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	return s.authRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

// Enable2FA enables two-factor authentication
func (s *AuthService) Enable2FA(ctx context.Context, userID int32) (string, string, []string, error) {
	// Generate 2FA secret
	secret := util.Generate2FASecret()

	// Generate QR code
	qrCode := util.Generate2FAQRCode(secret, "Discord")

	// Generate backup codes
	backupCodes := util.GenerateBackupCodes(10)

	// TODO: Store secret and backup codes in database

	return secret, qrCode, backupCodes, nil
}

// Verify2FA verifies 2FA code
func (s *AuthService) Verify2FA(ctx context.Context, userID int32, code string) (bool, error) {
	// TODO: Get user's 2FA secret from database
	// Verify code against secret
	valid := util.Verify2FACode("secret", code)

	if valid {
		// Enable 2FA for user
		s.authRepo.Enable2FA(ctx, userID, true)
	}

	return valid, nil
}
