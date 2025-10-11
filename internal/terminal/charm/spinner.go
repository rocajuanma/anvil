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
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SpinnerFrame represents a single frame of the spinner animation
type SpinnerFrame struct {
	frames []string
	index  int
}

// Spinner provides a beautiful animated spinner for long-running operations
type Spinner struct {
	frame   *SpinnerFrame
	message string
	style   lipgloss.Style
	done    chan bool
	running bool
}

// Common spinner frame sets (these are fun!)
var (
	// Dots: Classic rotating dots
	DotsFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// Line: Simple line spinner
	LineFrames = []string{"|", "/", "-", "\\"}

	// Arrow: Circular arrow
	ArrowFrames = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	// Box: Bouncing box
	BoxFrames = []string{"◰", "◳", "◲", "◱"}

	// Circle: Filling circle
	CircleFrames = []string{"◜", "◠", "◝", "◞", "◡", "◟"}

	// Moon: Moon phases
	MoonFrames = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}

	// Pulse: Pulsing effect
	PulseFrames = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}
)

// NewSpinner creates a new spinner with the specified frames and message
func NewSpinner(frames []string, message string) *Spinner {
	return &Spinner{
		frame: &SpinnerFrame{
			frames: frames,
			index:  0,
		},
		message: message,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D9FF")).
			Bold(true),
		done:    make(chan bool),
		running: false,
	}
}

// NewDotsSpinner creates a spinner with dots animation
func NewDotsSpinner(message string) *Spinner {
	return NewSpinner(DotsFrames, message)
}

// NewLineSpinner creates a spinner with line animation
func NewLineSpinner(message string) *Spinner {
	return NewSpinner(LineFrames, message)
}

// NewCircleSpinner creates a spinner with circle animation
func NewCircleSpinner(message string) *Spinner {
	return NewSpinner(CircleFrames, message)
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if s.running {
		return
	}

	s.running = true
	go s.animate()
}

// Stop stops the spinner animation and clears the line
func (s *Spinner) Stop() {
	if !s.running {
		return
	}

	s.running = false
	s.done <- true

	// Clear the line
	fmt.Print("\r\033[K")
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.Stop()
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF87")).
		Bold(true)
	fmt.Println(successStyle.Render("✓ " + message))
}

// Error stops the spinner and shows an error message
func (s *Spinner) Error(message string) {
	s.Stop()
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true)
	fmt.Println(errorStyle.Render("✗ " + message))
}

// Warning stops the spinner and shows a warning message
func (s *Spinner) Warning(message string) {
	s.Stop()
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)
	fmt.Println(warningStyle.Render("⚠ " + message))
}

// UpdateMessage updates the spinner message without stopping it
func (s *Spinner) UpdateMessage(message string) {
	s.message = message
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			s.render()
			s.frame.index = (s.frame.index + 1) % len(s.frame.frames)
		}
	}
}

// render displays the current frame of the spinner
func (s *Spinner) render() {
	frame := s.frame.frames[s.frame.index]
	output := s.style.Render(frame + " " + s.message)
	fmt.Print("\r" + output + " ")
}

// WithStyle sets a custom style for the spinner
func (s *Spinner) WithStyle(style lipgloss.Style) *Spinner {
	s.style = style
	return s
}

// WithColor sets the color for the spinner
func (s *Spinner) WithColor(color string) *Spinner {
	s.style = s.style.Foreground(lipgloss.Color(color))
	return s
}
