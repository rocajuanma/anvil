# Charm Integration Guide

## What Was Added

The Charm terminal package enhances your CLI with beautiful, modern UI components using [Lipgloss](https://github.com/charmbracelet/lipgloss).

### Files Created

```
internal/terminal/charm/
├── handler.go          # Enhanced output handler (replaces palantir visually)
├── spinner.go          # Animated spinner component
├── helpers.go          # Visual component utilities
├── brew_wrapper.go     # Convenience wrapper for brew operations
├── init.go            # Global initialization
├── demo.go            # Live demonstrations
├── examples.go        # Code examples and usage patterns
├── README.md          # Package documentation
└── INTEGRATION_GUIDE.md # This file
```

### What Changed

1. **main.go** - Added `charm.InitCharmOutput()` to enable enhanced output globally
2. **go.mod** - Added lipgloss and its dependencies

### What's Backward Compatible

✅ **Everything!** All your existing code works without changes. The enhancement is transparent:

```go
// Your existing code
o := palantir.GetGlobalOutputHandler()
o.PrintSuccess("Done!")

// Now automatically renders with:
// ✓ Done! (in beautiful green)
```

## How to Use

### 1. Basic Usage (Already Working)

Your existing `palantir` calls are automatically enhanced:

```go
o := palantir.GetGlobalOutputHandler()
o.PrintHeader("Installing Tools")    // ✨ Now with beautiful border
o.PrintSuccess("Installation done")  // ✓ Green with checkmark
o.PrintError("Failed")               // ✗ Red with X
o.PrintWarning("Optional skipped")   // ⚠ Yellow
o.PrintInfo("Processing")            // ℹ Blue
```

### 2. Add Spinners (Recommended for Long Operations)

#### Before (no visual feedback):
```go
func installPackage(name string) error {
    // Installation happens silently
    return brew.InstallPackageDirectly(name)
}
```

#### After (with beautiful spinner):
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

### 3. Integration Points

Look for these patterns in your code to add visual enhancements:

#### A. Brew Operations in `internal/brew/brew.go`

**Pattern to find:**
```go
// Long-running brew commands
system.RunCommand("brew", "install", packageName)
```

**Enhancement:**
```go
spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", packageName))
spinner.Start()
result, err := system.RunCommand("brew", "install", packageName)
if err != nil {
    spinner.Error("Installation failed")
    return err
}
spinner.Success("Installation complete")
```

#### B. Init Command in `cmd/initcmd/init.go`

**Pattern to find:**
```go
o.PrintStage("Installing required tools")
// Tool installation loop
```

**Enhancement:**
```go
o.PrintStage("Installing required tools")
for _, tool := range tools {
    spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", tool))
    spinner.Start()
    // ... install ...
    spinner.Success(fmt.Sprintf("%s ready", tool))
}
```

#### C. Install Command in `cmd/install/install.go`

**Current pattern:**
```go
func installSingleTool(toolName string) error {
    o := getOutputHandler()
    err := brew.InstallPackageDirectly(toolName)
    if err != nil {
        o.PrintError("Failed to install %s", toolName)
        return err
    }
    o.PrintSuccess("%s installed", toolName)
    return nil
}
```

**Enhanced:**
```go
func installSingleTool(toolName string) error {
    spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", toolName))
    spinner.Start()
    
    err := brew.InstallPackageDirectly(toolName)
    
    if err != nil {
        spinner.Error(fmt.Sprintf("Failed to install %s", toolName))
        return err
    }
    
    spinner.Success(fmt.Sprintf("%s installed successfully", toolName))
    return nil
}
```

#### D. Doctor Command in `cmd/doctor/doctor.go`

**Enhancement opportunity:**
```go
// When running validators
spinner := charm.NewCircleSpinner("Running health checks")
spinner.Start()
results := runValidators()
spinner.Success("Health check complete")

// Display results with visual components
fmt.Println(charm.RenderBox("Results", summary, "#00FF87"))
```

#### E. Config Commands in `cmd/config/`

**Pattern in pull/push:**
```go
o.PrintStage("Cloning repository")
// git clone
o.PrintSuccess("Repository cloned")
```

**Enhanced:**
```go
spinner := charm.NewDotsSpinner("Cloning repository")
spinner.Start()
// git clone
spinner.Success("Repository cloned")
```

### 4. Visual Components for Summaries

#### In Installation Summary:
```go
// Before
fmt.Printf("Installed: %d packages\n", count)
fmt.Printf("Failed: %d packages\n", failed)

// After
summary := fmt.Sprintf(
    "Installed: %s\nFailed: %s",
    charm.RenderBadge(fmt.Sprintf("%d", count), "#00FF87"),
    charm.RenderBadge(fmt.Sprintf("%d", failed), "#FF5F87"),
)
fmt.Println(charm.RenderBox("Installation Complete", summary, "#00D9FF"))
```

#### In Configuration Display:
```go
// Before
fmt.Printf("Version: %s\n", version)
fmt.Printf("Branch: %s\n", branch)

// After
fmt.Println(charm.RenderKeyValue("Version:", version))
fmt.Println(charm.RenderKeyValue("Branch:", branch))
```

## Recommended Integration Order

1. ✅ **Done**: Global initialization in `main.go`
2. 🎯 **High Impact**: Add spinners to `InstallPackageDirectly` in brew.go
3. 🎯 **High Impact**: Add spinners to installation commands
4. 📊 **Medium Impact**: Enhance progress indicators in batch operations
5. 🎨 **Polish**: Add visual components to summary outputs
6. 🎨 **Polish**: Enhance doctor command output with boxes and badges

## Testing Your Enhancements

### Quick Test
```bash
# Test enhanced output
./anvil --help

# Test installation (if safe)
./anvil install --dry-run git

# View demo
go run internal/terminal/charm/demo.go
```

### Gradual Rollout

1. Start with one command (e.g., `anvil init`)
2. Add spinners to long operations
3. Test thoroughly
4. Move to next command
5. Repeat

## Best Practices

### DO:
- ✅ Use spinners for operations > 1 second
- ✅ Always stop spinners before next output
- ✅ Use appropriate colors for message types
- ✅ Keep animations subtle and professional
- ✅ Test on different terminal types

### DON'T:
- ❌ Nest spinners (stop one before starting another)
- ❌ Use spinners for instant operations
- ❌ Mix spinner output with Printf statements
- ❌ Forget to stop spinners on errors
- ❌ Use overly flashy colors

## Troubleshooting

### Spinners not animating?
- Ensure operation takes > 80ms
- Check if running in non-TTY environment
- Verify terminal supports ANSI codes

### Colors not showing?
- Lipgloss auto-degrades for unsupported terminals
- Check `TERM` environment variable
- Run `o.IsSupported()` to verify

### Conflicts with existing output?
- Ensure spinners are stopped before other prints
- Use `\n` after spinner operations
- Clear line with `spinner.Stop()` before next output

## Examples in Action

See `examples.go` for comprehensive code examples and `demo.go` for live demonstrations.

Run the demo:
```go
import "github.com/rocajuanma/anvil/internal/terminal/charm"

charm.RunDemo()      // Full demo
charm.RunQuickDemo() // Quick test
```

## Future Enhancements

Potential additions:
- [ ] Bubble Tea for interactive components
- [ ] Progress bars with time estimates
- [ ] Tree view for dependency graphs
- [ ] Table formatting for structured data
- [ ] Chart rendering for statistics

## Support

If you have questions or need help:
1. Check `examples.go` for usage patterns
2. Review `README.md` for component details
3. Run `charm.RunDemo()` to see everything in action

