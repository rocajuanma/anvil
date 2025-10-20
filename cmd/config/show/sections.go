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
	"sort"
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

// showGroupsSection displays the groups section using shared rendering functions
func showGroupsSection(anvilConfig *config.AnvilConfig) error {
	// Load and prepare data using the same logic as install command
	groups, builtInGroupNames, customGroupNames, installedApps, err := loadAndPrepareAppData()
	if err != nil {
		return err
	}

	// Use the shared tree view renderer
	content := utils.RenderTreeView(groups, builtInGroupNames, customGroupNames, installedApps)

	// Display in box
	fmt.Println(charm.RenderBox("Groups", content, "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Use 'anvil install --list' to see all available groups")
	fmt.Println("  💡 Use 'anvil install --tree' to see groups in tree format")
	fmt.Println()

	return nil
}

// loadAndPrepareAppData loads all application data and prepares it for rendering
// This function is copied from the install package to maintain consistency
func loadAndPrepareAppData() (groups map[string][]string, builtInGroupNames []string, customGroupNames []string, installedApps []string, err error) {
	// Load groups from config
	groups, err = config.GetAvailableGroups()
	if err != nil {
		err = errors.NewConfigurationError(constants.OpShow, "load-data",
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
