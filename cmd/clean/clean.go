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

package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all content inside .anvil directories",
	Long:  constants.CLEAN_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCleanCommand(cmd, args); err != nil {
			getOutputHandler().PrintError("Clean failed: %v", err)
			return
		}
	},
}

// runCleanCommand executes the clean process
func runCleanCommand(cmd *cobra.Command, args []string) error {
	// Get command flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	force, _ := cmd.Flags().GetBool("force")
	output := getOutputHandler()
	output.PrintHeader("Cleaning Anvil Directories")

	// Get anvil directory path
	anvilDir, err := getAnvilDirectoryPath()
	if err != nil {
		return err
	}

	// Check if .anvil directory exists
	if _, err := os.Stat(anvilDir); os.IsNotExist(err) {
		output.PrintWarning("Directory %s does not exist. Nothing to clean.", anvilDir)
		return nil
	}

	// Get items to clean
	itemsToClean, err := getItemsToClean(anvilDir)
	if err != nil {
		return err
	}

	if len(itemsToClean) == 0 {
		output.PrintSuccess("No root directories found to clean. Only settings.yaml exists.")
		return nil
	}

	// Display what will be cleaned
	displayCleanPreview(output, itemsToClean)

	// Handle user confirmation and dry run
	if !handleUserConfirmation(output, force, dryRun, len(itemsToClean)) {
		return nil
	}

	if dryRun {
		output.PrintInfo("DRY RUN: Would clean contents of %d root directories", len(itemsToClean))
		return nil
	}

	// Perform the actual cleaning
	return performCleaning(output, itemsToClean)
}

// getAnvilDirectoryPath returns the path to the .anvil directory
func getAnvilDirectoryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", &errors.AnvilError{
			Op:      "clean",
			Command: "clean",
			Type:    errors.ErrorTypeFileSystem,
			Err:     fmt.Errorf("failed to get home directory: %w", err),
		}
	}
	return filepath.Join(homeDir, constants.AnvilConfigDir), nil
}

// getItemsToClean scans the anvil directory and returns items to clean
func getItemsToClean(anvilDir string) ([]string, error) {
	output := getOutputHandler()
	output.PrintStage("Scanning .anvil directory for content to clean")

	// Get all items in .anvil directory
	items, err := os.ReadDir(anvilDir)
	if err != nil {
		return nil, &errors.AnvilError{
			Op:      "clean",
			Command: "clean",
			Type:    errors.ErrorTypeFileSystem,
			Err:     fmt.Errorf("failed to read .anvil directory: %w", err),
		}
	}

	var itemsToClean []string
	for _, item := range items {
		// Skip settings.yaml
		if item.Name() == constants.ConfigFileName {
			continue
		}

		itemPath := filepath.Join(anvilDir, item.Name())
		itemsToClean = append(itemsToClean, itemPath)
	}

	return itemsToClean, nil
}

// displayCleanPreview shows what will be cleaned
func displayCleanPreview(output palantir.OutputHandler, itemsToClean []string) {
	output.PrintInfo("Found %d root directories to clean:", len(itemsToClean))
	output.PrintInfo("Directory structure to be cleaned:")

	// Build and display tree structure for each directory
	for _, itemPath := range itemsToClean {
		itemName := filepath.Base(itemPath)
		if info, err := os.Stat(itemPath); err == nil && info.IsDir() {
			// Count items in directory
			count, treeOutput := buildDirectoryTree(itemPath, itemName)
			output.PrintInfo("  📁 %s (%d)", itemName, count)
			fmt.Print(treeOutput)
		} else {
			output.PrintInfo("  📁 %s", itemName)
		}
	}
}

// handleUserConfirmation handles user confirmation and returns true if should proceed
func handleUserConfirmation(output palantir.OutputHandler, force, dryRun bool, itemCount int) bool {
	// Confirm deletion unless force flag is used
	if !force && !dryRun {
		confirmMsg := fmt.Sprintf("Are you sure you want to clean the contents of these %d root directories? This action cannot be undone", itemCount)
		if !output.Confirm(confirmMsg) {
			output.PrintInfo("Clean operation cancelled.")
			return false
		}
	}
	return true
}

// performCleaning executes the actual cleaning process
func performCleaning(output palantir.OutputHandler, itemsToClean []string) error {
	output.PrintStage("Cleaning directories and files")

	// Clean each item
	cleanedCount := 0
	for _, itemPath := range itemsToClean {
		if err := cleanItem(itemPath); err != nil {
			output.PrintWarning("Failed to clean %s: %v", filepath.Base(itemPath), err)
			continue
		}
		cleanedCount++
		displayCleanResult(output, itemPath)
	}

	output.PrintInfo("Successfully cleaned contents of %d/%d root directories", cleanedCount, len(itemsToClean))

	if cleanedCount < len(itemsToClean) {
		output.PrintWarning("Some root directories could not be cleaned. Check the warnings above.")
	}

	return nil
}

// displayCleanResult shows the result of cleaning a specific item
func displayCleanResult(output palantir.OutputHandler, itemPath string) {
	itemName := filepath.Base(itemPath)
	if info, err := os.Stat(itemPath); err == nil && info.IsDir() {
		if itemName == "dotfiles" {
			output.PrintSuccess("Removed dotfiles directory completely")
		} else {
			output.PrintSuccess("Cleaned contents of directory " + itemName)
		}
	} else {
		output.PrintSuccess("Cleaned " + itemName)
	}
}

// buildDirectoryTree builds a simple tree showing only immediate contents of a directory
func buildDirectoryTree(dirPath, dirName string) (int, string) {
	var output strings.Builder
	var count int

	// Read only the immediate contents of the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, ""
	}

	// Count and display only immediate entries
	for _, entry := range entries {
		count++
		output.WriteString(fmt.Sprintf("    ├── %s\n", entry.Name()))
	}

	return count, output.String()
}

// cleanItem removes the contents of a directory or the file itself
func cleanItem(itemPath string) error {
	info, err := os.Stat(itemPath)
	if err != nil {
		return fmt.Errorf("failed to stat item: %w", err)
	}

	if info.IsDir() {
		itemName := filepath.Base(itemPath)

		// Special handling for dotfiles directory - remove it completely
		if itemName == "dotfiles" {
			// Remove the entire dotfiles directory to ensure clean git repository state
			if err := os.RemoveAll(itemPath); err != nil {
				return fmt.Errorf("failed to remove dotfiles directory: %w", err)
			}
			return nil
		}

		// For other directories (temp/, archive/), remove contents but preserve the directory structure
		// This is important for directories that are needed by the tool but can be empty
		entries, err := os.ReadDir(itemPath)
		if err != nil {
			return fmt.Errorf("failed to read directory contents: %w", err)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(itemPath, entry.Name())
			if entry.IsDir() {
				// Remove subdirectory and all its contents
				if err := os.RemoveAll(entryPath); err != nil {
					return fmt.Errorf("failed to remove subdirectory %s: %w", entry.Name(), err)
				}
			} else {
				// Remove file
				if err := os.Remove(entryPath); err != nil {
					return fmt.Errorf("failed to remove file %s: %w", entry.Name(), err)
				}
			}
		}
	} else {
		// Remove single file
		if err := os.Remove(itemPath); err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}
	}

	return nil
}

func init() {
	CleanCmd.Flags().BoolP("dry-run", "n", false, "Show what would be cleaned without actually deleting")
	CleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}
