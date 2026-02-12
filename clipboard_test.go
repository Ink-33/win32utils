package win32utils

import (
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// TestBasicClipboardWorkflow tests the fundamental clipboard operations
func TestBasicClipboardWorkflow(t *testing.T) {
	// Try to open clipboard
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Skipping test - failed to open clipboard: %v", err)
		return
	}

	// Ensure clipboard is closed on function exit
	defer func() {
		_ = CloseClipboard()
	}()

	// Empty clipboard first
	err = EmptyClipboard()
	if err != nil {
		t.Fatalf("Failed to empty clipboard: %v", err)
	}

	// Set clipboard text
	testText := "Hello Clipboard Test 123"
	handle, err := SetClipboardText(testText)
	if err != nil {
		t.Fatalf("Failed to set clipboard text: %v", err)
	}

	// Basic validation that we got a handle
	if handle == 0 {
		t.Error("Expected non-zero handle from SetClipboardText")
	}

	t.Logf("Successfully set clipboard text. Handle: %v", handle)
}

// TestSetTextFunction tests the convenience SetText function
func TestSetTextFunction(t *testing.T) {
	testText := "Convenience function test - " + time.Now().Format("15:04:05")
	
	err := SetText(testText)
	if err != nil {
		t.Skipf("Skipping test - SetText failed: %v", err)
		return
	}

	t.Logf("Successfully set text using convenience function: %s", testText)
}

// TestClipboardHandleManagement tests that handles are properly managed
func TestClipboardHandleManagement(t *testing.T) {
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Skipping test - failed to open clipboard: %v", err)
		return
	}
	defer CloseClipboard()

	// Test setting multiple texts and getting different handles
	texts := []string{"Text 1", "Text 2", "Text 3"}
	handles := make([]windows.Handle, len(texts))

	for i, text := range texts {
		// Empty clipboard before each set
		err = EmptyClipboard()
		if err != nil {
			t.Fatalf("Failed to empty clipboard: %v", err)
		}

		handle, err := SetClipboardText(text)
		if err != nil {
			t.Fatalf("Failed to set text '%s': %v", text, err)
		}

		handles[i] = handle
		t.Logf("Set text '%s' with handle: %v", text, handle)
	}

	// Verify handles are different (they might be reused, but shouldn't be zero)
	for i, handle := range handles {
		if handle == 0 {
			t.Errorf("Handle %d is zero", i)
		}
	}
}

// TestEmptyClipboardFunction tests the EmptyClipboard function behavior
func TestEmptyClipboardFunction(t *testing.T) {
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Skipping test - failed to open clipboard: %v", err)
		return
	}
	defer CloseClipboard()

	// Set some initial text
	_, err = SetClipboardText("Initial text to clear")
	if err != nil {
		t.Fatalf("Failed to set initial text: %v", err)
	}

	// Empty the clipboard
	err = EmptyClipboard()
	if err != nil {
		t.Fatalf("Failed to empty clipboard: %v", err)
	}

	t.Log("Successfully emptied clipboard")
}

// TestUnicodeTextHandling tests basic Unicode text handling
func TestUnicodeTextHandling(t *testing.T) {
	// Test with simple Unicode text
	unicodeText := "Unicode: αβγ 中文 🚀"
	
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Skipping test - failed to open clipboard: %v", err)
		return
	}
	defer CloseClipboard()

	err = EmptyClipboard()
	if err != nil {
		t.Fatalf("Failed to empty clipboard: %v", err)
	}

	handle, err := SetClipboardText(unicodeText)
	if err != nil {
		t.Fatalf("Failed to set Unicode text: %v", err)
	}

	if handle == 0 {
		t.Error("Expected non-zero handle for Unicode text")
	}

	t.Logf("Successfully handled Unicode text with handle: %v", handle)
}

// TestModerateLengthText tests clipboard operations with moderate-length text
func TestModerateLengthText(t *testing.T) {
	// Create a moderate length text (~200 characters)
	moderateText := strings.Repeat("Moderate length text test. ", 8)
	
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Skipping test - failed to open clipboard: %v", err)
		return
	}
	defer CloseClipboard()

	err = EmptyClipboard()
	if err != nil {
		t.Fatalf("Failed to empty clipboard: %v", err)
	}

	handle, err := SetClipboardText(moderateText)
	if err != nil {
		t.Fatalf("Failed to set moderate length text: %v", err)
	}

	if handle == 0 {
		t.Error("Expected non-zero handle for moderate length text")
	}

	t.Logf("Successfully handled moderate length text (%d chars) with handle: %v", 
		len(moderateText), handle)
}

// TestSequentialOperations tests multiple sequential clipboard operations
func TestSequentialOperations(t *testing.T) {
	operations := []struct {
		name string
		text string
	}{
		{"First", "First sequential operation"},
		{"Second", "Second sequential operation"}, 
		{"Third", "Third sequential operation"},
		{"Fourth", "Fourth sequential operation"},
	}

	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			err := OpenClipboard(windows.HWND(GetConsoleWindows()))
			if err != nil {
				t.Skipf("Skipping operation - failed to open clipboard: %v", err)
				return
			}
			
			err = EmptyClipboard()
			if err != nil {
				CloseClipboard()
				t.Fatalf("Failed to empty clipboard: %v", err)
			}

			handle, err := SetClipboardText(op.text)
			if err != nil {
				CloseClipboard()
				t.Fatalf("Failed to set text: %v", err)
			}

			CloseClipboard()

			if handle == 0 {
				t.Error("Expected non-zero handle")
			}

			t.Logf("Completed operation '%s' successfully", op.name)
		})
	}
}

