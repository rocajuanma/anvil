# Homebrew Tap Setup Guide

This guide explains how to create a Homebrew tap for Anvil, allowing users to install with `brew install rocajuanma/anvil/anvil`.

> **Note**: This setup is for the public repository. Since Anvil is open source, the standard Homebrew formula will work directly.

## 🚀 Quick Setup Steps

### 1. Create the Tap Repository

Create a new **public** repository on GitHub named `homebrew-anvil`:
- Repository name: `homebrew-anvil`
- Description: "Homebrew tap for Anvil - macOS development environment automation"
- Make it **public** (required for Homebrew taps)

### 2. Create the Formula

In your `homebrew-anvil` repository, create the following structure:
```
homebrew-anvil/
└── Formula/
    └── anvil.rb
```

Create `Formula/anvil.rb` with the formula content (see example below).

### 3. Create the Formula

Create `Formula/anvil.rb` with this content:

```ruby
class Anvil < Formula
  desc "Complete macOS development environment automation tool"
  homepage "https://github.com/rocajuanma/anvil"
  url "https://github.com/rocajuanma/anvil/archive/v1.1.1.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256_OF_TARBALL"
  license "Apache-2.0"
  head "https://github.com/rocajuanma/anvil.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./main.go"
    
    # Generate shell completions (optional but recommended)
    generate_completions_from_executable(bin/"anvil", "completion")
  end

  def caveats
    <<~EOS
      🔨 Anvil has been installed!
      
      Get started:
        anvil init          # Initialize your environment
        anvil doctor        # Verify setup and dependencies
        anvil install dev   # Install development tools
      
      📚 Documentation: https://github.com/rocajuanma/anvil/docs
    EOS
  end

  test do
    assert_match "Anvil", shell_output("#{bin}/anvil --help")
    assert_match version.to_s, shell_output("#{bin}/anvil --version")
  end
end
```

### 4. Update the Formula for New Releases

For each new release, update:
```ruby
url "https://github.com/rocajuanma/anvil/archive/vX.Y.Z.tar.gz"
sha256 "NEW_SHA256_HASH"  # Get with: curl -sL URL | shasum -a 256
```

### 4. Test the Formula Locally

```bash
# Install your tap locally
brew tap rocajuanma/anvil

# Test installation
brew install --build-from-source anvil

# Test the binary
anvil --version
anvil doctor

# Uninstall for testing
brew uninstall anvil
brew untap rocajuanma/anvil
```

### 5. Automate Formula Updates

For automatic updates when you release new versions, add this GitHub Action to your **main anvil repository**:

**.github/workflows/update-homebrew.yml:**
```yaml
name: Update Homebrew Formula

on:
  release:
    types: [published]

jobs:
  update-homebrew:
    runs-on: ubuntu-latest
    steps:
    - name: Update Homebrew formula
      uses: mislav/bump-homebrew-formula@v2
      with:
        formula-name: anvil
        formula-path: Formula/anvil.rb
        homebrew-tap: rocajuanma/homebrew-anvil
        download-url: https://github.com/rocajuanma/anvil/archive/${{ github.event.release.tag_name }}.tar.gz
        commit-message: |
          anvil ${{ github.event.release.tag_name }}
          
          Created by ${{ github.event.release.html_url }}
      env:
        COMMITTER_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

### 6. Create GitHub Token

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Create a token with `public_repo` permissions
3. Add it as a secret named `HOMEBREW_TAP_TOKEN` in your main anvil repository

## 📋 User Installation Flow

Once set up, users can install Anvil with:

```bash
# Add your tap
brew tap rocajuanma/anvil

# Install anvil
brew install anvil

# Or do both in one command
brew install rocajuanma/anvil/anvil
```

## 🔄 Maintenance

### Updating for New Releases

When you create a new release (using the GitHub Actions workflow), the Homebrew formula will automatically update. If you need to update manually:

1. Update the `url` and `version` in `Formula/anvil.rb`
2. Get new SHA256: `curl -sL NEW_URL | shasum -a 256`
3. Update the `sha256` field
4. Commit and push to `homebrew-anvil` repository

### Formula Testing

Always test your formula changes:

```bash
# Audit the formula
brew audit --strict anvil

# Test installation
brew install --build-from-source anvil

# Test functionality
anvil --version
anvil doctor
```

## 🎯 Benefits for Users

- **Simple installation**: `brew install rocajuanma/anvil/anvil`
- **Automatic updates**: `brew upgrade` updates Anvil
- **Dependency management**: Homebrew handles Go installation if needed
- **Uninstallation**: `brew uninstall anvil`
- **Integration**: Works with existing Homebrew workflows

## 📚 Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Creating Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Formula Reference](https://rubydoc.brew.sh/Formula)
