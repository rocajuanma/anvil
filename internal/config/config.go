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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/system"
	"github.com/rocajuanma/anvil/internal/utils"
	"gopkg.in/yaml.v2"
)

// Configuration cache to avoid repeated file I/O operations
var (
	configCache      *AnvilConfig
	configCacheMutex sync.RWMutex
)

// ToolInstallConfig represents configuration for tool-specific installation
type ToolInstallConfig struct {
	PostInstallScript string            `yaml:"post_install_script,omitempty"`
	EnvironmentSetup  map[string]string `yaml:"environment_setup,omitempty"`
	ConfigCheck       bool              `yaml:"config_check,omitempty"`
	Dependencies      []string          `yaml:"dependencies,omitempty"`
}

// AnvilGroups represents grouped tool configurations
type AnvilGroups map[string][]string

// AnvilToolConfigs represents tool-specific configurations
type AnvilToolConfigs struct {
	Tools map[string]ToolInstallConfig `yaml:"tools"`
}

// GitConfig represents git configuration
type GitConfig struct {
	Username   string `yaml:"username"`
	Email      string `yaml:"email"`
	SSHKeyPath string `yaml:"ssh_key_path,omitempty"` // Reference to SSH private key
}

// GitHubConfig represents GitHub repository configuration for config sync
type GitHubConfig struct {
	ConfigRepo  string `yaml:"config_repo"`             // GitHub repository URL for configs (e.g., "username/dotfiles")
	Branch      string `yaml:"branch"`                  // Branch to use (default: "main")
	LocalPath   string `yaml:"local_path"`              // Local path where configs are stored/synced
	Token       string `yaml:"token,omitempty"`         // GitHub token (use env var reference)
	TokenEnvVar string `yaml:"token_env_var,omitempty"` // Environment variable name for token
}

// AnvilConfig represents the main anvil configuration
type AnvilConfig struct {
	Version     string            `yaml:"version"`
	Tools       AnvilTools        `yaml:"tools"`
	Groups      AnvilGroups       `yaml:"groups"`
	Configs     map[string]string `yaml:"configs"` // Maps app names to their local config paths
	Git         GitConfig         `yaml:"git"`
	GitHub      GitHubConfig      `yaml:"github"`
	ToolConfigs AnvilToolConfigs  `yaml:"tool_configs,omitempty"`
}

// AnvilTools represents tool configurations
type AnvilTools struct {
	RequiredTools []string `yaml:"required_tools"`
	OptionalTools []string `yaml:"optional_tools"`
	InstalledApps []string `yaml:"installed_apps"` // Tracks individually installed applications
}

// getCachedConfig returns the cached configuration or loads it if not cached
func getCachedConfig() (*AnvilConfig, error) {
	configCacheMutex.RLock()
	if configCache != nil {
		configCacheMutex.RUnlock()
		return configCache, nil
	}
	configCacheMutex.RUnlock()

	configCacheMutex.Lock()
	defer configCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if configCache != nil {
		return configCache, nil
	}

	var err error
	configCache, err = LoadConfig()
	return configCache, err
}

// withConfig executes a function with the cached config, handling common error patterns
func withConfig(fn func(*AnvilConfig) error) error {
	config, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return fn(config)
}

// withConfigAndSave executes a function with the cached config and saves it
func withConfigAndSave(fn func(*AnvilConfig) error) error {
	config, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := fn(config); err != nil {
		return err
	}
	return SaveConfig(config)
}

// ensureGroupsMap ensures the Groups map is initialized
func ensureGroupsMap(config *AnvilConfig) {
	if config.Groups == nil {
		config.Groups = make(map[string][]string)
	}
}

// ensureConfigsMap ensures the Configs map is initialized
func ensureConfigsMap(config *AnvilConfig) {
	if config.Configs == nil {
		config.Configs = make(map[string]string)
	}
}

