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

package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [directory]",
	Short: "Sync configuration state with system reality",
	Long:  constants.SYNC_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1), // Accept 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSyncCommand(cmd, args); err != nil {
			terminal.PrintError("Sync failed: %v", err)
			return
		}
	},
}

// runSyncCommand executes the configuration sync process
func runSyncCommand(cmd *cobra.Command, args []string) error {
	// Check for dry-run flag
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// If no arguments provided, sync the anvil settings
	if len(args) == 0 {
		return syncAnvilSettings(dryRun)
	}

	// Show specific app config sync (placeholder)
	targetDir := args[0]
	return syncAppConfig(targetDir, dryRun)
}

// syncAnvilSettings syncs the main anvil settings (installed_apps)
func syncAnvilSettings(dryRun bool) error {
	terminal.PrintHeader("Configuration Sync")

	// Ensure Homebrew is installed
	if !brew.IsBrewInstalled() {
		terminal.PrintError("Homebrew is not installed")
		terminal.PrintInfo("💡 Run 'anvil init' to install Homebrew")
		return fmt.Errorf("homebrew required for sync")
	}

	terminal.PrintInfo("📋 Analyzing configuration differences...")
	terminal.PrintInfo("")

	// Get apps from settings.yaml
	installedApps, err := config.GetInstalledApps()
	if err != nil {
		return errors.NewConfigurationError(constants.OpSync, "load-config", err)
	}

	if len(installedApps) == 0 {
		terminal.PrintInfo("No apps found in installed_apps to sync")
		terminal.PrintInfo("💡 Use 'anvil install [app-name]' to add apps to tracking")
		return nil
	}

	terminal.PrintInfo("Apps in settings.yaml: %s", strings.Join(installedApps, ", "))
	terminal.PrintInfo("")

	// Check which ones are missing
	var missingApps []string
	var alreadyInstalled []string

	for _, app := range installedApps {
		if brew.IsPackageInstalled(app) {
			alreadyInstalled = append(alreadyInstalled, app)
		} else {
			missingApps = append(missingApps, app)
		}
	}

	// Report status
	terminal.PrintInfo("installed_apps:")
	for _, app := range alreadyInstalled {
		terminal.PrintInfo("  ✅ %s (already installed)", app)
	}
	for _, app := range missingApps {
		terminal.PrintInfo("  ⬇️  %s (missing - will install)", app)
	}

	terminal.PrintInfo("")

	if len(missingApps) == 0 {
		terminal.PrintSuccess("✅ Configuration is already synced!")
		terminal.PrintInfo("All apps from settings.yaml are installed")
		return nil
	}

	terminal.PrintInfo("🔍 Found %d apps to install, %d already synced", len(missingApps), len(alreadyInstalled))
	terminal.PrintInfo("")

	if dryRun {
		terminal.PrintInfo("Dry run - would install: %s", strings.Join(missingApps, ", "))
		return nil
	}

	// Ask for confirmation
	if !terminal.Confirm("Proceed with sync?") {
		terminal.PrintInfo("Sync cancelled")
		return nil
	}

	terminal.PrintInfo("")

	// Install missing apps
	return installMissingApps(missingApps)
}

// installMissingApps installs the list of missing applications
func installMissingApps(missingApps []string) error {
	terminal.PrintInfo("🔧 Installing missing applications...")
	terminal.PrintInfo("")

	successCount := 0
	var failedApps []string

	for i, app := range missingApps {
		terminal.PrintProgress(i+1, len(missingApps), fmt.Sprintf("Installing %s", app))

		if err := brew.InstallPackage(app); err != nil {
			terminal.PrintError("Failed to install %s: %v", app, err)
			failedApps = append(failedApps, app)
		} else {
			terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", app))
			successCount++
		}
	}

	terminal.PrintInfo("")
	terminal.PrintHeader("Configuration Sync Complete!")

	if successCount > 0 {
		terminal.PrintSuccess(fmt.Sprintf("✅ Successfully installed %d of %d apps", successCount, len(missingApps)))
	}

	if len(failedApps) > 0 {
		terminal.PrintWarning("⚠️  Failed to install: %s", strings.Join(failedApps, ", "))
		terminal.PrintInfo("💡 Check app names or try installing manually with 'brew install [app-name]'")
		return fmt.Errorf("failed to install %d apps", len(failedApps))
	}

	return nil
}

// syncAppConfig handles app-specific config sync (placeholder functionality)
func syncAppConfig(targetDir string, dryRun bool) error {
	// Get config directory
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpSync, "load-config", err)
	}

	// Build path to the pulled config directory
	tempDir := filepath.Join(cfg.Directories.Config, "temp", targetDir)

	// Check if the directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		terminal.PrintError("Configuration directory '%s' not found", targetDir)
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 This could be because:")
		terminal.PrintInfo("   • The app name is incorrect")
		terminal.PrintInfo("   • The configuration was never pulled")
		terminal.PrintInfo("   • The directory name doesn't match what was pulled")
		terminal.PrintInfo("")
		terminal.PrintInfo("🔧 To fix this:")
		terminal.PrintInfo("   • Run 'anvil config pull %s' to download the configuration", targetDir)
		terminal.PrintInfo("   • Check available pulled configs in: %s", filepath.Join(cfg.Directories.Config, "temp"))
		return fmt.Errorf("configuration directory not found")
	}

	// Build and display the tree structure (reuse from show command)
	return showAppSyncStatus(tempDir, targetDir, dryRun)
}

