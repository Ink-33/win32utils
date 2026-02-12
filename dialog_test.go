//go:build windows

package win32utils

import (
	"testing"
)

func TestUsernamePasswordDialog(t *testing.T) {
	// Note: This test requires manual interaction since it shows a GUI dialog
	// In automated testing environments, this would need to be mocked or skipped
	
	t.Skip("Manual test - requires user interaction with dialog")
	
	/*
	// Example usage:
	username, password, cancelled, err := UsernamePasswordDialog(
		"Login",
		"Username:",
		"Password:",
		"testuser",
	)
	
	if err != nil {
		t.Fatalf("Dialog failed: %v", err)
	}
	
	if cancelled {
		t.Log("Dialog was cancelled by user")
	} else {
		t.Logf("Username: %s", username)
		t.Logf("Password length: %d", len(password))
		// Note: Don't log actual password for security reasons
	}
	*/
}

func TestTwoTextInputDialog(t *testing.T) {
	// Note: This test requires manual interaction since it shows a GUI dialog
	
	t.Skip("Manual test - requires user interaction with dialog")
	
	/*
	// Example usage:
	text1, text2, cancelled, err := TwoTextInputDialog(
		"Input Dialog",
		"Field 1:",
		"Field 2:",
		"default1",
		"default2",
	)
	
	if err != nil {
		t.Fatalf("Dialog failed: %v", err)
	}
	
	if cancelled {
		t.Log("Dialog was cancelled by user")
	} else {
		t.Logf("Text1: %s", text1)
		t.Logf("Text2: %s", text2)
	}
	*/
}