// ensureToolConfigsMap ensures the ToolConfigs.Tools map is initialized
func ensureToolConfigsMap(config *AnvilConfig) {
	if config.ToolConfigs.Tools == nil {
		config.ToolConfigs.Tools = make(map[string]ToolInstallConfig)
	}
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

// PopulateGitConfigFromSystem populates git configuration from local git settings and auto-detects SSH keys
func PopulateGitConfigFromSystem(gitConfig *GitConfig) error {
	// Always populate username from local git config
	if gitUser, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserName); err == nil && gitUser.Success {
		gitConfig.Username = strings.TrimSpace(gitUser.Output)
	}

	// Always populate email from local git config
	if gitEmail, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserEmail); err == nil && gitEmail.Success {
		gitConfig.Email = strings.TrimSpace(gitEmail.Output)
	}

	// Auto-detect SSH key path from common locations
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")

	// Common SSH key names in order of preference
	commonKeyNames := []string{
		"id_ed25519",
		"id_ed25519_personal",
		"id_rsa",
		"id_rsa_personal",
		"id_ecdsa",
	}

	// Find the first existing SSH key
	for _, keyName := range commonKeyNames {
		keyPath := filepath.Join(sshDir, keyName)
		if _, err := os.Stat(keyPath); err == nil {
			gitConfig.SSHKeyPath = keyPath
			break
		}
	}

	// If no common keys found, use the default path (will be created if needed)
	if gitConfig.SSHKeyPath == "" {
		gitConfig.SSHKeyPath = filepath.Join(sshDir, "id_ed25519")
	}

	return nil
}

// invalidateCache clears the configuration cache
func invalidateCache() {
	configCacheMutex.Lock()
	defer configCacheMutex.Unlock()
	configCache = nil
}

// LoadSampleConfig loads the sample configuration from the assets file
func LoadSampleConfig() (*AnvilConfig, error) {
	samplePath := getSampleConfigPath()
	sampleData, err := os.ReadFile(samplePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sample config file: %w", err)
	}

	var config AnvilConfig
	if err := yaml.Unmarshal(sampleData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sample config: %w", err)
	}

	// Set dynamic paths
	homeDir, _ := os.UserHomeDir()
	config.GitHub.LocalPath = filepath.Join(homeDir, constants.AnvilConfigDir, "dotfiles")

	// Populate Git configuration from system, including auto-detecting ssh_key_path
	if err := PopulateGitConfigFromSystem(&config.Git); err != nil {
		return nil, fmt.Errorf("failed to populate git configuration: %w", err)
	}

	return &config, nil
}

// GetConfigPath returns the path to the anvil configuration file
func GetConfigPath() string {
	return filepath.Join(getHomeDir(), constants.AnvilConfigDir, constants.ConfigFileName)
}

// getSampleConfigPath returns the path to the sample configuration file
func getSampleConfigPath() string {
	// Get the directory where the binary is located
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to current working directory
		wd, _ := os.Getwd()
		return filepath.Join(wd, "assets", "settings-sample.yaml")
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "assets", "settings-sample.yaml")
}

// CreateDirectories creates necessary directories for anvil
func CreateDirectories() error {
	configDir := GetConfigDirectory()

	// Only create the main config directory
	if err := utils.EnsureDirectory(configDir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", configDir, err)
	}

	return nil
}

// GenerateDefaultSettings generates the default settings.yaml file
func GenerateDefaultSettings() error {
	configPath := GetConfigPath()

	// Check if settings.yaml already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // File already exists, don't overwrite
	}

	// Load the sample configuration
	config, err := LoadSampleConfig()
	if err != nil {
		return fmt.Errorf("failed to load sample config: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, constants.FilePerm); err != nil {
		return fmt.Errorf("failed to write settings.yaml: %w", err)
	}

	return nil
}

