package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	voiceRepo "discord/internal/voice/repository"
)

type VoiceService struct {
	voiceRepo *voiceRepo.VoiceRepository
}

func NewVoiceService(voiceRepo *voiceRepo.VoiceRepository) *VoiceService {
	return &VoiceService{
		voiceRepo: voiceRepo,
	}
}

// JoinVoiceChannel allows a user to join a voice channel
func (s *VoiceService) JoinVoiceChannel(ctx context.Context, userID, channelID int32, serverID *int32) (repo.VoiceState, error) {
	// Check if user is already in a voice channel
	existingState, err := s.voiceRepo.GetUserVoiceState(ctx, userID)
	if err == nil {
		// User is already in a voice channel, leave it first
		if existingState.ChannelID != channelID {
			err = s.voiceRepo.DeleteVoiceState(ctx, userID, existingState.ChannelID)
			if err != nil {
				return repo.VoiceState{}, commonErrors.ErrInternalServer
			}
		} else {
			// Already in this channel
			return existingState, nil
		}
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrInternalServer
	}

	// Create voice state with default values
	voiceState, err := s.voiceRepo.CreateVoiceState(ctx, userID, channelID, serverID, sessionID, false, false, false, false)
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrInternalServer
	}

	return voiceState, nil
}

// LeaveVoiceChannel allows a user to leave a voice channel
func (s *VoiceService) LeaveVoiceChannel(ctx context.Context, userID, channelID int32) error {
	// Verify user is in the channel
	_, err := s.voiceRepo.GetVoiceState(ctx, userID, channelID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Delete voice state
	err = s.voiceRepo.DeleteVoiceState(ctx, userID, channelID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

// GetVoiceState retrieves a user's voice state in a channel
func (s *VoiceService) GetVoiceState(ctx context.Context, userID, channelID int32) (repo.VoiceState, error) {
	state, err := s.voiceRepo.GetVoiceState(ctx, userID, channelID)
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrNotFound
	}
	return state, nil
}

// GetUserVoiceState retrieves a user's current voice state
func (s *VoiceService) GetUserVoiceState(ctx context.Context, userID int32) (repo.VoiceState, error) {
	state, err := s.voiceRepo.GetUserVoiceState(ctx, userID)
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrNotFound
	}
	return state, nil
}

// GetChannelVoiceStates retrieves all voice states in a channel
func (s *VoiceService) GetChannelVoiceStates(ctx context.Context, channelID int32) ([]repo.VoiceState, error) {
	states, err := s.voiceRepo.GetChannelVoiceStates(ctx, channelID)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}
	return states, nil
}

// UpdateVoiceState updates a user's voice state properties
func (s *VoiceService) UpdateVoiceState(ctx context.Context, userID, channelID int32, isMuted, isDeafened, selfMute, selfDeaf, selfVideo, selfStream *bool) (repo.VoiceState, error) {
	// Verify user is in the channel
	_, err := s.voiceRepo.GetVoiceState(ctx, userID, channelID)
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrNotFound
	}

	// Update voice state
	state, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, isMuted, isDeafened, selfMute, selfDeaf, selfVideo, selfStream)
	if err != nil {
		return repo.VoiceState{}, commonErrors.ErrInternalServer
	}

	return state, nil
}

