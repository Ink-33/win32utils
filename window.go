package win32utils

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// HWND_MESSAGE is the parent handle for message-only windows.
// https://learn.microsoft.com/en-us/windows/win32/winmsg/window-features#message-only-windows
const HWND_MESSAGE windows.HWND = windows.HWND(^uintptr(2))

const (
	WM_DESTROY       uint32 = 0x0002
	WM_NCDESTROY     uint32 = 0x0082
	WM_CLOSE         uint32 = 0x0010
	WM_COMMAND       uint32 = 0x0111
	WM_GETTEXT       uint32 = 0x000D
	WM_SETTEXT       uint32 = 0x000C
	WM_SETFONT       uint32 = 0x0030
	WM_USER          uint32 = 0x0400
	WM_LBUTTONDOWN   uint32 = 0x0201
	WM_RBUTTONDOWN   uint32 = 0x0204
	WM_LBUTTONDBLCLK uint32 = 0x0203
	IDOK             int32  = 1
	IDCANCEL         int32  = 2
)

// WndProc is the window procedure callback.
// Return value meaning depends on msg.
type WndProc func(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr

var (
	wndProcMu     sync.RWMutex
	wndProcByHWND = map[windows.HWND]WndProc{}

	globalWndProcOnce sync.Once
	globalWndProcPtr  uintptr
)

func ensureGlobalWndProc() uintptr {
	globalWndProcOnce.Do(func() {
		globalWndProcPtr = syscall.NewCallback(globalWndProc)
	})
	return globalWndProcPtr
}

func setWndProc(hwnd windows.HWND, proc WndProc) {
	wndProcMu.Lock()
	defer wndProcMu.Unlock()
	if proc == nil {
		delete(wndProcByHWND, hwnd)
		return
	}
	wndProcByHWND[hwnd] = proc
}

func deleteWndProc(hwnd windows.HWND) {
	wndProcMu.Lock()
	defer wndProcMu.Unlock()
	delete(wndProcByHWND, hwnd)
}

func getWndProc(hwnd windows.HWND) (WndProc, bool) {
	wndProcMu.RLock()
	defer wndProcMu.RUnlock()
	proc, ok := wndProcByHWND[hwnd]
	return proc, ok
}

// DefWindowProcW calls the Win32 API DefWindowProcW.
func DefWindowProcW(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	r1, _, _ := User32.NewProc("DefWindowProcW").Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	return r1
}

// globalWndProc dispatches to the Go handler registered for hwnd.
// It also cleans up handler state on WM_NCDESTROY.
func globalWndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if proc, ok := getWndProc(hwnd); ok {
		ret := proc(hwnd, msg, wParam, lParam)
		if msg == WM_NCDESTROY {
			deleteWndProc(hwnd)
		}
		return ret
	}
	return DefWindowProcW(hwnd, msg, wParam, lParam)
}

// WNDCLASSEXW describes a window class.
// https://learn.microsoft.com/windows/win32/api/winuser/ns-winuser-wndclassexw
type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

func getModuleHandleCurrentProcess() (windows.Handle, error) {
	r1, _, _ := Kernel32.NewProc("GetModuleHandleW").Call(0)
	if r1 == 0 {
		return 0, windows.GetLastError()
	}
	return windows.Handle(r1), nil
}

func registerClassExW(wcx *WNDCLASSEXW) (uint16, error) {
	if wcx == nil {
		return 0, errors.New("registerClassExW: nil WNDCLASSEXW")
	}
	if wcx.CbSize == 0 {
		wcx.CbSize = uint32(unsafe.Sizeof(*wcx))
	}

	r1, _, _ := User32.NewProc("RegisterClassExW").Call(uintptr(unsafe.Pointer(wcx)))
	if r1 != 0 {
		return uint16(r1), nil
	}

	err := windows.GetLastError()
	// ERROR_CLASS_ALREADY_EXISTS (1410)
	if errno, ok := err.(syscall.Errno); ok && errno == 1410 {
		// Treat as success so callers can safely call registration multiple times.
		return 1, nil
	}
	return 0, err
}

