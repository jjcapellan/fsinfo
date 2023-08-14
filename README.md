# fsinfo: File System Information
fsinfo is a lightweight Go package designed to retrieve essential information about folders and drives in both Linux and Windows environments. It can be used to gather insights into the file system and storage devices on the host machine (e.g.: file browser for local webapp).

## Features
* GetFolderInfo: Retrieve information about a folder's contents, subfolders, and files.
* GetDrives: Obtain details about available drives on the system.

## Installation
You can install the fsinfo package using the following command:  
```shell
go get github.com/jjcapellan/fsinfo
```  
## Usage
Import the package into your Go code:
```go
import "github.com/jjcapellan/fsinfo"
```
### GetFolderInfo
The GetFolderInfo function provides information about the contents of a given folder. Here's how to use it:  
```go
folderInfo, err := fsinfo.GetFolderInfo("/path/to/folder")
if err != nil {
    fmt.Println("Error:", err)
    return
}

fmt.Println("Folder Path:", folderInfo.Path)
fmt.Println("Parent Directory:", folderInfo.Dir)
fmt.Println("Subfolders:")
for _, folder := range folderInfo.Folders {
    fmt.Println("  -", folder.Name, ":", folder.Path)
}
fmt.Println("Files:")
for _, file := range folderInfo.Files {
    fmt.Println("  -", file.Name, ":", file.Path)
}
```  
### GetDrives
The GetDrives function retrieves information about available drives on the system. Here's how to use it:  
```go
drives, err := fsinfo.GetDrives()
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Available Drives:")
for _, drive := range drives {
    fmt.Println("  -", drive.Name, ":", drive.Path)
}
```  
## License
This library is licensed under the terms of the [MIT open source license](LICENSE).