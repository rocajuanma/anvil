/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package system

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
		expectCode  int
	}{
		{
			name:        "successful echo command",
			command:     "echo",
			args:        []string{"hello world"},
			expectError: false,
			expectCode:  0,
		},
		{
			name:        "successful pwd command",
			command:     "pwd",
			args:        []string{},
			expectError: false,
			expectCode:  0,
		},
		{
			name:        "nonexistent command",
			command:     "nonexistent-command-12345",
			args:        []string{},
			expectError: true,
			expectCode:  0, // Command not found doesn't set exit code
		},
		{
			name:        "command with non-zero exit",
			command:     "sh",
			args:        []string{"-c", "exit 1"},
			expectError: true,
			expectCode:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunCommand(tt.command, tt.args...)

			if err != nil {
				t.Errorf("RunCommand() returned error: %v", err)
			}

			if result == nil {
				t.Fatal("RunCommand() returned nil result")
			}

			if tt.expectError && result.Success {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && !result.Success {
				t.Errorf("Expected command to succeed but it failed: %s", result.Error)
			}

			if tt.expectCode != 0 && result.ExitCode != tt.expectCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectCode, result.ExitCode)
			}

			if result.Command == "" {
				t.Error("Result command string should not be empty")
			}
		})
	}
}

func TestRunCommandWithTimeout(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		command       string
		args          []string
		expectError   bool
		shouldTimeout bool
	}{
		{
			name:          "quick command with long timeout",
			timeout:       5 * time.Second,
			command:       "echo",
			args:          []string{"test"},
			expectError:   false,
			shouldTimeout: false,
		},
		{
			name:          "slow command with short timeout",
			timeout:       100 * time.Millisecond,
			command:       "sleep",
			args:          []string{"1"},
			expectError:   true,
			shouldTimeout: true,
		},
		{
			name:          "quick command with very short timeout",
			timeout:       10 * time.Millisecond,
			command:       "echo",
			args:          []string{"fast"},
			expectError:   false,
			shouldTimeout: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			result, err := RunCommandWithTimeout(ctx, tt.command, tt.args...)

			if err != nil {
				t.Errorf("RunCommandWithTimeout() returned error: %v", err)
			}

			if result == nil {
				t.Fatal("RunCommandWithTimeout() returned nil result")
			}

			if tt.expectError && result.Success {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && !result.Success {
				t.Errorf("Expected command to succeed but it failed: %s", result.Error)
			}
		})
	}
}

func TestRunCommandWithOutput(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "successful echo command",
			command:     "echo",
			args:        []string{"test output"},
			expectError: false,
		},
		{
			name:        "nonexistent command",
			command:     "nonexistent-command-12345",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "command with non-zero exit",
			command:     "sh",
			args:        []string{"-c", "exit 1"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunCommandWithOutput(tt.command, tt.args...)

			if tt.expectError && err == nil {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected command to succeed but it failed: %v", err)
			}
		})
	}
}

func TestRunCommandWithOutputTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "quick command with long timeout",
			timeout:     5 * time.Second,
			command:     "echo",
			args:        []string{"test"},
			expectError: false,
		},
		{
			name:        "slow command with short timeout",
			timeout:     100 * time.Millisecond,
			command:     "sleep",
			args:        []string{"1"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := RunCommandWithOutputTimeout(ctx, tt.command, tt.args...)

			if tt.expectError && err == nil {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected command to succeed but it failed: %v", err)
			}
		})
	}
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "existing command - echo",
			command:  "echo",
			expected: true,
		},
		{
			name:     "existing command - pwd",
			command:  "pwd",
			expected: true,
		},
		{
			name:     "nonexistent command",
			command:  "nonexistent-command-12345",
			expected: false,
		},
		{
			name:     "empty command",
			command:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CommandExists(tt.command)
			if result != tt.expected {
				t.Errorf("CommandExists(%s) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestGetCommandPath(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "existing command - echo",
			command:     "echo",
			expectError: false,
		},
		{
			name:        "nonexistent command",
			command:     "nonexistent-command-12345",
			expectError: true,
		},
		{
			name:        "empty command",
			command:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetCommandPath(tt.command)

			if tt.expectError && err == nil {
				t.Error("Expected GetCommandPath to return error but it didn't")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected GetCommandPath to succeed but it failed: %v", err)
			}

			if !tt.expectError && path == "" {
				t.Error("Expected non-empty path for existing command")
			}
		})
	}
}

func TestRunCommandInDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		dir         string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "pwd in temp directory",
			dir:         tempDir,
			command:     "pwd",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "ls in temp directory",
			dir:         tempDir,
			command:     "ls",
			args:        []string{"-la"},
			expectError: false,
		},
		{
			name:        "nonexistent directory",
			dir:         "/nonexistent/directory/path",
			command:     "pwd",
			args:        []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunCommandInDirectory(tt.dir, tt.command, tt.args...)

			if err != nil {
				t.Errorf("RunCommandInDirectory() returned error: %v", err)
			}

			if result == nil {
				t.Fatal("RunCommandInDirectory() returned nil result")
			}

			if tt.expectError && result.Success {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && !result.Success {
				t.Errorf("Expected command to succeed but it failed: %s", result.Error)
			}

			if !tt.expectError && !strings.Contains(result.Command, tt.dir) {
				t.Errorf("Expected command string to contain directory %s", tt.dir)
			}
		})
	}
}

