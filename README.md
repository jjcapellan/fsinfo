# fsinfo: File System Information
fsinfo is a lightweight Go package designed to retrieve essential information about folders and drives in both Linux and Windows environments. It can be used to gather insights into the file system and storage devices on the host machine (e.g.: file browser for local webapp).

## Features
* GetFolderInfo: Retrieves information about a folder's contents, subfolders, and files.
* GetDrives: Obtain details about available drives on the system.
* GetHomePath: Retrieves the home directory path of the current user.
* FormatBytes: Converts a size in bytes into a human-readable string representation.
* SetHideDotFiles: Sets the visibility of dot files and folders when retrieving folder information.

## Documentation
This package documentation is indexed by *pkg.go.dev* [here](https://pkg.go.dev/github.com/jjcapellan/fsinfo#section-documentation).

## License
This library is licensed under the terms of the [MIT open source license](LICENSE).