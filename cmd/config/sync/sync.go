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
	"time"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/terminal"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [app-name]",
	Short: "Sync pulled configuration files to their local destinations",
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

	// Sync specific app config
	appName := args[0]
	return syncAppConfig(appName, dryRun)
}

// syncAnvilSettings syncs the main anvil settings.yaml file
func syncAnvilSettings(dryRun bool) error {
	terminal.PrintHeader("Configuration Sync: anvil settings")

	// Load current config to get temp directory
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpSync, "load-config", err)
	}

	// Check if pulled anvil settings exist
	tempSettingsPath := filepath.Join(cfg.Directories.Config, "temp", "anvil", "settings.yaml")
	if _, err := os.Stat(tempSettingsPath); os.IsNotExist(err) {
		terminal.PrintError("Pulled anvil settings not found")
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 No pulled settings found at: %s", tempSettingsPath)
		terminal.PrintInfo("🔧 To fix this:")
		terminal.PrintInfo("   • Run 'anvil config pull anvil' to download settings")
		terminal.PrintInfo("   • Ensure your repository has an 'anvil' directory with settings.yaml")
		return fmt.Errorf("config not pulled yet")
	}

	// Get current settings path
	currentSettingsPath := config.GetConfigPath()

	// Display sync information
	terminal.PrintInfo("Source: %s", tempSettingsPath)
	terminal.PrintInfo("Destination: %s", currentSettingsPath)
	terminal.PrintInfo("")

	if dryRun {
		terminal.PrintInfo("Dry run - would sync anvil settings")
		return nil
	}

	// Create archive before syncing
	archivePath, err := createArchiveDirectory("anvil-settings")
	if err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	terminal.PrintInfo("Archive: %s", archivePath)
	terminal.PrintInfo("")

	// Ask for confirmation
	if !terminal.Confirm("Override local settings.yaml? Old copy will be archived.") {
		terminal.PrintInfo("Sync cancelled")
		return nil
	}

	terminal.PrintInfo("")

	// Archive existing settings
	if err := archiveExistingConfig("anvil-settings", currentSettingsPath, archivePath); err != nil {
		return fmt.Errorf("failed to archive existing settings: %w", err)
	}

	// Copy new settings
	if err := copyFile(tempSettingsPath, currentSettingsPath); err != nil {
		return fmt.Errorf("failed to copy new settings: %w", err)
	}

	// Report success
	terminal.PrintSuccess("✅ Settings synced successfully")
	terminal.PrintInfo("📦 Old settings archived to: %s", archivePath)
	terminal.PrintInfo("💡 Manual recovery possible from archive directory (no auto-recover yet)")

	return nil
}

// syncAppConfig syncs configuration files for a specific app
func syncAppConfig(appName string, dryRun bool) error {
	terminal.PrintHeader(fmt.Sprintf("Configuration Sync: %s", appName))

	// Load current config
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpSync, "load-config", err)
	}

	// Check if pulled app config exists
	tempAppPath := filepath.Join(cfg.Directories.Config, "temp", appName)
	if _, err := os.Stat(tempAppPath); os.IsNotExist(err) {
		terminal.PrintError("Pulled %s configuration not found", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 No pulled config found at: %s", tempAppPath)
		terminal.PrintInfo("🔧 To fix this:")
		terminal.PrintInfo("   • Run 'anvil config pull %s' to download configuration", appName)
		terminal.PrintInfo("   • Ensure your repository has a '%s' directory", appName)
		return fmt.Errorf("config not pulled yet")
	}

	// Check if app config path is defined in settings
	if cfg.Configs == nil {
		return fmt.Errorf("no configs section found in settings.yaml")
	}

	localConfigPath, exists := cfg.Configs[appName]
	if !exists {
		terminal.PrintError("App config path not configured")
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 The app '%s' doesn't have a local config path defined", appName)
		terminal.PrintInfo("🔧 To fix this:")
		terminal.PrintInfo("   • Edit your settings.yaml file")
		terminal.PrintInfo("   • Add the following to the 'configs' section:")
		terminal.PrintInfo("")
		terminal.PrintInfo("configs:")
		terminal.PrintInfo("  %s: \"/path/to/%s/config\"", appName, appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("Example paths:")
		terminal.PrintInfo("  • ~/.config/%s", appName)
		terminal.PrintInfo("  • ~/Library/Application Support/%s", strings.Title(appName))
		return fmt.Errorf("app config path not defined")
	}

	// Display sync information
	terminal.PrintInfo("Source: %s", tempAppPath)
	terminal.PrintInfo("Destination: %s", localConfigPath)
	terminal.PrintInfo("")

	if dryRun {
		terminal.PrintInfo("Dry run - would sync %s configuration", appName)
		return nil
	}

	// Create archive before syncing
	archivePath, err := createArchiveDirectory(fmt.Sprintf("%s-configs", appName))
	if err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	terminal.PrintInfo("Archive: %s", archivePath)
	terminal.PrintInfo("")

	// Ask for confirmation
	if !terminal.Confirm(fmt.Sprintf("Override %s configs? Old copy will be archived.", appName)) {
		terminal.PrintInfo("Sync cancelled")
		return nil
	}

	terminal.PrintInfo("")

	// Archive existing config
	if err := archiveExistingConfig(fmt.Sprintf("%s-configs", appName), localConfigPath, archivePath); err != nil {
		return fmt.Errorf("failed to archive existing config: %w", err)
	}

	// Copy new config
	if err := copyDirRecursive(tempAppPath, localConfigPath); err != nil {
		return fmt.Errorf("failed to copy new config: %w", err)
	}

	// Report success
	terminal.PrintSuccess(fmt.Sprintf("✅ %s configs synced successfully", strings.Title(appName)))
	terminal.PrintInfo("📦 Old configs archived to: %s", archivePath)
	terminal.PrintInfo("💡 Manual recovery possible from archive directory (no auto-recover yet)")

	return nil
}

// createArchiveDirectory creates a timestamped archive directory
func createArchiveDirectory(prefix string) (string, error) {
	// Load config to get base directory
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	// Create timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	archiveName := fmt.Sprintf("%s-%s", prefix, timestamp)
	archivePath := filepath.Join(cfg.Directories.Config, "archive", archiveName)

	// Create archive directory
	if err := os.MkdirAll(archivePath, constants.DirPerm); err != nil {
		return "", err
	}

	return archivePath, nil
}

// archiveExistingConfig archives the existing configuration
func archiveExistingConfig(configType, sourcePath, archivePath string) error {
	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		// Nothing to archive
		return nil
	}

	// Determine destination in archive
	var destPath string
	if configType == "anvil-settings" {
		destPath = filepath.Join(archivePath, "settings.yaml")
	} else {
		// For app configs, preserve the directory structure
		destPath = archivePath
	}

	// Copy to archive
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return copyDirRecursive(sourcePath, destPath)
	} else {
		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), constants.DirPerm); err != nil {
			return err
		}
		return copyFile(sourcePath, destPath)
	}
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

// copyDirRecursive recursively copies a directory
func copyDirRecursive(src, dst string) error {
	// Remove destination if it exists
	if err := os.RemoveAll(dst); err != nil {
		return err
	}

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

func init() {
	// Add flags
	SyncCmd.Flags().Bool("dry-run", false, "Show what would be synced without making changes")
}
