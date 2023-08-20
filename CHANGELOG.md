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