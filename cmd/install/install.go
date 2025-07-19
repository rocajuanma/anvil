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

package install

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/installer"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install [group-name|app-name]",
	Short: "Install development tools and applications dynamically via Homebrew",
	Long:  constants.INSTALL_COMMAND_LONG_DESCRIPTION,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow no arguments if --list flag is used
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			return nil
		}
		// Otherwise, require exactly one argument
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check for list flag first
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			if err := listAvailableGroups(); err != nil {
				terminal.PrintError("Failed to list groups: %v", err)
			}
			return
		}

		if err := runInstallCommand(cmd, args[0]); err != nil {
			terminal.PrintError("Install failed: %v", err)
			return
		}
	},
}

// runInstallCommand executes the dynamic install process
func runInstallCommand(cmd *cobra.Command, target string) error {
	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		return errors.NewPlatformError(constants.OpInstall, target,
			fmt.Errorf("install command is only supported on macOS"))
	}

	// Check for dry-run flag
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		terminal.PrintInfo("Dry run mode - no actual installations will be performed")
	}

	// Check for concurrent flag
	concurrent, _ := cmd.Flags().GetBool("concurrent")
	maxWorkers, _ := cmd.Flags().GetInt("workers")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	// Ensure Homebrew is installed
	if !brew.IsBrewInstalled() {
		terminal.PrintInfo("Homebrew not found. Installing Homebrew...")
		if err := brew.InstallBrew(); err != nil {
			return errors.NewInstallationError(constants.OpInstall, "homebrew", err)
		}
		terminal.PrintSuccess("Homebrew installed successfully")
	}

	// Update Homebrew before installations
	terminal.PrintStage("Updating Homebrew...")
	if err := brew.UpdateBrew(); err != nil {
		terminal.PrintWarning("Failed to update Homebrew: %v", err)
		// Continue anyway, update failure shouldn't stop installation
	}

	// Try to get group tools first
	if tools, err := config.GetGroupTools(target); err == nil {
		return installGroup(target, tools, dryRun, concurrent, maxWorkers, timeout)
	}

	// If not a group, treat as individual application
	return installIndividualApp(target, dryRun)
}

// installGroup installs all tools in a group
func installGroup(groupName string, tools []string, dryRun bool, concurrent bool, maxWorkers int, timeout time.Duration) error {
	terminal.PrintHeader(fmt.Sprintf("Installing '%s' group", groupName))

	if len(tools) == 0 {
		return errors.NewInstallationError(constants.OpInstall, groupName,
			fmt.Errorf("group '%s' has no tools defined", groupName))
	}

	terminal.PrintInfo("Installing %d tools: %s", len(tools), strings.Join(tools, ", "))

	// Use concurrent installation if requested
	if concurrent {
		return installGroupConcurrent(groupName, tools, dryRun, maxWorkers, timeout)
	}

	// Use existing serial installation
	return installGroupSerial(groupName, tools, dryRun)
}

// installGroupConcurrent installs tools concurrently
func installGroupConcurrent(groupName string, tools []string, dryRun bool, maxWorkers int, timeout time.Duration) error {
	// Create output handler
	outputHandler := terminal.NewOutputHandler()

	// Create concurrent installer
	concurrentInstaller := installer.NewConcurrentInstaller(maxWorkers, outputHandler, dryRun)

	// Set timeout if provided
	if timeout > 0 {
		concurrentInstaller.SetTimeout(timeout)
	}

	// Create context with potential cancellation
	ctx := context.Background()

	// Install tools concurrently
	stats, err := concurrentInstaller.InstallTools(ctx, tools)

	// Track successfully installed apps
	if !dryRun && stats != nil && stats.SuccessfulTools > 0 {
		terminal.PrintInfo("Updating settings to track installed apps...")

		// For group installations, we don't track individual apps
		// since they're part of a group
		terminal.PrintInfo("Group installation tracking not implemented yet")
	}

	return err
}