// CreateMessageOnlyWindow creates a message-only window (parent = HWND_MESSAGE) and associates a Go WndProc.
//
// Notes:
// - The window is not visible.
// - The returned hwnd is only valid on the thread that created it for message processing.
// - To process messages, call MessageLoop (or use your own GetMessage loop).
func CreateMessageOnlyWindow(className, windowName string, proc WndProc) (windows.HWND, error) {
	if className == "" {
		return 0, errors.New("CreateMessageOnlyWindow: className is empty")
	}

	hInstance, err := getModuleHandleCurrentProcess()
	if err != nil {
		return 0, err
	}

	classNamePtr, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return 0, err
	}

	_, err = registerClassExW(&WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		LpfnWndProc:   ensureGlobalWndProc(),
		HInstance:     hInstance,
		LpszClassName: classNamePtr,
	})
	if err != nil {
		return 0, err
	}

	hwnd, err := CreateWindowExW(
		WindowExStyle{},
		className,
		windowName,
		WindowStyle{},
		0, 0, 0, 0,
		HWND_MESSAGE,
		0,
		hInstance,
		0,
	)
	if err != nil {
		return 0, err
	}

	if proc != nil {
		setWndProc(hwnd, proc)
	}
	return hwnd, nil
}

type POINT struct {
	X int32
	Y int32
}

type MSG struct {
	HWnd     windows.HWND
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       POINT
	LPrivate uint32
}

// GetMessageW wraps the Win32 API GetMessageW.
// Return values:
// - >0: message retrieved
// - 0: WM_QUIT received
// - <0: error
func GetMessageW(msg *MSG, hwnd windows.HWND, msgFilterMin, msgFilterMax uint32) (int32, error) {
	r1, _, _ := User32.NewProc("GetMessageW").Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	ret := int32(r1)
	if ret == -1 {
		return ret, windows.GetLastError()
	}
	return ret, nil
}

// TranslateMessage wraps the Win32 API TranslateMessage.
func TranslateMessage(msg *MSG) {
	_, _, _ = User32.NewProc("TranslateMessage").Call(uintptr(unsafe.Pointer(msg)))
}

// DispatchMessageW wraps the Win32 API DispatchMessageW.
func DispatchMessageW(msg *MSG) uintptr {
	r1, _, _ := User32.NewProc("DispatchMessageW").Call(uintptr(unsafe.Pointer(msg)))
	return r1
}

// PostQuitMessage wraps the Win32 API PostQuitMessage.
func PostQuitMessage(exitCode int32) {
	_, _, _ = User32.NewProc("PostQuitMessage").Call(uintptr(exitCode))
}

// DestroyWindow wraps the Win32 API DestroyWindow.
func DestroyWindow(hwnd windows.HWND) error {
	r1, _, _ := User32.NewProc("DestroyWindow").Call(uintptr(hwnd))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// MessageLoop runs the standard GetMessage/TranslateMessage/DispatchMessage loop.
// It returns the WM_QUIT exit code.
func MessageLoop() (int32, error) {
	var msg MSG
	for {
		ret, err := GetMessageW(&msg, 0, 0, 0)
		if ret == -1 {
			return 0, err
		}
		if ret == 0 {
			return int32(msg.WParam), nil
		}
		TranslateMessage(&msg)
		DispatchMessageW(&msg)
	}
}

// CreateCurrentProcessWindow creates a message-only window owned by the current process.
//
// This is useful as an HWND target for APIs that require a window handle (e.g. tray icon callbacks),
// without showing an actual top-level window.
//
// It registers the window class (if needed) using DefWindowProcW as the window procedure.
func CreateCurrentProcessWindow(className, windowName string) (windows.HWND, error) {
	if className == "" {
		return 0, errors.New("CreateCurrentProcessWindow: className is empty")
	}

	hInstance, err := getModuleHandleCurrentProcess()
	if err != nil {
		return 0, err
	}

	classNamePtr, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return 0, err
	}

	defProc := User32.NewProc("DefWindowProcW")
	if err := defProc.Find(); err != nil {
		return 0, err
	}

	_, err = registerClassExW(&WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		LpfnWndProc:   defProc.Addr(),
		HInstance:     hInstance,
		LpszClassName: classNamePtr,
	})
	if err != nil {
		return 0, err
	}

	// Create as a message-only window. For such windows, style/exStyle should be 0.
	return CreateWindowExW(
		WindowExStyle{},
		className,
		windowName,
		WindowStyle{},
		0, 0, 0, 0,
		HWND_MESSAGE,
		0,
		hInstance,
		0,
	)
}

