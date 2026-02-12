package win32utils

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// ShowWindow command constants
	SW_HIDE            = 0
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

// GetConsoleWindow retrieves the window handle used by the console associated with the calling process.
func GetConsoleWindow() (windows.HWND, error) {
	r1, _, _ := Kernel32.NewProc("GetConsoleWindow").Call()
	if r1 == 0 {
		return 0, errors.New("no console window found")
	}
	return windows.HWND(r1), nil
}

// ShowConsole shows the console window.
func ShowConsole() error {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		return err
	}
	return ShowWindow(hwnd, SW_SHOW)
}

// HideConsole hides the console window.
func HideConsole() error {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		return err
	}
	return ShowWindow(hwnd, SW_HIDE)
}

// ShowWindow wraps the Win32 API ShowWindow function.
// cmdShow: one of the SW_* constants
func ShowWindow(hwnd windows.HWND, cmdShow int32) error {
	r1, _, _ := User32.NewProc("ShowWindow").Call(
		uintptr(hwnd),
		uintptr(cmdShow),
	)
	// ShowWindow returns nonzero if the window was previously visible, zero otherwise
	// We don't treat this as an error condition
	_ = r1
	return nil
}

// IsConsoleVisible checks if the console window is currently visible.
func IsConsoleVisible() (bool, error) {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		return false, err
	}
	
	// Check if window is visible using IsWindowVisible
	r1, _, _ := User32.NewProc("IsWindowVisible").Call(uintptr(hwnd))
	return r1 != 0, nil
}

// ToggleConsole toggles the visibility of the console window.
// Returns the new visibility state.
func ToggleConsole() (bool, error) {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		return false, err
	}
	
	// Check current visibility
	r1, _, _ := User32.NewProc("IsWindowVisible").Call(uintptr(hwnd))
	isVisible := r1 != 0
	
	// Toggle the state
	var cmdShow int32
	if isVisible {
		cmdShow = SW_HIDE
	} else {
		cmdShow = SW_SHOW
	}
	
	err = ShowWindow(hwnd, cmdShow)
	if err != nil {
		return !isVisible, err
	}
	
	return !isVisible, nil
}

// GetConsoleTitle retrieves the title of the console window.
func GetConsoleTitle() (string, error) {
	// Buffer size - Windows console titles are typically limited
	const bufferSize = 1024
	buffer := make([]uint16, bufferSize)
	
	r1, _, _ := Kernel32.NewProc("GetConsoleTitleW").Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(bufferSize),
	)
	
	if r1 == 0 {
		return "", windows.GetLastError()
	}
	
	return windows.UTF16ToString(buffer[:r1]), nil
}

// SetConsoleTitle sets the title of the console window.
func SetConsoleTitle(title string) error {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	
	r1, _, _ := Kernel32.NewProc("SetConsoleTitleW").Call(uintptr(unsafe.Pointer(titlePtr)))
	if r1 == 0 {
		return windows.GetLastError()
	}
	
	return nil
}