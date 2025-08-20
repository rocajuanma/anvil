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
	"strings"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/github"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push [app-name]",
	Short: "Push configuration files to GitHub repository",
	Long:  constants.PUSH_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1), // Accept 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPushCommand(cmd, args); err != nil {
			terminal.PrintError("Push failed: %v", err)
			return
		}
	},
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
	terminal.PrintHeader(fmt.Sprintf("Push '%s' Configuration", appName))

	// Stage 1: Load and validate configuration
	terminal.PrintStage("Loading anvil configuration...")
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPush, "load-config", err)
	}

	// Validate GitHub configuration
	if anvilConfig.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPush, "missing-repo",
			fmt.Errorf("GitHub repository not configured. Please set 'github.config_repo' in your settings.yaml"))
	}
	terminal.PrintSuccess("Configuration loaded successfully")

	// Stage 2: Resolve app location
	terminal.PrintStage("Resolving app configuration location...")
	configPath, locationSource, err := config.ResolveAppLocation(appName)
	if err != nil {
		return handleAppLocationError(appName, err)
	}

	// Handle different location sources
	if locationSource == config.LocationTemp {
		terminal.PrintWarning("App '%s' found in temp directory but not configured in settings", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 To push app configurations, you need to configure the local path in settings.yaml:")
		terminal.PrintInfo("")
		terminal.PrintInfo("configs:")
		terminal.PrintInfo("  %s: /path/to/your/%s/configs", appName, appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("This ensures anvil knows where to find your local configurations.")
		terminal.PrintInfo("The temp directory (%s) contains pulled configs for review only.", configPath)
		return fmt.Errorf("app config path not configured in settings")
	}

	terminal.PrintSuccess("App configuration location resolved")
	terminal.PrintInfo("Config path: %s", configPath)

	// Stage 3: 🚨 SECURITY WARNING
	terminal.PrintWarning("🔒 SECURITY REMINDER: Configuration files contain sensitive data")
	terminal.PrintInfo("   • API keys, tokens, and credentials")
	terminal.PrintInfo("   • Personal file paths and system information")
	terminal.PrintInfo("   • Private development environment details")
	terminal.PrintInfo("")
	terminal.PrintInfo("🛡️  Anvil REQUIRES private repositories for security")
	terminal.PrintInfo("   • Repository '%s' must be PRIVATE", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("   • Public repositories will be BLOCKED")
	terminal.PrintInfo("   • Verify at: https://github.com/%s/settings", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("")

	// Stage 4: Authentication setup
	terminal.PrintStage("Setting up authentication...")
	var token string
	if anvilConfig.GitHub.TokenEnvVar != "" {
		token = os.Getenv(anvilConfig.GitHub.TokenEnvVar)
		if token == "" {
			terminal.PrintWarning("GitHub token not found in environment variable: %s", anvilConfig.GitHub.TokenEnvVar)
			terminal.PrintInfo("Proceeding with SSH authentication if available...")
		} else {
			terminal.PrintSuccess("GitHub token found in environment")
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

	terminal.PrintStage(fmt.Sprintf("Preparing to push %s configuration...", appName))
	terminal.PrintInfo("Repository: %s", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("Branch: %s", anvilConfig.GitHub.Branch)
	terminal.PrintInfo("App: %s", appName)
	terminal.PrintInfo("Local config path: %s", configPath)

	// NEW: Add diff output before confirmation
	terminal.PrintStage("Analyzing changes...")
	ctx := context.Background()
	targetPath := fmt.Sprintf("%s/", appName)
	diffSummary, err := githubClient.GetDiffPreview(ctx, configPath, targetPath)
	if err != nil {
		terminal.PrintWarning("Unable to generate diff preview: %v", err)
	} else {
		showDiffOutput(diffSummary)
	}

	// Stage 5: User confirmation
	terminal.PrintStage("Requesting user confirmation...")
	if !terminal.Confirm(fmt.Sprintf("Do you want to push your %s configurations to the repository?", appName)) {
		terminal.PrintInfo("Push cancelled by user")
		return nil
	}

	// Stage 6: Push configuration
	terminal.PrintStage(fmt.Sprintf("Pushing %s configuration to repository...", appName))
	result, err := githubClient.PushAppConfig(ctx, appName, configPath)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "push-app-config", err)
	}

	// Check if no changes were detected (result will be nil)
	if result == nil {
		// Configuration was up-to-date, success message already shown in PushAppConfig
		return nil
	}

	// Display success message for actual push
	terminal.PrintHeader("Push Complete!")
	terminal.PrintSuccess(fmt.Sprintf("%s configuration push completed successfully!", appName))
	terminal.PrintInfo("")
	terminal.PrintInfo("📋 Push Summary:")
	terminal.PrintInfo("  • Branch created: %s", result.BranchName)
	terminal.PrintInfo("  • Commit message: %s", result.CommitMessage)
	terminal.PrintInfo("  • Files committed: %v", result.FilesCommitted)
	terminal.PrintInfo("")
	terminal.PrintInfo("🔗 Repository: %s", result.RepositoryURL)
	terminal.PrintInfo("🌿 Branch: %s", result.BranchName)
	terminal.PrintInfo("")
	terminal.PrintSuccess("You can now create a Pull Request on GitHub to merge these changes!")
	terminal.PrintInfo("Direct link: %s/compare/%s...%s", result.RepositoryURL, anvilConfig.GitHub.Branch, result.BranchName)

	return nil
}

// handleAppLocationError provides helpful error messages for app location resolution failures
func handleAppLocationError(appName string, err error) error {
	if strings.Contains(err.Error(), "not found in configs or temp directory") {
		terminal.PrintError("App '%s' is not known to anvil", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 To push app configurations:")
		terminal.PrintInfo("")
		terminal.PrintInfo("1. Configure the app's local config path in settings.yaml:")
		terminal.PrintInfo("")
		terminal.PrintInfo("configs:")
		terminal.PrintInfo("  %s: /path/to/your/%s/configs", appName, appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("2. Or pull the app's configs first to discover it:")
		terminal.PrintInfo("   anvil config pull %s", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("3. Then configure the local path in settings.yaml")
		return fmt.Errorf("app not configured")
	}

	return fmt.Errorf("failed to resolve app location: %w", err)
}

// pushAnvilConfig pushes the anvil settings.yaml to the repository
func pushAnvilConfig() error {
	terminal.PrintHeader("Push Anvil Configuration")

	// Stage 1: Load and validate configuration
	terminal.PrintStage("Loading anvil configuration...")
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPush, "load-config", err)
	}

	// Validate GitHub configuration
	if anvilConfig.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPush, "missing-repo",
			fmt.Errorf("GitHub repository not configured. Please set 'github.config_repo' in your settings.yaml"))
	}
	terminal.PrintSuccess("Configuration loaded successfully")

	// 🚨 SECURITY WARNING: Remind users about private repository requirement
	terminal.PrintWarning("🔒 SECURITY REMINDER: Configuration files contain sensitive data")
	terminal.PrintInfo("   • API keys, tokens, and credentials")
	terminal.PrintInfo("   • Personal file paths and system information")
	terminal.PrintInfo("   • Private development environment details")
	terminal.PrintInfo("")
	terminal.PrintInfo("🛡️  Anvil REQUIRES private repositories for security")
	terminal.PrintInfo("   • Repository '%s' must be PRIVATE", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("   • Public repositories will be BLOCKED")
	terminal.PrintInfo("   • Verify at: https://github.com/%s/settings", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("")

	// Stage 2: Authentication setup
	terminal.PrintStage("Setting up authentication...")
	var token string
	if anvilConfig.GitHub.TokenEnvVar != "" {
		token = os.Getenv(anvilConfig.GitHub.TokenEnvVar)
		if token == "" {
			terminal.PrintWarning("GitHub token not found in environment variable: %s", anvilConfig.GitHub.TokenEnvVar)
			terminal.PrintInfo("Proceeding with SSH authentication if available...")
		} else {
			terminal.PrintSuccess("GitHub token found in environment")
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

	terminal.PrintStage("Preparing to push anvil configuration...")
	terminal.PrintInfo("Repository: %s", anvilConfig.GitHub.ConfigRepo)
	terminal.PrintInfo("Branch: %s", anvilConfig.GitHub.Branch)
	terminal.PrintInfo("Settings file: %s", settingsPath)

	// NEW: Add diff output before confirmation
	terminal.PrintStage("Analyzing changes...")
	ctx := context.Background()
	diffSummary, err := githubClient.GetDiffPreview(ctx, settingsPath, "anvil/settings.yaml")
	if err != nil {
		terminal.PrintWarning("Unable to generate diff preview: %v", err)
	} else {
		showDiffOutput(diffSummary)
	}

	// Stage 3: User confirmation
	terminal.PrintStage("Requesting user confirmation...")
	if !terminal.Confirm("Do you want to push your anvil settings to the repository?") {
		terminal.PrintInfo("Push cancelled by user")
		return nil
	}

	// Stage 4: Push configuration
	terminal.PrintStage("Pushing configuration to repository...")
	result, err := githubClient.PushAnvilConfig(ctx, settingsPath)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "push-config", err)
	}

	// Check if no changes were detected (result will be nil)
	if result == nil {
		// Configuration was up-to-date, success message already shown in PushAnvilConfig
		return nil
	}

	// Display success message for actual push
	terminal.PrintHeader("Push Complete!")
	terminal.PrintSuccess("Configuration push completed successfully!")
	terminal.PrintInfo("")
	terminal.PrintInfo("📋 Push Summary:")
	terminal.PrintInfo("  • Branch created: %s", result.BranchName)
	terminal.PrintInfo("  • Commit message: %s", result.CommitMessage)
	terminal.PrintInfo("  • Files committed: %v", result.FilesCommitted)
	terminal.PrintInfo("")
	terminal.PrintInfo("🔗 Repository: %s", result.RepositoryURL)
	terminal.PrintInfo("🌿 Branch: %s", result.BranchName)
	terminal.PrintInfo("")
	terminal.PrintSuccess("You can now create a Pull Request on GitHub to merge these changes!")
	terminal.PrintInfo("Direct link: %s/compare/%s...%s", result.RepositoryURL, anvilConfig.GitHub.Branch, result.BranchName)

	return nil
}

// showDiffOutput displays diff information using Git's native output
func showDiffOutput(diffSummary *github.DiffSummary) {
	if diffSummary.TotalFiles == 0 {
		terminal.PrintInfo("No changes detected")
		return
	}

	terminal.PrintInfo("")
	terminal.PrintHeader("📋 Changes to be pushed:")

	// Show Git's native stat output directly
	if diffSummary.GitStatOutput != "" {
		terminal.PrintInfo("")
		terminal.PrintInfo(diffSummary.GitStatOutput)
	}

	// For single small files, show full diff
	if diffSummary.TotalFiles == 1 && diffSummary.FullDiff != "" {
		lines := strings.Split(diffSummary.FullDiff, "\n")
		if len(lines) <= 50 {
			terminal.PrintInfo("")
			terminal.PrintInfo("📄 Full diff:")
			terminal.PrintInfo("")
			terminal.PrintInfo(diffSummary.FullDiff)
		} else {
			terminal.PrintInfo("")
			terminal.PrintInfo("📄 Diff preview (first 50 lines):")
			terminal.PrintInfo("")
			terminal.PrintInfo(strings.Join(lines[:50], "\n"))
			terminal.PrintInfo("")
			terminal.PrintInfo("... [diff truncated] ...")
		}
	}
	terminal.PrintInfo("")
}

func init() {
}
