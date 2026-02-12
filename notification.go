//go:build windows

package win32utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// System icon constants for notifications
const (
	IconInformation = iota // Info icon (blue)
	IconWarning            // Warning icon (yellow)
	IconError              // Error icon (red)
	IconSuccess            // Success/checkmark (green)
	IconQuestion           // Question mark
)

// ToastNotification represents a Windows Toast notification.
type ToastNotification struct {
	AppID    string // Application ID (e.g., "MyApp")
	Title    string
	Message  string
	LogoIcon string // Path to small square logo icon (128x128)
}

// ShowToast displays a Windows Toast notification.
// This uses PowerShell to interact with Windows Runtime APIs.
// AppID should be a simple identifier like "MyApp".
func (tn *ToastNotification) Show() error {
	if tn.AppID == "" {
		tn.AppID = "GoApp"
	}

	// Build the XML template - with logo if provided
	xmlTemplate := `<toast>
    <visual>
        <binding template="ToastText02"{{if .LogoIcon}} addImageQuery="true"{{end}}>{{if .LogoIcon}}
            <image id="1" placement="appLogo" src="{{.LogoIcon}}" alt="icon"/>{{end}}
            <text id="1">{{.Title}}</text>
            <text id="2">{{.Message}}</text>
        </binding>
    </visual>
</toast>`

	tmpl, err := template.New("toast").Parse(xmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tn)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	xml := buf.String()

	// Create PowerShell script
	psScript := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '%s'
$template = @"
%s
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
`, tn.AppID, xml)

	// Execute PowerShell script
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to show toast: %w", err)
	}

	return nil
}

// SimpleToast is a simplified version that doesn't require COM setup.
// It uses WinRT through PowerShell.
func SimpleToast(appID, title, message string) error {
	tn := &ToastNotification{
		AppID:   appID,
		Title:   title,
		Message: message,
	}
	return tn.Show()
}

// ToastNotificationBuilder provides a fluent interface for creating toasts.
type ToastNotificationBuilder struct {
	toast *ToastNotification
}

// NewToastBuilder creates a new toast builder.
func NewToastBuilder(appID string) *ToastNotificationBuilder {
	return &ToastNotificationBuilder{
		toast: &ToastNotification{AppID: appID},
	}
}

// Title sets the toast title.
func (b *ToastNotificationBuilder) Title(title string) *ToastNotificationBuilder {
	b.toast.Title = title
	return b
}

// Message sets the toast message.
func (b *ToastNotificationBuilder) Message(message string) *ToastNotificationBuilder {
	b.toast.Message = message
	return b
}

// Icon sets the toast icon path (appears as app logo).
func (b *ToastNotificationBuilder) Icon(iconPath string) *ToastNotificationBuilder {
	b.toast.LogoIcon = iconPath
	return b
}

// Show displays the toast notification.
func (b *ToastNotificationBuilder) Show() error {
	return b.toast.Show()
}

// AdvancedToastNotification supports more features (actions, sounds, images, etc).
type AdvancedToastNotification struct {
	AppID       string
	Title       string
	Message     string
	SubTitle    string // Secondary text
	LogoIcon    string // Small square icon (128x128) - appears top-left
	HeroImage   string // Large banner image (364x180) - appears at top
	InlineImage string // Inline image - appears in notification body
	Sound       string // "default", "silent", or path to sound file
	Duration    string // "short" or "long"
	Actions     []ToastAction
}

// ToastAction represents an action button on the toast.
type ToastAction struct {
	Content   string // Button label
	Arguments string // What to do when clicked
	Activate  string // "foreground" or "background"
}

// Show displays an advanced toast notification.
func (atn *AdvancedToastNotification) Show() error {
	if atn.AppID == "" {
		atn.AppID = "GoApp"
	}

	// Build advanced XML template with image support
	xmlTemplate := `<toast duration="{{.Duration}}">{{if .HeroImage}}
    <visual>
        <binding template="ToastGeneric" addImageQuery="true">
            <image placement="hero" src="{{.HeroImage}}" alt="hero"/>{{if .LogoIcon}}
            <image id="1" placement="appLogo" src="{{.LogoIcon}}" alt="logo" hint-crop="circle"/>{{end}}
            <text placement="attribution">{{.SubTitle}}</text>
            <text>{{.Title}}</text>
            <text hint-style="base" hint-wrap="true">{{.Message}}</text>
        </binding>
    </visual>{{else}}
    <visual>
        <binding template="ToastText04"{{if .LogoIcon}} addImageQuery="true"{{end}}>{{if .LogoIcon}}
            <image id="1" placement="appLogo" src="{{.LogoIcon}}" alt="logo"/>{{end}}
            <text id="1">{{.Title}}</text>{{if .SubTitle}}
            <text id="2">{{.SubTitle}}</text>
            <text id="3">{{.Message}}</text>{{else}}
            <text id="2">{{.Message}}</text>{{end}}{{if .InlineImage}}
            <image id="2" src="{{.InlineImage}}" alt="inline"/>{{end}}
        </binding>
    </visual>{{end}}{{if .Sound}}
    <audio src="{{.Sound}}" />{{end}}{{if .Actions}}
    <actions>{{range .Actions}}
        <action content="{{.Content}}" arguments="{{.Arguments}}" activationType="{{.Activate}}" />{{end}}
    </actions>{{end}}
</toast>`

	tmpl, err := template.New("advanced-toast").Parse(xmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, atn)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	xml := buf.String()

	// Create PowerShell script
	psScript := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '%s'
$template = @"
%s
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
`, atn.AppID, xml)

	// Execute PowerShell script
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to show advanced toast: %w", err)
	}

	return nil
}

