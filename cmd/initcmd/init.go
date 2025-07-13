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
	"os"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/rocajuanma/anvil/pkg/tools"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command, which bootstraps the Anvil CLI environment
// This command performs a complete initialization process including tool validation,
// directory creation, configuration generation, and environment checking.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap and initialize your Anvil CLI environment",
	Long:  constants.INIT_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInitCommand(); err != nil {
			terminal.PrintError("Initialization failed: %v", err)
			os.Exit(1)
		}
	},
}

// runInitCommand executes the complete initialization process for Anvil CLI
// This function orchestrates all initialization stages and handles errors gracefully
func runInitCommand() error {
	terminal.PrintHeader("Anvil Initialization")

	// Stage 1: Tool validation and installation
	// This stage ensures all required tools are available and installs missing ones
	terminal.PrintStage("Validating and installing required tools...")
	if err := tools.ValidateAndInstallTools(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "validate-tools", err)
	}
	terminal.PrintSuccess("All required tools are available")

	// Stage 2: Create necessary directories
	// This stage creates the Anvil configuration directory structure
	terminal.PrintStage("Creating necessary directories...")
	if err := config.CreateDirectories(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "create-directories", err)
	}
	terminal.PrintSuccess("Directories created successfully")

	// Stage 3: Generate default settings.yaml
	// This stage creates the main configuration file with sensible defaults
	terminal.PrintStage("Generating default settings.yaml...")
	if err := config.GenerateDefaultSettings(); err != nil {
		return constants.NewAnvilError(constants.OpInit, "generate-settings", err)
	}
	terminal.PrintSuccess("Default settings.yaml generated")

	// Stage 4: Check local environment configurations
	// This stage validates the local development environment and provides recommendations
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
	// This stage provides the user with completion status and actionable next steps
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
	terminal.PrintInfo("  • 'anvil --help' to see all available commands")
	terminal.PrintInfo("  • 'anvil setup' to install development tools")
	terminal.PrintInfo("  • 'anvil pull/push' to synchronize assets with GitHub")
	terminal.PrintInfo("  • Edit ~/.anvil/settings.yaml to customize your configuration")

	return nil
}

func init() {
	// Future flag definitions can be added here
	// Example flags that might be useful:
	// InitCmd.Flags().BoolP("force", "f", false, "Force re-initialization even if already initialized")
	// InitCmd.Flags().BoolP("minimal", "m", false, "Perform minimal initialization without optional tools")
	// InitCmd.Flags().StringP("config", "c", "", "Use custom configuration file")
}
