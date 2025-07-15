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
	"strings"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/figure"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// validFonts contains the list of supported fonts
var validFonts = []string{
	"standard", "doh", "big", "small", "banner", "block", "bubble", "digital",
	"ivrit", "lean", "mini", "script", "shadow", "slant", "speed", "term",
}

// isValidFont checks if the provided font is supported
func isValidFont(font string) bool {
	for _, validFont := range validFonts {
		if font == validFont {
			return true
		}
	}
	return false
}

var DrawCmd = &cobra.Command{
	Use:   "draw [font]",
	Short: "Uses go-figure to generate ASCII text",
	Long:  constants.DRAW_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Input validation is handled by cobra.ExactArgs(1)
		font := args[0]

		// Validate font
		if !isValidFont(font) {
			terminal.PrintError("Invalid font '%s'. Available fonts: %s", font, strings.Join(validFonts, ", "))
			return
		}

		// Draw the ASCII art
		figure.Draw("anvil", font)
	},
}

// GetValidFonts returns the list of valid fonts (useful for testing)
func GetValidFonts() []string {
	return validFonts
}

func init() {
	// Add help text showing available fonts
	DrawCmd.Long = DrawCmd.Long + "\n\nAvailable fonts: " + strings.Join(validFonts, ", ")
}
