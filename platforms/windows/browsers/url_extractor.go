//go:build windows
// +build windows

package browsers

import (
	"log/slog"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	WM_GETTEXT       = 0x000D
	WM_GETTEXTLENGTH = 0x000E
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procEnumChildWindows = user32.NewProc("EnumChildWindows")
	procGetClassName     = user32.NewProc("GetClassNameW")
	procSendMessage      = user32.NewProc("SendMessageW")
)

// BrowserURLExtractor extracts URLs from specific browser windows using HWND
type BrowserURLExtractor struct {
	logger *slog.Logger
}

func NewBrowserURLExtractor(logger *slog.Logger) *BrowserURLExtractor {
	if logger == nil {
		logger = slog.Default()
	}
	return &BrowserURLExtractor{
		logger: logger,
	}
}

// ExtractURLFromWindow extracts URL from specific browser window using HWND
func (e *BrowserURLExtractor) ExtractURLFromWindow(hwnd uintptr, browserName string, title string) string {
	// Try native UI Automation (fastest for modern Chrome/Chromium)
	start := time.Now()
	url := e.extractURLUsingUIA(hwnd)
	elapsed := time.Since(start)
	e.logger.Debug("UIA URL extraction", "Browser", browserName, "Title", title, "URL", url, "DurationMs", elapsed.Milliseconds())
	if url != "" {
		return url
	}
	// Fallback: Win32 child-window enumeration (works for browsers with native controls)
	start = time.Now()
	url = e.findAddressBarInWindow(hwnd)
	elapsed = time.Since(start)
	e.logger.Debug("Win32 URL extraction", "Browser", browserName, "Title", title, "URL", url, "DurationMs", elapsed.Milliseconds())
	return url
}

// Core method to find address bar in specific window using Win32 API
func (e *BrowserURLExtractor) findAddressBarInWindow(hwnd uintptr) string {
	var foundURL string

	// Known address bar class names for different browsers
	addressBarClasses := []string{
		"Chrome_OmniboxView", // Chrome/Edge modern
		"Edit",               // Generic edit control
		"RichEdit",           // Rich text edit
		"AddressBar",         // Generic address bar
		"OmniboxViewViews",   // Chrome variants
		"LocationBar",        // Alternative name
	}

	// Enumerate child windows to find address bar
	callback := syscall.NewCallback(func(childHwnd uintptr, lParam uintptr) uintptr {
		className := e.getClassName(childHwnd)

		// Check if this child window is an address bar
		for _, targetClass := range addressBarClasses {
			if strings.Contains(className, targetClass) {
				text := e.getWindowText(childHwnd)
				if text == "" {
					continue
				}
				normalized := normalizeURL(text)
				if isValidURL(normalized) {
					foundURL = normalized
					return 0 // Stop enumeration
				}
			}
		}

		return 1 // Continue enumeration
	})

	procEnumChildWindows.Call(hwnd, callback, 0)
	return foundURL
}

// Win32 API helper functions
func (e *BrowserURLExtractor) getClassName(hwnd uintptr) string {
	buf := make([]uint16, 256)
	ret, _, _ := procGetClassName.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func (e *BrowserURLExtractor) getWindowText(hwnd uintptr) string {
	// First get the length
	length, _, _ := procSendMessage.Call(hwnd, WM_GETTEXTLENGTH, 0, 0)
	if length == 0 {
		return ""
	}

	// Allocate buffer and get text
	buf := make([]uint16, length+1)
	ret, _, _ := procSendMessage.Call(hwnd, WM_GETTEXT, uintptr(len(buf)), uintptr(unsafe.Pointer(&buf[0])))
	if ret == 0 {
		return ""
	}

	return syscall.UTF16ToString(buf)
}

// normalizeURL prepends "https://" when the string looks like a bare domain
// (e.g. "github.com/mihn1/timekeeper") so downstream consumers always get a
// full URL. Returns "" for strings that don't look like URLs at all.
func normalizeURL(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	// Bare domain heuristic: must contain a dot and no spaces
	if strings.Contains(s, ".") && !strings.Contains(s, " ") {
		return "https://" + s
	}
	return ""
}

// isValidURL checks if s is a full URL with http/https scheme.
func isValidURL(s string) bool {
	if s == "" || len(s) < 7 {
		return false
	}
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
