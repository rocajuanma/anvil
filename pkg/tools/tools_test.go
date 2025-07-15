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

package tools

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/rocajuanma/anvil/pkg/constants"
)

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestGetRequiredTools(t *testing.T) {
	tools := GetRequiredTools()

	if len(tools) == 0 {
		t.Fatal("Expected at least one required tool")
	}

	expectedTools := map[string]bool{
		"Git":      false,
		"cURL":     false,
		"Homebrew": false,
	}

	for _, tool := range tools {
		if !tool.Required {
			t.Errorf("Tool %s should be required but isn't", tool.Name)
		}

		if tool.Name == "" {
			t.Error("Tool name should not be empty")
		}

		if tool.Command == "" {
			t.Error("Tool command should not be empty")
		}

		if tool.Description == "" {
			t.Error("Tool description should not be empty")
		}

		if tool.InstallWith == "" {
			t.Error("Tool InstallWith should not be empty")
		}

		if _, exists := expectedTools[tool.Name]; exists {
			expectedTools[tool.Name] = true
		}
	}

	// Verify all expected tools are present
	for toolName, found := range expectedTools {
		if !found {
			t.Errorf("Expected required tool %s was not found", toolName)
		}
	}
}

func TestGetOptionalTools(t *testing.T) {
	tools := GetOptionalTools()

	// Optional tools list can be empty, but if present should be valid
	for _, tool := range tools {
		if tool.Required {
			t.Errorf("Tool %s should not be required but is", tool.Name)
		}

		if tool.Name == "" {
			t.Error("Tool name should not be empty")
		}

		if tool.Command == "" {
			t.Error("Tool command should not be empty")
		}

		if tool.Description == "" {
			t.Error("Tool description should not be empty")
		}

		if tool.InstallWith == "" {
			t.Error("Tool InstallWith should not be empty")
		}
	}
}

func TestToolStructure(t *testing.T) {
	// Test Tool struct initialization
	tool := Tool{
		Name:        "Test Tool",
		Command:     "test-cmd",
		Required:    true,
		InstallWith: "brew",
		Description: "Test description",
	}

	if tool.Name != "Test Tool" {
		t.Errorf("Expected Name to be 'Test Tool', got '%s'", tool.Name)
	}

	if tool.Command != "test-cmd" {
		t.Errorf("Expected Command to be 'test-cmd', got '%s'", tool.Command)
	}

	if !tool.Required {
		t.Error("Expected Required to be true")
	}

	if tool.InstallWith != "brew" {
		t.Errorf("Expected InstallWith to be 'brew', got '%s'", tool.InstallWith)
	}

	if tool.Description != "Test description" {
		t.Errorf("Expected Description to be 'Test description', got '%s'", tool.Description)
	}
}

func TestGetToolInfo(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		expectError bool
	}{
		{
			name:        "find tool by name - Git",
			toolName:    "Git",
			expectError: false,
		},
		{
			name:        "find tool by command - git",
			toolName:    "git",
			expectError: false,
		},
		{
			name:        "find tool by name - cURL",
			toolName:    "cURL",
			expectError: false,
		},
		{
			name:        "find tool by command - curl",
			toolName:    "curl",
			expectError: false,
		},
		{
			name:        "find tool by name - Homebrew",
			toolName:    "Homebrew",
			expectError: false,
		},
		{
			name:        "find tool by command - brew",
			toolName:    "brew",
			expectError: false,
		},
		{
			name:        "nonexistent tool",
			toolName:    "nonexistent-tool",
			expectError: true,
		},
		{
			name:        "empty tool name",
			toolName:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := GetToolInfo(tt.toolName)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError && tool == nil {
				t.Error("Expected tool info but got nil")
			}

			if !tt.expectError && tool != nil {
				if tool.Name == "" {
					t.Error("Tool name should not be empty")
				}
				if tool.Command == "" {
					t.Error("Tool command should not be empty")
				}
			}
		})
	}
}

func TestListTools(t *testing.T) {
	// Test that ListTools doesn't panic and produces output
	output := captureOutput(func() {
		ListTools()
	})

	if output == "" {
		t.Error("Expected ListTools to produce output")
	}

	// Check for expected content
	expectedContent := []string{
		"Required Tools",
		"Optional Tools",
		"Git",
		"cURL",
		"Homebrew",
	}

	for _, content := range expectedContent {
		if !strings.Contains(output, content) {
			t.Errorf("Expected output to contain '%s'", content)
		}
	}
}

func TestCheckToolsStatus(t *testing.T) {
	// Skip test if not on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("CheckToolsStatus test requires macOS")
	}

	status, err := CheckToolsStatus()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if status == nil {
		t.Fatal("Expected status map but got nil")
	}

	// Check that all expected tools are in the status map
	expectedTools := []string{"Git", "cURL", "Homebrew"}
	for _, toolName := range expectedTools {
		if _, exists := status[toolName]; !exists {
			t.Errorf("Expected tool %s in status map", toolName)
		}
	}
}

func TestCheckToolsStatusNonMacOS(t *testing.T) {
	// Skip test if on macOS
	if runtime.GOOS == "darwin" {
		t.Skip("This test is for non-macOS systems")
	}

	status, err := CheckToolsStatus()
	if err == nil {
		t.Error("Expected error for non-macOS system")
	}

	if status != nil {
		t.Error("Expected nil status for non-macOS system")
	}

	if !strings.Contains(err.Error(), "only supported on macOS") {
		t.Error("Expected error message to mention macOS requirement")
	}
}

