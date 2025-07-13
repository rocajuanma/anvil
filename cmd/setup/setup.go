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

// SetupFlags encapsulates all setup command flags
type SetupFlags struct {
	Git      bool
	Zsh      bool
	Iterm2   bool
	Vscode   bool
	Slack    bool
	Chrome   bool
	Password bool
	List     bool
	DryRun   bool
}

// InstallConfig represents the configuration for installing a tool
type InstallConfig struct {
	PackageName  string
	PreCheck     func() bool
	PostInstall  func() error
	SkipIfExists bool
	Description  string
}

// installWithConfig installs a tool based on the provided configuration
func installWithConfig(config InstallConfig) error {
	if config.SkipIfExists && config.PreCheck() {
		return nil
	}

	if err := brew.InstallPackage(config.PackageName); err != nil {
		return constants.NewAnvilError(constants.OpSetup, config.Description, err)
	}

	if config.PostInstall != nil {
		if err := config.PostInstall(); err != nil {
			return constants.NewAnvilError(constants.OpSetup, config.Description, err)
		}
	}
	return nil
}

// runSetupCommand executes the setup process for groups or individual tools
func runSetupCommand(cmd *cobra.Command, args []string) {
	// Check if we're on macOS (required for most installations)
	if runtime.GOOS != "darwin" {
		terminal.PrintWarning("Setup command is currently optimized for macOS. Some features may not work on other platforms.")
	}

	// Extract flags from command
	flags := &SetupFlags{
		Git:      cmd.Flag("git").Changed,
		Zsh:      cmd.Flag("zsh").Changed,
		Iterm2:   cmd.Flag("iterm2").Changed,
		Vscode:   cmd.Flag("vscode").Changed,
		Slack:    cmd.Flag("slack").Changed,
		Chrome:   cmd.Flag("chrome").Changed,
		Password: cmd.Flag("1password").Changed,
		List:     cmd.Flag("list").Changed,
		DryRun:   cmd.Flag("dry-run").Changed,
	}

	// Handle list flag
	if flags.List {
		listAvailableGroups()
		return
	}

	// Handle dry run flag
	if flags.DryRun {
		terminal.PrintInfo("Dry run mode - no actual installations will be performed")
	}

	// Check individual tool flags
	if hasIndividualToolFlags(flags) {
		if err := runIndividualToolSetup(flags); err != nil {
			terminal.PrintError("Individual tool setup failed: %v", err)
			os.Exit(1)
		}
		return
	}

	// Handle group installation
	if len(args) == 0 {
		showUsageAndGroups()
		return
	}

	groupName := args[0]
	if err := runGroupSetup(groupName, flags); err != nil {
		terminal.PrintError("Group setup failed: %v", err)
		os.Exit(1)
	}
}

// hasIndividualToolFlags checks if any individual tool flags are set
func hasIndividualToolFlags(flags *SetupFlags) bool {
	return flags.Git || flags.Zsh || flags.Iterm2 || flags.Vscode || flags.Slack || flags.Chrome || flags.Password
}

// runIndividualToolSetup installs individual tools based on flags
func runIndividualToolSetup(flags *SetupFlags) error {
	terminal.PrintHeader("Individual Tool Setup")

	var toolsToInstall []string

	if flags.Git {
		toolsToInstall = append(toolsToInstall, "git")
	}
	if flags.Zsh {
		toolsToInstall = append(toolsToInstall, "zsh")
	}
	if flags.Iterm2 {
		toolsToInstall = append(toolsToInstall, "iterm2")
	}
	if flags.Vscode {
		toolsToInstall = append(toolsToInstall, "vscode")
	}
	if flags.Slack {
		toolsToInstall = append(toolsToInstall, "slack")
	}
	if flags.Chrome {
		toolsToInstall = append(toolsToInstall, "chrome")
	}
	if flags.Password {
		toolsToInstall = append(toolsToInstall, "1password")
	}

	if len(toolsToInstall) == 0 {
		return constants.NewAnvilError(constants.OpSetup, "individual-tools", fmt.Errorf("no tools specified for installation"))
	}

	terminal.PrintInfo("Installing individual tools: %s", strings.Join(toolsToInstall, ", "))

	var installErrors []string
	successCount := 0
	for i, tool := range toolsToInstall {
		terminal.PrintProgress(i+1, len(toolsToInstall), fmt.Sprintf("Installing %s", tool))

		if flags.DryRun {
			terminal.PrintInfo("Would install: %s", tool)
			successCount++
			continue
		}

		if err := installTool(tool); err != nil {
			installErrors = append(installErrors, fmt.Sprintf("%s: %v", tool, err))
			terminal.PrintError("Failed to install %s: %v", tool, err)
			continue
		}

		terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", tool))
		successCount++
	}

	terminal.PrintHeader("Individual Tool Setup Complete!")

	if len(installErrors) > 0 {
		return constants.NewAnvilError(constants.OpSetup, "individual-tools", fmt.Errorf("failed to install tools: %s", strings.Join(installErrors, ", ")))
	}

	return nil
}

