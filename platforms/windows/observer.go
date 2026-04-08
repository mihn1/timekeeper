//go:build windows
// +build windows

package windows

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/core"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/platforms/windows/browsers"
)

// Ensure interface compliance
var _ core.Observer = (*Observer)(nil)

// Observer captures foreground-window and title-change events via WinEvent hooks
// and converts them into AppSwitchEvents. Start() blocks in a Win32 message loop;
// the caller (core.StartTracking) is expected to run it in a goroutine.
type Observer struct {
	callback     func(models.AppSwitchEvent)
	logger       *slog.Logger
	isStandalone bool
	urlExtractor *browsers.BrowserURLExtractor

	mu        sync.Mutex
	lastApp   string
	lastTitle string

	cb             uintptr
	tid            uint32
	hookForeground windows.Handle
	hookNameChange windows.Handle

	// readyCh is closed once the message loop is running (o.tid is valid).
	// Stop() waits on it before posting WM_QUIT.
	readyCh chan struct{}
	doneCh  chan struct{}
}

func NewObserver(callback func(models.AppSwitchEvent), isStandalone bool, logger *slog.Logger) *Observer {
	if logger == nil {
		logger = slog.Default()
	}
	return &Observer{
		callback:     callback,
		logger:       logger,
		isStandalone: isStandalone,
		urlExtractor: browsers.NewBrowserURLExtractor(logger),
		readyCh:      make(chan struct{}),
		doneCh:       make(chan struct{}),
	}
}

// Start blocks in a Win32 message loop. core.StartTracking calls it via go obs.Start().
func (o *Observer) Start() error {
	o.logger.Info("Starting Windows event observer")
	return o.run()
}

// Stop signals the message loop to exit and waits for it to finish.
func (o *Observer) Stop() error {
	o.logger.Info("Stopping Windows event observer")
	// Wait until the loop is ready (readyCh closed) or has already exited (doneCh closed).
	select {
	case <-o.readyCh:
		if o.tid != 0 {
			procPostThreadMessageW.Call(uintptr(o.tid), uintptr(WM_QUIT), 0, 0)
		}
	case <-o.doneCh:
		// run() exited before the loop started (e.g. hook install failure)
		return nil
	}
	<-o.doneCh
	return nil
}

/* Win32 interop */
var (
	user32                   = windows.NewLazySystemDLL("user32.dll")
	kernel32                 = windows.NewLazySystemDLL("kernel32.dll")
	procSetWinEventHook      = user32.NewProc("SetWinEventHook")
	procUnhookWinEvent       = user32.NewProc("UnhookWinEvent")
	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadPID   = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess          = kernel32.NewProc("OpenProcess")
	procQueryFullImageNameW  = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle          = kernel32.NewProc("CloseHandle")
	procGetAncestor          = user32.NewProc("GetAncestor")
	procGetMessageW          = user32.NewProc("GetMessageW")
	procTranslateMessage     = user32.NewProc("TranslateMessage")
	procDispatchMessageW     = user32.NewProc("DispatchMessageW")
	procPostThreadMessageW   = user32.NewProc("PostThreadMessageW")
	procGetCurrentThreadId   = kernel32.NewProc("GetCurrentThreadId")
)

const (
	EVENT_SYSTEM_FOREGROUND           = 0x0003
	EVENT_OBJECT_NAMECHANGE           = 0x800C
	OBJID_WINDOW                      = 0x00000000
	WINEVENT_OUTOFCONTEXT             = 0x0000
	WINEVENT_SKIPOWNPROCESS           = 0x0002
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	WM_QUIT                           = 0x0012
	GA_ROOT                           = 2 // GetAncestor: root in parent chain
	GA_ROOTOWNER                      = 3 // GetAncestor: root in parent+owner chains
)

var globalObserver *Observer

func (o *Observer) run() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(o.doneCh)

	globalObserver = o
	o.cb = windows.NewCallback(winEventCallback)

	h1, _, _ := procSetWinEventHook.Call(
		uintptr(EVENT_SYSTEM_FOREGROUND), uintptr(EVENT_SYSTEM_FOREGROUND),
		0, o.cb, 0, 0, uintptr(WINEVENT_OUTOFCONTEXT|WINEVENT_SKIPOWNPROCESS),
	)
	h2, _, _ := procSetWinEventHook.Call(
		uintptr(EVENT_OBJECT_NAMECHANGE), uintptr(EVENT_OBJECT_NAMECHANGE),
		0, o.cb, 0, 0, uintptr(WINEVENT_OUTOFCONTEXT|WINEVENT_SKIPOWNPROCESS),
	)
	o.hookForeground = windows.Handle(h1)
	o.hookNameChange = windows.Handle(h2)
	if o.hookForeground == 0 && o.hookNameChange == 0 {
		o.logger.Error("Failed to install WinEvent hooks")
		close(o.readyCh) // unblock any concurrent Stop() call
		return nil
	}

	r, _, _ := procGetCurrentThreadId.Call()
	o.tid = uint32(r)
	close(o.readyCh) // signal: message loop is about to start, o.tid is valid

	var msg struct {
		hwnd   uintptr
		msg    uint32
		wParam uintptr
		lParam uintptr
		time   uint32
		pt     struct{ x, y int32 }
	}
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		switch int32(ret) {
		case -1:
			o.logger.Error("GetMessageW failed")
			o.cleanup()
			return nil
		case 0: // WM_QUIT
			o.cleanup()
			return nil
		default:
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}

