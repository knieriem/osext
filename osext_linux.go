// Package osext provides an Executable() implementation that tries
// to work more reliable regarding the presense of /proc/self/exe
// in case of UPX compressed executables.
package osext

import (
	"os"
)

// Executable wraps os.Executable; in case of a NotExist error it tries to
// read the path to the executable from environment variable "   " (three spaces),
// as it is configured by the startup code of UPX compressed executables.
func Executable() (string, error) {
	p, err := os.Executable()
	if err != nil {
		if !os.IsNotExist(err) {
			return p, err
		}

		// UPX sets env variable "   " to the link target of /proc/self/exe. In case
		// the UPX decompression code unmaps all pages, the kernel will
		// remove the symlink /proc/self/exe (this has been observed on linux_386
		// recently). In this case os.Executable() will fail, because it cannot read
		// the symlink; variable "   " might containt the original value, though.
		//
		// see https://github.com/upx/upx/blob/7a3637ff5a800b8bcbad20ae7f668d8c8449b014/doc/elf-to-mem.txt#L104
		// and https://github.com/upx/upx/blob/7a3637ff5a800b8bcbad20ae7f668d8c8449b014/src/stub/src/i386-linux.elf-fold.S#L124
		//
		if v := os.Getenv("   "); v != "" {
			_, err1 := os.Stat(v)
			if err1 != nil {
				return p, err
			}
			p = v
		}
	}
	return p, nil
}
