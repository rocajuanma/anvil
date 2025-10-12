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
	"sort"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/internal/brew"
	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/installer"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install [group-name|app-name] [--group-name group]",
	Short: "Install development tools and applications dynamically via Homebrew",
	Long:  constants.INSTALL_COMMAND_LONG_DESCRIPTION,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow no arguments if --list or --tree flag is used
		listFlag, _ := cmd.Flags().GetBool("list")
		treeFlag, _ := cmd.Flags().GetBool("tree")
		if listFlag || treeFlag {
			return nil
		}
		// Otherwise, require exactly one argument
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check for tree or list flag
		treeFlag, _ := cmd.Flags().GetBool("tree")
		listFlag, _ := cmd.Flags().GetBool("list")

		if treeFlag || listFlag {
			// Load and prepare data once
			groups, builtInGroupNames, customGroupNames, installedApps, err := loadAndPrepareAppData()
			if err != nil {
				getOutputHandler().PrintError("Failed to load application data: %v", err)
				return
			}

			// Choose rendering based on flag
			var content string
			var title = "Available Applications"
			if treeFlag {
				content = renderTreeView(groups, builtInGroupNames, customGroupNames, installedApps)
				title = fmt.Sprintf("%s (Tree View)", title)
			} else {
				content = renderListView(groups, builtInGroupNames, customGroupNames, installedApps)
				title = fmt.Sprintf("%s (List View)", title)
			}

			// Display in box
			fmt.Println(charm.RenderBox(title, content, "#00D9FF", false))
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
	if err := brew.EnsureBrewIsInstalled(); err != nil {
		return fmt.Errorf("install: %w", err)
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
	} else {
		tools = deduplicatedTools
	}

	o.PrintInfo("Installing %d tools: %s", len(tools), strings.Join(tools, ", "))

	if concurrent {
		return installGroupConcurrent(groupName, tools, dryRun, maxWorkers, timeout)
	}

	return installGroupSerial(groupName, tools, dryRun)
}

// deduplicateGroupTools removes duplicate tools within a group and updates the settings file
func deduplicateGroupTools(groupName string, tools []string) ([]string, error) {
	seen := make(map[string]struct{}, len(tools))
	deduplicatedTools := make([]string, 0, len(tools))
	var duplicatesFound []string

	// Deduplicate
	for _, tool := range tools {
		if _, exists := seen[tool]; !exists {
			seen[tool] = struct{}{}
			deduplicatedTools = append(deduplicatedTools, tool)
		} else {
			duplicatesFound = append(duplicatesFound, tool)
		}
	}

	// Return original list if no duplicates found
	if len(duplicatesFound) == 0 {
		return tools, nil
	}

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
	o := getOutputHandler()

	// Create new output handler to send into concurrent installer
	outputHandler := palantir.NewDefaultOutputHandler()
	concurrentInstaller := installer.NewConcurrentInstaller(maxWorkers, outputHandler, dryRun)

	if timeout > 0 {
		concurrentInstaller.SetTimeout(timeout)
	}

	// Create context with potential cancellation
	ctx := context.Background()
	stats, err := concurrentInstaller.InstallTools(ctx, tools)

	// Track successfully installed apps
	if !dryRun && stats != nil && stats.SuccessfulTools > 0 {
		o.PrintInfo("Updating settings to track installed apps...")
		o.PrintInfo("Group installation tracking not implemented yet")
	}

	return err
}

// toolStatus represents the status of a tool installation
type toolStatus struct {
	name   string
	status string // "pending", "installing", "done", "failed"
	emoji  string
}

