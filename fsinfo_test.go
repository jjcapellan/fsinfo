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
		if _v.Path != v.Path {
			t.Fatalf("Error in folderInfo.Folders[%d]\nPath: %s\nValue: %s\nExpected: %s", k, path1, _v.Path, v.Path)
		}
	}

	for k, v := range files {
		_v := files[k]
		if _v.Path != v.Path {
			t.Fatalf("Error in folderInfo.Files[%d]\nPath: %s\nValue: %s\nExpected: %s", k, path1, _v.Path, v.Path)
		}
	}

}

func TestHideDotFiles(t *testing.T) {
	path := "./test_assets"
	absPath := filepath.Join(wd, path)
	absPath = filepath.ToSlash(absPath)

	files := []File{
		//{Name: ".dotfile", Path: absPath + "/" + ".dotfile"}, <-- hideDotFiles == true
		{Name: "file1.txt", Path: absPath + "/" + "file1.txt"},
		{Name: "file2.txt", Path: absPath + "/" + "file2.txt"},
	}

	SetHideDotFiles(true)
	folderInfo, err := GetFolderInfo(path)
	if err != nil {
		t.Fatalf("Error retrieving folderInfo\nPath: %s\nError %s", path, err)
	}
	if len(folderInfo.Files) != 2 {
		t.Fatalf("Error in length of folderInfo.Files\nPath: %s\nLength: %d\nExpected: %d", path, len(folderInfo.Files), 2)
	}

	_files := folderInfo.Files

	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	sort.Slice(_files, func(i, j int) bool { return _files[i].Name < _files[j].Name })

	for k, v := range files {
		_v := files[k]
		if _v.Path != v.Path {
			t.Fatalf("Error in folderInfo.Files[%d]\nPath: %s\nValue: %s\nExpected: %s", k, path, _v.Path, v.Path)
		}
	}
}

func TestGetDrives(t *testing.T) {
	drives, err := GetDrives()
	if err != nil {
		t.Fatal("Error retrieving drives info")
	}
	buffer := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buffer, 5, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "\nName\tPath")
	for _, drive := range drives {
		fmt.Fprintln(tw, drive.Name+"\t"+drive.Path)
	}
	tw.Flush()
	t.Log(buffer.String())
}

func TestGetHomePath(t *testing.T) {
	path, err := GetHomePath()
	if err != nil {
		t.Fatalf("Error retrievng home path: %s", err)
	}
	t.Log(path)
}

func BenchmarkGetDrives(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetDrives()
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatBytes(8456213210)
	}
}
