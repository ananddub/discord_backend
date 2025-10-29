package service

import (
	"context"
	"time"

	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	userRepo "discord/internal/user/repository"
)

type UserService struct {
	userRepo *userRepo.UserRepository
}

func NewUserService(userRepo *userRepo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// User Management
func (s *UserService) GetUser(ctx context.Context, userID int32) (repo.User, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return repo.User{}, commonErrors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return repo.User{}, commonErrors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (repo.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return repo.User{}, commonErrors.ErrNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int32, fullName, profilePic, bio, colorCode, backgroundColor, backgroundPic *string) (repo.User, error) {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return repo.User{}, commonErrors.ErrNotFound
	}

	// Update user
	user, err := s.userRepo.UpdateUser(ctx, userID, fullName, profilePic, bio, colorCode, backgroundColor, backgroundPic)
	if err != nil {
		return repo.User{}, commonErrors.ErrInternalServer
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int32) error {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Delete user
	err = s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

func (s *UserService) SearchUsers(ctx context.Context, query string, limit, offset int32) ([]repo.User, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.SearchUsers(ctx, query, limit, offset)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}

	return users, nil
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]repo.User, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}

	return users, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, userID int32, hashedPassword string) error {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Update password
	err = s.userRepo.UpdateUserPassword(ctx, userID, hashedPassword)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

// User Presence & Status
func (s *UserService) GetUserPresence(ctx context.Context, userID int32) (repo.UserPresence, error) {
	presence, err := s.userRepo.GetUserPresence(ctx, userID)
	if err != nil {
		return repo.UserPresence{}, commonErrors.ErrNotFound
	}
	return presence, nil
}

func (s *UserService) UpdateStatus(ctx context.Context, userID int32, status string, customStatus, activity *string) error {
	// Validate status
	validStatuses := map[string]bool{
		"online":    true,
		"idle":      true,
		"dnd":       true,
		"offline":   true,
		"invisible": true,
	}

	if !validStatuses[status] {
		return commonErrors.ErrInvalidInput
	}

	// Update presence
	_, err := s.userRepo.UpsertUserPresence(ctx, userID, &status, customStatus, nil, activity)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	// Also update user table status
	err = s.userRepo.UpdateUserStatus(ctx, userID, status)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

func (s *UserService) SetCustomStatus(ctx context.Context, userID int32, customStatus, emoji *string, expiresAt *time.Time) error {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Validate custom status length
	if customStatus != nil && len(*customStatus) > 128 {
		return commonErrors.ErrInvalidInput
	}

	// Set custom status
	err = s.userRepo.SetCustomStatus(ctx, userID, customStatus, emoji, expiresAt)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

func (s *UserService) GetMultipleUserPresences(ctx context.Context, userIDs []int32) ([]repo.UserPresence, error) {
	if len(userIDs) == 0 {
		return []repo.UserPresence{}, nil
	}

	presences, err := s.userRepo.GetMultipleUserPresences(ctx, userIDs)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}

	return presences, nil
}

// User Blocking
func (s *UserService) BlockUser(ctx context.Context, userID, blockedUserID int32) error {
	// Can't block yourself
	if userID == blockedUserID {
		return commonErrors.ErrInvalidInput
	}

	// Verify blocked user exists
	_, err := s.userRepo.GetUser(ctx, blockedUserID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Block user
	err = s.userRepo.BlockUser(ctx, userID, blockedUserID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

func (s *UserService) UnblockUser(ctx context.Context, userID, blockedUserID int32) error {
	// Can't unblock yourself
	if userID == blockedUserID {
		return commonErrors.ErrInvalidInput
	}

	// Unblock user
	err := s.userRepo.UnblockUser(ctx, userID, blockedUserID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

func (s *UserService) GetBlockedUsers(ctx context.Context, userID int32) ([]int32, error) {
	friends, err := s.userRepo.GetBlockedUsers(ctx, userID)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}

	// Extract blocked user IDs
	blockedIDs := make([]int32, len(friends))
	for i, friend := range friends {
		blockedIDs[i] = friend.FriendID
	}

	return blockedIDs, nil
}

// User Settings (simple in-memory for now, can be stored in DB later)
type UserSettings struct {
	ShowCurrentActivity bool
	AllowDMs            bool
	EnableNotifications bool
	Theme               string
	Language            string
}

var userSettingsCache = make(map[int32]*UserSettings)

func (s *UserService) GetUserSettings(ctx context.Context, userID int32) (*UserSettings, error) {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, commonErrors.ErrNotFound
	}

	// Get from cache or return defaults
	if settings, exists := userSettingsCache[userID]; exists {
		return settings, nil
	}

	// Default settings
	return &UserSettings{
		ShowCurrentActivity: true,
		AllowDMs:            true,
		EnableNotifications: true,
		Theme:               "dark",
		Language:            "en",
	}, nil
}

func (s *UserService) UpdateUserSettings(ctx context.Context, userID int32, showCurrentActivity, allowDMs, enableNotifications *bool, theme, language *string) error {
	// Verify user exists
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Get current settings or defaults
	settings, _ := s.GetUserSettings(ctx, userID)

	// Update settings
	if showCurrentActivity != nil {
		settings.ShowCurrentActivity = *showCurrentActivity
	}
	if allowDMs != nil {
		settings.AllowDMs = *allowDMs
	}
	if enableNotifications != nil {
		settings.EnableNotifications = *enableNotifications
	}
	if theme != nil {
		if *theme != "light" && *theme != "dark" {
			return commonErrors.ErrInvalidInput
		}
		settings.Theme = *theme
	}
	if language != nil {
		settings.Language = *language
	}

	// Store in cache
	userSettingsCache[userID] = settings

	return nil
}