// LoadConfig loads the anvil configuration from settings.yaml
func LoadConfig() (*AnvilConfig, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings.yaml: %w", err)
	}

	var config AnvilConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings.yaml: %w", err)
	}

	// Validate and auto-correct GitHub configuration
	if ValidateAndFixGitHubConfig(&config) {
		// Save the corrected configuration back to file
		if err := SaveConfig(&config); err != nil {
			// Don't fail loading if we can't save the correction, just warn
			fmt.Printf("Warning: Could not save corrected GitHub configuration: %v\n", err)
		}
	}

	return &config, nil
}

// SaveConfig saves the anvil configuration to settings.yaml
func SaveConfig(config *AnvilConfig) error {
	configPath := GetConfigPath()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	if err := os.WriteFile(configPath, data, constants.FilePerm); err != nil {
		return fmt.Errorf("failed to write settings.yaml: %w", err)
	}

	// Invalidate cache after saving
	invalidateCache()

	return nil
}

// GetGroupTools returns the tools for a specific group
func GetGroupTools(groupName string) ([]string, error) {
	var result []string
	err := withConfig(func(config *AnvilConfig) error {
		// Check if the group exists in the Groups map
		if tools, exists := config.Groups[groupName]; exists {
			result = tools
			return nil
		}
		return fmt.Errorf("group '%s' not found", groupName)
	})
	return result, err
}

// GetAvailableGroups returns all available groups
func GetAvailableGroups() (map[string][]string, error) {
	var groups map[string][]string
	err := withConfig(func(config *AnvilConfig) error {
		groups = make(map[string][]string)
		// Add built-in groups
		for name, tools := range config.Groups {
			groups[name] = tools
		}
		return nil
	})
	return groups, err
}

// GetBuiltInGroups returns the list of built-in group names
func GetBuiltInGroups() []string {
	return []string{"dev", "new-laptop"}
}

// IsBuiltInGroup checks if a group name is a built-in group
func IsBuiltInGroup(groupName string) bool {
	builtInGroups := GetBuiltInGroups()
	for _, group := range builtInGroups {
		if group == groupName {
			return true
		}
	}
	return false
}

// AddCustomGroup adds a new custom group
func AddCustomGroup(name string, tools []string) error {
	return withConfigAndSave(func(config *AnvilConfig) error {
		ensureGroupsMap(config)
		config.Groups[name] = tools
		return nil
	})
}

// UpdateGroupTools updates the tools list for an existing group
func UpdateGroupTools(groupName string, tools []string) error {
	return withConfigAndSave(func(config *AnvilConfig) error {
		// Check if the group exists
		if _, exists := config.Groups[groupName]; !exists {
			return fmt.Errorf("group '%s' does not exist", groupName)
		}
		// Update the group with new tools list
		config.Groups[groupName] = tools
		return nil
	})
}

// AddAppToGroup adds an app to a group, creating the group if it doesn't exist
func AddAppToGroup(groupName string, appName string) error {
	return withConfigAndSave(func(config *AnvilConfig) error {
		ensureGroupsMap(config)

		// Check if the group exists
		if tools, exists := config.Groups[groupName]; exists {
			// Check if app is already in the group to avoid duplicates
			for _, tool := range tools {
				if tool == appName {
					return nil // App already in group, no need to add
				}
			}
			// Add app to existing group
			config.Groups[groupName] = append(tools, appName)
		} else {
			// Create new group with the app
			config.Groups[groupName] = []string{appName}
		}
		return nil
	})
}

