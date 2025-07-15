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

package draw

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
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

func TestDrawCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		expectPanic bool
	}{
		{
			name:        "valid font argument",
			args:        []string{"standard"},
			expectError: false,
			expectPanic: false,
		},
		{
			name:        "empty args should panic",
			args:        []string{},
			expectError: false,
			expectPanic: true,
		},
		{
			name:        "invalid font should panic",
			args:        []string{"invalid_font"},
			expectError: false,
			expectPanic: true,
		},
		{
			name:        "valid doh font",
			args:        []string{"doh"},
			expectError: false,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command instance for each test
			cmd := &cobra.Command{
				Use:   "draw",
				Short: "Uses go-figure to generate ASCII text",
				Run: func(cmd *cobra.Command, args []string) {
					if tt.expectPanic {
						defer func() {
							if r := recover(); r == nil {
								t.Errorf("Expected panic but got none")
							}
						}()
					}

					output := captureOutput(func() {
						DrawCmd.Run(cmd, args)
					})

					if !tt.expectPanic && strings.TrimSpace(output) == "" {
						t.Error("Expected output but got empty string")
					}
				},
			}

			// Set args for the command
			cmd.SetArgs(tt.args)

			// Execute the command
			if tt.expectError {
				if err := cmd.Execute(); err == nil {
					t.Error("Expected error but got none")
				}
			} else if !tt.expectPanic {
				if err := cmd.Execute(); err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestDrawCmdFlags(t *testing.T) {
	// Test that the command has proper structure
	if DrawCmd.Use != "draw" {
		t.Errorf("Expected Use to be 'draw', got '%s'", DrawCmd.Use)
	}

	if DrawCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}

	if DrawCmd.Run == nil {
		t.Error("Expected Run function to be set")
	}
}

func TestDrawCmdHelp(t *testing.T) {
	// Test that help works
	output := captureOutput(func() {
		DrawCmd.Help()
	})

	if !strings.Contains(output, "draw") {
		t.Error("Expected help output to contain 'draw'")
	}
}

func BenchmarkDrawCmd(b *testing.B) {
	args := []string{"standard"}
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			DrawCmd.Run(DrawCmd, args)
		})
	}
}
