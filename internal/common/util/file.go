package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

// AllowedImageTypes defines allowed image file types
var AllowedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// AllowedVideoTypes defines allowed video file types
var AllowedVideoTypes = []string{".mp4", ".webm", ".mov"}

// AllowedAudioTypes defines allowed audio file types
var AllowedAudioTypes = []string{".mp3", ".wav", ".ogg", ".m4a"}

// AllowedDocumentTypes defines allowed document file types
var AllowedDocumentTypes = []string{".pdf", ".doc", ".docx", ".txt", ".md"}

// MaxFileSize defines maximum file size (50MB)
const MaxFileSize = 50 * 1024 * 1024

// IsAllowedFileType checks if file type is allowed
func IsAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	allAllowed := append(AllowedImageTypes, AllowedVideoTypes...)
	allAllowed = append(allAllowed, AllowedAudioTypes...)
	allAllowed = append(allAllowed, AllowedDocumentTypes...)

	for _, allowed := range allAllowed {
		if ext == allowed {
			return true
		}
	}

	return false
}

// IsImageFile checks if file is an image
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range AllowedImageTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

// IsVideoFile checks if file is a video
func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range AllowedVideoTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

// IsAudioFile checks if file is an audio
func IsAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range AllowedAudioTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

// GetFileType returns the file type category
func GetFileType(filename string) string {
	switch {
	case IsImageFile(filename):
		return "image"
	case IsVideoFile(filename):
		return "video"
	case IsAudioFile(filename):
		return "audio"
	default:
		return "document"
	}
}

// ValidateFileSize checks if file size is within limits
func ValidateFileSize(size int64) error {
	if size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %s", FormatFileSize(MaxFileSize))
	}
	return nil
}

// SanitizeFilename removes potentially harmful characters from filename
func SanitizeFilename(filename string) string {
	// Remove path separators
	filename = filepath.Base(filename)

	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Remove special characters except dots, dashes, and underscores
	var safe strings.Builder
	for _, char := range filename {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' || char == '-' || char == '_' {
			safe.WriteRune(char)
		}
	}

	return safe.String()
}

// GenerateUniqueFilename generates a unique filename
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	base := strings.TrimSuffix(originalFilename, ext)
	base = SanitizeFilename(base)

	timestamp := TimeNow().Unix()
	randomID := GenerateID(8)

	return fmt.Sprintf("%s_%d_%s%s", base, timestamp, randomID, ext)
}
