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
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rocajuanma/anvil/cmd/clean"
	"github.com/rocajuanma/anvil/cmd/config"
	"github.com/rocajuanma/anvil/cmd/doctor"
	"github.com/rocajuanma/anvil/cmd/initcmd"
	"github.com/rocajuanma/anvil/cmd/install"
	"github.com/rocajuanma/anvil/cmd/update"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anvil",
	Short: "🔥 One CLI to rule them all.",
	Long:  fmt.Sprintf("%s\n\n%s", constants.AnvilLogo, constants.ANVIL_LONG_DESCRIPTION),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag was used
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			showVersionInfo()
			return
		}

		showWelcomeBanner()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// showWelcomeBanner displays the enhanced welcome banner
func showWelcomeBanner() {
	// Main banner
	bannerContent := fmt.Sprintf("%s\n🔥 One CLI to rule them all 🔥\n\tversion: %s\n\n", constants.AnvilLogo, version.GetVersion())
	fmt.Println(charm.RenderBox("", bannerContent, "#FF6B9D"))

	// Quick start guide
	quickStart := `
  anvil init              		Initialize your environment
  anvil install essentials      Install essential tools
  anvil doctor            		Check system health
  anvil config pull       		Sync your dotfiles
`
	fmt.Println(charm.RenderBox("Quick Start", quickStart, "#00D9FF"))

	// Footer
	fmt.Println()
	fmt.Println("  📚 Documentation: anvil --help")
	fmt.Println("  🐛 Issues: https://github.com/rocajuanma/anvil/issues")
	fmt.Println()
}

// showVersionInfo displays the version information with branding
func showVersionInfo() {
	versionContent := fmt.Sprintf("v%s", version.GetVersion())
	fmt.Println(charm.RenderBox("ANVIL CLI", versionContent, "#FF6B9D"))
}

func init() {
	rootCmd.AddCommand(initcmd.InitCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(doctor.DoctorCmd)
	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(update.UpdateCmd)

	// Add version flag
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Set custom help template
	rootCmd.SetHelpFunc(customHelpFunc)
}

// customHelpFunc provides an enhanced help display using Charm UI
func customHelpFunc(cmd *cobra.Command, args []string) {
	// Show logo for root command
	if cmd.Name() == "anvil" {
		fmt.Println(constants.AnvilLogo)
		fmt.Println()
	}

	// Description box - format multiline descriptions with indentation
	if cmd.Long != "" {
		// Remove the logo from long description if present (already shown above)
		description := strings.ReplaceAll(cmd.Long, constants.AnvilLogo, "")
		description = strings.TrimSpace(description)

		// Split into paragraphs and add indentation
		var formattedDesc strings.Builder
		formattedDesc.WriteString("\n")
		paragraphs := strings.Split(description, "\n\n")
		for i, para := range paragraphs {
			lines := strings.Split(para, "\n")
			for _, line := range lines {
				formattedDesc.WriteString("  " + strings.TrimSpace(line) + "\n")
			}
			if i < len(paragraphs)-1 {
				formattedDesc.WriteString("\n")
			}
		}
		formattedDesc.WriteString("\n")

		fmt.Println(charm.RenderBox("About", formattedDesc.String(), "#FF6B9D"))
	} else if cmd.Short != "" {
		fmt.Println(charm.RenderBox("", "\n  "+cmd.Short+"\n", "#FF6B9D"))
	}

	// Usage section
	if cmd.HasAvailableSubCommands() {
		usageContent := fmt.Sprintf("\n  %s [command] [flags]\n", cmd.Name())
		fmt.Println(charm.RenderBox("Usage", usageContent, "#00D9FF"))
	} else {
		usageContent := fmt.Sprintf("\n  %s\n", cmd.UseLine())
		fmt.Println(charm.RenderBox("Usage", usageContent, "#00D9FF"))
	}

	// Available Commands
	if cmd.HasAvailableSubCommands() {
		var commandsContent strings.Builder
		commandsContent.WriteString("\n")

		for _, subCmd := range cmd.Commands() {
			if !subCmd.Hidden {
				commandsContent.WriteString(fmt.Sprintf("  %-12s %s\n", subCmd.Name(), subCmd.Short))
			}
		}
		commandsContent.WriteString("\n")

		fmt.Println(charm.RenderBox("Available Commands", commandsContent.String(), "#00FF87"))
	}

	// Flags
	if cmd.HasAvailableFlags() {
		var flagsContent strings.Builder
		flagsContent.WriteString("\n")
		flagsContent.WriteString(cmd.Flags().FlagUsages())

		fmt.Println(charm.RenderBox("Flags", flagsContent.String(), "#FFD700"))
	}

	// Footer
	fmt.Println()
	if cmd.HasAvailableSubCommands() {
		fmt.Println("  💡 Use 'anvil [command] --help' for more information about a command")
	}
	if cmd.Name() == "anvil" {
		fmt.Println("  📚 Documentation: https://github.com/rocajuanma/anvil")
		fmt.Println("  🐛 Issues: https://github.com/rocajuanma/anvil/issues")
	}
	fmt.Println()
}
