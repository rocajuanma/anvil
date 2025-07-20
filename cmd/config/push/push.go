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
	"gopkg.in/yaml.v2"
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
	// Option 2: App-specific config push (in development)
	if len(args) > 0 {
		appName := args[0]
		return showAppPushInDevelopment(appName)
	}

	// Option 1: Anvil config push
	return pushAnvilConfig(cmd)
}

// showAppPushInDevelopment displays development message for app config push
func showAppPushInDevelopment(appName string) error {
	terminal.PrintHeader(fmt.Sprintf("Push '%s' Configuration", appName))
	terminal.PrintWarning("Application-specific configuration push is currently in development")
	terminal.PrintInfo("This feature will allow you to push %s configuration files to your GitHub repository", appName)
	terminal.PrintInfo("Expected functionality:")
	terminal.PrintInfo("  • Create timestamped branch: config-push-<DDMMYYYY>-<HHMM>")
	terminal.PrintInfo("  • Commit message: anvil[push]: %s", appName)
	terminal.PrintInfo("  • Push %s configs to /%s directory in repository", appName, appName)
	terminal.PrintInfo("  • Create pull request for review")
	terminal.PrintInfo("")
	terminal.PrintInfo("🚧 Status: In Development")
	terminal.PrintInfo("📅 Expected: Future release")
	terminal.PrintInfo("")
	terminal.PrintInfo("For now, use 'anvil config push' to push anvil settings only.")
	return nil
}

// pushAnvilConfig pushes the anvil settings.yaml to the repository
func pushAnvilConfig(cmd *cobra.Command) error {
	terminal.PrintHeader("Push Anvil Configuration")

	// Get command flags
	showFiltered, _ := cmd.Flags().GetBool("show-filtered")
	includeOverride, _ := cmd.Flags().GetStringSlice("include")
	excludeOverride, _ := cmd.Flags().GetStringSlice("exclude")

	// Load anvil configuration
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPush, "load-config", err)
	}

	// Apply flag overrides to sync configuration
	if len(includeOverride) > 0 {
		anvilConfig.SyncConfig.IncludeOverride = append(anvilConfig.SyncConfig.IncludeOverride, includeOverride...)
		terminal.PrintInfo("Override: Including sections: %s", strings.Join(includeOverride, ", "))
	}
	if len(excludeOverride) > 0 {
		anvilConfig.SyncConfig.ExcludeSections = append(anvilConfig.SyncConfig.ExcludeSections, excludeOverride...)
		terminal.PrintInfo("Override: Excluding sections: %s", strings.Join(excludeOverride, ", "))
	}

	// If show-filtered flag is set, display filtered configuration and exit
	if showFiltered {
		return showFilteredConfiguration(anvilConfig)
	}

	// Validate GitHub configuration
	if anvilConfig.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPush, "missing-repo",
			fmt.Errorf("GitHub repository not configured. Please set 'github.config_repo' in your settings.yaml"))
	}

	// Get GitHub token
	var token string
	if anvilConfig.GitHub.TokenEnvVar != "" {
		token = os.Getenv(anvilConfig.GitHub.TokenEnvVar)
		if token == "" {
			terminal.PrintWarning("GitHub token not found in environment variable: %s", anvilConfig.GitHub.TokenEnvVar)
			terminal.PrintInfo("Proceeding with SSH authentication if available...")
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

	// Confirm with user
	if !terminal.Confirm("Do you want to push your anvil settings to the repository?") {
		terminal.PrintInfo("Push cancelled by user")
		return nil
	}

	// Push configuration
	ctx := context.Background()
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

// showFilteredConfiguration displays what would be pushed with filtering applied
func showFilteredConfiguration(anvilConfig *config.AnvilConfig) error {
	terminal.PrintStage("Showing filtered configuration preview...")

	// Apply filtering
	filteredConfig, err := config.FilterForSync(anvilConfig)
	if err != nil {
		return fmt.Errorf("failed to apply filtering: %w", err)
	}

	// Convert to YAML for display
	yamlData, err := yaml.Marshal(filteredConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filtered config: %w", err)
	}

	terminal.PrintHeader("Filtered Configuration (what would be pushed)")

	if len(anvilConfig.SyncConfig.ExcludeSections) > 0 {
		terminal.PrintInfo("🚫 Excluded sections: %s", strings.Join(anvilConfig.SyncConfig.ExcludeSections, ", "))
	}
	if len(anvilConfig.SyncConfig.TemplateSections) > 0 {
		terminal.PrintInfo("📝 Templated sections: %s", strings.Join(anvilConfig.SyncConfig.TemplateSections, ", "))
	}
	if len(anvilConfig.SyncConfig.IncludeOverride) > 0 {
		terminal.PrintInfo("✅ Force included: %s", strings.Join(anvilConfig.SyncConfig.IncludeOverride, ", "))
	}

	terminal.PrintInfo("")
	terminal.PrintInfo("Filtered settings.yaml content:")
	terminal.PrintInfo("=" + strings.Repeat("=", 50))
	fmt.Println(string(yamlData))
	terminal.PrintInfo("=" + strings.Repeat("=", 50))

	terminal.PrintSuccess("Preview complete. Use 'anvil config push' without --show-filtered to actually push.")
	return nil
}

func init() {
	// Add flags for enhanced filtering functionality
	PushCmd.Flags().Bool("show-filtered", false, "Show what would be pushed (filtered configuration) without pushing")
	PushCmd.Flags().StringSlice("include", []string{}, "Force include sections (overrides exclude configuration)")
	PushCmd.Flags().StringSlice("exclude", []string{}, "Force exclude sections (adds to exclude configuration)")
}