// CreateWindowInstance creates a message-only window instance owned by the current process.
//
// This is a minimal helper that returns a valid HWND for the process. If your app needs to
// receive window messages continuously, run MessageLoop() on the same thread.
func CreateWindowInstance() (windows.HWND, error) {
	return CreateMessageOnlyWindow(
		"win32utils.MessageOnlyWindow",
		"win32utils",
		func(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
			switch msg {
			case WM_DESTROY:
				PostQuitMessage(0)
				return 0
			default:
				return DefWindowProcW(hwnd, msg, wParam, lParam)
			}
		},
	)
}

// WindowStyleBits represents dwStyle flags for CreateWindowExW.
// https://learn.microsoft.com/windows/win32/winmsg/window-styles
type WindowStyleBits uint32

// WindowExStyleBits represents dwExStyle flags for CreateWindowExW.
// https://learn.microsoft.com/windows/win32/winmsg/extended-window-styles
type WindowExStyleBits uint32

// WindowStyle is a small helper wrapper around dwStyle.
type WindowStyle struct {
	Bits WindowStyleBits
}

func (s WindowStyle) Uint32() uint32 { return uint32(s.Bits) }
func (s WindowStyle) With(bits WindowStyleBits) WindowStyle {
	return WindowStyle{Bits: s.Bits | bits}
}

func (s WindowStyle) Without(bits WindowStyleBits) WindowStyle {
	return WindowStyle{Bits: s.Bits &^ bits}
}
func (s WindowStyle) Has(bits WindowStyleBits) bool { return s.Bits&bits == bits }

// WindowExStyle is a small helper wrapper around dwExStyle.
type WindowExStyle struct {
	Bits WindowExStyleBits
}

func (s WindowExStyle) Uint32() uint32 { return uint32(s.Bits) }
func (s WindowExStyle) With(bits WindowExStyleBits) WindowExStyle {
	return WindowExStyle{Bits: s.Bits | bits}
}

func (s WindowExStyle) Without(bits WindowExStyleBits) WindowExStyle {
	return WindowExStyle{Bits: s.Bits &^ bits}
}
func (s WindowExStyle) Has(bits WindowExStyleBits) bool { return s.Bits&bits == bits }

// Common window styles (WS_*).
const (
	WS_OVERLAPPED   WindowStyleBits = 0x00000000
	WS_POPUP        WindowStyleBits = 0x80000000
	WS_CHILD        WindowStyleBits = 0x40000000
	WS_MINIMIZE     WindowStyleBits = 0x20000000
	WS_VISIBLE      WindowStyleBits = 0x10000000
	WS_DISABLED     WindowStyleBits = 0x08000000
	WS_CLIPSIBLINGS WindowStyleBits = 0x04000000
	WS_CLIPCHILDREN WindowStyleBits = 0x02000000
	WS_MAXIMIZE     WindowStyleBits = 0x01000000
	WS_CAPTION      WindowStyleBits = 0x00C00000
	WS_BORDER       WindowStyleBits = 0x00800000
	WS_DLGFRAME     WindowStyleBits = 0x00400000
	WS_VSCROLL      WindowStyleBits = 0x00200000
	WS_HSCROLL      WindowStyleBits = 0x00100000
	WS_SYSMENU      WindowStyleBits = 0x00080000
	WS_THICKFRAME   WindowStyleBits = 0x00040000
	WS_GROUP        WindowStyleBits = 0x00020000
	WS_TABSTOP      WindowStyleBits = 0x00010000
	WS_MINIMIZEBOX  WindowStyleBits = 0x00020000
	WS_MAXIMIZEBOX  WindowStyleBits = 0x00010000
)

// Common composite window styles.
const (
	WS_OVERLAPPEDWINDOW WindowStyleBits = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_POPUPWINDOW      WindowStyleBits = WS_POPUP | WS_BORDER | WS_SYSMENU
	WS_CHILDWINDOW      WindowStyleBits = WS_CHILD
)

