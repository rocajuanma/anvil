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
	"github.com/rocajuanma/anvil/pkg/interfaces"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() interfaces.OutputHandler {
	return terminal.GetGlobalOutputHandler()
}

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install [group-name|app-name] [--group-name group]",
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
				getOutputHandler().PrintError("Failed to list groups: %v", err)
			}
			return
		}

		if err := runInstallCommand(cmd, args[0]); err != nil {
			getOutputHandler().PrintError("Install failed: %v", err)
			return
		}
	},
}

// runInstallCommand executes the dynamic install process
func runInstallCommand(cmd *cobra.Command, target string) error {
	o := getOutputHandler()
	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		return errors.NewPlatformError(constants.OpInstall, target,
			fmt.Errorf("install command is only supported on macOS"))
	}

	// Check for dry-run flag
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		o.PrintInfo("Dry run mode - no actual installations will be performed")
	}

	// Check for concurrent flag
	concurrent, _ := cmd.Flags().GetBool("concurrent")
	maxWorkers, _ := cmd.Flags().GetInt("workers")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	// Ensure Homebrew is installed
	if !brew.IsBrewInstalled() {
		o.PrintInfo("Homebrew not found. Installing Homebrew...")
		if err := brew.InstallBrew(); err != nil {
			return errors.NewInstallationError(constants.OpInstall, "homebrew", err)
		}
		o.PrintSuccess("Homebrew installed successfully")
	}

	// Update Homebrew before installations
	o.PrintStage("Updating Homebrew...")
	if err := brew.UpdateBrew(); err != nil {
		o.PrintWarning("Failed to update Homebrew: %v", err)
		// Continue anyway, update failure shouldn't stop installation
	}

	// Try to get group tools first
	if tools, err := config.GetGroupTools(target); err == nil {
		return installGroup(target, tools, dryRun, concurrent, maxWorkers, timeout)
	}

	// If not a group, treat as individual application
	return installIndividualApp(target, dryRun, cmd)
}

// installGroup installs all tools in a group
func installGroup(groupName string, tools []string, dryRun bool, concurrent bool, maxWorkers int, timeout time.Duration) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Installing '%s' group", groupName))

	if len(tools) == 0 {
		return errors.NewInstallationError(constants.OpInstall, groupName,
			fmt.Errorf("group '%s' has no tools defined", groupName))
	}

	// Deduplicate tools within the group and update settings if needed
	deduplicatedTools, err := deduplicateGroupTools(groupName, tools)
	if err != nil {
		o.PrintWarning("Failed to deduplicate group tools: %v", err)
		// Continue with original tools list
		deduplicatedTools = tools
	} else {
		// Use the deduplicated tools list for installation
		tools = deduplicatedTools
	}

	o.PrintInfo("Installing %d tools: %s", len(tools), strings.Join(tools, ", "))

	// Use concurrent installation if requested
	if concurrent {
		return installGroupConcurrent(groupName, tools, dryRun, maxWorkers, timeout)
	}

	// Use existing serial installation
	return installGroupSerial(groupName, tools, dryRun)
}

// deduplicateGroupTools removes duplicate tools within a group and updates the settings file
func deduplicateGroupTools(groupName string, tools []string) ([]string, error) {
	// Track which tools we've seen
	seen := make(map[string]bool)
	var deduplicatedTools []string
	var duplicatesFound []string

	// Build deduplicated list
	for _, tool := range tools {
		if !seen[tool] {
			seen[tool] = true
			deduplicatedTools = append(deduplicatedTools, tool)
		} else {
			duplicatesFound = append(duplicatesFound, tool)
		}
	}

	// If no duplicates found, return original list
	if len(duplicatesFound) == 0 {
		return tools, nil
	}

	// Report found duplicates
	o := getOutputHandler()
	o.PrintWarning("Found duplicates in group '%s': %s", groupName, strings.Join(duplicatesFound, ", "))
	o.PrintInfo("Removing duplicates from settings file...")

	// Update the configuration with deduplicated tools
	if err := config.UpdateGroupTools(groupName, deduplicatedTools); err != nil {
		return tools, fmt.Errorf("failed to update group with deduplicated tools: %w", err)
	}

	o.PrintSuccess(fmt.Sprintf("Successfully removed %d duplicate(s) from group '%s'", len(duplicatesFound), groupName))
	return deduplicatedTools, nil
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
		o := getOutputHandler()
		o.PrintInfo("Updating settings to track installed apps...")

		// For group installations, we don't track individual apps
		// since they're part of a group
		o.PrintInfo("Group installation tracking not implemented yet")
	}

	return err
}

