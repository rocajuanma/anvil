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

package setup

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// SetupCmd represents the setup command
var SetupCmd = &cobra.Command{
	Use:   "setup [group-name|app-name]",
	Short: "Install development tools and applications dynamically via Homebrew",
	Long:  constants.SETUP_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSetupCommand(cmd, args[0]); err != nil {
			terminal.PrintError("Setup failed: %v", err)
			return
		}
	},
}

// runSetupCommand executes the dynamic setup process
func runSetupCommand(cmd *cobra.Command, target string) error {
	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		return constants.NewAnvilError(constants.OpSetup, target,
			fmt.Errorf("setup command is only supported on macOS"))
	}

	// Check for list flag
	listGroups, _ := cmd.Flags().GetBool("list")
	if listGroups {
		return listAvailableGroups()
	}

	// Check for dry-run flag
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		terminal.PrintInfo("Dry run mode - no actual installations will be performed")
	}

	// Ensure Homebrew is installed
	if !brew.IsBrewInstalled() {
		terminal.PrintInfo("Homebrew not found. Installing Homebrew...")
		if err := brew.InstallBrew(); err != nil {
			return constants.NewAnvilError(constants.OpSetup, "homebrew", err)
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
		return installGroup(target, tools, dryRun)
	}

	// If not a group, treat as individual application
	return installIndividualApp(target, dryRun)
}

// installGroup installs all tools in a group
func installGroup(groupName string, tools []string, dryRun bool) error {
	terminal.PrintHeader(fmt.Sprintf("Installing '%s' group", groupName))

	if len(tools) == 0 {
		return constants.NewAnvilError(constants.OpSetup, groupName,
			fmt.Errorf("group '%s' has no tools defined", groupName))
	}

	terminal.PrintInfo("Installing %d tools: %s", len(tools), strings.Join(tools, ", "))

	successCount := 0
	var errors []string

	for i, tool := range tools {
		terminal.PrintProgress(i+1, len(tools), fmt.Sprintf("Installing %s", tool))

		if dryRun {
			terminal.PrintInfo("Would install: %s", tool)
			successCount++
		} else {
			if err := installSingleTool(tool); err != nil {
				errorMsg := fmt.Sprintf("%s: %v", tool, err)
				errors = append(errors, errorMsg)
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

	if len(errors) > 0 {
		terminal.PrintWarning("Some installations failed:")
		for _, err := range errors {
			terminal.PrintError("  • %s", err)
		}
		return constants.NewAnvilError(constants.OpSetup, groupName,
			fmt.Errorf("failed to install %d tools", len(errors)))
	}

	return nil
}

// installIndividualApp installs a single application
func installIndividualApp(appName string, dryRun bool) error {
	terminal.PrintHeader(fmt.Sprintf("Installing '%s'", appName))

	// Validate app name
	if appName == "" {
		return constants.NewAnvilError(constants.OpSetup, appName,
			fmt.Errorf("application name cannot be empty"))
	}

	// Check if already installed
	if brew.IsPackageInstalled(appName) {
		terminal.PrintSuccess(fmt.Sprintf("%s is already installed", appName))
		return nil
	}

	// Try to install the application
	if dryRun {
		terminal.PrintInfo("Would install: %s", appName)
		return nil
	}

	if err := installSingleTool(appName); err != nil {
		// Provide helpful error message with suggestions
		return constants.NewAnvilError(constants.OpSetup, appName,
			fmt.Errorf("failed to install '%s'. Please verify the name is correct. You can search for packages using 'brew search %s'", appName, appName))
	}

	terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", appName))
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
	if err := brew.InstallPackage(toolName); err != nil {
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
	terminal.PrintInfo("To complete setup, run:")
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
		return constants.NewAnvilError(constants.OpSetup, "list",
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

	terminal.PrintInfo("\nUsage: anvil setup [group-name]")
	terminal.PrintInfo("Example: anvil setup dev")

	return nil
}

func init() {
	// Add flags for additional functionality
	SetupCmd.Flags().Bool("dry-run", false, "Show what would be installed without installing")
	SetupCmd.Flags().Bool("list", false, "List all available groups")
	SetupCmd.Flags().Bool("update", false, "Update Homebrew before installation")
}
