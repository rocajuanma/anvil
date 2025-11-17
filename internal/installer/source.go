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

package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rocajuanma/anvil/internal/config"
	"github.com/rocajuanma/anvil/internal/system"
	"github.com/rocajuanma/anvil/internal/terminal/charm"
	"github.com/rocajuanma/anvil/internal/utils"
)

// InstallFromSource installs an application from a source URL
func InstallFromSource(appName, sourceURL string) error {
	spinner := charm.NewDotsSpinner(fmt.Sprintf("Downloading %s from source", appName))
	spinner.Start()

	downloadedFile, err := downloadFile(sourceURL, appName)
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed to download %s", appName))
		return fmt.Errorf("failed to download %s: %w", appName, err)
	}
	spinner.Success(fmt.Sprintf("Downloaded %s", appName))

	spinner = charm.NewDotsSpinner(fmt.Sprintf("Installing %s", appName))
	spinner.Start()

	if err := installDownloadedFile(downloadedFile, appName); err != nil {
		spinner.Error(fmt.Sprintf("Failed to install %s", appName))
		return fmt.Errorf("failed to install %s: %w", appName, err)
	}

	spinner.Success(fmt.Sprintf("%s installed successfully", appName))
	return nil
}

// downloadFile downloads a file from URL to a temporary location
func downloadFile(fileURL, appName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "anvil-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	homeDir, _ := system.GetHomeDir()
	downloadsDir := filepath.Join(homeDir, "Downloads", "anvil-downloads")
	if err := utils.EnsureDirectory(downloadsDir); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %w", err)
	}

	fileName := getFileNameFromURL(fileURL, appName)
	filePath := filepath.Join(downloadsDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// getFileNameFromURL extracts filename from URL or uses app name
func getFileNameFromURL(fileURL, appName string) string {
	parsedURL, err := url.Parse(fileURL)
	if err == nil && parsedURL.Path != "" {
		fileName := filepath.Base(parsedURL.Path)
		if fileName != "" && fileName != "/" {
			return fileName
		}
	}

	ext := getExtensionFromURL(fileURL)
	return fmt.Sprintf("%s%s", appName, ext)
}

// getExtensionFromURL tries to detect file extension from URL
func getExtensionFromURL(fileURL string) string {
	parsedURL, err := url.Parse(fileURL)
	if err == nil {
		path := strings.ToLower(parsedURL.Path)
		extensions := []string{".dmg", ".pkg", ".zip", ".tar.gz", ".deb", ".rpm", ".AppImage", ".tar.bz2"}
		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				return ext
			}
		}
	}
	return ".zip"
}

// installDownloadedFile installs the downloaded file based on its type and OS
func installDownloadedFile(filePath, appName string) error {
	if system.IsMacOS() {
		return installOnMacOS(filePath, appName)
	} else if system.IsLinux() {
		return installOnLinux(filePath, appName)
	}
	return fmt.Errorf("unsupported operating system")
}

// installOnMacOS handles installation on macOS
func installOnMacOS(filePath, appName string) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".dmg":
		return installDMG(filePath, appName)
	case ".pkg":
		return installPKG(filePath)
	case ".zip":
		return installZIP(filePath, appName)
	default:
		return fmt.Errorf("unsupported file type: %s (supported: .dmg, .pkg, .zip)", ext)
	}
}

// installOnLinux handles installation on Linux
func installOnLinux(filePath, appName string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	baseName := strings.ToLower(filepath.Base(filePath))

	if strings.HasSuffix(baseName, ".tar.gz") {
		return installTarGz(filePath, appName)
	} else if strings.HasSuffix(baseName, ".tar.bz2") {
		return installTarBz2(filePath, appName)
	}

	switch ext {
	case ".deb":
		return installDEB(filePath)
	case ".rpm":
		return installRPM(filePath)
	case ".AppImage":
		return installAppImage(filePath, appName)
	case ".zip":
		return installZIP(filePath, appName)
	default:
		return fmt.Errorf("unsupported file type: %s (supported: .deb, .rpm, .AppImage, .zip, .tar.gz, .tar.bz2)", ext)
	}
}

