package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocajuanma/anvil/pkg/system"
	"gopkg.in/yaml.v2"
)

// AnvilConfig represents the main anvil configuration
type AnvilConfig struct {
	Version     string            `yaml:"version"`
	Directories AnvilDirectories  `yaml:"directories"`
	Tools       AnvilTools        `yaml:"tools"`
	Groups      AnvilGroups       `yaml:"groups"`
	Git         GitConfig         `yaml:"git"`
	Environment map[string]string `yaml:"environment"`
}

// AnvilDirectories represents directory configurations
type AnvilDirectories struct {
	Config string `yaml:"config"`
	Cache  string `yaml:"cache"`
	Data   string `yaml:"data"`
}

// AnvilTools represents tool configurations
type AnvilTools struct {
	RequiredTools []string `yaml:"required_tools"`
	OptionalTools []string `yaml:"optional_tools"`
}

// AnvilGroups represents grouped tool configurations
type AnvilGroups struct {
	Dev       []string            `yaml:"dev"`
	NewLaptop []string            `yaml:"new-laptop"`
	Custom    map[string][]string `yaml:"custom"`
}

// GitConfig represents git configuration
type GitConfig struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
}

// GetDefaultConfig returns the default anvil configuration
func GetDefaultConfig() *AnvilConfig {
	homeDir, _ := os.UserHomeDir()

	return &AnvilConfig{
		Version: "1.0.0",
		Directories: AnvilDirectories{
			Config: filepath.Join(homeDir, ".anvil"),
			Cache:  filepath.Join(homeDir, ".anvil", "cache"),
			Data:   filepath.Join(homeDir, ".anvil", "data"),
		},
		Tools: AnvilTools{
			RequiredTools: []string{"git", "curl"},
			OptionalTools: []string{"brew", "docker", "kubectl"},
		},
		Groups: AnvilGroups{
			Dev:       []string{"git", "zsh", "iterm2", "vscode"},
			NewLaptop: []string{"slack", "chrome", "1password"},
			Custom:    make(map[string][]string),
		},
		Git: GitConfig{
			Username: "",
			Email:    "",
		},
		Environment: make(map[string]string),
	}
}

// GetConfigPath returns the path to the anvil configuration file
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".anvil", "settings.yaml")
}

// CreateDirectories creates necessary directories for anvil
func CreateDirectories() error {
	config := GetDefaultConfig()

	directories := []string{
		config.Directories.Config,
		config.Directories.Cache,
		config.Directories.Data,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
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
	if gitUser, err := system.RunCommand("git", "config", "--global", "user.name"); err == nil && gitUser.Success {
		config.Git.Username = strings.TrimSpace(gitUser.Output)
	}

	if gitEmail, err := system.RunCommand("git", "config", "--global", "user.email"); err == nil && gitEmail.Success {
		config.Git.Email = strings.TrimSpace(gitEmail.Output)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
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

	return &config, nil
}

// SaveConfig saves the anvil configuration to settings.yaml
func SaveConfig(config *AnvilConfig) error {
	configPath := GetConfigPath()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings.yaml: %w", err)
	}

	return nil
}

// GetGroupTools returns the tools for a specific group
func GetGroupTools(groupName string) ([]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	switch groupName {
	case "dev":
		return config.Groups.Dev, nil
	case "new-laptop":
		return config.Groups.NewLaptop, nil
	default:
		if tools, exists := config.Groups.Custom[groupName]; exists {
			return tools, nil
		}
		return nil, fmt.Errorf("group '%s' not found", groupName)
	}
}

// GetAvailableGroups returns all available groups
func GetAvailableGroups() (map[string][]string, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	groups := make(map[string][]string)
	groups["dev"] = config.Groups.Dev
	groups["new-laptop"] = config.Groups.NewLaptop

	for name, tools := range config.Groups.Custom {
		groups[name] = tools
	}

	return groups, nil
}

// AddCustomGroup adds a new custom group
func AddCustomGroup(name string, tools []string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.Groups.Custom == nil {
		config.Groups.Custom = make(map[string][]string)
	}

	config.Groups.Custom[name] = tools

	return SaveConfig(config)
}

// CheckEnvironmentConfigurations checks local environment configurations
func CheckEnvironmentConfigurations() []string {
	var warnings []string

	// Check Git configuration
	if gitUser, err := system.RunCommand("git", "config", "--global", "user.name"); err != nil || !gitUser.Success || strings.TrimSpace(gitUser.Output) == "" {
		warnings = append(warnings, "Configure git user.name: git config --global user.name 'Your Name'")
	}

	if gitEmail, err := system.RunCommand("git", "config", "--global", "user.email"); err != nil || !gitEmail.Success || strings.TrimSpace(gitEmail.Output) == "" {
		warnings = append(warnings, "Configure git user.email: git config --global user.email 'your.email@example.com'")
	}

	// Check SSH keys
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")
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
			warnings = append(warnings, "No SSH keys found in ~/.ssh - consider generating SSH keys for GitHub")
		}
	}

	// Check for common environment variables
	envVars := []string{"EDITOR", "SHELL"}
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
	return filepath.Join(homeDir, ".anvil")
}

// GetCacheDirectory returns the anvil cache directory
func GetCacheDirectory() string {
	return filepath.Join(GetConfigDirectory(), "cache")
}

// GetDataDirectory returns the anvil data directory
func GetDataDirectory() string {
	return filepath.Join(GetConfigDirectory(), "data")
}
