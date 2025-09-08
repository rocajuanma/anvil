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
	"fmt"
	"os"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/interfaces"
)

// OutputLevel represents different levels of output
type OutputLevel int

const (
	LevelInfo OutputLevel = iota
	LevelWarning
	LevelError
	LevelSuccess
	LevelStage
	LevelHeader
)

// OutputConfig holds configuration for output formatting
type OutputConfig struct {
	UseColors     bool
	UseEmojis     bool
	UseFormatting bool
	DisableOutput bool
	VerboseMode   bool
}

// DefaultOutputHandler implements the OutputHandler interface
type DefaultOutputHandler struct {
	config *OutputConfig
}

// NewOutputHandler creates a new OutputHandler with default configuration
func NewOutputHandler() interfaces.OutputHandler {
	return &DefaultOutputHandler{
		config: &OutputConfig{
			UseColors:     true,
			UseEmojis:     true,
			UseFormatting: true,
			DisableOutput: false,
			VerboseMode:   false,
		},
	}
}

// FormatMessage formats a message according to the output level
func (oh *DefaultOutputHandler) FormatMessage(level OutputLevel, message string) string {
	if oh.config.DisableOutput {
		return ""
	}

	if !oh.IsSupported() {
		return message
	}

	var prefix, color string

	if oh.config.UseColors && oh.config.UseEmojis && oh.config.UseFormatting {
		switch level {
		case LevelHeader:
			prefix, color = "", ColorCyan
			return fmt.Sprintf("\n%s%s=== %s ===%s\n", ColorBold, color, message, ColorReset)
		case LevelStage:
			prefix, color = "🔧 ", ColorBlue
		case LevelSuccess:
			prefix, color = "✅ ", ColorGreen
		case LevelError:
			prefix, color = "❌ ", ColorRed
		case LevelWarning:
			prefix, color = "⚠️  ", ColorYellow
		case LevelInfo:
			prefix, color = "", ""
		}
	} else if oh.config.UseColors {
		switch level {
		case LevelHeader:
			prefix, color = "", ColorCyan
			return fmt.Sprintf("\n%s%s=== %s ===%s\n", ColorBold, color, message, ColorReset)
		case LevelStage:
			prefix, color = "[STAGE] ", ColorBlue
		case LevelSuccess:
			prefix, color = "[SUCCESS] ", ColorGreen
		case LevelError:
			prefix, color = "[ERROR] ", ColorRed
		case LevelWarning:
			prefix, color = "[WARNING] ", ColorYellow
		case LevelInfo:
			prefix, color = "", ""
		}
	} else {
		switch level {
		case LevelHeader:
			return fmt.Sprintf("\n=== %s ===\n", message)
		case LevelStage:
			prefix = "[STAGE] "
		case LevelSuccess:
			prefix = "[SUCCESS] "
		case LevelError:
			prefix = "[ERROR] "
		case LevelWarning:
			prefix = "[WARNING] "
		case LevelInfo:
			prefix = ""
		}
	}

	if level == LevelHeader {
		return fmt.Sprintf("\n%s%s=== %s ===%s\n", ColorBold, color, message, ColorReset)
	}

	if oh.config.UseColors && oh.config.UseFormatting {
		return fmt.Sprintf("%s%s%s%s%s\n", ColorBold, color, prefix, message, ColorReset)
	}

	return fmt.Sprintf("%s%s\n", prefix, message)
}

// PrintWithLevel prints a message with the specified level
func (oh *DefaultOutputHandler) PrintWithLevel(level OutputLevel, format string, args ...interface{}) {
	if oh.config.DisableOutput {
		return
	}

	message := fmt.Sprintf(format, args...)
	formatted := oh.FormatMessage(level, message)
	fmt.Print(formatted)
}

// Implementation of OutputHandler interface methods

func (oh *DefaultOutputHandler) PrintHeader(message string) {
	oh.PrintWithLevel(LevelHeader, message)
}

func (oh *DefaultOutputHandler) PrintStage(message string) {
	oh.PrintWithLevel(LevelStage, message)
}

func (oh *DefaultOutputHandler) PrintSuccess(message string) {
	oh.PrintWithLevel(LevelSuccess, message)
}

func (oh *DefaultOutputHandler) PrintError(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelError, format, args...)
}

func (oh *DefaultOutputHandler) PrintWarning(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelWarning, format, args...)
}

func (oh *DefaultOutputHandler) PrintInfo(format string, args ...interface{}) {
	oh.PrintWithLevel(LevelInfo, format, args...)
}

func (oh *DefaultOutputHandler) PrintAlreadyAvailable(format string, args ...interface{}) {
	if oh.config.DisableOutput {
		return
	}

	message := fmt.Sprintf(format, args...)

	if oh.config.UseColors && oh.config.UseEmojis && oh.config.UseFormatting {
		fmt.Printf("\n%s%s💙 %s%s\n", ColorBold, ColorBlue, message, ColorReset)
	} else if oh.config.UseColors {
		fmt.Printf("\n%s%s[AVAILABLE] %s%s\n", ColorBold, ColorBlue, message, ColorReset)
	} else {
		fmt.Printf("\n[AVAILABLE] %s\n", message)
	}
}

func (oh *DefaultOutputHandler) PrintProgress(current, total int, message string) {
	if oh.config.DisableOutput {
		return
	}

	percentage := float64(current) / float64(total) * 100

	if oh.config.UseColors && oh.config.UseFormatting {
		fmt.Printf("\r%s%s[%d/%d] %.0f%% - %s%s", ColorBold, ColorCyan, current, total, percentage, message, ColorReset)
	} else {
		fmt.Printf("\r[%d/%d] %.0f%% - %s", current, total, percentage, message)
	}

	if current == total {
		fmt.Println()
	}
}

func (oh *DefaultOutputHandler) Confirm(message string) bool {
	if oh.config.DisableOutput {
		return false
	}

	if oh.config.UseColors && oh.config.UseFormatting {
		fmt.Printf("%s%s? %s (y/N): %s", ColorBold, ColorYellow, message, ColorReset)
	} else {
		fmt.Printf("? %s (y/N): ", message)
	}

	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}

func (oh *DefaultOutputHandler) IsSupported() bool {
	return os.Getenv(constants.EnvTerm) != "dumb"
}

// SetVerbose enables or disables verbose mode
func (oh *DefaultOutputHandler) SetVerbose(verbose bool) {
	oh.config.VerboseMode = verbose
}

// SetColors enables or disables color output
func (oh *DefaultOutputHandler) SetColors(useColors bool) {
	oh.config.UseColors = useColors
}

// SetEmojis enables or disables emoji output
func (oh *DefaultOutputHandler) SetEmojis(useEmojis bool) {
	oh.config.UseEmojis = useEmojis
}

// Disable disables all output
func (oh *DefaultOutputHandler) Disable() {
	oh.config.DisableOutput = true
}

// Enable enables output
func (oh *DefaultOutputHandler) Enable() {
	oh.config.DisableOutput = false
}

// Global output handler instance
var globalOutputHandler interfaces.OutputHandler = NewOutputHandler()

// SetGlobalOutputHandler sets the global output handler
func SetGlobalOutputHandler(handler interfaces.OutputHandler) {
	globalOutputHandler = handler
}

// GetGlobalOutputHandler returns the global output handler
func GetGlobalOutputHandler() interfaces.OutputHandler {
	if globalOutputHandler == nil {
		globalOutputHandler = NewOutputHandler()
	}
	return globalOutputHandler
}
