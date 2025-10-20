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

package show

import (
	"fmt"
	"strings"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/utils"
)

// showAnvilSettingsSection displays specific sections of the anvil settings
func showAnvilSettingsSection(showGroups, showConfigs, showGit, showGitHub bool) error {
	o := getOutputHandler()

	// Stage 1: Locate settings file
	configPath := config.GetConfigPath()

	// Check settings file
	err := checkSettingsFileExists(o, configPath)
	if err != nil {
		return err
	}

	// Stage 2: Load configuration
	anvilConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewFileSystemError(constants.OpShow, "load-config", err)
	}

	// Stage 3: Display requested sections
	switch {
	case showGroups:
		if err := showGroupsSection(anvilConfig); err != nil {
			return err
		}
	case showConfigs:
		if err := showConfigsSection(anvilConfig); err != nil {
			return err
		}
	case showGit:
		if err := showGitSection(anvilConfig); err != nil {
			return err
		}
	case showGitHub:
		if err := showGitHubSection(anvilConfig); err != nil {
			return err
		}
	}

	return nil
}

// showGroups displays the groups section
func showGroupsSection(anvilConfig *config.AnvilConfig) error {
	var boxContent strings.Builder
	boxContent.WriteString("\n")

	// Get built-in and custom groups
	builtInGroups := config.GetBuiltInGroups()
	var customGroups []string
	for groupName := range anvilConfig.Groups {
		if !config.IsBuiltInGroup(groupName) {
			customGroups = append(customGroups, groupName)
		}
	}

	// Display built-in groups
	if len(builtInGroups) > 0 {
		boxContent.WriteString("  Built-in Groups:\n")
		for _, groupName := range builtInGroups {
			if tools, exists := anvilConfig.Groups[groupName]; exists {
				boxContent.WriteString(fmt.Sprintf("    📁 %s (%d tools)\n", groupName, len(tools)))
				for _, tool := range tools {
					boxContent.WriteString(fmt.Sprintf("      • %s\n", tool))
				}
			}
		}
		boxContent.WriteString("\n")
	}

	// Display custom groups
	if len(customGroups) > 0 {
		boxContent.WriteString("  Custom Groups:\n")
		for _, groupName := range customGroups {
			if tools, exists := anvilConfig.Groups[groupName]; exists {
				boxContent.WriteString(fmt.Sprintf("    📁 %s (%d tools)\n", groupName, len(tools)))
				for _, tool := range tools {
					boxContent.WriteString(fmt.Sprintf("      • %s\n", tool))
				}
			}
		}
		boxContent.WriteString("\n")
	}

	// Display in box
	fmt.Println(charm.RenderBox("Groups", boxContent.String(), "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Use 'anvil install --list' to see all available groups")
	fmt.Println("  💡 Use 'anvil install --tree' to see groups in tree format")
	fmt.Println()

	return nil
}

// showConfigsSection displays the configs section
func showConfigsSection(anvilConfig *config.AnvilConfig) error {
	var boxContent strings.Builder
	boxContent.WriteString("\n")

	if len(anvilConfig.Configs) == 0 {
		boxContent.WriteString("  No configured source directories found.\n")
		boxContent.WriteString("  Use 'anvil config push <app-name> <path>' to configure source directories.\n")
	} else {
		boxContent.WriteString("  Configured Source Directories:\n")
		for appName, path := range anvilConfig.Configs {
			boxContent.WriteString(fmt.Sprintf("    %s: %s\n", utils.ColorAppName(appName), path))
		}
	}

	boxContent.WriteString("\n")

	// Display in box
	fmt.Println(charm.RenderBox("Config Sources", boxContent.String(), "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Use 'anvil config push <app-name> <path>' to add source directories")
	fmt.Println("  💡 Use 'anvil config pull <app-name>' to pull configurations")
	fmt.Println()

	return nil
}

// showGitSection displays the git configuration section
func showGitSection(anvilConfig *config.AnvilConfig) error {
	var boxContent strings.Builder
	boxContent.WriteString("\n")

	boxContent.WriteString(fmt.Sprintf("    Username: %s\n", utils.BoldText(anvilConfig.Git.Username, "")))
	boxContent.WriteString(fmt.Sprintf("    Email: %s\n", utils.BoldText(anvilConfig.Git.Email, "")))
	if anvilConfig.Git.SSHKeyPath != "" {
		boxContent.WriteString(fmt.Sprintf("    SSH Key Path: %s\n", utils.BoldText(anvilConfig.Git.SSHKeyPath, "")))
	}

	boxContent.WriteString("\n")

	// Display in box
	fmt.Println(charm.RenderBox("Git Configuration", boxContent.String(), "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Git configuration is auto-populated from your local git settings")
	fmt.Println("  💡 SSH key path is auto-detected from common locations")
	fmt.Println()

	return nil
}

// showGitHubSection displays the GitHub configuration section
func showGitHubSection(anvilConfig *config.AnvilConfig) error {
	var boxContent strings.Builder
	boxContent.WriteString("\n")

	boxContent.WriteString(fmt.Sprintf("    Repository: %s\n", utils.BoldText(anvilConfig.GitHub.ConfigRepo, "")))
	boxContent.WriteString(fmt.Sprintf("    Branch: %s\n", utils.BoldText(anvilConfig.GitHub.Branch, "")))
	boxContent.WriteString(fmt.Sprintf("    Local Path: %s\n", utils.BoldText(anvilConfig.GitHub.LocalPath, "")))
	if anvilConfig.GitHub.Token != "" {
		boxContent.WriteString(fmt.Sprintf("    Token: %s\n", utils.BoldText(anvilConfig.GitHub.Token, "")))
	}
	if anvilConfig.GitHub.TokenEnvVar != "" {
		boxContent.WriteString(fmt.Sprintf("    Token Environment Variable: %s\n", utils.BoldText(anvilConfig.GitHub.TokenEnvVar, "")))
	}

	boxContent.WriteString("\n")

	// Display in box
	fmt.Println(charm.RenderBox("GitHub Configuration", boxContent.String(), "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Use 'anvil config import' to import configurations from GitHub")
	fmt.Println("  💡 Use 'anvil config sync' to sync with your GitHub repository")
	fmt.Println()

	return nil
}