// Common extended window styles (WS_EX_*).
const (
	WS_EX_DLGMODALFRAME   WindowExStyleBits = 0x00000001
	WS_EX_NOPARENTNOTIFY  WindowExStyleBits = 0x00000004
	WS_EX_TOPMOST         WindowExStyleBits = 0x00000008
	WS_EX_ACCEPTFILES     WindowExStyleBits = 0x00000010
	WS_EX_TRANSPARENT     WindowExStyleBits = 0x00000020
	WS_EX_MDICHILD        WindowExStyleBits = 0x00000040
	WS_EX_TOOLWINDOW      WindowExStyleBits = 0x00000080
	WS_EX_WINDOWEDGE      WindowExStyleBits = 0x00000100
	WS_EX_CLIENTEDGE      WindowExStyleBits = 0x00000200
	WS_EX_CONTEXTHELP     WindowExStyleBits = 0x00000400
	WS_EX_RIGHT           WindowExStyleBits = 0x00001000
	WS_EX_LEFT            WindowExStyleBits = 0x00000000
	WS_EX_RTLREADING      WindowExStyleBits = 0x00002000
	WS_EX_LTRREADING      WindowExStyleBits = 0x00000000
	WS_EX_LEFTSCROLLBAR   WindowExStyleBits = 0x00004000
	WS_EX_RIGHTSCROLLBAR  WindowExStyleBits = 0x00000000
	WS_EX_CONTROLPARENT   WindowExStyleBits = 0x00010000
	WS_EX_STATICEDGE      WindowExStyleBits = 0x00020000
	WS_EX_APPWINDOW       WindowExStyleBits = 0x00040000
	WS_EX_LAYERED         WindowExStyleBits = 0x00080000
	WS_EX_NOINHERITLAYOUT WindowExStyleBits = 0x00100000
	WS_EX_LAYOUTRTL       WindowExStyleBits = 0x00400000
	WS_EX_COMPOSITED      WindowExStyleBits = 0x02000000
	WS_EX_NOACTIVATE      WindowExStyleBits = 0x08000000
)

// CreateWindowExW wraps the Win32 API CreateWindowExW.
// https://learn.microsoft.com/windows/win32/api/winuser/nf-winuser-createwindowexw
//
// Parameters:
// - className: registered window class name
// - windowName: initial window title
// - parent/menu/instance can be 0 if not used.
// - param is the lpParam passed to WM_NCCREATE/WM_CREATE.
func CreateWindowExW(
	exStyle WindowExStyle,
	className string,
	windowName string,
	style WindowStyle,
	x, y, width, height int32,
	parent windows.HWND,
	menu windows.Handle,
	instance windows.Handle,
	param uintptr,
) (windows.HWND, error) {
	classNamePtr, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return 0, err
	}
	windowNamePtr, err := windows.UTF16PtrFromString(windowName)
	if err != nil {
		return 0, err
	}

	r1, _, _ := User32.NewProc("CreateWindowExW").Call(
		uintptr(exStyle.Uint32()),
		uintptr(unsafe.Pointer(classNamePtr)),
		uintptr(unsafe.Pointer(windowNamePtr)),
		uintptr(style.Uint32()),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		param,
	)
	if r1 == 0 {
		return 0, windows.GetLastError()
	}
	return windows.HWND(r1), nil
}

// Menu-related constants and APIs
const (
	MFT_STRING       uint32 = 0x00000000
	MF_ENABLED       uint32 = 0x00000000
	MF_GRAYED        uint32 = 0x00000001
	MF_DISABLED      uint32 = 0x00000002
	MF_SEPARATOR     uint32 = 0x00000800
	MFT_SEPARATOR    uint32 = 0x00000800
	TPM_LEFTALIGN    uint32 = 0x0000
	TPM_RIGHTALIGN   uint32 = 0x0008
	TPM_TOPALIGN     uint32 = 0x0000
	TPM_VCENTERALIGN uint32 = 0x0010
	TPM_BOTTOMALIGN  uint32 = 0x0020
	TPM_LEFTBUTTON   uint32 = 0x0000
	TPM_RIGHTBUTTON  uint32 = 0x0002
	TPM_RETURNCMD    uint32 = 0x0100
)

// CreatePopupMenu wraps the Win32 API CreatePopupMenu.
func CreatePopupMenu() (windows.Handle, error) {
	r1, _, _ := User32.NewProc("CreatePopupMenu").Call()
	if r1 == 0 {
		return 0, windows.GetLastError()
	}
	return windows.Handle(r1), nil
}

