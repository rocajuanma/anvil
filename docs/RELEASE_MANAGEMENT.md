# Release Management Guide

This guide covers how to create releases for Anvil, manage GitHub Actions workflows, and troubleshoot common issues.

## 🚀 **Creating a New Release**

### **Standard Release Process**

1. **Ensure your changes are merged to master:**
   ```bash
   git checkout master
   git pull origin master
   git status  # Should show "working tree clean"
   ```

2. **Create and push a version tag:**
   ```bash
   # Create a new tag (increment version appropriately)
   git tag v1.2.0
   
   # Push the tag to trigger GitHub Actions
   git push origin v1.2.0
   ```

3. **Monitor the GitHub Actions workflow:**
   - Go to: `https://github.com/0xjuanma/anvil/actions`
   - Watch the "Release" workflow complete
   - Check for any build errors

4. **Verify the release was created:**
   - Go to: `https://github.com/0xjuanma/anvil/releases`
   - Confirm binaries are attached
   - Test download links

### **Version Numbering Strategy**

Follow [Semantic Versioning](https://semver.org/):

- **Major (v2.0.0)**: Breaking changes, major new features
- **Minor (v1.2.0)**: New features, backward compatible
- **Patch (v1.1.1)**: Bug fixes, small improvements

**Examples:**
```bash
git tag v1.1.1   # Bug fix release
git tag v1.2.0   # New feature release  
git tag v2.0.0   # Major release with breaking changes
```

### **Pre-release and Beta Versions**

For testing purposes:
```bash
git tag v1.2.0-beta.1    # Beta release
git tag v1.2.0-rc.1      # Release candidate
git tag v1.2.0-alpha.1   # Alpha release
```

## 🔄 **Managing Release Tags**

### **Fixing a Failed Release**

If a release fails or needs to be updated:

1. **Delete the tag locally and remotely:**
   ```bash
   git tag -d v1.1.2
   git push origin --delete v1.1.2
   ```

2. **Delete the GitHub release (if created):**
   - Go to: `https://github.com/0xjuanma/anvil/releases`
   - Click on the problematic release
   - Click "Delete" button

3. **Fix the issue, commit, and recreate the tag:**
   ```bash
   # Fix the issue and commit
   git add .
   git commit -m "fix: resolve release issue"
   git push origin master
   
   # Recreate the tag
   git tag v1.1.2
   git push origin v1.1.2
   ```

### **Listing and Managing Tags**

```bash
# List all tags
git tag -l

# List tags with pattern
git tag -l "v1.1.*"

# Show tag details
git show v1.1.2

# Delete local tag
git tag -d v1.1.2

# Delete remote tag
git push origin --delete v1.1.2
```

## ⚙️ **GitHub Actions Workflow Details**

### **Workflow Triggers**

The release workflow (`.github/workflows/release.yml`) triggers on:

- **Tag push**: Any tag matching `v*.*.*` pattern
- **Manual trigger**: Via GitHub Actions UI

### **What the Workflow Does**

1. **Builds binaries for multiple platforms:**
   - macOS Intel (`darwin-amd64`)
   - macOS Apple Silicon (`darwin-arm64`)
   - Linux Intel (`linux-amd64`)
   - Linux ARM (`linux-arm64`)

2. **Generates security checksums** for all binaries

3. **Creates a GitHub release** with:
   - Release notes with installation instructions
   - All binary attachments
   - Checksums file
   - Installation script

### **Workflow Environment Variables**

The workflow uses these built-in variables:
- `$GITHUB_REF`: Contains the tag name
- `$GITHUB_TOKEN`: Automatic authentication
- `${{ github.event.release.tag_name }}`: Tag from release event

### **Customizing Build Flags**

Current build command:
```yaml
go build -ldflags="-s -w -X main.version=${{ steps.version.outputs.VERSION }}" -o dist/anvil-darwin-amd64 main.go
```

**Build flags explained:**
- `-s`: Strip symbol table
- `-w`: Strip debug info
- `-X main.version=...`: Set version variable

**To add more build info:**
```yaml
-ldflags="-s -w -X main.version=${{ steps.version.outputs.VERSION }} -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=${{ github.sha }}"
```

### **Adding New Platforms**

To add Windows support:
```yaml
# Add to build step
- name: Build binaries
  run: |
    # ... existing builds ...
    
    # Build for Windows
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${{ steps.version.outputs.VERSION }}" -o dist/anvil-windows-amd64.exe main.go

# Add to files list
files: |
  dist/anvil-darwin-amd64
  dist/anvil-darwin-arm64
  dist/anvil-linux-amd64
  dist/anvil-linux-arm64
  dist/anvil-windows-amd64.exe
  dist/checksums.txt
  install.sh
```

## 🧪 **Testing Releases**

### **Test Installation Locally**

Before creating a public release, test the install process:

```bash
# Test install script (after release is published)
curl -sSL https://github.com/0xjuanma/anvil/releases/latest/download/install.sh | bash

# Test manual download
curl -L https://github.com/0xjuanma/anvil/releases/latest/download/anvil-darwin-amd64 -o anvil-test
chmod +x anvil-test
./anvil-test --version
```

### **Verify Checksums**

```bash
# Download and verify
curl -L https://github.com/0xjuanma/anvil/releases/latest/download/anvil-darwin-amd64 -o anvil
curl -L https://github.com/0xjuanma/anvil/releases/latest/download/checksums.txt -o checksums.txt

# Verify checksum (macOS/Linux)
shasum -a 256 anvil
grep anvil-darwin-amd64 checksums.txt
```

### **Cross-Platform Testing**

Test on different platforms:
- **macOS Intel**: Download `anvil-darwin-amd64`
- **macOS Apple Silicon**: Download `anvil-darwin-arm64`  
- **Linux**: Download `anvil-linux-amd64`

## 🔧 **Troubleshooting Common Issues**

### **Workflow Fails to Run**

**Issue**: Tag push doesn't trigger workflow
**Solution**: 
- Ensure tag follows `v*.*.*` pattern
- Check workflow file is in `.github/workflows/`
- Verify GitHub Actions are enabled for the repo

### **Build Failures**

**Issue**: Go build fails
**Solution**:
- Check Go version in workflow (currently 1.17)
- Verify all dependencies are in `go.mod`
- Test build locally: `go build -o anvil main.go`

### **Permission Errors**

**Issue**: Can't create release or upload assets
**Solution**:
- GitHub Actions uses `GITHUB_TOKEN` automatically
- Ensure repository has "Actions" permissions enabled
- For private repos, check token permissions

### **Binary Download Issues**

**Issue**: Install script fails to download
**Solution**:
- Verify release exists and is published (not draft)
- Check binary names match in workflow and install script
- Ensure release is public (even for private repos)

### **Architecture Detection Problems**

**Issue**: Install script downloads wrong binary
**Solution**:
```bash
# Test architecture detection
uname -m    # Should return x86_64 or arm64
uname -s    # Should return Darwin or Linux

# Manually specify in install script if needed
```

## 📋 **Release Checklist**

Use this checklist for each release:

- [ ] **Pre-release:**
  - [ ] All changes merged to master
  - [ ] Version number decided (follow semver)
  - [ ] Local testing completed
  - [ ] Documentation updated

- [ ] **Release:**
  - [ ] Create and push tag: `git tag vX.Y.Z && git push origin vX.Y.Z`
  - [ ] Monitor GitHub Actions workflow
  - [ ] Verify release created successfully

- [ ] **Post-release:**
  - [ ] Test installation methods
  - [ ] Verify binary downloads work
  - [ ] Update any external documentation
  - [ ] Announce release (if applicable)

## 🎯 **Advanced Tips**

### **Automated Changelog Generation**

Add to workflow for automatic changelogs:
```yaml
- name: Generate changelog
  run: |
    git log $(git describe --tags --abbrev=0 HEAD^)..HEAD --pretty=format:"* %s" > CHANGELOG_LATEST.md
```

### **Draft Releases**

For review before publishing:
```yaml
draft: true  # Change to false when ready
```

### **Conditional Builds**

Skip builds for documentation-only changes:
```yaml
if: "!contains(github.event.head_commit.message, '[skip ci]')"
```

### **Parallel Builds**

Speed up builds with matrix strategy:
```yaml
strategy:
  matrix:
    goos: [darwin, linux]
    goarch: [amd64, arm64]
```

## 📚 **Related Documentation**

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)

---

**💡 Pro Tip**: Keep this document updated as you add new platforms, change build processes, or learn new techniques!
