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

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/system"
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
	SSHDir     string `yaml:"ssh_dir,omitempty"`      // Reference to .ssh directory
}

// GitHubConfig represents GitHub repository configuration for config sync
type GitHubConfig struct {
	ConfigRepo  string `yaml:"config_repo"`             // GitHub repository URL for configs (e.g., "username/dotfiles")
	Branch      string `yaml:"branch"`                  // Branch to use (default: "main")
	LocalPath   string `yaml:"local_path"`              // Local path where configs are stored/synced
	Token       string `yaml:"token,omitempty"`         // GitHub token (use env var reference)
	TokenEnvVar string `yaml:"token_env_var,omitempty"` // Environment variable name for token
}

// SyncConfig represents configuration for selective synchronization
type SyncConfig struct {
	ExcludeSections  []string          `yaml:"exclude_sections,omitempty"`  // Sections to exclude from sync
	TemplateSections []string          `yaml:"template_sections,omitempty"` // Sections to process as templates
	IncludeOverride  []string          `yaml:"include_override,omitempty"`  // Force include sections (overrides exclude)
	TemplateValues   map[string]string `yaml:"template_values,omitempty"`   // Template replacement values
}

// AnvilConfig represents the main anvil configuration
type AnvilConfig struct {
	Version     string            `yaml:"version"`
	SyncConfig  SyncConfig        `yaml:"_sync_config,omitempty"`
	Directories AnvilDirectories  `yaml:"directories"`
	Tools       AnvilTools        `yaml:"tools"`
	Groups      AnvilGroups       `yaml:"groups"`
	Git         GitConfig         `yaml:"git"`
	GitHub      GitHubConfig      `yaml:"github"`
	Environment map[string]string `yaml:"environment"`
	ToolConfigs AnvilToolConfigs  `yaml:"tool_configs,omitempty"`
}

