//go:build windows
// +build windows

package windows

import (
	"sync"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

// ---------------------------------------------------------------------------
// trimExt
// ---------------------------------------------------------------------------

func TestTrimExt(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"chrome.exe", "chrome"},
		{"notepad.exe", "notepad"},
		{"noext", "noext"},
		{"multi.part.exe", "multi.part"},
		{"", ""},
	}
	for _, c := range cases {
		if got := trimExt(c.in); got != c.want {
			t.Errorf("trimExt(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// utf16ToString
// ---------------------------------------------------------------------------

func TestUtf16ToString(t *testing.T) {
	encode := func(s string) []uint16 {
		runes := utf16.Encode([]rune(s))
		return append(runes, 0) // null-terminated
	}

	cases := []string{"hello", "world", "日本語", ""}
	for _, c := range cases {
		if got := utf16ToString(encode(c)); got != c {
			t.Errorf("utf16ToString(%q encoded) = %q, want %q", c, got, c)
		}
	}
}

func TestUtf16ToStringNoNull(t *testing.T) {
	// Slice without a null terminator — should return the full string
	s := utf16.Encode([]rune("abc"))
	if got := utf16ToString(s); got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestUtf16ToStringNullInMiddle(t *testing.T) {
	// Null byte should act as terminator; text after it is ignored
	s := []uint16{'h', 'i', 0, 'x', 'y'}
	if got := utf16ToString(s); got != "hi" {
		t.Errorf("got %q, want %q", got, "hi")
	}
}

// ---------------------------------------------------------------------------
// isBrowser
// ---------------------------------------------------------------------------

func TestIsBrowser(t *testing.T) {
	knownBrowsers := []string{constants.GOOGLE_CHROME, constants.BRAVE, constants.MICROSOFT_EDGE, constants.FIREFOX}
	for _, b := range knownBrowsers {
		if !isBrowser(b) {
			t.Errorf("isBrowser(%q) = false, want true", b)
		}
	}

	nonBrowsers := []string{"notepad", "explorer", "code", "slack", "chrome", ""}
	for _, nb := range nonBrowsers {
		if isBrowser(nb) {
			t.Errorf("isBrowser(%q) = true, want false", nb)
		}
	}
}

// ---------------------------------------------------------------------------
// normalizeAppName
// ---------------------------------------------------------------------------

func TestNormalizeAppName(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"chrome.exe", constants.GOOGLE_CHROME},
		{"Chrome.exe", constants.GOOGLE_CHROME}, // case-insensitive
		{"CHROME.EXE", constants.GOOGLE_CHROME},
		{"brave.exe", constants.BRAVE},
		{"msedge.exe", constants.MICROSOFT_EDGE},
		{"firefox.exe", constants.FIREFOX},
		{"notepad.exe", "notepad"},
		{"NOTEPAD.EXE", "notepad"},
		{"code.exe", "code"},
		{"slack.exe", "slack"},
		{"no_extension", "no_extension"},
	}
	for _, c := range cases {
		if got := normalizeAppName(c.in); got != c.want {
			t.Errorf("normalizeAppName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// NewObserver – initial state
// ---------------------------------------------------------------------------

func TestNewObserverInitialState(t *testing.T) {
	o := NewObserver(nil, true, nil)

	if o.callback != nil {
		t.Error("expected nil callback when none provided")
	}
	if o.lastApp != "" || o.lastTitle != "" {
		t.Error("expected empty lastApp/lastTitle on construction")
	}
	if o.tid != 0 {
		t.Error("expected zero tid before Start")
	}
	if o.readyCh == nil {
		t.Fatal("readyCh should be initialised")
	}
	if o.doneCh == nil {
		t.Fatal("doneCh should be initialised")
	}
	// Channels must not be closed yet
	select {
	case <-o.readyCh:
		t.Error("readyCh should be open before Start")
	default:
	}
	select {
	case <-o.doneCh:
		t.Error("doneCh should be open before Start")
	default:
	}
}

// ---------------------------------------------------------------------------
// handleWindowChange – helpers
// ---------------------------------------------------------------------------

func newTestObserver() (*Observer, *[]models.AppSwitchEvent) {
	var mu sync.Mutex
	var events []models.AppSwitchEvent
	cb := func(e models.AppSwitchEvent) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	}
	o := NewObserver(cb, false, nil)
	return o, &events
}

// collect returns a snapshot of captured events (safe to read after all goroutines finish).
func collect(events *[]models.AppSwitchEvent) []models.AppSwitchEvent {
	return append([]models.AppSwitchEvent(nil), *events...)
}

// ---------------------------------------------------------------------------
// handleWindowChange – basic event firing
// ---------------------------------------------------------------------------

func TestHandleWindowChange_FirstEvent(t *testing.T) {
	o, events := newTestObserver()
	o.handleWindowChange("notepad", `C:\notepad.exe`, "Untitled", 0)

	got := collect(events)
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	e := got[0]
	if e.AppName != "notepad" {
		t.Errorf("AppName = %q, want notepad", e.AppName)
	}
	if e.AdditionalData[constants.KEY_APP_DESC] != `C:\notepad.exe` {
		t.Errorf("KEY_APP_DESC = %q", e.AdditionalData[constants.KEY_APP_DESC])
	}
	if e.AdditionalData[constants.KEY_BROWSER_TITLE] != "Untitled" {
		t.Errorf("KEY_BROWSER_TITLE = %q", e.AdditionalData[constants.KEY_BROWSER_TITLE])
	}
}

func TestHandleWindowChange_EmptyTitle_KeyAbsent(t *testing.T) {
	o, events := newTestObserver()
	o.handleWindowChange("notepad", `C:\notepad.exe`, "", 0)

	got := collect(events)
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if _, ok := got[0].AdditionalData[constants.KEY_BROWSER_TITLE]; ok {
		t.Error("KEY_BROWSER_TITLE should be absent when title is empty")
	}
}

func TestHandleWindowChange_EventTimestamp(t *testing.T) {
	o, events := newTestObserver()
	before := time.Now().UTC()
	o.handleWindowChange("notepad", `C:\notepad.exe`, "Untitled", 0)
	after := time.Now().UTC()

	got := collect(events)
	if len(got) != 1 {
		t.Fatal("expected 1 event")
	}
	ts := got[0].StartTime
	if ts.Before(before) || ts.After(after) {
		t.Errorf("StartTime %v not in window [%v, %v]", ts, before, after)
	}
}

// ---------------------------------------------------------------------------
// handleWindowChange – deduplication
// ---------------------------------------------------------------------------

func TestHandleWindowChange_Deduplication_ExactRepeat(t *testing.T) {
	o, events := newTestObserver()

	for i := 0; i < 5; i++ {
		o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0)
	}

	if got := collect(events); len(got) != 1 {
		t.Errorf("expected 1 event after repeated identical calls, got %d", len(got))
	}
}

func TestHandleWindowChange_AppSwitch_FiresNewEvent(t *testing.T) {
	o, events := newTestObserver()

	o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0)
	o.handleWindowChange("code", `C:\code.exe`, "main.go", 0)

	got := collect(events)
	if len(got) != 2 {
		t.Fatalf("expected 2 events on app switch, got %d", len(got))
	}
	if got[1].AppName != "code" {
		t.Errorf("second event AppName = %q, want code", got[1].AppName)
	}
}

func TestHandleWindowChange_TitleChange_SameApp_FiresNewEvent(t *testing.T) {
	o, events := newTestObserver()

	o.handleWindowChange("notepad", `C:\notepad.exe`, "file1.txt", 0)
	o.handleWindowChange("notepad", `C:\notepad.exe`, "file2.txt", 0)

	if got := collect(events); len(got) != 2 {
		t.Errorf("expected 2 events on title change, got %d", len(got))
	}
}

func TestHandleWindowChange_ReturnToSameApp_FiresEvent(t *testing.T) {
	// A→B→A must fire 3 events (not deduplicated on the return)
	o, events := newTestObserver()

	o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0) // A
	o.handleWindowChange("code", `C:\code.exe`, "main.go", 0)       // B
	o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0) // A again

	got := collect(events)
	if len(got) != 3 {
		t.Errorf("expected 3 events for A→B→A, got %d", len(got))
	}
}

func TestHandleWindowChange_ExePathChange_SameAppTitle_NotDuplicated(t *testing.T) {
	// Same app name and title but different exe path: dedup is on (appName, title),
	// so this IS considered a duplicate and does NOT fire a second event.
	o, events := newTestObserver()

	o.handleWindowChange("notepad", `C:\Windows\notepad.exe`, "doc.txt", 0)
	o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0)

	got := collect(events)
	if len(got) != 1 {
		t.Errorf("expected 1 event (dedup on app+title, not path), got %d", len(got))
	}
}