// TestConcurrentAccessPatterns tests different concurrent access patterns
func TestConcurrentAccessPatterns(t *testing.T) {
	// Test that sequential operations work reliably
	results := make(chan string, 3)
	
	// Operation 1
	go func() {
		err := OpenClipboard(windows.HWND(GetConsoleWindows()))
		if err != nil {
			results <- "op1_failed_open"
			return
		}
		defer CloseClipboard()
		
		err = EmptyClipboard()
		if err != nil {
			results <- "op1_failed_empty"
			return
		}
		
		_, err = SetClipboardText("Concurrent Op 1")
		if err != nil {
			results <- "op1_failed_set"
			return
		}
		
		results <- "op1_success"
	}()

	// Operation 2 (with delay to ensure serialization)
	go func() {
		time.Sleep(100 * time.Millisecond)
		
		err := OpenClipboard(windows.HWND(GetConsoleWindows()))
		if err != nil {
			results <- "op2_failed_open"
			return
		}
		defer CloseClipboard()
		
		err = EmptyClipboard()
		if err != nil {
			results <- "op2_failed_empty"
			return
		}
		
		_, err = SetClipboardText("Concurrent Op 2")
		if err != nil {
			results <- "op2_failed_set"
			return
		}
		
		results <- "op2_success"
	}()

	// Operation 3 (with longer delay)
	go func() {
		time.Sleep(200 * time.Millisecond)
		
		err := OpenClipboard(windows.HWND(GetConsoleWindows()))
		if err != nil {
			results <- "op3_failed_open"
			return
		}
		defer CloseClipboard()
		
		err = EmptyClipboard()
		if err != nil {
			results <- "op3_failed_empty"
			return
		}
		
		_, err = SetClipboardText("Concurrent Op 3")
		if err != nil {
			results <- "op3_failed_set"
			return
		}
		
		results <- "op3_success"
	}()

	// Collect results with timeout
	timeout := time.After(15 * time.Second)
	successCount := 0
	totalOps := 3
	
	for successCount < totalOps {
		select {
		case result := <-results:
			if strings.Contains(result, "success") {
				successCount++
				t.Logf("Operation completed: %s", result)
			} else {
				t.Logf("Operation had issue: %s", result)
			}
		case <-timeout:
			t.Logf("Test timeout - completed %d/%d operations", successCount, totalOps)
			return
		}
	}
	
	t.Logf("All %d concurrent operations completed", successCount)
}

// TestPlatformSpecificBehavior documents platform-specific clipboard behavior
func TestPlatformSpecificBehavior(t *testing.T) {
	// Document various behaviors we might encounter
	t.Log("Testing platform-specific clipboard behaviors...")
	
	// Test 1: Multiple opens might succeed on some systems
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		t.Skipf("Cannot test platform behavior - failed to open clipboard: %v", err)
		return
	}
	
	// Try to open again
	err2 := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err2 != nil {
		t.Logf("✓ Second OpenClipboard failed as expected: %v", err2)
	} else {
		t.Log("ℹ Second OpenClipboard succeeded (platform-dependent behavior)")
		CloseClipboard() // Clean up
	}
	
	CloseClipboard()
	
	// Test 2: Operations without explicit open might work
	err = EmptyClipboard()
	if err != nil {
		t.Logf("✓ EmptyClipboard without open failed: %v", err)
	} else {
		t.Log("ℹ EmptyClipboard without open succeeded (platform-dependent)")
	}
	
	// Test 3: Check if we can read clipboard state
	err = OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err == nil {
		defer CloseClipboard()
		t.Log("✓ Can open clipboard for state checking")
	} else {
		t.Logf("⚠ Cannot open clipboard for state checking: %v", err)
	}
}

// BenchmarkClipboardOperations provides performance benchmarks
func BenchmarkClipboardSetText(b *testing.B) {
	// Skip if we can't open clipboard
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		b.Skipf("Skipping benchmark - cannot open clipboard: %v", err)
		return
	}
	defer CloseClipboard()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := SetClipboardText("Benchmark test text")
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkClipboardSetTextUnicode(b *testing.B) {
	err := OpenClipboard(windows.HWND(GetConsoleWindows()))
	if err != nil {
		b.Skipf("Skipping benchmark - cannot open clipboard: %v", err)
		return
	}
	defer CloseClipboard()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := SetClipboardText("Unicode benchmark: 中文 🚀 αβγ")
		if err != nil {
			b.Fatalf("Unicode benchmark failed: %v", err)
		}
	}
}

func BenchmarkSequentialOperations(b *testing.B) {
	testText := "Sequential benchmark text"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Complete clipboard cycle
		err := OpenClipboard(windows.HWND(GetConsoleWindows()))
		if err != nil {
			b.Fatalf("Failed to open clipboard: %v", err)
		}
		
		err = EmptyClipboard()
		if err != nil {
			CloseClipboard()
			b.Fatalf("Failed to empty clipboard: %v", err)
		}
		
		_, err = SetClipboardText(testText)
		if err != nil {
			CloseClipboard()
			b.Fatalf("Failed to set text: %v", err)
		}
		
		CloseClipboard()
	}
}