// AnvilDirectories represents directory configurations
type AnvilDirectories struct {
	Config string `yaml:"config"`
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

// invalidateCache clears the configuration cache
func invalidateCache() {
	configCacheMutex.Lock()
	defer configCacheMutex.Unlock()
	configCache = nil
}

// GetDefaultConfig returns the default anvil configuration
func GetDefaultConfig() *AnvilConfig {
	homeDir, _ := os.UserHomeDir()

	return &AnvilConfig{
		Version: "1.0.0",
		Directories: AnvilDirectories{
			Config: filepath.Join(homeDir, constants.AnvilConfigDir),
		},
		Tools: AnvilTools{
			RequiredTools: []string{constants.PkgGit, constants.CurlCommand},
			OptionalTools: []string{constants.BrewCommand, constants.PkgDocker, constants.PkgKubectl},
			InstalledApps: []string{}, // Initialize empty slice for tracking
		},
		Groups: AnvilGroups{
			"dev":        {constants.PkgGit, constants.PkgZsh, constants.PkgIterm2, constants.PkgVSCode},
			"new-laptop": {constants.PkgSlack, constants.PkgChrome, constants.Pkg1Password},
		},
		Git: GitConfig{
			Username:   "",
			Email:      "",
			SSHKeyPath: filepath.Join(homeDir, constants.SSHDir, "id_rsa"),
			SSHDir:     filepath.Join(homeDir, constants.SSHDir),
		},
		GitHub: GitHubConfig{
			ConfigRepo:  "", // User needs to populate this
			Branch:      "main",
			LocalPath:   filepath.Join(homeDir, constants.AnvilConfigDir, "dotfiles"),
			TokenEnvVar: "GITHUB_TOKEN", // Recommend using env var for token
		},
		Environment: make(map[string]string),
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

// GetConfigPath returns the path to the anvil configuration file
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, constants.AnvilConfigDir, constants.ConfigFileName)
}

// CreateDirectories creates necessary directories for anvil
func CreateDirectories() error {
	config := GetDefaultConfig()

	// Only create the main config directory
	if err := os.MkdirAll(config.Directories.Config, constants.DirPerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", config.Directories.Config, err)
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

	config := GetDefaultConfig()

	// Try to populate git configuration from system
	if gitUser, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserName); err == nil && gitUser.Success {
		config.Git.Username = strings.TrimSpace(gitUser.Output)
	}

	if gitEmail, err := system.RunCommand(constants.GitCommand, constants.GitConfig, constants.GitGlobal, constants.GitUserEmail); err == nil && gitEmail.Success {
		config.Git.Email = strings.TrimSpace(gitEmail.Output)
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

// LoadConfigFromPath loads the anvil configuration from a specific path
func LoadConfigFromPath(configPath string) (*AnvilConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", configPath, err)
	}

	var config AnvilConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from %s: %w", configPath, err)
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
	config, err := getCachedConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Check if the group exists in the Groups map
	if tools, exists := config.Groups[groupName]; exists {
		return tools, nil
	}

	return nil, fmt.Errorf("group '%s' not found", groupName)
}

// GetAvailableGroups returns all available groups
func GetAvailableGroups() (map[string][]string, error) {
	config, err := getCachedConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	groups := make(map[string][]string)

	// Add built-in groups
	for name, tools := range config.Groups {
		groups[name] = tools
	}

	return groups, nil
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
	config, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.Groups == nil {
		config.Groups = make(map[string][]string)
	}

	config.Groups[name] = tools

	return SaveConfig(config)
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
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, constants.AnvilConfigDir)
}

// GetToolConfig returns the configuration for a specific tool
func GetToolConfig(toolName string) (*ToolInstallConfig, error) {
	config, err := getCachedConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if toolConfig, exists := config.ToolConfigs.Tools[toolName]; exists {
		return &toolConfig, nil
	}

	// Return default config if not found
	return &ToolInstallConfig{
		PostInstallScript: "",
		EnvironmentSetup:  make(map[string]string),
		ConfigCheck:       false,
		Dependencies:      []string{},
	}, nil
}

// SetToolConfig sets the configuration for a specific tool
func SetToolConfig(toolName string, config ToolInstallConfig) error {
	anvilConfig, err := getCachedConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if anvilConfig.ToolConfigs.Tools == nil {
		anvilConfig.ToolConfigs.Tools = make(map[string]ToolInstallConfig)
	}

	anvilConfig.ToolConfigs.Tools[toolName] = config

	return SaveConfig(anvilConfig)
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
	config, err := LoadConfig()
	if err == nil && (config.Git.Username == "" || config.Git.Email == "") {
		return fmt.Errorf("git is not fully configured - consider setting username and email")
	}
	return nil
}

// FilterForSync creates a filtered version of the configuration for synchronization
// Excludes specified sections and processes templates according to sync configuration
func FilterForSync(config *AnvilConfig) (*AnvilConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create a deep copy of the configuration
	filteredConfig := &AnvilConfig{
		Version:     config.Version,
		SyncConfig:  config.SyncConfig,
		Directories: config.Directories,
		Tools:       config.Tools,
		Groups:      config.Groups,
		Git:         config.Git,
		GitHub:      config.GitHub,
		Environment: make(map[string]string),
		ToolConfigs: config.ToolConfigs,
	}

	// Copy environment map
	for k, v := range config.Environment {
		filteredConfig.Environment[k] = v
	}

	// Apply filtering based on sync configuration
	excludeSections := config.SyncConfig.ExcludeSections
	templateSections := config.SyncConfig.TemplateSections
	includeOverride := config.SyncConfig.IncludeOverride

	// Process each section according to filtering rules
	for _, section := range excludeSections {
		if !contains(includeOverride, section) {
			if err := excludeSection(filteredConfig, section); err != nil {
				return nil, fmt.Errorf("failed to exclude section %s: %w", section, err)
			}
		}
	}

	// Apply templates to specified sections
	for _, section := range templateSections {
		if err := applyTemplateToSection(filteredConfig, section); err != nil {
			return nil, fmt.Errorf("failed to apply template to section %s: %w", section, err)
		}
	}

	// Remove sync config from filtered version (shouldn't be synced)
	filteredConfig.SyncConfig = SyncConfig{}

	return filteredConfig, nil
}

// ApplyTemplates applies template values to a configuration
// Used when pulling configurations to replace template placeholders with actual values
func ApplyTemplates(config *AnvilConfig, templateValues map[string]string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Apply templates to Git section
	if err := applyTemplateValues(&config.Git.Username, templateValues); err != nil {
		return fmt.Errorf("failed to apply template to git username: %w", err)
	}
	if err := applyTemplateValues(&config.Git.Email, templateValues); err != nil {
		return fmt.Errorf("failed to apply template to git email: %w", err)
	}

	// Apply templates to Environment section
	for key, value := range config.Environment {
		if err := applyTemplateValues(&value, templateValues); err != nil {
			return fmt.Errorf("failed to apply template to environment %s: %w", key, err)
		}
		config.Environment[key] = value
	}

	return nil
}

// excludeSection removes a section from the configuration based on section path
func excludeSection(config *AnvilConfig, sectionPath string) error {
	switch sectionPath {
	case "git":
		config.Git = GitConfig{}
	case "environment":
		config.Environment = make(map[string]string)
	case "environment.machine_specific":
		// Remove machine_specific keys from environment
		delete(config.Environment, "machine_specific")
	case "tool_configs":
		config.ToolConfigs = AnvilToolConfigs{}
	case "directories":
		config.Directories = AnvilDirectories{}
	default:
		// For nested paths like "environment.KEY", remove specific environment variable
		if strings.HasPrefix(sectionPath, "environment.") {
			envKey := strings.TrimPrefix(sectionPath, "environment.")
			delete(config.Environment, envKey)
		} else {
			return fmt.Errorf("unsupported section path: %s", sectionPath)
		}
	}
	return nil
}

// applyTemplateToSection applies template placeholders to a specific section
func applyTemplateToSection(config *AnvilConfig, section string) error {
	switch section {
	case "git":
		config.Git.Username = "{{ REPLACE_USERNAME }}"
		config.Git.Email = "{{ REPLACE_EMAIL }}"
		if config.Git.SSHKeyPath != "" {
			config.Git.SSHKeyPath = "{{ REPLACE_SSH_KEY_PATH }}"
		}
	case "environment":
		// Apply templates to all environment variables
		for key, value := range config.Environment {
			if strings.Contains(value, "/") { // Likely a path
				config.Environment[key] = "{{ REPLACE_" + strings.ToUpper(key) + " }}"
			}
		}
	default:
		return fmt.Errorf("template not supported for section: %s", section)
	}
	return nil
}

// applyTemplateValues replaces template placeholders with actual values
func applyTemplateValues(target *string, templateValues map[string]string) error {
	if target == nil {
		return nil
	}

	original := *target
	result := original

	// Replace common template placeholders
	replacements := map[string]string{
		"{{ REPLACE_USERNAME }}":     templateValues["username"],
		"{{ REPLACE_EMAIL }}":        templateValues["email"],
		"{{ REPLACE_SSH_KEY_PATH }}": templateValues["ssh_key_path"],
	}

	// Apply custom template values
	for placeholder, value := range templateValues {
		templateKey := fmt.Sprintf("{{ REPLACE_%s }}", strings.ToUpper(placeholder))
		if value != "" {
			replacements[templateKey] = value
		}
	}

	// Perform replacements
	for placeholder, value := range replacements {
		if value != "" {
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	*target = result
	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PromptForTemplateValues prompts the user for template values needed for configuration
func PromptForTemplateValues(config *AnvilConfig) (map[string]string, error) {
	templateValues := make(map[string]string)

	// Check what templates are needed based on configuration content
	if needsGitTemplate(config) {
		username := promptForInput("Enter your git username", "")
		if username != "" {
			templateValues["username"] = username
		}

		email := promptForInput("Enter your git email", "")
		if email != "" {
			templateValues["email"] = email
		}
	}

	// Check for environment template needs
	for key, value := range config.Environment {
		if strings.Contains(value, "{{ REPLACE_") {
			promptKey := strings.ToLower(key)
			promptValue := promptForInput(fmt.Sprintf("Enter value for %s", key), "")
			if promptValue != "" {
				templateValues[promptKey] = promptValue
			}
		}
	}

	return templateValues, nil
}

// needsGitTemplate checks if git section needs template values
func needsGitTemplate(config *AnvilConfig) bool {
	return strings.Contains(config.Git.Username, "{{ REPLACE_") ||
		strings.Contains(config.Git.Email, "{{ REPLACE_")
}

// promptForInput prompts user for input with a default value
func promptForInput(prompt, defaultValue string) string {
	if defaultValue != "" {
		prompt = fmt.Sprintf("%s [%s]", prompt, defaultValue)
	}

	// For now, return empty string - in real implementation this would use terminal.Prompt
	// This allows the system to work without breaking existing functionality
	return ""
}
