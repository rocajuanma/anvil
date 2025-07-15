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

// PrintHeader prints a formatted header message
func PrintHeader(message string) {
	fmt.Printf("\n%s%s=== %s ===%s\n\n", ColorBold, ColorCyan, message, ColorReset)
}

// PrintStage prints a stage message with blue color
func PrintStage(message string) {
	fmt.Printf("%s%s🔧 %s%s\n", ColorBold, ColorBlue, message, ColorReset)
}

// PrintSuccess prints a success message with green color
func PrintSuccess(message string) {
	fmt.Printf("%s%s✅ %s%s\n", ColorBold, ColorGreen, message, ColorReset)
}

// PrintError prints an error message with red color and exits if requested
func PrintError(format string, args ...interface{}) {
	fmt.Printf("%s%s❌ %s%s\n", ColorBold, ColorRed, fmt.Sprintf(format, args...), ColorReset)
}

// PrintWarning prints a warning message with yellow color
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("%s%s⚠️  %s%s\n", ColorBold, ColorYellow, fmt.Sprintf(format, args...), ColorReset)
}

// PrintInfo prints an info message with normal color
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("%s%s\n", fmt.Sprintf(format, args...), ColorReset)
}

// PrintProgress prints a progress indicator
func PrintProgress(current, total int, message string) {
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\r%s%s[%d/%d] %.0f%% - %s%s", ColorBold, ColorCyan, current, total, percentage, message, ColorReset)
	if current == total {
		fmt.Println()
	}
}

// Confirm prompts the user for confirmation
func Confirm(message string) bool {
	fmt.Printf("%s%s? %s (y/N): %s", ColorBold, ColorYellow, message, ColorReset)

	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}

// IsTerminalSupported checks if the terminal supports colored output
func IsTerminalSupported() bool {
	return os.Getenv(constants.EnvTerm) != "dumb"
}
