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

package charm

import (
	"fmt"
	"time"

	"github.com/rocajuanma/palantir"
)

// RunDemo demonstrates all the Charm terminal features
// This is useful for testing and showcasing the enhanced UI
func RunDemo() {
	o := palantir.GetGlobalOutputHandler()

	// Banner
	fmt.Println(RenderBanner("ANVIL CHARM DEMO"))
	fmt.Println()

	// Section 1: Basic Output Types
	fmt.Println(RenderBox("Basic Output Types", "Enhanced versions of standard output", "#00D9FF"))
	o.PrintHeader("This is a header")
	o.PrintStage("This is a stage marker")
	o.PrintSuccess("This is a success message")
	o.PrintError("This is an error message")
	o.PrintWarning("This is a warning message")
	o.PrintInfo("This is an info message")
	o.PrintAlreadyAvailable("This item is already available")
	fmt.Println()

	// Section 2: Progress Bars
	fmt.Println(RenderBox("Progress Indicators", "Visual progress tracking", "#FFD700"))
	for i := 1; i <= 5; i++ {
		o.PrintProgress(i, 5, fmt.Sprintf("Installing package %d", i))
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println()

	// Section 3: Spinners
	fmt.Println(RenderBox("Spinners", "Animated loading indicators", "#00FF87"))

	spinner := NewDotsSpinner("Processing with dots spinner")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Success("Dots spinner completed!")

	spinner = NewLineSpinner("Processing with line spinner")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Success("Line spinner completed!")

	spinner = NewCircleSpinner("Processing with circle spinner")
	spinner.Start()
	time.Sleep(2 * time.Second)
	spinner.Success("Circle spinner completed!")
	fmt.Println()

	// Section 4: Visual Components
	fmt.Println(RenderBox("Visual Components", "Various UI elements", "#C792EA"))

	// Lists
	items := []string{
		"Git installed and configured",
		"Homebrew up to date",
		"Configuration validated",
		"All dependencies ready",
	}
	fmt.Println(RenderList(items, "✓", "#00FF87"))

	// Key-Value Pairs
	fmt.Println(RenderKeyValue("Version:", "2.0.0"))
	fmt.Println(RenderKeyValue("Status:", "Ready"))
	fmt.Println(RenderKeyValue("Platform:", "macOS"))
	fmt.Println()

	// Badges
	fmt.Print(RenderBadge("SUCCESS", "#00FF87"))
	fmt.Print(" ")
	fmt.Print(RenderBadge("WARNING", "#FFD700"))
	fmt.Print(" ")
	fmt.Print(RenderBadge("ERROR", "#FF5F87"))
	fmt.Print(" ")
	fmt.Println(RenderBadge("INFO", "#00D9FF"))
	fmt.Println()

	// Steps
	steps := []string{
		"Initialize environment",
		"Install dependencies",
		"Configure settings",
		"Verify installation",
	}
	fmt.Println(RenderSteps(steps))

	// Code and Highlights
	fmt.Println(RenderHighlight("IMPORTANT:", "#FFD700"), "Remember to commit your changes")
	fmt.Println("Run:", RenderCode("anvil install git"))
	fmt.Println()

	// Quote
	quote := RenderQuote(
		"The best time to plant a tree was 20 years ago. The second best time is now.",
		"Chinese Proverb",
	)
	fmt.Println(quote)
	fmt.Println()

	// Status indicators
	fmt.Println(RenderStatus("All systems operational", true))
	fmt.Println(RenderStatus("Service degraded", false))
	fmt.Println()

	// Percentages
	fmt.Println("Coverage:", RenderPercentage(95.5))
	fmt.Println("Health:", RenderPercentage(67.3))
	fmt.Println("Quality:", RenderPercentage(42.1))
	fmt.Println()

	// Separator
	fmt.Println(RenderSeparator(60, "─", "#666666"))

	// Final message
	fmt.Println(RenderBox(
		"Demo Complete!",
		"All Charm terminal features demonstrated.\nThese components are now available throughout Anvil.",
		"#FF6B9D",
	))
}

// RunQuickDemo shows a quick demonstration suitable for testing
func RunQuickDemo() {
	o := palantir.GetGlobalOutputHandler()

	fmt.Println(RenderBanner("CHARM QUICK DEMO"))

	o.PrintHeader("Testing Enhanced Output")
	o.PrintStage("Running quick test")

	spinner := NewDotsSpinner("Processing")
	spinner.Start()
	time.Sleep(1 * time.Second)
	spinner.Success("Test completed!")

	o.PrintSuccess("All features working correctly")

	fmt.Println()
	fmt.Println(RenderBox("Status", "Charm terminal is ready!", "#00FF87"))
}