// installGroupSerial installs tools serially (existing logic)
func installGroupSerial(groupName string, tools []string, dryRun bool) error {
	successCount := 0
	var installErrors []string

	for i, tool := range tools {
		terminal.PrintProgress(i+1, len(tools), fmt.Sprintf("Installing %s", tool))

		if dryRun {
			terminal.PrintInfo("Would install: %s", tool)
			successCount++
		} else {
			if err := installSingleTool(tool); err != nil {
				errorMsg := fmt.Sprintf("%s: %v", tool, err)
				installErrors = append(installErrors, errorMsg)
				terminal.PrintError("Failed to install %s: %v", tool, err)
			} else {
				successCount++
				terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", tool))
			}
		}
	}

	// Print summary
	terminal.PrintHeader("Group Installation Complete")
	terminal.PrintInfo("Successfully installed %d of %d tools", successCount, len(tools))

	if len(installErrors) > 0 {
		terminal.PrintWarning("Some installations failed:")
		for _, err := range installErrors {
			terminal.PrintError("  • %s", err)
		}
		return errors.NewInstallationError(constants.OpInstall, groupName,
			fmt.Errorf("failed to install %d tools", len(installErrors)))
	}

	return nil
}

// installIndividualApp installs a single application
func installIndividualApp(appName string, dryRun bool) error {
	terminal.PrintHeader(fmt.Sprintf("Installing '%s'", appName))

	// Validate app name
	if appName == "" {
		return errors.NewInstallationError(constants.OpInstall, appName,
			fmt.Errorf("application name cannot be empty"))
	}

	// Check if already available (via any method - Homebrew or manual installation)
	alreadyAvailable := brew.IsApplicationAvailable(appName)
	var wasNewlyInstalled bool

	if alreadyAvailable {
		terminal.PrintAlreadyAvailable("%s is already available on the system", appName)
		wasNewlyInstalled = false
	} else {
		// Try to install the application
		if dryRun {
			terminal.PrintInfo("Would install: %s", appName)
			return nil
		}

		if err := installSingleTool(appName); err != nil {
			// Provide helpful error message with suggestions
			return errors.NewInstallationError(constants.OpInstall, appName,
				fmt.Errorf("failed to install '%s'. Please verify the name is correct. You can search for packages using 'brew search %s'", appName, appName))
		}

		terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", appName))
		wasNewlyInstalled = true
	}

	// Only track the app in settings if it was newly installed and not already tracked
	if !dryRun && wasNewlyInstalled {
		// Check if already tracked to avoid duplicates
		if isTracked, err := config.IsAppTracked(appName); err != nil {
			terminal.PrintWarning("Failed to check if %s is already tracked: %v", appName, err)
		} else if isTracked {
			terminal.PrintInfo("%s is already tracked in settings", appName)
		} else {
			terminal.PrintInfo("Updating settings to track %s...", appName)
			if err := config.AddInstalledApp(appName); err != nil {
				terminal.PrintWarning("Failed to update settings file: %v", err)
				// Don't return error here as the installation was successful
			} else {
				terminal.PrintSuccess(fmt.Sprintf("Settings updated - %s is now tracked", appName))
			}
		}
	}

	return nil
}

// installSingleTool installs a single tool, handling special cases dynamically
func installSingleTool(toolName string) error {
	// Get tool-specific configuration
	toolConfig, err := config.GetToolConfig(toolName)
	if err != nil {
		terminal.PrintWarning("Failed to get tool config for %s: %v", toolName, err)
		// Continue with default installation
	}

	// Install the tool via brew
	if err := brew.InstallPackageWithCheck(toolName); err != nil {
		return fmt.Errorf("failed to install %s: %w", toolName, err)
	}

	// Handle post-install script if configured
	if toolConfig != nil && toolConfig.PostInstallScript != "" {
		terminal.PrintInfo("Running post-install script for %s...", toolName)
		if err := runPostInstallScript(toolConfig.PostInstallScript); err != nil {
			terminal.PrintWarning("Failed to run post-install script for %s: %v", toolName, err)
			// Don't fail the whole installation for this
		}
	}

	// Handle config check if configured
	if toolConfig != nil && toolConfig.ConfigCheck {
		if err := checkToolConfiguration(toolName); err != nil {
			terminal.PrintWarning("Configuration check failed for %s: %v", toolName, err)
		}
	}

	return nil
}

