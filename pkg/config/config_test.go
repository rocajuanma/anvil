package config

import (
	"os"
	"testing"

	"github.com/rocajuanma/anvil/pkg/constants"
)

// setupTestConfig creates a test configuration with temporary directories
func setupTestConfig(t *testing.T) (string, func()) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Mock the config path
	originalHome := os.Getenv("HOME")
	cleanup := func() {
		os.Setenv("HOME", originalHome)
	}

	os.Setenv("HOME", tempDir)

	// Create a test config
	config := GetDefaultConfig()
	config.Directories.Config = tempDir

	// Create the necessary directories
	err := CreateDirectories()
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Save the initial config
	err = SaveConfig(config)
	if err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	return tempDir, cleanup
}

func TestAddInstalledApp(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test adding a new app
	testApp := "test-app"
	err := AddInstalledApp(testApp)
	if err != nil {
		t.Fatalf("Failed to add installed app: %v", err)
	}

	// Load config and verify the app was added
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 1 {
		t.Errorf("Expected 1 installed app, got %d", len(updatedConfig.Tools.InstalledApps))
	}

	if updatedConfig.Tools.InstalledApps[0] != testApp {
		t.Errorf("Expected app '%s', got '%s'", testApp, updatedConfig.Tools.InstalledApps[0])
	}
}

func TestAddInstalledAppDuplicate(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add an existing app to the config
	err := AddInstalledApp("existing-app")
	if err != nil {
		t.Fatalf("Failed to add existing app: %v", err)
	}

	// Test adding the same app again
	err = AddInstalledApp("existing-app")
	if err != nil {
		t.Fatalf("Failed to add installed app: %v", err)
	}

	// Load config and verify no duplicate was added
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 1 {
		t.Errorf("Expected 1 installed app (no duplicate), got %d", len(updatedConfig.Tools.InstalledApps))
	}
}

func TestAddInstalledAppSkipRequired(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test adding a required tool - should not be added to installed apps
	err := AddInstalledApp(constants.PkgGit)
	if err != nil {
		t.Fatalf("Failed to add required tool: %v", err)
	}

	// Load config and verify the required tool was not added to installed apps
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 0 {
		t.Errorf("Expected 0 installed apps (required tool should not be tracked), got %d", len(updatedConfig.Tools.InstalledApps))
	}
}

func TestGetInstalledApps(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add test apps
	testApps := []string{"app1", "app2", "app3"}
	for _, app := range testApps {
		err := AddInstalledApp(app)
		if err != nil {
			t.Fatalf("Failed to add test app %s: %v", app, err)
		}
	}

	// Test getting installed apps
	installedApps, err := GetInstalledApps()
	if err != nil {
		t.Fatalf("Failed to get installed apps: %v", err)
	}

	if len(installedApps) != 3 {
		t.Errorf("Expected 3 installed apps, got %d", len(installedApps))
	}

	expectedApps := []string{"app1", "app2", "app3"}
	for i, app := range installedApps {
		if app != expectedApps[i] {
			t.Errorf("Expected app '%s' at index %d, got '%s'", expectedApps[i], i, app)
		}
	}
}

func TestIsAppTracked(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add a tracked app
	err := AddInstalledApp("tracked-app")
	if err != nil {
		t.Fatalf("Failed to add tracked app: %v", err)
	}

	// Test checking if app is tracked
	isTracked, err := IsAppTracked("tracked-app")
	if err != nil {
		t.Fatalf("Failed to check if app is tracked: %v", err)
	}

	if !isTracked {
		t.Error("Expected app to be tracked")
	}

	// Test checking if app is not tracked
	isTracked, err = IsAppTracked("untracked-app")
	if err != nil {
		t.Fatalf("Failed to check if app is tracked: %v", err)
	}

	if isTracked {
		t.Error("Expected app to not be tracked")
	}

	// Test checking if required tool is tracked
	isTracked, err = IsAppTracked(constants.PkgGit)
	if err != nil {
		t.Fatalf("Failed to check if required tool is tracked: %v", err)
	}

	if !isTracked {
		t.Error("Expected required tool to be tracked")
	}
}

func TestRemoveInstalledApp(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Add test apps
	testApps := []string{"app1", "app2", "app3"}
	for _, app := range testApps {
		err := AddInstalledApp(app)
		if err != nil {
			t.Fatalf("Failed to add test app %s: %v", app, err)
		}
	}

	// Test removing an app
	err := RemoveInstalledApp("app2")
	if err != nil {
		t.Fatalf("Failed to remove installed app: %v", err)
	}

	// Load config and verify the app was removed
	updatedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load updated config: %v", err)
	}

	if len(updatedConfig.Tools.InstalledApps) != 2 {
		t.Errorf("Expected 2 installed apps after removal, got %d", len(updatedConfig.Tools.InstalledApps))
	}

	expectedApps := []string{"app1", "app3"}
	for i, app := range updatedConfig.Tools.InstalledApps {
		if app != expectedApps[i] {
			t.Errorf("Expected app '%s' at index %d, got '%s'", expectedApps[i], i, app)
		}
	}
}
