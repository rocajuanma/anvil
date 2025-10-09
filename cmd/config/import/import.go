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

package importcmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var ImportCmd = &cobra.Command{
	Use:   "import [file-or-url]",
	Short: "Import groups from a local file or URL",
	Long:  "Import tool groups from a local YAML file or remote URL into your anvil configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		importPath := args[0]
		if err := runImportCommand(cmd, importPath); err != nil {
			getOutputHandler().PrintError("Import failed: %v", err)
			return
		}
	},
}

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// ImportConfig represents the structure for importing configurations
type ImportConfig struct {
	Groups map[string][]string `yaml:"groups"`
}

// runImportCommand executes the group import process
func runImportCommand(cmd *cobra.Command, importPath string) error {
	output := getOutputHandler()
	output.PrintHeader("Import Groups from File")

	// Stage 1: Fetch and validate source file
	output.PrintStage("Fetching source file...")
	tempFile, cleanup, err := fetchFile(importPath)
	if err != nil {
		return errors.NewFileSystemError(constants.OpConfig, "fetch-file", err)
	}
	defer cleanup()
	output.PrintSuccess("Source file fetched successfully")

	// Stage 2: Parse and validate import data
	output.PrintStage("Parsing import file...")
	importData, err := parseImportFile(tempFile)
	if err != nil {
		return errors.NewConfigurationError(constants.OpConfig, "parse-import", err)
	}

	if len(importData.Groups) == 0 {
		return errors.NewConfigurationError(constants.OpConfig, "no-groups",
			fmt.Errorf("no valid groups found in import file"))
	}
	output.PrintSuccess("Import file parsed successfully")

	// Stage 3: Validate group structure
	output.PrintStage("Validating group structure...")
	if err := validateImportGroups(importData.Groups); err != nil {
		return errors.NewConfigurationError(constants.OpConfig, "validate-groups", err)
	}
	output.PrintSuccess("Group structure validation passed")

	// Stage 4: Check for conflicts with existing groups
	output.PrintStage("Checking for conflicts...")
	currentConfig, err := config.LoadConfig()
	if err != nil {
		return errors.NewConfigurationError(constants.OpConfig, "load-config", err)
	}

	conflicts := checkGroupConflicts(importData.Groups, currentConfig.Groups)
	if len(conflicts) > 0 {
		return errors.NewConfigurationError(constants.OpConfig, "group-conflicts",
			fmt.Errorf("groups already exist: %s", strings.Join(conflicts, ", ")))
	}
	output.PrintSuccess("No conflicts detected")

	// Stage 5: Display import summary
	output.PrintStage("Preparing import summary...")
	displayImportSummary(importData.Groups)

	// Stage 6: Confirm import
	if !output.Confirm("Proceed with importing these groups?") {
		output.PrintInfo("Import cancelled by user")
		return nil
	}

	// Stage 7: Import groups
	output.PrintStage("Importing groups...")
	if err := importGroups(currentConfig, importData.Groups); err != nil {
		return errors.NewConfigurationError(constants.OpConfig, "import-groups", err)
	}
	output.PrintSuccess("Groups imported successfully")

	output.PrintInfo("\n✨ Import completed! %d groups have been added to your configuration.", len(importData.Groups))
	return nil
}

// fetchFile downloads a file from URL or copies from local path to a temporary file
func fetchFile(sourcePath string) (string, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if isURL(sourcePath) {
		return fetchFromURL(ctx, sourcePath)
	}

	// Handle local file
	return fetchFromLocal(sourcePath)
}

// isURL checks if the given string is a URL
func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// fetchFromURL downloads file from URL to a temporary file
func fetchFromURL(ctx context.Context, fileURL string) (string, func(), error) {
	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent
	req.Header.Set("User-Agent", "anvil-cli/1.0")

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "anvil-import-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Copy content to temporary file
	_, err = io.Copy(tempFile, resp.Body)
	tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}

