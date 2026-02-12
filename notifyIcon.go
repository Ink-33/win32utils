package win32utils

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Shell_NotifyIconW messages.
// https://learn.microsoft.com/windows/win32/api/shellapi/nf-shellapi-shell_notifyiconw
const (
	NIM_ADD        uint32 = 0x00000000
	NIM_MODIFY     uint32 = 0x00000001
	NIM_DELETE     uint32 = 0x00000002
	NIM_SETFOCUS   uint32 = 0x00000003
	NIM_SETVERSION uint32 = 0x00000004
)

// NOTIFYICONDATAW.uFlags
const (
	NIF_MESSAGE  uint32 = 0x00000001
	NIF_ICON     uint32 = 0x00000002
	NIF_TIP      uint32 = 0x00000004
	NIF_STATE    uint32 = 0x00000008
	NIF_INFO     uint32 = 0x00000010
	NIF_GUID     uint32 = 0x00000020
	NIF_REALTIME uint32 = 0x00000040
	NIF_SHOWTIP  uint32 = 0x00000080
)

// NOTIFYICONDATAW.uVersion (NIM_SETVERSION)
const (
	NOTIFYICON_VERSION   uint32 = 3
	NOTIFYICON_VERSION_4 uint32 = 4
)

// NOTIFYICONDATAW is the wide-char version of NOTIFYICONDATA.
// This definition matches the Windows SDK layout for modern Windows.
// Note: uTimeout and uVersion share the same field in the C union; here it is exposed as TimeoutOrVersion.
type NOTIFYICONDATAW struct {
	CbSize           uint32
	HWnd             windows.HWND
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            windows.Handle
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	TimeoutOrVersion uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GUIDItem         windows.GUID
	HBalloonIcon     windows.Handle
}

// ShellNotifyIconW calls the Win32 API Shell_NotifyIconW.
// If lpData.CbSize is 0, it will be initialized to sizeof(NOTIFYICONDATAW).
func ShellNotifyIconW(dwMessage uint32, lpData *NOTIFYICONDATAW) error {
	if lpData != nil && lpData.CbSize == 0 {
		lpData.CbSize = uint32(unsafe.Sizeof(*lpData))
	}

	r1, _, _ := Shell32.NewProc("Shell_NotifyIconW").Call(
		uintptr(dwMessage),
		uintptr(unsafe.Pointer(lpData)),
	)
	if r1 != 0 {
		return nil
	}

	if err := windows.GetLastError(); err != windows.ERROR_SUCCESS {
		return err
	}
	return errors.New("ShellNotifyIconW failed")
}