// AdvancedBuilder provides fluent interface for advanced toasts.
type AdvancedBuilder struct {
	toast *AdvancedToastNotification
}

// NewAdvancedToastBuilder creates a new advanced toast builder.
func NewAdvancedToastBuilder(appID string) *AdvancedBuilder {
	return &AdvancedBuilder{
		toast: &AdvancedToastNotification{
			AppID:    appID,
			Duration: "long",
		},
	}
}

// Title sets the toast title.
func (b *AdvancedBuilder) Title(title string) *AdvancedBuilder {
	b.toast.Title = title
	return b
}

// Message sets the toast message.
func (b *AdvancedBuilder) Message(message string) *AdvancedBuilder {
	b.toast.Message = message
	return b
}

// SubTitle sets the subtitle.
func (b *AdvancedBuilder) SubTitle(subtitle string) *AdvancedBuilder {
	b.toast.SubTitle = subtitle
	return b
}

// LogoIcon sets the small square app logo (128x128).
func (b *AdvancedBuilder) LogoIcon(iconPath string) *AdvancedBuilder {
	b.toast.LogoIcon = iconPath
	return b
}

// HeroImage sets the large banner image (364x180) at the top of the toast.
func (b *AdvancedBuilder) HeroImage(imagePath string) *AdvancedBuilder {
	b.toast.HeroImage = imagePath
	return b
}

// InlineImage sets an inline image in the notification body.
func (b *AdvancedBuilder) InlineImage(imagePath string) *AdvancedBuilder {
	b.toast.InlineImage = imagePath
	return b
}

// Sound sets the notification sound (e.g., "default", "silent").
func (b *AdvancedBuilder) Sound(sound string) *AdvancedBuilder {
	b.toast.Sound = sound
	return b
}

// Duration sets how long the toast shows ("short" or "long").
func (b *AdvancedBuilder) Duration(duration string) *AdvancedBuilder {
	if duration == "short" || duration == "long" {
		b.toast.Duration = duration
	}
	return b
}

// AddAction adds an action button to the toast.
func (b *AdvancedBuilder) AddAction(label, args, activationType string) *AdvancedBuilder {
	b.toast.Actions = append(b.toast.Actions, ToastAction{
		Content:   label,
		Arguments: args,
		Activate:  activationType,
	})
	return b
}

// Show displays the advanced toast.
func (b *AdvancedBuilder) Show() error {
	return b.toast.Show()
}

// System icon helpers
// GetSystemIconPath returns a path to a common Windows system icon.
// This uses the ms-appx:// scheme which is more reliable for Toast notifications.
// Usage: GetSystemIconPath("info") for info icon, "warning", "error", "success", "question"
func GetSystemIconPath(iconType string) string {
	// Use ms-appx:/// to reference system resources
	// These paths work more reliably in Toast notifications
	switch iconType {
	case "info":
		// Use info icon from system
		return "ms-appx:///Assets/InformationBadge.png"
	case "warning":
		// Use warning icon
		return "ms-appx:///Assets/WarningBadge.png"
	case "error":
		// Use error icon
		return "ms-appx:///Assets/ErrorBadge.png"
	case "success":
		// Use checkmark icon
		return "ms-appx:///Assets/CheckmarkBadge.png"
	case "question":
		// Use question icon
		return "ms-appx:///Assets/QuestionBadge.png"
	default:
		return ""
	}
}

// GetIconPathFromFile returns a properly formatted file:/// URI for a local icon file.
// Use this for custom PNG/JPG icon files.
// Example: GetIconPathFromFile("C:\\Users\\user\\icon.png")
func GetIconPathFromFile(filePath string) string {
	// Convert Windows path to file:/// URI format
	// Replace backslashes with forward slashes and URL-encode if needed
	if filePath == "" {
		return ""
	}

	// Ensure file exists
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}

	// Convert to forward slashes for URL
	uri := filePath
	// Windows paths need to be formatted as file:///C:/path/to/file
	// The easiest way is to replace backslashes and add the protocol
	if len(uri) > 0 && uri[0] != '/' {
		// Add leading slash for drive letter
		uri = "/" + uri
	}

	// Replace backslashes with forward slashes
	for i := 0; i < len(uri); i++ {
		if uri[i] == '\\' {
			uri = uri[:i] + "/" + uri[i+1:]
		}
	}

	return "file://" + uri
}

// Quick helper functions

// NotifyProgress sends a progress notification (simplified).
func NotifyProgress(title, message string, progress int) error {
	// Note: True progress bars require more complex XML and WinRT APIs
	msg := fmt.Sprintf("%s: %d%%", message, progress)
	return SimpleToast("GoApp.Progress", title, msg)
}

// NotifyWarning sends a warning notification with icon.
func NotifyWarning(title, message string) error {
	return SimpleColoredNotification("GoApp", "warning", title, message)
}

// NotifySuccess sends a success notification with icon.
func NotifySuccess(title, message string) error {
	return SimpleColoredNotification("GoApp", "success", title, message)
}

// NotifyError sends an error notification with icon.
func NotifyError(title, message string) error {
	return SimpleColoredNotification("GoApp", "error", title, message)
}
