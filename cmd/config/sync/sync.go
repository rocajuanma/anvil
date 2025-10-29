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

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/utils"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [app-name]",
	Short: "Sync pulled configuration files to their local destinations",
	Long:  constants.SYNC_COMMAND_LONG_DESCRIPTION,
	Args:  cobra.MaximumNArgs(1), // Accept 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSyncCommand(cmd, args); err != nil {
			palantir.GetGlobalOutputHandler().PrintError("Sync failed: %v", err)
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
	o := palantir.GetGlobalOutputHandler()
	o.PrintHeader("Configuration Sync: Anvil settings")

	// Check if pulled anvil settings exist
	tempSettingsPath := filepath.Join(config.GetAnvilConfigDirectory(), "temp", constants.ANVIL, constants.ANVIL_CONFIG_FILE)
	if _, err := os.Stat(tempSettingsPath); os.IsNotExist(err) {
		o.PrintError("Pulled anvil settings not found\n")
		o.PrintInfo("💡 No pulled settings found at: %s", tempSettingsPath)
		o.PrintInfo("🔧 To fix this:")
		o.PrintInfo("   • Run 'anvil config pull anvil' to download settings")
		o.PrintInfo("   • Ensure your repository has an 'anvil' directory with settings.yaml")
		return fmt.Errorf("config not pulled yet")
	}

	// Get current settings path
	currentSettingsPath := config.GetAnvilConfigPath()

	// Display sync information
	o.PrintInfo("Source: %s", tempSettingsPath)
	o.PrintInfo("Destination: %s\n", currentSettingsPath)

	if dryRun {
		o.PrintInfo("Dry run - would sync anvil settings")
		return nil
	}

	// Create archive before syncing
	archivePath, err := createArchiveDirectory("anvil-settings")
	if err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	o.PrintInfo("Archive: %s\n", archivePath)

	// Ask for confirmation
	if !o.Confirm(fmt.Sprintf("Override local %s? Old copy will be archived.", constants.ANVIL_CONFIG_FILE)) {
		o.PrintInfo("Sync cancelled")
		return nil
	}

	o.PrintInfo("")

	spinner := charm.NewDotsSpinner("Syncing anvil settings")
	spinner.Start()

	// Archive existing settings
	if err := archiveExistingConfig("anvil-settings", currentSettingsPath, archivePath); err != nil {
		spinner.Error("Failed to archive existing settings")
		return fmt.Errorf("failed to archive existing settings: %w", err)
	}

	// Copy new settings
	if err := utils.CopyFileSimple(tempSettingsPath, currentSettingsPath); err != nil {
		spinner.Error("Failed to copy new settings")
		return fmt.Errorf("failed to copy new settings: %w", err)
	}

	spinner.Success("Settings synced successfully")

	// Report success
	o.PrintSuccess("✅ Settings synced successfully")
	o.PrintInfo("📦 Old settings archived to: %s", archivePath)
	o.PrintInfo("💡 Manual recovery possible from archive directory (no auto-recover yet)")

	return nil
}

// syncAppConfig syncs configuration files for a specific app
func syncAppConfig(appName string, dryRun bool) error {
	output := palantir.GetGlobalOutputHandler()
	output.PrintHeader(fmt.Sprintf("Configuration Sync: %s", appName))

	// Load current config
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpSync, "load-config", err)
	}

	// Check if pulled app config exists
	tempAppPath := filepath.Join(config.GetAnvilConfigDirectory(), "temp", appName)
	if _, err := os.Stat(tempAppPath); os.IsNotExist(err) {
		output.PrintError("Pulled %s configuration not found\n", appName)
		output.PrintInfo("💡 No pulled config found at: %s", tempAppPath)
		output.PrintInfo("🔧 To fix this:")
		output.PrintInfo("   • Run 'anvil config pull %s' to download configuration", appName)
		output.PrintInfo("   • Ensure your repository has a '%s' directory", appName)
		return fmt.Errorf("config not pulled yet")
	}

	// Check if app config path is defined in settings
	if cfg.Configs == nil {
		return fmt.Errorf("no configs section found in %s", constants.ANVIL_CONFIG_FILE)
	}

	localConfigPath, exists := cfg.Configs[appName]
	if !exists {
		output.PrintError("App config path not configured\n")
		output.PrintInfo("💡 The app '%s' doesn't have a local config path defined", appName)
		output.PrintInfo("🔧 To fix this:")
		output.PrintInfo("   • Edit your %s file", constants.ANVIL_CONFIG_FILE)
		output.PrintInfo("   • Add the following to the 'configs' section:\n")
		output.PrintInfo("configs:")
		output.PrintInfo("  %s: \"/path/to/%s/config\"\n", appName, appName)
		output.PrintInfo("Example paths:")
		output.PrintInfo("  • ~/.config/%s", appName)
		output.PrintInfo("  • ~/Library/Application Support/%s", strings.Title(appName))
		return fmt.Errorf("app config path not defined")
	}

	// Display sync information
	output.PrintInfo("Source: %s", tempAppPath)
	output.PrintInfo("Destination: %s\n", localConfigPath)

	if dryRun {
		output.PrintInfo("Dry run - would sync %s configuration", appName)
		return nil
	}

	// Create archive before syncing
	archivePath, err := createArchiveDirectory(fmt.Sprintf("%s-configs", appName))
	if err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	output.PrintInfo("Archive: %s\n", archivePath)

	// Ask for confirmation
	if !output.Confirm(fmt.Sprintf("Override %s configs? Old copy will be archived.", appName)) {
		output.PrintInfo("Sync cancelled")
		return nil
	}

	output.PrintInfo("")

	spinner := charm.NewDotsSpinner(fmt.Sprintf("Syncing %s configuration", appName))
	spinner.Start()

	// Archive existing config
	if err := archiveExistingConfig(fmt.Sprintf("%s-configs", appName), localConfigPath, archivePath); err != nil {
		spinner.Error("Failed to archive existing config")
		return fmt.Errorf("failed to archive existing config: %w", err)
	}

	// Copy new config
	if err := copyDirRecursive(tempAppPath, localConfigPath); err != nil {
		spinner.Error("Failed to copy new config")
		return fmt.Errorf("failed to copy new config: %w", err)
	}

	spinner.Success(fmt.Sprintf("%s configuration synced successfully", strings.Title(appName)))

	// Report success
	output.PrintSuccess(fmt.Sprintf("✅ %s configs synced successfully", strings.Title(appName)))
	output.PrintInfo("📦 Old configs archived to: %s", archivePath)
	output.PrintInfo("💡 Manual recovery possible from archive directory (no auto-recover yet)")

	return nil
}

// copyDirRecursive recursively copies a directory using the consolidated utils.CopyDirectorySimple
func copyDirRecursive(src, dst string) error {
	return utils.CopyDirectorySimple(src, dst)
}

func init() {
	// Add flags
	SyncCmd.Flags().Bool("dry-run", false, "Show what would be synced without making changes")
}
