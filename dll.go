package win32utils

import "golang.org/x/sys/windows"

var Kernel32 = windows.NewLazySystemDLL("kernel32.dll")
var User32 = windows.NewLazySystemDLL("user32.dll")
