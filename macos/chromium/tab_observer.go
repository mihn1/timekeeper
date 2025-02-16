package chromium

/*
#include <stdio.h>
#include <stdlib.h>

void startTabObserver(int ßpid, char * name);
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

func StartTabObserver(pid int, browser string, t *core.TimeKeeper) {
	mu.Lock()
	if timekeeper == nil {
		timekeeper = t
	}
	mu.Unlock()

	log.Printf("🚀 Listening for tab changes in %v...", browser)

	cBrowser := C.CString(browser)
	defer C.free(unsafe.Pointer(cBrowser))

	C.startTabObserver(C.int(pid), cBrowser)
}

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
