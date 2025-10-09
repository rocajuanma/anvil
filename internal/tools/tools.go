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

package tools

import (
	"fmt"
	"runtime"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/palantir"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

// Tool represents a macOS system tool
type Tool struct {
	Name        string
	Command     string
	Required    bool
	InstallWith string
	Description string
}

// GetRequiredTools returns the list of required tools for anvil on macOS
// Note: Homebrew is handled separately as a prerequisite in ValidateAndInstallTools()
func GetRequiredTools() []Tool {
	return []Tool{
		{
			Name:        "Git",
			Command:     constants.GitCommand,
			Required:    true,
			InstallWith: "brew",
			Description: "Version control system",
		},
		{
			Name:        "cURL",
			Command:     constants.CurlCommand,
			Required:    true,
			InstallWith: "system",
			Description: "Command line tool for transferring data",
		},
	}
}

// GetOptionalTools returns the list of optional tools for anvil on macOS
func GetOptionalTools() []Tool {
	return []Tool{
		{
			Name:        "Docker",
			Command:     "docker",
			Required:    false,
			InstallWith: "brew",
			Description: "Container runtime",
		},
		{
			Name:        "kubectl",
			Command:     "kubectl",
			Required:    false,
			InstallWith: "brew",
			Description: "Kubernetes command-line tool",
		},
	}
}

// ValidateAndInstallTools validates and installs required tools on macOS
func ValidateAndInstallTools() error {
	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("Anvil only supports macOS")
	}

	// Phase 1: Install Homebrew as a prerequisite (required for other tool installations)
	if err := brew.EnsureBrewIsInstalled(); err != nil {
		return fmt.Errorf("tools: %w", err)
	}

	// Phase 2: Validate and install other required tools (using Homebrew when needed)
	requiredTools := GetRequiredTools()
	for _, tool := range requiredTools {
		if err := validateTool(tool); err != nil {
			return fmt.Errorf("failed to validate required tool %s: %w", tool.Name, err)
		}
	}

	// Validate optional tools (don't fail if they're not available)
	optionalTools := GetOptionalTools()
	for _, tool := range optionalTools {
		if err := validateTool(tool); err != nil {
			getOutputHandler().PrintWarning("Optional tool %s is not available: %v", tool.Name, err)
		}
	}

	return nil
}

// validateTool validates a single tool on macOS
func validateTool(tool Tool) error {
	o := getOutputHandler()
	if system.CommandExists(tool.Command) {
		o.PrintInfo("✓ %s is available", tool.Name)
		return nil
	}

	if !tool.Required {
		o.PrintWarning("○ %s is not installed (optional)", tool.Name)
		return nil
	}

	// Try to install the tool
	o.PrintInfo("Installing %s...", tool.Name)

	switch tool.InstallWith {
	case "brew":
		if err := brew.InstallPackage(tool.Command); err != nil {
			return fmt.Errorf("failed to install %s with brew: %w", tool.Name, err)
		}
	case "system":
		// cURL should be available by default on macOS
		return fmt.Errorf("%s is not available on this macOS system", tool.Name)
	default:
		return fmt.Errorf("unknown installation method for %s", tool.Name)
	}

	// Verify installation
	if !system.CommandExists(tool.Command) {
		return fmt.Errorf("%s was not successfully installed", tool.Name)
	}

	o.PrintSuccess(fmt.Sprintf("%s installed successfully", tool.Name))
	return nil
}

// GetToolInfo returns information about a specific tool
func GetToolInfo(toolName string) (*Tool, error) {
	allTools := append(GetRequiredTools(), GetOptionalTools()...)

	for _, tool := range allTools {
		if tool.Name == toolName || tool.Command == toolName {
			return &tool, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found", toolName)
}

// CheckToolsStatus checks the status of all tools on macOS
func CheckToolsStatus() (map[string]bool, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("tool status check only supported on macOS")
	}

	status := make(map[string]bool)

	allTools := append(GetRequiredTools(), GetOptionalTools()...)
	for _, tool := range allTools {
		status[tool.Name] = system.CommandExists(tool.Command)
	}

	return status, nil
}
