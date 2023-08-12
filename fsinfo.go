package fsinfo

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type FolderInfo struct {
	Path    string   `json:"path"`    // Absolute folder's path
	Dir     string   `json:"parent"`  // Directory where is the folder
	Folders []Folder `json:"folders"` // Folders contained in current folder
	Files   []File   `json:"files"`   // Files contained in current folder
}

type Folder struct {
	Name string // Folder's name
	Path string // Absolute folder's path
}

type File struct {
	Name string // File's name
	Path string // Absolute file's path
}

type DriveInfo struct {
	Name string
	Path string
}

var CURRENT_DIR, _ = os.Getwd()

func GetFolderInfo(path string) (*FolderInfo, error) {

	folderInfo := &FolderInfo{}
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(CURRENT_DIR, path)
	}
	dir := filepath.Dir(path)

	// In windows, filepath.Clean() could return something like "c:."
	// and backslashes could cause problems for example in golang templates
	if runtime.GOOS == "windows" {
		path = strings.TrimRight(path, ".")
		path = filepath.ToSlash(path)
		dir = strings.TrimRight(dir, ".")
		dir = filepath.ToSlash(dir)
	}

	folderInfo.Path = path
	folderInfo.Dir = dir

	folders := []Folder{}
	files := []File{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, v := range entries {
		if v.IsDir() {
			folder := Folder{}
			folder.Name = v.Name()
			folder.Path = path + "/" + folder.Name
			folders = append(folders, folder)
			continue
		}
		file := File{}
		file.Name = v.Name()
		file.Path = path + "/" + file.Name
		files = append(files, file)
	}

	folderInfo.Files = files
	folderInfo.Folders = folders

	return folderInfo, err
}

func GetDrives() []DriveInfo {
	drives := []DriveInfo{}

	if runtime.GOOS == "linux" {

		data, _ := os.ReadFile("/proc/self/mountinfo")

		var rgx = regexp.MustCompile(`\s\/media\/\S+\b`)
		matches := rgx.FindAllString(string(data), -1)
		for _, v := range matches {
			drive := DriveInfo{}

			path := strings.ReplaceAll(v, `\040`, " ")
			path = strings.TrimSpace(path)
			drive.Path = path

			parts := strings.Split(path, "/")
			name := parts[len(parts)-1]
			if isUUID(name) {
				f, _ := os.Open(path)
				info, _ := f.Stat()
				name = "Volume of " + convSize(info.Size())
				f.Close()
			}

			drive.Name = name

			drives = append(drives, drive)
		}
	}

	return drives
}

func isUUID(str string) bool {
	l := len(str)
	if l != 9 && l != 16 {
		return false
	}

	// Hex number --> XXXXXXXXXXXXXXXX
	if l == 16 {
		_, err := hex.DecodeString(str)
		if err != nil {
			return false
		}
	}

	// Hex number --> XXXX-XXXX
	if l == 9 {
		if str[4] != '-' {
			return false
		}
		n1 := str[:4]
		n2 := str[5:]
		_, err := hex.DecodeString(n1)
		if err != nil {
			return false
		}
		_, err = hex.DecodeString(n2)
		if err != nil {
			return false
		}
	}

	return true
}

func convSize(size int64) string {
	sizeF := float64(size)
	sufix := "byte"
	// bytes to kb
	if sizeF > 999 {
		sizeF /= 1024
		sufix = "Kb"
	}
	// Kb to Mb
	if sizeF > 999 {
		sizeF /= 1024
		sufix = "Mb"
	}
	// Mb to Gb
	if sizeF > 999 {
		sizeF /= 1024
		sufix = "Gb"
	}

	size = int64(sizeF)

	return strconv.Itoa(int(size)) + sufix
}