// installDMG mounts DMG, copies .app to Applications, and unmounts
func installDMG(filePath, appName string) error {
	spinner := charm.NewDotsSpinner("Mounting DMG")
	spinner.Start()

	mountResult, err := system.RunCommand("hdiutil", "attach", filePath, "-nobrowse", "-quiet")
	if err != nil || !mountResult.Success {
		spinner.Error("Failed to mount DMG")
		return fmt.Errorf("failed to mount DMG: %s", mountResult.Error)
	}

	mountPath := extractMountPath(mountResult.Output)
	if mountPath == "" {
		system.RunCommand("hdiutil", "detach", mountPath, "-quiet")
		spinner.Error("Failed to find mount path")
		return fmt.Errorf("failed to extract mount path from DMG")
	}
	spinner.Success("DMG mounted")

	defer func() {
		system.RunCommand("hdiutil", "detach", mountPath, "-quiet")
	}()

	spinner = charm.NewDotsSpinner("Finding application")
	spinner.Start()

	appPath := findAppInDirectory(mountPath, appName)
	if appPath == "" {
		spinner.Error("Application not found")
		return fmt.Errorf("failed to find .app in DMG")
	}
	spinner.Success("Application found")

	spinner = charm.NewDotsSpinner("Installing to Applications")
	spinner.Start()

	homeDir, _ := system.GetHomeDir()
	applicationsDir := filepath.Join(homeDir, "Applications")
	if err := utils.EnsureDirectory(applicationsDir); err != nil {
		spinner.Error("Failed to create Applications directory")
		return fmt.Errorf("failed to create Applications directory: %w", err)
	}

	appNameFromPath := filepath.Base(appPath)
	destPath := filepath.Join(applicationsDir, appNameFromPath)

	if err := utils.CopyDirectorySimple(appPath, destPath); err != nil {
		spinner.Error("Failed to copy application")
		return fmt.Errorf("failed to copy application: %w", err)
	}

	spinner.Success("Application installed")
	return nil
}

// installPKG installs a .pkg file using installer command
func installPKG(filePath string) error {
	spinner := charm.NewDotsSpinner("Installing package")
	spinner.Start()

	result, err := system.RunCommand("sudo", "installer", "-pkg", filePath, "-target", "/")
	if err != nil || !result.Success {
		spinner.Error("Failed to install package")
		return fmt.Errorf("failed to install package: %s", result.Error)
	}

	spinner.Success("Package installed")
	return nil
}

// installZIP extracts ZIP and handles contents
func installZIP(filePath, appName string) error {
	spinner := charm.NewDotsSpinner("Extracting ZIP")
	spinner.Start()

	extractDir := filepath.Join(filepath.Dir(filePath), appName+"-extracted")
	if err := utils.EnsureDirectory(extractDir); err != nil {
		spinner.Error("Failed to create extract directory")
		return fmt.Errorf("failed to create extract directory: %w", err)
	}

	result, err := system.RunCommand("unzip", "-q", filePath, "-d", extractDir)
	if err != nil || !result.Success {
		spinner.Error("Failed to extract ZIP")
		return fmt.Errorf("failed to extract ZIP: %s", result.Error)
	}
	spinner.Success("ZIP extracted")

	if system.IsMacOS() {
		return handleExtractedContentsMacOS(extractDir, appName)
	}
	return handleExtractedContentsLinux(extractDir, appName)
}

// installDEB installs a .deb package
func installDEB(filePath string) error {
	spinner := charm.NewDotsSpinner("Installing DEB package")
	spinner.Start()

	result, err := system.RunCommand("sudo", "dpkg", "-i", filePath)
	if err != nil || !result.Success {
		spinner.Error("Failed to install DEB package")
		return fmt.Errorf("failed to install DEB package: %s", result.Error)
	}

	result, err = system.RunCommand("sudo", "apt-get", "-f", "install", "-y")
	if err != nil || !result.Success {
		spinner.Warning("Dependency resolution had issues")
	}

	spinner.Success("DEB package installed")
	return nil
}

// installRPM installs an .rpm package
func installRPM(filePath string) error {
	spinner := charm.NewDotsSpinner("Installing RPM package")
	spinner.Start()

	var result *system.CommandResult
	var err error

	if system.CommandExists("dnf") {
		result, err = system.RunCommand("sudo", "dnf", "install", "-y", filePath)
	} else if system.CommandExists("yum") {
		result, err = system.RunCommand("sudo", "yum", "install", "-y", filePath)
	} else {
		result, err = system.RunCommand("sudo", "rpm", "-i", filePath)
	}

	if err != nil || !result.Success {
		spinner.Error("Failed to install RPM package")
		return fmt.Errorf("failed to install RPM package: %s", result.Error)
	}

	spinner.Success("RPM package installed")
	return nil
}

