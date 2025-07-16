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

// Package github provides GitHub integration for configuration management
package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/system"
)

// GitHubClient handles GitHub operations for config management
type GitHubClient struct {
	RepoURL    string
	Branch     string
	LocalPath  string
	Token      string
	SSHKeyPath string
	Username   string
	Email      string
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(repoURL, branch, localPath, token, sshKeyPath, username, email string) *GitHubClient {
	return &GitHubClient{
		RepoURL:    repoURL,
		Branch:     branch,
		LocalPath:  localPath,
		Token:      token,
		SSHKeyPath: sshKeyPath,
		Username:   username,
		Email:      email,
	}
}

// CloneRepository clones the repository if it doesn't exist locally
func (gc *GitHubClient) CloneRepository(ctx context.Context) error {
	// Check if local path already exists and is a valid git repository
	if gc.isValidGitRepository() {
		return nil // Repository already exists and is valid
	}

	// Remove any existing directory that might be corrupted
	if err := os.RemoveAll(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPull, "remove-existing", err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(gc.LocalPath), constants.DirPerm); err != nil {
		return errors.NewFileSystemError(constants.OpPull, "mkdir-parent", err)
	}

	// Determine clone URL format (HTTPS with token or SSH)
	cloneURL := gc.getCloneURL()

	// Clone the repository
	args := []string{"clone", "--branch", gc.Branch, cloneURL, gc.LocalPath}
	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, args...)
	if err != nil {
		// Enhanced error message for branch issues
		if strings.Contains(result.Error, "Remote branch") || strings.Contains(result.Error, "not found") {
			return gc.createBranchNotFoundError("clone", result.Error)
		}
		return errors.NewInstallationError(constants.OpPull, "git-clone",
			fmt.Errorf("failed to clone repository: %s, error: %w", result.Error, err))
	}

	// Verify the repository was cloned successfully
	if !gc.isValidGitRepository() {
		return errors.NewInstallationError(constants.OpPull, "verify-clone",
			fmt.Errorf("repository clone completed but directory is not a valid git repository: %s", gc.LocalPath))
	}

	return nil
}

// PullChanges pulls the latest changes from the remote repository
func (gc *GitHubClient) PullChanges(ctx context.Context) error {
	// Verify the repository exists and is valid
	if !gc.isValidGitRepository() {
		return errors.NewFileSystemError(constants.OpPull, "invalid-repo",
			fmt.Errorf("local repository at %s is not valid or doesn't exist", gc.LocalPath))
	}

	// Ensure we're in the correct directory
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPull, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPull, "chdir",
			fmt.Errorf("cannot change to repository directory %s: %w", gc.LocalPath, err))
	}

	// Configure git user if provided
	if err := gc.configureGitUser(ctx); err != nil {
		return err
	}

	// Fetch latest changes
	fetchResult, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "fetch", "origin", gc.Branch)
	if err != nil {
		// Enhanced error message for branch issues during fetch
		if strings.Contains(fetchResult.Error, "couldn't find remote ref") || strings.Contains(fetchResult.Error, "not found") {
			return gc.createBranchNotFoundError("fetch", fetchResult.Error)
		}
		return errors.NewInstallationError(constants.OpPull, "git-fetch",
			fmt.Errorf("failed to fetch changes: %s, error: %w", fetchResult.Error, err))
	}

	// Pull changes
	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "pull", "origin", gc.Branch)
	if err != nil {
		// Enhanced error message for branch issues during pull
		if strings.Contains(result.Error, "couldn't find remote ref") || strings.Contains(result.Error, "not found") {
			return gc.createBranchNotFoundError("pull", result.Error)
		}
		return errors.NewInstallationError(constants.OpPull, "git-pull",
			fmt.Errorf("failed to pull changes: %s, error: %w", result.Error, err))
	}

	return nil
}

// PushChanges commits and pushes local changes to the remote repository
func (gc *GitHubClient) PushChanges(ctx context.Context, commitMessage string) error {
	// Ensure we're in the correct directory
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Configure git user if provided
	if err := gc.configureGitUser(ctx); err != nil {
		return err
	}

	// Add all changes
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "add", "."); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-add", err)
	}

	// Check if there are changes to commit
	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "diff", "--cached", "--exit-code")
	if err == nil {
		// No changes to commit
		return nil
	}

	// Commit changes
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "commit", "-m", commitMessage); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-commit", err)
	}

	// Push changes
	result, err = system.RunCommandWithTimeout(ctx, constants.GitCommand, "push", "origin", gc.Branch)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-push",
			fmt.Errorf("failed to push changes: %s, error: %w", result.Error, err))
	}

	return nil
}

// CreateRepository creates a new GitHub repository if it doesn't exist
func (gc *GitHubClient) CreateRepository(ctx context.Context, repoName, description string) error {
	// This would require GitHub API integration
	// For now, we'll assume the repository exists or provide instructions
	return fmt.Errorf("repository creation not implemented - please create the repository manually on GitHub: %s", gc.RepoURL)
}

// ValidateRepository checks if the repository is accessible and the specified branch exists
func (gc *GitHubClient) ValidateRepository(ctx context.Context) error {
	// First, try to fetch repository information
	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "ls-remote", gc.getCloneURL(), "HEAD")
	if err != nil {
		return errors.NewNetworkError(constants.OpConfig, "git-ls-remote",
			fmt.Errorf("cannot access repository %s: %s, error: %w", gc.RepoURL, result.Error, err))
	}

	// Check if the specified branch exists
	branchResult, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "ls-remote", "--heads", gc.getCloneURL(), gc.Branch)
	if err != nil {
		return errors.NewNetworkError(constants.OpConfig, "git-ls-remote-branch",
			fmt.Errorf("failed to check branch %s in repository %s: %s, error: %w", gc.Branch, gc.RepoURL, branchResult.Error, err))
	}

	// If the branch result is empty, the branch doesn't exist
	if strings.TrimSpace(branchResult.Output) == "" {
		return gc.createBranchNotFoundError("validation", fmt.Sprintf("branch '%s' not found in remote repository", gc.Branch))
	}

	return nil
}

