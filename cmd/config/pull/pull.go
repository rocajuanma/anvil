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

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/github"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/utils"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull [directory]",
	Short: "Pull configuration files from a specific directory in GitHub repository",
	Long:  constants.PULL_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPullCommand(cmd, args); err != nil {
			palantir.GetGlobalOutputHandler().PrintError("Pull failed: %v", err)
			return
		}
	},
}

// runPullCommand executes the configuration pull process for a specific directory
func runPullCommand(cmd *cobra.Command, args []string) error {
	// Default to "anvil" if no argument provided
	targetDir := constants.ANVIL
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpPull, "load-config", err)
	}

	// Validate GitHub configuration
	if err := validateGitHubConfig(cfg); err != nil {
		return err
	}
	output := palantir.GetGlobalOutputHandler()
	output.PrintHeader(fmt.Sprintf("Pulling Configuration Directory: %s", targetDir))
	output.PrintInfo("Repository: %s", cfg.GitHub.ConfigRepo)
	output.PrintInfo("Branch: %s", cfg.GitHub.Branch)
	output.PrintInfo("Target directory: %s", targetDir)
	output.PrintInfo("")

	// Stage 1: Authentication check
	output.PrintStage("Checking authentication...")
	token := ""
	if cfg.GitHub.TokenEnvVar != "" {
		token = os.Getenv(cfg.GitHub.TokenEnvVar)
		if token != "" {
			output.PrintSuccess(fmt.Sprintf("GitHub token found in environment variable: %s", cfg.GitHub.TokenEnvVar))
		} else {
			output.PrintWarning("No GitHub token found in %s - will attempt SSH authentication", cfg.GitHub.TokenEnvVar)
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

	// Stage 2: Repository validation
	output.PrintStage("Stage 2: Validating repository access...")
	spinner := charm.NewCircleSpinner("Validating repository access and branch configuration")
	spinner.Start()
	if err := githubClient.ValidateRepository(ctx); err != nil {
		spinner.Error("Repository validation failed")
		// Provide additional context for repository validation errors
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			fmt.Println("")
			output.PrintError("%s", err.Error())
			fmt.Println("")
			output.PrintInfo("🔄 The repository exists but the configured branch is not available.")
			output.PrintInfo("    You may need to:")
			output.PrintInfo("    • Update the branch in your %s", constants.ANVIL_CONFIG_FILE)
			output.PrintInfo("    • Or check the available branches in your repository")
			return fmt.Errorf("repository validation failed due to branch configuration issue")
		}
		return fmt.Errorf("failed to validate repository: %w", err)
	}
	spinner.Success("Repository access confirmed")

	// Stage 3: Clone/update repository
	output.PrintStage("Stage 3: Cloning or updating repository...")
	spinner = charm.NewDotsSpinner("Cloning or updating repository")
	spinner.Start()
	if err := githubClient.CloneRepository(ctx); err != nil {
		spinner.Error("Clone failed")
		// Provide additional context for clone errors
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			fmt.Println("")
			output.PrintError("%s", err.Error())
			fmt.Println("")
			output.PrintInfo("🔄 The repository exists but the configured branch is not available during clone.")
			output.PrintInfo("    You may need to:")
			output.PrintInfo("    • Update the branch in your %s", constants.ANVIL_CONFIG_FILE)
			output.PrintInfo("    • Or delete the local repository at: %s", cfg.GitHub.LocalPath)
			output.PrintInfo("      (It will be re-cloned with the correct branch)")
			return fmt.Errorf("clone failed due to branch configuration issue")
		}
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	spinner.Success("Repository ready")

	// Stage 4: Pull latest changes
	output.PrintStage("Stage 4: Pulling latest changes...")
	spinner = charm.NewDotsSpinner("Pulling latest changes")
	spinner.Start()
	if err := githubClient.PullChanges(ctx); err != nil {
		spinner.Error("Pull failed")
		// Provide additional context for branch configuration errors during pull
		if strings.Contains(err.Error(), "Branch Configuration Error") {
			output.PrintError("%s", err.Error())
			fmt.Println("")
			output.PrintInfo("🔄 The local repository exists but the configured branch is not available.")
			output.PrintInfo("    You may need to:")
			output.PrintInfo("    • Update the branch in your %s", constants.ANVIL_CONFIG_FILE)
			output.PrintInfo("    • Or delete the local repository at: %s", cfg.GitHub.LocalPath)
			output.PrintInfo("      (It will be re-cloned with the correct branch)")
			return fmt.Errorf("pull failed due to branch configuration issue")
		}
		return fmt.Errorf("failed to pull changes: %w", err)
	}
	spinner.Success("Repository updated")

	// Stage 5: Copy configuration directory
	output.PrintStage("Stage 5: Copying configuration directory...")
	spinner = charm.NewDotsSpinner(fmt.Sprintf("Copying %s directory", targetDir))
	spinner.Start()
	tempDir, err := copyDirectoryToTemp(cfg, targetDir)
	if err != nil {
		spinner.Error("Failed to copy configuration")
		return err
	}
	spinner.Success("Configuration directory copied to temp location")

	// Display completion message
	output.PrintHeader("Pull Complete!")
	output.PrintInfo("Configuration directory '%s' has been pulled from: %s", targetDir, cfg.GitHub.ConfigRepo)
	output.PrintInfo("Files are available at: %s", tempDir)

	// List the files that were copied
	if err := listCopiedFiles(tempDir); err == nil {
		// Files listed successfully
	} else {
		output.PrintWarning("Could not list copied files: %v", err)
	}

	// Provide next steps
	fmt.Println("")
	output.PrintInfo("Next steps:")
	output.PrintInfo("  • Review the pulled configuration files in: %s", tempDir)
	output.PrintInfo("  • Apply/copy configurations to their destination as needed")
	output.PrintInfo("  • Use 'anvil config push' to upload any local changes")

	return nil
}

