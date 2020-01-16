// +build !linux

package osext

import (
	"os"
)

// Executable, on this platform, wraps os.Executable without doing any
// extra processing.
func Executable() (string, error) {
	return os.Executable()
}
