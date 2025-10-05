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

package update

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Anvil to the latest version",
	Long:  constants.UPDATE_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUpdateCommand(cmd); err != nil {
			getOutputHandler().PrintError("Update failed: %v", err)
			return
		}
	},
}

// runUpdateCommand executes the update process
func runUpdateCommand(cmd *cobra.Command) error {
	o := getOutputHandler()
	// Ensure we're running on macOS (following existing project pattern)
	if runtime.GOOS != "darwin" {
		return errors.NewPlatformError(constants.OpUpdate, "anvil",
			fmt.Errorf("update command is only supported on macOS"))
	}

	o.PrintHeader("Updating Anvil to Latest Version")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	result, err := updateAnvil(cmd.Context(), dryRun)

	if err != nil {
		return errors.NewInstallationError(constants.OpUpdate, "anvil",
			fmt.Errorf("failed to execute update script: %w", err))
	}

	// For dry-run mode, result will be nil. Return early
	if dryRun {
		return nil
	}

	if !result.Success {
		return errors.NewInstallationError(constants.OpUpdate, "anvil",
			fmt.Errorf("update script failed with exit code %d: %s", result.ExitCode, result.Output))
	}

	o.PrintSuccess("Anvil has been successfully updated!")
	o.PrintInfo("Run 'anvil --version' to verify the new version")
	o.PrintInfo("You may need to restart your terminal session for changes to take effect")

	return nil
}

// updateAnvil updates Anvil to the latest version
// it uses the curl command to download the latest installation script from GitHub releases
func updateAnvil(ctx context.Context, dryRun bool) (*system.CommandResult, error) {
	o := getOutputHandler()

	if dryRun {
		o.PrintInfo("Dry run mode - would update Anvil to the latest version")
		o.PrintInfo("Command that would be executed:")
		o.PrintInfo("curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash")
		return nil, nil
	}

	// Check if curl is available
	if !system.CommandExists("curl") {
		return nil, errors.NewAnvilErrorWithType(constants.OpUpdate, "curl", errors.ErrorTypeInstallation,
			fmt.Errorf("curl is required for updating Anvil but is not available"))
	}

	o.PrintStage("Downloading and executing update script...")
	o.PrintInfo("Fetching latest version from GitHub releases...")

	// Execute the update command using the existing system package
	result, err := system.RunCommandWithTimeout(
		ctx,
		"bash",
		"-c",
		"curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash",
	)

	return result, err
}
func init() {
	UpdateCmd.Flags().Bool("dry-run", false, "Show what would be updated without actually updating")
}
