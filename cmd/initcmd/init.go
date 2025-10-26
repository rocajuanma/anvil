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
	"strings"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/tools"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// InitCmd represents the init command for macOS environment setup
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Anvil CLI environment for macOS",
	Long:  constants.INIT_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInitCommand(); err != nil {
			getOutputHandler().PrintError("Initialization failed: %v", err)
			os.Exit(1)
		}
	},
}

// runInitCommand executes the complete initialization process for Anvil CLI on macOS
func runInitCommand() error {
	// Display initialization banner
	fmt.Println(charm.RenderBox("🔨 ANVIL INITIALIZATION", "", "#00D9FF", true))
	fmt.Println()

	o := getOutputHandler()

	// Stage 1: Tool validation and installation
	o.PrintStage("Stage 1: Tool Validation")
	spinner := charm.NewCircleSpinner("Validating and installing required tools")
	spinner.Start()
	if err := tools.ValidateAndInstallTools(); err != nil {
		spinner.Error("Tool validation failed")
		return errors.NewValidationError(constants.OpInit, "validate-tools", err)
	}
	spinner.Success("All required tools are available")

	// Stage 2: Create necessary directories
	o.PrintStage("Stage 2: Directory Creation")
	spinner = charm.NewDotsSpinner("Creating necessary directories")
	spinner.Start()
	if err := config.CreateDirectories(); err != nil {
		spinner.Error("Failed to create directories")
		return errors.NewFileSystemError(constants.OpInit, "create-directories", err)
	}
	spinner.Success("Directories created successfully")

	// Stage 3: Generate default settings.yaml
	o.PrintStage("Stage 3: Settings Generation")
	spinner = charm.NewDotsSpinner("Generating default settings.yaml")
	spinner.Start()
	if err := config.GenerateDefaultSettings(); err != nil {
		spinner.Error("Failed to generate settings")
		return errors.NewConfigurationError(constants.OpInit, "generate-settings", err)
	}
	spinner.Success("Default settings.yaml generated")

	// Stage 4: Check local environment configurations
	o.PrintStage("Stage 4: Environment Check")
	spinner = charm.NewLineSpinner("Checking local environment configurations")
	spinner.Start()
	warnings := config.CheckEnvironmentConfigurations()
	if len(warnings) > 0 {
		spinner.Warning("Environment configuration warnings found")
		for _, warning := range warnings {
			o.PrintWarning("  - %s", warning)
		}
	} else {
		spinner.Success("Environment configurations are properly set")
	}

	// Stage 5: Print completion message and next steps
	o.PrintHeader("Initialization Complete!")
	o.PrintInfo("Anvil has been successfully initialized and is ready to use.")
	o.PrintInfo("Configuration files have been created in: %s", config.GetConfigDirectory())

	// Provide specific guidance if there are configuration warnings
	if len(warnings) > 0 {
		o.PrintInfo("\nRecommended next steps to complete your setup:")
		for _, warning := range warnings {
			o.PrintInfo("  • %s", warning)
		}
		o.PrintInfo("\nThese steps are optional but recommended for the best experience.")
	}

	// Final usage guidance
	o.PrintInfo("\nYou can now use:")
	o.PrintInfo("  • 'anvil install [group]' to install development tool groups")
	o.PrintInfo("  • 'anvil install [app]' to install any individual application")
	o.PrintInfo("  • Edit %s/settings.yaml to customize your configuration", config.GetConfigDirectory())

	// GitHub configuration warning
	o.PrintWarning("\n⚙️  Configuration Management Setup Required:")
	o.PrintInfo("  • Edit the 'github.config_repo' field in settings.yaml to enable config pull/push")
	o.PrintInfo("  • Example: 'github.config_repo: username/dotfiles'")
	o.PrintInfo("  • Set GITHUB_TOKEN environment variable for authentication")
	o.PrintInfo("  • Run 'anvil config pull' once configured to sync your dotfiles")

	// Show available groups dynamically
	if groups, err := config.GetAvailableGroups(); err == nil {
		builtInGroups := config.GetBuiltInGroups()
		o.PrintInfo("\nAvailable groups: %s", strings.Join(builtInGroups, ", "))
		if len(groups) > len(builtInGroups) {
			o.PrintInfo("Custom groups: %d defined", len(groups)-len(builtInGroups))
		}
	} else {
		o.PrintInfo("\nAvailable groups: dev, essentials")
	}
	o.PrintInfo("Example: 'anvil install dev' or 'anvil install firefox'")

	return nil
}

func init() {
	// Add flags for additional functionality
	InitCmd.Flags().Bool("skip-tools", false, "Skip tool validation and installation")
}
