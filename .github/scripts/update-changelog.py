#!/usr/bin/env python3
"""
Simple changelog automation script for Anvil CLI.

This script moves [Unreleased] content to a versioned section
and creates a new empty [Unreleased] section.
"""

import argparse
import re
import sys
from datetime import datetime
from pathlib import Path


def update_changelog_file(changelog_path, version, date):
    """Update the changelog file with the new version."""
    
    # Read current changelog
    try:
        with open(changelog_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except FileNotFoundError:
        print(f"Error: Changelog file not found: {changelog_path}", file=sys.stderr)
        return False
    except Exception as e:
        print(f"Error reading changelog: {e}", file=sys.stderr)
        return False
    
    # Find the [Unreleased] section
    unreleased_pattern = r'^## \[Unreleased\](.*?)(?=^## \[|\Z)'
    match = re.search(unreleased_pattern, content, re.MULTILINE | re.DOTALL)
    
    if not match:
        print("No [Unreleased] section found in changelog", file=sys.stderr)
        return False
    
    unreleased_content = match.group(1).strip()
    if not unreleased_content:
        print("No content found in [Unreleased] section")
        return False
    
    # Create the new version section
    new_version_section = f"## [{version}] - {date}\n\n{unreleased_content}\n\n"
    
    # Create new empty [Unreleased] section
    new_unreleased_section = """## [Unreleased]

### Added

### Changed

### Fixed

"""
    
    # Replace the [Unreleased] section with new content
    new_content = re.sub(
        unreleased_pattern,
        new_unreleased_section + new_version_section,
        content,
        flags=re.MULTILINE | re.DOTALL
    )
    
    # Write updated changelog
    try:
        with open(changelog_path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"✅ Successfully updated changelog for version {version}")
        return True
    except Exception as e:
        print(f"Error writing changelog: {e}", file=sys.stderr)
        return False


def main():
    parser = argparse.ArgumentParser(description="Update changelog after release")
    parser.add_argument("--version", required=True, help="Release version (e.g., 1.2.0)")
    parser.add_argument("--date", required=True, help="Release date (YYYY-MM-DD)")
    parser.add_argument("--changelog-path", default="docs/CHANGELOG.md", 
                       help="Path to changelog file")
    
    args = parser.parse_args()
    
    # Validate inputs
    try:
        datetime.strptime(args.date, '%Y-%m-%d')
    except ValueError:
        print(f"Error: Invalid date format: {args.date}. Use YYYY-MM-DD", file=sys.stderr)
        sys.exit(1)
    
    if not Path(args.changelog_path).exists():
        print(f"Error: Changelog file not found: {args.changelog_path}", file=sys.stderr)
        sys.exit(1)
    
    # Update changelog
    success = update_changelog_file(args.changelog_path, args.version, args.date)
    
    if not success:
        sys.exit(1)
    
    print(f"Changelog update completed for version {args.version}")


if __name__ == "__main__":
    main()