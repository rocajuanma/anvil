package system

import (
	"runtime"
)

func getType() string {
	return runtime.GOOS
}

func IsMacOS() bool {
	return getType() == "darwin"
}

func IsLinux() bool {
	return getType() == "linux"
}
