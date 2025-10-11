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

// validateString validates a string with common rules
func validateString(value, fieldName string, maxLength int, pattern string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s too long (max %d characters)", fieldName, maxLength)
	}

	if pattern != "" {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
		if !regex.MatchString(value) {
			return fmt.Errorf("invalid %s format", fieldName)
		}
	}

	return nil
}

// validateEmail validates an email address
func validateEmail(email string) error {
	if email == "" {
		return nil // Empty email is allowed
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return validateString(email, "email", 100, pattern)
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
	if err := validateString(groupName, "group name", 50, `^[a-zA-Z0-9_-]+$`); err != nil {
		return fmt.Errorf("group name '%s' contains invalid characters. Only alphanumeric, underscore, and dash are allowed", groupName)
	}
	return nil
}

// ValidateAppName validates an application name
func (cv *ConfigValidator) ValidateAppName(appName string) error {
	if err := validateString(appName, "application name", 100, `^[a-zA-Z0-9_.-]+$`); err != nil {
		return fmt.Errorf("application name '%s' contains invalid characters. Only alphanumeric, underscore, dot, and dash are allowed", appName)
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
		if err := validateString(git.Username, "git username", 100, ""); err != nil {
			return err
		}
	}

	if err := validateEmail(git.Email); err != nil {
		return fmt.Errorf("invalid git email format: %s", git.Email)
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
