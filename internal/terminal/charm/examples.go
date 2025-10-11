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

/*
USAGE EXAMPLES:

1. BASIC OUTPUT (Enhanced versions of your existing palantir calls):

   o := palantir.GetGlobalOutputHandler()
   o.PrintHeader("Starting Installation")  // Now with beautiful border and colors!
   o.PrintSuccess("Package installed")     // Green with checkmark
   o.PrintError("Installation failed")     // Red with X
   o.PrintWarning("Skipping optional step") // Yellow with warning
   o.PrintInfo("Downloading dependencies")  // Blue with info icon

2. SPINNERS (For long-running operations):

   // Simple spinner
   spinner := charm.NewDotsSpinner("Installing package")
   spinner.Start()
   // ... do work ...
   spinner.Success("Package installed!")

   // Or use the brew wrapper
   brewSpinner := charm.NewBrewSpinner()
   brewSpinner.InstallPackage("git", func() error {
       return brew.InstallPackageDirectly("git")
   })

3. PROGRESS BARS (Already enhanced in PrintProgress):

   o.PrintProgress(5, 10, "Installing packages")
   // Now shows: [5/10] 50% ████████████░░░░░░░░ Installing packages

4. VISUAL COMPONENTS:

   // Box with content
   content := charm.RenderBox("Configuration", "All settings are valid", "#00FF87")
   fmt.Println(content)

   // List of items
   items := []string{"git installed", "brew updated", "config valid"}
   list := charm.RenderList(items, "✓", "#00FF87")
   fmt.Println(list)

   // Key-value pairs
   fmt.Println(charm.RenderKeyValue("Version:", "2.0.0"))
   fmt.Println(charm.RenderKeyValue("Status:", "Ready"))

   // Banner
   fmt.Println(charm.RenderBanner("ANVIL CLI"))

   // Status badge
   fmt.Println(charm.RenderBadge("READY", "#00FF87"))
   fmt.Println(charm.RenderBadge("ERROR", "#FF5F87"))

   // Steps
   steps := []string{
       "Install Homebrew",
       "Configure git",
       "Install packages",
   }
   fmt.Println(charm.RenderSteps(steps))

   // Separator
   fmt.Println(charm.RenderSeparator(50, "─", "#666666"))

5. CODE FORMATTING:

   // Highlight important text
   fmt.Println(charm.RenderHighlight("IMPORTANT", "#FFD700"))

   // Show code
   fmt.Println(charm.RenderCode("brew install git"))

   // Quote
   quote := charm.RenderQuote(
       "Make it work, make it right, make it fast",
       "Kent Beck",
   )
   fmt.Println(quote)

6. INTEGRATION IN YOUR COMMANDS:

   In cmd/install/install.go:

   func installSingleTool(toolName string) error {
       o := getOutputHandler()

       // Show what we're doing
       o.PrintStage(fmt.Sprintf("Installing %s", toolName))

       // Use spinner for the actual installation
       spinner := charm.NewDotsSpinner(fmt.Sprintf("Downloading %s", toolName))
       spinner.Start()

       err := brew.InstallPackageDirectly(toolName)

       if err != nil {
           spinner.Error(fmt.Sprintf("Failed to install %s", toolName))
           return err
       }

       spinner.Success(fmt.Sprintf("%s installed successfully", toolName))
       return nil
   }

7. DIFFERENT SPINNER STYLES:

   // Choose based on operation type
   spinner := charm.NewDotsSpinner("Processing...")    // Default, professional
   spinner := charm.NewLineSpinner("Loading...")       // Minimal
   spinner := charm.NewCircleSpinner("Searching...")   // Smooth arc

   // Custom frames
   customFrames := []string{"◐", "◓", "◑", "◒"}
   spinner := charm.NewSpinner(customFrames, "Custom animation")

8. CONDITIONAL STYLING:

   // Only use enhanced output if supported
   if o.IsSupported() {
       fmt.Println(charm.RenderBanner("Welcome to Anvil"))
   } else {
       fmt.Println("Welcome to Anvil")
   }

BEST PRACTICES:

1. Use spinners for operations > 1 second
2. Use PrintProgress for batch operations with known count
3. Use PrintStage to mark major phases
4. Use appropriate colors:
   - Green (#00FF87) for success
   - Red (#FF5F87) for errors
   - Yellow (#FFD700) for warnings
   - Blue (#87CEEB) for info
   - Cyan (#00D9FF) for progress

5. Always stop spinners before printing other output
6. Use boxes for important summaries
7. Use badges for status indicators
8. Keep animations subtle and professional

MIGRATION GUIDE:

The beauty is you don't need to change anything! Just by initializing
charm.InitCharmOutput() in main.go, all your existing palantir calls
automatically get enhanced styling. Then progressively add spinners
and other visual components where they make sense.

*/
