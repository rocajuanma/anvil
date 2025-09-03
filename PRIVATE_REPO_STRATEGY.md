# Private Repository Distribution Strategy

This guide explains how to distribute Anvil binaries publicly while keeping your source code private.

## ✅ **What Works with Private Repos**

### **GitHub Releases (Fully Supported)**
- ✅ Release binaries are **publicly downloadable**
- ✅ Release notes are **publicly visible**
- ✅ GitHub Actions can build and publish releases
- ✅ Install script works from release attachments
- ✅ Direct binary downloads work perfectly

### **Installation Methods That Work**
```bash
# ✅ Download install script from releases
curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash

# ✅ Direct binary download
curl -L https://github.com/rocajuanma/anvil/releases/latest/download/anvil-darwin-universal -o anvil

# ✅ Go install (if you want to allow source access)
go install github.com/rocajuanma/anvil@latest
```

## ❌ **Limitations with Private Repos**

### **Homebrew Traditional Tap**
- ❌ Standard Homebrew taps require **public source access**
- ✅ **Alternative**: Binary-only formula (downloads pre-built binary)

### **Raw File Access**
- ❌ `raw.githubusercontent.com` URLs don't work for private repos
- ✅ **Solution**: Attach files to releases instead

## 🛠 **Recommended Setup for Private Repos**

### **1. Use Release-Based Distribution**

Your GitHub Actions workflow will:
1. Build binaries privately
2. Attach binaries to public releases
3. Include `install.sh` in release attachments

Users install with:
```bash
curl -sSL https://github.com/rocajuanma/anvil/releases/latest/download/install.sh | bash
```

### **2. Documentation Strategy**

Since users can't see your README directly, provide documentation via:

**Option A: Rich Release Notes**
Include full installation and usage instructions in every release.

**Option B: Public Documentation Repository**
Create `rocajuanma/anvil-docs` (public) with:
- Installation guides
- Usage examples  
- Configuration documentation
- Link from release notes

**Option C: GitHub Pages or Website**
Host documentation at `rocajuanma.github.io/anvil` or your domain.

### **3. Binary-Only Homebrew Formula**

Create a formula that downloads pre-built binaries instead of building from source:

```ruby
class Anvil < Formula
  desc "Complete macOS development environment automation tool"
  homepage "https://github.com/rocajuanma/anvil"
  
  if Hardware::CPU.intel?
    url "https://github.com/rocajuanma/anvil/releases/latest/download/anvil-darwin-amd64"
    sha256 "INTEL_SHA256_HERE"
  else
    url "https://github.com/rocajuanma/anvil/releases/latest/download/anvil-darwin-arm64"  
    sha256 "ARM_SHA256_HERE"
  end
  
  license "Apache-2.0"

  def install
    bin.install "anvil-darwin-amd64" => "anvil" if Hardware::CPU.intel?
    bin.install "anvil-darwin-arm64" => "anvil" if Hardware::CPU.arm?
  end

  test do
    assert_match "Anvil", shell_output("#{bin}/anvil --help")
  end
end
```

## 🔒 **Security Considerations**

### **What Remains Private**
- ✅ Source code and implementation details
- ✅ Commit history and development process
- ✅ Private dependencies or configurations
- ✅ Team discussions in issues/PRs

### **What Becomes Public**
- ⚠️ Binary files (but not source code)
- ⚠️ Release notes and version information
- ⚠️ Repository name and basic metadata

### **Best Practices**
1. **Review release notes** - Don't include sensitive information
2. **Use meaningful commit messages** for tags (they're visible in releases)
3. **Consider binary obfuscation** if needed
4. **Monitor downloads** via GitHub insights

## 📊 **Example: Successful Private + Public Distribution**

Many commercial tools use this approach:

- **Docker Desktop** - Private source, public binaries
- **1Password CLI** - Private source, public releases
- **Tailscale** - Private source, public binaries

Your setup will work exactly the same way.

## 🚀 **Next Steps**

1. **Keep your repo private**
2. **Set up the GitHub Actions workflow** (already created)
3. **Create your first release** to test the process
4. **Choose a documentation strategy** (release notes vs. separate repo)
5. **Optionally set up binary-only Homebrew formula**

The GitHub Releases approach is **perfect** for your needs - you get professional binary distribution while keeping your source code completely private.
