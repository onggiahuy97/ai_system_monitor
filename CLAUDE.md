# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based system monitoring tool that tracks user activity on macOS. The project monitors the active window/application and watches for file system changes to analyze productivity patterns using AI.

## Architecture

- **Single-file application**: All logic is contained in `main.go`
- **macOS-specific**: Uses AppleScript via `osascript` to get active window information
- **File system monitoring**: Uses `fsnotify` package to watch for file changes
- **Concurrent design**: Runs file watcher in a separate goroutine

## Key Components

- `getActiveWindowInfo()`: Executes AppleScript to retrieve the frontmost application name and window title
- File system watcher: Monitors the project directory for file modifications
- Main goroutine: Blocks indefinitely to keep the program running

## Development Commands

```bash
# Run the application
go run main.go

# Build the application
go build -o ai_system_monitor main.go

# Install dependencies
go mod tidy

# Update dependencies
go mod download
```

## Dependencies

- `github.com/fsnotify/fsnotify v1.9.0` - File system notifications
- `github.com/shirou/gopsutil v3.21.11+incompatible` - System and process utilities (currently unused)

## Platform Requirements

- macOS only (due to AppleScript dependency)
- Accessibility permissions may be required for window title access
- Go 1.23.4+

## Development Notes

- The file watcher is currently hardcoded to monitor `/Users/huyong97/personal/ai_system_monitor`
- The application runs indefinitely until manually terminated
- Error handling is minimal - fatal errors will crash the program