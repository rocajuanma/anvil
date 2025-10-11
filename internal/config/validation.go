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
	"regexp"
	"strings"

	"github.com/rocajuanma/anvil/internal/interfaces"
	"github.com/rocajuanma/palantir"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// ConfigValidator implements the Validator interface for configuration validation
type ConfigValidator struct {
	config *AnvilConfig
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(config *AnvilConfig) interfaces.Validator {
	return &ConfigValidator{
		config: config,
	}
}

// ValidateGroupName validates a group name
func (cv *ConfigValidator) ValidateGroupName(groupName string) error {
	if groupName == "" {
		return fmt.Errorf("group name cannot be empty")
	}

	// Check if group name contains invalid characters
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, groupName); !matched {
		return fmt.Errorf("group name '%s' contains invalid characters. Only alphanumeric, underscore, and dash are allowed", groupName)
	}

	// Check if group name is too long
	if len(groupName) > 50 {
		return fmt.Errorf("group name '%s' is too long (max 50 characters)", groupName)
	}

	return nil
}

// ValidateAppName validates an application name
func (cv *ConfigValidator) ValidateAppName(appName string) error {
	if appName == "" {
		return fmt.Errorf("application name cannot be empty")
	}

	// Check if app name contains invalid characters
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, appName); !matched {
		return fmt.Errorf("application name '%s' contains invalid characters. Only alphanumeric, underscore, dot, and dash are allowed", appName)
	}

	// Check if app name is too long
	if len(appName) > 100 {
		return fmt.Errorf("application name '%s' is too long (max 100 characters)", appName)
	}

	return nil
}

// ValidateFont validates a font name
func (cv *ConfigValidator) ValidateFont(font string) error {
	if font == "" {
		return fmt.Errorf("font name cannot be empty")
	}

	validFonts := []string{
		"standard", "doh", "big", "small", "banner", "block", "bubble", "digital",
		"ivrit", "lean", "mini", "script", "shadow", "slant", "speed", "term",
	}

	for _, validFont := range validFonts {
		if font == validFont {
			return nil
		}
	}

	return fmt.Errorf("invalid font '%s'. Valid fonts are: %s", font, strings.Join(validFonts, ", "))
}

// ValidateConfig validates the entire configuration
func (cv *ConfigValidator) ValidateConfig(config interface{}) error {
	anvilConfig, ok := config.(*AnvilConfig)
	if !ok {
		return fmt.Errorf("invalid config type: expected *AnvilConfig")
	}

	// Validate version
	if err := cv.validateVersion(anvilConfig.Version); err != nil {
		return fmt.Errorf("version validation failed: %w", err)
	}

	// Validate tools
	if err := cv.validateTools(&anvilConfig.Tools); err != nil {
		return fmt.Errorf("tools validation failed: %w", err)
	}

	// Validate groups
	if err := cv.validateGroups(&anvilConfig.Groups); err != nil {
		return fmt.Errorf("groups validation failed: %w", err)
	}

	// Validate git configuration
	if err := cv.validateGitConfig(&anvilConfig.Git); err != nil {
		return fmt.Errorf("git config validation failed: %w", err)
	}

	// Validate tool configs
	if err := cv.validateToolConfigs(&anvilConfig.ToolConfigs); err != nil {
		return fmt.Errorf("tool configs validation failed: %w", err)
	}

	return nil
}

// validateVersion validates the version string
func (cv *ConfigValidator) validateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Check semantic version format
	if matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, version); !matched {
		return fmt.Errorf("version '%s' is not in valid semantic version format (e.g., 1.0.0)", version)
	}

	return nil
}

// validateTools validates tool configurations
func (cv *ConfigValidator) validateTools(tools *AnvilTools) error {
	if len(tools.RequiredTools) == 0 {
		return fmt.Errorf("at least one required tool must be specified")
	}

	// Validate required tools
	for _, tool := range tools.RequiredTools {
		if err := cv.ValidateAppName(tool); err != nil {
			return fmt.Errorf("invalid required tool name: %w", err)
		}
	}

	// Validate optional tools
	for _, tool := range tools.OptionalTools {
		if err := cv.ValidateAppName(tool); err != nil {
			return fmt.Errorf("invalid optional tool name: %w", err)
		}
	}

	// Validate installed apps
	for _, app := range tools.InstalledApps {
		if err := cv.ValidateAppName(app); err != nil {
			return fmt.Errorf("invalid installed app name: %w", err)
		}
	}

	// Check for duplicates
	if err := cv.validateNoDuplicateTools(tools); err != nil {
		return err
	}

	return nil
}

