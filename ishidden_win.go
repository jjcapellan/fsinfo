//go:build windows

package fsinfo

import (
	"io/fs"
	"syscall"
)

func isHidden(file fs.DirEntry) bool {
	info, _ := file.Info()

	if info.IsDir() {
		return false
	}

	attributes := info.Sys().(*syscall.Win32FileAttributeData).FileAttributes
	// https://go.dev/src/syscall/types_windows.go
	if attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0 {
		return true
	}

	return false
}
