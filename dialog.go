//go:build windows

package win32utils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/windows"
)

// TwoTextInputDialog displays a modal dialog with two text input fields.
// Returns (text1, text2, cancelled, error).
func TwoTextInputDialog(title, label1, label2 string, defaultValue1, defaultValue2 string) (string, string, bool, error) {
	hInstance, err := getModuleHandleCurrentProcess()
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get module handle: %w", err)
	}

	// Create a top-level window (not message-only) for the dialog
	// Use DPI-scaled dimensions
	dialogWidth := ScaleSize(380)
	dialogHeight := ScaleSize(260) // Increased for larger font
	
	dialogHWnd, err := CreateWindowExW(
		WindowExStyle{}.With(WS_EX_DLGMODALFRAME | WS_EX_TOPMOST),
		"dialog_input",
		title,
		WindowStyle{}.With(WS_OVERLAPPED | WS_SYSMENU | WS_CAPTION),
		ScaleX(100), ScaleY(100), dialogWidth, dialogHeight,
		0,
		0,
		hInstance,
		0,
	)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to create dialog: %w", err)
	}

	// Create controls - all with DPI scaling
	label1Hwnd, _ := CreateWindowExW(
		WindowExStyle{},
		"STATIC",
		label1,
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD),
		ScaleX(10), ScaleY(10), ScaleX(100), ScaleY(20),
		dialogHWnd,
		0,
		hInstance,
		0,
	)

	edit1Hwnd, _ := CreateWindowExW(
		WindowExStyle{}.With(WS_EX_CLIENTEDGE),
		"EDIT",
		defaultValue1,
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD | WS_TABSTOP),
		ScaleX(120), ScaleY(10), ScaleX(245), ScaleY(26),
		dialogHWnd,
		windows.Handle(1001),
		hInstance,
		0,
	)

	label2Hwnd, _ := CreateWindowExW(
		WindowExStyle{},
		"STATIC",
		label2,
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD),
		ScaleX(10), ScaleY(50), ScaleX(100), ScaleY(20),
		dialogHWnd,
		0,
		hInstance,
		0,
	)

	edit2Hwnd, _ := CreateWindowExW(
		WindowExStyle{}.With(WS_EX_CLIENTEDGE),
		"EDIT",
		defaultValue2,
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD | WS_TABSTOP),
		ScaleX(120), ScaleY(50), ScaleX(245), ScaleY(26),
		dialogHWnd,
		windows.Handle(1002),
		hInstance,
		0,
	)

	okHwnd, _ := CreateWindowExW(
		WindowExStyle{},
		"BUTTON",
		"OK",
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD),
		ScaleX(120), ScaleY(110), ScaleX(110), ScaleY(30),
		dialogHWnd,
		windows.Handle(IDOK),
		hInstance,
		0,
	)

	cancelHwnd, _ := CreateWindowExW(
		WindowExStyle{},
		"BUTTON",
		"Cancel",
		WindowStyle{}.With(WS_VISIBLE | WS_CHILD),
		ScaleX(240), ScaleY(110), ScaleX(110), ScaleY(30),
		dialogHWnd,
		windows.Handle(IDCANCEL),
		hInstance,
		0,
	)
	_ = cancelHwnd

	// Create and apply font to all controls
	// Font size: 11pt (scaled for DPI)
	fontHeight := ScaleSize(-14) // negative value for character height (11pt ≈ 14 pixels at 96 DPI)
	uiFont, fontErr := CreateFontW(
		fontHeight,
		0,                   // width (0 = auto)
		0,                   // escapement
		0,                   // orientation
		FW_NORMAL,           // weight
		false, false, false, // italic, underline, strikeOut
		DEFAULT_CHARSET,
		OUT_DEFAULT_PRECIS,
		CLIP_DEFAULT_PRECIS,
		PROOF_QUALITY,
		FF_DONTCARE,
		"Segoe UI", // Modern Windows font
	)
	if fontErr == nil && uiFont != 0 {
		// Apply font to all controls
		SetWindowFontW(label1Hwnd, uiFont, false)
		SetWindowFontW(edit1Hwnd, uiFont, false)
		SetWindowFontW(label2Hwnd, uiFont, false)
		SetWindowFontW(edit2Hwnd, uiFont, false)
		SetWindowFontW(okHwnd, uiFont, false)
		SetWindowFontW(cancelHwnd, uiFont, false)
	}

	// State variables
	var result1, result2 string
	cancelled := false
	var done int32 = 0

	oldProc := setDialogWndProc(dialogHWnd, func(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
		switch msg {
		case WM_COMMAND:
			id := int32(wParam & 0xFFFF)
			if id == IDOK {
				text1, _ := GetWindowTextW(edit1Hwnd)
				text2, _ := GetWindowTextW(edit2Hwnd)
				result1 = text1
				result2 = text2
				atomic.StoreInt32(&done, 1)
				PostMessageW(hwnd, WM_CLOSE, 0, 0)
				return 0
			} else if id == IDCANCEL {
				cancelled = true
				atomic.StoreInt32(&done, 1)
				PostMessageW(hwnd, WM_CLOSE, 0, 0)
				return 0
			}

		case WM_CLOSE:
			DestroyWindow(hwnd)
			return 0

		case WM_DESTROY:
			atomic.StoreInt32(&done, 1)
			return 0
		}
		return DefWindowProcW(hwnd, msg, wParam, lParam)
	})
	defer setDialogWndProc(dialogHWnd, oldProc)

	// Show dialog and set focus
	ShowWindowW(dialogHWnd, 5) // SW_SHOW
	SetFocus(okHwnd)

	// Modal message loop - standard GetMessageW/DispatchMessageW pattern
	const timeoutMs = 30000 // 30 second timeout
	startTick := GetTickCount()
	
	for atomic.LoadInt32(&done) == 0 {
		// Check timeout
		elapsed := GetTickCount() - startTick
		if elapsed > timeoutMs {
			break
		}

		// Get next message from queue - use 0 to get messages from all windows
		// DispatchMessageW will automatically route to the correct window procedure
		var msg MSG
		ret, _ := GetMessageW(&msg, 0, 0, 0)
		
		if ret == 0 {
			// WM_QUIT received
			break
		}
		if ret == -1 {
			// Error
			break
		}
		
		// Translate key messages (like Alt+key combinations)
		TranslateMessage(&msg)
		// Dispatch message to the appropriate window procedure
		DispatchMessageW(&msg)
	}

	// Cleanup - ensure window is destroyed
	if IsWindowW(dialogHWnd) {
		DestroyWindow(dialogHWnd)
	}
	
	// Clean up the window procedure mapping
	setDialogWndProc(dialogHWnd, nil)

	return result1, result2, cancelled, nil
}

