# Charm Terminal - Quick Reference Card

## 🎯 Most Common Use Cases

### 1. Add Spinner to Long Operation
```go
import "github.com/rocajuanma/anvil/internal/terminal/charm"

spinner := charm.NewDotsSpinner("Installing package")
spinner.Start()
// ... do work ...
spinner.Success("Package installed!")  // or .Error() or .Warning()
```

### 2. Enhanced Output (Automatic)
```go
o := palantir.GetGlobalOutputHandler()
o.PrintSuccess("Done")   // ✓ Done (green)
o.PrintError("Failed")   // ✗ Failed (red)
o.PrintWarning("Skip")   // ⚠ Skip (yellow)
o.PrintInfo("Loading")   // ℹ Loading (blue)
```

### 3. Show Box with Content
```go
fmt.Println(charm.RenderBox("Title", "Content here", "#00FF87"))
```

### 4. Display List
```go
items := []string{"Item 1", "Item 2", "Item 3"}
fmt.Println(charm.RenderList(items, "•", "#87CEEB"))
```

### 5. Key-Value Pairs
```go
fmt.Println(charm.RenderKeyValue("Version:", "2.0.0"))
fmt.Println(charm.RenderKeyValue("Status:", "Ready"))
```

## 🎨 Spinner Types

```go
NewDotsSpinner(msg)    // ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ (default, professional)
NewLineSpinner(msg)    // |/-\ (minimal)
NewCircleSpinner(msg)  // ◜◠◝◞◡◟ (smooth)
```

## 🌈 Standard Colors

```go
"#FF6B9D"  // Pink - Headers, titles
"#00FF87"  // Green - Success
"#FF5F87"  // Red - Errors
"#FFD700"  // Yellow - Warnings
"#00D9FF"  // Cyan - Progress
"#87CEEB"  // Blue - Info
"#C792EA"  // Purple - Code, available
```

## 📦 Visual Components

```go
// Banner
RenderBanner("ANVIL CLI")

// Badge
RenderBadge("SUCCESS", "#00FF87")

// Steps
steps := []string{"Step 1", "Step 2"}
RenderSteps(steps)

// Code
RenderCode("brew install git")

// Highlight
RenderHighlight("IMPORTANT", "#FFD700")

// Separator
RenderSeparator(50, "─", "#666666")

// Status
RenderStatus("All good", true)  // Green dot
RenderStatus("Error", false)    // Red dot
```

## 🔧 Brew Wrapper (Convenience)

```go
brewSpinner := charm.NewBrewSpinner()

brewSpinner.InstallPackage("git", func() error {
    return brew.InstallPackageDirectly("git")
})

brewSpinner.UpdateBrew(func() error {
    return brew.Update()
})
```

## ⚡ Quick Integration Pattern

### Before:
```go
func install(pkg string) error {
    return brew.InstallPackageDirectly(pkg)
}
```

### After:
```go
func install(pkg string) error {
    spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", pkg))
    spinner.Start()
    defer spinner.Stop()  // Safety
    
    err := brew.InstallPackageDirectly(pkg)
    
    if err != nil {
        spinner.Error("Installation failed")
        return err
    }
    
    spinner.Success("Installation complete")
    return nil
}
```

## 💡 Best Practices

✅ DO:
- Use spinners for operations > 1 second
- Always stop spinners (Success/Error/Warning)
- Use appropriate colors for message types
- Keep animations professional

❌ DON'T:
- Nest spinners (stop before starting another)
- Use spinners for instant operations
- Mix spinner output with raw Printf
- Forget to handle spinner on error

## 🐛 Common Issues

**Spinner not animating?**
```go
// Ensure operation takes time
time.Sleep(2 * time.Second)  // For testing
```

**Want to update spinner message?**
```go
spinner.UpdateMessage("New message")
```

**Need custom style?**
```go
spinner.WithColor("#FF6B9D")
spinner.WithStyle(customLipglossStyle)
```

## 📚 Full Docs

- `README.md` - Complete package documentation
- `INTEGRATION_GUIDE.md` - Step-by-step integration
- `examples.go` - Code examples
- `demo.go` - Run `charm.RunDemo()` for live demo

## 🎬 See It In Action

```go
import "github.com/rocajuanma/anvil/internal/terminal/charm"

// Quick test
charm.RunQuickDemo()

// Full demo
charm.RunDemo()
```

---

**Remember:** All existing `palantir` calls are automatically enhanced.
Just add spinners and visual components where they add value! 🚀

