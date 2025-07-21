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

package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
)

// PushConfigResult represents the result of a config push operation
type PushConfigResult struct {
	BranchName     string
	CommitMessage  string
	RepositoryURL  string
	FilesCommitted []string
}

// verifyRepositoryPrivacy ensures the repository is private before allowing push operations
func (gc *GitHubClient) verifyRepositoryPrivacy(ctx context.Context) error {
	// First test git access (this should work for private repos with proper auth)
	gitURL := fmt.Sprintf("https://github.com/%s.git", gc.RepoURL)
	result, err := system.RunCommandWithTimeout(ctx, "git", "ls-remote", gitURL, "HEAD")

	if err != nil || !result.Success {
		return fmt.Errorf("🚨 SECURITY BLOCK: Cannot verify repository privacy - authentication failed\n"+
			"Repository: %s\n"+
			"Anvil REQUIRES private repositories for configuration data.\n"+
			"Configure proper authentication (GITHUB_TOKEN or SSH keys) before pushing", gc.RepoURL)
	}

	// Test if repository is publicly accessible (this should FAIL for private repos)
	repoURL := fmt.Sprintf("https://github.com/%s", gc.RepoURL)
	httpResult, httpErr := system.RunCommandWithTimeout(ctx, "curl", "-s", "-f", "-I", repoURL)

	if httpErr == nil && httpResult.Success {
		// 🚨 CRITICAL: Repository is public - BLOCK the push
		terminal.PrintError("🚨 SECURITY VIOLATION: Configuration push BLOCKED")
		terminal.PrintError("")
		terminal.PrintError("Repository '%s' is PUBLIC", gc.RepoURL)
		terminal.PrintError("❌ Configuration files contain sensitive data")
		terminal.PrintError("❌ PUBLIC repositories expose API keys, paths, and personal information")
		terminal.PrintError("❌ This could lead to security breaches and data leaks")
		terminal.PrintError("")
		terminal.PrintError("🔒 REQUIRED ACTION: Make repository PRIVATE")
		terminal.PrintError("   Visit: https://github.com/%s/settings", gc.RepoURL)
		terminal.PrintError("   Go to: Danger Zone → Change repository visibility → Private")
		terminal.PrintError("")
		terminal.PrintError("🛡️  Anvil will NEVER push configuration data to public repositories")

		return fmt.Errorf("SECURITY BLOCK: Repository is public. Configuration push denied for security")
	}

	// Repository appears to be private and git access works - safe to proceed
	terminal.PrintSuccess("🔒 Repository privacy verified - safe to push configuration data")
	return nil
}

// PushAnvilConfig pushes the anvil settings.yaml to the repository
func (gc *GitHubClient) PushAnvilConfig(ctx context.Context, settingsPath string) (*PushConfigResult, error) {
	// 🚨 CRITICAL SECURITY CHECK: Verify repository is private before ANY push operations
	if err := gc.verifyRepositoryPrivacy(ctx); err != nil {
		return nil, err
	}

	// Ensure repository is ready
	if err := gc.ensureRepositoryReady(ctx); err != nil {
		return nil, err
	}

	// Check if there are differences before proceeding
	hasChanges, err := gc.hasConfigChanges(settingsPath, "anvil/settings.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to check for config changes: %w", err)
	}

	if !hasChanges {
		terminal.PrintSuccess("Configuration is up-to-date!")
		terminal.PrintInfo("Local anvil settings match the remote repository.")
		terminal.PrintInfo("No changes to push.")
		return nil, nil
	}

	terminal.PrintInfo("Differences detected between local and remote configuration")

	// Generate branch name with timestamp
	branchName := generateTimestampedBranchName("config-push")

	// Create and checkout new branch
	if err := gc.createAndCheckoutBranch(ctx, branchName); err != nil {
		return nil, err
	}

	// Copy anvil settings to repo
	targetDir := filepath.Join(gc.LocalPath, "anvil")
	if err := os.MkdirAll(targetDir, constants.DirPerm); err != nil {
		return nil, errors.NewFileSystemError(constants.OpPush, "mkdir-anvil", err)
	}

	targetFile := filepath.Join(targetDir, "settings.yaml")
	if err := copyFile(settingsPath, targetFile); err != nil {
		return nil, errors.NewFileSystemError(constants.OpPush, "copy-settings", err)
	}

	// Commit changes
	commitMessage := "anvil[push]: anvil"
	if err := gc.commitChanges(ctx, commitMessage); err != nil {
		return nil, err
	}

	// Push branch
	if err := gc.pushBranch(ctx, branchName); err != nil {
		return nil, err
	}

	result := &PushConfigResult{
		BranchName:     branchName,
		CommitMessage:  commitMessage,
		RepositoryURL:  gc.getRepositoryURL(),
		FilesCommitted: []string{"anvil/settings.yaml"},
	}

	return result, nil
}