func TestValidateAndInstallToolsNonMacOS(t *testing.T) {
	// Skip test if on macOS
	if runtime.GOOS == "darwin" {
		t.Skip("This test is for non-macOS systems")
	}

	err := ValidateAndInstallTools()
	if err == nil {
		t.Error("Expected error for non-macOS system")
	}

	if !strings.Contains(err.Error(), "only supports macOS") {
		t.Error("Expected error message to mention macOS requirement")
	}
}

func TestToolRequiredVsOptional(t *testing.T) {
	requiredTools := GetRequiredTools()
	optionalTools := GetOptionalTools()

	// Verify that required tools are marked as required
	for _, tool := range requiredTools {
		if !tool.Required {
			t.Errorf("Required tool %s should be marked as required", tool.Name)
		}
	}

	// Verify that optional tools are not marked as required
	for _, tool := range optionalTools {
		if tool.Required {
			t.Errorf("Optional tool %s should not be marked as required", tool.Name)
		}
	}
}

func TestToolInstallationMethods(t *testing.T) {
	allTools := append(GetRequiredTools(), GetOptionalTools()...)

	validInstallMethods := map[string]bool{
		"brew":   true,
		"script": true,
		"system": true,
	}

	for _, tool := range allTools {
		if !validInstallMethods[tool.InstallWith] {
			t.Errorf("Tool %s has invalid installation method: %s", tool.Name, tool.InstallWith)
		}
	}
}

func TestToolCommandConstants(t *testing.T) {
	// Test that tools use proper constants
	requiredTools := GetRequiredTools()

	gitFound := false
	curlFound := false
	brewFound := false

	for _, tool := range requiredTools {
		switch tool.Command {
		case constants.GitCommand:
			gitFound = true
			if tool.Name != "Git" {
				t.Errorf("Expected Git tool to have name 'Git', got '%s'", tool.Name)
			}
		case constants.CurlCommand:
			curlFound = true
			if tool.Name != "cURL" {
				t.Errorf("Expected cURL tool to have name 'cURL', got '%s'", tool.Name)
			}
		case constants.BrewCommand:
			brewFound = true
			if tool.Name != "Homebrew" {
				t.Errorf("Expected Homebrew tool to have name 'Homebrew', got '%s'", tool.Name)
			}
		}
	}

	if !gitFound {
		t.Error("Expected to find Git tool with constants.GitCommand")
	}
	if !curlFound {
		t.Error("Expected to find cURL tool with constants.CurlCommand")
	}
	if !brewFound {
		t.Error("Expected to find Homebrew tool with constants.BrewCommand")
	}
}

func TestToolUniqueness(t *testing.T) {
	allTools := append(GetRequiredTools(), GetOptionalTools()...)

	nameMap := make(map[string]bool)
	commandMap := make(map[string]bool)

	for _, tool := range allTools {
		// Check for duplicate names
		if nameMap[tool.Name] {
			t.Errorf("Duplicate tool name found: %s", tool.Name)
		}
		nameMap[tool.Name] = true

		// Check for duplicate commands
		if commandMap[tool.Command] {
			t.Errorf("Duplicate tool command found: %s", tool.Command)
		}
		commandMap[tool.Command] = true
	}
}

func TestGetToolInfoCaseInsensitive(t *testing.T) {
	// Test that tool lookup works with different cases
	tests := []struct {
		name        string
		toolName    string
		expectError bool
	}{
		{
			name:        "lowercase git",
			toolName:    "git",
			expectError: false,
		},
		{
			name:        "uppercase GIT",
			toolName:    "GIT",
			expectError: true, // Should not match because it's case sensitive
		},
		{
			name:        "mixed case Git",
			toolName:    "Git",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := GetToolInfo(tt.toolName)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError && tool == nil {
				t.Error("Expected tool info but got nil")
			}
		})
	}
}

func BenchmarkGetRequiredTools(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRequiredTools()
	}
}

func BenchmarkGetOptionalTools(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetOptionalTools()
	}
}

func BenchmarkGetToolInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetToolInfo("Git")
	}
}

func BenchmarkListTools(b *testing.B) {
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			ListTools()
		})
	}
}

// Test helper functions
func TestToolValidation(t *testing.T) {
	// Test that all tools have proper validation
	allTools := append(GetRequiredTools(), GetOptionalTools()...)

	for _, tool := range allTools {
		// Test that all required fields are present
		if tool.Name == "" {
			t.Error("Tool name should not be empty")
		}
		if tool.Command == "" {
			t.Error("Tool command should not be empty")
		}
		if tool.InstallWith == "" {
			t.Error("Tool InstallWith should not be empty")
		}
		if tool.Description == "" {
			t.Error("Tool description should not be empty")
		}

		// Test that command doesn't contain spaces (should be single command)
		if strings.Contains(tool.Command, " ") {
			t.Errorf("Tool command should not contain spaces: %s", tool.Command)
		}
	}
}

func TestToolConsistency(t *testing.T) {
	// Test that tool data is consistent
	requiredTools := GetRequiredTools()

	// Git should be required and installable with brew
	for _, tool := range requiredTools {
		if tool.Name == "Git" {
			if !tool.Required {
				t.Error("Git should be required")
			}
			if tool.InstallWith != "brew" {
				t.Errorf("Git should be installable with brew, got %s", tool.InstallWith)
			}
			if tool.Command != "git" {
				t.Errorf("Git command should be 'git', got %s", tool.Command)
			}
		}
	}
}