// showAppSyncStatus displays app config sync status (reuses tree logic from show)
func showAppSyncStatus(basePath, targetDir string, dryRun bool) error {
	// Import the tree building logic from show command
	root, err := buildTree(basePath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpSync, "build-tree", err)
	}

	terminal.PrintHeader(fmt.Sprintf("Configuration Sync: %s", targetDir))
	terminal.PrintInfo("Path: %s", basePath)
	terminal.PrintInfo("")

	// If there's only one file at root level, show single file status
	if len(root.Children) == 1 && !root.Children[0].IsDir {
		terminal.PrintInfo("📄 Configuration file: %s", root.Children[0].Name)
	} else {
		// Display the tree structure
		terminal.PrintInfo("📁 Configuration structure:")
		terminal.PrintInfo("")

		// Sort children for consistent display
		sortChildren(root)

		// Print the tree starting from root
		printTreeNode(root, "", true, true)
	}

	terminal.PrintInfo("")
	terminal.PrintWarning("🚧 Syncing %s configurations is in active development", targetDir)
	terminal.PrintInfo("💡 This feature will automatically apply configuration files to their destinations")

	if dryRun {
		terminal.PrintInfo("Dry run - would sync configuration files when feature is ready")
	}

	return nil
}

// TreeNode represents a node in the file tree (reused from show)
type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*TreeNode
}

// buildTree recursively builds a tree structure from the filesystem (reused from show)
func buildTree(dirPath string) (*TreeNode, error) {
	root := &TreeNode{
		Name:     filepath.Base(dirPath),
		Path:     dirPath,
		IsDir:    true,
		Children: []*TreeNode{},
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		// Get relative path from root
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Split the path into components
		parts := strings.Split(relPath, string(filepath.Separator))

		// Find or create the parent node
		current := root
		for i, part := range parts[:len(parts)-1] {
			found := false
			for _, child := range current.Children {
				if child.Name == part && child.IsDir {
					current = child
					found = true
					break
				}
			}
			if !found {
				// Create intermediate directory
				newDir := &TreeNode{
					Name:     part,
					Path:     filepath.Join(dirPath, strings.Join(parts[:i+1], string(filepath.Separator))),
					IsDir:    true,
					Children: []*TreeNode{},
				}
				current.Children = append(current.Children, newDir)
				current = newDir
			}
		}

		// Add the final node
		finalNode := &TreeNode{
			Name:  parts[len(parts)-1],
			Path:  path,
			IsDir: info.IsDir(),
		}
		if info.IsDir() {
			finalNode.Children = []*TreeNode{}
		}
		current.Children = append(current.Children, finalNode)

		return nil
	})

	return root, err
}

// sortChildren recursively sorts all children in the tree (reused from show)
func sortChildren(node *TreeNode) {
	if node.Children == nil {
		return
	}

	// Sort children: directories first, then files, both alphabetically
	for i := 0; i < len(node.Children)-1; i++ {
		for j := i + 1; j < len(node.Children); j++ {
			// Compare: directories first, then alphabetically
			if (!node.Children[i].IsDir && node.Children[j].IsDir) ||
				(node.Children[i].IsDir == node.Children[j].IsDir && node.Children[i].Name > node.Children[j].Name) {
				node.Children[i], node.Children[j] = node.Children[j], node.Children[i]
			}
		}
	}

	// Recursively sort children
	for _, child := range node.Children {
		sortChildren(child)
	}
}

// printTreeNode prints a tree node with ASCII art and colors (reused from show)
func printTreeNode(node *TreeNode, prefix string, isLast bool, isRoot bool) {
	if !isRoot {
		// Choose the appropriate tree character
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}

		// Color the output based on file type
		var coloredName string
		if node.IsDir {
			coloredName = fmt.Sprintf("%s%s%s%s", terminal.ColorBold, terminal.ColorBlue, node.Name, terminal.ColorReset)
		} else {
			// Color files based on extension
			ext := strings.ToLower(filepath.Ext(node.Name))
			switch ext {
			case ".json", ".yaml", ".yml", ".toml":
				coloredName = fmt.Sprintf("%s%s%s", terminal.ColorGreen, node.Name, terminal.ColorReset)
			case ".md", ".txt", ".log":
				coloredName = fmt.Sprintf("%s%s%s", terminal.ColorCyan, node.Name, terminal.ColorReset)
			case ".sh", ".zsh", ".bash":
				coloredName = fmt.Sprintf("%s%s%s", terminal.ColorYellow, node.Name, terminal.ColorReset)
			default:
				coloredName = node.Name
			}
		}

		// Print the current node
		fmt.Printf("%s%s%s\n", prefix, treeChar, coloredName)
	}

	// Print children
	if node.Children != nil {
		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1

			// Calculate prefix for child
			var childPrefix string
			if isRoot {
				childPrefix = ""
			} else {
				if isLast {
					childPrefix = prefix + "    "
				} else {
					childPrefix = prefix + "│   "
				}
			}

			printTreeNode(child, childPrefix, isChildLast, false)
		}
	}
}

func init() {
	// Add flags
	SyncCmd.Flags().Bool("dry-run", false, "Show what would be synced without making changes")
}
