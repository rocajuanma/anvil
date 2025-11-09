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

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirectoryMergeBehavior(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(sourceDir, "remote.txt"), []byte("remote"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "local.txt"), []byte("local"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "shared.txt"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "shared.txt"), []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyDirectorySimple(sourceDir, destDir)
	if err != nil {
		t.Fatalf("CopyDirectorySimple failed: %v", err)
	}

	tests := []struct {
		file    string
		content string
		desc    string
	}{
		{"remote.txt", "remote", "added from remote"},
		{"local.txt", "local", "preserved local-only file"},
		{"shared.txt", "new", "updated from remote"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(destDir, tt.file))
			if err != nil {
				t.Fatalf("File %s not found (%s): %v", tt.file, tt.desc, err)
			}
			if string(content) != tt.content {
				t.Errorf("File %s (%s): got %q, want %q", tt.file, tt.desc, string(content), tt.content)
			}
		})
	}
}

func TestCopyDirectoryNestedStructure(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	sourceNested := filepath.Join(sourceDir, "nested", "deep")
	destNested := filepath.Join(destDir, "nested")

	if err := os.MkdirAll(sourceNested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destNested, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(sourceNested, "remote.conf"), []byte("remote"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destNested, "local.conf"), []byte("local"), 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyDirectorySimple(sourceDir, destDir)
	if err != nil {
		t.Fatalf("CopyDirectorySimple failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "nested", "deep", "remote.conf")); err != nil {
		t.Error("Remote nested file not copied")
	}
	if _, err := os.Stat(filepath.Join(destDir, "nested", "local.conf")); err != nil {
		t.Error("Local nested file not preserved")
	}
}

func TestCopyFileSimple(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	destFile := filepath.Join(tempDir, "dest.txt")

	content := []byte("test content")
	if err := os.WriteFile(sourceFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyFileSimple(sourceFile, destFile)
	if err != nil {
		t.Fatalf("CopyFileSimple failed: %v", err)
	}

	got, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("Content mismatch: got %q, want %q", string(got), string(content))
	}
}

func TestCopyFileOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	destFile := filepath.Join(tempDir, "dest.txt")

	if err := os.WriteFile(sourceFile, []byte("new"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destFile, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyFileSimple(sourceFile, destFile)
	if err != nil {
		t.Fatalf("CopyFileSimple failed: %v", err)
	}

	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new" {
		t.Errorf("File not overwritten: got %q, want %q", string(content), "new")
	}
}

func TestCopyDirectoryHiddenFiles(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(sourceDir, ".hidden"), []byte("hidden"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "visible.txt"), []byte("visible"), 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyDirectorySimple(sourceDir, destDir)
	if err != nil {
		t.Fatalf("CopyDirectorySimple failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, ".hidden")); err != nil {
		t.Error("Hidden file not copied")
	}
	if _, err := os.Stat(filepath.Join(destDir, "visible.txt")); err != nil {
		t.Error("Visible file not copied")
	}
}

func TestCopyDirectoryEmptySource(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := CopyDirectorySimple(sourceDir, destDir)
	if err != nil {
		t.Fatalf("CopyDirectorySimple failed on empty source: %v", err)
	}

	if _, err := os.Stat(destDir); err != nil {
		t.Error("Destination directory not created")
	}
}

func TestCopyDirectorySourceNotExists(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "nonexistent")
	destDir := filepath.Join(tempDir, "dest")

	err := CopyDirectorySimple(sourceDir, destDir)
	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}
