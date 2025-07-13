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
			Cache:  filepath.Join(homeDir, constants.AnvilConfigDir, constants.CacheSubDir),
			Data:   filepath.Join(homeDir, constants.AnvilConfigDir, constants.DataSubDir),
		},
		Tools: AnvilTools{
			RequiredTools: []string{constants.PkgGit, constants.CurlCommand},
			OptionalTools: []string{constants.BrewCommand, constants.PkgDocker, constants.PkgKubectl},
		},
		Groups: AnvilGroups{
			Dev:       []string{constants.PkgGit, constants.PkgZsh, constants.PkgIterm2, constants.PkgVSCode},
			NewLaptop: []string{constants.PkgSlack, constants.PkgChrome, constants.Pkg1Password},
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
	return filepath.Join(homeDir, constants.AnvilConfigDir, constants.ConfigFileName)
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
		if err := os.MkdirAll(dir, constants.DirPerm); err != nil {
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
	config, err := getCachedConfig()
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
	config, err := getCachedConfig()
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

// GetCacheDirectory returns the anvil cache directory
func GetCacheDirectory() string {
	return filepath.Join(GetConfigDirectory(), constants.CacheSubDir)
}

// GetDataDirectory returns the anvil data directory
func GetDataDirectory() string {
	return filepath.Join(GetConfigDirectory(), constants.DataSubDir)
}
