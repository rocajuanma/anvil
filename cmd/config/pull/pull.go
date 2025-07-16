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

package pull

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/github"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull [directory]",
	Short: "Pull configuration files from a specific directory in GitHub repository",
	Long:  constants.PULL_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.ExactArgs(1), // Require exactly one argument (directory name)
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPullCommand(cmd, args); err != nil {
			terminal.PrintError("Pull failed: %v", err)
			return
		}
	},
}

// runPullCommand executes the configuration pull process for a specific directory
func runPullCommand(cmd *cobra.Command, args []string) error {
	// Get the directory to pull
	targetDir := args[0]

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPull, "load-config", err)
	}

	// Validate GitHub configuration
	if err := validateGitHubConfig(cfg); err != nil {
		return err
	}

	terminal.PrintHeader(fmt.Sprintf("Pulling Configuration Directory: %s", targetDir))
	terminal.PrintInfo("Repository: %s", cfg.GitHub.ConfigRepo)
	terminal.PrintInfo("Branch: %s", cfg.GitHub.Branch)
	terminal.PrintInfo("Target directory: %s", targetDir)

	// Get GitHub token from environment variable if configured
	token := ""
	if cfg.GitHub.TokenEnvVar != "" {
		token = os.Getenv(cfg.GitHub.TokenEnvVar)
		if token != "" {
			terminal.PrintInfo("✅ GitHub token found in environment variable: %s", cfg.GitHub.TokenEnvVar)
		} else {
			terminal.PrintWarning("⚠️  No GitHub token found in %s - will attempt SSH authentication", cfg.GitHub.TokenEnvVar)
		}
	}

	// Create GitHub client
	githubClient := github.NewGitHubClient(
		cfg.GitHub.ConfigRepo,
		cfg.GitHub.Branch,
		cfg.GitHub.LocalPath,
		token,
		cfg.Git.SSHKeyPath,
		cfg.Git.Username,
		cfg.Git.Email,
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Check if repository is accessible and branch exists
	terminal.PrintStage("Validating repository access and branch configuration...")
	if err := githubClient.ValidateRepository(ctx); err != nil {
		// Provide additional context for branch configuration errors
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			terminal.PrintError("\n" + err.Error())
			terminal.PrintInfo("\n💡 Quick Help:")
			terminal.PrintInfo("   • Most repositories use 'main' or 'master' as the default branch")
			terminal.PrintInfo("   • Check your repository on GitHub to see available branches")
			terminal.PrintInfo("   • Your settings file is at: %s/settings.yaml", config.GetConfigDirectory())
			return fmt.Errorf("branch configuration validation failed")
		}
		return fmt.Errorf("repository validation failed: %w", err)
	}
	terminal.PrintSuccess("Repository and branch configuration validated")

	// Clone repository if it doesn't exist locally
	terminal.PrintStage("Setting up local repository...")
	if err := githubClient.CloneRepository(ctx); err != nil {
		// Provide additional context for branch configuration errors during clone
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			terminal.PrintError("\n" + err.Error())
			return fmt.Errorf("clone failed due to branch configuration issue")
		}
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	terminal.PrintSuccess("Local repository ready")

	// Pull latest changes
	terminal.PrintStage("Pulling latest changes...")
	if err := githubClient.PullChanges(ctx); err != nil {
		// Provide additional context for branch configuration errors during pull
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			terminal.PrintError("\n" + err.Error())
			terminal.PrintInfo("\n🔄 The local repository exists but the configured branch is not available.")
			terminal.PrintInfo("    You may need to:")
			terminal.PrintInfo("    • Update the branch in your settings.yaml")
			terminal.PrintInfo("    • Or delete the local repository at: %s", cfg.GitHub.LocalPath)
			terminal.PrintInfo("      (It will be re-cloned with the correct branch)")
			return fmt.Errorf("pull failed due to branch configuration issue")
		}
		return fmt.Errorf("failed to pull changes: %w", err)
	}
	terminal.PrintSuccess("Repository updated")

	// Copy specific directory to temp location
	terminal.PrintStage("Copying configuration directory...")
	tempDir, err := copyDirectoryToTemp(cfg, targetDir)
	if err != nil {
		return err
	}
	terminal.PrintSuccess("Configuration directory copied to temp location")

	// Display completion message
	terminal.PrintHeader("Pull Complete!")
	terminal.PrintInfo("Configuration directory '%s' has been pulled from: %s", targetDir, cfg.GitHub.ConfigRepo)
	terminal.PrintInfo("Files are available at: %s", tempDir)

	// List the files that were copied
	if err := listCopiedFiles(tempDir); err == nil {
		// Files listed successfully
	} else {
		terminal.PrintWarning("Could not list copied files: %v", err)
	}

	// Provide next steps
	terminal.PrintInfo("\nNext steps:")
	terminal.PrintInfo("  • Review the pulled configuration files in: %s", tempDir)
	terminal.PrintInfo("  • Apply/copy configurations to their destination as needed")
	terminal.PrintInfo("  • Use 'anvil config push' to upload any local changes")

	return nil
}

