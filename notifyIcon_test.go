//go:build windows

package win32utils

import (
	"testing"
	"unsafe"
)

func TestShellNotifyIconW_FillsCbSize(t *testing.T) {
	var data NOTIFYICONDATAW
	if data.CbSize != 0 {
		t.Fatalf("expected initial CbSize to be 0, got %d", data.CbSize)
	}

	// The API call may fail depending on environment (no tray, permissions, etc.).
	// This test focuses on wrapper behavior: auto-filling cbSize.
	_ = ShellNotifyIconW(NIM_DELETE, &data)

	want := uint32(unsafe.Sizeof(data))
	if data.CbSize != want {
		t.Fatalf("CbSize not filled: got %d, want %d", data.CbSize, want)
	}
}

func TestShellNotifyIconW_DoesNotOverrideCbSize(t *testing.T) {
	data := NOTIFYICONDATAW{CbSize: 123}
	_ = ShellNotifyIconW(NIM_DELETE, &data)
	if data.CbSize != 123 {
		t.Fatalf("CbSize was overridden: got %d, want %d", data.CbSize, 123)
	}
}