var (
	dialogWndProcMu sync.RWMutex
	dialogWndProcs  = map[windows.HWND]WndProc{}
)

func setDialogWndProc(hwnd windows.HWND, proc WndProc) WndProc {
	dialogWndProcMu.Lock()
	defer dialogWndProcMu.Unlock()
	old := dialogWndProcs[hwnd]
	if proc == nil {
		delete(dialogWndProcs, hwnd)
	} else {
		dialogWndProcs[hwnd] = proc
	}
	return old
}

func getDialogWndProc(hwnd windows.HWND) (WndProc, bool) {
	dialogWndProcMu.RLock()
	defer dialogWndProcMu.RUnlock()
	proc, ok := dialogWndProcs[hwnd]
	return proc, ok
}

func dialogGlobalWndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if proc, ok := getDialogWndProc(hwnd); ok {
		return proc(hwnd, msg, wParam, lParam)
	}
	return DefWindowProcW(hwnd, msg, wParam, lParam)
}

func init() {
	// Register a single window class for all dialog windows we create
	dialogClassOnce.Do(func() {
		hInstance, _ := getModuleHandleCurrentProcess()
		classNamePtr, _ := windows.UTF16PtrFromString("dialog_input")

		wcx := &WNDCLASSEXW{
			LpfnWndProc:   syscall.NewCallback(dialogGlobalWndProc),
			HInstance:     hInstance,
			LpszClassName: classNamePtr,
			HbrBackground: windows.Handle(uintptr(5) + 1), // COLOR_BTNFACE + 1
		}
		_, _ = registerClassExW(wcx)
	})
}

var (
	dialogClassOnce sync.Once
)
