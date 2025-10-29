package util

import (
	"fmt"
	"time"
)

// TimeNow returns current timestamp
func TimeNow() time.Time {
	return time.Now()
}

// TimeToUnix converts time to unix timestamp
func TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// UnixToTime converts unix timestamp to time
func UnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// FormatTime formats time in human readable format
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimeAgo returns "time ago" format
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case duration < 365*24*time.Hour:
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(duration.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// IsExpired checks if timestamp is expired
func IsExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// AddDuration adds duration to current time
func AddDuration(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