// ---------------------------------------------------------------------------
// handleWindowChange – edge cases
// ---------------------------------------------------------------------------

func TestHandleWindowChange_EmptyAppName_NoEvent(t *testing.T) {
	o, events := newTestObserver()
	o.handleWindowChange("", `C:\unknown.exe`, "Title", 0)

	if got := collect(events); len(got) != 0 {
		t.Errorf("expected no event for empty app name, got %d", len(got))
	}
}

func TestHandleWindowChange_NilCallback_NoPanic(t *testing.T) {
	o := NewObserver(nil, false, nil)
	// Must not panic regardless of app name
	o.handleWindowChange("notepad", `C:\notepad.exe`, "Untitled", 0)
}

func TestHandleWindowChange_NilCallbackBrowser_NoPanic(t *testing.T) {
	o := NewObserver(nil, false, nil)
	// Browser path also hits URL extraction; with nil callback must not panic
	o.handleWindowChange(constants.GOOGLE_CHROME, `C:\chrome.exe`, "Google", 0)
}

func TestHandleWindowChange_BrowserApp_EventFires(t *testing.T) {
	// Verify that browser apps fire events just like other apps.
	// We do NOT assert on KEY_BROWSER_URL because hwnd=0 causes
	// EnumChildWindows(NULL) to scan real windows — the result is
	// environment-dependent and must not be asserted in a unit test.
	o, events := newTestObserver()
	o.handleWindowChange(constants.GOOGLE_CHROME, `C:\chrome.exe`, "Google - Chrome", 0)

	got := collect(events)
	if len(got) != 1 {
		t.Fatalf("expected 1 event for browser app, got %d", len(got))
	}
	if got[0].AppName != constants.GOOGLE_CHROME {
		t.Errorf("AppName = %q, want %q", got[0].AppName, constants.GOOGLE_CHROME)
	}
	if got[0].AdditionalData[constants.KEY_BROWSER_TITLE] != "Google - Chrome" {
		t.Errorf("KEY_BROWSER_TITLE = %q", got[0].AdditionalData[constants.KEY_BROWSER_TITLE])
	}
}

// ---------------------------------------------------------------------------
// handleWindowChange – concurrency safety
// ---------------------------------------------------------------------------

func TestHandleWindowChange_Concurrent_NoRaceNoDeadlock(t *testing.T) {
	o, events := newTestObserver()

	// Each goroutine uses a distinct app name so dedup never suppresses events.
	apps := []string{"notepad", "code", "slack", "explorer", "winword"}
	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			o.handleWindowChange(a, `C:\`+a+`.exe`, "title", 0)
		}(app)
	}
	wg.Wait()

	// With 5 distinct app names and no prior state, at least 1 event must fire.
	// The exact count is scheduling-dependent but should be > 0.
	got := collect(events)
	if len(got) == 0 {
		t.Error("expected at least one event from concurrent calls")
	}
}

func TestHandleWindowChange_Concurrent_SameApp_AtMostOne(t *testing.T) {
	o, events := newTestObserver()

	// All goroutines submit the same (app, title): dedup means ≤1 event.
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.handleWindowChange("notepad", `C:\notepad.exe`, "doc.txt", 0)
		}()
	}
	wg.Wait()

	got := collect(events)
	if len(got) > 1 {
		t.Errorf("concurrent dedup should produce ≤1 event, got %d", len(got))
	}
}