// installAppImage makes AppImage executable and optionally installs it
func installAppImage(filePath, appName string) error {
	spinner := charm.NewDotsSpinner("Setting up AppImage")
	spinner.Start()

	homeDir, _ := system.GetHomeDir()
	appImageDir := filepath.Join(homeDir, "Applications")
	if err := utils.EnsureDirectory(appImageDir); err != nil {
		spinner.Error("Failed to create Applications directory")
		return fmt.Errorf("failed to create Applications directory: %w", err)
	}

	destPath := filepath.Join(appImageDir, filepath.Base(filePath))
	if err := utils.CopyFileSimple(filePath, destPath); err != nil {
		spinner.Error("Failed to copy AppImage")
		return fmt.Errorf("failed to copy AppImage: %w", err)
	}

	result, err := system.RunCommand("chmod", "+x", destPath)
	if err != nil || !result.Success {
		spinner.Error("Failed to make AppImage executable")
		return fmt.Errorf("failed to make AppImage executable: %s", result.Error)
	}

	spinner.Success("AppImage installed")
	return nil
}

// installTarGz extracts and installs .tar.gz archive
func installTarGz(filePath, appName string) error {
	return installTarArchive(filePath, appName, "tar", "-xzf")
}

// installTarBz2 extracts and installs .tar.bz2 archive
func installTarBz2(filePath, appName string) error {
	return installTarArchive(filePath, appName, "tar", "-xjf")
}

// installTarArchive extracts tar archive and handles contents
func installTarArchive(filePath, appName, command, flags string) error {
	spinner := charm.NewDotsSpinner("Extracting archive")
	spinner.Start()

	extractDir := filepath.Join(filepath.Dir(filePath), appName+"-extracted")
	if err := utils.EnsureDirectory(extractDir); err != nil {
		spinner.Error("Failed to create extract directory")
		return fmt.Errorf("failed to create extract directory: %w", err)
	}

	result, err := system.RunCommand(command, flags, filePath, "-C", extractDir)
	if err != nil || !result.Success {
		spinner.Error("Failed to extract archive")
		return fmt.Errorf("failed to extract archive: %s", result.Error)
	}
	spinner.Success("Archive extracted")

	return handleExtractedContentsLinux(extractDir, appName)
}

// handleExtractedContentsMacOS handles extracted contents on macOS
func handleExtractedContentsMacOS(extractDir, appName string) error {
	appPath := findAppInDirectory(extractDir, appName)
	if appPath == "" {
		return fmt.Errorf("failed to find .app in extracted contents")
	}

	homeDir, _ := system.GetHomeDir()
	applicationsDir := filepath.Join(homeDir, "Applications")
	if err := utils.EnsureDirectory(applicationsDir); err != nil {
		return fmt.Errorf("failed to create Applications directory: %w", err)
	}

	appNameFromPath := filepath.Base(appPath)
	destPath := filepath.Join(applicationsDir, appNameFromPath)

	return utils.CopyDirectorySimple(appPath, destPath)
}

// handleExtractedContentsLinux handles extracted contents on Linux
func handleExtractedContentsLinux(extractDir, appName string) error {
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return fmt.Errorf("failed to read extract directory: %w", err)
	}

	if len(entries) == 1 && entries[0].IsDir() {
		appDir := filepath.Join(extractDir, entries[0].Name())
		homeDir, _ := system.GetHomeDir()
		destDir := filepath.Join(homeDir, ".local", "share", "applications", entries[0].Name())

		if err := utils.EnsureDirectory(filepath.Dir(destDir)); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		return utils.CopyDirectorySimple(appDir, destDir)
	}

	homeDir, _ := system.GetHomeDir()
	destDir := filepath.Join(homeDir, ".local", "share", "applications", appName)
	if err := utils.EnsureDirectory(filepath.Dir(destDir)); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return utils.CopyDirectorySimple(extractDir, destDir)
}

// extractMountPath extracts the mount path from hdiutil output
func extractMountPath(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "/Volumes/") {
					return part
				}
			}
		}
	}
	return ""
}

// findAppInDirectory searches for .app bundle in directory
func findAppInDirectory(dir, appName string) string {
	var foundApp string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && strings.HasSuffix(path, ".app") {
			baseName := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ".app"))
			searchName := strings.ToLower(appName)
			if baseName == searchName || strings.Contains(baseName, searchName) {
				foundApp = path
				return filepath.SkipDir
			}
		}
		return nil
	})
	return foundApp
}

// GetSourceURL returns the source URL for an app if it exists
func GetSourceURL(appName string) (string, bool, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", false, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Sources == nil {
		return "", false, nil
	}

	sourceURL, exists := cfg.Sources[appName]
	if !exists || sourceURL == "" {
		return "", false, nil
	}

	return sourceURL, true, nil
}
