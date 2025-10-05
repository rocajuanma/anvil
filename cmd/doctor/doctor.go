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

package doctor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/pkg/constants"
	"github.com/rocajuanma/anvil/pkg/errors"
	"github.com/rocajuanma/anvil/pkg/validators"
	"github.com/rocajuanma/palantir"
	"github.com/spf13/cobra"
)

// getOutputHandler returns the global output handler for terminal operations
func getOutputHandler() palantir.OutputHandler {
	return palantir.GetGlobalOutputHandler()
}

var DoctorCmd = &cobra.Command{
	Use:   "doctor [category|check]",
	Short: "Run health checks and validate anvil environment",
	Long:  constants.DOCTOR_COMMAND_LONG_DESCRIPTION,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDoctorCommand(cmd, args); err != nil {
			getOutputHandler().PrintError("Doctor failed: %v", err)
			return
		}
	},
}

// runDoctorCommand executes the doctor validation process
func runDoctorCommand(cmd *cobra.Command, args []string) error {
	// Get command flags
	listChecks, _ := cmd.Flags().GetBool("list")
	fix, _ := cmd.Flags().GetBool("fix")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Create doctor engine with terminal output
	engine := validators.NewDoctorEngine(nil)

	// Handle list command
	if listChecks {
		return showAvailableChecks(engine)
	}

	// Handle fix command
	if fix {
		if len(args) > 0 {
			return runFixCheck(engine, args[0])
		} else {
			return runFixAll(engine, "")
		}
	}

	// Handle positional arguments
	if len(args) == 0 {
		// No arguments: run all checks
		return runAllChecks(engine, verbose)
	}

	target := args[0]

	// Check if it's a category first
	categories := []string{"environment", "dependencies", "configuration", "connectivity"}
	for _, category := range categories {
		if target == category {
			return runCategoryChecks(engine, category, verbose)
		}
	}

	// Otherwise treat it as a specific check
	return runSingleCheck(engine, target, verbose)
}

// showAvailableChecks displays all available checks organized by category
func showAvailableChecks(engine *validators.DoctorEngine) error {
	o := getOutputHandler()
	o.PrintHeader("Available Health Checks")

	o.PrintInfo("🏷️  CATEGORIES (run all checks in a group):\n")

	checks := engine.ListChecks()

	categoryDescriptions := map[string]string{
		"environment":   "Verify anvil initialization and directory structure",
		"dependencies":  "Check required tools and Homebrew installation",
		"configuration": "Validate git and GitHub settings",
		"connectivity":  "Test GitHub access and repository connections",
	}

	totalChecks := 0
	for _, category := range []string{"environment", "dependencies", "configuration", "connectivity"} {
		if checkNames, exists := checks[category]; exists {
			o.PrintStage(fmt.Sprintf("anvil doctor %s", category))
			o.PrintInfo("    %s", categoryDescriptions[category])
			o.PrintInfo("    Includes: %s", strings.Join(checkNames, ", "))
			o.PrintInfo("    (%d checks)\n", len(checkNames))
			totalChecks += len(checkNames)
		}
	}

	o.PrintInfo("🔍 SPECIFIC CHECKS (run individual validators):\n")

	for _, category := range []string{"environment", "dependencies", "configuration", "connectivity"} {
		if checkNames, exists := checks[category]; exists {
			o.PrintStage(fmt.Sprintf("%s checks:", strings.Title(category)))
			for _, checkName := range checkNames {
				o.PrintInfo("  anvil doctor %s", checkName)
			}
			o.PrintInfo("")
		}
	}

	o.PrintInfo("💡 USAGE EXAMPLES:\n")
	o.PrintInfo("  anvil doctor                    # Run all %d checks", totalChecks)
	o.PrintInfo("  anvil doctor environment        # Run 3 environment checks")
	o.PrintInfo("  anvil doctor git-config         # Run only git configuration check")
	o.PrintInfo("  anvil doctor --fix              # Auto-fix detected issues")
	o.PrintInfo("  anvil doctor dependencies --fix # Auto-fix dependency issues")

	return nil
}

// runSingleCheck executes a specific health check
func runSingleCheck(engine *validators.DoctorEngine, checkName string, verbose bool) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Running Check: %s", checkName))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	o.PrintStage(fmt.Sprintf("Executing %s check...", checkName))

	result := engine.RunCheckWithProgress(ctx, checkName, verbose)
	displayResults([]*validators.ValidationResult{result}, verbose)

	if result.Status == validators.FAIL {
		return errors.NewValidationError(constants.OpDoctor, checkName, fmt.Errorf(result.Message))
	}

	return nil
}

