package util

import (
	"time"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
)

// ConvertPgTimestampToMillis converts pgtype.Timestamp to unix milliseconds
func ConvertPgTimestampToMillis(ts pgtype.Timestamp) int64 {
	if !ts.Valid {
		return 0
	}
	return ts.Time.UnixMilli()
}

// ConvertTimeToMillis converts time.Time to unix milliseconds
func ConvertTimeToMillis(t time.Time) int64 {
	return t.UnixMilli()
}

// IsSynced checks if data is synced (no updates)
func IsSynced(count int64) bool {
	return count == 0
}

// FormatSyncMessage returns appropriate sync message
func FormatSyncMessage(entityName string, isSynced bool, count int64) string {
	if isSynced {
		return "No new " + entityName + " updates"
	}
	return entityName + " synced successfully"
}

// ExtractMessageIDs extracts message IDs from sync messages
func ExtractMessageIDs(messages []repo.SyncMessagesRow) []int32 {
	ids := make([]int32, len(messages))
	for i, msg := range messages {
		ids[i] = msg.ID
	}
	return ids
}

// FormatFriendData formats friend data for response
func FormatFriendData(friend repo.SyncFriendsRow) map[string]interface{} {
	return map[string]interface{}{
		"id":          friend.ID,
		"user_id":     friend.UserID,
		"friend_id":   friend.FriendID,
		"status":      friend.Status,
		"alias_name":  friend.AliasName,
		"is_favorite": friend.IsFavorite,
		"created_at":  ConvertPgTimestampToMillis(friend.CreatedAt),
		"updated_at":  ConvertPgTimestampToMillis(friend.UpdatedAt),
		"friend": map[string]interface{}{
			"id":               friend.FriendUserID,
			"username":         friend.FriendUsername,
			"email":            friend.FriendEmail,
			"full_name":        friend.FriendFullName,
			"profile_pic":      friend.FriendProfilePic,
			"bio":              friend.FriendBio,
			"color_code":       friend.FriendColorCode,
			"background_color": friend.FriendBackgroundColor,
			"background_pic":   friend.FriendBackgroundPic,
			"status":           friend.FriendStatus,
			"custom_status":    friend.FriendCustomStatus,
			"is_bot":           friend.FriendIsBot,
			"is_verified":      friend.FriendIsVerified,
			"created_at":       ConvertPgTimestampToMillis(friend.FriendCreatedAt),
			"updated_at":       ConvertPgTimestampToMillis(friend.FriendUpdatedAt),
		},
	}
}

// FormatMessageData formats message data for response
func FormatMessageData(message repo.SyncMessagesRow) map[string]interface{} {
	return map[string]interface{}{
		"id":                  message.ID,
		"channel_id":          message.ChannelID,
		"sender_id":           message.SenderID,
		"content":             message.Content,
		"message_type":        message.MessageType,
		"reply_to_message_id": message.ReplyToMessageID,
		"is_edited":           message.IsEdited,
		"is_pinned":           message.IsPinned,
		"mention_everyone":    message.MentionEveryone,
		"created_at":          ConvertPgTimestampToMillis(message.CreatedAt),
		"updated_at":          ConvertPgTimestampToMillis(message.UpdatedAt),
		"edited_at":           ConvertPgTimestampToMillis(message.EditedAt),
		"sender": map[string]interface{}{
			"username":    message.SenderUsername,
			"profile_pic": message.SenderProfilePic,
		},
	}
}

// FormatServerData formats server data for response
func FormatServerData(server repo.Server) map[string]interface{} {
	return map[string]interface{}{
		"id":           server.ID,
		"name":         server.Name,
		"icon":         server.Icon,
		"banner":       server.Banner,
		"description":  server.Description,
		"owner_id":     server.OwnerID,
		"region":       server.Region,
		"member_count": server.MemberCount,
		"is_verified":  server.IsVerified,
		"vanity_url":   server.VanityUrl,
		"created_at":   ConvertPgTimestampToMillis(server.CreatedAt),
		"updated_at":   ConvertPgTimestampToMillis(server.UpdatedAt),
	}
}

// FormatChannelData formats channel data for response
func FormatChannelData(channel repo.Channel) map[string]interface{} {
	return map[string]interface{}{
		"id":             channel.ID,
		"server_id":      channel.ServerID,
		"category_id":    channel.CategoryID,
		"name":           channel.Name,
		"type":           channel.Type,
		"position":       channel.Position,
		"topic":          channel.Topic,
		"is_nsfw":        channel.IsNsfw,
		"slowmode_delay": channel.SlowmodeDelay,
		"user_limit":     channel.UserLimit,
		"bitrate":        channel.Bitrate,
		"is_private":     channel.IsPrivate,
		"created_at":     ConvertPgTimestampToMillis(channel.CreatedAt),
		"updated_at":     ConvertPgTimestampToMillis(channel.UpdatedAt),
	}
}

// FormatUserProfileData formats user profile data for response
func FormatUserProfileData(user repo.SyncUserProfileRow) map[string]interface{} {
	return map[string]interface{}{
		"id":               user.ID,
		"username":         user.Username,
		"email":            user.Email,
		"full_name":        user.FullName,
		"profile_pic":      user.ProfilePic,
		"bio":              user.Bio,
		"color_code":       user.ColorCode,
		"background_color": user.BackgroundColor,
		"background_pic":   user.BackgroundPic,
		"status":           user.Status,
		"custom_status":    user.CustomStatus,
		"is_bot":           user.IsBot,
		"is_verified":      user.IsVerified,
		"is_2fa_enabled":   user.Is2faEnabled,
		"created_at":       ConvertPgTimestampToMillis(user.CreatedAt),
		"updated_at":       ConvertPgTimestampToMillis(user.UpdatedAt),
	}
}

// CalculateTotalUpdates calculates total updates across all entities
func CalculateTotalUpdates(counts ...int64) int64 {
	var total int64
	for _, count := range counts {
		total += count
	}
	return total
}

// BuildSyncResponse builds a standard sync response
func BuildSyncResponse(data interface{}, serverTimestamp int64, isSynced bool, totalUpdates int, message string) map[string]interface{} {
	return map[string]interface{}{
		"data":             data,
		"server_timestamp": serverTimestamp,
		"is_synced":        isSynced,
		"total_updates":    totalUpdates,
		"message":          message,
	}
}
