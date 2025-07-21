# Doctor Command

The `anvil doctor` command provides comprehensive health checks to validate your development environment and troubleshoot common configuration issues.

## Overview

The doctor command performs systematic validation across four key areas. You can run checks at different levels of granularity:

### 🏷️ **Categories** (groups of related checks)

- **environment** - Verify anvil initialization and directory structure (3 checks)
- **dependencies** - Check required tools and Homebrew installation (3 checks)
- **configuration** - Validate git and GitHub settings (3 checks)
- **connectivity** - Test GitHub access and repository connections (3 checks)

### 🔍 **Specific Checks** (individual validators)

Run `anvil doctor --list` to see all 12 available individual checks like `git-config`, `homebrew`, `github-access`, etc.

## Usage

### Basic Commands

```bash
# Run all health checks (12 total)
anvil doctor

# List available categories and checks with explanations
anvil doctor --list

# Run all checks in a category
anvil doctor environment        # 3 environment checks
anvil doctor dependencies       # 3 dependency checks
anvil doctor configuration      # 3 configuration checks
anvil doctor connectivity       # 3 connectivity checks

# Run a specific individual check
anvil doctor git-config
anvil doctor homebrew
anvil doctor github-access

# Show detailed output
anvil doctor --verbose
anvil doctor environment --verbose
```

### Auto-Fix Functionality

```bash
# Auto-fix all fixable issues
anvil doctor --fix

# Auto-fix issues in a specific category
anvil doctor dependencies --fix

# Auto-fix a specific check
anvil doctor homebrew --fix
```

## Understanding Categories vs Specific Checks

**Categories** are groups of related checks that test a particular area:

- When you run `anvil doctor environment`, it runs 3 checks: `init-run`, `settings-file`, and `directory-structure`
- When you run `anvil doctor dependencies`, it runs 3 checks: `homebrew`, `required-tools`, and `optional-tools`

**Specific checks** are individual validators that test one particular thing:

- When you run `anvil doctor git-config`, it only checks if your git configuration is properly set
- When you run `anvil doctor homebrew`, it only checks if Homebrew is installed and functional

Use categories when you want to check an entire area, and use specific checks when you want to focus on one particular aspect.

## Health Check Categories

### Environment Checks

| Check                 | Description                                     | Auto-Fix |
| --------------------- | ----------------------------------------------- | -------- |
| `init-run`            | Verify anvil initialization has been completed  | ❌       |
| `settings-file`       | Validate settings.yaml file exists and is valid | ✅       |
| `directory-structure` | Verify anvil directory structure is correct     | ✅       |

### Dependencies Checks

| Check            | Description                                 | Auto-Fix |
| ---------------- | ------------------------------------------- | -------- |
| `homebrew`       | Verify Homebrew is installed and functional | ✅       |
| `required-tools` | Check all required tools are installed      | ✅       |
| `optional-tools` | Check status of optional tools              | ❌       |

### Configuration Checks

| Check           | Description                                 | Auto-Fix |
| --------------- | ------------------------------------------- | -------- |
| `git-config`    | Verify git configuration is properly set    | ❌       |
| `github-config` | Verify GitHub configuration is properly set | ❌       |
| `sync-config`   | Verify sync configuration is valid          | ❌       |

### Connectivity Checks

| Check               | Description                                           | Auto-Fix |
| ------------------- | ----------------------------------------------------- | -------- |
| `github-access`     | Verify GitHub API access and authentication           | ❌       |
| `repository-access` | Verify configured repository exists and is accessible | ❌       |
| `git-connectivity`  | Verify git operations are functional                  | ❌       |

## Check Results

Each check returns one of four statuses:

- **✅ PASS** - Check completed successfully
- **⚠️ WARN** - Check passed but has warnings or recommendations
- **❌ FAIL** - Check failed and requires attention
- **⏭️ SKIP** - Check was skipped (usually due to missing configuration)

## Common Issues and Solutions

### Environment Issues

