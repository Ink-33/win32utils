package win32utils

import (
	"fmt"
	"testing"

	"golang.org/x/sys/windows"
)

func TestMain(m *testing.M) {
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		panic(err)
	}

	defer func() {
		err = CloseClipboard()
		if err != nil {
			panic(err)
		}
	}()

	err = EmptyClipboard()
	if err != nil {
		panic(err)
	}

	h, err := SetClipboardText("你好 Win32\r\nこんにちは Win32")
	fmt.Printf("err: %v\n", err)
	if err != nil {
		panic(err)
	}
	datab, err := GetClipboardDataText()
	fmt.Printf("err: %v\n", err)
	fmt.Printf("datab: %v\n", datab)
	fmt.Printf("h: %v\n", h)
}
