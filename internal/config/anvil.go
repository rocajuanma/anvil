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

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/system"
	"gopkg.in/yaml.v2"
)

// AnvilConfig represents the main anvil configuration
type AnvilConfig struct {
	Version string            `yaml:"version"`
	Tools   AnvilTools        `yaml:"tools"`
	Groups  AnvilGroups       `yaml:"groups"`
	Configs map[string]string `yaml:"configs"` // Maps app names to their local config paths
	Git     GitConfig         `yaml:"git"`
	GitHub  GitHubConfig      `yaml:"github"`
}

// GetAnvilConfigDirectory returns the path to the anvil config directory
func GetAnvilConfigDirectory() string {
	homeDir, _ := system.GetHomeDir()
	return filepath.Join(homeDir, constants.ANVIL_CONFIG_DIR)
}

// GetAnvilConfigPath returns the path to the anvil config file
func GetAnvilConfigPath() string {
	return fmt.Sprintf("%s/%s", GetAnvilConfigDirectory(), constants.ANVIL_CONFIG_FILE)
}

// LoadConfig loads the anvil configuration from settings.yaml
func LoadConfig() (*AnvilConfig, error) {
	configPath := GetAnvilConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", constants.ANVIL_CONFIG_FILE, err)
	}

	var config AnvilConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", constants.ANVIL_CONFIG_FILE, err)
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

// LoadSampleConfigWithVersion loads the sample configuration with a specific version
func LoadSampleConfigWithVersion(version string) (*AnvilConfig, error) {
	// Use embedded sample config data
	configData := string(sampleConfigData)

	// Replace version placeholder if version is provided
	if version != "" {
		configData = strings.ReplaceAll(configData, "{{APP_VERSION}}", version)
	}

	var config AnvilConfig
	if err := yaml.Unmarshal([]byte(configData), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sample config: %w", err)
	}

	// Set dynamic paths
	homeDir, _ := system.GetHomeDir()
	config.GitHub.LocalPath = filepath.Join(homeDir, constants.ANVIL_CONFIG_DIR, "dotfiles")

	// Populate Git configuration from system, including auto-detecting ssh_key_path
	if err := PopulateGitConfigFromSystem(&config.Git); err != nil {
		return nil, fmt.Errorf("failed to populate git configuration: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the anvil configuration to settings.yaml
func SaveConfig(config *AnvilConfig) error {
	configPath := GetAnvilConfigPath()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	if err := os.WriteFile(configPath, data, constants.FilePerm); err != nil {
		return fmt.Errorf("failed to write %s: %w", constants.ANVIL_CONFIG_FILE, err)
	}

	// Invalidate cache after saving
	invalidateCache()

	return nil
}
