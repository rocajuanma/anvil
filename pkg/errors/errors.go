package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	// ErrorTypeGeneral represents general errors
	ErrorTypeGeneral ErrorType = iota
	// ErrorTypePlatform represents platform-specific errors
	ErrorTypePlatform
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation
	// ErrorTypeConfiguration represents configuration errors
	ErrorTypeConfiguration
	// ErrorTypeInstallation represents installation errors
	ErrorTypeInstallation
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork
	// ErrorTypeFileSystem represents file system errors
	ErrorTypeFileSystem
)

// String returns a string representation of the error type
func (et ErrorType) String() string {
	switch et {
	case ErrorTypePlatform:
		return "platform"
	case ErrorTypeValidation:
		return "validation"
	case ErrorTypeConfiguration:
		return "configuration"
	case ErrorTypeInstallation:
		return "installation"
	case ErrorTypeNetwork:
		return "network"
	case ErrorTypeFileSystem:
		return "filesystem"
	default:
		return "general"
	}
}

// AnvilError represents a structured error with operation, command, and type context
type AnvilError struct {
	Op      string    // The operation being performed (init, setup, config, etc.)
	Command string    // The specific command or subcommand
	Type    ErrorType // The category of error
	Err     error     // The underlying error
	Context string    // Additional context information
}

// Error implements the error interface with improved formatting
func (e *AnvilError) Error() string {
	var parts []string

	// Build the error prefix
	if e.Command != "" {
		parts = append(parts, fmt.Sprintf("anvil %s %s", e.Op, e.Command))
	} else {
		parts = append(parts, fmt.Sprintf("anvil %s", e.Op))
	}

	// Add error type if not general
	if e.Type != ErrorTypeGeneral {
		parts = append(parts, fmt.Sprintf("[%s]", e.Type.String()))
	}

	// Add context if available
	if e.Context != "" {
		parts = append(parts, fmt.Sprintf("(%s)", e.Context))
	}

	// Join with the error message
	prefix := strings.Join(parts, " ")
	return fmt.Sprintf("%s: %v", prefix, e.Err)
}

// Unwrap returns the underlying error
func (e *AnvilError) Unwrap() error {
	return e.Err
}

// Is checks if the error matches the target error type
func (e *AnvilError) Is(target error) bool {
	if t, ok := target.(*AnvilError); ok {
		return e.Type == t.Type && e.Op == t.Op && e.Command == t.Command
	}
	return false
}

// NewAnvilError creates a new AnvilError with general type
func NewAnvilError(op, command string, err error) *AnvilError {
	return &AnvilError{
		Op:      op,
		Command: command,
		Type:    ErrorTypeGeneral,
		Err:     err,
	}
}

// NewAnvilErrorWithType creates a new AnvilError with specified type
func NewAnvilErrorWithType(op, command string, errType ErrorType, err error) *AnvilError {
	return &AnvilError{
		Op:      op,
		Command: command,
		Type:    errType,
		Err:     err,
	}
}

// NewAnvilErrorWithContext creates a new AnvilError with additional context
func NewAnvilErrorWithContext(op, command, context string, errType ErrorType, err error) *AnvilError {
	return &AnvilError{
		Op:      op,
		Command: command,
		Type:    errType,
		Context: context,
		Err:     err,
	}
}

// Helper functions for common error scenarios

// NewPlatformError creates a platform-specific error
func NewPlatformError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypePlatform, err)
}

// NewValidationError creates a validation error
func NewValidationError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypeValidation, err)
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypeConfiguration, err)
}

// NewInstallationError creates an installation error
func NewInstallationError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypeInstallation, err)
}

// NewNetworkError creates a network error
func NewNetworkError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypeNetwork, err)
}

// NewFileSystemError creates a file system error
func NewFileSystemError(op, command string, err error) *AnvilError {
	return NewAnvilErrorWithType(op, command, ErrorTypeFileSystem, err)
}

// ErrorMatches checks if an error matches specific criteria
func ErrorMatches(err error, op, command string, errType ErrorType) bool {
	if anvilErr, ok := err.(*AnvilError); ok {
		return anvilErr.Op == op && anvilErr.Command == command && anvilErr.Type == errType
	}
	return false
}

// GetErrorType extracts the error type from an AnvilError
func GetErrorType(err error) ErrorType {
	if anvilErr, ok := err.(*AnvilError); ok {
		return anvilErr.Type
	}
	return ErrorTypeGeneral
}
