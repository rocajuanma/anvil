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

	"github.com/rocajuanma/anvil/internal/constants"
	"github.com/rocajuanma/anvil/internal/errors"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/validators"
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

	spinner := charm.NewLineSpinner(fmt.Sprintf("Executing %s check", checkName))
	spinner.Start()

	result := engine.RunCheckWithProgress(ctx, checkName, verbose)

	if result.Status == validators.PASS {
		spinner.Success(fmt.Sprintf("%s check passed", checkName))
	} else if result.Status == validators.WARN {
		spinner.Warning(fmt.Sprintf("%s check completed with warnings", checkName))
	} else {
		spinner.Error(fmt.Sprintf("%s check failed", checkName))
	}

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

	spinner := charm.NewLineSpinner(fmt.Sprintf("Executing %d checks in %s category", len(categoryValidators), category))
	spinner.Start()

	results := engine.RunCategoryWithProgress(ctx, category, verbose)

	// Count status
	passed, warned, failed := 0, 0, 0
	for _, result := range results {
		if result.Status == validators.PASS {
			passed++
		} else if result.Status == validators.WARN {
			warned++
		} else if result.Status == validators.FAIL {
			failed++
		}
	}

	if failed > 0 {
		spinner.Error(fmt.Sprintf("%s checks completed: %d failed", category, failed))
	} else if warned > 0 {
		spinner.Warning(fmt.Sprintf("%s checks completed: %d warnings", category, warned))
	} else {
		spinner.Success(fmt.Sprintf("All %s checks passed", category))
	}

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

// checkStatus represents the status of an individual health check
type checkStatus struct {
	name    string
	status  string // "pending", "checking", "pass", "warn", "fail"
	emoji   string
	message string
}