// validateGitHubConfig validates that GitHub configuration is properly set up
func validateGitHubConfig(cfg *config.AnvilConfig) error {
	if cfg.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf("github.config_repo is not configured. Please edit %s/settings.yaml and set github.config_repo to your repository (e.g., 'username/dotfiles')",
				config.GetConfigDirectory()))
	}

	if cfg.GitHub.Branch == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf(`github.branch is not configured.

📝 To fix this:
  1. Edit your settings.yaml file at: %s/settings.yaml
  2. Set the 'github.branch' field to your repository's default branch
  3. Common branch names: 'main', 'master', 'develop'
  
Example:
  github:
    branch: "main"  # ← Set this to your repository's default branch`,
				config.GetConfigDirectory()))
	}

	if cfg.GitHub.LocalPath == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf("github.local_path is not configured"))
	}

	// Provide guidance about branch configuration
	terminal.PrintInfo("🔧 Using branch: %s", cfg.GitHub.Branch)
	if cfg.GitHub.Branch != "main" && cfg.GitHub.Branch != "master" {
		terminal.PrintWarning("⚠️  Note: You're using branch '%s'. Make sure this branch exists in your repository.", cfg.GitHub.Branch)
		terminal.PrintInfo("💡 Common default branches are 'main' or 'master'")
	}

	// Check if git is available
	if cfg.Git.Username == "" || cfg.Git.Email == "" {
		terminal.PrintWarning("⚠️  Git user configuration is incomplete. Consider setting git.username and git.email in settings.yaml")
	}

	return nil
}

// copyDirectoryToTemp copies a specific directory from the repo to a temporary location
func copyDirectoryToTemp(cfg *config.AnvilConfig, targetDir string) (string, error) {
	// Source directory in the cloned repo
	sourceDir := filepath.Join(cfg.GitHub.LocalPath, targetDir)

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return "", errors.NewConfigurationError(constants.OpPull, "source-directory",
			fmt.Errorf("directory '%s' does not exist in repository %s", targetDir, cfg.GitHub.ConfigRepo))
	}

	// Create temp directory inside anvil config
	tempBasedir := filepath.Join(cfg.Directories.Config, "temp")
	if err := os.MkdirAll(tempBasedir, constants.DirPerm); err != nil {
		return "", errors.NewFileSystemError(constants.OpPull, "create-temp-dir", err)
	}

	// Destination directory
	destDir := filepath.Join(tempBasedir, targetDir)

	// Remove existing destination if it exists
	if err := os.RemoveAll(destDir); err != nil {
		return "", errors.NewFileSystemError(constants.OpPull, "remove-existing", err)
	}

	// Copy directory recursively
	if err := copyDirRecursive(sourceDir, destDir); err != nil {
		return "", errors.NewFileSystemError(constants.OpPull, "copy-directory", err)
	}

	return destDir, nil
}

// copyDirRecursive recursively copies a directory
func copyDirRecursive(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		} else {
			// Copy file
			return copyFile(path, destPath)
		}
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), constants.DirPerm); err != nil {
		return err
	}

	// Write destination file
	return os.WriteFile(dst, data, constants.FilePerm)
}

// listCopiedFiles lists the files that were copied to the temp directory
func listCopiedFiles(tempDir string) error {
	terminal.PrintInfo("\nCopied files:")

	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only show files
		if !info.IsDir() {
			relPath, err := filepath.Rel(tempDir, path)
			if err != nil {
				relPath = path
			}
			terminal.PrintInfo("  • %s", relPath)
		}

		return nil
	})
}

func init() {
	// Add flags for additional functionality
	PullCmd.Flags().Bool("force", false, "Force pull even if local changes exist")
	PullCmd.Flags().String("branch", "", "Override the branch to pull from")
}
