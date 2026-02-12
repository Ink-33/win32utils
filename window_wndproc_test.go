//go:build windows

package win32utils

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestGlobalWndProcDispatchAndCleanup(t *testing.T) {
	hwnd := windows.HWND(0x1234)
	called := 0

	setWndProc(hwnd, func(h windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
		called++
		if h != hwnd {
			t.Fatalf("unexpected hwnd: got %v want %v", h, hwnd)
		}
		if msg == 0x0400 {
			return 99
		}
		return 0
	})

	ret := globalWndProc(hwnd, 0x0400, 1, 2)
	if ret != 99 {
		t.Fatalf("unexpected return value: got %d want %d", ret, 99)
	}
	if called != 1 {
		t.Fatalf("handler call count mismatch: got %d want %d", called, 1)
	}

	_ = globalWndProc(hwnd, WM_NCDESTROY, 0, 0)
	if called != 2 {
		t.Fatalf("handler should be called for WM_NCDESTROY")
	}

	if _, ok := getWndProc(hwnd); ok {
		t.Fatalf("handler should be removed on WM_NCDESTROY")
	}
}