// installGroupSerial installs tools serially using unified installation logic
func installGroupSerial(groupName string, tools []string, dryRun bool) error {
	o := getOutputHandler()

	successCount := 0
	var installErrors []string

	// Initialize tool statuses
	toolStatuses := make([]toolStatus, len(tools))
	for i, tool := range tools {
		toolStatuses[i] = toolStatus{
			name:   tool,
			status: "pending",
			emoji:  "⋯",
		}
	}

	for i, tool := range tools {
		// Update status to installing
		toolStatuses[i].status = "installing"
		toolStatuses[i].emoji = "⠋"

		// Print dashboard
		printInstallDashboard(groupName, toolStatuses, i+1, len(tools))

		// Use unified installation logic
		_, err := installSingleToolUnified(tool, dryRun)

		if err != nil {
			toolStatuses[i].status = "failed"
			toolStatuses[i].emoji = "✗"
			errorMsg := fmt.Sprintf("%s: %v", tool, err)
			installErrors = append(installErrors, errorMsg)
			o.PrintError("%s: %v", tool, err)
		} else {
			toolStatuses[i].status = "done"
			toolStatuses[i].emoji = "✓"
			successCount++
		}

		// Print final dashboard state
		printInstallDashboard(groupName, toolStatuses, i+1, len(tools))
	}

	return reportGroupInstallationResults(groupName, successCount, len(tools), installErrors)
}

// printInstallDashboard displays the current installation progress
func printInstallDashboard(groupName string, statuses []toolStatus, current, total int) {
	var content strings.Builder
	content.WriteString("\n")

	// Show each tool with its status
	for i, status := range statuses {
		var statusText string
		switch status.status {
		case "done":
			statusText = fmt.Sprintf("%-20s %s %-15s", status.name, status.emoji, "Installed")
		case "failed":
			statusText = fmt.Sprintf("%-20s %s %-15s", status.name, status.emoji, "Failed")
		case "installing":
			statusText = fmt.Sprintf("%-20s %s %-15s", status.name, status.emoji, "Installing...")
		default:
			statusText = fmt.Sprintf("%-20s %s %-15s", status.name, status.emoji, "Pending")
		}

		content.WriteString(fmt.Sprintf("  [%d/%d] %s\n", i+1, total, statusText))
	}

	content.WriteString("\n")

	// Calculate progress
	percentage := (current * 100) / total
	barWidth := 30
	filled := (percentage * barWidth) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	content.WriteString(fmt.Sprintf("  Progress: %d%% %s\n", percentage, bar))

	// Clear previous output and print new dashboard
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top
	fmt.Println(charm.RenderBox(fmt.Sprintf("Installing '%s' group (%d tools)", groupName, total), content.String(), "#00D9FF", false))
}

// installIndividualApp installs a single application using unified installation logic
func installIndividualApp(appName string, dryRun bool, cmd *cobra.Command) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Installing '%s'", appName))

	// Validate app name is not empty
	if appName == "" {
		return errors.NewInstallationError(constants.OpInstall, appName,
			fmt.Errorf("application name cannot be empty"))
	}

	wasNewlyInstalled, err := installSingleToolUnified(appName, dryRun)
	if err != nil {
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

	// Install the tool via brew (availability already checked by caller)
	if err := brew.InstallPackageDirectly(toolName); err != nil {
		return err
	}

	// Handle special cases for specific tools
	if toolName == "zsh" {
		spinner := charm.NewLineSpinner("Installing Oh My Zsh")
		spinner.Start()
		ohMyZshScript := `sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`
		if err := runPostInstallScript(ohMyZshScript); err != nil {
			spinner.Warning("Oh My Zsh setup skipped")
			o.PrintWarning("Post-install script failed for %s: %v", toolName, err)
		} else {
			spinner.Success("Oh My Zsh installed successfully")
		}
	}

	// Handle config check for git
	if toolName == "git" {
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
		return false, err
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
		return nil
	}

	spinner := charm.NewDotsSpinner(fmt.Sprintf("Tracking %s in settings", appName))
	spinner.Start()

	if err := config.AddInstalledApp(appName); err != nil {
		spinner.Warning("Failed to update settings")
		o.PrintWarning("Failed to update settings file: %v", err)
		return nil // Don't fail installation for tracking issues
	}

	spinner.Success(fmt.Sprintf("%s tracked in settings", appName))
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

// loadAndPrepareAppData loads all application data and prepares it for rendering
func loadAndPrepareAppData() (groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string, err error) {
	// Load groups from config
	groups, err = config.GetAvailableGroups()
	if err != nil {
		err = errors.NewConfigurationError(constants.OpInstall, "load-data",
			fmt.Errorf("failed to load groups: %w", err))
		return
	}

	// Get built-in group names
	builtInGroupNames = config.GetBuiltInGroups()

	// Extract and sort custom group names
	for groupName := range groups {
		if !config.IsBuiltInGroup(groupName) {
			customGroupNames = append(customGroupNames, groupName)
		}
	}
	sort.Strings(customGroupNames)

	// Load and sort installed apps
	installedApps, err = config.GetInstalledApps()
	if err != nil {
		// Don't fail on installed apps error, just log warning
		getOutputHandler().PrintWarning("Failed to load installed apps: %v", err)
		installedApps = []string{}
		err = nil // Reset error since we can continue
	} else {
		sort.Strings(installedApps)
	}

	return
}

// Color helper functions for consistent formatting
func colorSectionHeader(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorBold+palantir.ColorCyan, text, palantir.ColorReset)
}

