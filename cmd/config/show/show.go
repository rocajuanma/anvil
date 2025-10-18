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
	"strings"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
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

func init() {
	ShowCmd.Flags().Bool("raw", false, "Show raw file content without formatting")
}

// runShowCommand executes the configuration show process
func runShowCommand(cmd *cobra.Command, args []string) error {
	raw, _ := cmd.Flags().GetBool("raw")

	// If no arguments provided, show the anvil settings.yaml
	if len(args) == 0 {
		return showAnvilSettings(raw)
	}

	// Show specific pulled configuration directory
	targetDir := args[0]
	return showPulledConfig(targetDir)
}

// showAnvilSettings displays the main anvil settings.yaml file
func showAnvilSettings(raw bool) error {
	o := getOutputHandler()

	// Stage 1: Locate settings file
	configPath := config.GetConfigPath()

	// Check if settings file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		o.PrintError("Anvil settings file not found at: %s", configPath)
		o.PrintInfo("💡 Run 'anvil init' to create the initial settings file")
		return fmt.Errorf("settings file not found")
	}

	// Stage 2: Read and display content
	content, err := os.ReadFile(configPath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpShow, "read-settings", err)
	}

	// If raw flag is set, show raw content
	if raw {
		o.PrintHeader("Anvil Settings Configuration (Raw)")
		o.PrintInfo("File: %s\n", configPath)
		fmt.Print(string(content))
		return nil
	}

	// Enhanced view with box
	var boxContent strings.Builder
	boxContent.WriteString("\n")
	boxContent.WriteString(fmt.Sprintf("  Location: %s\n", configPath))
	boxContent.WriteString("\n")

	// Parse and display YAML content with slight indentation
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line != "" {
			boxContent.WriteString("  " + line + "\n")
		}
	}

	boxContent.WriteString("\n")

	// Display in box
	fmt.Println(charm.RenderBox("anvil settings.yaml", boxContent.String(), "#00FF87", false))

	// Footer with helpful info
	fmt.Println()
	fmt.Println("  💡 Edit with: nano " + configPath)
	fmt.Println("  💡 Show raw: anvil config show --raw")
	fmt.Println()

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
}
