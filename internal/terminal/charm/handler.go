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

	"github.com/charmbracelet/lipgloss"
	"github.com/0xjuanma/palantir"
)

// CharmOutputHandler wraps the palantir OutputHandler with enhanced lipgloss styling
type CharmOutputHandler struct {
	baseHandler palantir.OutputHandler
	styles      *StyleConfig
}

// StyleConfig holds all the lipgloss styles for different output types
type StyleConfig struct {
	Header           lipgloss.Style
	Stage            lipgloss.Style
	Success          lipgloss.Style
	Error            lipgloss.Style
	Warning          lipgloss.Style
	Info             lipgloss.Style
	AlreadyAvailable lipgloss.Style
	Progress         lipgloss.Style
	Confirm          lipgloss.Style
}

// NewCharmOutputHandler creates a new enhanced output handler with beautiful lipgloss styling
func NewCharmOutputHandler() palantir.OutputHandler {
	return &CharmOutputHandler{
		baseHandler: palantir.NewDefaultOutputHandler(),
		styles:      createDefaultStyles(),
	}
}

// NewCharmOutputHandlerWithBase creates a new enhanced output handler wrapping an existing handler
func NewCharmOutputHandlerWithBase(base palantir.OutputHandler) palantir.OutputHandler {
	return &CharmOutputHandler{
		baseHandler: base,
		styles:      createDefaultStyles(),
	}
}

// createDefaultStyles creates beautiful default styles using lipgloss
func createDefaultStyles() *StyleConfig {
	return &StyleConfig{
		// Header: Large, bold, gradient-colored banner
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B9D")).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 2).
			MarginTop(1).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B9D")),

		// Stage: Cyan with arrow, indicates progress stage
		Stage: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D9FF")).
			PaddingLeft(1),

		// Success: Green with checkmark, celebratory
		Success: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87")).
			PaddingLeft(1),

		// Error: Red with X mark, critical attention
		Error: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")).
			PaddingLeft(1),

		// Warning: Yellow/Orange with warning sign
		Warning: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			PaddingLeft(1),

		// Info: Blue with info icon, neutral information
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEEB")).
			PaddingLeft(1),

		// AlreadyAvailable: Purple/Magenta, showing existing state
		AlreadyAvailable: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C792EA")).
			Italic(true).
			PaddingLeft(1),

		// Progress: Cyan with bold counter
		Progress: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CED1")).
			PaddingLeft(1),

		// Confirm: Yellow with question mark
		Confirm: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500")),
	}
}

// PrintHeader prints a beautiful header with borders
func (c *CharmOutputHandler) PrintHeader(message string) {
	fmt.Println(c.styles.Header.Render("✨ " + message + " ✨"))
}

// PrintStage prints a stage message with an arrow
func (c *CharmOutputHandler) PrintStage(message string) {
	fmt.Println(c.styles.Stage.Render("▸ " + message))
}

// PrintSuccess prints a success message with a checkmark
func (c *CharmOutputHandler) PrintSuccess(message string) {
	fmt.Println(c.styles.Success.Render("✓ " + message))
}

// PrintError prints an error message with an X mark
func (c *CharmOutputHandler) PrintError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(c.styles.Error.Render("✗ " + message))
}

// PrintWarning prints a warning message with a warning sign
func (c *CharmOutputHandler) PrintWarning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(c.styles.Warning.Render("⚠ " + message))
}

// PrintInfo prints an info message with an info icon
func (c *CharmOutputHandler) PrintInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(c.styles.Info.Render("ℹ " + message))
}

// PrintAlreadyAvailable prints a message for already available items
func (c *CharmOutputHandler) PrintAlreadyAvailable(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(c.styles.AlreadyAvailable.Render("◆ " + message))
}

// PrintProgress prints a progress indicator with percentage
func (c *CharmOutputHandler) PrintProgress(current, total int, message string) {
	percentage := float64(current) / float64(total) * 100
	progressBar := createProgressBar(current, total, 20)

	progressText := fmt.Sprintf("[%d/%d] %.0f%% %s", current, total, percentage, progressBar)
	fmt.Printf("\r%s %s", c.styles.Progress.Render(progressText), message)

	// Print newline if this is the last item
	if current == total {
		fmt.Println()
	}
}

// createProgressBar creates a visual progress bar
func createProgressBar(current, total, width int) string {
	if total == 0 {
		return ""
	}

	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return bar
}

// Confirm prompts the user for confirmation
func (c *CharmOutputHandler) Confirm(message string) bool {
	fmt.Print(c.styles.Confirm.Render("? " + message + " (y/N): "))

	var response string
	fmt.Scanln(&response)

	switch response {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}

// IsSupported checks if the terminal supports colors
func (c *CharmOutputHandler) IsSupported() bool {
	return c.baseHandler.IsSupported()
}

// Disable disables all output
func (c *CharmOutputHandler) Disable() {
	c.baseHandler.Disable()
}