// installGroupSerial installs tools serially using unified installation logic
func installGroupSerial(groupName string, tools []string, dryRun bool) error {
	successCount := 0
	var installErrors []string

	for i, tool := range tools {
		getOutputHandler().PrintProgress(i+1, len(tools), fmt.Sprintf("\nInstalling %s", tool))

		// Use unified installation logic - this ensures consistent behavior with availability checking
		_, err := installSingleToolUnified(tool, dryRun)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: %v", tool, err)
			installErrors = append(installErrors, errorMsg)
			getOutputHandler().PrintError("Failed to install %s: %v", tool, err)
		} else {
			successCount++
		}
	}

	// Use unified error reporting
	return reportGroupInstallationResults(groupName, successCount, len(tools), installErrors)
}

// installIndividualApp installs a single application using unified installation logic
func installIndividualApp(appName string, dryRun bool, cmd *cobra.Command) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Installing '%s'", appName))

	// Validate app name
	if appName == "" {
		return errors.NewInstallationError(constants.OpInstall, appName,
			fmt.Errorf("application name cannot be empty"))
	}

	// Use unified installation logic
	wasNewlyInstalled, err := installSingleToolUnified(appName, dryRun)
	if err != nil {
		// Provide helpful error message with suggestions
		return errors.NewInstallationError(constants.OpInstall, appName,
			fmt.Errorf("failed to install '%s'. Please verify the name is correct. You can search for packages using 'brew search %s'", appName, appName))
	}

	// Only track the app in settings if it was newly installed and not dry-run
	if !dryRun && wasNewlyInstalled {
		// Check if --group-name flag is provided
		groupName, _ := cmd.Flags().GetString("group-name")
		if groupName != "" {
			// Add app to the specified group
			if err := config.AddAppToGroup(groupName, appName); err != nil {
				o.PrintWarning("Failed to add %s to group '%s': %v", appName, groupName, err)
				// Continue with normal tracking as fallback
				return trackAppInSettings(appName)
			}
			o.PrintSuccess(fmt.Sprintf("Added %s to group '%s'", appName, groupName))
			return nil
		} else {
			// Normal tracking in installed_apps
			return trackAppInSettings(appName)
		}
	}

	return nil
}

// installSingleTool installs a single tool, handling special cases dynamically
func installSingleTool(toolName string) error {
	o := getOutputHandler()
	// Get tool-specific configuration
	toolConfig, err := config.GetToolConfig(toolName)
	if err != nil {
		o.PrintWarning("Failed to get tool config for %s: %v", toolName, err)
		// Continue with default installation
	}

	// Install the tool via brew
	if err := brew.InstallPackageWithCheck(toolName); err != nil {
		return fmt.Errorf("failed to install %s: %w", toolName, err)
	}

	// Handle post-install script if configured
	if toolConfig != nil && toolConfig.PostInstallScript != "" {
		o.PrintInfo("Running post-install script for %s...", toolName)
		if err := runPostInstallScript(toolConfig.PostInstallScript); err != nil {
			o.PrintWarning("Failed to run post-install script for %s: %v", toolName, err)
			// Don't fail the whole installation for this
		}
	}

	// Handle config check if configured
	if toolConfig != nil && toolConfig.ConfigCheck {
		if err := checkToolConfiguration(toolName); err != nil {
			o.PrintWarning("Configuration check failed for %s: %v", toolName, err)
		}
	}

	return nil
}

// installSingleToolUnified provides unified installation logic for all installation modes
// This is the core function that ensures consistent behavior across individual, serial, and concurrent installations
func installSingleToolUnified(toolName string, dryRun bool) (wasNewlyInstalled bool, err error) {
	o := getOutputHandler()
	// ALWAYS check availability first using the latest IsApplicationAvailable logic
	if brew.IsApplicationAvailable(toolName) {
		o.PrintAlreadyAvailable("%s is already available on the system", toolName)
		return false, nil
	}

	// Handle installation based on mode
	if dryRun {
		o.PrintInfo("Would install: %s", toolName)
		return true, nil
	}

	// Perform real installation using existing logic
	if err := installSingleTool(toolName); err != nil {
		return false, fmt.Errorf("failed to install %s: %w", toolName, err)
	}

	o.PrintSuccess(fmt.Sprintf("%s installed successfully", toolName))
	return true, nil
}

// trackAppInSettings handles adding newly installed apps to settings
func trackAppInSettings(appName string) error {
	o := getOutputHandler()
	// Check if already tracked to avoid duplicates
	if isTracked, err := config.IsAppTracked(appName); err != nil {
		o.PrintWarning("Failed to check if %s is already tracked: %v", appName, err)
		return nil // Don't fail installation for tracking issues
	} else if isTracked {
		o.PrintInfo("%s is already tracked in settings", appName)
		return nil
	}

	o.PrintInfo("Updating settings to track %s...", appName)
	if err := config.AddInstalledApp(appName); err != nil {
		o.PrintWarning("Failed to update settings file: %v", err)
		return nil // Don't fail installation for tracking issues
	}

	o.PrintSuccess(fmt.Sprintf("Settings updated - %s is now tracked", appName))
	return nil
}

