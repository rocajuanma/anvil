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

// Package terminal provides legacy terminal output functions that delegate
// to the modern OutputHandler interface for backward compatibility.
//
// DEPRECATED: These functions are deprecated in favor of using the
// OutputHandler interface directly from pkg/terminal/output.go.
// They are kept for backward compatibility but will be removed in future versions.
package terminal

import (
	"github.com/rocajuanma/anvil/pkg/interfaces"
)

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Legacy functions that delegate to the modern OutputHandler interface
// for backward compatibility.

// getDefaultOutputHandler returns the global default output handler
func getDefaultOutputHandler() interfaces.OutputHandler {
	return GetGlobalOutputHandler()
}

// PrintHeader prints a formatted header message
// DEPRECATED: Use OutputHandler.PrintHeader() instead
func PrintHeader(message string) {
	getDefaultOutputHandler().PrintHeader(message)
}

// PrintStage prints a stage message with blue color
// DEPRECATED: Use OutputHandler.PrintStage() instead
func PrintStage(message string) {
	getDefaultOutputHandler().PrintStage(message)
}

// PrintSuccess prints a success message with green color
// DEPRECATED: Use OutputHandler.PrintSuccess() instead
func PrintSuccess(message string) {
	getDefaultOutputHandler().PrintSuccess(message)
}

// PrintError prints an error message with red color
// DEPRECATED: Use OutputHandler.PrintError() instead
func PrintError(format string, args ...interface{}) {
	getDefaultOutputHandler().PrintError(format, args...)
}

// PrintWarning prints a warning message with yellow color
// DEPRECATED: Use OutputHandler.PrintWarning() instead
func PrintWarning(format string, args ...interface{}) {
	getDefaultOutputHandler().PrintWarning(format, args...)
}

// PrintInfo prints an info message with normal color
// DEPRECATED: Use OutputHandler.PrintInfo() instead
func PrintInfo(format string, args ...interface{}) {
	getDefaultOutputHandler().PrintInfo(format, args...)
}

// PrintAlreadyAvailable prints a message indicating something is already available
// DEPRECATED: Use OutputHandler.PrintAlreadyAvailable() instead
func PrintAlreadyAvailable(format string, args ...interface{}) {
	getDefaultOutputHandler().PrintAlreadyAvailable(format, args...)
}

// PrintProgress prints a progress indicator
// DEPRECATED: Use OutputHandler.PrintProgress() instead
func PrintProgress(current, total int, message string) {
	getDefaultOutputHandler().PrintProgress(current, total, message)
}

// Confirm prompts the user for confirmation
// DEPRECATED: Use OutputHandler.Confirm() instead
func Confirm(message string) bool {
	return getDefaultOutputHandler().Confirm(message)
}

// IsTerminalSupported checks if the terminal supports colored output
// DEPRECATED: Use OutputHandler.IsSupported() instead
func IsTerminalSupported() bool {
	return getDefaultOutputHandler().IsSupported()
}
