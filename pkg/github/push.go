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

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"gopkg.in/yaml.v2"
)

// PushConfigResult represents the result of a config push operation
type PushConfigResult struct {
	BranchName     string
	CommitMessage  string
	RepositoryURL  string
	FilesCommitted []string
}

// PushAnvilConfig pushes the anvil settings.yaml to the repository
func (gc *GitHubClient) PushAnvilConfig(ctx context.Context, settingsPath string) (*PushConfigResult, error) {
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

	// Apply smart filtering before copying
	if err := gc.copyFilteredSettings(settingsPath, targetFile); err != nil {
		return nil, errors.NewFileSystemError(constants.OpPush, "copy-filtered-settings", err)
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

// copyFilteredSettings loads, filters, and saves configuration with smart filtering
func (gc *GitHubClient) copyFilteredSettings(sourcePath, targetPath string) error {
	// Load the original configuration
	anvilConfig, err := config.LoadConfigFromPath(sourcePath)
	if err != nil {
		// If loading with filtering fails, fall back to direct copy for backward compatibility
		terminal.PrintWarning("Unable to load config for filtering, using direct copy")
		return copyFile(sourcePath, targetPath)
	}

	// Apply filtering if sync config is present
	var filteredConfig *config.AnvilConfig
	if len(anvilConfig.SyncConfig.ExcludeSections) > 0 || len(anvilConfig.SyncConfig.TemplateSections) > 0 {
		terminal.PrintInfo("Applying smart filtering based on sync configuration...")

		filteredConfig, err = config.FilterForSync(anvilConfig)
		if err != nil {
			terminal.PrintWarning("Filtering failed, using original config: %v", err)
			filteredConfig = anvilConfig
		} else {
			// Log what was filtered
			if len(anvilConfig.SyncConfig.ExcludeSections) > 0 {
				terminal.PrintInfo("Excluded sections: %s", strings.Join(anvilConfig.SyncConfig.ExcludeSections, ", "))
			}
			if len(anvilConfig.SyncConfig.TemplateSections) > 0 {
				terminal.PrintInfo("Templated sections: %s", strings.Join(anvilConfig.SyncConfig.TemplateSections, ", "))
			}
		}
	} else {
		// No filtering configured, use original config
		filteredConfig = anvilConfig
	}

	// Write filtered configuration to target file
	filteredData, err := yaml.Marshal(filteredConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal filtered config: %w", err)
	}

	if err := os.WriteFile(targetPath, filteredData, constants.FilePerm); err != nil {
		return fmt.Errorf("failed to write filtered config: %w", err)
	}

	return nil
}