// runAllChecks executes all available health checks
func runAllChecks(engine *validators.DoctorEngine, verbose bool) error {
	o := getOutputHandler()
	o.PrintHeader("Running Anvil Health Check")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Get all validators and initialize status tracking
	allValidators := engine.GetAllValidators()
	totalChecks := len(allValidators)

	// Group validators by category
	categories := map[string][]string{
		"environment":   {},
		"dependencies":  {},
		"configuration": {},
		"connectivity":  {},
	}

	for _, v := range allValidators {
		categories[v.Category()] = append(categories[v.Category()], v.Name())
	}

	// Initialize check statuses
	checkStatuses := make(map[string]*checkStatus)
	for _, v := range allValidators {
		checkStatuses[v.Name()] = &checkStatus{
			name:   v.Name(),
			status: "pending",
			emoji:  "⋯",
		}
	}

	// Run checks with a spinner
	spinner := charm.NewLineSpinner(fmt.Sprintf("Running %d health checks", totalChecks))
	spinner.Start()

	results := engine.RunAll(ctx)

	// Count results for spinner message
	passed, warned, failed := 0, 0, 0
	for _, result := range results {
		switch result.Status {
		case validators.PASS:
			passed++
		case validators.WARN:
			warned++
		case validators.FAIL:
			failed++
		}
	}

	// Update spinner based on results
	if failed > 0 {
		spinner.Error(fmt.Sprintf("Completed: %d passed, %d failed", passed, failed))
	} else if warned > 0 {
		spinner.Warning(fmt.Sprintf("Completed: %d passed, %d warnings", passed, warned))
	} else {
		spinner.Success(fmt.Sprintf("All %d checks passed!", totalChecks))
	}

	// Update statuses based on results
	for _, result := range results {
		if cs, exists := checkStatuses[result.Name]; exists {
			cs.message = result.Message
			switch result.Status {
			case validators.PASS:
				cs.status = "pass"
				cs.emoji = "✓"
			case validators.WARN:
				cs.status = "warn"
				cs.emoji = "⚠"
			case validators.FAIL:
				cs.status = "fail"
				cs.emoji = "✗"
			case validators.SKIP:
				cs.status = "skip"
				cs.emoji = "○"
			}
		}
	}

	// Print organized results by category
	fmt.Println()
	for _, category := range []string{"environment", "dependencies", "configuration", "connectivity"} {
		if checkNames, exists := categories[category]; exists && len(checkNames) > 0 {
			printCategoryResults(category, checkNames, checkStatuses, results, verbose)
		}
	}

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
	spinner := charm.NewDotsSpinner(fmt.Sprintf("Attempting to fix %s", checkName))
	spinner.Start()
	if err := engine.FixCheck(ctx, checkName); err != nil {
		spinner.Error("Fix failed")
		o.PrintError("Fix failed: %v", err)
		return err
	}

	spinner.Success("Fix completed!")

	// Verify fix
	spinner = charm.NewLineSpinner("Verifying fix")
	spinner.Start()
	newResult := engine.RunCheck(ctx, checkName)
	spinner.Success("Verification complete")
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

// printCategoryResults prints results for a single category in a clean format
func printCategoryResults(category string, checkNames []string, statuses map[string]*checkStatus, results []*validators.ValidationResult, verbose bool) {
	// Count status for this category
	passed, warned, failed := 0, 0, 0
	for _, name := range checkNames {
		if cs, exists := statuses[name]; exists {
			switch cs.status {
			case "pass":
				passed++
			case "warn":
				warned++
			case "fail":
				failed++
			}
		}
	}

	categoryStatus := getCategoryStatus(passed, warned, failed, 0)
	categoryTitle := strings.Title(category)

	// Print category header with emoji
	fmt.Printf("  %s %s\n", categoryStatus, charm.RenderHighlight(categoryTitle, "#00D9FF"))

	// Print each check result
	for _, name := range checkNames {
		if cs, exists := statuses[name]; exists {
			// Get full result for this check
			var result *validators.ValidationResult
			for _, r := range results {
				if r.Name == name {
					result = r
					break
				}
			}

			if result != nil {
				fmt.Printf("    %s %s\n", cs.emoji, result.Message)

				// Show fix hint for failed/warned checks
				if (result.Status == validators.FAIL || result.Status == validators.WARN) && result.FixHint != "" && !verbose {
					fmt.Printf("      💡 %s\n", result.FixHint)
				}

				// Show details in verbose mode
				if verbose && len(result.Details) > 0 {
					for _, detail := range result.Details {
						fmt.Printf("        %s\n", detail)
					}
				}
			}
		}
	}

	fmt.Println()
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
	passed, warned, failed, _ := validators.GetSummary(results)
	total := len(results)

	// Get category breakdown
	categoryStats := getCategoryBreakdown(results)

	// Build dashboard box
	var dashboardContent strings.Builder
	dashboardContent.WriteString("\n")

	// Category status bars
	for _, category := range []string{"environment", "dependencies", "configuration", "connectivity"} {
		if stats, exists := categoryStats[category]; exists {
			status := getCategoryStatus(stats.passed, stats.warned, stats.failed, stats.skipped)

			// Calculate percentage for progress bar
			percentage := 0
			if stats.total > 0 {
				percentage = (stats.passed * 100) / stats.total
			}

			// Create progress bar (20 chars wide)
			barWidth := 20
			filled := (percentage * barWidth) / 100
			bar := strings.Repeat("━", filled) + strings.Repeat("━", barWidth-filled)

			categoryTitle := strings.Title(category)
			dashboardContent.WriteString(fmt.Sprintf("  %s %-15s %s  %d/%d passing\n",
				status, categoryTitle, bar, stats.passed, stats.total))
		}
	}

	dashboardContent.WriteString("\n")
	dashboardContent.WriteString(fmt.Sprintf("  Overall: %d/%d checks passing\n", passed, total))
	dashboardContent.WriteString("\n")

	// Render dashboard box
	fmt.Println(charm.RenderBox("Summary", dashboardContent.String(), "#00D9FF", true))

	// Show fixable issues in a separate box
	fixableIssues := validators.GetFixableIssues(results)
	if len(fixableIssues) > 0 {
		var fixContent strings.Builder
		for _, issue := range fixableIssues {
			fixContent.WriteString(fmt.Sprintf("  • %s\n", issue.Name))
		}
		fixContent.WriteString("\n")
		fixContent.WriteString("  Run 'anvil doctor --fix' to automatically fix them\n")

		fmt.Println(charm.RenderBox("🔧 Auto-fixable Issues", fixContent.String(), "#FFD700", true))
	}

	// Overall status badge
	fmt.Println()
	if failed > 0 {
		fmt.Println("  " + charm.RenderBadge("ISSUES FOUND", "#FF5F87"))
	} else if warned > 0 {
		fmt.Println("  " + charm.RenderBadge("MINOR ISSUES", "#FFD700"))
	} else {
		fmt.Println("  " + charm.RenderBadge("HEALTHY", "#00FF87"))
	}
	fmt.Println()
}

// categoryStats holds statistics for a category
type categoryStats struct {
	passed  int
	warned  int
	failed  int
	skipped int
	total   int
}

// getCategoryBreakdown returns stats broken down by category
func getCategoryBreakdown(results []*validators.ValidationResult) map[string]categoryStats {
	breakdown := make(map[string]categoryStats)

	for _, result := range results {
		stats := breakdown[result.Category]
		stats.total++

		switch result.Status {
		case validators.PASS:
			stats.passed++
		case validators.WARN:
			stats.warned++
		case validators.FAIL:
			stats.failed++
		case validators.SKIP:
			stats.skipped++
		}

		breakdown[result.Category] = stats
	}

	return breakdown
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
