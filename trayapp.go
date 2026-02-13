//go:build windows

package win32utils

import (
	"fmt"
	"sync"

	"golang.org/x/sys/windows"
)

// TrayAppConfig holds configuration for a tray application.
type TrayAppConfig struct {
	AppID         string          // Application ID for notifications
	AppName       string          // Display name
	IconID        uint32          // System icon ID (e.g., IDI_INFORMATION)
	IconTip       string          // Tooltip when hovering over tray icon
	MenuItems     []*TrayMenuItem // Menu items configuration
	OnLeftClick   func()          // Callback when tray icon is left-clicked
	OnDoubleClick func()          // Callback when tray icon is double-clicked
}

// TrayMenuItem represents a menu item in the tray context menu.
type TrayMenuItem struct {
	Label       string // Display text
	OnClick     func() // Callback when clicked
	IsSeparator bool   // If true, this is a separator line
	Icon        string // Optional: emoji or icon character
}

// TrayApp is a simplified tray application wrapper.
type TrayApp struct {
	config *TrayAppConfig
	tray   *TrayIcon
	menu   *PopupMenu
	mu     sync.RWMutex
	done   bool
	hIcon  uintptr
}

// NewTrayApp creates a new tray application with the given configuration.
func NewTrayApp(config *TrayAppConfig) (*TrayApp, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if config.AppID == "" {
		config.AppID = "TrayApp"
	}
	if config.AppName == "" {
		config.AppName = "Tray Application"
	}
	if config.IconTip == "" {
		config.IconTip = config.AppName
	}

	// Load system icon
	hIcon, _, _ := User32.NewProc("LoadIconW").Call(0, uintptr(config.IconID))
	if hIcon == 0 {
		return nil, fmt.Errorf("failed to load icon (ID: %d)", config.IconID)
	}

	app := &TrayApp{
		config: config,
		hIcon:  hIcon,
	}

	// Create tray icon
	tray, err := NewTrayIcon(1, func(mouseMsg uint32) {
		switch mouseMsg {
		case WM_LBUTTONDOWN:
			if app.config.OnLeftClick != nil {
				app.config.OnLeftClick()
			}
		case WM_LBUTTONDBLCLK:
			if app.config.OnDoubleClick != nil {
				app.config.OnDoubleClick()
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tray icon: %w", err)
	}

	app.tray = tray

	// Add tray icon
	if err := tray.Add(windows.Handle(hIcon), config.IconTip); err != nil {
		return nil, fmt.Errorf("failed to add tray icon: %w", err)
	}

	// Setup menu
	menu, err := tray.SetupMenu()
	if err != nil {
		return nil, fmt.Errorf("failed to setup menu: %w", err)
	}
	app.menu = menu

	// Add menu items
	if err := app.buildMenu(); err != nil {
		return nil, fmt.Errorf("failed to build menu: %w", err)
	}

	return app, nil
}

// buildMenu reconstructs the context menu from config.
func (app *TrayApp) buildMenu() error {
	// Clear existing menu items (this is a simplification - in reality we'd need to
	// destroy and recreate the menu, but PopupMenu handles this)
	if app.menu == nil {
		return fmt.Errorf("menu not initialized")
	}

	for _, item := range app.config.MenuItems {
		if item.IsSeparator {
			if err := app.menu.AddSeparator(); err != nil {
				return err
			}
		} else {
			label := item.Label
			if item.Icon != "" {
				label = item.Icon + " " + label
			}

			onClick := item.OnClick // Capture for closure
			_, err := app.menu.AddItem(label, func(itemID int32) {
				if onClick != nil {
					onClick()
				}
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddMenuItem adds a new menu item to the tray app menu.
func (app *TrayApp) AddMenuItem(label string, onClick func()) error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.menu == nil {
		return fmt.Errorf("menu not initialized")
	}

	_, err := app.menu.AddItem(label, func(itemID int32) {
		if onClick != nil {
			onClick()
		}
	})
	return err
}

// AddMenuItemWithEmoji adds a menu item with an emoji prefix.
func (app *TrayApp) AddMenuItemWithEmoji(emoji, label string, onClick func()) error {
	fullLabel := label
	if emoji != "" {
		fullLabel = emoji + " " + label
	}
	return app.AddMenuItem(fullLabel, onClick)
}

// AddMenuSeparator adds a separator line to the menu.
func (app *TrayApp) AddMenuSeparator() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.menu == nil {
		return fmt.Errorf("menu not initialized")
	}

	return app.menu.AddSeparator()
}

// NotificationDuration specifies how long a notification should be displayed.
type NotificationDuration string

const (
	// DurationShort - notification auto-closes after ~5 seconds
	DurationShort NotificationDuration = "short"
	// DurationLong - notification auto-closes after ~10 seconds (default)
	DurationLong NotificationDuration = "long"
)

// ShowNotification displays a simple notification with default duration (long).
func (app *TrayApp) ShowNotification(title, message string) error {
	return app.ShowNotificationEx(title, message, DurationLong)
}

// ShowNotificationEx displays a notification with custom auto-close duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationEx(title, message string, duration NotificationDuration) error {
	return SimpleToast(app.config.AppID, title, message)
}

// ShowNotificationWithEmoji displays a notification with an emoji and default duration.
func (app *TrayApp) ShowNotificationWithEmoji(emoji, title, message string) error {
	return app.ShowNotificationWithEmojiEx(emoji, title, message, DurationLong)
}

// ShowNotificationWithEmojiEx displays a notification with emoji and custom duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationWithEmojiEx(emoji, title, message string, duration NotificationDuration) error {
	d := string(duration)
	if d != "short" && d != "long" {
		d = "long"
	}
	return NewAdvancedToastBuilder(app.config.AppID).
		Title(emoji + " " + title).
		Message(message).
		Duration(d).
		Show()
}

// ShowNotificationSuccess displays a success notification with default duration.
func (app *TrayApp) ShowNotificationSuccess(title, message string) error {
	return app.ShowNotificationSuccessEx(title, message, DurationLong)
}

// ShowNotificationSuccessEx displays a success notification with custom duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationSuccessEx(title, message string, duration NotificationDuration) error {
	return app.ShowNotificationWithEmojiEx("✅", title, message, duration)
}

// ShowNotificationWarning displays a warning notification with default duration.
func (app *TrayApp) ShowNotificationWarning(title, message string) error {
	return app.ShowNotificationWarningEx(title, message, DurationLong)
}

// ShowNotificationWarningEx displays a warning notification with custom duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationWarningEx(title, message string, duration NotificationDuration) error {
	return app.ShowNotificationWithEmojiEx("⚠️", title, message, duration)
}

// ShowNotificationError displays an error notification with default duration.
func (app *TrayApp) ShowNotificationError(title, message string) error {
	return app.ShowNotificationErrorEx(title, message, DurationLong)
}

// ShowNotificationErrorEx displays an error notification with custom duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationErrorEx(title, message string, duration NotificationDuration) error {
	return app.ShowNotificationWithEmojiEx("❌", title, message, duration)
}

// ShowNotificationInfo displays an info notification with default duration.
func (app *TrayApp) ShowNotificationInfo(title, message string) error {
	return app.ShowNotificationInfoEx(title, message, DurationLong)
}

// ShowNotificationInfoEx displays an info notification with custom duration.
// duration: "short" (~5 seconds) or "long" (~10 seconds)
func (app *TrayApp) ShowNotificationInfoEx(title, message string, duration NotificationDuration) error {
	return app.ShowNotificationWithEmojiEx("ℹ️", title, message, duration)
}

// ShowDialog displays a modal text input dialog.
// Returns (text1, text2, cancelled, error).
func (app *TrayApp) ShowDialog(dialogTitle, label1, label2, default1, default2 string) (string, string, bool, error) {
	return TwoTextInputDialog(dialogTitle, label1, label2, default1, default2)
}

// ShowUsernamePasswordDialog displays a modal dialog with username and password input fields.
// The password field is masked with asterisks.
// Returns (username, password, cancelled, error).
func (app *TrayApp) ShowUsernamePasswordDialog(title, usernameLabel, passwordLabel string, defaultUsername, defaultPassword string) (string, string, bool, error) {
	return UsernamePasswordDialog(title, usernameLabel, passwordLabel, defaultUsername, defaultPassword)
}

// MessageBoxW displays a message box with the given caption, title, and flags.
func (app *TrayApp) MessageBoxW(caption, title string, flags uint) int {
	return MessageBoxW(uintptr(app.tray.hwnd), caption, title, flags)
}

// Run starts the message loop (blocking until PostQuitMessage is called).
// This typically runs in main().
func (app *TrayApp) Run() (int32, error) {
	app.mu.Lock()
	if app.done {
		app.mu.Unlock()
		return 1, fmt.Errorf("app already running or closed")
	}
	app.mu.Unlock()

	return MessageLoop()
}

// Close closes the tray application and cleans up resources.
func (app *TrayApp) Close() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.done {
		return nil
	}
	app.done = true

	if app.menu != nil {
		_ = app.menu.Destroy()
	}
	if app.tray != nil {
		return app.tray.Close()
	}
	return nil
}

