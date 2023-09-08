package win32utils

import "golang.org/x/sys/windows"

const GMEM_MOVEABLE uintptr = 0x0002

func GlobalAlloc(flags uint, size uint) (handle windows.Handle, err error) {
	r1, _, _ := Kernel32.NewProc("GlobalAlloc").Call(uintptr(flags), uintptr(size))
	if r1 == 0 {
		return 0, windows.GetLastError()
	}
	return windows.Handle(r1), nil
}

func GlobalLock(hMem windows.Handle) (pointer uintptr, err error) {
	r1, _, _ := Kernel32.NewProc("GlobalLock").Call(uintptr(hMem))
	if r1 == 0 {
		return 0, windows.GetLastError()
	}
	return r1, nil
}

func GlobalUnlock(hMem windows.Handle) (err error) {
	r1, _, _ := Kernel32.NewProc("GlobalUnlock").Call(uintptr(hMem))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}
