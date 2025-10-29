package util

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"fmt"
)

// FormatVoiceState converts database voice state to proto format
func FormatVoiceState(state repo.VoiceState) *schema.VoiceState {
	pbState := &schema.VoiceState{
		UserId:    state.UserID,
		ChannelId: state.ChannelID,
		SessionId: state.SessionID,
	}

	if state.ServerID.Valid {
		pbState.ServerId = state.ServerID.Int32
	}

	if state.IsMuted.Valid {
		pbState.IsMuted = state.IsMuted.Bool
	}
	if state.IsDeafened.Valid {
		pbState.IsDeafened = state.IsDeafened.Bool
	}
	if state.SelfMute.Valid {
		pbState.SelfMute = state.SelfMute.Bool
	}
	if state.SelfDeaf.Valid {
		pbState.SelfDeaf = state.SelfDeaf.Bool
	}
	if state.SelfVideo.Valid {
		pbState.SelfVideo = state.SelfVideo.Bool
	}
	if state.SelfStream.Valid {
		pbState.SelfStream = state.SelfStream.Bool
	}
	if state.Suppress.Valid {
		pbState.Suppress = state.Suppress.Bool
	}

	if state.JoinedAt.Valid {
		pbState.JoinedAt = state.JoinedAt.Time.Unix()
	}

	return pbState
}

// FormatVoiceStateList converts list of voice states to proto format
func FormatVoiceStateList(states []repo.VoiceState) []*schema.VoiceState {
	result := make([]*schema.VoiceState, len(states))
	for i, state := range states {
		result[i] = FormatVoiceState(state)
	}
	return result
}

// GetVoiceStateStatus returns human-readable status of voice state
func GetVoiceStateStatus(state repo.VoiceState) string {
	status := "Connected"

	modifiers := []string{}

	if state.SelfMute.Valid && state.SelfMute.Bool {
		modifiers = append(modifiers, "Self Muted")
	}
	if state.SelfDeaf.Valid && state.SelfDeaf.Bool {
		modifiers = append(modifiers, "Self Deafened")
	}
	if state.IsMuted.Valid && state.IsMuted.Bool {
		modifiers = append(modifiers, "Server Muted")
	}
	if state.IsDeafened.Valid && state.IsDeafened.Bool {
		modifiers = append(modifiers, "Server Deafened")
	}
	if state.SelfVideo.Valid && state.SelfVideo.Bool {
		modifiers = append(modifiers, "Video On")
	}
	if state.SelfStream.Valid && state.SelfStream.Bool {
		modifiers = append(modifiers, "Streaming")
	}
	if state.Suppress.Valid && state.Suppress.Bool {
		modifiers = append(modifiers, "Suppressed")
	}

	if len(modifiers) > 0 {
		status += " ("
		for i, mod := range modifiers {
			if i > 0 {
				status += ", "
			}
			status += mod
		}
		status += ")"
	}

	return status
}

// IsUserAudible checks if user can be heard in voice channel
func IsUserAudible(state repo.VoiceState) bool {
	// User is not audible if they are muted (self or server) or deafened (self or server)
	if state.IsMuted.Valid && state.IsMuted.Bool {
		return false
	}
	if state.IsDeafened.Valid && state.IsDeafened.Bool {
		return false
	}
	if state.SelfMute.Valid && state.SelfMute.Bool {
		return false
	}
	if state.SelfDeaf.Valid && state.SelfDeaf.Bool {
		return false
	}
	return true
}

// IsUserListening checks if user can hear others in voice channel
func IsUserListening(state repo.VoiceState) bool {
	// User is not listening if they are deafened (self or server)
	if state.IsDeafened.Valid && state.IsDeafened.Bool {
		return false
	}
	if state.SelfDeaf.Valid && state.SelfDeaf.Bool {
		return false
	}
	return true
}

// HasVideo checks if user has video enabled
func HasVideo(state repo.VoiceState) bool {
	return state.SelfVideo.Valid && state.SelfVideo.Bool
}

// IsStreaming checks if user is screen sharing
func IsStreaming(state repo.VoiceState) bool {
	return state.SelfStream.Valid && state.SelfStream.Bool
}

