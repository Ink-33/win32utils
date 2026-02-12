package win32utils

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 在测试开始前保存原始控制台标题
	originalTitle, err := GetConsoleTitle()
	if err != nil {
		log.Printf("Warning: Could not get original console title: %v", err)
		originalTitle = ""
	}

	// 运行测试
	exitCode := m.Run()

	// 测试结束后恢复原始控制台标题
	if originalTitle != "" {
		if err := SetConsoleTitle(originalTitle); err != nil {
			log.Printf("Warning: Could not restore console title: %v", err)
		}
	}

	os.Exit(exitCode)
}

func TestGetConsoleWindow(t *testing.T) {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}
	if hwnd == 0 {
		t.Error("Expected non-zero HWND")
	}
	t.Logf("Console window handle: %v", hwnd)
}

func TestShowHideConsole(t *testing.T) {
	// 获取初始状态
	initialVisible, err := IsConsoleVisible()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}
	t.Logf("Initial console visibility: %v", initialVisible)

	// 测试隐藏控制台
	err = HideConsole()
	if err != nil {
		t.Fatalf("HideConsole failed: %v", err)
	}

	// 等待一小段时间让状态更新
	time.Sleep(100 * time.Millisecond)

	visible, err := IsConsoleVisible()
	if err != nil {
		t.Fatalf("IsConsoleVisible failed: %v", err)
	}
	if visible {
		t.Error("Console should be hidden but is still visible")
	}
	t.Log("Console successfully hidden")

	// 测试显示控制台
	err = ShowConsole()
	if err != nil {
		t.Fatalf("ShowConsole failed: %v", err)
	}

	// 等待一小段时间让状态更新
	time.Sleep(100 * time.Millisecond)

	visible, err = IsConsoleVisible()
	if err != nil {
		t.Fatalf("IsConsoleVisible failed: %v", err)
	}
	if !visible {
		t.Error("Console should be visible but is still hidden")
	}
	t.Log("Console successfully shown")

	// 恢复初始状态
	if initialVisible {
		err = ShowConsole()
	} else {
		err = HideConsole()
	}
	if err != nil {
		t.Logf("Warning: Could not restore initial console state: %v", err)
	}
}

func TestToggleConsole(t *testing.T) {
	// 获取初始状态
	initialVisible, err := IsConsoleVisible()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}
	t.Logf("Initial console visibility: %v", initialVisible)

	// 测试切换功能
	newState, err := ToggleConsole()
	if err != nil {
		t.Fatalf("ToggleConsole failed: %v", err)
	}
	t.Logf("After toggle, console visibility: %v", newState)

	// 验证状态确实改变了
	if newState == initialVisible {
		t.Error("ToggleConsole did not change the visibility state")
	}

	// 再次切换回来
	finalState, err := ToggleConsole()
	if err != nil {
		t.Fatalf("Second ToggleConsole failed: %v", err)
	}
	t.Logf("After second toggle, console visibility: %v", finalState)

	// 验证状态回到了初始状态
	if finalState != initialVisible {
		t.Error("Double toggle did not restore original visibility state")
	}
}

func TestConsoleTitle(t *testing.T) {
	// 获取原始标题
	originalTitle, err := GetConsoleTitle()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}
	t.Logf("Original console title: '%s'", originalTitle)

	// 设置新标题
	testTitle := "Win32Utils Test Console"
	err = SetConsoleTitle(testTitle)
	if err != nil {
		t.Fatalf("SetConsoleTitle failed: %v", err)
	}

	// 验证标题已设置
	newTitle, err := GetConsoleTitle()
	if err != nil {
		t.Fatalf("GetConsoleTitle failed: %v", err)
	}
	if newTitle != testTitle {
		t.Errorf("Expected title '%s', got '%s'", testTitle, newTitle)
	}
	t.Logf("Console title successfully set to: '%s'", newTitle)

	// 恢复原始标题
	err = SetConsoleTitle(originalTitle)
	if err != nil {
		t.Logf("Warning: Could not restore original console title: %v", err)
	} else {
		restoredTitle, err := GetConsoleTitle()
		if err != nil {
			t.Logf("Warning: Could not verify title restoration: %v", err)
		} else if restoredTitle != originalTitle {
			t.Logf("Warning: Title not properly restored. Expected '%s', got '%s'", originalTitle, restoredTitle)
		}
	}
}

func TestShowWindowCommands(t *testing.T) {
	hwnd, err := GetConsoleWindow()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	// 测试各种ShowWindow命令
	commands := []struct {
		name    string
		cmd     int32
		visible bool
	}{
		{"SW_HIDE", SW_HIDE, false},
		{"SW_SHOW", SW_SHOW, true},
		{"SW_SHOWNORMAL", SW_SHOWNORMAL, true},
		{"SW_SHOWMINIMIZED", SW_SHOWMINIMIZED, true}, // minimized windows are still "visible"
		{"SW_SHOWMAXIMIZED", SW_SHOWMAXIMIZED, true},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			err := ShowWindow(hwnd, cmd.cmd)
			if err != nil {
				t.Errorf("%s failed: %v", cmd.name, err)
				return
			}

			// 给系统一点时间处理命令
			time.Sleep(50 * time.Millisecond)

			visible, err := IsConsoleVisible()
			if err != nil {
				t.Errorf("IsConsoleVisible failed after %s: %v", cmd.name, err)
				return
			}

			// 注意：对于最小化状态，IsWindowVisible可能仍然返回true
			// 因为我们只测试基本的显示/隐藏功能，这里主要验证没有错误
			t.Logf("%s executed successfully, console visible: %v", cmd.name, visible)
		})
	}
}

func BenchmarkConsoleOperations(b *testing.B) {
	b.Run("GetConsoleWindow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GetConsoleWindow()
		}
	})

	b.Run("IsConsoleVisible", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = IsConsoleVisible()
		}
	})

	b.Run("ShowConsole", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ShowConsole()
		}
	})

	b.Run("HideConsole", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = HideConsole()
		}
	})
}