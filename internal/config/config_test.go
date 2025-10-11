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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rocajuanma/anvil/internal/constants"
)

// setupTestConfig creates a test configuration with temporary directories
func setupTestConfig(t *testing.T) (string, func()) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Mock the config path
	originalHome := os.Getenv("HOME")
	cleanup := func() {
		os.Setenv("HOME", originalHome)
	}

	os.Setenv("HOME", tempDir)

	// Create a minimal test config instead of relying on sample file
	config := createTestConfig()

	// Create the necessary directories
	err := CreateDirectories()
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Save the initial config
	err = SaveConfig(config)
	if err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	return tempDir, cleanup
}

// createTestConfig creates a minimal test configuration
func createTestConfig() *AnvilConfig {
	return &AnvilConfig{
		Version: "2.0.0",
		Tools: AnvilTools{
			RequiredTools: []string{constants.PkgGit, constants.CurlCommand},
			OptionalTools: []string{constants.BrewCommand, constants.PkgDocker, constants.PkgKubectl},
			InstalledApps: []string{},
		},
		Groups: AnvilGroups{
			"dev":        {constants.PkgGit, constants.PkgZsh, constants.PkgIterm2, constants.PkgVSCode},
			"essentials": {constants.PkgSlack, constants.PkgChrome, constants.Pkg1Password},
		},
		Configs: make(map[string]string),
		Git: GitConfig{
			Username:   "Test User",
			Email:      "test@example.com",
			SSHKeyPath: "/tmp/test_ssh_key",
		},
		GitHub: GitHubConfig{
			ConfigRepo:  "",
			Branch:      "main",
			LocalPath:   "/tmp/test_dotfiles",
			TokenEnvVar: "GITHUB_TOKEN",
		},
		ToolConfigs: AnvilToolConfigs{
			Tools: map[string]ToolInstallConfig{
				constants.PkgZsh: {
					PostInstallScript: constants.OhMyZshInstallCmd,
					ConfigCheck:       false,
					Dependencies:      []string{},
				},
				constants.PkgGit: {
					ConfigCheck:  true,
					Dependencies: []string{},
				},
			},
		},
	}
}

func TestAddInstalledApp(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test adding a new app
	testApp := "test-app"
	err := AddInstalledApp(testApp)
	if err != nil {
		t.Fatalf("Failed to add installed app: %v", err)
	}

	// Load config and verify the app was added
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 1 {
		t.Errorf("Expected 1 installed app, got %d", len(updatedConfig.Tools.InstalledApps))
	}

	if updatedConfig.Tools.InstalledApps[0] != testApp {
		t.Errorf("Expected app '%s', got '%s'", testApp, updatedConfig.Tools.InstalledApps[0])
	}
}

func TestAddInstalledAppDuplicate(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add an existing app to the config
	err := AddInstalledApp("existing-app")
	if err != nil {
		t.Fatalf("Failed to add existing app: %v", err)
	}

	// Test adding the same app again
	err = AddInstalledApp("existing-app")
	if err != nil {
		t.Fatalf("Failed to add installed app: %v", err)
	}

	// Load config and verify no duplicate was added
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 1 {
		t.Errorf("Expected 1 installed app (no duplicate), got %d", len(updatedConfig.Tools.InstalledApps))
	}
}

func TestAddInstalledAppSkipRequired(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test adding a required tool - should not be added to installed apps
	err := AddInstalledApp(constants.PkgGit)
	if err != nil {
		t.Fatalf("Failed to add required tool: %v", err)
	}

	// Load config and verify the required tool was not added to installed apps
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 0 {
		t.Errorf("Expected 0 installed apps (required tool should not be tracked), got %d", len(updatedConfig.Tools.InstalledApps))
	}
}

func TestGetInstalledApps(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add test apps
	testApps := []string{"app1", "app2", "app3"}
	for _, app := range testApps {
		err := AddInstalledApp(app)
		if err != nil {
			t.Fatalf("Failed to add test app %s: %v", app, err)
		}
	}

	// Test getting installed apps
	installedApps, err := GetInstalledApps()
	if err != nil {
		t.Fatalf("Failed to get installed apps: %v", err)
	}

	if len(installedApps) != 3 {
		t.Errorf("Expected 3 installed apps, got %d", len(installedApps))
	}

	expectedApps := []string{"app1", "app2", "app3"}
	for i, app := range installedApps {
		if app != expectedApps[i] {
			t.Errorf("Expected app '%s' at index %d, got '%s'", expectedApps[i], i, app)
		}
	}
}

func TestIsAppTracked(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add a tracked app
	err := AddInstalledApp("tracked-app")
	if err != nil {
		t.Fatalf("Failed to add tracked app: %v", err)
	}

	// Test checking if app is tracked
	isTracked, err := IsAppTracked("tracked-app")
	if err != nil {
		t.Fatalf("Failed to check if app is tracked: %v", err)
	}

	if !isTracked {
		t.Error("Expected app to be tracked")
	}

	// Test checking if app is not tracked
	isTracked, err = IsAppTracked("untracked-app")
	if err != nil {
		t.Fatalf("Failed to check if app is tracked: %v", err)
	}

	if isTracked {
		t.Error("Expected app to not be tracked")
	}

	// Test checking if required tool is tracked
	isTracked, err = IsAppTracked(constants.PkgGit)
	if err != nil {
		t.Fatalf("Failed to check if required tool is tracked: %v", err)
	}

	if !isTracked {
		t.Error("Expected required tool to be tracked")
	}
}

