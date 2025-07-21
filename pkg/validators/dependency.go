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

package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/rocajuanma/anvil/pkg/brew"
	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/system"
)

// BrewValidator checks if Homebrew is installed and functional
type BrewValidator struct{}

func (v *BrewValidator) Name() string        { return "homebrew" }
func (v *BrewValidator) Category() string    { return "dependencies" }
func (v *BrewValidator) Description() string { return "Verify Homebrew is installed and functional" }
func (v *BrewValidator) CanFix() bool        { return true }

func (v *BrewValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	// Check if brew is installed
	if !brew.IsBrewInstalled() {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Homebrew is not installed",
			Details:  []string{"Homebrew is required for app installation"},
			FixHint:  "Homebrew will be installed automatically",
			AutoFix:  true,
		}
	}

	// Check if brew is functional by running brew --version
	result, err := system.RunCommand("brew", "--version")
	if err != nil || !result.Success {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  "Homebrew is not functional",
			Details:  []string{"brew --version failed"},
			FixHint:  "Try running 'brew doctor' to diagnose issues",
			AutoFix:  false,
		}
	}

	// Check if brew needs updating (warn only)
	updateResult, err := system.RunCommand("brew", "outdated", "--quiet")
	if err == nil && strings.TrimSpace(updateResult.Output) != "" {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   WARN,
			Message:  "Homebrew has available updates",
			Details:  []string{"Run 'brew update && brew upgrade' to update"},
			FixHint:  "Homebrew will be updated",
			AutoFix:  true,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  "Homebrew is installed and functional",
		Details:  []string{strings.TrimSpace(result.Output)},
		AutoFix:  false,
	}
}

func (v *BrewValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	if !brew.IsBrewInstalled() {
		// Install Homebrew
		if err := brew.InstallBrew(); err != nil {
			return fmt.Errorf("failed to install Homebrew: %w", err)
		}
	} else {
		// Update Homebrew
		result, err := system.RunCommand("brew", "update")
		if err != nil {
			return fmt.Errorf("failed to update Homebrew: %w", err)
		}
		if !result.Success {
			return fmt.Errorf("brew update failed")
		}
	}
	return nil
}

// RequiredToolsValidator checks if all required tools are installed
type RequiredToolsValidator struct{}

func (v *RequiredToolsValidator) Name() string     { return "required-tools" }
func (v *RequiredToolsValidator) Category() string { return "dependencies" }
func (v *RequiredToolsValidator) Description() string {
	return "Verify all required tools are installed"
}
func (v *RequiredToolsValidator) CanFix() bool { return true }

func (v *RequiredToolsValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	requiredTools := cfg.Tools.RequiredTools
	if len(requiredTools) == 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   PASS,
			Message:  "No required tools configured",
			AutoFix:  false,
		}
	}

	var missingTools []string
	var installedTools []string

	for _, tool := range requiredTools {
		if brew.IsApplicationAvailable(tool) {
			installedTools = append(installedTools, tool)
		} else {
			missingTools = append(missingTools, tool)
		}
	}

	if len(missingTools) > 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   FAIL,
			Message:  fmt.Sprintf("Missing required tools: %s", strings.Join(missingTools, ", ")),
			Details:  []string{fmt.Sprintf("Installed: %d/%d", len(installedTools), len(requiredTools))},
			FixHint:  "Missing tools will be installed automatically",
			AutoFix:  true,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  fmt.Sprintf("All required tools installed (%d/%d)", len(installedTools), len(requiredTools)),
		Details:  installedTools,
		AutoFix:  false,
	}
}

func (v *RequiredToolsValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	requiredTools := cfg.Tools.RequiredTools
	var installErrors []string

	for _, tool := range requiredTools {
		if !brew.IsApplicationAvailable(tool) {
			if err := brew.InstallPackageWithCheck(tool); err != nil {
				installErrors = append(installErrors, fmt.Sprintf("%s: %v", tool, err))
			}
		}
	}

	if len(installErrors) > 0 {
		return fmt.Errorf("failed to install some tools: %s", strings.Join(installErrors, "; "))
	}

	return nil
}

// OptionalToolsValidator checks the status of optional tools
type OptionalToolsValidator struct{}

func (v *OptionalToolsValidator) Name() string        { return "optional-tools" }
func (v *OptionalToolsValidator) Category() string    { return "dependencies" }
func (v *OptionalToolsValidator) Description() string { return "Check status of optional tools" }
func (v *OptionalToolsValidator) CanFix() bool        { return false }

func (v *OptionalToolsValidator) Validate(ctx context.Context, cfg *config.AnvilConfig) *ValidationResult {
	optionalTools := cfg.Tools.OptionalTools
	if len(optionalTools) == 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   PASS,
			Message:  "No optional tools configured",
			AutoFix:  false,
		}
	}

	var installedTools []string
	var missingTools []string

	for _, tool := range optionalTools {
		if brew.IsApplicationAvailable(tool) {
			installedTools = append(installedTools, tool)
		} else {
			missingTools = append(missingTools, tool)
		}
	}

	if len(missingTools) > 0 {
		return &ValidationResult{
			Name:     v.Name(),
			Category: v.Category(),
			Status:   WARN,
			Message:  fmt.Sprintf("Optional tools missing: %s", strings.Join(missingTools, ", ")),
			Details:  []string{fmt.Sprintf("Installed: %d/%d", len(installedTools), len(optionalTools))},
			FixHint:  fmt.Sprintf("Run 'anvil install %s' to install missing tools", strings.Join(missingTools, " ")),
			AutoFix:  false,
		}
	}

	return &ValidationResult{
		Name:     v.Name(),
		Category: v.Category(),
		Status:   PASS,
		Message:  fmt.Sprintf("All optional tools installed (%d/%d)", len(installedTools), len(optionalTools)),
		Details:  installedTools,
		AutoFix:  false,
	}
}

func (v *OptionalToolsValidator) Fix(ctx context.Context, cfg *config.AnvilConfig) error {
	// Optional tools are not automatically fixed
	return fmt.Errorf("optional tools must be installed manually using 'anvil install [tool]'")
}
