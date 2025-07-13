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
	"os"
	"runtime"
	"strings"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// SetupCmd represents the setup command, which installs tools and applications in groups
// This command reads from the Anvil configuration to install predefined sets of tools
var SetupCmd = &cobra.Command{
	Use:   "setup [group]",
	Short: "Install development tools and applications in predefined groups",
	Long:  constants.SETUP_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSetupCommand(cmd, args)
	},
}

// Command flags for individual tool installation
var (
	gitFlag      bool
	zshFlag      bool
	iterm2Flag   bool
	vscodeFlag   bool
	slackFlag    bool
	chromeFlag   bool
	passwordFlag bool
	listFlag     bool
	dryRunFlag   bool
)

// runSetupCommand executes the setup process for groups or individual tools
func runSetupCommand(cmd *cobra.Command, args []string) {
	// Check if we're on macOS (required for most installations)
	if runtime.GOOS != "darwin" {
		terminal.PrintWarning("Setup command is currently optimized for macOS. Some features may not work on other platforms.")
	}

	// Handle list flag
	if listFlag {
		listAvailableGroups()
		return
	}

	// Handle dry run flag
	if dryRunFlag {
		terminal.PrintInfo("Dry run mode - no actual installations will be performed")
	}

	// Check individual tool flags
	if hasIndividualToolFlags(cmd) {
		runIndividualToolSetup(cmd)
		return
	}

	// Handle group installation
	if len(args) == 0 {
		showUsageAndGroups()
		return
	}

	groupName := args[0]
	runGroupSetup(groupName)
}

// hasIndividualToolFlags checks if any individual tool flags are set
func hasIndividualToolFlags(cmd *cobra.Command) bool {
	return gitFlag || zshFlag || iterm2Flag || vscodeFlag || slackFlag || chromeFlag || passwordFlag
}

// runIndividualToolSetup installs individual tools based on flags
func runIndividualToolSetup(cmd *cobra.Command) {
	terminal.PrintHeader("Individual Tool Setup")

	var toolsToInstall []string

	if gitFlag {
		toolsToInstall = append(toolsToInstall, "git")
	}
	if zshFlag {
		toolsToInstall = append(toolsToInstall, "zsh")
	}
	if iterm2Flag {
		toolsToInstall = append(toolsToInstall, "iterm2")
	}
	if vscodeFlag {
		toolsToInstall = append(toolsToInstall, "vscode")
	}
	if slackFlag {
		toolsToInstall = append(toolsToInstall, "slack")
	}
	if chromeFlag {
		toolsToInstall = append(toolsToInstall, "chrome")
	}
	if passwordFlag {
		toolsToInstall = append(toolsToInstall, "1password")
	}

	if len(toolsToInstall) == 0 {
		terminal.PrintError("No tools specified for installation")
		return
	}

	terminal.PrintInfo("Installing individual tools: %s", strings.Join(toolsToInstall, ", "))

	for i, tool := range toolsToInstall {
		terminal.PrintProgress(i+1, len(toolsToInstall), fmt.Sprintf("Installing %s", tool))

		if dryRunFlag {
			terminal.PrintInfo("Would install: %s", tool)
			continue
		}

		if err := installTool(tool); err != nil {
			terminal.PrintError("Failed to install %s: %v", tool, err)
			continue
		}

		terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", tool))
	}

	terminal.PrintHeader("Individual Tool Setup Complete!")
}

// runGroupSetup installs tools for a specific group
func runGroupSetup(groupName string) {
	terminal.PrintHeader(fmt.Sprintf("Setting up '%s' group", groupName))

	// Get tools for the group
	tools, err := config.GetGroupTools(groupName)
	if err != nil {
		terminal.PrintError("Failed to get tools for group '%s': %v", groupName, err)
		terminal.PrintInfo("Use 'anvil setup --list' to see available groups")
		return
	}

	if len(tools) == 0 {
		terminal.PrintWarning("No tools configured for group '%s'", groupName)
		return
	}

	terminal.PrintInfo("Installing tools for group '%s': %s", groupName, strings.Join(tools, ", "))

	// Install each tool in the group
	successCount := 0
	for i, tool := range tools {
		terminal.PrintProgress(i+1, len(tools), fmt.Sprintf("Installing %s", tool))

		if dryRunFlag {
			terminal.PrintInfo("Would install: %s", tool)
			successCount++
			continue
		}

		if err := installTool(tool); err != nil {
			terminal.PrintError("Failed to install %s: %v", tool, err)
			continue
		}

		terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", tool))
		successCount++
	}

	// Print summary
	terminal.PrintHeader("Group Setup Complete!")
	terminal.PrintInfo("Successfully installed %d of %d tools in group '%s'", successCount, len(tools), groupName)

	if successCount < len(tools) {
		terminal.PrintWarning("Some tools failed to install. Check the output above for details.")
	}
}