// AppendMenuW wraps the Win32 API AppendMenuW.
func AppendMenuW(hMenu windows.Handle, flags, id uint32, label string) error {
	labelPtr, err := windows.UTF16PtrFromString(label)
	if err != nil {
		return err
	}
	r1, _, _ := User32.NewProc("AppendMenuW").Call(
		uintptr(hMenu),
		uintptr(flags),
		uintptr(id),
		uintptr(unsafe.Pointer(labelPtr)),
	)
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// TrackPopupMenuEx wraps the Win32 API TrackPopupMenuEx.
// Returns the selected menu item ID, or 0 if cancelled.
func TrackPopupMenuEx(hMenu windows.Handle, flags uint32, x, y int32, hwnd windows.HWND) (int32, error) {
	r1, _, _ := User32.NewProc("TrackPopupMenuEx").Call(
		uintptr(hMenu),
		uintptr(flags),
		uintptr(x),
		uintptr(y),
		uintptr(hwnd),
		0, // lpTPMParams (optional)
	)
	return int32(r1), nil
}

// DestroyMenu wraps the Win32 API DestroyMenu.
func DestroyMenu(hMenu windows.Handle) error {
	r1, _, _ := User32.NewProc("DestroyMenu").Call(uintptr(hMenu))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// GetWindowTextW wraps the Win32 API GetWindowTextW.
func GetWindowTextW(hwnd windows.HWND) (string, error) {
	buf := make([]uint16, 256)
	r1, _, _ := User32.NewProc("GetWindowTextW").Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r1 == 0 {
		return "", windows.GetLastError()
	}
	return windows.UTF16ToString(buf[:r1]), nil
}

// SetWindowTextW wraps the Win32 API SetWindowTextW.
func SetWindowTextW(hwnd windows.HWND, text string) error {
	textPtr, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return err
	}
	r1, _, _ := User32.NewProc("SetWindowTextW").Call(uintptr(hwnd), uintptr(unsafe.Pointer(textPtr)))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// SetForegroundWindow wraps the Win32 API SetForegroundWindow.
func SetForegroundWindow(hwnd windows.HWND) error {
	r1, _, _ := User32.NewProc("SetForegroundWindow").Call(uintptr(hwnd))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// SetFocus wraps the Win32 API SetFocus.
func SetFocus(hwnd windows.HWND) windows.HWND {
	r1, _, _ := User32.NewProc("SetFocus").Call(uintptr(hwnd))
	return windows.HWND(r1)
}

// PostMessageW wraps the Win32 API PostMessageW.
func PostMessageW(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) error {
	r1, _, _ := User32.NewProc("PostMessageW").Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// ShowWindowW wraps the Win32 API ShowWindow.
func ShowWindowW(hwnd windows.HWND, cmdShow int32) error {
	r1, _, _ := User32.NewProc("ShowWindow").Call(uintptr(hwnd), uintptr(cmdShow))
	if r1 == 0 {
		return windows.GetLastError()
	}
	return nil
}

// IsWindowW wraps the Win32 API IsWindow.
func IsWindowW(hwnd windows.HWND) bool {
	r1, _, _ := User32.NewProc("IsWindow").Call(uintptr(hwnd))
	return r1 != 0
}

// GetTickCount wraps the Win32 API GetTickCount.
func GetTickCount() uint32 {
	r1, _, _ := Kernel32.NewProc("GetTickCount").Call()
	return uint32(r1)
}

// IsChildWindowW wraps the Win32 API IsChild.
func IsChildWindowW(parent, child windows.HWND) bool {
	r1, _, _ := User32.NewProc("IsChild").Call(uintptr(parent), uintptr(child))
	return r1 != 0
}

// PeekMessageW wraps the Win32 API PeekMessageW.
// flags: 0=PM_NOREMOVE, 1=PM_REMOVE
func PeekMessageW(lpMsg *MSG, hwnd windows.HWND, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) int32 {
	r1, _, _ := User32.NewProc("PeekMessageW").Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		uintptr(wRemoveMsg),
	)
	return int32(r1)
}

// SleepW wraps the Win32 API Sleep.
func SleepW(dwMilliseconds uint32) {
	Kernel32.NewProc("Sleep").Call(uintptr(dwMilliseconds))
}

// GetDPIScaleFactor returns the DPI scale factor relative to 96 DPI (standard).
// For example, on a 150% DPI display, returns 1.5.
func GetDPIScaleFactor() float64 {
	// Try GetSystemDpiForProcess (Windows 10+)
	r1, _, _ := User32.NewProc("GetSystemDpiForProcess").Call(^uintptr(0)) // GetCurrentProcess()
	if r1 > 0 {
		dpi := int32(r1)
		if dpi > 0 {
			return float64(dpi) / 96.0
		}
	}

	// Fallback to GetDpiForSystem (older Windows versions)
	r1, _, _ = User32.NewProc("GetDpiForSystem").Call()
	if r1 > 0 {
		dpi := int32(r1)
		if dpi > 0 {
			return float64(dpi) / 96.0
		}
	}

	// Default to no scaling
	return 1.0
}

// ScaleX scales a horizontal coordinate for current DPI.
func ScaleX(x int32) int32 {
	return int32(float64(x) * GetDPIScaleFactor())
}

// ScaleY scales a vertical coordinate for current DPI.
func ScaleY(y int32) int32 {
	return int32(float64(y) * GetDPIScaleFactor())
}

// ScaleSize scales a size value for current DPI.
func ScaleSize(size int32) int32 {
	return int32(float64(size) * GetDPIScaleFactor())
}

// CreateFontW creates a new font with the specified parameters.
// height: height in logical units (negative values for character height)
// width: average character width (0 = let Windows calculate)
// escapement: angle in tenths of degrees
// orientation: font rotation angle in tenths of degrees
// weight: FW_* constants (e.g., FW_NORMAL=400, FW_BOLD=700)
// italic, underline, strikeOut: bool flags
// charset: DEFAULT_CHARSET=1
// outPrecision, clipPrecision: OUT_*, CLIP_* constants
// quality: DRAFT_QUALITY, PROOF_QUALITY, etc.
// pitchAndFamily: pitch and family flags
// faceName: font name (e.g., "Arial", "Segoe UI")
func CreateFontW(height, width, escapement, orientation, weight int32,
	italic, underline, strikeOut bool, charset, outPrecision, clipPrecision, quality, pitchAndFamily byte,
	faceName string,
) (windows.Handle, error) {
	faceName16, _ := windows.UTF16PtrFromString(faceName)

	italicVal := byte(0)
	if italic {
		italicVal = 1
	}
	underlineVal := byte(0)
	if underline {
		underlineVal = 1
	}
	strikeOutVal := byte(0)
	if strikeOut {
		strikeOutVal = 1
	}

	r1, _, _ := Gdi32.NewProc("CreateFontW").Call(
		uintptr(height), uintptr(width), uintptr(escapement), uintptr(orientation),
		uintptr(weight),
		uintptr(italicVal), uintptr(underlineVal), uintptr(strikeOutVal),
		uintptr(charset), uintptr(outPrecision), uintptr(clipPrecision), uintptr(quality),
		uintptr(pitchAndFamily), uintptr(unsafe.Pointer(faceName16)),
	)

	if r1 == 0 {
		return 0, errors.New("CreateFontW failed")
	}
	return windows.Handle(r1), nil
}

// SetWindowFontW sets the font for a window.
// Sends WM_SETFONT message to the window.
func SetWindowFontW(hwnd windows.HWND, hFont windows.Handle, redraw bool) {
	redrawVal := uintptr(0)
	if redraw {
		redrawVal = 1
	}
	User32.NewProc("SendMessageW").Call(
		uintptr(hwnd),
		uintptr(WM_SETFONT),
		uintptr(hFont),
		redrawVal,
	)
}

const (
	// Font weights
	FW_NORMAL = 400
	FW_BOLD   = 700

	// Charset
	DEFAULT_CHARSET = 1

	// Output precision
	OUT_DEFAULT_PRECIS = 0

	// Clip precision
	CLIP_DEFAULT_PRECIS = 0

	// Quality
	DEFAULT_QUALITY = 0
	PROOF_QUALITY   = 2

	// Pitch and family
	FF_DONTCARE = 0
)