// reportGroupInstallationResults provides unified error reporting for group installations
func reportGroupInstallationResults(groupName string, successCount, totalCount int, installErrors []string) error {
	// Print summary
	o := getOutputHandler()
	o.PrintHeader("Group Installation Complete")
	o.PrintInfo("Successfully installed %d of %d tools", successCount, totalCount)

	if len(installErrors) > 0 {
		o.PrintWarning("Some installations failed:")
		for _, err := range installErrors {
			o.PrintError("  • %s", err)
		}
		return errors.NewInstallationError(constants.OpInstall, groupName,
			fmt.Errorf("failed to install %d tools", len(installErrors)))
	}

	return nil
}

// runPostInstallScript runs a post-install script for a tool
func runPostInstallScript(script string) error {
	// For now, just provide instructions to the user
	o := getOutputHandler()
	o.PrintInfo("To complete installation, run:")
	o.PrintInfo("  %s", script)
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
		o := getOutputHandler()
		o.PrintInfo("Git installed successfully")
		o.PrintWarning("Consider configuring git with:")
		o.PrintInfo("  git config --global user.name 'Your Name'")
		o.PrintInfo("  git config --global user.email 'your.email@example.com'")
	}
	return nil
}

// listAvailableGroups shows all available groups and their tools
func listAvailableGroups() error {
	o := getOutputHandler()
	o.PrintHeader("Available Groups")

	groups, err := config.GetAvailableGroups()
	if err != nil {
		return errors.NewConfigurationError(constants.OpInstall, "list",
			fmt.Errorf("failed to load groups: %w", err))
	}

	builtInGroups := config.GetBuiltInGroups()

	// Show built-in groups first
	o.PrintInfo("Built-in Groups:")
	for _, groupName := range builtInGroups {
		if tools, exists := groups[groupName]; exists {
			o.PrintInfo("  • %s: %s", groupName, strings.Join(tools, ", "))
		}
	}

	// Show custom groups
	hasCustomGroups := false
	for groupName := range groups {
		if !config.IsBuiltInGroup(groupName) {
			if !hasCustomGroups {
				o.PrintInfo("\nCustom Groups:")
				hasCustomGroups = true
			}
			o.PrintInfo("  • %s: %s", groupName, strings.Join(groups[groupName], ", "))
		}
	}

	if !hasCustomGroups {
		o.PrintInfo("\nNo custom groups defined.")
		o.PrintInfo("Add custom groups in ~/%s/%s", constants.AnvilConfigDir, constants.ConfigFileName)
	}

	// Show individually tracked installed apps
	installedApps, err := config.GetInstalledApps()
	if err != nil {
		o.PrintWarning("Failed to load installed apps: %v", err)
	} else if len(installedApps) > 0 {
		o.PrintInfo("\nIndividually Tracked Apps:")
		o.PrintInfo("  %s", strings.Join(installedApps, ", "))
	} else {
		o.PrintInfo("\nNo individually tracked apps.")
		o.PrintInfo("Apps installed via 'anvil install [app-name]' will be tracked automatically.")
	}

	o.PrintInfo("\nUsage:")
	o.PrintInfo("  anvil install [group-name] - Install all apps in a group")
	o.PrintInfo("  anvil install [app-name]   - Install individual app (auto-tracked)")
	o.PrintInfo("  anvil install [app-name] --group-name [group] - Install app and add to group")
	o.PrintInfo("Examples:")
	o.PrintInfo("  anvil install dev")
	o.PrintInfo("  anvil install 1password")
	o.PrintInfo("  anvil install firefox --group-name essentials")
	o.PrintInfo("  anvil install final-cut --group-name editing")

	return nil
}

func init() {
	// Add flags for additional functionality
	InstallCmd.Flags().Bool("dry-run", false, "Show what would be installed without installing")
	InstallCmd.Flags().Bool("list", false, "List all available groups")
	InstallCmd.Flags().Bool("update", false, "Update Homebrew before installation")
	InstallCmd.Flags().String("group-name", "", "Add the installed app to a group (creates group if it doesn't exist)")

	// Add concurrent installation flags
	InstallCmd.Flags().Bool("concurrent", false, "Enable concurrent installation for improved performance")
	InstallCmd.Flags().Int("workers", 0, "Number of concurrent workers (default: number of CPU cores)")
	InstallCmd.Flags().Duration("timeout", 0, "Timeout for individual tool installations (default: 10 minutes)")
}