// installTool installs a specific tool
func installTool(toolName string) error {
	switch toolName {
	case "git":
		return installGit()
	case "zsh":
		return installZsh()
	case "iterm2":
		return installIterm2()
	case "vscode":
		return installVSCode()
	case "slack":
		return installSlack()
	case "chrome":
		return installChrome()
	case "1password":
		return install1Password()
	default:
		// Try to install as a brew package
		return brew.InstallPackage(toolName)
	}
}

// installGit installs and configures Git
func installGit() error {
	if system.CommandExists("git") {
		return nil // Already installed
	}

	return brew.InstallPackage("git")
}

// installZsh installs and configures Zsh
func installZsh() error {
	if err := brew.InstallPackage("zsh"); err != nil {
		return fmt.Errorf("failed to install zsh: %w", err)
	}

	// Install oh-my-zsh if not present
	homeDir, _ := os.UserHomeDir()
	ohmyzshDir := fmt.Sprintf("%s/.oh-my-zsh", homeDir)

	if _, err := os.Stat(ohmyzshDir); os.IsNotExist(err) {
		terminal.PrintInfo("Installing oh-my-zsh...")
		installCmd := `sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`

		result, err := system.RunCommand("sh", "-c", installCmd)
		if err != nil || !result.Success {
			return fmt.Errorf("failed to install oh-my-zsh: %v", err)
		}
	}

	return nil
}

// installIterm2 installs iTerm2
func installIterm2() error {
	return brew.InstallPackage("iterm2")
}

// installVSCode installs Visual Studio Code
func installVSCode() error {
	return brew.InstallPackage("visual-studio-code")
}

// installSlack installs Slack
func installSlack() error {
	return brew.InstallPackage("slack")
}

// installChrome installs Google Chrome
func installChrome() error {
	return brew.InstallPackage("google-chrome")
}

// install1Password installs 1Password
func install1Password() error {
	return brew.InstallPackage("1password")
}

// listAvailableGroups lists all available groups and their tools
func listAvailableGroups() {
	terminal.PrintHeader("Available Setup Groups")

	groups, err := config.GetAvailableGroups()
	if err != nil {
		terminal.PrintError("Failed to load groups: %v", err)
		return
	}

	for groupName, tools := range groups {
		terminal.PrintInfo("Group: %s", groupName)
		for _, tool := range tools {
			terminal.PrintInfo("  • %s", tool)
		}
		terminal.PrintInfo("")
	}

	terminal.PrintInfo("Usage:")
	terminal.PrintInfo("  anvil setup <group>     - Install all tools in a group")
	terminal.PrintInfo("  anvil setup --git       - Install only Git")
	terminal.PrintInfo("  anvil setup --zsh       - Install only Zsh with oh-my-zsh")
	terminal.PrintInfo("  anvil setup --dry-run   - Show what would be installed without installing")
}

// showUsageAndGroups shows usage information and available groups
func showUsageAndGroups() {
	terminal.PrintHeader("Anvil Setup Command")
	terminal.PrintInfo("Usage: anvil setup [group] [flags]")
	terminal.PrintInfo("")
	terminal.PrintInfo("Examples:")
	terminal.PrintInfo("  anvil setup dev         - Install development tools")
	terminal.PrintInfo("  anvil setup new-laptop  - Install new laptop essentials")
	terminal.PrintInfo("  anvil setup --git       - Install only Git")
	terminal.PrintInfo("  anvil setup --list      - List all available groups")
	terminal.PrintInfo("")

	listAvailableGroups()
}

func init() {
	// Individual tool flags
	SetupCmd.Flags().BoolVar(&gitFlag, "git", false, "Install and configure Git")
	SetupCmd.Flags().BoolVar(&zshFlag, "zsh", false, "Install Zsh with oh-my-zsh configuration")
	SetupCmd.Flags().BoolVar(&iterm2Flag, "iterm2", false, "Install iTerm2 terminal emulator")
	SetupCmd.Flags().BoolVar(&vscodeFlag, "vscode", false, "Install Visual Studio Code")
	SetupCmd.Flags().BoolVar(&slackFlag, "slack", false, "Install Slack communication app")
	SetupCmd.Flags().BoolVar(&chromeFlag, "chrome", false, "Install Google Chrome browser")
	SetupCmd.Flags().BoolVar(&passwordFlag, "1password", false, "Install 1Password password manager")

	// Utility flags
	SetupCmd.Flags().BoolVar(&listFlag, "list", false, "List all available groups and tools")
	SetupCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Show what would be installed without actually installing")
}
