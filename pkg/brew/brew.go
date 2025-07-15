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
	"github.com/rocajuanma/anvil/pkg/system"
	"github.com/rocajuanma/anvil/pkg/terminal"
)

// BrewPackage represents a brew package
type BrewPackage struct {
	Name        string
	Version     string
	Description string
	Installed   bool
}

// IsBrewInstalled checks if Homebrew is installed
func IsBrewInstalled() bool {
	return system.CommandExists(constants.BrewCommand)
}

// InstallBrew installs Homebrew if not already installed
func InstallBrew() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("Homebrew is only supported on macOS")
	}

	if IsBrewInstalled() {
		return nil
	}

	terminal.PrintInfo("Installing Homebrew...")

	// Official Homebrew installation command
	installCmd := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`

	result, err := system.RunCommand("/bin/bash", "-c", installCmd)
	if err != nil {
		return fmt.Errorf("failed to run brew installation command: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("brew installation failed: %s", result.Error)
	}

	return nil
}

// UpdateBrew updates Homebrew and its formulae
func UpdateBrew() error {
	if !IsBrewInstalled() {
		return fmt.Errorf("Homebrew is not installed")
	}

	terminal.PrintInfo("Updating Homebrew...")

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

	terminal.PrintInfo("Installing %s...", packageName)

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewInstall, packageName)
	if err != nil {
		return fmt.Errorf("failed to run brew install: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("failed to install %s: %s", packageName, result.Error)
	}

	return nil
}

// IsPackageInstalled checks if a package is installed
func IsPackageInstalled(packageName string) bool {
	if !IsBrewInstalled() {
		return false
	}

	result, err := system.RunCommand(constants.BrewCommand, constants.BrewList, "--formula", packageName)
	if err != nil {
		return false
	}

	return result.Success
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
		terminal.PrintProgress(i+1, len(packages), fmt.Sprintf("Installing %s", pkg))

		if IsPackageInstalled(pkg) {
			terminal.PrintInfo("%s is already installed", pkg)
			continue
		}

		if err := InstallPackage(pkg); err != nil {
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