// PushAppConfig pushes application configuration files to the repository
func (gc *GitHubClient) PushAppConfig(ctx context.Context, appName, configPath string) (*PushConfigResult, error) {
	// 🚨 CRITICAL SECURITY CHECK: Verify repository is private before ANY push operations
	if err := gc.verifyRepositoryPrivacy(ctx); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("application config push is still in development")
}

// ensureRepositoryReady ensures the repository is cloned and up to date
func (gc *GitHubClient) ensureRepositoryReady(ctx context.Context) error {
	// Clone repository if it doesn't exist
	if err := gc.CloneRepository(ctx); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Switch back to main branch and pull latest changes
	if err := gc.switchToMainBranch(ctx); err != nil {
		return fmt.Errorf("failed to switch to main branch: %w", err)
	}

	if err := gc.PullChanges(ctx); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	return nil
}

// switchToMainBranch switches to the main branch specified in config
func (gc *GitHubClient) switchToMainBranch(ctx context.Context) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Checkout main branch
	_, err = system.RunCommandWithTimeout(ctx, constants.GitCommand, "checkout", gc.Branch)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-checkout-main", err)
	}

	return nil
}

// createAndCheckoutBranch creates a new branch and checks it out
func (gc *GitHubClient) createAndCheckoutBranch(ctx context.Context, branchName string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Create and checkout new branch
	_, err = system.RunCommandWithTimeout(ctx, constants.GitCommand, "checkout", "-b", branchName)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-checkout-new-branch", err)
	}

	terminal.PrintInfo("Created and switched to branch: %s", branchName)
	return nil
}

// commitChanges adds and commits all changes in the repository
func (gc *GitHubClient) commitChanges(ctx context.Context, commitMessage string) error {
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
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-diff-check", err)
	}

	if result.ExitCode == 0 {
		// Exit code 0 means no differences
		return fmt.Errorf("no changes to commit")
	}

	// Exit code 1 means there are differences - proceed with commit
	terminal.PrintInfo("Changes detected, proceeding with commit...")

	// Commit changes
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "commit", "-m", commitMessage); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-commit", err)
	}

	terminal.PrintSuccess(fmt.Sprintf("Committed changes: %s", commitMessage))
	return nil
}

// pushBranch pushes the current branch to origin
func (gc *GitHubClient) pushBranch(ctx context.Context, branchName string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Push branch to origin
	result, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "push", "--set-upstream", "origin", branchName)
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-push",
			fmt.Errorf("failed to push branch: %s, error: %w", result.Error, err))
	}

	terminal.PrintSuccess(fmt.Sprintf("Pushed branch '%s' to origin", branchName))
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// generateTimestampedBranchName generates a branch name with current date and time
func generateTimestampedBranchName(prefix string) string {
	now := time.Now()
	dateStr := now.Format("02012006") // DDMMYYYY
	timeStr := now.Format("1504")     // HHMM (24h format)
	return fmt.Sprintf("%s-%s-%s", prefix, dateStr, timeStr)
}

// hasConfigChanges checks if there are differences between local and remote config files
func (gc *GitHubClient) hasConfigChanges(localFilePath, repoRelativePath string) (bool, error) {
	repoFilePath := filepath.Join(gc.LocalPath, repoRelativePath)

	// Check if the remote file exists
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		// Remote file doesn't exist, so we have changes to push
		terminal.PrintInfo("Remote file does not exist, will create new file")
		return true, nil
	}

	// Read local file
	localContent, err := os.ReadFile(localFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to read local file %s: %w", localFilePath, err)
	}

	// Read remote file
	remoteContent, err := os.ReadFile(repoFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to read remote file %s: %w", repoFilePath, err)
	}

	// Compare file contents
	areEqual := string(localContent) == string(remoteContent)
	return !areEqual, nil
}

// getRepositoryURL returns the GitHub repository URL for display
func (gc *GitHubClient) getRepositoryURL() string {
	if strings.Contains(gc.RepoURL, "://") {
		return gc.RepoURL
	}
	return fmt.Sprintf("https://github.com/%s", gc.RepoURL)
}