// CheckEnvironmentConfigurations checks local environment configurations
func CheckEnvironmentConfigurations() []string {
	var warnings []string

	// Check Git configuration
	if gitUser, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserName); err != nil || !gitUser.Success || strings.TrimSpace(gitUser.Output) == "" {
		warnings = append(warnings, fmt.Sprintf("Configure git user.name: %s %s %s %s 'Your Name'", constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserName))
	}

	if gitEmail, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserEmail); err != nil || !gitEmail.Success || strings.TrimSpace(gitEmail.Output) == "" {
		warnings = append(warnings, fmt.Sprintf("Configure git user.email: %s %s %s %s 'your.email@example.com'", constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserEmail))
	}

	// Check SSH keys
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, constants.SSHDir)
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		warnings = append(warnings, "Set up SSH keys for GitHub: ssh-keygen -t ed25519 -C 'your.email@example.com'")
	} else {
		// Check for common SSH key files
		keyFiles := []string{"id_rsa", "id_ed25519", "id_ecdsa"}
		hasKey := false
		for _, keyFile := range keyFiles {
			if _, err := os.Stat(filepath.Join(sshDir, keyFile)); err == nil {
				hasKey = true
				break
			}
		}
		if !hasKey {
			warnings = append(warnings, fmt.Sprintf("No SSH keys found in ~/%s - consider generating SSH keys for GitHub", constants.SSHDir))
		}
	}

	// Check for common environment variables
	envVars := []string{constants.EnvEditor, constants.EnvShell}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			warnings = append(warnings, fmt.Sprintf("Consider setting %s environment variable", envVar))
		}
	}

	return warnings
}

// GetConfigDirectory returns the anvil configuration directory
func GetConfigDirectory() string {
	return filepath.Join(getHomeDir(), constants.AnvilConfigDir)
}

// GetToolConfig returns the configuration for a specific tool
func GetToolConfig(toolName string) (*ToolInstallConfig, error) {
	var result *ToolInstallConfig
	err := withConfig(func(config *AnvilConfig) error {
		if toolConfig, exists := config.ToolConfigs.Tools[toolName]; exists {
			result = &toolConfig
		} else {
			result = getDefaultToolConfig()
		}
		return nil
	})
	return result, err
}

// GetToolConfigs returns configurations for multiple tools in a single operation
func GetToolConfigs(toolNames []string) (map[string]*ToolInstallConfig, error) {
	var configs map[string]*ToolInstallConfig
	err := withConfig(func(config *AnvilConfig) error {
		configs = make(map[string]*ToolInstallConfig, len(toolNames))
		for _, toolName := range toolNames {
			if toolConfig, exists := config.ToolConfigs.Tools[toolName]; exists {
				configs[toolName] = &toolConfig
			} else {
				configs[toolName] = getDefaultToolConfig()
			}
		}
		return nil
	})
	return configs, err
}

// getDefaultToolConfig returns a default tool configuration
func getDefaultToolConfig() *ToolInstallConfig {
	return &ToolInstallConfig{
		PostInstallScript: "",
		EnvironmentSetup:  make(map[string]string),
		ConfigCheck:       false,
		Dependencies:      []string{},
	}
}

// SetToolConfig sets the configuration for a specific tool
func SetToolConfig(toolName string, toolConfig ToolInstallConfig) error {
	return withConfigAndSave(func(config *AnvilConfig) error {
		ensureToolConfigsMap(config)
		config.ToolConfigs.Tools[toolName] = toolConfig
		return nil
	})
}

// AddInstalledApp adds an app to the installed apps list if it's not already there
func AddInstalledApp(appName string) error {
	config, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if app is already in the installed apps list
	for _, installedApp := range config.Tools.InstalledApps {
		if installedApp == appName {
			return nil // App already tracked, no need to add
		}
	}

	// Check if app is already in required or optional tools to avoid duplicates
	for _, tool := range config.Tools.RequiredTools {
		if tool == appName {
			return nil // App is a required tool, no need to track separately
		}
	}

	for _, tool := range config.Tools.OptionalTools {
		if tool == appName {
			return nil // App is an optional tool, no need to track separately
		}
	}

	// Check if app is in any group to avoid duplicates
	groups, err := GetAvailableGroups()
	if err == nil {
		for _, tools := range groups {
			for _, tool := range tools {
				if tool == appName {
					return nil // App is in a group, no need to track separately
				}
			}
		}
	}

	// Add the app to the installed apps list
	config.Tools.InstalledApps = append(config.Tools.InstalledApps, appName)

	// Save the updated configuration
	return SaveConfig(config)
}