func (o *Observer) cleanup() {
	if o.hookForeground != 0 {
		procUnhookWinEvent.Call(uintptr(o.hookForeground))
	}
	if o.hookNameChange != 0 {
		procUnhookWinEvent.Call(uintptr(o.hookNameChange))
	}
	globalObserver = nil
}

func winEventCallback(hWinEventHook uintptr, event uint32, hwnd uintptr, idObject, idChild int32, thread uint32, ms uint32) uintptr {
	o := globalObserver
	if o == nil || hwnd == 0 {
		return 0
	}
	if event == EVENT_OBJECT_NAMECHANGE && idObject != OBJID_WINDOW {
		return 0
	}
	if event == EVENT_OBJECT_NAMECHANGE {
		fg, _, _ := procGetForegroundWindow.Call()
		if fg != hwnd {
			return 0
		}
	}
	appName, exePath, title := collectWindowInfo(hwnd)
	o.handleWindowChange(appName, exePath, title, hwnd)
	return 0
}

// handleWindowChange applies deduplication and fires the callback when the
// active window changes. Extracted from winEventCallback for testability.
func (o *Observer) handleWindowChange(appName, exePath, title string, hwnd uintptr) {
	if appName == "" {
		return
	}

	o.mu.Lock()
	same := appName == o.lastApp && title == o.lastTitle
	if !same {
		o.lastApp, o.lastTitle = appName, title
	}
	o.mu.Unlock()

	if same || o.callback == nil {
		return
	}

	add := map[string]string{constants.KEY_APP_DESC: exePath}
	if title != "" {
		add[constants.KEY_BROWSER_TITLE] = title
	}

	if isBrowser(appName) {
		if url := o.urlExtractor.ExtractURLFromWindow(hwnd, appName, title); url != "" {
			add[constants.KEY_BROWSER_URL] = url
		}
	}

	o.callback(models.AppSwitchEvent{AppName: appName, StartTime: time.Now().UTC(), AdditionalData: add})
}

func collectWindowInfo(hwnd uintptr) (appName, exePath, title string) {
	// GA_ROOTOWNER walks both the parent and owner chains to find the
	// ultimate top-level owning window. This correctly handles:
	//   - Child windows (e.g. embedded controls): parent chain leads up
	//   - Owned top-level windows (e.g. WebView2 host surfaces): owner chain leads up
	// Using GA_ROOT (parent chain only) misses the owned-window case and
	// can return a WebView2 or shell window instead of the real app window.
	root, _, _ := procGetAncestor.Call(hwnd, GA_ROOTOWNER)
	if root == 0 {
		root = hwnd
	}
	title = getWindowTitle(root)
	pid := getWindowPID(root)
	if pid == 0 {
		return "", "", title
	}
	exePath = getProcessImage(pid)
	if exePath == "" {
		return "", "", title
	}
	appName = normalizeAppName(filepath.Base(exePath))

	if globalObserver != nil {
		globalObserver.logger.Debug("collectWindowInfo",
			"hwnd", hwnd,
			"root", root,
			"pid", pid,
			"exe", exePath,
			"app", appName,
			"title", title,
		)
	}

	return appName, exePath, title
}

// normalizeAppName maps a Windows executable base name to a human-readable app name.
func normalizeAppName(base string) string {
	lower := strings.ToLower(base)
	switch lower {
	case "chrome.exe":
		return constants.GOOGLE_CHROME
	case "brave.exe":
		return constants.BRAVE
	case "msedge.exe":
		return constants.MICROSOFT_EDGE
	case "firefox.exe":
		return constants.FIREFOX
	default:
		return trimExt(lower)
	}
}

func getWindowTitle(hwnd uintptr) string {
	l, _, _ := procGetWindowTextLengthW.Call(hwnd)
	n := l + 8 + 1
	if n < 64 {
		n = 64
	}
	buf := make([]uint16, n)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(n))
	return utf16ToString(buf)
}

func getWindowPID(hwnd uintptr) uint32 {
	var pid uint32
	procGetWindowThreadPID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func getProcessImage(pid uint32) string {
	h, _, _ := procOpenProcess.Call(uintptr(PROCESS_QUERY_LIMITED_INFORMATION), 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandle.Call(h)
	buf := make([]uint16, windows.MAX_PATH)
	size := uint32(len(buf))
	r1, _, _ := procQueryFullImageNameW.Call(h, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r1 == 0 {
		return ""
	}
	return utf16ToString(buf[:size])
}

func utf16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[:i]
			break
		}
	}
	return string(utf16.Decode(s))
}

func trimExt(name string) string {
	ext := filepath.Ext(name)
	if ext != "" {
		return name[:len(name)-len(ext)]
	}
	return name
}

// isBrowser reports whether appName corresponds to a tracked web browser.
func isBrowser(appName string) bool {
	switch appName {
	case constants.GOOGLE_CHROME, constants.BRAVE, constants.MICROSOFT_EDGE, constants.FIREFOX:
		return true
	}
	return false
}
