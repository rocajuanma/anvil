package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestAnvilError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AnvilError
		expected string
	}{
		{
			name: "general error with command",
			err: &AnvilError{
				Op:      "init",
				Command: "validate",
				Type:    ErrorTypeGeneral,
				Err:     errors.New("validation failed"),
			},
			expected: "anvil init validate: validation failed",
		},
		{
			name: "platform error without command",
			err: &AnvilError{
				Op:   "setup",
				Type: ErrorTypePlatform,
				Err:  errors.New("unsupported platform"),
			},
			expected: "anvil setup [platform]: unsupported platform",
		},
		{
			name: "installation error with context",
			err: &AnvilError{
				Op:      "setup",
				Command: "homebrew",
				Type:    ErrorTypeInstallation,
				Context: "brew install git",
				Err:     errors.New("installation failed"),
			},
			expected: "anvil setup homebrew [installation] (brew install git): installation failed",
		},
		{
			name: "configuration error with all fields",
			err: &AnvilError{
				Op:      "config",
				Command: "pull",
				Type:    ErrorTypeConfiguration,
				Context: "loading settings.yaml",
				Err:     errors.New("file not found"),
			},
			expected: "anvil config pull [configuration] (loading settings.yaml): file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("AnvilError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewAnvilError(t *testing.T) {
	err := NewAnvilError("init", "validate", fmt.Errorf("test error"))

	if err.Op != "init" {
		t.Errorf("Expected Op to be 'init', got %s", err.Op)
	}

	if err.Command != "validate" {
		t.Errorf("Expected Command to be 'validate', got %s", err.Command)
	}

	if err.Type != ErrorTypeGeneral {
		t.Errorf("Expected Type to be ErrorTypeGeneral, got %v", err.Type)
	}

	if err.Err.Error() != "test error" {
		t.Errorf("Expected underlying error to be 'test error', got %s", err.Err.Error())
	}
}

func TestNewPlatformError(t *testing.T) {
	err := NewPlatformError("setup", "install", fmt.Errorf("unsupported OS"))

	if err.Type != ErrorTypePlatform {
		t.Errorf("Expected Type to be ErrorTypePlatform, got %v", err.Type)
	}

	expected := "anvil setup install [platform]: unsupported OS"
	if err.Error() != expected {
		t.Errorf("Expected error string to be '%s', got '%s'", expected, err.Error())
	}
}

func TestErrorMatches(t *testing.T) {
	err := NewAnvilErrorWithType("init", "validate", ErrorTypeValidation, fmt.Errorf("validation failed"))

	if !ErrorMatches(err, "init", "validate", ErrorTypeValidation) {
		t.Error("Expected ErrorMatches to return true for matching error")
	}

	if ErrorMatches(err, "setup", "validate", ErrorTypeValidation) {
		t.Error("Expected ErrorMatches to return false for different operation")
	}

	if ErrorMatches(err, "init", "install", ErrorTypeValidation) {
		t.Error("Expected ErrorMatches to return false for different command")
	}

	if ErrorMatches(err, "init", "validate", ErrorTypePlatform) {
		t.Error("Expected ErrorMatches to return false for different error type")
	}
}

func TestGetErrorType(t *testing.T) {
	err := NewInstallationError("setup", "homebrew", fmt.Errorf("installation failed"))

	if GetErrorType(err) != ErrorTypeInstallation {
		t.Errorf("Expected GetErrorType to return ErrorTypeInstallation, got %v", GetErrorType(err))
	}

	// Test with non-AnvilError
	regularErr := fmt.Errorf("regular error")
	if GetErrorType(regularErr) != ErrorTypeGeneral {
		t.Errorf("Expected GetErrorType to return ErrorTypeGeneral for non-AnvilError, got %v", GetErrorType(regularErr))
	}
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		errType  ErrorType
		expected string
	}{
		{ErrorTypeGeneral, "general"},
		{ErrorTypePlatform, "platform"},
		{ErrorTypeValidation, "validation"},
		{ErrorTypeConfiguration, "configuration"},
		{ErrorTypeInstallation, "installation"},
		{ErrorTypeNetwork, "network"},
		{ErrorTypeFileSystem, "filesystem"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.errType.String(); got != tt.expected {
				t.Errorf("ErrorType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