// MuteUser server mutes a user (requires permission)
func (s *VoiceService) MuteUser(ctx context.Context, userID, channelID int32) error {
	muted := true
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, &muted, nil, nil, nil, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// UnmuteUser server unmutes a user (requires permission)
func (s *VoiceService) UnmuteUser(ctx context.Context, userID, channelID int32) error {
	muted := false
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, &muted, nil, nil, nil, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// DeafenUser server deafens a user (requires permission)
func (s *VoiceService) DeafenUser(ctx context.Context, userID, channelID int32) error {
	deafened := true
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, &deafened, nil, nil, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// UndeafenUser server undeafens a user (requires permission)
func (s *VoiceService) UndeafenUser(ctx context.Context, userID, channelID int32) error {
	deafened := false
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, &deafened, nil, nil, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// ToggleSelfMute toggles user's self mute
func (s *VoiceService) ToggleSelfMute(ctx context.Context, userID, channelID int32, mute bool) error {
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, nil, &mute, nil, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// ToggleSelfDeaf toggles user's self deafen
func (s *VoiceService) ToggleSelfDeaf(ctx context.Context, userID, channelID int32, deaf bool) error {
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, nil, nil, &deaf, nil, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// ToggleSelfVideo toggles user's video
func (s *VoiceService) ToggleSelfVideo(ctx context.Context, userID, channelID int32, video bool) error {
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, nil, nil, nil, &video, nil)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// ToggleSelfStream toggles user's screen share
func (s *VoiceService) ToggleSelfStream(ctx context.Context, userID, channelID int32, stream bool) error {
	_, err := s.voiceRepo.UpdateVoiceState(ctx, userID, channelID, nil, nil, nil, nil, nil, &stream)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// DisconnectUser disconnects a user from voice (requires permission)
func (s *VoiceService) DisconnectUser(ctx context.Context, userID int32) error {
	err := s.voiceRepo.DeleteUserVoiceStates(ctx, userID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}
	return nil
}

// generateSessionID generates a unique session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetChannelUserCount returns the number of users in a voice channel
func (s *VoiceService) GetChannelUserCount(ctx context.Context, channelID int32) (int, error) {
	states, err := s.voiceRepo.GetChannelVoiceStates(ctx, channelID)
	if err != nil {
		return 0, commonErrors.ErrInternalServer
	}
	return len(states), nil
}

// IsUserInVoiceChannel checks if a user is in a specific voice channel
func (s *VoiceService) IsUserInVoiceChannel(ctx context.Context, userID, channelID int32) bool {
	_, err := s.voiceRepo.GetVoiceState(ctx, userID, channelID)
	return err == nil
}

// MoveUser moves a user from one voice channel to another (requires permission)
func (s *VoiceService) MoveUser(ctx context.Context, userID, fromChannelID, toChannelID int32, serverID *int32) error {
	// Get current state
	currentState, err := s.voiceRepo.GetVoiceState(ctx, userID, fromChannelID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Delete old state
	err = s.voiceRepo.DeleteVoiceState(ctx, userID, fromChannelID)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	// Create new state with same properties
	selfMute := currentState.SelfMute.Valid && currentState.SelfMute.Bool
	selfDeaf := currentState.SelfDeaf.Valid && currentState.SelfDeaf.Bool
	isMuted := currentState.IsMuted.Valid && currentState.IsMuted.Bool
	isDeafened := currentState.IsDeafened.Valid && currentState.IsDeafened.Bool

	_, err = s.voiceRepo.CreateVoiceState(ctx, userID, toChannelID, serverID, currentState.SessionID, isMuted, isDeafened, selfMute, selfDeaf)
	if err != nil {
		return commonErrors.ErrInternalServer
	}

	return nil
}

// ValidateVoicePermissions checks if user has permission to join voice channel
func (s *VoiceService) ValidateVoicePermissions(ctx context.Context, userID, channelID int32) error {
	// TODO: Implement permission checks with server/channel permissions
	// For now, allow all users
	return nil
}

// GetVoiceChannelStats returns statistics for a voice channel
func (s *VoiceService) GetVoiceChannelStats(ctx context.Context, channelID int32) (map[string]interface{}, error) {
	states, err := s.voiceRepo.GetChannelVoiceStates(ctx, channelID)
	if err != nil {
		return nil, commonErrors.ErrInternalServer
	}

	stats := map[string]interface{}{
		"total_users":     len(states),
		"muted_users":     0,
		"deafened_users":  0,
		"video_users":     0,
		"streaming_users": 0,
	}

	for _, state := range states {
		if state.IsMuted.Valid && state.IsMuted.Bool {
			stats["muted_users"] = stats["muted_users"].(int) + 1
		}
		if state.IsDeafened.Valid && state.IsDeafened.Bool {
			stats["deafened_users"] = stats["deafened_users"].(int) + 1
		}
		if state.SelfVideo.Valid && state.SelfVideo.Bool {
			stats["video_users"] = stats["video_users"].(int) + 1
		}
		if state.SelfStream.Valid && state.SelfStream.Bool {
			stats["streaming_users"] = stats["streaming_users"].(int) + 1
		}
	}

	return stats, nil
}

// FormatVoiceStateInfo returns formatted voice state information
func (s *VoiceService) FormatVoiceStateInfo(state repo.VoiceState) string {
	status := fmt.Sprintf("User %d in Channel %d", state.UserID, state.ChannelID)

	if state.SelfMute.Valid && state.SelfMute.Bool {
		status += " [Self Muted]"
	}
	if state.SelfDeaf.Valid && state.SelfDeaf.Bool {
		status += " [Self Deafened]"
	}
	if state.IsMuted.Valid && state.IsMuted.Bool {
		status += " [Server Muted]"
	}
	if state.IsDeafened.Valid && state.IsDeafened.Bool {
		status += " [Server Deafened]"
	}
	if state.SelfVideo.Valid && state.SelfVideo.Bool {
		status += " [Video On]"
	}
	if state.SelfStream.Valid && state.SelfStream.Bool {
		status += " [Streaming]"
	}

	return status
}
