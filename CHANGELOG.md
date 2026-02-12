# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Console Management Features** - New functionality for controlling console window visibility and properties:
  - `GetConsoleWindow()` - Retrieves the window handle of the console associated with the calling process
  - `ShowConsole()` - Shows the console window
  - `HideConsole()` - Hides the console window
  - `ToggleConsole()` - Toggles the visibility of the console window
  - `IsConsoleVisible()` - Checks if the console window is currently visible
  - `GetConsoleTitle()` - Retrieves the title of the console window
  - `SetConsoleTitle()` - Sets the title of the console window
  - `ShowWindow()` - Generic window show/hide function with various display options
  - Support for all standard ShowWindow commands (SW_HIDE, SW_SHOW, SW_SHOWNORMAL, etc.)

### Changed
- Updated README.md with comprehensive documentation for new console management features
- Enhanced project documentation with usage examples and API references

### Fixed
- Corrected package declaration in clipboard_test.go
- Fixed file naming issue (cliboard_test.go → clipboard_test.go)

## [1.0.0] - 2024-XX-XX

### Added
- Initial release with core Win32 utilities
- System tray icon management
- Menu support with emoji icons
- DPI scaling awareness
- Text input dialogs
- Toast notifications
- Advanced TrayApp builder API
- Pure Go implementation without CGO dependencies

[Unreleased]: https://github.com/Ink-33/win32utils/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Ink-33/win32utils/releases/tag/v1.0.0