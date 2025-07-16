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
	"regexp"
	"strings"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/interfaces"
)

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

	// Validate directories
	if err := cv.validateDirectories(&anvilConfig.Directories); err != nil {
		return fmt.Errorf("directory validation failed: %w", err)
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

// validateDirectories validates directory configurations
func (cv *ConfigValidator) validateDirectories(dirs *AnvilDirectories) error {
	if dirs.Config == "" {
		return fmt.Errorf("config directory cannot be empty")
	}

	// Check if directory is an absolute path
	if !filepath.IsAbs(dirs.Config) {
		return fmt.Errorf("config directory must be an absolute path: %s", dirs.Config)
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
	// Validate built-in groups
	if len(groups.Dev) == 0 {
		return fmt.Errorf("dev group cannot be empty")
	}

	if len(groups.NewLaptop) == 0 {
		return fmt.Errorf("new-laptop group cannot be empty")
	}

	// Validate dev group tools
	for _, tool := range groups.Dev {
		if err := cv.ValidateAppName(tool); err != nil {
			return fmt.Errorf("invalid tool in dev group: %w", err)
		}
	}

	// Validate new-laptop group tools
	for _, tool := range groups.NewLaptop {
		if err := cv.ValidateAppName(tool); err != nil {
			return fmt.Errorf("invalid tool in new-laptop group: %w", err)
		}
	}

	// Validate custom groups
	for groupName, tools := range groups.Custom {
		if err := cv.ValidateGroupName(groupName); err != nil {
			return fmt.Errorf("invalid custom group name: %w", err)
		}

		if len(tools) == 0 {
			return fmt.Errorf("custom group '%s' cannot be empty", groupName)
		}

		for _, tool := range tools {
			if err := cv.ValidateAppName(tool); err != nil {
				return fmt.Errorf("invalid tool in custom group '%s': %w", groupName, err)
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

// ValidateDirectoryAccess validates that directories exist and are accessible
func ValidateDirectoryAccess(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("cannot access directory %s: %w", dirPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dirPath)
	}

	// Check if directory is writable
	testFile := filepath.Join(dirPath, ".anvil_test")
	if err := os.WriteFile(testFile, []byte("test"), constants.FilePerm); err != nil {
		return fmt.Errorf("directory is not writable: %s", dirPath)
	}

	// Clean up test file
	os.Remove(testFile)

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