// fetchFromLocal copies local file to temporary file for consistent handling
func fetchFromLocal(filePath string) (string, func(), error) {
	// Validate file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return "", nil, fmt.Errorf("file does not exist: %s", filePath)
		}
		return "", nil, fmt.Errorf("cannot access file: %w", err)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "anvil-import-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFile.Close()

	// Copy file content
	sourceData, err := os.ReadFile(filePath)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(tempFile.Name(), sourceData, constants.FilePerm); err != nil {
		os.Remove(tempFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	cleanup := func() {
		os.Remove(tempFile.Name())
	}

	return tempFile.Name(), cleanup, nil
}

// parseImportFile parses the import file and extracts only group data
func parseImportFile(filePath string) (*ImportConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	// Parse as generic map first to extract only groups
	var rawData map[string]interface{}
	if err := yaml.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract only groups section
	groupsData, exists := rawData["groups"]
	if !exists {
		return &ImportConfig{Groups: make(map[string][]string)}, nil
	}

	// Convert to proper structure
	groupsMap, ok := groupsData.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("groups section has invalid format")
	}

	importConfig := &ImportConfig{
		Groups: make(map[string][]string),
	}

	for groupName, groupTools := range groupsMap {
		groupNameStr, ok := groupName.(string)
		if !ok {
			continue // Skip invalid group names
		}

		toolsList, ok := groupTools.([]interface{})
		if !ok {
			continue // Skip invalid tool lists
		}

		var tools []string
		for _, tool := range toolsList {
			if toolStr, ok := tool.(string); ok {
				tools = append(tools, toolStr)
			}
		}

		if len(tools) > 0 {
			importConfig.Groups[groupNameStr] = tools
		}
	}

	return importConfig, nil
}

// validateImportGroups validates the structure of imported groups
func validateImportGroups(groups map[string][]string) error {
	if len(groups) == 0 {
		return fmt.Errorf("no groups found to import")
	}

	validator := config.NewConfigValidator(nil)

	for groupName, tools := range groups {
		// Validate group name
		if err := validator.ValidateGroupName(groupName); err != nil {
			return fmt.Errorf("invalid group name '%s': %w", groupName, err)
		}

		// Validate group is not empty
		if len(tools) == 0 {
			return fmt.Errorf("group '%s' cannot be empty", groupName)
		}

		// Validate each tool name
		for _, tool := range tools {
			if err := validator.ValidateAppName(tool); err != nil {
				return fmt.Errorf("invalid tool '%s' in group '%s': %w", tool, groupName, err)
			}
		}
	}

	return nil
}

// checkGroupConflicts checks if any imported groups already exist
func checkGroupConflicts(importGroups map[string][]string, existingGroups config.AnvilGroups) []string {
	var conflicts []string
	for groupName := range importGroups {
		if _, exists := existingGroups[groupName]; exists {
			conflicts = append(conflicts, groupName)
		}
	}
	sort.Strings(conflicts)
	return conflicts
}

// displayImportSummary shows a tree view of groups that will be imported
func displayImportSummary(groups map[string][]string) {
	output := getOutputHandler()
	output.PrintInfo("\n📋 Import Summary:")
	output.PrintInfo("═══════════════════")

	// Sort group names for consistent output
	var groupNames []string
	for groupName := range groups {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)

	totalGroups := len(groups)
	totalApps := 0

	for _, groupName := range groupNames {
		tools := groups[groupName]
		totalApps += len(tools)

		// Display group with tree structure
		output.PrintInfo("├── 📁 %s (%d tools)", groupName, len(tools))

		// Sort tools for consistent output
		sortedTools := make([]string, len(tools))
		copy(sortedTools, tools)
		sort.Strings(sortedTools)

		for i, tool := range sortedTools {
			if i == len(sortedTools)-1 {
				output.PrintInfo("│   └── 🔧 %s", tool)
			} else {
				output.PrintInfo("│   ├── 🔧 %s", tool)
			}
		}
		output.PrintInfo("│")
	}

	output.PrintInfo("📊 Total: %d groups, %d applications", totalGroups, totalApps)
	output.PrintInfo("")
}

// importGroups adds the imported groups to the current configuration
func importGroups(currentConfig *config.AnvilConfig, importGroups map[string][]string) error {
	// Add new groups to existing configuration
	for groupName, tools := range importGroups {
		currentConfig.Groups[groupName] = tools
	}

	// Save updated configuration
	return config.SaveConfig(currentConfig)
}