// Exit triggers application exit (calls PostQuitMessage).
func (app *TrayApp) Exit() {
	PostQuitMessage(0)
}

// TrayAppBuilder provides a fluent interface for building a TrayApp.
type TrayAppBuilder struct {
	config *TrayAppConfig
}

// NewTrayAppBuilder creates a new builder for TrayApp.
func NewTrayAppBuilder(appID string) *TrayAppBuilder {
	return &TrayAppBuilder{
		config: &TrayAppConfig{
			AppID:     appID,
			AppName:   appID,
			IconID:    32516, // IDI_INFORMATION
			MenuItems: make([]*TrayMenuItem, 0),
		},
	}
}

// Name sets the application name/display name.
func (b *TrayAppBuilder) Name(name string) *TrayAppBuilder {
	b.config.AppName = name
	b.config.IconTip = name
	return b
}

// IconID sets the system icon ID to display.
func (b *TrayAppBuilder) IconID(iconID uint32) *TrayAppBuilder {
	b.config.IconID = iconID
	return b
}

// IconTip sets the tooltip text for the tray icon.
func (b *TrayAppBuilder) IconTip(tip string) *TrayAppBuilder {
	b.config.IconTip = tip
	return b
}

// OnLeftClick sets the callback for left-click events.
func (b *TrayAppBuilder) OnLeftClick(callback func()) *TrayAppBuilder {
	b.config.OnLeftClick = callback
	return b
}