func colorBoldText(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorBold, text, palantir.ColorReset)
}

func colorAppName(text string) string {
	return fmt.Sprintf("%s%s%s", palantir.ColorGreen, text, palantir.ColorReset)
}

func colorGroupNameWithIcon(text string) string {
	return fmt.Sprintf("%s 📁", colorBoldText(text))
}

// renderListView renders applications in a flat list format
func renderListView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
	var content strings.Builder
	content.WriteString("\n")

	// Show built-in groups first
	content.WriteString(colorSectionHeader("Built-in Groups") + "\n\n")
	for _, groupName := range builtInGroupNames {
		if tools, exists := groups[groupName]; exists {
			content.WriteString(fmt.Sprintf("  %s  %s\n", colorGroupNameWithIcon(groupName), strings.Join(tools, ", ")))
		}
	}

	// Show custom groups
	if len(customGroupNames) > 0 {
		content.WriteString("\n" + colorSectionHeader("Custom Groups") + "\n\n")
		for _, groupName := range customGroupNames {
			content.WriteString(fmt.Sprintf("  %s  %s\n", colorGroupNameWithIcon(groupName), strings.Join(groups[groupName], ", ")))
		}
	} else {
		content.WriteString(fmt.Sprintf("\n%sNo custom groups defined%s\n", palantir.ColorBold+palantir.ColorYellow, palantir.ColorReset))
		content.WriteString(fmt.Sprintf("  Add custom groups in ~/%s/%s\n", constants.AnvilConfigDir, constants.ConfigFileName))
	}

	// Show individually tracked installed apps
	if len(installedApps) > 0 {
		content.WriteString("\n" + colorSectionHeader("Individually Tracked Apps") + "\n\n")
		for _, app := range installedApps {
			content.WriteString(fmt.Sprintf("  %s\n", colorAppName(app)))
		}
	}

	content.WriteString("\n")
	return content.String()
}

// AppTreeNode represents a node in the applications tree
type AppTreeNode struct {
	Name     string
	IsGroup  bool
	Apps     []string
	Children []*AppTreeNode
}