// runGroupSetup installs tools for a specific group
func runGroupSetup(groupName string, flags *SetupFlags) error {
	terminal.PrintHeader(fmt.Sprintf("Setting up '%s' group", groupName))

	// Get tools for the group
	tools, err := config.GetGroupTools(groupName)
	if err != nil {
		return constants.NewAnvilError(constants.OpSetup, groupName, err)
	}

	if len(tools) == 0 {
		return constants.NewAnvilError(constants.OpSetup, groupName, fmt.Errorf("no tools configured for group"))
	}

	terminal.PrintInfo("Installing tools for group '%s': %s", groupName, strings.Join(tools, ", "))

	// Install each tool in the group
	successCount := 0
	var installErrors []string
	for i, tool := range tools {
		terminal.PrintProgress(i+1, len(tools), fmt.Sprintf("Installing %s", tool))

		if flags.DryRun {
			terminal.PrintInfo("Would install: %s", tool)
			successCount++
			continue
		}

		if err := installTool(tool); err != nil {
			installErrors = append(installErrors, fmt.Sprintf("%s: %v", tool, err))
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
		return constants.NewAnvilError(constants.OpSetup, groupName, fmt.Errorf("failed to install %d tools: %s", len(tools)-successCount, strings.Join(installErrors, ", ")))
	}

	return nil
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
		if err := brew.InstallPackage(toolName); err != nil {
			return constants.NewAnvilError(constants.OpSetup, toolName, err)
		}
		return nil
	}
}

// installGit installs and configures Git
func installGit() error {
	return installWithConfig(InstallConfig{
		PackageName:  "git",
		PreCheck:     func() bool { return system.CommandExists("git") },
		SkipIfExists: true,
		Description:  "Git",
	})
}

// installZsh installs and configures Zsh
func installZsh() error {
	return installWithConfig(InstallConfig{
		PackageName: "zsh",
		PostInstall: func() error {
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
		},
		Description: "Zsh with oh-my-zsh",
	})
}

// installIterm2 installs iTerm2
func installIterm2() error {
	return installWithConfig(InstallConfig{
		PackageName: "iterm2",
		Description: "iTerm2",
	})
}

// installVSCode installs Visual Studio Code
func installVSCode() error {
	return installWithConfig(InstallConfig{
		PackageName: "visual-studio-code",
		Description: "Visual Studio Code",
	})
}

// installSlack installs Slack
func installSlack() error {
	return installWithConfig(InstallConfig{
		PackageName: "slack",
		Description: "Slack",
	})
}

// installChrome installs Google Chrome
func installChrome() error {
	return installWithConfig(InstallConfig{
		PackageName: "google-chrome",
		Description: "Google Chrome",
	})
}

// install1Password installs 1Password
func install1Password() error {
	return installWithConfig(InstallConfig{
		PackageName: "1password",
		Description: "1Password",
	})
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
	SetupCmd.Flags().Bool("git", false, "Install and configure Git")
	SetupCmd.Flags().Bool("zsh", false, "Install Zsh with oh-my-zsh configuration")
	SetupCmd.Flags().Bool("iterm2", false, "Install iTerm2 terminal emulator")
	SetupCmd.Flags().Bool("vscode", false, "Install Visual Studio Code")
	SetupCmd.Flags().Bool("slack", false, "Install Slack communication app")
	SetupCmd.Flags().Bool("chrome", false, "Install Google Chrome browser")
	SetupCmd.Flags().Bool("1password", false, "Install 1Password password manager")

	// Utility flags
	SetupCmd.Flags().Bool("list", false, "List all available groups and tools")
	SetupCmd.Flags().Bool("dry-run", false, "Show what would be installed without actually installing")
}
