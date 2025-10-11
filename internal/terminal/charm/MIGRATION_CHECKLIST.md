# Charm Migration Checklist

Quick checklist for integrating Charm enhancements into your commands.

## ✅ Setup (Already Complete)

- [x] Added `charm.InitCharmOutput()` to `main.go`
- [x] Added lipgloss dependency to `go.mod`
- [x] Created charm package in `internal/terminal/charm/`
- [x] All tests passing

## 🎯 Integration Tasks

### High Priority Commands

#### [ ] `cmd/install/install.go`
- [ ] Add spinner to `installSingleTool()`
  ```go
  spinner := charm.NewDotsSpinner(fmt.Sprintf("Installing %s", toolName))
  spinner.Start()
  // ... install ...
  spinner.Success("Installed!")
  ```
- [ ] Enhance installation summary with badges/box
- [ ] Add visual progress for batch installations

#### [ ] `internal/brew/brew.go`
- [ ] Add spinner to `InstallPackageDirectly()`
- [ ] Add spinner to `InstallBrew()`
- [ ] Add spinner to `UpdateBrew()`
- [ ] Enhance package search with spinner

#### [ ] `cmd/initcmd/init.go`
- [ ] Add spinners during tool installation loop
- [ ] Use box for configuration summary
- [ ] Enhance group listing with styled list

### Medium Priority Commands

#### [ ] `cmd/config/pull/pull.go`
- [ ] Add spinner for git clone
- [ ] Add spinner for git pull
- [ ] Use box for summary

#### [ ] `cmd/config/push/push.go`
- [ ] Add spinner for git operations
- [ ] Use box for diff display
- [ ] Enhance PR link with highlight

#### [ ] `cmd/config/show/show.go`
- [ ] Use key-value pairs for config display
- [ ] Add boxes for sections
- [ ] Use list for groups

#### [ ] `cmd/doctor/doctor.go`
- [ ] Add spinner during validation
- [ ] Use status indicators for checks
- [ ] Use box for results summary
- [ ] Add badges for pass/fail counts

### Low Priority (Polish)

#### [ ] `cmd/clean/clean.go`
- [ ] Add spinner for cleanup operations
- [ ] Use box for summary

#### [ ] `cmd/config/import/import.go`
- [ ] Add spinner for download
- [ ] Use list for imported groups

#### [ ] `cmd/config/sync/sync.go`
- [ ] Add spinner for sync operation
- [ ] Use box for results

## 🔍 Pattern Recognition

### Find These Patterns:

1. **Long brew operations:**
   ```go
   // BEFORE
   system.RunCommand("brew", "install", pkg)
   
   // AFTER
   spinner := charm.NewDotsSpinner("Installing " + pkg)
   spinner.Start()
   system.RunCommand("brew", "install", pkg)
   spinner.Success("Installed!")
   ```

2. **Git operations:**
   ```go
   // BEFORE
   o.PrintStage("Cloning repository")
   git.Clone(repo)
   o.PrintSuccess("Repository cloned")
   
   // AFTER
   spinner := charm.NewDotsSpinner("Cloning repository")
   spinner.Start()
   git.Clone(repo)
   spinner.Success("Repository cloned")
   ```

3. **Summaries:**
   ```go
   // BEFORE
   fmt.Printf("Success: %d, Failed: %d\n", success, failed)
   
   // AFTER
   summary := fmt.Sprintf(
       "Success: %s | Failed: %s",
       charm.RenderBadge(fmt.Sprintf("%d", success), "#00FF87"),
       charm.RenderBadge(fmt.Sprintf("%d", failed), "#FF5F87"),
   )
   fmt.Println(charm.RenderBox("Results", summary, "#00D9FF"))
   ```

4. **Configuration displays:**
   ```go
   // BEFORE
   fmt.Printf("Version: %s\n", version)
   fmt.Printf("Branch: %s\n", branch)
   
   // AFTER
   fmt.Println(charm.RenderKeyValue("Version:", version))
   fmt.Println(charm.RenderKeyValue("Branch:", branch))
   ```

## 📝 Testing Checklist

After each integration:

- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] Spinner starts and stops correctly
- [ ] Colors display properly
- [ ] Visual components render correctly
- [ ] Error handling works
- [ ] No spinner overlaps
- [ ] Terminal degradation works (test with TERM=dumb)

## 🎨 Quick Reference

### Import Statement
```go
import "github.com/rocajuanma/anvil/internal/terminal/charm"
```

### Common Patterns
```go
// Spinner
spinner := charm.NewDotsSpinner("message")
spinner.Start()
// work...
spinner.Success("done")

// Box
charm.RenderBox("title", "content", "#00FF87")

// Badge
charm.RenderBadge("text", "#00FF87")

// Key-Value
charm.RenderKeyValue("key", "value")

// List
charm.RenderList(items, "•", "#87CEEB")
```

### Color Codes
- Success: `#00FF87`
- Error: `#FF5F87`
- Warning: `#FFD700`
- Info: `#00D9FF`
- Neutral: `#87CEEB`

## 🚦 Progress Tracking

Mark completed tasks with [x]:

**Week 1:**
- [ ] Install command enhanced
- [ ] Brew operations enhanced
- [ ] Init command enhanced

**Week 2:**
- [ ] Config commands enhanced
- [ ] Doctor command enhanced

**Week 3:**
- [ ] Remaining commands polished
- [ ] All visual components integrated

## 💡 Tips

1. Start with one command at a time
2. Test thoroughly before moving to next
3. Keep spinners simple (dots or line)
4. Use boxes sparingly (for summaries)
5. Always stop spinners before next output
6. Use appropriate colors for message types
7. Test on different terminal types

## 🐛 Common Issues & Solutions

**Spinner not animating?**
→ Operation too fast. Only use for >1 second operations.

**Text overlapping?**
→ Stop spinner before printing other output.

**Colors not showing?**
→ Terminal might not support. Check with `o.IsSupported()`.

**Spinner not stopping?**
→ Always defer `spinner.Stop()` or use Success/Error/Warning.

## 📚 Documentation

- `QUICK_REFERENCE.md` - Quick patterns
- `INTEGRATION_GUIDE.md` - Detailed guide
- `README.md` - Full API reference
- `examples.go` - Code examples

## ✨ Done!

When all tasks are complete:
- [ ] All commands enhanced
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Ready for release!

---

**Track your progress and enjoy building a beautiful CLI!** 🚀

