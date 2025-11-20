# Contributing to Anvil CLI

Thank you for your interest in contributing to Anvil! This guide will help you get started.

## Ways to Contribute

- **Bug Reports** - Help us identify and fix issues
- **Feature Requests** - Suggest new functionality
- **Documentation** - Improve guides and examples
- **Code Contributions** - Fix bugs, implement features
- **Testing** - Add test cases and improve coverage

## Getting Started

1. **Check existing issues** to avoid duplicate work
2. **Start small** with documentation or minor bug fixes
3. **Ask questions** if anything is unclear

## Development Setup

### Prerequisites

- **Go 1.17+** - [Download](https://golang.org/dl/)
- **Git** - Version control

### Setup

```bash
# Fork the repository on GitHub first
git clone https://github.com/yourusername/anvil.git
cd anvil

# Add upstream remote
git remote add upstream https://github.com/0xjuanma/anvil.git

# Build the project
go build -o anvil main.go
```

## Contributing Guidelines

### Bug Reports

Include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, Anvil version)

### Feature Requests

Include:
- Clear description of the feature
- Use case and problem it solves
- Proposed solution

### Pull Requests

Before submitting:
- Search existing PRs to avoid duplicates
- Create an issue for significant changes
- Write tests for new functionality
- Update documentation as needed
- Test thoroughly

## Development Workflow

1. **Create feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make changes** following the code style guidelines

3. **Test changes**
   ```bash
   go test ./...
   go build -o anvil-dev main.go
   ```

4. **Commit with conventional messages**
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve bug in install command"
   ```

5. **Submit pull request**

## Code Style

Follow standard Go conventions and use the existing code patterns in the project.