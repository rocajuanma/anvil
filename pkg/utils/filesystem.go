/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package utils provides common utility functions for file system operations.
//
// This package consolidates duplicated file operations across the codebase
// to follow the DRY (Don't Repeat Yourself) principle.
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rocajuanma/anvil/pkg/constants"
)

// CopyFileOptions holds options for file copying operations
type CopyFileOptions struct {
	CreateDirs   bool        // Create destination directory if it doesn't exist
	Overwrite    bool        // Overwrite destination file if it exists
	PreservePerm bool        // Preserve source file permissions
	FileMode     os.FileMode // File mode to use if not preserving permissions
}

// DefaultCopyFileOptions returns sensible default options for file copying
func DefaultCopyFileOptions() CopyFileOptions {
	return CopyFileOptions{
		CreateDirs:   true,
		Overwrite:    true,
		PreservePerm: false,
		FileMode:     constants.FilePerm,
	}
}

// CopyFile copies a file from src to dst with configurable options.
// This consolidates the three different copyFile implementations across the codebase.
func CopyFile(src, dst string, options CopyFileOptions) error {
	// Check if source file exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source file error: %w", err)
	}

	if srcInfo.IsDir() {
		return fmt.Errorf("source path is a directory, not a file: %s", src)
	}

	// Check if destination exists and handle overwrite
	if _, err := os.Stat(dst); err == nil && !options.Overwrite {
		return fmt.Errorf("destination file already exists: %s", dst)
	}

	// Create destination directory if needed
	if options.CreateDirs {
		if err := EnsureDirectory(filepath.Dir(dst)); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Determine file mode
	fileMode := options.FileMode
	if options.PreservePerm {
		fileMode = srcInfo.Mode()
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

// CopyFileSimple is a convenience function that uses default options.
// This matches the behavior of most existing copyFile implementations.
func CopyFileSimple(src, dst string) error {
	return CopyFile(src, dst, DefaultCopyFileOptions())
}

// CopyDirectoryOptions holds options for directory copying operations
type CopyDirectoryOptions struct {
	Overwrite     bool        // Remove destination directory before copying
	PreservePerms bool        // Preserve source permissions
	IncludeHidden bool        // Include hidden files and directories
	FileMode      os.FileMode // Default file mode if not preserving permissions
	DirMode       os.FileMode // Default directory mode if not preserving permissions
}

// DefaultCopyDirectoryOptions returns sensible default options for directory copying
func DefaultCopyDirectoryOptions() CopyDirectoryOptions {
	return CopyDirectoryOptions{
		Overwrite:     true,
		PreservePerms: false,
		IncludeHidden: true,
		FileMode:      constants.FilePerm,
		DirMode:       constants.DirPerm,
	}
}

// CopyDirectory recursively copies a directory from src to dst with configurable options.
// This consolidates the two different copyDirRecursive implementations across the codebase.
func CopyDirectory(src, dst string, options CopyDirectoryOptions) error {
	// Check if source directory exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source directory error: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", src)
	}

	// Remove destination if it exists and overwrite is enabled
	if options.Overwrite {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove existing destination: %w", err)
		}
	}

	// Walk the source directory and copy everything
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}

		// Skip hidden files/directories if not including them
		if !options.IncludeHidden && isHidden(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory
			dirMode := options.DirMode
			if options.PreservePerms {
				dirMode = info.Mode()
			}
			return os.MkdirAll(destPath, dirMode)
		} else {
			// Copy file
			fileOptions := CopyFileOptions{
				CreateDirs:   true,
				Overwrite:    options.Overwrite,
				PreservePerm: options.PreservePerms,
				FileMode:     options.FileMode,
			}
			return CopyFile(path, destPath, fileOptions)
		}
	})
}

// CopyDirectorySimple is a convenience function that uses default options.
// This matches the behavior of most existing copyDirRecursive implementations.
func CopyDirectorySimple(src, dst string) error {
	return CopyDirectory(src, dst, DefaultCopyDirectoryOptions())
}

// isHidden checks if a file/directory name represents a hidden item
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// EnsureDirectory creates a directory and all necessary parent directories
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, constants.DirPerm)
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// PathExists checks if a path (file or directory) exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
