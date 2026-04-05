//go:build windows
// +build windows

package browsers

import (
	"testing"
)

// ---------------------------------------------------------------------------
// isValidURL
// ---------------------------------------------------------------------------

func TestIsValidURL(t *testing.T) {
	valid := []string{
		"http://example.com",
		"https://example.com",
		"https://www.google.com/search?q=hello",
		"http://localhost:8080/path",
		"https://github.com/mihn1/timekeeper",
	}
	for _, u := range valid {
		if !isValidURL(u) {
			t.Errorf("isValidURL(%q) = false, want true", u)
		}
	}
}

func TestIsValidURL_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"http",          // too short / no host
		"ftp://example", // wrong scheme
		"just a title",
		"example.com",       // no scheme
		"file:///some/path", // wrong scheme
		"https",             // exactly 5 chars, < 7
		"https:/",           // 7 chars but still invalid host
	}
	for _, u := range invalid {
		if isValidURL(u) {
			t.Errorf("isValidURL(%q) = true, want false", u)
		}
	}
}

// ---------------------------------------------------------------------------
// findAddressBarInWindow with hwnd=0 must not panic
// ---------------------------------------------------------------------------

func TestFindAddressBarInWindow_NullHwnd(t *testing.T) {
	e := NewBrowserURLExtractor(nil)
	// hwnd=0 — EnumChildWindows with NULL parent enumerates top-level windows.
	// The call should return an empty string gracefully, not panic.
	got := e.findAddressBarInWindow(0)
	// We don't assert a specific value — just that it doesn't crash and
	// any returned string is either empty or a valid URL.
	if got != "" && !isValidURL(got) {
		t.Errorf("findAddressBarInWindow(0) returned non-empty non-URL: %q", got)
	}
}

// ---------------------------------------------------------------------------
// ExtractURLFromWindow – hwnd=0 should return empty or valid URL, no panic
// ---------------------------------------------------------------------------

func TestExtractURLFromWindow_NullHwnd(t *testing.T) {
	e := NewBrowserURLExtractor(nil)
	got := e.ExtractURLFromWindow(0, "Google Chrome", "Test Title")
	if got != "" && !isValidURL(got) {
		t.Errorf("unexpected non-URL result %q", got)
	}
}

// ---------------------------------------------------------------------------
// normalizeURL
// ---------------------------------------------------------------------------

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"https://github.com", "https://github.com"},
		{"http://localhost:8080", "http://localhost:8080"},
		{"github.com/mihn1/timekeeper", "https://github.com/mihn1/timekeeper"},
		{"www.google.com", "https://www.google.com"},
		{"example.com", "https://example.com"},
		{"", ""},
		{"just text", ""},
		{"no-dots", ""},
		{"  https://trimmed.com  ", "https://trimmed.com"},
		{"  github.com  ", "https://github.com"},
	}
	for _, tt := range tests {
		got := normalizeURL(tt.input)
		if got != tt.want {
			t.Errorf("normalizeURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
