package util

import (
	"time"
)

// DMChannelType represents the type of DM channel
type DMChannelType int

const (
	DMChannelTypeDirect DMChannelType = iota // 1-on-1 DM
	DMChannelTypeGroup                       // Group DM
)

// MaxGroupDMParticipants is the maximum number of participants in a group DM
const MaxGroupDMParticipants = 10

// ValidateGroupDMSize checks if the number of participants is valid for a group DM
func ValidateGroupDMSize(participantCount int) bool {
	return participantCount >= 2 && participantCount <= MaxGroupDMParticipants
}

// GenerateGroupDMName generates a default name for a group DM based on participant names
func GenerateGroupDMName(participantNames []string) string {
	if len(participantNames) == 0 {
		return "Group DM"
	}

	if len(participantNames) <= 3 {
		name := ""
		for i, pName := range participantNames {
			if i > 0 {
				name += ", "
			}
			name += pName
		}
		return name
	}

	// If more than 3, show first 3 and add count
	name := participantNames[0] + ", " + participantNames[1] + ", " + participantNames[2]
	remaining := len(participantNames) - 3
	return name + " and " + string(rune(remaining)) + " others"
}

// IsRecentlyActive checks if a DM channel was active in the last N days
func IsRecentlyActive(lastMessageAt time.Time, days int) bool {
	if lastMessageAt.IsZero() {
		return false
	}
	threshold := time.Now().AddDate(0, 0, -days)
	return lastMessageAt.After(threshold)
}

// CalculateUnreadCount calculates unread message count
// This is a placeholder - actual implementation would query the database
func CalculateUnreadCount(lastReadMessageID, latestMessageID int32) int32 {
	if latestMessageID <= lastReadMessageID {
		return 0
	}
	return latestMessageID - lastReadMessageID
}

// DMChannelState represents the state of a DM channel
type DMChannelState int

const (
	DMChannelStateActive DMChannelState = iota
	DMChannelStateArchived
	DMChannelStateClosed
)

// CanSendMessage checks if messages can be sent in the DM channel state
func CanSendMessage(state DMChannelState) bool {
	return state == DMChannelStateActive
}

// DMPermissions represents DM-specific permissions
type DMPermissions struct {
	CanSendMessages       bool
	CanAttachFiles        bool
	CanAddReactions       bool
	CanManageChannel      bool // Only for group DM owner
	CanAddParticipants    bool // Only for group DM
	CanRemoveParticipants bool // Only for group DM owner
}

// GetDefaultDMPermissions returns default permissions for DM participants
func GetDefaultDMPermissions(isOwner bool) DMPermissions {
	return DMPermissions{
		CanSendMessages:       true,
		CanAttachFiles:        true,
		CanAddReactions:       true,
		CanManageChannel:      isOwner,
		CanAddParticipants:    true, // In group DMs, all members can add
		CanRemoveParticipants: isOwner,
	}
}

// ValidateDMChannelName validates group DM channel name
func ValidateDMChannelName(name string) bool {
	if name == "" {
		return false
	}
	// Name should be between 1 and 100 characters
	return len(name) >= 1 && len(name) <= 100
}

// FormatDMChannelName formats the channel name for display
func FormatDMChannelName(name string, isGroup bool, participantNames []string) string {
	if name != "" {
		return name
	}

	if !isGroup && len(participantNames) > 0 {
		// For 1-on-1 DMs, show the other person's name
		return participantNames[0]
	}

	// For group DMs without a name, generate one
	return GenerateGroupDMName(participantNames)
}

// DMNotificationType represents types of DM notifications
type DMNotificationType int

const (
	DMNotificationNewMessage DMNotificationType = iota
	DMNotificationMention
	DMNotificationReply
	DMNotificationGroupInvite
	DMNotificationGroupRemoval
)

// ShouldNotify determines if a notification should be sent
func ShouldNotify(notificationType DMNotificationType, isMuted bool, isFriend bool) bool {
	if isMuted {
		return false
	}

	// Only send notifications if users are friends or for important events
	switch notificationType {
	case DMNotificationGroupInvite, DMNotificationGroupRemoval:
		return true
	case DMNotificationMention, DMNotificationReply:
		return true
	case DMNotificationNewMessage:
		return isFriend
	default:
		return false
	}
}

// DMMessageStatus represents the status of a DM message
type DMMessageStatus int

const (
	DMMessageStatusSending DMMessageStatus = iota
	DMMessageStatusSent
	DMMessageStatusDelivered
	DMMessageStatusRead
	DMMessageStatusFailed
)

// GetMessageStatusText returns human-readable message status
func GetMessageStatusText(status DMMessageStatus) string {
	switch status {
	case DMMessageStatusSending:
		return "Sending..."
	case DMMessageStatusSent:
		return "Sent"
	case DMMessageStatusDelivered:
		return "Delivered"
	case DMMessageStatusRead:
		return "Read"
	case DMMessageStatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}
