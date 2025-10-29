package util

import (
	"encoding/json"
	"fmt"
)

// ToJSON converts object to JSON string
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON parses JSON string to object
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// PrettyJSON returns formatted JSON string
func PrettyJSON(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Contains checks if slice contains element
func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Remove removes element from slice
func Remove[T comparable](slice []T, element T) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if item != element {
			result = append(result, item)
		}
	}
	return result
}

// Unique returns unique elements from slice
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Map applies function to each element
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// Filter filters slice by predicate
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Chunk splits slice into chunks of specified size
func Chunk[T any](slice []T, size int) [][]T {
	var chunks [][]T

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// IntPtr returns pointer to int
func IntPtr(i int) *int {
	return &i
}

// Int32Ptr returns pointer to int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// StringPtr returns pointer to string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns pointer to bool
func BoolPtr(b bool) *bool {
	return &b
}

// Paginate returns paginated slice
func Paginate[T any](slice []T, page, pageSize int) []T {
	start := (page - 1) * pageSize
	if start > len(slice) {
		return []T{}
	}

	end := start + pageSize
	if end > len(slice) {
		end = len(slice)
	}

	return slice[start:end]
}

// Ternary returns trueVal if condition is true, else falseVal
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// Coalesce returns first non-zero value
func Coalesce[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

// Min returns minimum value
func Min[T interface {
	~int | ~int32 | ~int64 | ~float64
}](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns maximum value
func Max[T interface {
	~int | ~int32 | ~int64 | ~float64
}](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp clamps value between min and max
func Clamp[T interface {
	~int | ~int32 | ~int64 | ~float64
}](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// FormatError formats error with context
func FormatError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