func TestRunCommandInDirectoryWithTimeout(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		timeout     time.Duration
		dir         string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "quick command in temp directory",
			timeout:     5 * time.Second,
			dir:         tempDir,
			command:     "pwd",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "slow command with short timeout",
			timeout:     100 * time.Millisecond,
			dir:         tempDir,
			command:     "sleep",
			args:        []string{"1"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			result, err := RunCommandInDirectoryWithTimeout(ctx, tt.dir, tt.command, tt.args...)

			if err != nil {
				t.Errorf("RunCommandInDirectoryWithTimeout() returned error: %v", err)
			}

			if result == nil {
				t.Fatal("RunCommandInDirectoryWithTimeout() returned nil result")
			}

			if tt.expectError && result.Success {
				t.Error("Expected command to fail but it succeeded")
			}

			if !tt.expectError && !result.Success {
				t.Errorf("Expected command to succeed but it failed: %s", result.Error)
			}
		})
	}
}

func TestGetEnvironmentVariable(t *testing.T) {
	testKey := "TEST_ANVIL_VAR"
	testValue := "test_value"
	defaultValue := "default_value"

	// Clean up any existing value
	originalValue := os.Getenv(testKey)
	defer func() {
		if originalValue != "" {
			os.Setenv(testKey, originalValue)
		} else {
			os.Unsetenv(testKey)
		}
	}()

	tests := []struct {
		name     string
		key      string
		default_ string
		setValue string
		expected string
	}{
		{
			name:     "existing environment variable",
			key:      testKey,
			default_: defaultValue,
			setValue: testValue,
			expected: testValue,
		},
		{
			name:     "nonexistent environment variable",
			key:      testKey,
			default_: defaultValue,
			setValue: "",
			expected: defaultValue,
		},
		{
			name:     "empty environment variable",
			key:      testKey,
			default_: defaultValue,
			setValue: "",
			expected: defaultValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset the environment variable
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnvironmentVariable(tt.key, tt.default_)
			if result != tt.expected {
				t.Errorf("GetEnvironmentVariable(%s, %s) = %s, expected %s", tt.key, tt.default_, result, tt.expected)
			}
		})
	}
}

func TestSetEnvironmentVariable(t *testing.T) {
	testKey := "TEST_ANVIL_SET_VAR"
	testValue := "test_set_value"

	// Clean up
	originalValue := os.Getenv(testKey)
	defer func() {
		if originalValue != "" {
			os.Setenv(testKey, originalValue)
		} else {
			os.Unsetenv(testKey)
		}
	}()

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "set normal variable",
			key:   testKey,
			value: testValue,
		},
		{
			name:  "set empty value",
			key:   testKey,
			value: "",
		},
		{
			name:  "set variable with spaces",
			key:   testKey,
			value: "value with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetEnvironmentVariable(tt.key, tt.value)
			if err != nil {
				t.Errorf("SetEnvironmentVariable(%s, %s) returned error: %v", tt.key, tt.value, err)
			}

			// Verify the value was set
			actualValue := os.Getenv(tt.key)
			if actualValue != tt.value {
				t.Errorf("Expected environment variable %s to be %s, got %s", tt.key, tt.value, actualValue)
			}
		})
	}
}

func TestCommandResultStruct(t *testing.T) {
	// Test CommandResult struct initialization and fields
	result := &CommandResult{
		Command:  "test command",
		ExitCode: 1,
		Output:   "test output",
		Error:    "test error",
		Success:  false,
	}

	if result.Command != "test command" {
		t.Errorf("Expected Command to be 'test command', got '%s'", result.Command)
	}

	if result.ExitCode != 1 {
		t.Errorf("Expected ExitCode to be 1, got %d", result.ExitCode)
	}

	if result.Output != "test output" {
		t.Errorf("Expected Output to be 'test output', got '%s'", result.Output)
	}

	if result.Error != "test error" {
		t.Errorf("Expected Error to be 'test error', got '%s'", result.Error)
	}

	if result.Success != false {
		t.Errorf("Expected Success to be false, got %v", result.Success)
	}
}

func BenchmarkRunCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunCommand("echo", "benchmark test")
	}
}

func BenchmarkCommandExists(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CommandExists("echo")
	}
}

func BenchmarkGetEnvironmentVariable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetEnvironmentVariable("PATH", "default")
	}
}

// Test platform-specific behavior
func TestPlatformSpecificCommands(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		t.Run("windows dir command", func(t *testing.T) {
			result, err := RunCommand("dir")
			if err != nil {
				t.Errorf("RunCommand(dir) returned error: %v", err)
			}
			if result == nil {
				t.Fatal("RunCommand(dir) returned nil result")
			}
		})
	case "darwin", "linux":
		t.Run("unix ls command", func(t *testing.T) {
			result, err := RunCommand("ls")
			if err != nil {
				t.Errorf("RunCommand(ls) returned error: %v", err)
			}
			if result == nil {
				t.Fatal("RunCommand(ls) returned nil result")
			}
		})
	}
}