**Settings file not found**

```bash
# Solution: Run anvil init
anvil init
```

**Directory structure incomplete**

```bash
# Solution: Auto-fix creates missing directories
anvil doctor directory-structure --fix
```

### Dependency Issues

**Homebrew not installed**

```bash
# Solution: Auto-fix installs Homebrew
anvil doctor homebrew --fix
```

**Required tools missing**

```bash
# Solution: Auto-fix installs missing tools
anvil doctor required-tools --fix
```

### Configuration Issues

**Git configuration incomplete**

```bash
# Solution: Manually configure in settings.yaml
# Edit ~/.anvil/settings.yaml:
git:
  username: "Your Name"
  email: "your.email@example.com"
```

**GitHub configuration incomplete**

```bash
# Solution: Configure GitHub repository in settings.yaml
# Edit ~/.anvil/settings.yaml:
github:
  config_repo: "username/repository"
  branch: "main"
  token_env_var: "GITHUB_TOKEN"
```

### Connectivity Issues

**GitHub authentication failed**

```bash
# Solution: Set up GitHub token or SSH keys
export GITHUB_TOKEN="your_token_here"

# Or configure SSH keys:
ssh-keygen -t ed25519 -C "your.email@example.com"
# Add to GitHub: https://github.com/settings/keys
```

**Repository not accessible**

```bash
# Solution: Check repository name and permissions
# Ensure repository exists and you have access
# Update config_repo in settings.yaml if needed
```

## Output Examples

### Full Health Check

```bash
$ anvil doctor

=== Running Anvil Health Check ===

🔍 Validating environment, dependencies, configuration, and connectivity...

✅ Environment
  ✅ Anvil initialization complete
  ✅ Settings file is valid
  ✅ Directory structure is correct

⚠️ Dependencies
  ✅ Homebrew is installed and functional
  ✅ All required tools installed (3/3)
  ⚠️ Optional tools missing: docker
      💡 Run 'anvil install docker' to install missing tools

✅ Configuration
  ✅ Git configuration is complete
  ✅ GitHub configuration is complete
  ⏭️ No sync configuration present (optional)

✅ Connectivity
  ✅ GitHub API access confirmed
  ✅ Repository is accessible
  ✅ Git operations are functional

=== Health Check Summary ===

Total checks: 12
✅ Passed: 10
⚠️ Warnings: 1
⏭️ Skipped: 1

✅ Overall status: Healthy
```

### Category-Specific Check

```bash
$ anvil doctor dependencies

=== Running Dependencies Health Checks ===

⚠️ Dependencies
  ✅ Homebrew is installed and functional
  ❌ Missing required tools: git, curl
      💡 Missing tools will be installed automatically
  ✅ All optional tools installed (2/2)

=== Health Check Summary ===

Total checks: 3
✅ Passed: 2
❌ Failed: 1

🔧 1 issues can be auto-fixed
Run 'anvil doctor --fix' to automatically fix them
❌ Overall status: Issues found
```

## Integration with Other Commands

The doctor command integrates seamlessly with other anvil commands:

```bash
# After init, verify setup
anvil init
anvil doctor

# Before config operations, check connectivity
anvil doctor connectivity
anvil config pull

# After installation, verify dependencies
anvil install dev
anvil doctor dependencies
```

## Best Practices

1. **Run after initialization**: Always run `anvil doctor` after `anvil init` to verify setup
2. **Pre-flight checks**: Use category-specific checks before major operations
3. **Regular maintenance**: Run periodic health checks to catch configuration drift
4. **Auto-fix safely**: Review auto-fix suggestions before applying them
5. **Verbose output**: Use `--verbose` for detailed troubleshooting information

## Troubleshooting

If the doctor command itself fails:

```bash
# Check basic anvil functionality
anvil --help

# Verify Go installation
go version

# Rebuild anvil
go build -o anvil main.go
```

For issues with specific validators, check the relevant package documentation and ensure all dependencies are properly installed.
