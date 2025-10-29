package util

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"fmt"
	"strings"
)

// FormatUserInfo converts database user to proto format with all fields
func FormatUserInfo(user repo.User) *schema.User {
	pbUser := &schema.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	// All optional string fields
	if user.FullName.Valid {
		pbUser.FullName = user.FullName.String
	}
	if user.ProfilePic.Valid {
		pbUser.ProfilePic = user.ProfilePic.String
	}
	if user.Bio.Valid {
		pbUser.Bio = user.Bio.String
	}
	if user.ColorCode.Valid {
		pbUser.ColorCode = user.ColorCode.String
	}
	if user.BackgroundColor.Valid {
		pbUser.BackgroundColor = user.BackgroundColor.String
	}
	if user.BackgroundPic.Valid {
		pbUser.BackgroundPic = user.BackgroundPic.String
	}
	if user.CustomStatus.Valid {
		pbUser.CustomStatus = user.CustomStatus.String
	}

	// Boolean fields
	if user.IsBot.Valid {
		pbUser.IsBot = user.IsBot.Bool
	}
	if user.IsVerified.Valid {
		pbUser.IsVerified = user.IsVerified.Bool
	}

	// Timestamps
	if user.CreatedAt.Valid {
		pbUser.CreatedAt = user.CreatedAt.Time.Unix()
	}
	if user.UpdatedAt.Valid {
		pbUser.UpdatedAt = user.UpdatedAt.Time.Unix()
	}

	return pbUser
}

// FormatUserPresence converts database presence to readable format
func FormatUserPresence(presence repo.UserPresence) map[string]interface{} {
	result := make(map[string]interface{})

	result["user_id"] = presence.UserID

	if presence.Status.Valid {
		result["status"] = presence.Status.String
	}
	if presence.CustomStatus.Valid {
		result["custom_status"] = presence.CustomStatus.String
	}
	if presence.CustomStatusEmoji.Valid {
		result["custom_status_emoji"] = presence.CustomStatusEmoji.String
	}
	if presence.Activity.Valid {
		result["activity"] = presence.Activity.String
	}
	if presence.LastSeen.Valid {
		result["last_seen"] = presence.LastSeen.Time.Unix()
	}

	return result
}

// ValidateUsername checks if username meets requirements
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 32 {
		return fmt.Errorf("username must be at most 32 characters")
	}

	// Only alphanumeric, underscore, and hyphen allowed
	for _, char := range username {
		if !isValidUsernameChar(char) {
			return fmt.Errorf("username contains invalid characters")
		}
	}

	return nil
}

func isValidUsernameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_' || char == '-'
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) error {
	if len(email) < 5 {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}

	return nil
}

// ValidateStatus checks if status is valid
func ValidateStatus(status string) bool {
	validStatuses := map[string]bool{
		"online":    true,
		"idle":      true,
		"dnd":       true,
		"offline":   true,
		"invisible": true,
	}
	return validStatuses[status]
}

// FormatUserList converts list of database users to proto format
func FormatUserList(users []repo.User) []*schema.User {
	result := make([]*schema.User, len(users))
	for i, user := range users {
		result[i] = FormatUserInfo(user)
	}
	return result
}

// SanitizeUserInput removes potentially harmful characters from input
func SanitizeUserInput(input string) string {
	// Remove null bytes and control characters
	result := strings.Map(func(r rune) rune {
		if r == 0 || (r < 32 && r != '\n' && r != '\r' && r != '\t') {
			return -1
		}
		return r
	}, input)

	return strings.TrimSpace(result)
}

// GenerateUserDisplayName creates a display name from user data
func GenerateUserDisplayName(user repo.User) string {
	if user.FullName.Valid && user.FullName.String != "" {
		return user.FullName.String
	}
	return user.Username
}

// IsUserOnline checks if user is online based on presence
func IsUserOnline(presence repo.UserPresence) bool {
	if !presence.Status.Valid {
		return false
	}

	status := presence.Status.String
	return status == "online" || status == "idle" || status == "dnd"
}

// FormatUserTag creates a user tag (username + discriminator)
func FormatUserTag(username string, discriminator int) string {
	return fmt.Sprintf("%s#%04d", username, discriminator)
}

// ValidateBio checks if bio meets requirements
func ValidateBio(bio string) error {
	if len(bio) > 190 {
		return fmt.Errorf("bio must be at most 190 characters")
	}
	return nil
}

// ValidateCustomStatus checks if custom status meets requirements
func ValidateCustomStatus(status string) error {
	if len(status) > 128 {
		return fmt.Errorf("custom status must be at most 128 characters")
	}
	return nil
}

// GetUserAge calculates user account age in days
func GetUserAge(user repo.User) int {
	if !user.CreatedAt.Valid {
		return 0
	}
	return int(user.CreatedAt.Time.Unix() / 86400)
}