// GetVoiceChannelSummary returns summary information for a voice channel
func GetVoiceChannelSummary(states []repo.VoiceState) map[string]interface{} {
	summary := map[string]interface{}{
		"total_users":     len(states),
		"audible_users":   0,
		"listening_users": 0,
		"video_users":     0,
		"streaming_users": 0,
		"muted_users":     0,
		"deafened_users":  0,
	}

	for _, state := range states {
		if IsUserAudible(state) {
			summary["audible_users"] = summary["audible_users"].(int) + 1
		}
		if IsUserListening(state) {
			summary["listening_users"] = summary["listening_users"].(int) + 1
		}
		if HasVideo(state) {
			summary["video_users"] = summary["video_users"].(int) + 1
		}
		if IsStreaming(state) {
			summary["streaming_users"] = summary["streaming_users"].(int) + 1
		}
		if (state.IsMuted.Valid && state.IsMuted.Bool) || (state.SelfMute.Valid && state.SelfMute.Bool) {
			summary["muted_users"] = summary["muted_users"].(int) + 1
		}
		if (state.IsDeafened.Valid && state.IsDeafened.Bool) || (state.SelfDeaf.Valid && state.SelfDeaf.Bool) {
			summary["deafened_users"] = summary["deafened_users"].(int) + 1
		}
	}

	return summary
}

// ValidateSessionID checks if session ID format is valid
func ValidateSessionID(sessionID string) bool {
	// Session ID should be 32 characters hex string
	if len(sessionID) != 32 {
		return false
	}

	for _, c := range sessionID {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}

	return true
}

// FormatVoiceChannelInfo returns formatted voice channel information
func FormatVoiceChannelInfo(channelID int32, states []repo.VoiceState) string {
	summary := GetVoiceChannelSummary(states)
	return fmt.Sprintf(
		"Voice Channel %d: %d users (%d audible, %d listening, %d video, %d streaming)",
		channelID,
		summary["total_users"],
		summary["audible_users"],
		summary["listening_users"],
		summary["video_users"],
		summary["streaming_users"],
	)
}

// GetUserVoiceStatus returns user's voice status string
func GetUserVoiceStatus(state repo.VoiceState) string {
	if state.SelfDeaf.Valid && state.SelfDeaf.Bool {
		return "deafened"
	}
	if state.SelfMute.Valid && state.SelfMute.Bool {
		return "muted"
	}
	if state.IsDeafened.Valid && state.IsDeafened.Bool {
		return "server_deafened"
	}
	if state.IsMuted.Valid && state.IsMuted.Bool {
		return "server_muted"
	}
	return "speaking"
}

// CanUserSpeak checks if user has permission to speak (not suppressed or muted)
func CanUserSpeak(state repo.VoiceState) bool {
	if state.Suppress.Valid && state.Suppress.Bool {
		return false
	}
	return IsUserAudible(state)
}

// GetVoiceStateChanges compares two voice states and returns what changed
func GetVoiceStateChanges(oldState, newState repo.VoiceState) []string {
	changes := []string{}

	if oldState.IsMuted.Bool != newState.IsMuted.Bool {
		if newState.IsMuted.Bool {
			changes = append(changes, "server muted")
		} else {
			changes = append(changes, "server unmuted")
		}
	}

	if oldState.IsDeafened.Bool != newState.IsDeafened.Bool {
		if newState.IsDeafened.Bool {
			changes = append(changes, "server deafened")
		} else {
			changes = append(changes, "server undeafened")
		}
	}

	if oldState.SelfMute.Bool != newState.SelfMute.Bool {
		if newState.SelfMute.Bool {
			changes = append(changes, "self muted")
		} else {
			changes = append(changes, "self unmuted")
		}
	}

	if oldState.SelfDeaf.Bool != newState.SelfDeaf.Bool {
		if newState.SelfDeaf.Bool {
			changes = append(changes, "self deafened")
		} else {
			changes = append(changes, "self undeafened")
		}
	}

	if oldState.SelfVideo.Bool != newState.SelfVideo.Bool {
		if newState.SelfVideo.Bool {
			changes = append(changes, "video enabled")
		} else {
			changes = append(changes, "video disabled")
		}
	}

	if oldState.SelfStream.Bool != newState.SelfStream.Bool {
		if newState.SelfStream.Bool {
			changes = append(changes, "screen share started")
		} else {
			changes = append(changes, "screen share stopped")
		}
	}

	return changes
}