// validateNoDuplicateTools checks for duplicate tool names
func (cv *ConfigValidator) validateNoDuplicateTools(tools *AnvilTools) error {
	allTools := make(map[string]bool)

	// Check required tools
	for _, tool := range tools.RequiredTools {
		if allTools[tool] {
			return fmt.Errorf("duplicate tool found: %s", tool)
		}
		allTools[tool] = true
	}

	// Check optional tools
	for _, tool := range tools.OptionalTools {
		if allTools[tool] {
			return fmt.Errorf("duplicate tool found: %s", tool)
		}
		allTools[tool] = true
	}

	// Check installed apps
	for _, app := range tools.InstalledApps {
		if allTools[app] {
			return fmt.Errorf("duplicate app found: %s", app)
		}
		allTools[app] = true
	}

	return nil
}

// validateGroups validates group configurations
func (cv *ConfigValidator) validateGroups(groups *AnvilGroups) error {
	if groups == nil || *groups == nil {
		return fmt.Errorf("groups configuration is nil")
	}

	groupsMap := *groups

	// Validate that required built-in groups exist
	devGroup, devExists := groupsMap["dev"]
	if !devExists || len(devGroup) == 0 {
		return fmt.Errorf("dev group is required and cannot be empty")
	}

	newLaptopGroup, newLaptopExists := groupsMap["essentials"]
	if !newLaptopExists || len(newLaptopGroup) == 0 {
		return fmt.Errorf("essentials group is required and cannot be empty")
	}

	// Validate all groups
	for groupName, tools := range groupsMap {
		if err := cv.ValidateGroupName(groupName); err != nil {
			return fmt.Errorf("invalid group name: %w", err)
		}

		if len(tools) == 0 {
			return fmt.Errorf("group '%s' cannot be empty", groupName)
		}

		for _, tool := range tools {
			if err := cv.ValidateAppName(tool); err != nil {
				return fmt.Errorf("invalid tool in group '%s': %w", groupName, err)
			}
		}
	}

	return nil
}

// validateGitConfig validates git configuration
func (cv *ConfigValidator) validateGitConfig(git *GitConfig) error {
	if git.Username != "" {
		if len(git.Username) > 100 {
			return fmt.Errorf("git username too long (max 100 characters)")
		}
	}

	if git.Email != "" {
		// Basic email validation
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, git.Email); !matched {
			return fmt.Errorf("invalid git email format: %s", git.Email)
		}
	}

	return nil
}

// validateToolConfigs validates tool-specific configurations
func (cv *ConfigValidator) validateToolConfigs(configs *AnvilToolConfigs) error {
	if configs.Tools == nil {
		return nil // Optional section
	}

	for toolName, toolConfig := range configs.Tools {
		if err := cv.ValidateAppName(toolName); err != nil {
			return fmt.Errorf("invalid tool config name: %w", err)
		}

		if err := cv.validateToolConfig(toolName, &toolConfig); err != nil {
			return fmt.Errorf("invalid config for tool '%s': %w", toolName, err)
		}
	}

	return nil
}

// validateToolConfig validates a single tool configuration
func (cv *ConfigValidator) validateToolConfig(toolName string, config *ToolInstallConfig) error {
	// Validate post-install script
	if config.PostInstallScript != "" {
		if len(config.PostInstallScript) > 500 {
			return fmt.Errorf("post-install script too long (max 500 characters)")
		}
	}

	// Validate environment setup
	for key, value := range config.EnvironmentSetup {
		if key == "" {
			return fmt.Errorf("environment variable name cannot be empty")
		}

		if matched, _ := regexp.MatchString(`^[A-Z_][A-Z0-9_]*$`, key); !matched {
			return fmt.Errorf("invalid environment variable name: %s", key)
		}

		if len(value) > 1000 {
			return fmt.Errorf("environment variable value too long (max 1000 characters)")
		}
	}

	// Validate dependencies
	for _, dep := range config.Dependencies {
		if err := cv.ValidateAppName(dep); err != nil {
			return fmt.Errorf("invalid dependency name: %w", err)
		}
	}

	return nil
}

