#!/usr/bin/env python3
"""
Changelog automation script for Anvil CLI.

This script automatically moves [Unreleased] content to a versioned section
and creates a new empty [Unreleased] section for future changes.
"""

import argparse
import re
import sys
from datetime import datetime
from pathlib import Path


def parse_changelog(content):
    """Parse changelog content and extract sections."""
    lines = content.split('\n')
    sections = {}
    current_section = None
    current_content = []
    
    for line in lines:
        # Check if this is a version header
        version_match = re.match(r'^## \[(.*?)\](?:\s*-\s*(.*))?', line)
        if version_match:
            # Save previous section
            if current_section is not None:
                sections[current_section] = '\n'.join(current_content)
            
            # Start new section
            current_section = version_match.group(1)
            current_content = []
            
            # Include the date if present
            if version_match.group(2):
                sections[f"{current_section}_date"] = version_match.group(2)
        else:
            if current_section is not None:
                current_content.append(line)
    
    # Save the last section
    if current_section is not None:
        sections[current_section] = '\n'.join(current_content)
    
    return sections


def extract_header_content(content):
    """Extract content before the first version section."""
    lines = content.split('\n')
    header_lines = []
    
    for line in lines:
        if re.match(r'^## \[', line):
            break
        header_lines.append(line)
    
    return '\n'.join(header_lines).rstrip()


def extract_footer_content(content):
    """Extract content after all version sections."""
    lines = content.split('\n')
    footer_lines = []
    in_footer = False
    
    # Look for the end of version sections (start of footer)
    for i, line in enumerate(lines):
        if re.match(r'^---\s*$', line) and not in_footer:
            in_footer = True
            footer_lines = lines[i:]
            break
    
    return '\n'.join(footer_lines) if footer_lines else ""


def create_new_changelog(header, unreleased_content, new_version, release_date, existing_sections, footer):
    """Create the updated changelog content."""
    new_lines = []
    
    # Add header
    new_lines.append(header)
    new_lines.append("")
    
    # Add new empty Unreleased section
    new_lines.append("## [Unreleased]")
    new_lines.append("")
    new_lines.append("### Added")
    new_lines.append("")
    new_lines.append("### Changed")
    new_lines.append("")
    new_lines.append("### Fixed")
    new_lines.append("")
    
    # Add the new version section with the unreleased content
    if unreleased_content.strip():
        new_lines.append(f"## [{new_version}] - {release_date}")
        new_lines.append("")
        # Clean up the unreleased content
        unreleased_lines = unreleased_content.split('\n')
        cleaned_lines = []
        for line in unreleased_lines:
            if line.strip():  # Skip empty lines at start/end
                cleaned_lines.append(line)
        
        if cleaned_lines:
            new_lines.extend(cleaned_lines)
            new_lines.append("")
    
    # Add existing sections (skip Unreleased as we already processed it)
    for section_key in existing_sections:
        if section_key != 'Unreleased' and not section_key.endswith('_date'):
            section_content = existing_sections[section_key]
            
            # Check if this section has a date
            date_key = f"{section_key}_date"
            if date_key in existing_sections:
                section_header = f"## [{section_key}] - {existing_sections[date_key]}"
            else:
                section_header = f"## [{section_key}]"
            
            new_lines.append(section_header)
            new_lines.append("")
            
            # Add section content
            section_lines = section_content.split('\n')
            for line in section_lines:
                if line.strip():  # Skip empty lines
                    new_lines.append(line)
            new_lines.append("")
    
    # Add footer if present
    if footer.strip():
        new_lines.append(footer)
    
    return '\n'.join(new_lines)


def update_changelog_file(changelog_path, version, date):
    """Update the changelog file with the new version."""
    
    # Read current changelog
    try:
        with open(changelog_path, 'r', encoding='utf-8') as f:
            current_content = f.read()
    except FileNotFoundError:
        print(f"Error: Changelog file not found: {changelog_path}", file=sys.stderr)
        return False
    except Exception as e:
        print(f"Error reading changelog: {e}", file=sys.stderr)
        return False
    
    # Parse changelog sections
    sections = parse_changelog(current_content)
    
    # Check if Unreleased section exists and has content
    if 'Unreleased' not in sections:
        print("No [Unreleased] section found in changelog", file=sys.stderr)
        return False
    
    unreleased_content = sections['Unreleased'].strip()
    if not unreleased_content:
        print("No content found in [Unreleased] section")
        return False
    
    # Extract header and footer
    header = extract_header_content(current_content)
    footer = extract_footer_content(current_content)
    
    # Create new changelog content
    new_content = create_new_changelog(
        header, 
        unreleased_content, 
        version, 
        date, 
        sections, 
        footer
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
