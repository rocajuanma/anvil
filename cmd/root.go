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
}
