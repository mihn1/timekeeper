package browsers

/*
#include <stdio.h>
#include <stdlib.h>

int startTabObserver(int pid, const char *browserName, const char *script);
void cleanupAllObservers(void);
void cleanupObserverByName(const char *name);
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

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
)

var mu sync.Mutex = sync.Mutex{}
var callbackMap = make(map[string]func(models.AppSwitchEvent))

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

	callbackMap = make(map[string]func(models.AppSwitchEvent))

	C.cleanupAllObservers()
}

//export goTabChangeCallback
func goTabChangeCallback(info *C.char, browserName *C.char) {
	tabInfoRaw := C.GoString(info)
	// log.Printf("TAB EVENT FROM GO: %s", tabInfo)

	idx := strings.IndexByte(tabInfoRaw, '|')
	if idx == -1 {
		return
	}

	url := tabInfoRaw[:idx]
	title := tabInfoRaw[idx+1:]
	tabInfo := map[string]string{
		constants.KEY_BROWSER_URL:   url,
		constants.KEY_BROWSER_TITLE: title,
	}

	if callback, ok := callbackMap[C.GoString(browserName)]; ok {
		callback(models.AppSwitchEvent{
			AppName:        C.GoString(browserName),
			StartTime:      time.Now().UTC(),
			AdditionalData: tabInfo,
		})
	}
}

func StartTabObserver(pid int, browserName string, callback func(models.AppSwitchEvent)) bool {
	log.Printf("🚀 Listening for tab changes in %v...", browserName)

	mu.Lock()
	defer mu.Unlock()
	if _, ok := callbackMap[browserName]; ok {
		log.Printf("Observer already running for %v", browserName)
		return false
	}

	callbackMap[browserName] = callback

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
		return false
	}

	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))

	res := C.startTabObserver(C.int(pid), cBrowserName, cScript)
	return res == 1
}

func StopTabObserver(browserName string) {
	log.Printf("🛑 Stopping tab observer for %v...", browserName)
	mu.Lock()
	defer mu.Unlock()

	delete(callbackMap, browserName)

	cBrowserName := C.CString(browserName)
	defer C.free(unsafe.Pointer(cBrowserName))

	C.cleanupObserverByName(cBrowserName)
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
