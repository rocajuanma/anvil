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

package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

// mockStdin mocks stdin for testing interactive functions
func mockStdin(input string, f func()) {
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		defer w.Close()
		fmt.Fprint(w, input)
	}()

	f()
	os.Stdin = oldStdin
}

func TestPrintHeader(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:     "simple header",
			message:  "Test Header",
			expected: []string{"=== Test Header ==="},
		},
		{
			name:     "header with spaces",
			message:  "Header With Spaces",
			expected: []string{"=== Header With Spaces ==="},
		},
		{
			name:     "empty header",
			message:  "",
			expected: []string{"===  ==="},
		},
		{
			name:     "header with special characters",
			message:  "Header!@#$%",
			expected: []string{"=== Header!@#$% ==="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintHeader(tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorCyan) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestPrintStage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:     "simple stage",
			message:  "Processing...",
			expected: []string{"🔧 Processing..."},
		},
		{
			name:     "stage with details",
			message:  "Installing package xyz",
			expected: []string{"🔧 Installing package xyz"},
		},
		{
			name:     "empty stage",
			message:  "",
			expected: []string{"🔧 "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintStage(tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorBlue) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestPrintSuccess(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected []string
	}{
		{
			name:     "simple success",
			message:  "Operation completed",
			expected: []string{"✅ Operation completed"},
		},
		{
			name:     "success with details",
			message:  "Successfully installed 5 packages",
			expected: []string{"✅ Successfully installed 5 packages"},
		},
		{
			name:     "empty success",
			message:  "",
			expected: []string{"✅ "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintSuccess(tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorGreen) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected []string
	}{
		{
			name:     "simple error",
			format:   "Something went wrong",
			args:     []interface{}{},
			expected: []string{"❌ Something went wrong"},
		},
		{
			name:     "error with formatting",
			format:   "Failed to process %d items",
			args:     []interface{}{5},
			expected: []string{"❌ Failed to process 5 items"},
		},
		{
			name:     "error with multiple args",
			format:   "Error in %s: %v",
			args:     []interface{}{"function", "invalid input"},
			expected: []string{"❌ Error in function: invalid input"},
		},
		{
			name:     "empty error",
			format:   "",
			args:     []interface{}{},
			expected: []string{"❌ "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintError(tt.format, tt.args...)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorRed) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestPrintWarning(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected []string
	}{
		{
			name:     "simple warning",
			format:   "This is a warning",
			args:     []interface{}{},
			expected: []string{"⚠️  This is a warning"},
		},
		{
			name:     "warning with formatting",
			format:   "Warning: %d files will be overwritten",
			args:     []interface{}{3},
			expected: []string{"⚠️  Warning: 3 files will be overwritten"},
		},
		{
			name:     "empty warning",
			format:   "",
			args:     []interface{}{},
			expected: []string{"⚠️  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintWarning(tt.format, tt.args...)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorYellow) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestPrintInfo(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected []string
	}{
		{
			name:     "simple info",
			format:   "Information message",
			args:     []interface{}{},
			expected: []string{"Information message"},
		},
		{
			name:     "info with formatting",
			format:   "Processing %d of %d items",
			args:     []interface{}{1, 10},
			expected: []string{"Processing 1 of 10 items"},
		},
		{
			name:     "empty info",
			format:   "",
			args:     []interface{}{},
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintInfo(tt.format, tt.args...)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color reset
			if !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color reset")
			}
		})
	}
}

func TestPrintProgress(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		message  string
		expected []string
	}{
		{
			name:     "progress at start",
			current:  1,
			total:    10,
			message:  "Installing packages",
			expected: []string{"[1/10]", "10%", "Installing packages"},
		},
		{
			name:     "progress at middle",
			current:  5,
			total:    10,
			message:  "Processing items",
			expected: []string{"[5/10]", "50%", "Processing items"},
		},
		{
			name:     "progress at end",
			current:  10,
			total:    10,
			message:  "Completed",
			expected: []string{"[10/10]", "100%", "Completed"},
		},
		{
			name:     "progress with single item",
			current:  1,
			total:    1,
			message:  "Single task",
			expected: []string{"[1/1]", "100%", "Single task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintProgress(tt.current, tt.total, tt.message)
			})

			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', got: %s", expected, output)
				}
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorCyan) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestProgressNewlineAtEnd(t *testing.T) {
	// Test that progress prints a newline when current == total
	output := captureOutput(func() {
		PrintProgress(10, 10, "Done")
	})

	if !strings.HasSuffix(output, "\n") {
		t.Error("Expected progress to end with newline when current == total")
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		input    string
		expected bool
	}{
		{
			name:     "confirm with y",
			message:  "Do you want to continue?",
			input:    "y\n",
			expected: true,
		},
		{
			name:     "confirm with Y",
			message:  "Do you want to continue?",
			input:    "Y\n",
			expected: true,
		},
		{
			name:     "confirm with yes",
			message:  "Do you want to continue?",
			input:    "yes\n",
			expected: true,
		},
		{
			name:     "confirm with Yes",
			message:  "Do you want to continue?",
			input:    "Yes\n",
			expected: true,
		},
		{
			name:     "confirm with n",
			message:  "Do you want to continue?",
			input:    "n\n",
			expected: false,
		},
		{
			name:     "confirm with N",
			message:  "Do you want to continue?",
			input:    "N\n",
			expected: false,
		},
		{
			name:     "confirm with no",
			message:  "Do you want to continue?",
			input:    "no\n",
			expected: false,
		},
		{
			name:     "confirm with empty",
			message:  "Do you want to continue?",
			input:    "\n",
			expected: false,
		},
		{
			name:     "confirm with invalid input",
			message:  "Do you want to continue?",
			input:    "invalid\n",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			var output string

			mockStdin(tt.input, func() {
				output = captureOutput(func() {
					result = Confirm(tt.message)
				})
			})

			if result != tt.expected {
				t.Errorf("Expected Confirm() to return %v, got %v", tt.expected, result)
			}

			// Check that the message and prompt are displayed
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected output to contain message '%s', got: %s", tt.message, output)
			}

			if !strings.Contains(output, "(y/N):") {
				t.Errorf("Expected output to contain prompt '(y/N):', got: %s", output)
			}

			// Check for color codes
			if !strings.Contains(output, ColorBold) || !strings.Contains(output, ColorYellow) || !strings.Contains(output, ColorReset) {
				t.Error("Expected output to contain color codes")
			}
		})
	}
}

func TestIsTerminalSupported(t *testing.T) {
	// Save original TERM value
	originalTerm := os.Getenv(constants.EnvTerm)
	defer func() {
		if originalTerm != "" {
			os.Setenv(constants.EnvTerm, originalTerm)
		} else {
			os.Unsetenv(constants.EnvTerm)
		}
	}()

	tests := []struct {
		name      string
		termValue string
		expected  bool
	}{
		{
			name:      "terminal supported - normal term",
			termValue: "xterm-256color",
			expected:  true,
		},
		{
			name:      "terminal supported - basic term",
			termValue: "xterm",
			expected:  true,
		},
		{
			name:      "terminal not supported - dumb term",
			termValue: "dumb",
			expected:  false,
		},
		{
			name:      "terminal supported - empty term",
			termValue: "",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.termValue == "" {
				os.Unsetenv(constants.EnvTerm)
			} else {
				os.Setenv(constants.EnvTerm, tt.termValue)
			}

			result := IsTerminalSupported()
			if result != tt.expected {
				t.Errorf("Expected IsTerminalSupported() to return %v for TERM=%s, got %v", tt.expected, tt.termValue, result)
			}
		})
	}
}

func TestColorConstants(t *testing.T) {
	// Test that color constants are properly defined
	colorTests := []struct {
		name     string
		color    string
		expected string
	}{
		{"ColorReset", ColorReset, "\033[0m"},
		{"ColorRed", ColorRed, "\033[31m"},
		{"ColorGreen", ColorGreen, "\033[32m"},
		{"ColorYellow", ColorYellow, "\033[33m"},
		{"ColorBlue", ColorBlue, "\033[34m"},
		{"ColorPurple", ColorPurple, "\033[35m"},
		{"ColorCyan", ColorCyan, "\033[36m"},
		{"ColorWhite", ColorWhite, "\033[37m"},
		{"ColorBold", ColorBold, "\033[1m"},
	}

	for _, tt := range colorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.color)
			}
		})
	}
}

func BenchmarkPrintHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			PrintHeader("Benchmark Header")
		})
	}
}

func BenchmarkPrintInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			PrintInfo("Benchmark info message")
		})
	}
}

func BenchmarkPrintProgress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			PrintProgress(5, 10, "Benchmark progress")
		})
	}
}

func BenchmarkIsTerminalSupported(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsTerminalSupported()
	}
}
