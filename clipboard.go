package win32utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const CF_TEXT uintptr = 1
const CF_UNICODETEXT uintptr = 13
const CF_LOCALE uintptr = 16

func OpenClipboard(hwnd windows.HWND) error {
	r1, _, _ := User32.NewProc("OpenClipboard").Call(uintptr(hwnd))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}
func CloseClipboard() error {
	r1, _, _ := User32.NewProc("CloseClipboard").Call()
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}
func EmptyClipboard() error {
	r1, _, _ := User32.NewProc("EmptyClipboard").Call()
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}
func SetClipboardText(text string) (handle windows.Handle, err error) {
	proc := User32.NewProc("SetClipboardData")
	u16text, err := windows.UTF16FromString(text)
	if err != nil {
		return 0, err
	}

	h, err := GlobalAlloc(uint(GMEM_MOVEABLE),
		uint(len(u16text))*uint(unsafe.Sizeof(u16text[0])))
	if err != nil {
		return 0, err
	}

	p, err := GlobalLock(h)
	if err != nil {
		return 0, err
	}

	dst := unsafe.Slice((*uint16)(unsafe.Pointer(p)), len(u16text))
	copy(dst, u16text)
	
	err = GlobalUnlock(h)
	if err != nil {
		return 0, err
	}

	r1, _, _ := proc.Call(CF_UNICODETEXT, uintptr(h))
	if r1 == 0 {
		return 0, windows.GetLastError()
	}

	return windows.Handle(r1), nil
}

func GetClipboardDataText() (string, error) {
	r1, _, _ := User32.NewProc("GetClipboardData").Call(CF_UNICODETEXT)
	if r1 == 0 {
		return "", windows.GetLastError()
	}

	p, err := GlobalLock(windows.Handle(r1))
	if err != nil {
		return "", err
	}
	defer GlobalUnlock(windows.Handle(r1))

	return windows.UTF16PtrToString((*uint16)(unsafe.Pointer(p))), nil
}

func SetText(text string) error {
	err := OpenClipboard(windows.GetShellWindow())
	if err != nil {
		return err
	}

	err = EmptyClipboard()
	if err != nil {
		return err
	}

	_, err = SetClipboardText(text)
	if err != nil {
		return err
	}

	err = CloseClipboard()
	if err != nil {
		return err
	}

	return nil
}