func TestRemoveInstalledApp(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add test apps
	testApps := []string{"app1", "app2", "app3"}
	for _, app := range testApps {
		err := AddInstalledApp(app)
		if err != nil {
			t.Fatalf("Failed to add test app %s: %v", app, err)
		}
	}

	// Test removing an app
	err := RemoveInstalledApp("app2")
	if err != nil {
		t.Fatalf("Failed to remove installed app: %v", err)
	}

	// Load config and verify the app was removed
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 2 {
		t.Errorf("Expected 2 installed apps after removal, got %d", len(updatedConfig.Tools.InstalledApps))
	}

	expectedApps := []string{"app1", "app3"}
	for i, app := range updatedConfig.Tools.InstalledApps {
		if app != expectedApps[i] {
			t.Errorf("Expected app '%s' at index %d, got '%s'", expectedApps[i], i, app)
		}
	}
}

func TestAddAppToGroup(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test adding app to a new group
	groupName := "test-group"
	appName := "test-app"

	err := AddAppToGroup(groupName, appName)
	if err != nil {
		t.Fatalf("Failed to add app to group: %v", err)
	}

	// Load config and verify the group was created with the app
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	// Default config has 2 groups ("dev" and "essentials"), so we expect 3 total
	if len(updatedConfig.Groups) != 3 {
		t.Errorf("Expected 3 groups (2 default + 1 new), got %d", len(updatedConfig.Groups))
	}

	if tools, exists := updatedConfig.Groups[groupName]; !exists {
		t.Errorf("Expected group '%s' to exist", groupName)
	} else if len(tools) != 1 {
		t.Errorf("Expected 1 tool in group, got %d", len(tools))
	} else if tools[0] != appName {
		t.Errorf("Expected app '%s' in group, got '%s'", appName, tools[0])
	}
}

func TestAddAppToExistingGroup(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a group with an initial app
	groupName := "existing-group"
	initialApp := "initial-app"

	err := AddAppToGroup(groupName, initialApp)
	if err != nil {
		t.Fatalf("Failed to create initial group: %v", err)
	}

	// Add another app to the existing group
	newApp := "new-app"
	err = AddAppToGroup(groupName, newApp)
	if err != nil {
		t.Fatalf("Failed to add app to existing group: %v", err)
	}

	// Load config and verify both apps are in the group
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if tools, exists := updatedConfig.Groups[groupName]; !exists {
		t.Errorf("Expected group '%s' to exist", groupName)
	} else if len(tools) != 2 {
		t.Errorf("Expected 2 tools in group, got %d", len(tools))
	} else {
		// Check that both apps are present
		foundInitial := false
		foundNew := false
		for _, tool := range tools {
			if tool == initialApp {
				foundInitial = true
			}
			if tool == newApp {
				foundNew = true
			}
		}
		if !foundInitial {
			t.Errorf("Expected to find '%s' in group", initialApp)
		}
		if !foundNew {
			t.Errorf("Expected to find '%s' in group", newApp)
		}
	}
}

func TestAddAppToGroupDuplicate(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create a group with an app
	groupName := "duplicate-group"
	appName := "duplicate-app"

	err := AddAppToGroup(groupName, appName)
	if err != nil {
		t.Fatalf("Failed to create initial group: %v", err)
	}

	// Try to add the same app again
	err = AddAppToGroup(groupName, appName)
	if err != nil {
		t.Fatalf("Failed to add duplicate app to group: %v", err)
	}

	// Load config and verify no duplicate was added
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if tools, exists := updatedConfig.Groups[groupName]; !exists {
		t.Errorf("Expected group '%s' to exist", groupName)
	} else if len(tools) != 1 {
		t.Errorf("Expected 1 tool in group (no duplicates), got %d", len(tools))
	} else if tools[0] != appName {
		t.Errorf("Expected app '%s' in group, got '%s'", appName, tools[0])
	}
}

// Test helper functions and DRY improvements
func TestHelperFunctions(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test GetConfigDirectory
	configDir := GetConfigDirectory()
	if configDir == "" {
		t.Error("Expected GetConfigDirectory to return a non-empty string")
	}

	// Test GetConfigPath
	configPath := GetConfigPath()
	if configPath == "" {
		t.Error("Expected GetConfigPath to return a non-empty string")
	}

	// Test GetBuiltInGroups
	builtInGroups := GetBuiltInGroups()
	expectedGroups := []string{"dev", "essentials"}
	if len(builtInGroups) != len(expectedGroups) {
		t.Errorf("Expected %d built-in groups, got %d", len(expectedGroups), len(builtInGroups))
	}

	// Test IsBuiltInGroup
	if !IsBuiltInGroup("dev") {
		t.Error("Expected 'dev' to be a built-in group")
	}
	if !IsBuiltInGroup("essentials") {
		t.Error("Expected 'essentials' to be a built-in group")
	}
	if IsBuiltInGroup("custom-group") {
		t.Error("Expected 'custom-group' to not be a built-in group")
	}
}

