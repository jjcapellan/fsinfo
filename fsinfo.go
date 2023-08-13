package fsinfo

import (
	"bufio"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
var DRIVES []DriveInfo

func Init() string {
	out, _ := exec.Command("df", "-H").Output()
	return string(out)
}

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

/**
* GetDrives gets the names and full paths of storage volumes,
* and returns it as []DriveInfo.
* This function could have an impact on performance, especially on Windows,
* so it is recommended to use it only if necessary.
* In windows:
* DriveInfo.Name --> drive letter (e.g.: "c:", "d:", "e:")
* DriveInfo.Path --> drive letter + forward slash (e.g.: "c:/")
* In linux:
* DriveInfo.Name --> drive label || "Volume of size xx" (e.g.: "VolMusic" || "Volume of size 120G")
* DriveInfo.Path -->
 */
func GetDrives() []DriveInfo {

	if runtime.GOOS == "linux" {
		return getLinuxDrives()
	}

	if runtime.GOOS == "windows" {
		return getWindowsDrives()
	}

	return []DriveInfo{}
}

func getLinuxDrives() []DriveInfo {
	drives := []DriveInfo{}

	// This command reports mounts information
	out, _ := exec.Command("df", "-H").Output()
	str := string(out)

	// All paths of storage volumes starts with string " /media/"
	scanner := bufio.NewScanner(strings.NewReader(str))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, " /media/") {
			lines = append(lines, line)
		}
	}

	// Regex to extract size, name and path from each line
	var rgxSize = regexp.MustCompile(`\w+(\s{2}\w+%)`)
	var rgxName = regexp.MustCompile(`[\w\s]+$`)
	var rgxPath = regexp.MustCompile(`\s/media/.+$`)

	for _, line := range lines {
		drive := DriveInfo{}
		// Name (Disk label || "Volume of size x")
		name := rgxName.FindString(line)
		if isUUID(name) {
			size, _, _ := strings.Cut(rgxSize.FindString(line), "  ")
			name = "Volume of " + size
		}
		drive.Name = name
		// Path
		drive.Path = strings.TrimSpace(rgxPath.FindString(line))
		drives = append(drives, drive)
	}

	return drives
}

func getWindowsDrives() []DriveInfo {
	drives := []DriveInfo{}
	for _, letter := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		fsEntry, err := os.Open(string(letter) + `:\`)
		if err == nil {
			drive := DriveInfo{}
			drive.Name = string(letter) + ":"
			drive.Path = drive.Name + `/`
			drives = append(drives, drive)
			fsEntry.Close()
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
