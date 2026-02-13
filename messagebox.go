package win32utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_HELP              = 0x00004000
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_RETRYCANCEL       = 0x00000005
	MB_YESNO             = 0x00000004
	MB_YESNOCANCEL       = 0x00000003

	MB_ICONEXCLAMATION = 0x00000030
	MB_ICONWARNING     = MB_ICONEXCLAMATION
	MB_ICONINFORMATION = 0x00000040
	MB_ICONASTERISK    = MB_ICONINFORMATION
	MB_ICONQUESTION    = 0x00000020
	MB_ICONSTOP        = 0x00000010
	MB_ICONERROR       = MB_ICONSTOP
	MB_ICONHAND        = MB_ICONSTOP

	MB_DEFBUTTON1 = 0x00000000
	MB_DEFBUTTON2 = 0x00000100
	MB_DEFBUTTON3 = 0x00000200
	MB_DEFBUTTON4 = 0x00000300

	MB_APPLMODAL   = 0x00000000
	MB_SYSTEMMODAL = 0x00001000
	MB_TASKMODAL   = 0x00002000

	MB_DEFAULT_DESKTOP_ONLY = 0x00020000
	MB_TOPMOST              = 0x00040000
	MB_RIGHT                = 0x00080000
	MB_RTLREADING           = 0x00100000
)

const (
	IDABORT    = 3
	IDCANCEL   = 2
	IDCONTINUE = 11
	IDIGNORE   = 5
	IDNO       = 7
	IDOK       = 1
	IDRETRY    = 4
	IDTRYAGAIN = 10
	IDYES      = 6
)

// RunningByDoubleClick Check if run directly by double-clicking
func RunningByDoubleClick() bool {
	lp := Kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var ids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&ids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}

// MessageBoxW of Win32 API. Check https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw for more detail.
func MessageBoxW(hwnd uintptr, caption, title string, flags uint) int {
	captionPtr, _ := windows.UTF16PtrFromString(caption)
	titlePtr, _ := windows.UTF16PtrFromString(title)
	ret, _, _ := User32.NewProc("MessageBoxW").Call(
		hwnd,
		uintptr(unsafe.Pointer(captionPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags))

	return int(ret)
}

// GetConsoleWindows retrieves the window handle used by the console associated with the calling process.
func GetConsoleWindows() (hwnd uintptr) {
	hwnd, _, _ = Kernel32.NewProc("GetConsoleWindow").Call()
	return
}

// ToHighDPI tries to raise DPI awareness context to DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED
func ToHighDPI() {
	systemAware := ^uintptr(2) + 1
	unawareGDIScaled := ^uintptr(5) + 1
	proc := User32.NewProc("SetThreadDpiAwarenessContext")
	if proc.Find() != nil {
		return
	}
	for i := unawareGDIScaled; i <= systemAware; i++ {
		_, _, _ = User32.NewProc("SetThreadDpiAwarenessContext").Call(i)
	}
}

// ToHighDPIEx tries to raise DPI awareness context to DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED for the process
func ToHighDPIEx() {
	systemAware := ^uintptr(2) + 1
	unawareGDIScaled := ^uintptr(5) + 1
	proc := User32.NewProc("SetProcessDpiAwarenessContext")
	if proc.Find() != nil {
		return
	}
	for i := unawareGDIScaled; i <= systemAware; i++ {
		_, _, _ = User32.NewProc("SetProcessDpiAwarenessContext").Call(i)
	}
}
