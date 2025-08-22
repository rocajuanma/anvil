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

package push

import (
	"fmt"
	"strings"

	"github.com/rocajuanma/anvil/pkg/config"
	"github.com/rocajuanma/anvil/pkg/github"
	"github.com/rocajuanma/anvil/pkg/terminal"
)

// handleAppLocationError provides helpful error messages for app location resolution failures
func handleAppLocationError(appName string, err error) error {
	if strings.Contains(err.Error(), "not found in configs or temp directory") {
		terminal.PrintError("App '%s' is not known to anvil", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("💡 To push app configurations:")
		terminal.PrintInfo("")
		terminal.PrintInfo("1. Configure the app's local config path in settings.yaml:")
		terminal.PrintInfo("")
		terminal.PrintInfo("configs:")
		terminal.PrintInfo("  %s: /path/to/your/%s/configs", appName, appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("2. Or pull the app's configs first to discover it:")
		terminal.PrintInfo("   anvil config pull %s", appName)
		terminal.PrintInfo("")
		terminal.PrintInfo("3. Then configure the local path in settings.yaml")
		return fmt.Errorf("app not configured")
	}

	return fmt.Errorf("failed to resolve app location: %w", err)
}

// showSecurityWarning displays a security warning about private repositories
func showSecurityWarning(privateRepo string) {
	// 🚨 SECURITY WARNING: Remind users about private repository requirement
	terminal.PrintWarning("🔒 SECURITY REMINDER: Configuration files contain sensitive data")
	terminal.PrintInfo("   • API keys, tokens, and credentials")
	terminal.PrintInfo("   • Personal file paths and system information")
	terminal.PrintInfo("   • Private development environment details")
	terminal.PrintInfo("")
	terminal.PrintInfo("🛡️  Anvil ENFORCES private repositories for security")
	terminal.PrintInfo("   • Repository '%s' must be PRIVATE", privateRepo)
	terminal.PrintInfo("   • Public repositories will be BLOCKED\n")
}

// displaySuccessMessage displays a success message after the push operation
func displaySuccessMessage(appName string, result *github.PushConfigResult, diffSummary *github.DiffSummary, anvilConfig *config.AnvilConfig) {
	// Display full success message for actual push
	terminal.PrintHeader("Push Complete!")
	terminal.PrintSuccess(fmt.Sprintf("%s configuration push completed successfully!\n", appName))
	terminal.PrintInfo("📋 Push Summary:")
	terminal.PrintInfo("  • Branch created: %s", result.BranchName)
	terminal.PrintInfo("  • Commit message: %s", result.CommitMessage)
	terminal.PrintInfo("  • Files committed: \n\n%s", diffSummary.GitStatOutput)
	terminal.PrintInfo("🔗 Repository: %s", result.RepositoryURL)
	terminal.PrintInfo("🌿 Branch: %s\n", result.BranchName)
	terminal.PrintSuccess("You can now create a Pull Request on GitHub to merge these changes!")
	terminal.PrintInfo("Direct link: %s/compare/%s...%s", result.RepositoryURL, anvilConfig.GitHub.Branch, result.BranchName)
}

// showDiffOutput displays diff information using Git's native output
func showDiffOutput(diffSummary *github.DiffSummary) {
	if diffSummary.TotalFiles == 0 {
		terminal.PrintInfo("No changes detected")
		return
	}

	terminal.PrintInfo("")
	terminal.PrintHeader("📋 Changes to be pushed:")

	// Show Git's native stat output directly
	if diffSummary.GitStatOutput != "" {
		terminal.PrintInfo("")
		terminal.PrintInfo(diffSummary.GitStatOutput)
	}

	// For single small files, show full diff
	if diffSummary.TotalFiles == 1 && diffSummary.FullDiff != "" {
		lines := strings.Split(diffSummary.FullDiff, "\n")
		if len(lines) <= 50 {
			terminal.PrintInfo("")
			terminal.PrintInfo("📄 Full diff:")
			terminal.PrintInfo("")
			terminal.PrintInfo(diffSummary.FullDiff)
		} else {
			terminal.PrintInfo("")
			terminal.PrintInfo("📄 Diff preview (first 50 lines):")
			terminal.PrintInfo("")
			terminal.PrintInfo(strings.Join(lines[:50], "\n"))
			terminal.PrintInfo("")
			terminal.PrintInfo("... [diff truncated] ...")
		}
	}
	terminal.PrintInfo("")
}
