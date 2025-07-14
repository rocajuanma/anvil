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

package initcmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/rocajuanma/anvil/pkg/tools"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command for macOS environment setup
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Anvil CLI environment for macOS",
	Long:  constants.INIT_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInitCommand(); err != nil {
			terminal.PrintError("Initialization failed: %v", err)
			os.Exit(1)
		}
	},
}

// runInitCommand executes the complete initialization process for Anvil CLI on macOS
func runInitCommand() error {
	terminal.PrintHeader("Anvil Initialization")

	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		return constants.NewAnvilError(constants.OpInit, "platform",
			fmt.Errorf("Anvil is only supported on macOS"))
	}

	// Stage 1: Tool validation and installation
	terminal.PrintStage("Validating and installing required tools...")
	if err := tools.ValidateAndInstallTools(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "validate-tools", err)
	}
	terminal.PrintSuccess("All required tools are available")

	// Stage 2: Create necessary directories
	terminal.PrintStage("Creating necessary directories...")
	if err := config.CreateDirectories(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "create-directories", err)
	}
	terminal.PrintSuccess("Directories created successfully")

	// Stage 3: Generate default settings.yaml
	terminal.PrintStage("Generating default settings.yaml...")
	if err := config.GenerateDefaultSettings(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "generate-settings", err)
	}
	terminal.PrintSuccess("Default settings.yaml generated")

	// Stage 4: Check local environment configurations
	terminal.PrintStage("Checking local environment configurations...")
	warnings := config.CheckEnvironmentConfigurations()
	if len(warnings) > 0 {
		terminal.PrintWarning("Environment configuration warnings:")
		for _, warning := range warnings {
			terminal.PrintWarning("  - %s", warning)
		}
	} else {
		terminal.PrintSuccess("Environment configurations are properly set")
	}

	// Stage 5: Print completion message and next steps
	terminal.PrintHeader("Initialization Complete!")
	terminal.PrintInfo("Anvil has been successfully initialized and is ready to use.")
	terminal.PrintInfo("Configuration files have been created in: %s", config.GetConfigDirectory())

	// Provide specific guidance if there are configuration warnings
	if len(warnings) > 0 {
		terminal.PrintInfo("\nRecommended next steps to complete your setup:")
		for _, warning := range warnings {
			terminal.PrintInfo("  • %s", warning)
		}
		terminal.PrintInfo("\nThese steps are optional but recommended for the best experience.")
	}

	// Final usage guidance
	terminal.PrintInfo("\nYou can now use:")
	terminal.PrintInfo("  • 'anvil setup [group]' to install development tool groups")
	terminal.PrintInfo("  • 'anvil setup [app]' to install any individual application")
	terminal.PrintInfo("  • Edit %s/settings.yaml to customize your configuration", config.GetConfigDirectory())

	// Show available groups dynamically
	if groups, err := config.GetAvailableGroups(); err == nil {
		builtInGroups := config.GetBuiltInGroups()
		terminal.PrintInfo("\nAvailable groups: %s", strings.Join(builtInGroups, ", "))
		if len(groups) > len(builtInGroups) {
			terminal.PrintInfo("Custom groups: %d defined", len(groups)-len(builtInGroups))
		}
	} else {
		terminal.PrintInfo("\nAvailable groups: dev, new-laptop")
	}
	terminal.PrintInfo("Example: 'anvil setup dev' or 'anvil setup firefox'")

	return nil
}

func init() {
	// Add flags for additional functionality
	InitCmd.Flags().Bool("force", false, "Force re-initialization even if already initialized")
	InitCmd.Flags().Bool("skip-tools", false, "Skip tool validation and installation")
}
