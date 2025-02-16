package browsers

/*
#include <stdio.h>
#include <stdlib.h>

void startTabObserver(int pid, const char *browserName, const char *script);
void cleanupAllObservers(void);
*/
import "C"
import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/internal/models"
)

var timekeeper *core.TimeKeeper
var mu sync.Mutex = sync.Mutex{}

func init() {
	// Set up signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
}

func cleanup() {
	mu.Lock()
	defer mu.Unlock()

	// Call Objective-C cleanup
	C.cleanupAllObservers()
}

//export goTabChangeCallback
func goTabChangeCallback(info *C.char, browserName *C.char) {
	tabInfoRaw := C.GoString(info)
	// log.Printf("TAB EVENT FROM GO: %s", tabInfo)

	idx := strings.IndexByte(tabInfoRaw, '|')
	if idx == -1 {
		log.Println("Can't parse chrome's tab info")
		return
	}

	url := tabInfoRaw[:idx]
	title := tabInfoRaw[idx+1:]
	tabInfo := map[string]string{
		constants.KEY_BROWSER_URL:   url,
		constants.KEY_BROWSER_TITLE: title,
	}

	timekeeper.PushEvent(models.AppSwitchEvent{
		AppName:        C.GoString(browserName),
		StartTime:      time.Now().UTC(),
		AdditionalData: tabInfo,
	})
}

func StartTabObserver(pid int, browserName string, t *core.TimeKeeper) {
	mu.Lock()
	if timekeeper == nil {
		timekeeper = t
	}
	mu.Unlock()

	log.Printf("🚀 Listening for tab changes in %v...", browserName)

	cBrowserName := C.CString(browserName)
	defer C.free(unsafe.Pointer(cBrowserName))

	var script string
	switch browserName {
	case constants.GOOGLE_CHROME, constants.BRAVE:
		script = chromiumScript
	case constants.SAFARI:
		script = safariScript
	default:
		log.Printf("Unsupported browser: %v", browserName)
		return
	}

	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))

	C.startTabObserver(C.int(pid), cBrowserName, cScript)
}

const (
	safariScript = `
tell application "Safari"
	set frontTab to current tab of front window
	set tabTitle to name of frontTab
	set tabURL to URL of frontTab
	return tabURL & "|" & tabTitle
end tell
`
	chromiumScript = `
tell application "Google Chrome"
	set frontTab to active tab of front window
	set tabTitle to title of frontTab
	set tabURL to URL of frontTab
	return tabURL & "|" & tabTitle
end tell
`
)
