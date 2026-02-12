//go:build windows

package win32utils

import "testing"

func TestWindowStyle_WithWithoutHas(t *testing.T) {
	s := WindowStyle{}
	s = s.With(WS_VISIBLE | WS_SYSMENU)
	if !s.Has(WS_VISIBLE) {
		t.Fatalf("expected WS_VISIBLE to be set")
	}
	if !s.Has(WS_SYSMENU) {
		t.Fatalf("expected WS_SYSMENU to be set")
	}

	s2 := s.Without(WS_VISIBLE)
	if s2.Has(WS_VISIBLE) {
		t.Fatalf("expected WS_VISIBLE to be cleared")
	}
	if !s2.Has(WS_SYSMENU) {
		t.Fatalf("expected WS_SYSMENU to remain set")
	}
}

func TestWindowExStyle_WithWithoutHas(t *testing.T) {
	ex := WindowExStyle{}
	ex = ex.With(WS_EX_TOPMOST | WS_EX_TOOLWINDOW)
	if !ex.Has(WS_EX_TOPMOST) || !ex.Has(WS_EX_TOOLWINDOW) {
		t.Fatalf("expected exstyle bits to be set")
	}
	ex2 := ex.Without(WS_EX_TOPMOST)
	if ex2.Has(WS_EX_TOPMOST) {
		t.Fatalf("expected WS_EX_TOPMOST to be cleared")
	}
	if !ex2.Has(WS_EX_TOOLWINDOW) {
		t.Fatalf("expected WS_EX_TOOLWINDOW to remain set")
	}
}