// OnDoubleClick sets the callback for double-click events.
func (b *TrayAppBuilder) OnDoubleClick(callback func()) *TrayAppBuilder {
	b.config.OnDoubleClick = callback
	return b
}

// AddMenuItem adds a menu item.
func (b *TrayAppBuilder) AddMenuItem(label string, onClick func()) *TrayAppBuilder {
	b.config.MenuItems = append(b.config.MenuItems, &TrayMenuItem{
		Label:   label,
		OnClick: onClick,
	})
	return b
}

// AddMenuItemWithEmoji adds a menu item with an emoji prefix.
func (b *TrayAppBuilder) AddMenuItemWithEmoji(emoji, label string, onClick func()) *TrayAppBuilder {
	b.config.MenuItems = append(b.config.MenuItems, &TrayMenuItem{
		Label:   label,
		Icon:    emoji,
		OnClick: onClick,
	})
	return b
}

// AddMenuSeparator adds a separator line.
func (b *TrayAppBuilder) AddMenuSeparator() *TrayAppBuilder {
	b.config.MenuItems = append(b.config.MenuItems, &TrayMenuItem{
		IsSeparator: true,
	})
	return b
}

// Build creates and returns the TrayApp.
func (b *TrayAppBuilder) Build() (*TrayApp, error) {
	return NewTrayApp(b.config)
}
