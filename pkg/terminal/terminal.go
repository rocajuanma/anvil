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
