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
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/interfaces"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/rocajuanma/anvil/pkg/utils"
)

// PushConfigResult represents the result of a config push operation
type PushConfigResult struct {
	BranchName     string
	CommitMessage  string
	RepositoryURL  string
	FilesCommitted []string
}

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() interfaces.OutputHandler {
	return terminal.GetGlobalOutputHandler()
}

// verifyRepositoryPrivacy ensures the repository is private before allowing push operations
func (gc *GitHubClient) verifyRepositoryPrivacy(ctx context.Context) error {
	// First test git access using the client's authentication method
	authenticatedURL := gc.getCloneURL()
	result, err := system.RunCommandWithTimeout(ctx, "git", "ls-remote", authenticatedURL, "HEAD")

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
		output := getOutputHandler()
		output.PrintError("🚨 SECURITY VIOLATION: Configuration push BLOCKED")
		output.PrintError("")
		output.PrintError("Repository '%s' is PUBLIC", gc.RepoURL)
		output.PrintError("❌ Configuration files contain sensitive data")
		output.PrintError("❌ PUBLIC repositories expose API keys, paths, and personal information")
		output.PrintError("❌ This could lead to security breaches and data leaks")
		output.PrintError("")
		output.PrintError("🔒 REQUIRED ACTION: Make repository PRIVATE")
		output.PrintError("   Visit: https://github.com/%s/settings", gc.RepoURL)
		output.PrintError("   Go to: Danger Zone → Change repository visibility → Private")
		output.PrintError("")
		output.PrintError("🛡️  Anvil will NEVER push configuration data to public repositories")

		return fmt.Errorf("SECURITY BLOCK: Repository is public. Configuration push denied for security")
	}

	// Repository appears to be private and git access works - safe to proceed
	getOutputHandler().PrintSuccess("🔒 Repository privacy verified - safe to push configuration data")
	return nil
}

// PushConfig pushes configuration files to the repository (unified function for both anvil and app configs)
func (gc *GitHubClient) PushConfig(ctx context.Context, appName, configPath string) (*PushConfigResult, error) {
	// 🚨 CRITICAL SECURITY CHECK: Verify repository is private before ANY push operations
	if err := gc.verifyRepositoryPrivacy(ctx); err != nil {
		return nil, err
	}

	// Ensure repository is ready
	if err := gc.ensureRepositoryReady(ctx); err != nil {
		return nil, err
	}

	// Check if there are differences before proceeding
	targetPath := fmt.Sprintf("%s/", appName) // App configs go in a directory named after the app

	// Get the output handler
	output := getOutputHandler()

	// For new apps, we need to check if the target directory exists in the repo
	repoTargetPath := filepath.Join(gc.LocalPath, targetPath)
	if _, err := os.Stat(repoTargetPath); os.IsNotExist(err) {
		// Target doesn't exist in repo - this is a new app
		// Verify the local path actually exists and has content
		if localInfo, err := os.Stat(configPath); err == nil {
			if localInfo.IsDir() {
				// Check if directory has files
				entries, err := os.ReadDir(configPath)
				if err == nil && len(entries) > 0 {
					output.PrintInfo("New app '%s' detected - will be added to repository", appName)
				} else {
					output.PrintSuccess("Configuration is up-to-date!")
					output.PrintInfo("Local %s configs match the remote repository.", appName)
					output.PrintInfo("No changes to push.")
					return nil, nil
				}
			} else if localInfo.Size() > 0 {
				output.PrintInfo("New app '%s' detected - will be added to repository", appName)
			} else {
				output.PrintSuccess("Configuration is up-to-date!")
				output.PrintInfo("Local %s configs match the remote repository.", appName)
				output.PrintInfo("No changes to push.")
				return nil, nil
			}
		} else {
			return nil, fmt.Errorf("local config path is invalid: %w", err)
		}
	} else {
		// Target exists in repo - check for changes
		hasChanges, err := gc.hasAppConfigChanges(configPath, targetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check for config changes: %w", err)
		}

		if !hasChanges {
			output.PrintSuccess("Configuration is up-to-date!")
			output.PrintInfo("Local %s configs match the remote repository.", appName)
			output.PrintInfo("No changes to push.")
			return nil, nil
		}
	}

	output.PrintInfo("Differences detected between local and remote %s configuration", appName)

	// Generate branch name with timestamp
	branchName := generateTimestampedBranchName("config-push")

	// Create and checkout new branch
	if err := gc.createAndCheckoutBranch(ctx, branchName); err != nil {
		return nil, err
	}

	// Copy configs to repo
	targetDir := filepath.Join(gc.LocalPath, appName)
	if err := os.MkdirAll(targetDir, constants.DirPerm); err != nil {
		return nil, errors.NewFileSystemError(constants.OpPush, "mkdir-app", err)
	}

	// Copy the config path (file or directory) to the target directory
	if err := gc.copyConfigToRepo(configPath, targetDir); err != nil {
		return nil, err
	}

	// Commit changes
	commitMessage := fmt.Sprintf("anvil[push]: %s", appName)
	if err := gc.commitChanges(ctx, commitMessage); err != nil {
		return nil, err
	}

	// Push branch
	if err := gc.pushBranch(ctx, branchName); err != nil {
		return nil, err
	}

	// Determine files committed
	filesCommitted, err := gc.getCommittedFiles(targetDir, appName)
	if err != nil {
		filesCommitted = []string{fmt.Sprintf("%s/", appName)} // Fallback
	}

	result := &PushConfigResult{
		BranchName:     branchName,
		CommitMessage:  commitMessage,
		RepositoryURL:  gc.getRepositoryURL(),
		FilesCommitted: filesCommitted,
	}

	return result, nil
}

