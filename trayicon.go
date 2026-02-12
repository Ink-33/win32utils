package win32utils

import (
	"fmt"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// MenuItemCallback is called when a menu item is selected.
type MenuItemCallback func(itemID int32)

// PopupMenu manages a context menu for the tray icon.
type PopupMenu struct {
	mu      sync.Mutex
	hMenu   windows.Handle
	items   map[int32]MenuItemCallback
	nextID  int32
}

// NewPopupMenu creates a new popup menu.
func NewPopupMenu() (*PopupMenu, error) {
	hMenu, err := CreatePopupMenu()
	if err != nil {
		return nil, fmt.Errorf("failed to create popup menu: %w", err)
	}
	return &PopupMenu{
		hMenu:  hMenu,
		items:  make(map[int32]MenuItemCallback),
		nextID: 1000, // Start from 1000 to avoid conflicts
	}, nil
}

// AddItem adds a menu item with a callback.
func (m *PopupMenu) AddItem(label string, callback MenuItemCallback) (int32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextID
	m.nextID++

	if err := AppendMenuW(m.hMenu, MFT_STRING|MF_ENABLED, uint32(id), label); err != nil {
		return 0, err
	}

	if callback != nil {
		m.items[id] = callback
	}
	return id, nil
}

// AddSeparator adds a menu separator.
func (m *PopupMenu) AddSeparator() error {
	return AppendMenuW(m.hMenu, MFT_SEPARATOR, 0, "")
}

// GetCallback retrieves the callback for a menu item.
func (m *PopupMenu) GetCallback(itemID int32) (MenuItemCallback, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	callback, ok := m.items[itemID]
	return callback, ok
}

// Show displays the menu at the given coordinates.
// Returns the selected item ID.
func (m *PopupMenu) Show(hwnd windows.HWND, x, y int32) (int32, error) {
	return TrackPopupMenuEx(
		m.hMenu,
		TPM_RIGHTALIGN|TPM_TOPALIGN|TPM_RETURNCMD,
		x, y,
		hwnd,
	)
}

// Destroy destroys the menu and cleans up resources.
func (m *PopupMenu) Destroy() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := DestroyMenu(m.hMenu); err != nil {
		return err
	}
	m.items = nil
	m.hMenu = 0
	return nil
}

// TrayIconCallback is called when the user interacts with the tray icon.
// mouseMsg contains the mouse event (WM_LBUTTONDOWN, WM_RBUTTONDOWN, etc.).
type TrayIconCallback func(mouseMsg uint32)

// TrayIcon manages a system tray icon instance.
type TrayIcon struct {
	hwnd     windows.HWND
	uid      uint32
	callback TrayIconCallback
	msgID    uint32 // Custom WM_USER-based message ID
	menu     *PopupMenu
}

// NewTrayIcon creates a new TrayIcon instance.
// callback is optional and will be called when the user interacts with the tray icon.
func NewTrayIcon(uid uint32, callback TrayIconCallback) (*TrayIcon, error) {
	if uid == 0 {
		uid = 1 // Default to 1 if not specified
	}

	// Custom message ID for tray icon callback
	msgID := WM_USER + 1

	ti := &TrayIcon{
		uid:      uid,
		callback: callback,
		msgID:    msgID,
	}

	// Create message-only window
	hwnd, err := CreateMessageOnlyWindow(
		"win32utils.TrayIcon",
		"Tray Icon",
		func(h windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
			return defaultTrayWndProc(h, msg, wParam, lParam, ti)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message window: %w", err)
	}

	ti.hwnd = hwnd
	return ti, nil
}

// Add registers the tray icon with the system tray.
func (ti *TrayIcon) Add(hIcon windows.Handle, tip string) error {
	iconData := ti.buildNotifyIconData()
	iconData.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	iconData.HIcon = hIcon
	iconData.UCallbackMessage = ti.msgID

	// Copy tooltip (max 128 chars including null terminator)
	tipUtf16, _ := windows.UTF16FromString(tip)
	for i, ch := range tipUtf16 {
		if i >= len(iconData.SzTip) {
			break
		}
		iconData.SzTip[i] = ch
	}

	return ShellNotifyIconW(NIM_ADD, iconData)
}

// Remove removes the tray icon from the system tray.
func (ti *TrayIcon) Remove() error {
	iconData := ti.buildNotifyIconData()
	return ShellNotifyIconW(NIM_DELETE, iconData)
}

// Update updates the tray icon (icon, tip, or other properties).
func (ti *TrayIcon) Update(hIcon windows.Handle, tip string) error {
	iconData := ti.buildNotifyIconData()
	iconData.UFlags = NIF_ICON | NIF_TIP
	iconData.HIcon = hIcon

	// Copy tooltip
	tipUtf16, _ := windows.UTF16FromString(tip)
	for i, ch := range tipUtf16 {
		if i >= len(iconData.SzTip) {
			break
		}
		iconData.SzTip[i] = ch
	}

	return ShellNotifyIconW(NIM_MODIFY, iconData)
}

// Close removes the icon and destroys the associated window.
func (ti *TrayIcon) Close() error {
	if err := ti.Remove(); err != nil {
		// Ignore error if icon is already removed
	}
	if ti.menu != nil {
		_ = ti.menu.Destroy()
		ti.menu = nil
	}
	return DestroyWindow(ti.hwnd)
}

// HWND returns the underlying window handle (useful for advanced use cases).
func (ti *TrayIcon) HWND() windows.HWND {
	return ti.hwnd
}

// Menu returns the popup menu associated with this tray icon.
// Create one with SetupMenu() if needed.
func (ti *TrayIcon) Menu() *PopupMenu {
	return ti.menu
}

// SetupMenu creates and associates a popup menu with this tray icon.
func (ti *TrayIcon) SetupMenu() (*PopupMenu, error) {
	if ti.menu != nil {
		return ti.menu, nil
	}
	menu, err := NewPopupMenu()
	if err != nil {
		return nil, err
	}
	ti.menu = menu
	return menu, nil
}

func (ti *TrayIcon) buildNotifyIconData() *NOTIFYICONDATAW {
	return &NOTIFYICONDATAW{
		HWnd: ti.hwnd,
		UID:  ti.uid,
	}
}

// defaultTrayWndProc is the default window procedure for a tray icon window.
func defaultTrayWndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr, ti *TrayIcon) uintptr {
	switch msg {
	case ti.msgID:
		if ti.callback != nil {
			mouseMsg := uint32(lParam)
			ti.callback(mouseMsg)

			// Show popup menu on right-click
			if mouseMsg == WM_RBUTTONDOWN && ti.menu != nil {
				// Get cursor position
				pt := POINT{}
				r1, _, _ := User32.NewProc("GetCursorPos").Call(uintptr(unsafe.Pointer(&pt)))
				if r1 != 0 {
					selectedID, _ := ti.menu.Show(hwnd, pt.X, pt.Y)
					if selectedID != 0 {
						if callback, ok := ti.menu.GetCallback(selectedID); ok && callback != nil {
							callback(selectedID)
						}
					}
				}
			}
		}
		return 0

	case WM_COMMAND:
		// Handle menu item selection
		cmd := int32(wParam & 0xFFFF)
		if ti.menu != nil {
			if callback, ok := ti.menu.GetCallback(cmd); ok && callback != nil {
				callback(cmd)
			}
		}
		return 0

	case WM_DESTROY:
		PostQuitMessage(0)
		return 0

	default:
		return DefWindowProcW(hwnd, msg, wParam, lParam)
	}
}
