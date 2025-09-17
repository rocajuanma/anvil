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

package brew

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/interfaces"
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() interfaces.OutputHandler {
	return terminal.GetGlobalOutputHandler()
}

// BrewPackage represents a brew package
type BrewPackage struct {
	Name        string
	Version     string
	Description string
	Installed   bool
}

// EnsureBrewIsInstalled ensures Homebrew is installed
func EnsureBrewIsInstalled() error {
	if !IsBrewInstalled() {
		getOutputHandler().PrintInfo("Homebrew not found. Installing Homebrew...")
		if err := InstallBrew(); err != nil {
			return fmt.Errorf("failed to install Homebrew: %w", err)
		}
		getOutputHandler().PrintSuccess("Homebrew installed successfully")
	}

	return nil
}

// IsBrewInstalled checks if Homebrew is installed
func IsBrewInstalled() bool {
	return system.CommandExists(constants.BrewCommand)
}

// IsBrewInstalledAtPath checks if Homebrew is installed at known paths
func IsBrewInstalledAtPath() bool {
	brewPaths := []string{
		"/opt/homebrew/bin/brew", // Apple Silicon
		"/usr/local/bin/brew",    // Intel
	}

	for _, path := range brewPaths {
		result, err := system.RunCommand("test", "-x", path)
		if err == nil && result.Success {
			return true
		}
	}

	return system.CommandExists("brew")
}

// InstallBrew installs Homebrew if not already installed
func InstallBrew() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("Homebrew is only supported on macOS")
	}

	if IsBrewInstalled() {
		return nil
	}

	xcodeResult, err := system.RunCommand("xcode-select", "-p")
	if err != nil || !xcodeResult.Success {
		return fmt.Errorf("Xcode Command Line Tools required for Homebrew installation. Install with: xcode-select --install")
	}

	getOutputHandler().PrintInfo("Installing Homebrew...")

	installCmd := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

	result, err := system.RunCommand("/bin/bash", "-c", installCmd)
	if err != nil {
		return fmt.Errorf("failed to run brew installation command: %w", err)
	}

	if !result.Success {
		errorDetails := result.Error
		if result.Output != "" {
			errorDetails = fmt.Sprintf("%s\nOutput: %s", result.Error, strings.TrimSpace(result.Output))
		}
		return fmt.Errorf("brew installation failed: %s", errorDetails)
	}

	if !IsBrewInstalledAtPath() {
		return fmt.Errorf("Homebrew installation completed but brew command not accessible")
	}

	return nil
}

// UpdateBrew updates Homebrew and its formulae
func UpdateBrew() error {
	if !IsBrewInstalled() {
		return fmt.Errorf("Homebrew is not installed")
	}

	getOutputHandler().PrintInfo("Updating Homebrew...")

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewUpdate)
	if err != nil {
		return fmt.Errorf("failed to run brew update: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("brew update failed: %s", result.Error)
	}

	return nil
}

// InstallPackage installs a package using Homebrew
func InstallPackage(packageName string) error {
	if !IsBrewInstalled() {
		return fmt.Errorf("Homebrew is not installed")
	}

	getOutputHandler().PrintInfo("Installing %s...", packageName)

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewInstall, packageName)
	if err != nil {
		return fmt.Errorf("failed to run brew install: %w", err)
	}

	if !result.Success {
		// Include actual brew output for better diagnostics
		var errorDetails string
		if result.Output != "" {
			errorDetails = fmt.Sprintf("brew output: %s", strings.TrimSpace(result.Output))
		} else {
			errorDetails = fmt.Sprintf("system error: %s", result.Error)
		}
		return fmt.Errorf("failed to install %s: %s", packageName, errorDetails)
	}

	return nil
}

// IsPackageInstalled checks if a package is installed (both formulas and casks)
func IsPackageInstalled(packageName string) bool {
	if !IsBrewInstalled() {
		return false
	}

	// First try to check if it's installed as a formula
	result, err := system.RunCommand(constants.BrewCommand, constants.BrewList, "--formula", packageName)
	if err == nil && result.Success {
		return true
	}

	// If not found as formula, check if it's installed as a cask
	result, err = system.RunCommand(constants.BrewCommand, constants.BrewList, "--cask", packageName)
	if err == nil && result.Success {
		return true
	}

	return false
}

