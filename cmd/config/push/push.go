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

package push

import (
	"context"
	"fmt"
	"os"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/github"
	"github.com/rocajuanma/anvil/pkg/interfaces"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() interfaces.OutputHandler {
	return terminal.GetGlobalOutputHandler()
}

var PushCmd = &cobra.Command{
	Use:   "push [app-name]",
	Short: "Push configuration files to GitHub repository",
	Long:  constants.PUSH_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1), // Accept 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPushCommand(cmd, args); err != nil {
			getOutputHandler().PrintError("Push failed: %v", err)
			return
		}
	},
}

// isNewAppAddition checks if this is a new app that exists locally but not in remote
func isNewAppAddition(appName string, anvilConfig *config.AnvilConfig) bool {
	// Check if app exists in local configs but not in remote
	if localPath, exists := anvilConfig.Configs[appName]; exists {
		if _, err := os.Stat(localPath); err == nil {
			// App exists locally and is configured
			return true
		}
	}
	return false
}

// runPushCommand executes the configuration push process
func runPushCommand(cmd *cobra.Command, args []string) error {
	// Option 2: App-specific config push
	if len(args) > 0 {
		appName := args[0]
		return pushAppConfig(appName)
	}

	// Option 1: Anvil config push
	return pushAnvilConfig()
}

// pushAppConfig pushes application-specific configuration to the repository
func pushAppConfig(appName string) error {
	output := getOutputHandler()
	output.PrintHeader(fmt.Sprintf("Push '%s' Configuration", appName))

	// Stage 1: Load and validate configuration
	output.PrintStage("Loading anvil configuration...")
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPush, "load-config", err)
	}

	// Validate GitHub configuration
	if anvilConfig.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPush, "missing-repo",
			fmt.Errorf("GitHub repository not configured. Please set 'github.config_repo' in your settings.yaml"))
	}
	output.PrintSuccess("Configuration loaded successfully")

	// Stage 2: Resolve app location
	output.PrintStage("Resolving app configuration location...")
	configPath, locationSource, err := config.ResolveAppLocation(appName)
	if err != nil {
		// Check if this is a new app addition
		if isNewAppAddition(appName, anvilConfig) {
			output.PrintInfo("🆕 New app '%s' detected - will be added to repository", appName)
			// Get the configured path for new apps
			if localPath, exists := anvilConfig.Configs[appName]; exists {
				configPath = localPath
			} else {
				return handleAppLocationError(appName, err)
			}
		} else {
			return handleAppLocationError(appName, err)
		}
	}

	// Handle different location sources
	if locationSource == config.LocationTemp {
		output.PrintWarning("App '%s' found in temp directory but not configured in settings\n", appName)
		output.PrintInfo("💡 To push app configurations, you need to configure the local path in settings.yaml:\n")
		output.PrintInfo("configs:")
		output.PrintInfo("  %s: /path/to/your/%s/configs\n", appName, appName)
		output.PrintInfo("This ensures anvil knows where to find your local configurations.")
		output.PrintInfo("The temp directory (%s) contains pulled configs for review only.", configPath)
		return fmt.Errorf("app config path not configured in settings")
	}

	output.PrintSuccess("App configuration location resolved")
	output.PrintInfo("Config path: %s", configPath)

	// Show new app information if this is a new addition
	if isNewAppAddition(appName, anvilConfig) {
		showNewAppInfo(appName, configPath)
	}

	// Stage 3: 🚨 SECURITY WARNING
	showSecurityWarning(anvilConfig.GitHub.ConfigRepo)

	// Stage 4: Authentication setup
	output.PrintStage("Setting up authentication...")
	var token string
	if anvilConfig.GitHub.TokenEnvVar != "" {
		token = os.Getenv(anvilConfig.GitHub.TokenEnvVar)
		if token == "" {
			output.PrintWarning("GitHub token not found in environment variable: %s", anvilConfig.GitHub.TokenEnvVar)
			output.PrintInfo("Proceeding with SSH authentication if available...")
		} else {
			output.PrintSuccess("GitHub token found in environment")
		}
	}

	// Create GitHub client
	githubClient := github.NewGitHubClient(
		anvilConfig.GitHub.ConfigRepo,
		anvilConfig.GitHub.Branch,
		anvilConfig.GitHub.LocalPath,
		token,
		anvilConfig.Git.SSHKeyPath,
		anvilConfig.Git.Username,
		anvilConfig.Git.Email,
	)

	output.PrintStage(fmt.Sprintf("Preparing to push %s configuration...", appName))
	output.PrintInfo("Repository: %s", anvilConfig.GitHub.ConfigRepo)
	output.PrintInfo("Branch: %s", anvilConfig.GitHub.Branch)
	output.PrintInfo("App: %s", appName)
	output.PrintInfo("Local config path: %s", configPath)

	// NEW: Add diff output before confirmation
	output.PrintStage("Analyzing changes...")
	ctx := context.Background()
	targetPath := fmt.Sprintf("%s/", appName)
	diffSummary, err := githubClient.GetDiffPreview(ctx, configPath, targetPath)
	if err != nil {
		output.PrintWarning("Unable to generate diff preview: %v", err)
	} else {
		showDiffOutput(diffSummary)
	}

	// Stage 5: User confirmation
	output.PrintStage("Requesting user confirmation...")
	if !output.Confirm(fmt.Sprintf("Do you want to push your %s configurations to the repository?", appName)) {
		output.PrintInfo("Push cancelled by user")
		return nil
	}

	// Stage 6: Push configuration
	output.PrintStage(fmt.Sprintf("Pushing %s configuration to repository...", appName))
	result, err := githubClient.PushAppConfig(ctx, appName, configPath)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "push-app-config", err)
	}

	// Check if no changes were detected (result will be nil)
	if result == nil {
		// Configuration was up-to-date, success message already shown in PushAppConfig
		return nil
	}

	displaySuccessMessage(appName, result, diffSummary, anvilConfig)

	return nil
}