// validateGitHubConfig validates that GitHub configuration is properly set up
func validateGitHubConfig(cfg *config.AnvilConfig) error {
	if cfg.GitHub.ConfigRepo == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf("github.config_repo is not configured. Please edit %s/%s and set github.config_repo to your repository (e.g., 'username/dotfiles')",
				config.GetAnvilConfigDirectory(), constants.ANVIL_CONFIG_FILE))
	}

	if cfg.GitHub.Branch == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf(`github.branch is not configured.

📝 To fix this:
  1. Edit your %s file at: %s/%s
  2. Set the 'github.branch' field to your repository's default branch
  3. Common branch names: 'main', 'master', 'develop'
  
Example:
  github:
    branch: "main"  # ← Set this to your repository's default branch`,
				constants.ANVIL_CONFIG_FILE, config.GetAnvilConfigDirectory(), constants.ANVIL_CONFIG_FILE))
	}

	if cfg.GitHub.LocalPath == "" {
		return errors.NewConfigurationError(constants.OpPull, "validate-config",
			fmt.Errorf("github.local_path is not configured"))
	}

	output := palantir.GetGlobalOutputHandler()
	// Provide guidance about branch configuration
	output.PrintInfo("🔧 Using branch: %s", cfg.GitHub.Branch)
	if cfg.GitHub.Branch != "main" && cfg.GitHub.Branch != "master" {
		output.PrintWarning("⚠️  Note: You're using branch '%s'. Make sure this branch exists in your repository.", cfg.GitHub.Branch)
		output.PrintInfo("💡 Common default branches are 'main' or 'master'")
	}

	// Check if git is available
	if cfg.Git.Username == "" || cfg.Git.Email == "" {
		output.PrintWarning(fmt.Sprintf("⚠️  Git user configuration is incomplete. Consider setting git.username and git.email in %s", constants.ANVIL_CONFIG_FILE))
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
	tempBasedir := filepath.Join(config.GetAnvilConfigDirectory(), "temp")
	if err := utils.EnsureDirectory(tempBasedir); err != nil {
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

// copyDirRecursive recursively copies a directory using the consolidated utils.CopyDirectorySimple
func copyDirRecursive(src, dst string) error {
	return utils.CopyDirectorySimple(src, dst)
}

// listCopiedFiles lists the files that were copied to the temp directory
func listCopiedFiles(tempDir string) error {
	palantir.GetGlobalOutputHandler().PrintInfo("\nCopied files:")

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
			palantir.GetGlobalOutputHandler().PrintInfo("  • %s", relPath)
		}

		return nil
	})
}

func init() {
	// Add flags for additional functionality
	PullCmd.Flags().Bool("force", false, "Force pull even if local changes exist")
	PullCmd.Flags().String("branch", "", "Override the branch to pull from")
}