// GetInstalledPackages returns a list of installed packages
func GetInstalledPackages() ([]BrewPackage, error) {
	if !IsBrewInstalled() {
		return nil, fmt.Errorf("Homebrew is not installed")
	}

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewList, "--formula")
	if err != nil {
		return nil, fmt.Errorf("failed to run brew list: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get installed packages: %s", result.Error)
	}

	var packages []BrewPackage
	lines := strings.Split(strings.TrimSpace(result.Output), "\n")

	for _, line := range lines {
		if line != "" {
			packages = append(packages, BrewPackage{
				Name:      line,
				Installed: true,
			})
		}
	}

	return packages, nil
}

// InstallPackages installs multiple packages
func InstallPackages(packages []string) error {
	if !IsBrewInstalled() {
		return fmt.Errorf("Homebrew is not installed")
	}

	for i, pkg := range packages {
		getOutputHandler().PrintProgress(i+1, len(packages), fmt.Sprintf("Installing %s", pkg))

		if IsPackageInstalled(pkg) {
			getOutputHandler().PrintInfo("%s is already installed", pkg)
			continue
		}

		if err := InstallPackageWithCheck(pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	return nil
}

// GetPackageInfo gets information about a package
func GetPackageInfo(packageName string) (*BrewPackage, error) {
	if !IsBrewInstalled() {
		return nil, fmt.Errorf("Homebrew is not installed")
	}

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewInfo, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to run brew info: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get info for %s: %s", packageName, result.Error)
	}

	// Parse the output to extract package information
	lines := strings.Split(result.Output, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no information found for %s", packageName)
	}

	pkg := &BrewPackage{
		Name:      packageName,
		Installed: IsPackageInstalled(packageName),
	}

	// Extract version and description from the first line
	firstLine := lines[0]
	if strings.Contains(firstLine, ":") {
		parts := strings.Split(firstLine, ":")
		if len(parts) > 1 {
			pkg.Description = strings.TrimSpace(parts[1])
		}
	}

	return pkg, nil
}

// IsApplicationAvailable checks if an application is available on the system
// Uses a hybrid approach: Homebrew detection -> intelligent search -> system-wide fallback
func IsApplicationAvailable(packageName string) bool {
	// Step 1: Quick Homebrew check (fastest)
	if IsPackageInstalled(packageName) {
		return true
	}

	// Step 2: Check if it's an installed Homebrew cask
	if isBrewCaskInstalled(packageName) {
		return true
	}

	// Step 3: Search for the cask and get actual install location
	if checkBrewCaskAvailable(packageName) {
		return true
	}

	// Step 4: Intelligent /Applications directory search
	if searchApplicationsDirectory(packageName) {
		return true
	}

	// Step 5: System-wide Spotlight search fallback
	if spotlightSearch(packageName) {
		return true
	}

	// Step 6: Final check - command-line tools in PATH
	result, err := system.RunCommand("which", packageName)
	return err == nil && result.Success
}

// isBrewCaskInstalled checks if package is in brew's installed cask list
func isBrewCaskInstalled(packageName string) bool {
	result, err := system.RunCommand(constants.BrewCommand, "list", "--cask")
	if err != nil {
		return false
	}

	// Check if packageName is in the output
	return strings.Contains(result.Output, packageName)
}

// checkBrewCaskAvailable searches for cask and checks if it's installed at the location brew expects
func checkBrewCaskAvailable(packageName string) bool {
	// Search for the cask to get its actual name
	result, err := system.RunCommand(constants.BrewCommand, "search", "--cask", packageName)
	if err != nil {
		return false
	}

	lines := strings.Split(strings.TrimSpace(result.Output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip headers and empty lines
		if line == "" || strings.Contains(line, "==>") {
			continue
		}

		// Found exact match or close match
		if line == packageName || strings.Contains(line, packageName) {
			// Get cask info to find install location
			infoResult, infoErr := system.RunCommand(constants.BrewCommand, "info", "--cask", line)
			if infoErr == nil && strings.Contains(infoResult.Output, "/Applications/") {
				// Extract app path from brew info output
				if appPath := extractAppPath(infoResult.Output); appPath != "" {
					result, err := system.RunCommand("test", "-d", appPath)
					return err == nil && result.Success
				}
			}
		}
	}
	return false
}

// searchApplicationsDirectory performs intelligent search in /Applications
func searchApplicationsDirectory(packageName string) bool {
	// Transform package name to likely app names
	possibleNames := generateAppNames(packageName)

	for _, appName := range possibleNames {
		appPath := "/Applications/" + appName
		result, err := system.RunCommand("test", "-d", appPath)
		if err == nil && result.Success {
			return true
		}
	}
	return false
}

// spotlightSearch uses macOS Spotlight to find applications system-wide
func spotlightSearch(packageName string) bool {
	// Use mdfind to search for applications containing the package name
	query := fmt.Sprintf("kMDItemKind == 'Application' && kMDItemFSName == '*%s*'", packageName)
	result, err := system.RunCommand("mdfind", query)

	if err != nil {
		return false
	}

	// If mdfind returns any results, the app exists somewhere
	return strings.TrimSpace(result.Output) != ""
}

// generateAppNames creates possible application names from package name
func generateAppNames(packageName string) []string {
	var names []string

	// Direct name with .app
	names = append(names, packageName+".app")
	names = append(names, strings.Title(packageName)+".app")

	// Handle hyphenated names
	if strings.Contains(packageName, "-") {
		// Convert hyphens to spaces and title case
		spacedName := strings.ReplaceAll(packageName, "-", " ")
		names = append(names, strings.Title(spacedName)+".app")

		// Remove hyphens entirely
		noDashName := strings.ReplaceAll(packageName, "-", "")
		names = append(names, strings.Title(noDashName)+".app")
	}

	// Handle common specific cases
	specialCases := map[string][]string{
		"visual-studio-code": {"Visual Studio Code.app"},
		"google-chrome":      {"Google Chrome.app"},
		"1password":          {"1Password 7 - Password Manager.app", "1Password.app"},
		"iterm2":             {"iTerm.app"},
	}

	if special, exists := specialCases[packageName]; exists {
		names = append(names, special...)
	}

	return names
}

// extractAppPath extracts the application path from brew info output
func extractAppPath(brewOutput string) string {
	lines := strings.Split(brewOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Applications/") && strings.Contains(line, ".app") {
			// Extract path - look for pattern like "/Applications/AppName.app"
			start := strings.Index(line, "/Applications/")
			if start == -1 {
				continue
			}

			end := strings.Index(line[start:], ".app")
			if end == -1 {
				continue
			}

			return line[start : start+end+4] // +4 for ".app"
		}
	}
	return ""
}

// isCaskPackage dynamically determines if a package is a Homebrew cask
func isCaskPackage(packageName string) bool {
	// First check if it exists as a cask
	result, err := system.RunCommand(constants.BrewCommand, "search", "--cask", packageName)
	if err == nil && result.Success {
		lines := strings.Split(strings.TrimSpace(result.Output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Skip headers, empty lines, and error messages
			if line == "" || strings.Contains(line, "==>") || strings.Contains(line, "Error:") || strings.Contains(line, "Warning:") {
				continue
			}
			// Only consider exact matches for casks to avoid false positives
			if line == packageName {
				return true
			}
		}
	}

	// Default to false (formula) if not a cask
	return false
}

// InstallPackageWithCheck installs a package only if it's not already available
func InstallPackageWithCheck(packageName string) error {
	if !IsBrewInstalled() {
		return fmt.Errorf("Homebrew is not installed")
	}

	// Check if application is already available (via any method)
	if IsApplicationAvailable(packageName) {
		getOutputHandler().PrintAlreadyAvailable("%s is already available on the system", packageName)
		return nil
	}

	// Dynamically determine if this is a cask (GUI app) or formula (CLI tool)
	isCask := isCaskPackage(packageName)

	getOutputHandler().PrintInfo("Installing %s...", packageName)

	var result *system.CommandResult
	var err error

	if isCask {
		// Install as cask
		result, err = system.RunCommand(constants.BrewCommand, constants.BrewInstall, "--cask", packageName)
	} else {
		// Install as formula
		result, err = system.RunCommand(constants.BrewCommand, constants.BrewInstall, packageName)
	}

	if err != nil {
		return fmt.Errorf("failed to run brew install: %w", err)
	}

	if !result.Success {
		// Check if the error is because the app already exists
		if strings.Contains(result.Error, "already an App at") {
			getOutputHandler().PrintAlreadyAvailable("%s is already installed manually, skipping Homebrew installation", packageName)
			return nil
		}

		// Return the actual brew output for clearer error messages
		if result.Output != "" {
			return fmt.Errorf("brew: %s", strings.TrimSpace(result.Output))
		} else {
			return fmt.Errorf("installation failed: %s", result.Error)
		}
	}

	return nil
}
