# v0.1.3
## Changed
- `Folder` struct includes 1 new property: *ModTime* (time.Time).
## Fixed
- `GetFolderInfo` now correctly handles situations where information is retrieved for non-existent files or folders.

# v0.1.2
## Changed
- `FormatBytes` : The string returned by this function has changed its format. The string has precision of 1 decimal and trailing zeros are removed. (e.g. : "126 B", "26,2 KB")

# v0.1.1
## Added
- `FormatBytes` utility function converts a size in bytes into a human-readable string representation.

## Changed
- `File` struct includes 2 new properties: *Size* (int64) and *ModTime* (time.Time).  

# v0.1.0
## Added:
- `GetFolderInfo` function retrieves information about a folder's contents and its parent directory.
- `GetDrives` function retrieves information about available drives on the system.
- `GetHomePath` function retrieves the home directory path of the current user.

## Notes:
- The package aims to provide essential folder and drive information for local web applications and file browsers.
- The initial version (v0.1.0) introduces core functionality and basic documentation.