// getCloneURL returns the appropriate clone URL based on available authentication
func (gc *GitHubClient) getCloneURL() string {
	if gc.Token != "" {
		// Use HTTPS with token
		if strings.HasPrefix(gc.RepoURL, "https://") {
			return strings.Replace(gc.RepoURL, "https://", fmt.Sprintf("https://%s@", gc.Token), 1)
		}
		// Convert repo format like "username/repo" to HTTPS with token
		if !strings.Contains(gc.RepoURL, "://") {
			return fmt.Sprintf("https://%s@github.com/%s.git", gc.Token, gc.RepoURL)
		}
	}

	// Use SSH if available
	if gc.SSHKeyPath != "" {
		if _, err := os.Stat(gc.SSHKeyPath); err == nil {
			// Convert to SSH format
			if strings.HasPrefix(gc.RepoURL, "https://github.com/") {
				repoPath := strings.TrimPrefix(gc.RepoURL, "https://github.com/")
				repoPath = strings.TrimSuffix(repoPath, ".git")
				return fmt.Sprintf("git@github.com:%s.git", repoPath)
			}
			if !strings.Contains(gc.RepoURL, "://") {
				return fmt.Sprintf("git@github.com:%s.git", gc.RepoURL)
			}
		}
	}

	// Default to HTTPS
	if !strings.Contains(gc.RepoURL, "://") {
		return fmt.Sprintf("https://github.com/%s.git", gc.RepoURL)
	}
	return gc.RepoURL
}

// configureGitUser configures git user for the repository
func (gc *GitHubClient) configureGitUser(ctx context.Context) error {
	if gc.Username != "" {
		if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "config", "user.name", gc.Username); err != nil {
			return errors.NewConfigurationError(constants.OpConfig, "git-config-user", err)
		}
	}

	if gc.Email != "" {
		if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "config", "user.email", gc.Email); err != nil {
			return errors.NewConfigurationError(constants.OpConfig, "git-config-email", err)
		}
	}

	return nil
}

// GetRepositoryStatus returns the current status of the local repository
func (gc *GitHubClient) GetRepositoryStatus(ctx context.Context) (string, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return "", errors.NewFileSystemError(constants.OpConfig, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return "", errors.NewFileSystemError(constants.OpConfig, "chdir", err)
	}

	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "status", "--porcelain")
	if err != nil {
		return "", errors.NewInstallationError(constants.OpConfig, "git-status", err)
	}

	return result.Output, nil
}

// isValidGitRepository checks if the local path contains a valid git repository
func (gc *GitHubClient) isValidGitRepository() bool {
	// Check if directory exists
	if _, err := os.Stat(gc.LocalPath); os.IsNotExist(err) {
		return false
	}

	// Check if .git directory exists
	gitDir := filepath.Join(gc.LocalPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false
	}

	// Try to run a simple git command to verify it's a valid repo
	originalDir, err := os.Getwd()
	if err != nil {
		return false
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return false
	}

	// Run git status to verify it's a valid repository
	_, err = system.RunCommand(constants.GitCommand, "status", "--porcelain")
	return err == nil
}

// createBranchNotFoundError creates a detailed error message when a branch is not found
func (gc *GitHubClient) createBranchNotFoundError(operation, gitError string) error {
	availableBranches := gc.getAvailableBranches()

	errorMsg := fmt.Sprintf(`
❌ Branch Configuration Error

The branch '%s' does not exist in repository '%s'.

Git error from %s operation: %s

🔍 IMPORTANT: Check your branch configuration in settings.yaml!

Current configuration:
  - Repository: %s
  - Branch: %s

%s

📝 To fix this issue:
  1. Edit your settings.yaml file (usually at ~/.anvil/settings.yaml)
  2. Update the 'github.branch' field to match an existing branch
  3. Or create the branch '%s' in your repository
  4. Save the file and try the pull command again

Example settings.yaml section:
  github:
    config_repo: "%s"
    branch: "main"  # ← Update this to an existing branch
    local_path: "~/.anvil/repo"`,
		gc.Branch, gc.RepoURL, operation, gitError,
		gc.RepoURL, gc.Branch,
		availableBranches,
		gc.Branch, gc.RepoURL)

	return errors.NewConfigurationError(constants.OpPull, "branch-not-found", fmt.Errorf(errorMsg))
}

// getAvailableBranches attempts to list available branches from the remote repository
func (gc *GitHubClient) getAvailableBranches() string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "ls-remote", "--heads", gc.getCloneURL())
	if err != nil {
		return "\n⚠️  Could not retrieve available branches. Check repository access."
	}

	if result.Output == "" {
		return "\n⚠️  No branches found in the repository."
	}

	lines := strings.Split(strings.TrimSpace(result.Output), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		// Extract branch name from "commit_hash refs/heads/branch_name"
		parts := strings.Fields(line)
		if len(parts) >= 2 && strings.HasPrefix(parts[1], "refs/heads/") {
			branchName := strings.TrimPrefix(parts[1], "refs/heads/")
			branches = append(branches, branchName)
		}
	}

	if len(branches) == 0 {
		return "\n⚠️  Could not parse available branches."
	}

	branchList := strings.Join(branches, "\n    - ")
	return fmt.Sprintf("\n✅ Available branches in repository:\n    - %s", branchList)
}
