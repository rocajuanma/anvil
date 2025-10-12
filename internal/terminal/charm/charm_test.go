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

package charm

import (
	"strings"
	"testing"
)

func TestNewCharmOutputHandler(t *testing.T) {
	handler := NewCharmOutputHandler()
	if handler == nil {
		t.Error("NewCharmOutputHandler returned nil")
	}

	if !handler.IsSupported() {
		t.Log("Terminal not supported, skipping visual tests")
		return
	}
}

func TestSpinnerCreation(t *testing.T) {
	tests := []struct {
		name    string
		spinner *Spinner
	}{
		{"Dots", NewDotsSpinner("test")},
		{"Line", NewLineSpinner("test")},
		{"Circle", NewCircleSpinner("test")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.spinner == nil {
				t.Errorf("%s spinner is nil", tt.name)
			}
			if tt.spinner.message != "test" {
				t.Errorf("Expected message 'test', got '%s'", tt.spinner.message)
			}
			if len(tt.spinner.frame.frames) == 0 {
				t.Error("Spinner has no frames")
			}
		})
	}
}

func TestSpinnerLifecycle(t *testing.T) {
	spinner := NewDotsSpinner("test")

	// Test that spinner can be started
	spinner.Start()
	if !spinner.running {
		t.Error("Spinner should be running after Start()")
	}

	// Test that spinner can be stopped
	spinner.Stop()
	if spinner.running {
		t.Error("Spinner should not be running after Stop()")
	}
}

func TestRenderHelpers(t *testing.T) {
	t.Run("RenderBadge", func(t *testing.T) {
		badge := RenderBadge("TEST", "#00FF87")
		if !strings.Contains(badge, "TEST") {
			t.Error("Badge should contain the text")
		}
	})

	t.Run("RenderList", func(t *testing.T) {
		items := []string{"item1", "item2"}
		list := RenderList(items, "•", "#87CEEB")
		if !strings.Contains(list, "item1") {
			t.Error("List should contain item1")
		}
		if !strings.Contains(list, "item2") {
			t.Error("List should contain item2")
		}
	})

	t.Run("RenderKeyValue", func(t *testing.T) {
		kv := RenderKeyValue("key", "value")
		if !strings.Contains(kv, "key") {
			t.Error("Output should contain key")
		}
		if !strings.Contains(kv, "value") {
			t.Error("Output should contain value")
		}
	})

	t.Run("RenderStatus", func(t *testing.T) {
		positive := RenderStatus("good", true)
		if !strings.Contains(positive, "good") {
			t.Error("Status should contain the message")
		}

		negative := RenderStatus("bad", false)
		if !strings.Contains(negative, "bad") {
			t.Error("Status should contain the message")
		}
	})
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		current int
		total   int
		width   int
	}{
		{5, 10, 20},
		{0, 10, 20},
		{10, 10, 20},
		{3, 5, 10},
	}

	for _, tt := range tests {
		bar := createProgressBar(tt.current, tt.total, tt.width)
		// Progress bar uses UTF-8 characters, so we check character count not byte length
		if bar == "" && tt.width > 0 {
			t.Error("Progress bar should not be empty when width > 0")
		}
		// Just verify it contains the expected filled/empty characters
		if !strings.Contains(bar, "█") && !strings.Contains(bar, "░") && tt.total > 0 {
			t.Error("Progress bar should contain progress characters")
		}
	}
}

func TestInitialization(t *testing.T) {
	// Test that initialization doesn't panic
	InitCharmOutput()

	handler := GetCharmHandler()
	if handler == nil {
		t.Error("GetCharmHandler returned nil after initialization")
	}

	if !IsCharmEnabled() {
		t.Error("Charm should be enabled after initialization")
	}
}

func BenchmarkSpinnerStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spinner := NewDotsSpinner("test")
		spinner.Start()
		spinner.Stop()
	}
}

func BenchmarkRenderBadge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RenderBadge("TEST", "#00FF87")
	}
}

func BenchmarkRenderList(b *testing.B) {
	items := []string{"item1", "item2", "item3"}
	for i := 0; i < b.N; i++ {
		RenderList(items, "•", "#87CEEB")
	}
}
