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
	"os"
	"strings"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/system"
)

// GitHubAccessValidator checks if GitHub API is accessible
type GitHubAccessValidator struct{}

func (v *GitHubAccessValidator) Name() string     { return "github-access" }
func (v *GitHubAccessValidator) Category() string { return "connectivity" }
func (v *GitHubAccessValidator) Description() string {
	return "Verify GitHub API access and authentication"
}
func (v *GitHubAccessValidator) CanFix() bool { return false }

func (v *GitHubAccessValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	// Skip if no GitHub config
	if cfg.GitHub.ConfigRepo == "" {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   SKIP,
			Message:  "No GitHub configuration found",
			FixHint:  "Configure GitHub repository in settings.yaml",
			AutoFix:  false,
		}
	}

	// Check if GitHub token is available
	var token string
	if cfg.GitHub.TokenEnvVar != "" {
		token = os.Getenv(cfg.GitHub.TokenEnvVar)
	}

	if token == "" {
		// Test SSH access as fallback
		result, err := system.RunCommand("ssh", "-T", "git@github.com")
		if err != nil || !strings.Contains(result.Output, "successfully authenticated") {
			return &ValidationResult{
				Name:     v.Name(),
				Category: v.Category(),
				Status:   FAIL,
				Message:  "No GitHub authentication available",
				Details:  []string{"No token found", "SSH authentication failed"},
				FixHint:  fmt.Sprintf("Set %s environment variable or configure SSH keys", cfg.GitHub.TokenEnvVar),
				AutoFix:  false,
			}
		}

		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   PASS,
			Message:  "GitHub SSH access confirmed",
			Details:  []string{"SSH authentication successful"},
			AutoFix:  false,
		}
	}

	// Test GitHub API with token
	result, err := system.RunCommand("curl", "-s", "-f", "-H", fmt.Sprintf("Authorization: token %s", token), "https://api.github.com/user")
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "GitHub API access failed",
			Details:  []string{"Token authentication failed"},
			FixHint:  fmt.Sprintf("Check %s environment variable", cfg.GitHub.TokenEnvVar),
			AutoFix:  false,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "GitHub API access confirmed",
		Details:  []string{"Token authentication successful"},
		AutoFix:  false,
	}
}

func (v *GitHubAccessValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	return fmt.Errorf("GitHub access issues must be fixed manually by setting up authentication")
}

// RepositoryValidator checks if the configured repository exists and is accessible
type RepositoryValidator struct{}

func (v *RepositoryValidator) Name() string     { return "repository-access" }
func (v *RepositoryValidator) Category() string { return "connectivity" }
func (v *RepositoryValidator) Description() string {
	return "Verify configured repository exists and is accessible"
}
func (v *RepositoryValidator) CanFix() bool { return false }

func (v *RepositoryValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	// Skip if no GitHub config
	if cfg.GitHub.ConfigRepo == "" {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   SKIP,
			Message:  "No GitHub repository configured",
			FixHint:  "Configure GitHub repository in settings.yaml",
			AutoFix:  false,
		}
	}

	repoURL := fmt.Sprintf("https://github.com/%s", cfg.GitHub.ConfigRepo)

	// Check if repository exists (public check)
	result, err := system.RunCommand("curl", "-s", "-f", "-I", repoURL)
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Repository not accessible",
			Details:  []string{fmt.Sprintf("Repository: %s", cfg.GitHub.ConfigRepo), "Public access failed"},
			FixHint:  "Check repository name and visibility settings",
			AutoFix:  false,
		}
	}

	// Test git clone access (dry run)
	gitURL := fmt.Sprintf("https://github.com/%s.git", cfg.GitHub.ConfigRepo)
	result, err = system.RunCommand("git", "ls-remote", gitURL, "HEAD")
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   WARN,
			Message:  "Repository exists but git access limited",
			Details:  []string{fmt.Sprintf("Repository: %s", cfg.GitHub.ConfigRepo), "May require authentication for git operations"},
			FixHint:  "Ensure you have access to the repository",
			AutoFix:  false,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "Repository is accessible",
		Details:  []string{fmt.Sprintf("Repository: %s", cfg.GitHub.ConfigRepo), "Git operations available"},
		AutoFix:  false,
	}
}

func (v *RepositoryValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	return fmt.Errorf("repository access issues must be fixed manually")
}

// GitConnectivityValidator checks if git operations work properly
type GitConnectivityValidator struct{}

func (v *GitConnectivityValidator) Name() string     { return "git-connectivity" }
func (v *GitConnectivityValidator) Category() string { return "connectivity" }
func (v *GitConnectivityValidator) Description() string {
	return "Verify git operations are functional"
}
func (v *GitConnectivityValidator) CanFix() bool { return false }

func (v *GitConnectivityValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	var details []string
	var warnings []string

	// Check if git is available
	result, err := system.RunCommand("git", "--version")
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Git is not available",
			Details:  []string{"git command not found"},
			FixHint:  "Install git using 'anvil install git'",
			AutoFix:  false,
		}
	}
	details = append(details, strings.TrimSpace(result.Output))

	// Check git configuration
	username, err := system.RunCommand("git", "config", "--global", "user.name")
	if err != nil || strings.TrimSpace(username.Output) == "" {
		warnings = append(warnings, "global git username not set")
	} else {
		details = append(details, "Global username: "+strings.TrimSpace(username.Output))
	}

	email, err := system.RunCommand("git", "config", "--global", "user.email")
	if err != nil || strings.TrimSpace(email.Output) == "" {
		warnings = append(warnings, "global git email not set")
	} else {
		details = append(details, "Global email: "+strings.TrimSpace(email.Output))
	}

	// Test basic git operations
	// This is a simple test that doesn't require a repository
	result, err = system.RunCommand("git", "config", "--list")
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Git configuration is not functional",
			Details:  details,
			FixHint:  "Check git installation and configuration",
			AutoFix:  false,
		}
	}

	if len(warnings) > 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   WARN,
			Message:  "Git is functional but configuration incomplete: " + strings.Join(warnings, ", "),
			Details:  details,
			FixHint:  "Configure git username and email globally or in settings.yaml",
			AutoFix:  false,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "Git operations are functional",
		Details:  details,
		AutoFix:  false,
	}
}

func (v *GitConnectivityValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	return fmt.Errorf("git connectivity issues must be fixed manually")
}
