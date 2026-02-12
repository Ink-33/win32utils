package win32utils

import "golang.org/x/sys/windows"

var Kernel32 = windows.NewLazySystemDLL("kernel32.dll")
var User32 = windows.NewLazySystemDLL("user32.dll")
var Shell32 = windows.NewLazySystemDLL("shell32.dll")
var Gdi32 = windows.NewLazySystemDLL("gdi32.dll")