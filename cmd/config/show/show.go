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

package show

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/interfaces"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() interfaces.OutputHandler {
	return terminal.GetGlobalOutputHandler()
}

var ShowCmd = &cobra.Command{
	Use:   "show [directory]",
	Short: "Show configuration files from anvil settings or pulled directories",
	Long:  constants.SHOW_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1), // Accept 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		if err := runShowCommand(cmd, args); err != nil {
			getOutputHandler().PrintError("Show failed: %v", err)
			return
		}
	},
}

// runShowCommand executes the configuration show process
func runShowCommand(cmd *cobra.Command, args []string) error {
	// If no arguments provided, show the anvil settings.yaml
	if len(args) == 0 {
		return showAnvilSettings()
	}

	// Show specific pulled configuration directory
	targetDir := args[0]
	return showPulledConfig(targetDir)
}

// showAnvilSettings displays the main anvil settings.yaml file
func showAnvilSettings() error {
	o := getOutputHandler()
	o.PrintHeader("Anvil Settings Configuration")

	// Stage 1: Locate settings file
	o.PrintStage("Locating anvil settings file...")
	configPath := config.GetConfigPath()

	// Check if settings file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		o.PrintError("Anvil settings file not found at: %s", configPath)
		o.PrintInfo("💡 Run 'anvil init' to create the initial settings file")
		return fmt.Errorf("settings file not found")
	}
	o.PrintSuccess("Settings file located")

	o.PrintInfo("File: %s\n", configPath)

	// Stage 2: Read and display content
	o.PrintStage("Reading configuration content...")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpShow, "read-settings", err)
	}
	o.PrintSuccess("Configuration content loaded\n")

	fmt.Print(string(content))
	return nil
}

// showPulledConfig displays configuration files from a pulled directory
func showPulledConfig(targetDir string) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Configuration Directory: %s", targetDir))

	// Stage 1: Load anvil configuration
	o.PrintStage("Loading anvil configuration...")
	o.PrintSuccess("Configuration loaded")

	// Stage 2: Locate pulled configuration directory
	o.PrintStage("Locating pulled configuration directory...")
	tempDir := filepath.Join(config.GetConfigDirectory(), "temp", targetDir)

	// Check if the directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		o.PrintError("Configuration directory '%s' not found\n", targetDir)
		o.PrintInfo("💡 This could be because:")
		o.PrintInfo("   • The app name is incorrect")
		o.PrintInfo("   • The configuration was never pulled")
		o.PrintInfo("   • Use 'anvil config pull %s' to pull this configuration first", targetDir)
		o.PrintInfo("")

		// Show available pulled configurations
		tempBasePath := filepath.Join(config.GetConfigDirectory(), "temp")
		if entries, err := os.ReadDir(tempBasePath); err == nil && len(entries) > 0 {
			o.PrintInfo("Available pulled configurations:")
			for _, entry := range entries {
				if entry.IsDir() {
					o.PrintInfo("  • %s", entry.Name())
				}
			}
		} else {
			o.PrintInfo("No configurations have been pulled yet.")
			o.PrintInfo("Use 'anvil config pull <directory>' to pull configurations from your repository.")
		}

		return fmt.Errorf("configuration directory not found")
	}
	o.PrintSuccess("Configuration directory located")
	o.PrintInfo("Directory: %s\n", tempDir)

	// Stage 3: Display directory contents
	o.PrintStage("Reading configuration files...")
	err := showDirectoryTree(tempDir, targetDir)
	if err != nil {
		return err
	}
	o.PrintSuccess("Configuration files displayed")

	return nil
}

// TreeNode represents a node in the file tree
type TreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*TreeNode
}

// showDirectoryTree displays a tree structure of files/directories
func showDirectoryTree(basePath, targetDir string) error {
	// Build the tree structure
	root, err := buildTree(basePath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpShow, "build-tree", err)
	}

	// If there's only one file at root level, display its content directly
	if len(root.Children) == 1 && !root.Children[0].IsDir {
		return showSingleFile(root.Children[0].Path, targetDir)
	}
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Configuration Directory: %s", targetDir))
	o.PrintInfo("Path: %s\n", basePath)

	// Display the tree structure
	o.PrintInfo("Directory structure:\n")

	// Sort children for consistent display
	sortChildren(root)

	// Print the tree starting from root
	printTreeNode(root, "", true, true)

	o.PrintInfo("\n💡 To view a specific file, you can use:")
	o.PrintInfo("   • cat %s/[filename]", basePath)
	o.PrintInfo("   • Or navigate to the directory and explore manually")

	return nil
}

// buildTree recursively builds a tree structure from the filesystem
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

// sortChildren recursively sorts all children in the tree (directories first, then files, both alphabetically)
func sortChildren(node *TreeNode) {
	if node.Children == nil {
		return
	}

	// Sort children: directories first, then files, both alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir // directories come first
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortChildren(child)
	}
}

// printTreeNode prints a tree node with ASCII art and colors
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

// showSingleFile displays the content of a single configuration file
func showSingleFile(filePath, targetDir string) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Configuration: %s", targetDir))
	o.PrintInfo("File: %s\n", filepath.Base(filePath))

	// Read and display the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpShow, "read-config-file", err)
	}

	fmt.Print(string(content))
	return nil
}

func init() {
	// Add any flags if needed in the future
}