// PushAppConfig is a wrapper for backwards compatibility - delegates to unified PushConfig
func (gc *GitHubClient) PushAppConfig(ctx context.Context, appName, configPath string) (*PushConfigResult, error) {
	return gc.PushConfig(ctx, appName, configPath)
}

// PushAnvilConfig is a wrapper for backwards compatibility - delegates to unified PushConfig
func (gc *GitHubClient) PushAnvilConfig(ctx context.Context, settingsPath string) (*PushConfigResult, error) {
	return gc.PushConfig(ctx, "anvil", settingsPath)
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

	// Ensure repository is in a clean state before starting push operations
	if err := gc.ensureCleanState(ctx); err != nil {
		return fmt.Errorf("failed to ensure clean repository state: %w", err)
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

	getOutputHandler().PrintInfo("Created and switched to branch: %s", branchName)
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
	getOutputHandler().PrintInfo("Changes detected, proceeding with commit...")

	// Commit changes
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "commit", "-m", commitMessage); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-commit", err)
	}

	getOutputHandler().PrintSuccess(fmt.Sprintf("Committed changes: %s", commitMessage))
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

	getOutputHandler().PrintSuccess(fmt.Sprintf("Pushed branch '%s' to origin", branchName))
	return nil
}

// copyFile copies a file from src to dst using the consolidated utils.CopyFileSimple
func copyFile(src, dst string) error {
	return utils.CopyFileSimple(src, dst)
}

// generateTimestampedBranchName generates a branch name with current date and time
func generateTimestampedBranchName(prefix string) string {
	now := time.Now()
	dateStr := now.Format("02012006") // DDMMYYYY
	timeStr := now.Format("1504")     // HHMM (24h format)
	return fmt.Sprintf("%s-%s-%s", prefix, dateStr, timeStr)
}

// getRepositoryURL returns the GitHub repository URL for display
func (gc *GitHubClient) getRepositoryURL() string {
	if strings.Contains(gc.RepoURL, "://") {
		return gc.RepoURL
	}
	return fmt.Sprintf("https://github.com/%s", gc.RepoURL)
}

// hasAppConfigChanges checks if the local app config differs from the remote
func (gc *GitHubClient) hasAppConfigChanges(localConfigPath, targetPath string) (bool, error) {
	// Check if the target directory exists in the repo
	repoTargetPath := filepath.Join(gc.LocalPath, targetPath)

	// If target doesn't exist in repo, this is a new app
	if _, err := os.Stat(repoTargetPath); os.IsNotExist(err) {
		// Verify the local path actually exists and has content
		if localInfo, err := os.Stat(localConfigPath); err == nil {
			if localInfo.IsDir() {
				// Check if directory has files
				entries, err := os.ReadDir(localConfigPath)
				if err == nil && len(entries) > 0 {
					return true, nil // New app with content
				}
			} else if localInfo.Size() > 0 {
				return true, nil // New file with content
			}
		}
		return false, fmt.Errorf("local config path is empty or invalid")
	}

	// Compare the local config with the repo version
	return gc.hasFileOrDirChanges(localConfigPath, repoTargetPath)
}

// hasFileOrDirChanges compares a local file or directory with a repo version
func (gc *GitHubClient) hasFileOrDirChanges(localPath, repoPath string) (bool, error) {
	localInfo, err := os.Stat(localPath)
	if err != nil {
		return false, fmt.Errorf("failed to stat local path %s: %w", localPath, err)
	}

	repoInfo, err := os.Stat(repoPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // Repo version doesn't exist, so there are changes
		}
		return false, fmt.Errorf("failed to stat repo path %s: %w", repoPath, err)
	}

	// If one is a file and the other is a directory, there are changes
	if localInfo.IsDir() != repoInfo.IsDir() {
		return true, nil
	}

	if localInfo.IsDir() {
		// Compare directories recursively
		return gc.hasDirectoryChanges(localPath, repoPath)
	} else {
		// Compare files
		return gc.hasFileChanges(localPath, repoPath)
	}
}

