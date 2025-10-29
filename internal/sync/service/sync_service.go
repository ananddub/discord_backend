package service

import (
	"context"
	"fmt"

	"discord/internal/sync/repository"
)

type SyncService struct {
	repo *repository.SyncRepository
}

func NewSyncService(repo *repository.SyncRepository) *SyncService {
	return &SyncService{
		repo: repo,
	}
}

// SyncFriends syncs friends data for user
func (s *SyncService) SyncFriends(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) (interface{}, error) {
	// Get server timestamp
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	// Get updated friends
	friends, err := s.repo.SyncFriends(ctx, userID, lastUpdatedAt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to sync friends: %w", err)
	}

	// Get pending requests
	pendingRequests, err := s.repo.SyncPendingFriendRequests(ctx, userID, lastUpdatedAt, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to sync pending requests: %w", err)
	}

	// Get count
	count, err := s.repo.CountUpdatedFriends(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count updated friends: %w", err)
	}

	isSynced := count == 0 && len(pendingRequests) == 0
	message := "Friends synced successfully"
	if isSynced {
		message = "No new friend updates"
	}

	return map[string]interface{}{
		"friends":          friends,
		"pending_requests": pendingRequests,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    count,
		"message":          message,
	}, nil
}

// SyncMessages syncs messages for user
func (s *SyncService) SyncMessages(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	// Get all user messages
	messages, err := s.repo.SyncUserMessages(ctx, userID, lastUpdatedAt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to sync messages: %w", err)
	}

	// Get count
	count, err := s.repo.CountUpdatedMessages(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count updated messages: %w", err)
	}

	isSynced := count == 0
	hasMore := int64(len(messages)) >= int64(limit)
	message := "Messages synced successfully"
	if isSynced {
		message = "No new messages"
	}

	return map[string]interface{}{
		"messages":         messages,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    count,
		"has_more":         hasMore,
		"message":          message,
	}, nil
}

// SyncServers syncs servers for user
func (s *SyncService) SyncServers(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	servers, err := s.repo.SyncServers(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sync servers: %w", err)
	}

	count, err := s.repo.CountUpdatedServers(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count updated servers: %w", err)
	}

	isSynced := count == 0
	message := "Servers synced successfully"
	if isSynced {
		message = "No server updates"
	}

	return map[string]interface{}{
		"servers":          servers,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    count,
		"message":          message,
	}, nil
}

// SyncChannels syncs channels for user
func (s *SyncService) SyncChannels(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	channels, err := s.repo.SyncChannels(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sync channels: %w", err)
	}

	count, err := s.repo.CountUpdatedChannels(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count updated channels: %w", err)
	}

	isSynced := count == 0
	message := "Channels synced successfully"
	if isSynced {
		message = "No channel updates"
	}

	return map[string]interface{}{
		"channels":         channels,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    count,
		"message":          message,
	}, nil
}

// SyncUserProfile syncs user profile
func (s *SyncService) SyncUserProfile(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	user, err := s.repo.SyncUserProfile(ctx, userID, lastUpdatedAt)
	if err != nil {
		// User not updated
		lastUpdate, _ := s.repo.GetUserLastUpdate(ctx, userID)
		return map[string]interface{}{
			"user":                nil,
			"server_timestamp":    serverTimestamp,
			"is_synced":           true,
			"last_profile_update": lastUpdate.Time.UnixMilli(),
			"message":             "Profile is up to date",
		}, nil
	}

	return map[string]interface{}{
		"user":                user,
		"server_timestamp":    serverTimestamp,
		"is_synced":           false,
		"last_profile_update": user.UpdatedAt.Time.UnixMilli(),
		"message":             "Profile updated",
	}, nil
}

// SyncVoiceChannels syncs voice channels for user
func (s *SyncService) SyncVoiceChannels(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	voiceChannels, err := s.repo.SyncVoiceChannels(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sync voice channels: %w", err)
	}

	isSynced := len(voiceChannels) == 0
	message := "Voice channels synced successfully"
	if isSynced {
		message = "No voice channel updates"
	}

	return map[string]interface{}{
		"voice_channels":   voiceChannels,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    len(voiceChannels),
		"message":          message,
	}, nil
}

