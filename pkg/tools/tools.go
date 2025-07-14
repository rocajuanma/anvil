package tools

import (
	"fmt"
	"runtime"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
)

// Tool represents a system tool
type Tool struct {
	Name        string
	Command     string
	Required    bool
	InstallWith string
	Description string
}

// GetRequiredTools returns the list of required tools for anvil
func GetRequiredTools() []Tool {
	return []Tool{
		{
			Name:        "Git",
			Command:     "git",
			Required:    true,
			InstallWith: "brew",
			Description: "Version control system",
		},
		{
			Name:        "cURL",
			Command:     "curl",
			Required:    true,
			InstallWith: "system",
			Description: "Command line tool for transferring data",
		},
	}
}

// GetOptionalTools returns the list of optional tools for anvil
func GetOptionalTools() []Tool {
	return []Tool{
		{
			Name:        "Homebrew",
			Command:     "brew",
			Required:    false,
			InstallWith: "script",
			Description: "Package manager for macOS",
		},
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

// ValidateAndInstallTools validates and installs required tools
func ValidateAndInstallTools() error {
	// First, ensure Homebrew is installed on macOS
	if runtime.GOOS == "darwin" {
		if !brew.IsBrewInstalled() {
			terminal.PrintInfo("Homebrew not found. Installing Homebrew...")
			if err := brew.InstallBrew(); err != nil {
				return fmt.Errorf("failed to install Homebrew: %w", err)
			}
			terminal.PrintSuccess("Homebrew installed successfully")
		}
	}

	// Validate required tools
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
			terminal.PrintWarning("Optional tool %s is not available: %v", tool.Name, err)
		}
	}

	return nil
}

// validateTool validates a single tool
func validateTool(tool Tool) error {
	if system.CommandExists(tool.Command) {
		terminal.PrintInfo("✓ %s is available", tool.Name)
		return nil
	}

	if !tool.Required {
		terminal.PrintWarning("○ %s is not installed (optional)", tool.Name)
		return nil
	}

	// Try to install the tool
	terminal.PrintInfo("Installing %s...", tool.Name)

	switch tool.InstallWith {
	case "brew":
		if runtime.GOOS == "darwin" {
			if err := brew.InstallPackage(tool.Command); err != nil {
				return fmt.Errorf("failed to install %s with brew: %w", tool.Name, err)
			}
		} else {
			return fmt.Errorf("%s requires manual installation on %s", tool.Name, runtime.GOOS)
		}
	case "script":
		if tool.Command == "brew" {
			if err := brew.InstallBrew(); err != nil {
				return fmt.Errorf("failed to install %s: %w", tool.Name, err)
			}
		} else {
			return fmt.Errorf("unsupported script installation for %s", tool.Name)
		}
	case "system":
		// System tools should be available by default
		return fmt.Errorf("%s is not available on this system", tool.Name)
	default:
		return fmt.Errorf("unknown installation method for %s", tool.Name)
	}

	// Verify installation
	if !system.CommandExists(tool.Command) {
		return fmt.Errorf("%s was not successfully installed", tool.Name)
	}

	terminal.PrintSuccess(fmt.Sprintf("%s installed successfully", tool.Name))
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

// ListTools lists all available tools
func ListTools() {
	terminal.PrintHeader("Required Tools")
	for _, tool := range GetRequiredTools() {
		status := "❌ Not installed"
		if system.CommandExists(tool.Command) {
			status = "✅ Installed"
		}
		terminal.PrintInfo("- %s (%s): %s - %s", tool.Name, tool.Command, tool.Description, status)
	}

	terminal.PrintHeader("Optional Tools")
	for _, tool := range GetOptionalTools() {
		status := "❌ Not installed"
		if system.CommandExists(tool.Command) {
			status = "✅ Installed"
		}
		terminal.PrintInfo("- %s (%s): %s - %s", tool.Name, tool.Command, tool.Description, status)
	}
}

// CheckToolsStatus checks the status of all tools
func CheckToolsStatus() (map[string]bool, error) {
	status := make(map[string]bool)

	allTools := append(GetRequiredTools(), GetOptionalTools()...)
	for _, tool := range allTools {
		status[tool.Name] = system.CommandExists(tool.Command)
	}

	return status, nil
}
