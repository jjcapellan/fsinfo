# v0.1.0

Features:
- Added the `GetFolderInfo` function to retrieve information about a folder's contents and its parent directory.
- Added the `GetDrives` function to retrieve information about available drives on the system.
- Added the `GetHomePath` function to retrieve the home directory path of the current user.

Documentation:
- Created documentation for the `GetFolderInfo` function, including its purpose, parameters, return values, and example usage.
- Created documentation for the `GetDrives` function, explaining its purpose, return values, and example usage.
- Created documentation for the `GetHomePath` function, detailing its purpose, return values, and example usage.

Code:
- Implemented the `GetFolderInfo` function with support for both Linux and Windows environments.
- Implemented the `GetDrives` function to gather information about drives on both Linux and Windows systems.
- Implemented the `GetHomePath` function to retrieve the home directory path of the current user.
- Created data structures (`FolderInfo`, `Folder`, `File`, `DriveInfo`) to hold relevant information.

Notes:
- The package aims to provide essential folder and drive information for local web applications and file browsers.
- The initial version (v0.1.0) introduces core functionality and basic documentation.