// SyncTextChannels syncs text channels for user
func (s *SyncService) SyncTextChannels(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	textChannels, err := s.repo.SyncTextChannels(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sync text channels: %w", err)
	}

	isSynced := len(textChannels) == 0
	message := "Text channels synced successfully"
	if isSynced {
		message = "No text channel updates"
	}

	return map[string]interface{}{
		"text_channels":    textChannels,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    len(textChannels),
		"message":          message,
	}, nil
}

// SyncDirectMessages syncs DM channels for user
func (s *SyncService) SyncDirectMessages(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	directMessages, err := s.repo.SyncDirectMessages(ctx, userID, lastUpdatedAt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to sync direct messages: %w", err)
	}

	isSynced := len(directMessages) == 0
	hasMore := int32(len(directMessages)) >= limit
	message := "Direct messages synced successfully"
	if isSynced {
		message = "No new direct messages"
	}

	return map[string]interface{}{
		"direct_messages":  directMessages,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    len(directMessages),
		"has_more":         hasMore,
		"message":          message,
	}, nil
}

// SyncPermissions syncs permissions for user
func (s *SyncService) SyncPermissions(ctx context.Context, userID int32, lastUpdatedAt int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	permissions, err := s.repo.SyncPermissions(ctx, userID, lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sync permissions: %w", err)
	}

	isSynced := len(permissions) == 0
	message := "Permissions synced successfully"
	if isSynced {
		message = "No permission updates"
	}

	return map[string]interface{}{
		"permissions":      permissions,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    len(permissions),
		"message":          message,
	}, nil
}

// SyncAll syncs all data for user
func (s *SyncService) SyncAll(ctx context.Context, userID int32, timestamps map[string]int64) (interface{}, error) {
	serverTimestamp, err := s.repo.GetServerTimestamp(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server timestamp: %w", err)
	}

	// Sync all entities in parallel would be ideal, but for simplicity doing sequential
	friends, _ := s.SyncFriends(ctx, userID, timestamps["friends"], 100, 0)
	messages, _ := s.SyncMessages(ctx, userID, timestamps["messages"], 100, 0)
	servers, _ := s.SyncServers(ctx, userID, timestamps["servers"])
	channels, _ := s.SyncChannels(ctx, userID, timestamps["channels"])
	userProfile, _ := s.SyncUserProfile(ctx, userID, timestamps["user_profile"])
	voiceChannels, _ := s.SyncVoiceChannels(ctx, userID, timestamps["voice_channels"])
	textChannels, _ := s.SyncTextChannels(ctx, userID, timestamps["text_channels"])
	directMessages, _ := s.SyncDirectMessages(ctx, userID, timestamps["direct_messages"], 50, 0)
	permissions, _ := s.SyncPermissions(ctx, userID, timestamps["permissions"])

	// Calculate if fully synced
	isFullySynced := true
	totalUpdates := 0

	if friendsData, ok := friends.(map[string]interface{}); ok {
		if !friendsData["is_synced"].(bool) {
			isFullySynced = false
		}
		totalUpdates += int(friendsData["total_updates"].(int64))
	}

	if messagesData, ok := messages.(map[string]interface{}); ok {
		if !messagesData["is_synced"].(bool) {
			isFullySynced = false
		}
		totalUpdates += int(messagesData["total_updates"].(int64))
	}

	return map[string]interface{}{
		"friends":          friends,
		"messages":         messages,
		"servers":          servers,
		"channels":         channels,
		"user_profile":     userProfile,
		"voice_channels":   voiceChannels,
		"text_channels":    textChannels,
		"direct_messages":  directMessages,
		"permissions":      permissions,
		"server_timestamp": serverTimestamp,
		"is_fully_synced":  isFullySynced,
		"total_updates":    totalUpdates,
		"message":          "Sync completed",
	}, nil
}
