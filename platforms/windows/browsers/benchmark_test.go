//go:build windows

package browsers

import (
	"testing"

	"github.com/mihn1/timekeeper/constants"
)

// TestURLExtraction is a smoke test that exercises the extractor with hwnd=0.
// With a null HWND both UIA and Win32 methods return empty gracefully.
func TestURLExtraction(t *testing.T) {
	e := NewBrowserURLExtractor(nil)

	browsers := []string{constants.GOOGLE_CHROME, constants.BRAVE, constants.MICROSOFT_EDGE, constants.FIREFOX}
	for _, browser := range browsers {
		url := e.ExtractURLFromWindow(0, browser, "test")
		t.Logf("%s (hwnd=0): %q", browser, url)
		if url != "" && !isValidURL(url) {
			t.Errorf("%s: got non-URL %q", browser, url)
		}
	}
}