// Test tool configuration functions
func TestToolConfigFunctions(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test GetToolConfig for existing tool
	toolConfig, err := GetToolConfig(constants.PkgGit)
	if err != nil {
		t.Fatalf("Failed to get tool config: %v", err)
	}
	if toolConfig == nil {
		t.Error("Expected tool config to not be nil")
	}
	if !toolConfig.ConfigCheck {
		t.Error("Expected git tool config to have ConfigCheck enabled")
	}

	// Test GetToolConfig for non-existing tool (should return default)
	toolConfig, err = GetToolConfig("non-existing-tool")
	if err != nil {
		t.Fatalf("Failed to get tool config for non-existing tool: %v", err)
	}
	if toolConfig == nil {
		t.Error("Expected default tool config to not be nil")
	}
	if toolConfig.ConfigCheck {
		t.Error("Expected default tool config to have ConfigCheck disabled")
	}

	// Test GetToolConfigs for multiple tools
	toolNames := []string{constants.PkgGit, constants.PkgZsh, "non-existing-tool"}
	configs, err := GetToolConfigs(toolNames)
	if err != nil {
		t.Fatalf("Failed to get tool configs: %v", err)
	}
	if len(configs) != len(toolNames) {
		t.Errorf("Expected %d tool configs, got %d", len(toolNames), len(configs))
	}

	// Test SetToolConfig
	newToolConfig := ToolInstallConfig{
		PostInstallScript: "echo 'test'",
		ConfigCheck:       true,
		Dependencies:      []string{"dependency1"},
	}
	err = SetToolConfig("test-tool", newToolConfig)
	if err != nil {
		t.Fatalf("Failed to set tool config: %v", err)
	}

	// Verify the tool config was saved
	savedConfig, err := GetToolConfig("test-tool")
	if err != nil {
		t.Fatalf("Failed to get saved tool config: %v", err)
	}
	if savedConfig.PostInstallScript != newToolConfig.PostInstallScript {
		t.Errorf("Expected PostInstallScript '%s', got '%s'", newToolConfig.PostInstallScript, savedConfig.PostInstallScript)
	}
}

// Test app configuration functions
func TestAppConfigFunctions(t *testing.T) {
	tempDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test SetAppConfigPath
	appName := "test-app"
	configPath := filepath.Join(tempDir, "test-config")

	// Create the directory to avoid path existence check failure
	err := os.MkdirAll(configPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test config directory: %v", err)
	}

	err = SetAppConfigPath(appName, configPath)
	if err != nil {
		t.Fatalf("Failed to set app config path: %v", err)
	}

	// Test GetConfiguredApps
	apps, err := GetConfiguredApps()
	if err != nil {
		t.Fatalf("Failed to get configured apps: %v", err)
	}
	if len(apps) != 1 {
		t.Errorf("Expected 1 configured app, got %d", len(apps))
	}
	if apps[0] != appName {
		t.Errorf("Expected app '%s', got '%s'", appName, apps[0])
	}

	// Test GetAppConfigPath
	retrievedPath, found, err := GetAppConfigPath(appName)
	if err != nil {
		t.Fatalf("Failed to get app config path: %v", err)
	}
	if !found {
		t.Error("Expected app config path to be found")
	}
	if retrievedPath != configPath {
		t.Errorf("Expected config path '%s', got '%s'", configPath, retrievedPath)
	}
}

// Test group management functions
func TestGroupManagementFunctions(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test AddCustomGroup
	groupName := "test-group"
	tools := []string{"tool1", "tool2"}
	err := AddCustomGroup(groupName, tools)
	if err != nil {
		t.Fatalf("Failed to add custom group: %v", err)
	}

	// Verify the group was added
	groupTools, err := GetGroupTools(groupName)
	if err != nil {
		t.Fatalf("Failed to get group tools: %v", err)
	}
	if len(groupTools) != len(tools) {
		t.Errorf("Expected %d tools in group, got %d", len(tools), len(groupTools))
	}

	// Test UpdateGroupTools
	newTools := []string{"tool3", "tool4", "tool5"}
	err = UpdateGroupTools(groupName, newTools)
	if err != nil {
		t.Fatalf("Failed to update group tools: %v", err)
	}

	// Verify the group was updated
	updatedTools, err := GetGroupTools(groupName)
	if err != nil {
		t.Fatalf("Failed to get updated group tools: %v", err)
	}
	if len(updatedTools) != len(newTools) {
		t.Errorf("Expected %d tools in updated group, got %d", len(newTools), len(updatedTools))
	}

	// Test GetAvailableGroups
	groups, err := GetAvailableGroups()
	if err != nil {
		t.Fatalf("Failed to get available groups: %v", err)
	}
	if len(groups) < 3 { // Should have at least 2 built-in + 1 custom
		t.Errorf("Expected at least 3 groups, got %d", len(groups))
	}
}