// runCategoryChecks executes all checks in a specific category
func runCategoryChecks(engine *validators.DoctorEngine, category string, verbose bool) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Running %s Health Checks", strings.Title(category)))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get validators for this category to show count
	categoryValidators := engine.GetValidatorsByCategory(category)
	if len(categoryValidators) == 0 {
		o.PrintError("Category '%s' not found", category)
		return errors.NewValidationError(constants.OpDoctor, category, fmt.Errorf("category not found"))
	}

	o.PrintStage(fmt.Sprintf("Executing %d checks in %s category...", len(categoryValidators), category))

	results := engine.RunCategoryWithProgress(ctx, category, verbose)
	displayResults(results, verbose)
	printSummary(results)

	// Check if any critical failures
	for _, result := range results {
		if result.Status == validators.FAIL {
			return errors.NewValidationError(constants.OpDoctor, category, fmt.Errorf("validation failures detected"))
		}
	}

	return nil
}

// runAllChecks executes all available health checks
func runAllChecks(engine *validators.DoctorEngine, verbose bool) error {
	o := getOutputHandler()
	o.PrintHeader("Running Anvil Health Check")
	o.PrintInfo("🔍 Validating environment, dependencies, configuration, and connectivity...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Get all validators to show total count
	allValidators := engine.GetAllValidators()
	totalChecks := len(allValidators)

	o.PrintStage(fmt.Sprintf("Executing %d health checks...", totalChecks))

	results := engine.RunAllWithProgress(ctx, verbose)
	displayResults(results, verbose)
	printSummary(results)

	return nil
}

// runFixCheck attempts to fix a specific check
func runFixCheck(engine *validators.DoctorEngine, checkName string) error {
	o := getOutputHandler()
	o.PrintHeader(fmt.Sprintf("Fixing Check: %s", checkName))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// First run the check to see current status
	result := engine.RunCheck(ctx, checkName)
	displayResults([]*validators.ValidationResult{result}, false)

	if result.Status == validators.PASS && checkName != "git-config" {
		o.PrintSuccess("Check is already passing, no fix needed")
		return nil
	}

	if !result.AutoFix {
		o.PrintWarning("This check cannot be automatically fixed")
		o.PrintInfo("Manual fix required: %s", result.FixHint)
		return nil
	}

	// Confirm with user
	if !o.Confirm(fmt.Sprintf("Attempt to fix '%s'?", checkName)) {
		o.PrintInfo("Fix cancelled by user")
		return nil
	}

	// Attempt fix
	o.PrintInfo("Attempting to fix %s...", checkName)
	if err := engine.FixCheck(ctx, checkName); err != nil {
		o.PrintError("Fix failed: %v", err)
		return err
	}

	o.PrintSuccess("Fix completed!")

	// Verify fix
	o.PrintInfo("Verifying fix...")
	newResult := engine.RunCheck(ctx, checkName)
	displayResults([]*validators.ValidationResult{newResult}, false)

	if newResult.Status == validators.PASS {
		o.PrintSuccess("✅ Check is now passing!")
	} else {
		o.PrintWarning("⚠️  Check still has issues after fix attempt")
	}

	return nil
}

// runFixAll attempts to fix all auto-fixable issues
func runFixAll(engine *validators.DoctorEngine, category string) error {
	o := getOutputHandler()
	o.PrintHeader("Auto-fixing Issues")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Run checks to find fixable issues
	var results []*validators.ValidationResult
	if category != "" {
		results = engine.RunCategory(ctx, category)
	} else {
		results = engine.RunAll(ctx)
	}

	fixableIssues := validators.GetFixableIssues(results)
	if len(fixableIssues) == 0 {
		if category != "" {
			o.PrintSuccess(fmt.Sprintf("No auto-fixable issues found in %s category!", category))
		} else {
			o.PrintSuccess("No auto-fixable issues found!")
		}
		return nil
	}

	if category != "" {
		o.PrintInfo("Found %d auto-fixable issues in %s category:", len(fixableIssues), category)
	} else {
		o.PrintInfo("Found %d auto-fixable issues:", len(fixableIssues))
	}
	for _, issue := range fixableIssues {
		o.PrintInfo("  • %s: %s", issue.Name, issue.Message)
	}

	confirmMessage := "Attempt to fix all auto-fixable issues?"
	if category != "" {
		confirmMessage = fmt.Sprintf("Attempt to fix all auto-fixable issues in %s category?", category)
	}

	if !o.Confirm(confirmMessage) {
		o.PrintInfo("Fix cancelled by user")
		return nil
	}

	var fixedCount, failedCount int
	for _, issue := range fixableIssues {
		o.PrintInfo("Fixing %s...", issue.Name)
		if err := engine.FixCheck(ctx, issue.Name); err != nil {
			o.PrintError("Failed to fix %s: %v", issue.Name, err)
			failedCount++
		} else {
			o.PrintSuccess(fmt.Sprintf("Fixed %s", issue.Name))
			fixedCount++
		}
	}

	o.PrintInfo("Fix complete: %d succeeded, %d failed", fixedCount, failedCount)
	return nil
}

// displayResults shows validation results in a formatted table
func displayResults(results []*validators.ValidationResult, verbose bool) {
	categories := validators.FormatResultsTable(results)

	for _, category := range []string{"environment", "dependencies", "configuration", "connectivity"} {
		if categoryResults, exists := categories[category]; exists {
			displayCategory(category, categoryResults, verbose)
		}
	}
}

// displayCategory shows results for a specific category
func displayCategory(category string, results []*validators.ValidationResult, verbose bool) {
	// Count statuses
	passed, warned, failed, skipped := 0, 0, 0, 0
	for _, result := range results {
		switch result.Status {
		case validators.PASS:
			passed++
		case validators.WARN:
			warned++
		case validators.FAIL:
			failed++
		case validators.SKIP:
			skipped++
		}
	}

	// Choose category status
	categoryStatus := getCategoryStatus(passed, warned, failed, skipped)
	o := getOutputHandler()
	o.PrintStage(fmt.Sprintf("%s %s", categoryStatus, strings.Title(category)))

	for _, result := range results {
		statusEmoji := getStatusEmoji(result.Status)
		o.PrintInfo("  %s %s", statusEmoji, result.Message)

		if verbose && len(result.Details) > 0 {
			for _, detail := range result.Details {
				o.PrintInfo("      %s", detail)
			}
		}

		if result.Status != validators.PASS && result.FixHint != "" {
			o.PrintInfo("      💡 %s", result.FixHint)
		}
	}

	o.PrintInfo("")
}

// printSummary shows overall health check summary
func printSummary(results []*validators.ValidationResult) {
	passed, warned, failed, skipped := validators.GetSummary(results)
	total := len(results)

	o := getOutputHandler()
	o.PrintHeader("Health Check Summary")
	o.PrintInfo("Total checks: %d", total)
	o.PrintInfo("✅ Passed: %d", passed)
	if warned > 0 {
		o.PrintInfo("⚠️  Warnings: %d", warned)
	}
	if failed > 0 {
		o.PrintInfo("❌ Failed: %d", failed)
	}
	if skipped > 0 {
		o.PrintInfo("⏭️  Skipped: %d", skipped)
	}

	// Show fixable issues
	fixableIssues := validators.GetFixableIssues(results)
	if len(fixableIssues) > 0 {
		o.PrintInfo("\n🔧 %d issues can be auto-fixed", len(fixableIssues))
		o.PrintInfo("Run 'anvil doctor --fix' to automatically fix them")
	}

	// Overall status
	if failed > 0 {
		o.PrintWarning("❌ Overall status: Issues found")
	} else if warned > 0 {
		o.PrintInfo("⚠️  Overall status: Minor issues")
	} else {
		o.PrintSuccess("✅ Overall status: Healthy")
	}
}

// Helper functions for display formatting
func getCategoryEmoji(category string) string {
	switch category {
	case "environment":
		return "🏠"
	case "dependencies":
		return "📦"
	case "configuration":
		return "⚙️"
	case "connectivity":
		return "🌐"
	default:
		return "🔍"
	}
}

func getCategoryStatus(passed, warned, failed, skipped int) string {
	if failed > 0 {
		return "❌"
	} else if warned > 0 {
		return "⚠️ "
	} else {
		return "✅"
	}
}

func getStatusEmoji(status validators.ValidationStatus) string {
	switch status {
	case validators.PASS:
		return "✅"
	case validators.WARN:
		return "⚠️ "
	case validators.FAIL:
		return "❌"
	case validators.SKIP:
		return "⏭️ "
	default:
		return "❓"
	}
}

func init() {
	// Add flags for enhanced doctor functionality
	DoctorCmd.Flags().Bool("list", false, "List all available health checks")
	DoctorCmd.Flags().Bool("fix", false, "Attempt to automatically fix issues")
	DoctorCmd.Flags().Bool("verbose", false, "Show detailed output")
}