// pushAnvilConfig pushes the anvil settings.yaml to the repository
func pushAnvilConfig() error {
	output := getOutputHandler()
	output.PrintHeader("Push Anvil Configuration")

	// Stage 1: Load and validate configuration
	output.PrintStage("Loading anvil configuration...")
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPush, "load-config", err)
	}

	// Validate GitHub configuration
	if anvilConfig.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPush, "missing-repo",
			fmt.Errorf("GitHub repository not configured. Please set 'github.config_repo' in your settings.yaml"))
	}
	output.PrintSuccess("Configuration loaded successfully")

	showSecurityWarning(anvilConfig.GitHub.ConfigRepo)

	// Stage 2: Authentication setup
	output.PrintStage("Setting up authentication...")
	var token string
	if anvilConfig.GitHub.TokenEnvVar != "" {
		token = os.Getenv(anvilConfig.GitHub.TokenEnvVar)
		if token == "" {
			output.PrintWarning("GitHub token not found in environment variable: %s", anvilConfig.GitHub.TokenEnvVar)
			output.PrintInfo("Proceeding with SSH authentication if available...\n")
		} else {
			output.PrintSuccess("GitHub token found in environment\n")
		}
	}

	// Create GitHub client
	githubClient := github.NewGitHubClient(
		anvilConfig.GitHub.ConfigRepo,
		anvilConfig.GitHub.Branch,
		anvilConfig.GitHub.LocalPath,
		token,
		anvilConfig.Git.SSHKeyPath,
		anvilConfig.Git.Username,
		anvilConfig.Git.Email,
	)

	// Get settings file path
	settingsPath := config.GetConfigPath()

	output.PrintStage("Preparing to push anvil configuration...")
	output.PrintInfo("Repository: %s", anvilConfig.GitHub.ConfigRepo)
	output.PrintInfo("Branch: %s", anvilConfig.GitHub.Branch)
	output.PrintInfo("Settings file: %s", settingsPath)

	// NEW: Add diff output before confirmation
	output.PrintStage("Analyzing changes...")
	ctx := context.Background()
	diffSummary, err := githubClient.GetDiffPreview(ctx, settingsPath, "anvil/settings.yaml")
	if err != nil {
		output.PrintWarning("Unable to generate diff preview: %v", err)
	} else {
		showDiffOutput(diffSummary)
	}

	// Stage 3: User confirmation
	output.PrintStage("Requesting user confirmation...")
	if !output.Confirm("Do you want to push your anvil settings to the repository?") {
		output.PrintInfo("Push cancelled by user")
		return nil
	}

	// Stage 4: Push configuration
	output.PrintStage("Pushing configuration to repository...")
	result, err := githubClient.PushAnvilConfig(ctx, settingsPath)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "push-config", err)
	}

	// Check if no changes were detected (result will be nil)
	if result == nil {
		// Configuration was up-to-date, success message already shown in PushAnvilConfig
		return nil
	}

	displaySuccessMessage("anvil", result, diffSummary, anvilConfig)

	return nil
}

func init() {
}