// GetInstalledApps returns the list of individually installed applications
func GetInstalledApps() ([]string, error) {
	config, err := getCachedConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return config.Tools.InstalledApps, nil
}

// IsAppTracked checks if an app is being tracked in any category
func IsAppTracked(appName string) (bool, error) {
	config, err := getCachedConfig()
	if err != nil {
		return false, fmt.Errorf("failed to load config: %w", err)
	}

	// Check in required tools
	for _, tool := range config.Tools.RequiredTools {
		if tool == appName {
			return true, nil
		}
	}

	// Check in optional tools
	for _, tool := range config.Tools.OptionalTools {
		if tool == appName {
			return true, nil
		}
	}

	// Check in installed apps
	for _, app := range config.Tools.InstalledApps {
		if app == appName {
			return true, nil
		}
	}

	// Check in groups
	groups, err := GetAvailableGroups()
	if err == nil {
		for _, tools := range groups {
			for _, tool := range tools {
				if tool == appName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// RemoveInstalledApp removes an app from the installed apps list
func RemoveInstalledApp(appName string) error {
	config, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Find and remove the app from the installed apps list
	for i, app := range config.Tools.InstalledApps {
		if app == appName {
			config.Tools.InstalledApps = append(config.Tools.InstalledApps[:i], config.Tools.InstalledApps[i+1:]...)
			return SaveConfig(config)
		}
	}

	return nil // App not found, nothing to remove
}

// LocationSource represents where an app config location was found
type LocationSource int

const (
	LocationConfigs LocationSource = iota // Found in configs section of settings.yaml
	LocationTemp                          // Found in temp directory (pulled but not configured)
)

// String returns a string representation of the location source
func (ls LocationSource) String() string {
	switch ls {
	case LocationConfigs:
		return "configs"
	case LocationTemp:
		return "temp"
	default:
		return "unknown"
	}
}

// GetAppConfigPath checks if an app has a configured local path in the configs section
func GetAppConfigPath(appName string) (string, bool, error) {
	config, err := getCachedConfig()
	if err != nil {
		return "", false, fmt.Errorf("failed to load config: %w", err)
	}

	if config.Configs == nil {
		return "", false, nil
	}

	path, exists := config.Configs[appName]
	if !exists {
		return "", false, nil
	}

	// Verify the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", false, fmt.Errorf("configured path for %s does not exist: %s", appName, path)
	}

	return path, true, nil
}

// GetTempAppPath checks if an app directory exists in the temp directory (from previous pull)
func GetTempAppPath(appName string) (string, bool, error) {
	tempPath := filepath.Join(GetConfigDirectory(), "temp", appName)
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		return "", false, nil
	}

	return tempPath, true, nil
}

// ResolveAppLocation finds the config location for an app following the priority order
func ResolveAppLocation(appName string) (string, LocationSource, error) {
	// Priority 1: Check configs section in settings.yaml
	if path, found, err := GetAppConfigPath(appName); err != nil {
		return "", LocationConfigs, err
	} else if found {
		return path, LocationConfigs, nil
	}

	// Priority 2: Check temp directory (pulled configs)
	if path, found, err := GetTempAppPath(appName); err != nil {
		return "", LocationTemp, err
	} else if found {
		return path, LocationTemp, nil
	}

	// Not found anywhere
	return "", LocationConfigs, fmt.Errorf("app '%s' not found in configs or temp directory", appName)
}

// SetAppConfigPath sets the config path for an app in the configs section
func SetAppConfigPath(appName, configPath string) error {
	return withConfigAndSave(func(config *AnvilConfig) error {
		ensureConfigsMap(config)
		config.Configs[appName] = configPath
		return nil
	})
}

// GetConfiguredApps returns a list of all apps that have configured paths
func GetConfiguredApps() ([]string, error) {
	var apps []string
	err := withConfig(func(config *AnvilConfig) error {
		for appName := range config.Configs {
			apps = append(apps, appName)
		}
		return nil
	})
	return apps, err
}
