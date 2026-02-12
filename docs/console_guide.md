# Console Management Guide

This guide explains how to use the console management features added to win32utils.

## Overview

The console management functionality allows you to control the visibility and properties of your application's console window on Windows systems.

## Available Functions

### Basic Console Operations

```go
// Get the console window handle
hwnd, err := win32utils.GetConsoleWindow()
if err != nil {
    log.Printf("No console window available: %v", err)
}

// Show the console window
err = win32utils.ShowConsole()
if err != nil {
    log.Printf("Failed to show console: %v", err)
}

// Hide the console window
err = win32utils.HideConsole()
if err != nil {
    log.Printf("Failed to hide console: %v", err)
}

// Toggle console visibility (returns new state)
isVisible, err := win32utils.ToggleConsole()
if err != nil {
    log.Printf("Failed to toggle console: %v", err)
} else {
    fmt.Printf("Console is now %s\n", map[bool]string{true: "visible", false: "hidden"}[isVisible])
}

// Check if console is currently visible
visible, err := win32utils.IsConsoleVisible()
if err != nil {
    log.Printf("Failed to check visibility: %v", err)
}
```

### Console Title Management

```go
// Get current console title
title, err := win32utils.GetConsoleTitle()
if err != nil {
    log.Printf("Failed to get console title: %v", err)
} else {
    fmt.Printf("Current title: %s\n", title)
}

// Set console title
err = win32utils.SetConsoleTitle("My Application Console")
if err != nil {
    log.Printf("Failed to set console title: %v", err)
}
```

### Advanced Window Control

```go
// Using the generic ShowWindow function with specific commands
hwnd, _ := win32utils.GetConsoleWindow()

// Hide window
win32utils.ShowWindow(hwnd, win32utils.SW_HIDE)

// Show normally
win32utils.ShowWindow(hwnd, win32utils.SW_SHOWNORMAL)

// Show maximized
win32utils.ShowWindow(hwnd, win32utils.SW_SHOWMAXIMIZED)

// Show minimized
win32utils.ShowWindow(hwnd, win32utils.SW_SHOWMINIMIZED)
```

## Available ShowWindow Commands

| Constant | Value | Description |
|----------|-------|-------------|
| `SW_HIDE` | 0 | Hides the window and activates another window |
| `SW_SHOWNORMAL` | 1 | Activates and displays a window |
| `SW_SHOWMINIMIZED` | 2 | Activates the window and displays it as an icon |
| `SW_SHOWMAXIMIZED` | 3 | Activates the window and displays it as a maximized window |
| `SW_SHOWNOACTIVATE` | 4 | Displays a window in its most recent size and position |
| `SW_SHOW` | 5 | Activates the window and displays it in its current size and position |
| `SW_MINIMIZE` | 6 | Minimizes the specified window and activates the next top-level window |
| `SW_SHOWMINNOACTIVE` | 7 | Displays the window as an icon |
| `SW_SHOWNA` | 8 | Displays the window in its current size and position |
| `SW_RESTORE` | 9 | Activates and displays the window |
| `SW_SHOWDEFAULT` | 10 | Sets the show state based on the SW_ value specified in the STARTUPINFO structure |
| `SW_FORCEMINIMIZE` | 11 | Minimizes a window, even if the thread that owns the window is not responding |

## Example Usage

Here's a complete example demonstrating various console management operations:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "repo.smlk.org/win32utils"
)

func main() {
    // Check if we have a console
    hwnd, err := win32utils.GetConsoleWindow()
    if err != nil {
        log.Fatal("No console window available")
    }
    fmt.Printf("Console window handle: %v\n", hwnd)
    
    // Save original title
    originalTitle, _ := win32utils.GetConsoleTitle()
    
    // Demonstrate title management
    win32utils.SetConsoleTitle("Demo Application")
    fmt.Println("Console title updated")
    
    // Demonstrate visibility control
    fmt.Println("Hiding console for 3 seconds...")
    win32utils.HideConsole()
    time.Sleep(3 * time.Second)
    
    fmt.Println("Showing console again...")
    win32utils.ShowConsole()
    
    // Demonstrate toggle
    fmt.Println("Toggling console visibility 3 times...")
    for i := 0; i < 3; i++ {
        newState, _ := win32utils.ToggleConsole()
        fmt.Printf("Toggle %d: Console is now %s\n", i+1, 
            map[bool]string{true: "visible", false: "hidden"}[newState])
        time.Sleep(1 * time.Second)
    }
    
    // Restore original title
    win32utils.SetConsoleTitle(originalTitle)
    fmt.Println("Demo completed!")
}
```

## Important Notes

1. **Console Availability**: These functions only work when your application has an associated console window. GUI applications without consoles will return errors.

2. **Thread Safety**: Most console operations are safe to call from any thread, but visibility changes should ideally be made from the main thread.

3. **Error Handling**: Always check for errors, especially when dealing with console operations in production code.

4. **Resource Cleanup**: The functions automatically manage Windows resources and don't require manual cleanup.

5. **Testing**: Run the provided test suite (`console_test.go`) to verify functionality on your system.

## Running Tests

To test the console functionality:

```bash
go test console_test.go console.go dll.go winbase.go -v
```

To run benchmarks:

```bash
go test console_test.go console.go dll.go winbase.go -bench=. -benchmem
```

## See Also

- [API Reference](API.md) - Complete API documentation
- [Examples](../examples/) - More usage examples
- [Changelog](../CHANGELOG.md) - Version history