// runPostInstallScript runs a post-install script for a tool
func runPostInstallScript(script string) error {
	// For now, just provide instructions to the user
	terminal.PrintInfo("To complete installation, run:")
	terminal.PrintInfo("  %s", script)
	return nil
}

// checkToolConfiguration checks if a tool is properly configured
func checkToolConfiguration(toolName string) error {
	switch toolName {
	case constants.PkgGit:
		return checkGitConfiguration()
	default:
		return nil
	}
}

// checkGitConfiguration checks if git is properly configured
func checkGitConfiguration() error {
	config, err := config.LoadConfig()
	if err == nil && (config.Git.Username == "" || config.Git.Email == "") {
		terminal.PrintInfo("Git installed successfully")
		terminal.PrintWarning("Consider configuring git with:")
		terminal.PrintInfo("  git config --global user.name 'Your Name'")
		terminal.PrintInfo("  git config --global user.email 'your.email@example.com'")
	}
	return nil
}

// listAvailableGroups shows all available groups and their tools
func listAvailableGroups() error {
	terminal.PrintHeader("Available Groups")

	groups, err := config.GetAvailableGroups()
	if err != nil {
		return errors.NewConfigurationError(constants.OpInstall, "list",
			fmt.Errorf("failed to load groups: %w", err))
	}

	builtInGroups := config.GetBuiltInGroups()

	// Show built-in groups first
	terminal.PrintInfo("Built-in Groups:")
	for _, groupName := range builtInGroups {
		if tools, exists := groups[groupName]; exists {
			terminal.PrintInfo("  • %s: %s", groupName, strings.Join(tools, ", "))
		}
	}

	// Show custom groups
	hasCustomGroups := false
	for groupName := range groups {
		if !config.IsBuiltInGroup(groupName) {
			if !hasCustomGroups {
				terminal.PrintInfo("\nCustom Groups:")
				hasCustomGroups = true
			}
			terminal.PrintInfo("  • %s: %s", groupName, strings.Join(groups[groupName], ", "))
		}
	}

	if !hasCustomGroups {
		terminal.PrintInfo("\nNo custom groups defined.")
		terminal.PrintInfo("Add custom groups in ~/.anvil/settings.yaml")
	}

	// Show individually tracked installed apps
	installedApps, err := config.GetInstalledApps()
	if err != nil {
		terminal.PrintWarning("Failed to load installed apps: %v", err)
	} else if len(installedApps) > 0 {
		terminal.PrintInfo("\nIndividually Tracked Apps:")
		terminal.PrintInfo("  %s", strings.Join(installedApps, ", "))
	} else {
		terminal.PrintInfo("\nNo individually tracked apps.")
		terminal.PrintInfo("Apps installed via 'anvil install [app-name]' will be tracked automatically.")
	}

	terminal.PrintInfo("\nUsage:")
	terminal.PrintInfo("  anvil install [group-name] - Install all apps in a group")
	terminal.PrintInfo("  anvil install [app-name]   - Install individual app (auto-tracked)")
	terminal.PrintInfo("Examples:")
	terminal.PrintInfo("  anvil install dev")
	terminal.PrintInfo("  anvil install 1password")

	return nil
}

func init() {
	// Add flags for additional functionality
	InstallCmd.Flags().Bool("dry-run", false, "Show what would be installed without installing")
	InstallCmd.Flags().Bool("list", false, "List all available groups")
	InstallCmd.Flags().Bool("update", false, "Update Homebrew before installation")

	// Add concurrent installation flags
	InstallCmd.Flags().Bool("concurrent", false, "Enable concurrent installation for improved performance")
	InstallCmd.Flags().Int("workers", 0, "Number of concurrent workers (default: number of CPU cores)")
	InstallCmd.Flags().Duration("timeout", 0, "Timeout for individual tool installations (default: 10 minutes)")
}