// ValidateFileAccess validates that a file exists and is accessible
func ValidateFileAccess(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("cannot access file %s: %w", filePath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	return nil
}

// ValidateConfigFile validates the entire configuration file
func ValidateConfigFile(configPath string) error {
	// Check if file exists and is accessible
	if err := ValidateFileAccess(configPath); err != nil {
		return err
	}

	// Load and validate configuration
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	validator := NewConfigValidator(config)
	return validator.ValidateConfig(config)
}

// ValidateAndFixGitHubConfig validates and automatically fixes GitHub configuration
func ValidateAndFixGitHubConfig(config *AnvilConfig) bool {
	fixed := false

	if config.GitHub.ConfigRepo != "" {
		originalRepo := config.GitHub.ConfigRepo
		normalizedRepo := normalizeGitHubRepo(config.GitHub.ConfigRepo)

		if normalizedRepo != originalRepo {
			config.GitHub.ConfigRepo = normalizedRepo
			o := getOutputHandler()
			o.PrintInfo("🔧 Auto-corrected GitHub repository URL:")
			o.PrintInfo("   From: %s", originalRepo)
			o.PrintInfo("   To:   %s", normalizedRepo)
			o.PrintInfo("   Expected format: 'username/repository' (without domain)")
			fixed = true
		}
	}

	return fixed
}

// normalizeGitHubRepo converts various GitHub URL formats to the standard "username/repository" format
func normalizeGitHubRepo(repoURL string) string {
	if repoURL == "" {
		return repoURL
	}

	// Remove quotes if present
	repoURL = strings.Trim(repoURL, `"'`)

	// Handle different GitHub URL formats
	patterns := []struct {
		regex   *regexp.Regexp
		example string
	}{
		// HTTPS URLs
		{regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)(?:\.git)?/?$`), "https://github.com/username/repo"},
		{regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/.*$`), "https://github.com/username/repo/..."},

		// SSH URLs
		{regexp.MustCompile(`^git@github\.com:([^/]+/[^/]+)(?:\.git)?/?$`), "git@github.com:username/repo"},

		// Domain without protocol
		{regexp.MustCompile(`^github\.com/([^/]+/[^/]+)(?:\.git)?/?$`), "github.com/username/repo"},
		{regexp.MustCompile(`^github\.com/([^/]+/[^/]+)/.*$`), "github.com/username/repo/..."},

		// Already in correct format (username/repo)
		{regexp.MustCompile(`^([^/]+/[^/]+)$`), "username/repo"},
	}

	for _, pattern := range patterns {
		if matches := pattern.regex.FindStringSubmatch(repoURL); len(matches) > 1 {
			// Extract username/repository part
			userRepo := matches[1]
			// Remove .git suffix if present
			userRepo = strings.TrimSuffix(userRepo, ".git")
			return userRepo
		}
	}

	// If no pattern matches, return as-is (might be invalid, but let validation catch it)
	return repoURL
}

// validateGitHubRepoFormat validates that the repository is in the correct format
func validateGitHubRepoFormat(repo string) error {
	if repo == "" {
		return nil // Empty is handled elsewhere
	}

	// Expected format: username/repository
	repoPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
	if !repoPattern.MatchString(repo) {
		return fmt.Errorf(`invalid repository format: '%s'
Expected format: 'username/repository' (e.g., 'octocat/Hello-World')

Supported input formats that will be auto-corrected:
  • https://github.com/username/repository
  • https://github.com/username/repository.git
  • git@github.com:username/repository.git
  • github.com/username/repository

Your repository will be auto-corrected to the proper format when the config is loaded.`, repo)
	}

	return nil
}