// renderTreeView renders applications in a hierarchical tree format
func renderTreeView(groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string) string {
	// Create root node
	root := &AppTreeNode{
		Name:     "Applications",
		IsGroup:  false,
		Children: []*AppTreeNode{},
	}

	// Add built-in groups section
	if len(builtInGroupNames) > 0 {
		builtInNode := &AppTreeNode{
			Name:     "Built-in Groups",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, groupName := range builtInGroupNames {
			if tools, exists := groups[groupName]; exists {
				groupNode := &AppTreeNode{
					Name:    groupName,
					IsGroup: true,
					Apps:    tools,
				}
				builtInNode.Children = append(builtInNode.Children, groupNode)
			}
		}

		if len(builtInNode.Children) > 0 {
			root.Children = append(root.Children, builtInNode)
		}
	}

	// Add custom groups section
	if len(customGroupNames) > 0 {
		customNode := &AppTreeNode{
			Name:     "Custom Groups",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, groupName := range customGroupNames {
			groupNode := &AppTreeNode{
				Name:    groupName,
				IsGroup: true,
				Apps:    groups[groupName],
			}
			customNode.Children = append(customNode.Children, groupNode)
		}

		root.Children = append(root.Children, customNode)
	}

	// Add individually tracked apps section
	if len(installedApps) > 0 {
		individualNode := &AppTreeNode{
			Name:     "Individually Tracked Apps",
			IsGroup:  false,
			Children: []*AppTreeNode{},
		}

		for _, appName := range installedApps {
			appNode := &AppTreeNode{
				Name:    appName,
				IsGroup: false,
			}
			individualNode.Children = append(individualNode.Children, appNode)
		}

		root.Children = append(root.Children, individualNode)
	}

	// Build tree content
	var content strings.Builder
	content.WriteString("\n")
	buildTreeString(&content, root, "", true, true)
	content.WriteString("\n")

	return content.String()
}

// buildTreeString writes an app tree node to a string builder with ASCII art and colors
func buildTreeString(builder *strings.Builder, node *AppTreeNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		// Choose the appropriate tree character
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}

		// Color the output based on node type
		var coloredName string
		if node.IsGroup {
			// Groups are colored in bold blue
			coloredName = fmt.Sprintf("%s%s%s 📁 %s", palantir.ColorBold, palantir.ColorBlue, node.Name, palantir.ColorReset)
		} else if len(node.Children) > 0 {
			// Category headers (Built-in Groups, Custom Groups, etc.) in bold cyan
			coloredName = fmt.Sprintf("%s%s%s%s", palantir.ColorBold, palantir.ColorCyan, node.Name, palantir.ColorReset)
		} else {
			// Individual apps in green
			coloredName = fmt.Sprintf("%s%s%s", palantir.ColorGreen, node.Name, palantir.ColorReset)
		}

		// Write the current node
		builder.WriteString(fmt.Sprintf("%s%s%s\n", prefix, treeChar, coloredName))
	}

	// Write apps within a group
	if node.IsGroup && len(node.Apps) > 0 {
		for i, app := range node.Apps {
			isAppLast := i == len(node.Apps)-1

			// Calculate prefix for app
			var appPrefix string
			if isLast {
				appPrefix = prefix + "    "
			} else {
				appPrefix = prefix + "│   "
			}

			var appTreeChar string
			if isAppLast {
				appTreeChar = "└── "
			} else {
				appTreeChar = "├── "
			}

			// Color individual apps in green
			coloredApp := fmt.Sprintf("%s%s%s", palantir.ColorGreen, app, palantir.ColorReset)
			builder.WriteString(fmt.Sprintf("%s%s%s\n", appPrefix, appTreeChar, coloredApp))
		}
	}

	// Write children
	if node.Children != nil {
		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1

			// Calculate prefix for child
			var childPrefix string
			if isRoot {
				childPrefix = ""
			} else {
				if isLast {
					childPrefix = prefix + "    "
				} else {
					childPrefix = prefix + "│   "
				}
			}

			buildTreeString(builder, child, childPrefix, isChildLast, false)
		}
	}
}

func init() {
	// Add flags for additional functionality
	InstallCmd.Flags().Bool("dry-run", false, "Show what would be installed without installing")
	InstallCmd.Flags().Bool("list", false, "List all available groups")
	InstallCmd.Flags().Bool("tree", false, "Display all applications in a tree format")
	InstallCmd.Flags().Bool("update", false, "Update Homebrew before installation")
	InstallCmd.Flags().String("group-name", "", "Add the installed app to a group (creates group if it doesn't exist)")

	// Add concurrent installation flags
	InstallCmd.Flags().Bool("concurrent", false, "Enable concurrent installation for improved performance")
	InstallCmd.Flags().Int("workers", 0, "Number of concurrent workers (default: number of CPU cores)")
	InstallCmd.Flags().Duration("timeout", 0, "Timeout for individual tool installations (default: 10 minutes)")
}
