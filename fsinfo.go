// Package fsinfo is a lightweight Go package designed to retrieve essential information about folders and drives
// in both Linux and Windows environments. It can be used to gather insights into the file system and storage
// devices on the host machine (e.g.: file browser for local webapp).
package fsinfo

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FolderInfo contains information about a folder and its contents.
type FolderInfo struct {
	Path    string   `json:"path"`    // Absolute folder's path
	Dir     string   `json:"parent"`  // Directory where the folder is located
	Folders []Folder `json:"folders"` // Subfolders contained within the folder
	Files   []File   `json:"files"`   // Files contained within the folder
}

type Folder struct {
	Name    string // Folder's name
	Path    string // Absolute folder's path
	ModTime time.Time
}

type File struct {
	Name    string // File's name
	Path    string // Absolute file's path
	Size    int64
	ModTime time.Time
}

type DriveInfo struct {
	Name string
	Path string
}

var current_dir, _ = os.Getwd()

// GetFolderInfo retrieves information about the contents of a specified folder and its parent directory.
// It returns a pointer to a FolderInfo struct that contains details about the folder, its subfolders,
// and the files it contains.
//
// Parameters:
//  - path (string): An absolute or relative path to the folder for which information is to be retrieved.
//
// Returns:
//  - *FolderInfo: A pointer to a FolderInfo struct with information about the folder and its contents.
//  - error: An error, if any, that occurred during the retrieval of the folder information.
//
// Example usage:
//   folderPath := "/path/to/folder"
//   folderInfo, err := GetFolderInfo(folderPath)
//   if err != nil {
//       fmt.Println("Error:", err)
//       return
//   }
//   fmt.Println("Folder Path:", folderInfo.Path)
//   fmt.Println("Parent Directory:", folderInfo.Dir)
//   fmt.Println("Subfolders:")
//   for _, subfolder := range folderInfo.Folders {
//       fmt.Println("  Name:", subfolder.Name)
//       fmt.Println("  Path:", subfolder.Path)
//   }
//   fmt.Println("Files:")
//   for _, file := range folderInfo.Files {
//       fmt.Println("  Name:", file.Name)
//       fmt.Println("  Path:", file.Path)
//   }
func GetFolderInfo(path string) (*FolderInfo, error) {

	folderInfo := &FolderInfo{}
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(current_dir, path)
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
		// err -> ErrNotExist
		finfo, err := v.Info()
		if err != nil {
			continue
		}

		if v.IsDir() {
			folder := Folder{}
			folder.Name = v.Name()
			folder.Path = path + "/" + folder.Name
			folder.ModTime = finfo.ModTime()
			folders = append(folders, folder)
			continue
		}

		file := File{}
		file.Name = finfo.Name()
		file.Path = path + "/" + file.Name
		file.Size = finfo.Size()
		file.ModTime = finfo.ModTime()
		files = append(files, file)
	}

	folderInfo.Files = files
	folderInfo.Folders = folders

	return folderInfo, err
}

// GetDrives retrieves information about the available drives on the system.
// It returns a slice of DriveInfo structs containing details about the available drives.
//
// Returns:
//  - []DriveInfo: A slice of DriveInfo structs with information about the available drives.
//  - error: An error if the information about drives cannot be retrieved.
//
// Example usage:
//   drives, err := GetDrives()
//   if err != nil {
//       fmt.Println("Error:", err)
//       return
//   }
//   fmt.Println("Available Drives:")
//   for _, drive := range drives {
//       fmt.Println("  Name:", drive.Name)
//       fmt.Println("  Path:", drive.Path)
//   }
func GetDrives() ([]DriveInfo, error) {

	if runtime.GOOS == "linux" {
		return getLinuxDrives()
	}

	if runtime.GOOS == "windows" {
		return getWindowsDrives(), nil
	}

	return []DriveInfo{}, errors.New("os not supported")
}

// GetHomePath retrieves the home directory path of the current user.
// The function returns the home directory path and an error if encountered.
//
// Returns:
//  - string: The home directory path of the current user.
//  - error: An error, if any, that occurred during the retrieval of the home directory path.
//
// Example usage:
//   homePath, err := GetHomePath()
//   if err != nil {
//       fmt.Println("Error:", err)
//       return
//   }
//   fmt.Println("Home Directory Path:", homePath)
func GetHomePath() (string, error) {

	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path = filepath.ToSlash(path)

	return path, nil
}

// getLinuxDrives gathers information about mounted drives on a Linux system.
// It utilizes the /proc/self/mountinfo file to retrieve information about mounted drives.
// It returns a list of DriveInfo structures containing the name and path of the drives,
// and an error if the information cannot be accessed.
func getLinuxDrives() ([]DriveInfo, error) {
	drives := []DriveInfo{}

	// Root file system
	drives = append(drives, DriveInfo{Name: "File system", Path: "/"})

	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return drives, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		path := fields[4]
		if !strings.HasPrefix(path, "/media/") {
			continue
		}
		path = strings.ReplaceAll(path, `\040`, " ")
		name := path[strings.LastIndexByte(path, '/')+1:]
		drives = append(drives, DriveInfo{Name: name, Path: path})
	}

	return drives, nil
}

// getWindowsDrives retrieves information about available drives on a Windows system.
// It iterates through drive letters and checks for their existence using os.Open.
// Returns a slice of DriveInfo structs containing details about the drives.
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

// FormatBytes converts a size in bytes into a human-readable string representation.
// The function takes a size in bytes as input and returns a string that includes
// the formatted size in appropriate units (bytes, kilobytes, megabytes, gigabytes, terabytes, or petabytes).
// If the size is zero, the string "0B" will be returned.
//
// Parameters:
//   - size: The size in bytes to be formatted.
//
// Returns:
//   A string representing the formatted size with the corresponding unit.
//
// Example:
//   size := int64(1234567890)
//   formattedSize := FormatBytes(size)
//   fmt.Println("File size:", formattedSize)
//
// Expected Output:
//   File size: 1.15 GB
func FormatBytes(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	var idx int
	fSize := float64(size)

	for fSize >= 1024 && idx < len(units)-1 {
		fSize /= 1024
		idx++
	}

	str := fmt.Sprintf("%.1f %s", fSize, units[idx])
	return strings.Replace(str, ".0", "", 1)
}
