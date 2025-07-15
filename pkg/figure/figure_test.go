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

package figure

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
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

func TestDraw(t *testing.T) {
	tests := []struct {
		name          string
		word          string
		font          string
		expectedEmpty bool
		shouldPanic   bool
	}{
		{
			name:          "draw with standard font",
			word:          "test",
			font:          "standard",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with doh font",
			word:          "anvil",
			font:          "doh",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with big font",
			word:          "hello",
			font:          "big",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with small font",
			word:          "world",
			font:          "small",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with empty word",
			word:          "",
			font:          "standard",
			expectedEmpty: true,
			shouldPanic:   false,
		},
		{
			name:          "draw with single character",
			word:          "a",
			font:          "standard",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with spaces",
			word:          "hello world",
			font:          "standard",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with numbers",
			word:          "123",
			font:          "standard",
			expectedEmpty: false,
			shouldPanic:   false,
		},
		{
			name:          "draw with special characters",
			word:          "test!@#",
			font:          "standard",
			expectedEmpty: false,
			shouldPanic:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Draw() should have panicked but didn't")
					}
				}()
			}

			output := captureOutput(func() {
				Draw(tt.word, tt.font)
			})

			if tt.expectedEmpty {
				if strings.TrimSpace(output) != "" {
					t.Errorf("Expected empty output for empty word, got: %s", output)
				}
			} else {
				if strings.TrimSpace(output) == "" {
					t.Errorf("Expected non-empty output, got empty string")
				}

				// Verify the output contains some form of ASCII art
				lines := strings.Split(output, "\n")
				if len(lines) < 2 {
					t.Errorf("Expected multi-line ASCII art output, got: %s", output)
				}
			}
		})
	}
}

func TestDrawNoPanic(t *testing.T) {
	// Test that Draw doesn't panic with various valid inputs
	testCases := []struct {
		word string
		font string
	}{
		{"", ""},
		{"test", ""},
		{"", "standard"},
		{"very_long_word_that_might_cause_issues", "standard"},
		{"test123", "standard"},
		{"Test-With-Dashes", "standard"},
	}

	for _, tc := range testCases {
		t.Run("word:"+tc.word+"_font:"+tc.font, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Draw() panicked with word='%s' font='%s': %v", tc.word, tc.font, r)
				}
			}()

			captureOutput(func() {
				Draw(tc.word, tc.font)
			})
		})
	}
}

func TestDrawPanicCases(t *testing.T) {
	// Test cases that are expected to panic
	panicCases := []struct {
		name string
		word string
		font string
	}{
		{
			name: "invalid font should panic",
			word: "test",
			font: "invalid_font_name",
		},
	}

	for _, tc := range panicCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected Draw() to panic with word='%s' font='%s', but it didn't", tc.word, tc.font)
				}
			}()

			captureOutput(func() {
				Draw(tc.word, tc.font)
			})
		})
	}
}

func TestDrawFontVariations(t *testing.T) {
	// Test different font variations produce different outputs
	word := "test"
	fonts := []string{"standard", "doh", "big", "small"}

	outputs := make(map[string]string)

	for _, font := range fonts {
		output := captureOutput(func() {
			Draw(word, font)
		})
		outputs[font] = output
	}

	// Verify that different fonts produce different outputs
	for i, font1 := range fonts {
		for j, font2 := range fonts {
			if i != j && outputs[font1] == outputs[font2] {
				t.Errorf("Fonts %s and %s produced identical output, expected different", font1, font2)
			}
		}
	}
}

func BenchmarkDraw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			Draw("benchmark", "standard")
		})
	}
}

func BenchmarkDrawLongText(b *testing.B) {
	longText := "this is a very long text to test performance"
	for i := 0; i < b.N; i++ {
		captureOutput(func() {
			Draw(longText, "standard")
		})
	}
}
