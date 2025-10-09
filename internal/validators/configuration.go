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

package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/rocajuanma/anvil/internal/config"
)

// GitConfigValidator checks if git configuration is properly set
type GitConfigValidator struct{}

func (v *GitConfigValidator) Name() string        { return "git-config" }
func (v *GitConfigValidator) Category() string    { return "configuration" }
func (v *GitConfigValidator) Description() string { return "Verify git configuration is properly set" }
func (v *GitConfigValidator) CanFix() bool        { return true }

func (v *GitConfigValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	var issues []string
	var details []string

	// Check username
	if cfg.Git.Username == "" {
		issues = append(issues, "username not set")
	} else {
		details = append(details, "Username: "+cfg.Git.Username)
	}

	// Check email
	if cfg.Git.Email == "" {
		issues = append(issues, "email not set")
	} else {
		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(cfg.Git.Email) {
			issues = append(issues, "email format invalid")
		} else {
			details = append(details, "Email: "+cfg.Git.Email)
		}
	}

	// Check SSH key path if specified
	if cfg.Git.SSHKeyPath != "" {
		details = append(details, "SSH Key: "+cfg.Git.SSHKeyPath)
	}

	if len(issues) > 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Git configuration incomplete: " + strings.Join(issues, ", "),
			Details:  details,
			FixHint:  "Run 'anvil doctor git-config --fix' to regenerate from local git configuration",
			AutoFix:  true,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "Git configuration is complete",
		Details:  details,
		FixHint:  "Run 'anvil doctor git-config --fix' to refresh from local git configuration",
		AutoFix:  true, // Always allow regenerating from local git config
	}
}

func (v *GitConfigValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	// Load current configuration
	currentConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load current configuration: %w", err)
	}

	// Store original values for comparison and logging
	originalUsername := currentConfig.Git.Username
	originalEmail := currentConfig.Git.Email
	originalSSHKeyPath := currentConfig.Git.SSHKeyPath

	// ALWAYS regenerate ALL git configuration from local git settings
	// This handles typos, invalid paths, and incorrect values
	if err := config.PopulateGitConfigFromSystem(&currentConfig.Git); err != nil {
		return fmt.Errorf("failed to read local git configuration: %w", err)
	}

	// Check if local git config is available
	if currentConfig.Git.Username == "" && currentConfig.Git.Email == "" {
		return fmt.Errorf("no local git configuration found (git config --global user.name/user.email not set)")
	}

	// Log what we're updating (for user visibility)
	changes := []string{}
	if currentConfig.Git.Username != originalUsername {
		if originalUsername == "" {
			changes = append(changes, fmt.Sprintf("username: (empty) → %s", currentConfig.Git.Username))
		} else {
			changes = append(changes, fmt.Sprintf("username: %s → %s", originalUsername, currentConfig.Git.Username))
		}
	}
	if currentConfig.Git.Email != originalEmail {
		if originalEmail == "" {
			changes = append(changes, fmt.Sprintf("email: (empty) → %s", currentConfig.Git.Email))
		} else {
			changes = append(changes, fmt.Sprintf("email: %s → %s", originalEmail, currentConfig.Git.Email))
		}
	}
	if currentConfig.Git.SSHKeyPath != originalSSHKeyPath {
		if originalSSHKeyPath == "" {
			changes = append(changes, fmt.Sprintf("ssh_key_path: (empty) → %s", currentConfig.Git.SSHKeyPath))
		} else {
			changes = append(changes, fmt.Sprintf("ssh_key_path: %s → %s", originalSSHKeyPath, currentConfig.Git.SSHKeyPath))
		}
	}

	// Always save the updated configuration (even if no visible changes)
	// This ensures the latest auto-detected values are persisted
	if err := config.SaveConfig(currentConfig); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}

	// Provide user feedback about what was updated
	if len(changes) > 0 {
		fmt.Printf("\n🔧 Updated git configuration:\n")
		for _, change := range changes {
			fmt.Printf("  • %s\n", change)
		}
	} else {
		fmt.Printf("\n✅ Git configuration verified and refreshed from local git config\n")
	}

	return nil
}

// GitHubConfigValidator checks if GitHub configuration is properly set
type GitHubConfigValidator struct{}

func (v *GitHubConfigValidator) Name() string     { return "github-config" }
func (v *GitHubConfigValidator) Category() string { return "configuration" }
func (v *GitHubConfigValidator) Description() string {
	return "Verify GitHub configuration is properly set"
}
func (v *GitHubConfigValidator) CanFix() bool { return false }

func (v *GitHubConfigValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	var issues []string
	var details []string

	// Check if config_repo is set
	if cfg.GitHub.ConfigRepo == "" {
		issues = append(issues, "config_repo not set")
	} else {
		// Validate repository format (should be "username/repository")
		repoRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
		if !repoRegex.MatchString(cfg.GitHub.ConfigRepo) {
			issues = append(issues, "config_repo format invalid (should be 'username/repository')")
		} else {
			details = append(details, "Repository: "+cfg.GitHub.ConfigRepo)
		}
	}

	// Check branch
	if cfg.GitHub.Branch == "" {
		issues = append(issues, "branch not set")
	} else {
		details = append(details, "Branch: "+cfg.GitHub.Branch)
	}

	// Check token environment variable
	if cfg.GitHub.TokenEnvVar == "" {
		issues = append(issues, "token_env_var not set")
	} else {
		details = append(details, "Token env var: "+cfg.GitHub.TokenEnvVar)
	}

	if len(issues) > 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "GitHub configuration incomplete: " + strings.Join(issues, ", "),
			Details:  details,
			FixHint:  "Set GitHub configuration manually in settings.yaml",
			AutoFix:  false,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "GitHub configuration is complete",
		Details:  details,
		AutoFix:  false,
	}
}

func (v *GitHubConfigValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	// For now, GitHub configuration fixes must be done manually
	return fmt.Errorf("GitHub configuration must be set manually in settings.yaml")
}

// SyncConfigValidator checks if sync configuration is valid
type SyncConfigValidator struct{}

func (v *SyncConfigValidator) Name() string        { return "sync-config" }
func (v *SyncConfigValidator) Category() string    { return "configuration" }
func (v *SyncConfigValidator) Description() string { return "Verify sync configuration is valid" }
func (v *SyncConfigValidator) CanFix() bool        { return false }

func (v *SyncConfigValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	// For now, skip sync config validation until the field is properly added
	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   SKIP,
		Message:  "Sync configuration validation not yet implemented",
		FixHint:  "Add _sync_config section to settings.yaml for selective synchronization",
		AutoFix:  false,
	}
}

func (v *SyncConfigValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	return fmt.Errorf("sync configuration issues must be fixed manually in settings.yaml")
}
