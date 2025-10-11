# Charm Terminal Package

Beautiful, animated terminal output for Anvil CLI using [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Overview

This package enhances your CLI's visual output with:
- 🎨 Beautiful colors and styling
- ✨ Animated spinners for long operations
- 📊 Visual progress bars
- 🎁 Boxes, badges, and other UI components
- 🔄 Drop-in replacement for existing output

## Features

### 1. Enhanced Output Handler

Automatically enhances all your existing `palantir` output calls:

```go
o := palantir.GetGlobalOutputHandler()
o.PrintHeader("Installing Tools")    // ✨ Beautiful header with border
o.PrintSuccess("Installation done")  // ✓ Green with checkmark
o.PrintError("Failed to install")    // ✗ Red with X mark
o.PrintWarning("Skipping optional")  // ⚠ Yellow warning
o.PrintInfo("Processing files")      // ℹ Blue info
```

### 2. Animated Spinners

Show progress for long-running operations:

```go
spinner := charm.NewDotsSpinner("Installing package")
spinner.Start()
// ... do work ...
spinner.Success("Package installed!")
```

Available spinner types:
- **Dots**: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` (professional, default)
- **Line**: `|/-\` (minimal)
- **Circle**: `◜◠◝◞◡◟` (smooth)
- **Arrow**: `←↖↑↗→↘↓↙` (directional)
- **Box**: `◰◳◲◱` (playful)
- **Moon**: `🌑🌒🌓🌔🌕🌖🌗🌘` (creative)

### 3. Progress Bars

Visual progress indication:

```go
o.PrintProgress(3, 10, "Installing packages")
// Shows: [3/10] 30% ██████░░░░░░░░░░░░░░ Installing packages
```

### 4. Visual Components

#### Boxes
```go
content := charm.RenderBox("Configuration", "All settings valid", "#00FF87")
fmt.Println(content)
```

#### Lists
```go
items := []string{"git installed", "brew updated", "config valid"}
list := charm.RenderList(items, "✓", "#00FF87")
fmt.Println(list)
```

#### Key-Value Pairs
```go
fmt.Println(charm.RenderKeyValue("Version:", "2.0.0"))
fmt.Println(charm.RenderKeyValue("Status:", "Ready"))
```

#### Badges
```go
fmt.Println(charm.RenderBadge("READY", "#00FF87"))
fmt.Println(charm.RenderBadge("ERROR", "#FF5F87"))
```

#### Steps
```go
steps := []string{
    "Install Homebrew",
    "Configure git",
    "Install packages",
}
fmt.Println(charm.RenderSteps(steps))
```

#### Banners
```go
fmt.Println(charm.RenderBanner("ANVIL CLI"))
```

#### Code & Highlights
```go
fmt.Println(charm.RenderCode("brew install git"))
fmt.Println(charm.RenderHighlight("IMPORTANT", "#FFD700"))
```

## Quick Start

### 1. Initialize (Already done in main.go)

```go
import "github.com/rocajuanma/anvil/internal/terminal/charm"

func main() {
    charm.InitCharmOutput()
    // Your existing code...
}
```

### 2. Use Spinners in Long Operations

```go
func installPackage(name string) error {
    spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", name))
    spinner.Start()
    
    err := brew.InstallPackageDirectly(name)
    
    if err != nil {
        spinner.Error(fmt.Sprintf("Failed to install %s", name))
        return err
    }
    
    spinner.Success(fmt.Sprintf("%s installed!", name))
    return nil
}
```

### 3. Use Brew Wrapper (Convenience)

```go
brewSpinner := charm.NewBrewSpinner()
err := brewSpinner.InstallPackage("git", func() error {
    return brew.InstallPackageDirectly("git")
})
```

## Color Palette

Standard colors used throughout:

| Color | Hex | Usage |
|-------|-----|-------|
| Pink | `#FF6B9D` | Headers, important items |
| Green | `#00FF87` | Success messages |
| Red | `#FF5F87` | Error messages |
| Yellow | `#FFD700` | Warnings, highlights |
| Cyan | `#00D9FF` | Progress, stages |
| Blue | `#87CEEB` | Info, neutral text |
| Purple | `#C792EA` | Already available, code |
| Orange | `#FFA500` | Confirmations |
| Gray | `#666666` | Separators, subtle text |

## Best Practices

1. **Use spinners for operations > 1 second**
   - Package installations
   - Git operations
   - File downloads
   - Long computations

2. **Use progress bars for batch operations**
   - Installing multiple packages
   - Processing multiple files
   - Validating multiple items

3. **Use appropriate output types**
   - `PrintStage`: Major phase transitions
   - `PrintSuccess`: Successful completion
   - `PrintError`: Failures
   - `PrintWarning`: Non-critical issues
   - `PrintInfo`: General information

4. **Always stop spinners**
   ```go
   defer spinner.Stop()  // Or use Success/Error/Warning
   ```

5. **Use visual components for summaries**
   - Boxes for configuration details
   - Lists for results
   - Badges for status
   - Steps for instructions

## Migration from Plain Palantir

**Good news**: You don't need to change existing code! All your `palantir` calls are automatically enhanced once you call `charm.InitCharmOutput()` in main.

Then progressively add spinners and visual components where they provide value.

## Examples

See `examples.go` for comprehensive usage examples.

## Architecture

```
charm/
├── handler.go        # Enhanced output handler (wraps palantir)
├── spinner.go        # Animated spinner component
├── helpers.go        # Visual component helpers
├── brew_wrapper.go   # Convenience wrapper for brew ops
├── init.go          # Global initialization
├── examples.go      # Usage examples
└── README.md        # This file
```

## Dependencies

- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [rocajuanma/palantir](https://github.com/rocajuanma/palantir) - Base output handler

## Future Enhancements

Potential additions:
- Bubble Tea integration for interactive components
- Multi-line progress indicators
- Tree view rendering
- Table formatting
- Chart/graph rendering

