package errors

import "fmt"

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// DatabaseError represents database operation errors
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database %s failed: %v", e.Operation, e.Err)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// AuthenticationError represents authentication failures
type AuthenticationError struct {
	Reason string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(reason string) *AuthenticationError {
	return &AuthenticationError{
		Reason: reason,
	}
}

// PermissionError represents permission/authorization errors
type PermissionError struct {
	Action   string
	Resource string
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("permission denied: cannot %s %s", e.Action, e.Resource)
}

// NewPermissionError creates a new permission error
func NewPermissionError(action, resource string) *PermissionError {
	return &PermissionError{
		Action:   action,
		Resource: resource,
	}
}
