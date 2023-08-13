package fsinfo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"text/tabwriter"
)

var wd, _ = os.Getwd()

func TestGetFolderInfoRel(t *testing.T) {

	path1 := "./test_assets"
	absPath := filepath.Join(wd, path1)
	absPath = filepath.ToSlash(absPath)
	dir := filepath.Dir(absPath)
	dir = filepath.ToSlash(dir)

	folders := []Folder{
		{Name: "folder1", Path: absPath + "/" + "folder1"},
		{Name: "folder2", Path: absPath + "/" + "folder2"},
	}
	files := []File{
		{Name: ".dotfile", Path: absPath + "/" + ".dotfile"},
		{Name: "file1.txt", Path: absPath + "/" + "file1.txt"},
		{Name: "file2.txt", Path: absPath + "/" + "file2.txt"},
	}

	folderInfo, err := GetFolderInfo(path1)
	if err != nil {
		t.Fatalf("Error retrieving folderInfo\nPath: %s\nError %s", path1, err)
	}

	if folderInfo.Path != absPath {
		t.Fatalf("Error folderInfo.Path\nPath: %s\n.Path: %s\nExpected: %s", path1, folderInfo.Path, absPath)
	}

	if folderInfo.Dir != dir {
		t.Fatalf("Error folderInfo.Dir\nPath: %s\n.Dir: %s\nExpected: %s", path1, folderInfo.Dir, dir)
	}

	if len(folderInfo.Folders) != 2 {
		t.Fatalf("Error in length of folderInfo.Folders\nPath: %s\nLength: %d\nExpected: %d", path1, len(folderInfo.Folders), 2)
	}

	if len(folderInfo.Files) != 3 {
		t.Fatalf("Error in length of folderInfo.Files\nPath: %s\nLength: %d\nExpected: %d", path1, len(folderInfo.Files), 3)
	}

	_folders := folderInfo.Folders
	_files := folderInfo.Files

	sort.Slice(folders, func(i, j int) bool { return folders[i].Name < folders[j].Name })
	sort.Slice(_folders, func(i, j int) bool { return _folders[i].Name < _folders[j].Name })
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	sort.Slice(_files, func(i, j int) bool { return _files[i].Name < _files[j].Name })

	for k, v := range folders {
		_v := folders[k]
		if _v != v {
			t.Fatalf("Error in folderInfo.Folders[%d]\nPath: %s\nValue: %s\nExpected: %s", k, path1, _v, v)
		}
	}

	for k, v := range files {
		_v := files[k]
		if _v != v {
			t.Fatalf("Error in folderInfo.Files[%d]\nPath: %s\nValue: %s\nExpected: %s", k, path1, _v, v)
		}
	}

}

func TestGetDrives(t *testing.T) {
	drives := GetDrives()
	buffer := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buffer, 5, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "\nName\tPath")
	for _, drive := range drives {
		fmt.Fprintln(tw, drive.Name+"\t"+drive.Path)
	}
	tw.Flush()
	t.Log(buffer.String())
}
