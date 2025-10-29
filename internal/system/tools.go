package system

import (
	"os"
	"runtime"
)

// GetType returns the type of the operating system
func getType() string {
	return runtime.GOOS
}

func IsMacOS() bool {
	return getType() == "darwin"
}

func IsLinux() bool {
	return getType() == "linux"
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}
