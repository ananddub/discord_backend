package util

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
)

// ConvertMessageToProto converts a repo.Message to proto.Message
func ConvertMessageToProto(message repo.Message) *schema.Message {
	pbMessage := &schema.Message{
		Id:        message.ID,
		ChannelId: message.ChannelID,
		SenderId:  message.SenderID,
		Content:   message.Content,
		IsEdited:  message.IsEdited.Bool,
		IsPinned:  message.IsPinned.Bool,
		CreatedAt: message.CreatedAt.Time.Unix(),
	}

	if message.ReplyToMessageID.Valid {
		pbMessage.ReplyToMessageId = message.ReplyToMessageID.Int32
	}

	if message.MentionEveryone.Valid {
		pbMessage.MentionEveryone = message.MentionEveryone.Bool
	}

	if message.EditedAt.Valid {
		pbMessage.EditedAt = message.EditedAt.Time.Unix()
	}

	if message.UpdatedAt.Valid {
		pbMessage.UpdatedAt = message.UpdatedAt.Time.Unix()
	}

	return pbMessage
}

// ConvertReactionToProto converts a repo.MessageReaction to proto format
func ConvertReactionToProto(reaction repo.MessageReaction) map[string]interface{} {
	result := map[string]interface{}{
		"id":         reaction.ID,
		"message_id": reaction.MessageID,
		"user_id":    reaction.UserID,
		"emoji":      reaction.Emoji,
		"created_at": reaction.CreatedAt.Time.Unix(),
	}

	if reaction.EmojiID.Valid {
		result["emoji_id"] = reaction.EmojiID.String
	}

	return result
}

// ConvertAttachmentToProto converts a repo.MessageAttachment to proto format
func ConvertAttachmentToProto(attachment repo.MessageAttachment) map[string]interface{} {
	result := map[string]interface{}{
		"id":         attachment.ID,
		"message_id": attachment.MessageID,
		"file_url":   attachment.FileUrl,
		"file_name":  attachment.FileName,
		"file_type":  attachment.FileType,
		"file_size":  attachment.FileSize,
		"created_at": attachment.CreatedAt.Time.Unix(),
	}

	if attachment.Width.Valid {
		result["width"] = attachment.Width.Int32
	}

	if attachment.Height.Valid {
		result["height"] = attachment.Height.Int32
	}

	return result
}

// ValidateMessageContent validates message content
func ValidateMessageContent(content string) bool {
	// Check if content is not empty
	if content == "" {
		return false
	}

	// Check max length (Discord limit is 2000 characters)
	if len(content) > 2000 {
		return false
	}

	return true
}

// ValidateMessageType validates message type
func ValidateMessageType(messageType string) bool {
	validTypes := []string{"TEXT", "IMAGE", "VIDEO", "FILE", "AUDIO", "SYSTEM"}
	for _, t := range validTypes {
		if t == messageType {
			return true
		}
	}
	return false
}

// IsSystemMessage checks if a message is a system message
func IsSystemMessage(messageType string) bool {
	return messageType == "SYSTEM"
}

// CanEditMessage checks if a user can edit a message
func CanEditMessage(message repo.Message, userID int32) bool {
	// User can only edit their own messages
	return message.SenderID == userID
}

// CanDeleteMessage checks if a user can delete a message
func CanDeleteMessage(message repo.Message, userID int32, hasModPermission bool) bool {
	// User can delete their own messages or if they have mod permissions
	return message.SenderID == userID || hasModPermission
}

// CanPinMessage checks if a user can pin a message
func CanPinMessage(hasManageMessagesPermission bool) bool {
	return hasManageMessagesPermission
}

// GroupReactionsByEmoji groups reactions by emoji
func GroupReactionsByEmoji(reactions []repo.MessageReaction) map[string][]int32 {
	grouped := make(map[string][]int32)

	for _, reaction := range reactions {
		if _, exists := grouped[reaction.Emoji]; !exists {
			grouped[reaction.Emoji] = []int32{}
		}
		grouped[reaction.Emoji] = append(grouped[reaction.Emoji], reaction.UserID)
	}

	return grouped
}

// FilterMessagesByUser filters messages by user ID
func FilterMessagesByUser(messages []repo.Message, userID int32) []repo.Message {
	var filtered []repo.Message
	for _, msg := range messages {
		if msg.SenderID == userID {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}

// GetPinnedMessageIDs extracts IDs from pinned messages
func GetPinnedMessageIDs(messages []repo.Message) []int32 {
	var ids []int32
	for _, msg := range messages {
		if msg.IsPinned.Bool {
			ids = append(ids, msg.ID)
		}
	}
	return ids
}

// HasMention checks if a message mentions a specific user
func HasMention(content string, username string) bool {
	// Simple check - in production, use regex or proper parsing
	mention := "@" + username
	return contains(content, mention)
}

// HasEveryone checks if a message mentions everyone
func HasEveryone(message repo.Message) bool {
	return message.MentionEveryone.Bool
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