// hasDirectoryChanges recursively compares two directories
func (gc *GitHubClient) hasDirectoryChanges(localDir, repoDir string) (bool, error) {
	// Get all files in both directories
	localFiles := make(map[string]os.FileInfo)
	repoFiles := make(map[string]os.FileInfo)

	// Walk local directory
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}
		localFiles[relPath] = info
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to walk local directory: %w", err)
	}

	// Walk repo directory
	err = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(repoDir, path)
		if err != nil {
			return err
		}
		repoFiles[relPath] = info
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to walk repo directory: %w", err)
	}

	// Check if file lists differ
	if len(localFiles) != len(repoFiles) {
		return true, nil
	}

	// Compare each file
	for relPath, localInfo := range localFiles {
		_, exists := repoFiles[relPath]
		if !exists {
			return true, nil
		}

		// Skip directories for content comparison
		if localInfo.IsDir() {
			continue
		}

		// Compare file contents
		localFilePath := filepath.Join(localDir, relPath)
		repoFilePath := filepath.Join(repoDir, relPath)
		hasChanges, err := gc.hasFileChanges(localFilePath, repoFilePath)
		if err != nil {
			return false, err
		}
		if hasChanges {
			return true, nil
		}
	}

	return false, nil
}

// hasFileChanges compares two files for differences
func (gc *GitHubClient) hasFileChanges(localFile, repoFile string) (bool, error) {
	localContent, err := os.ReadFile(localFile)
	if err != nil {
		return false, fmt.Errorf("failed to read local file %s: %w", localFile, err)
	}

	repoContent, err := os.ReadFile(repoFile)
	if err != nil {
		return false, fmt.Errorf("failed to read repo file %s: %w", repoFile, err)
	}

	return !bytes.Equal(localContent, repoContent), nil
}

// copyConfigToRepo copies a file or directory to the repository
func (gc *GitHubClient) copyConfigToRepo(sourcePath, targetDir string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat source path %s: %w", sourcePath, err)
	}

	if sourceInfo.IsDir() {
		// Copy directory contents to target directory
		return gc.copyDirectoryContents(sourcePath, targetDir)
	} else {
		// Copy single file to target directory
		fileName := filepath.Base(sourcePath)
		targetFile := filepath.Join(targetDir, fileName)
		return copyFile(sourcePath, targetFile)
	}
}

// copyDirectoryContents recursively copies directory contents
func (gc *GitHubClient) copyDirectoryContents(sourceDir, targetDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			// Copy file
			return copyFile(path, targetPath)
		}
	})
}

// getCommittedFiles returns a list of files that were committed in the target directory
func (gc *GitHubClient) getCommittedFiles(targetDir, appName string) ([]string, error) {
	var files []string

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Get relative path from the repo root
			relPath, err := filepath.Rel(gc.LocalPath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk target directory: %w", err)
	}

	if len(files) == 0 {
		// Fallback to just showing the app directory
		files = []string{fmt.Sprintf("%s/", appName)}
	}

	return files, nil
}

// isSingleSmallFile determines if we should include the full diff output
func (gc *GitHubClient) isSingleSmallFile(statOutput string) bool {
	// Only get full diff for single files with reasonable size
	return strings.Contains(statOutput, "1 file changed") &&
		strings.Count(statOutput, "+")+strings.Count(statOutput, "-") <= 50
}

// extractFileCount parses the file count from Git's stat output
func (gc *GitHubClient) extractFileCount(statOutput string) int {
	if strings.TrimSpace(statOutput) == "" {
		return 0
	}

	// Parse "1 file changed" or "2 files changed"
	if strings.Contains(statOutput, "1 file changed") {
		return 1
	}

	// Use regex to extract number from "X files changed"
	re := regexp.MustCompile(`(\d+) files changed`)
	matches := re.FindStringSubmatch(statOutput)
	if len(matches) >= 2 {
		if count, err := strconv.Atoi(matches[1]); err == nil {
			return count
		}
	}

	return 0
}

// ensureCleanState ensures the repository is in a clean state before push operations
func (gc *GitHubClient) ensureCleanState(ctx context.Context) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Check if there are any staged changes
	stagedResult, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "diff", "--cached", "--exit-code")
	if err != nil && stagedResult.ExitCode != 0 {
		// There are staged changes, reset them
		if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "reset", "HEAD"); err != nil {
			return errors.NewInstallationError(constants.OpPush, "git-reset", err)
		}
	}

	// Check if there are any untracked files
	statusResult, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "status", "--porcelain")
	if err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-status", err)
	}

	// If there are untracked files, clean them
	if strings.TrimSpace(statusResult.Output) != "" {
		if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "clean", "-fd"); err != nil {
			return errors.NewInstallationError(constants.OpPush, "git-clean", err)
		}
	}

	return nil
}

// CleanupStagedChanges removes any staged changes from the repository
// This is called when a push operation is cancelled to ensure clean state
func (gc *GitHubClient) CleanupStagedChanges(ctx context.Context) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return errors.NewFileSystemError(constants.OpPush, "getwd", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(gc.LocalPath); err != nil {
		return errors.NewFileSystemError(constants.OpPush, "chdir", err)
	}

	// Reset any staged changes
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "reset", "HEAD"); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-reset", err)
	}

	// Clean any untracked files that might have been created during diff preview
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "clean", "-fd"); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-clean", err)
	}

	// Switch back to main branch to ensure we're in a clean state
	if _, err := system.RunCommandWithTimeout(ctx, constants.GitCommand, "checkout", gc.Branch); err != nil {
		return errors.NewInstallationError(constants.OpPush, "git-checkout-main", err)
	}

	return nil
}
