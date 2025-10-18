# Doctor Command

The `anvil doctor` command provides comprehensive health checks to validate your development environment and troubleshoot common configuration issues.

## Overview

The doctor command performs systematic validation across four key areas. You can run checks at different levels of granularity:

### Categories (groups of related checks)

- **environment** - Verify anvil initialization and directory structure (3 checks)
- **dependencies** - Check required tools and Homebrew installation (2 checks)
- **configuration** - Validate git and GitHub settings (3 checks)
- **connectivity** - Test GitHub access and repository connections (3 checks)

### Specific Checks (individual validators)

Run `anvil doctor --list` to see all available individual checks like `git-config`, `homebrew`, `github-auth`, etc.

## Key Features

### Real-time Progress

- Live feedback as each validation runs
- Progress counters showing remaining checks
- Stage indicators for current phase
- No hanging terminals

### Two Output Modes

- **Default mode** - Brief but informative progress messages
- **Verbose mode (`--verbose`)** - Detailed descriptions and step-by-step results

### Secure Authentication

- Non-interactive operations - no credential prompts
- Environment-based auth using `GITHUB_TOKEN` or SSH keys
- Private repository enforcement with clear warnings

## Usage

### Basic Commands

```bash
# Run all health checks (11 total) with real-time progress
anvil doctor

# List available categories and checks with explanations
anvil doctor --list

# Run all checks in a category with progress feedback
anvil doctor environment        # 3 environment checks
anvil doctor dependencies       # 2 dependency checks
anvil doctor configuration      # 3 configuration checks
anvil doctor connectivity       # 3 connectivity checks

# Run a specific individual check with detailed feedback
anvil doctor git-config
anvil doctor homebrew
anvil doctor github-auth

# Show detailed output with descriptions and step-by-step results
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

- When you run `anvil doctor environment`, it runs 3 checks: `anvil-init`, `settings-valid`, and `directory-structure`
- When you run `anvil doctor dependencies`, it runs 2 checks: `homebrew` and `required-tools`

**Specific checks** are individual validators that test one particular thing:

- When you run `anvil doctor git-config`, it only checks if your git configuration is properly set
- When you run `anvil doctor homebrew`, it only checks if Homebrew is installed and functional

Use categories when you want to check an entire area, and use specific checks when you want to focus on one particular aspect.

## Health Check Categories

### Environment Checks

| Check                 | Description                                     | Auto-Fix |
| --------------------- | ----------------------------------------------- | -------- |
| `anvil-init`          | Verify anvil initialization has been completed  | No       |
| `settings-valid`      | Validate settings.yaml structure and content    | No       |
| `directory-structure` | Check ~/.anvil directory structure              | No       |

### Dependencies Checks

| Check            | Description                                         | Auto-Fix |
| ---------------- | --------------------------------------------------- | -------- |
| `homebrew`       | Verify Homebrew installation and updates            | Yes      |
| `required-tools` | Check git and curl are installed                    | No       |

### Configuration Checks

| Check           | Description                                  | Auto-Fix |
| --------------- | -------------------------------------------- | -------- |
| `git-config`    | Validate git user.name and user.email        | Yes      |
| `github-config` | Verify GitHub repository configuration       | No       |
| `sync-config`   | Check config sync settings (not implemented) | No       |

### Connectivity Checks

| Check             | Description                              | Auto-Fix |
| ----------------- | ---------------------------------------- | -------- |
| `github-auth`     | Test GitHub authentication and access    | No       |
| `github-repo`     | Verify repository accessibility          | No       |
| `git-operations`  | Test git clone and pull operations       | No       |

## Check Results

Each check returns one of four statuses:

- **PASS** - Check completed successfully
- **WARN** - Check passed but has warnings or recommendations
- **FAIL** - Check failed and requires attention
- **SKIP** - Check was skipped (usually due to missing configuration)

## Understanding Output

The doctor command provides real-time progress feedback and organized results by category. Use `--verbose` for detailed information about each check.

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
2. **Use real-time feedback**: Watch the progress indicators to understand timing
3. **Pre-flight checks**: Use category-specific checks before major operations
4. **Regular maintenance**: Run periodic health checks to catch configuration drift
5. **Auto-fix safely**: Review auto-fix suggestions before applying them
6. **Verbose for debugging**: Use `--verbose` for detailed troubleshooting information
7. **Monitor authentication**: Watch for authentication method details in verbose mode

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

For issues with specific validators, check the relevant package documentation and ensure all dependencies are properly installed. The verbose mode (`--verbose`) provides detailed execution information for debugging.

## Security Notes

- **No credential prompts**: All git operations are non-interactive
- **Environment-based auth**: Uses `GITHUB_TOKEN` environment variable or SSH keys
- **Private repository validation**: Warns about public repositories for security
- **Secure by default**: All authentication methods prioritize security
