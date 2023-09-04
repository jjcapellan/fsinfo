//go:build !windows

package fsinfo

import "io/fs"

func isHidden(file fs.DirEntry) bool {
	return file.Name()[0] == '.'
}
