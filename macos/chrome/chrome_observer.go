package chrome

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#include <stdio.h>
#include <stdlib.h> // Required for free()

// Import function from Objective-C file
void startTabObserver(int pid);

// Forward declare the Go function
// extern void goTabChangeCallback(const char* info);
*/
import "C"
import (
	"log"
	"strings"
	"time"

	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/internal/datatypes"
	"github.com/mihn1/timekeeper/internal/models"
)

var timekeeper *core.TimeKeeper

//export goTabChangeCallback
func goTabChangeCallback(info *C.char) {
	tabInfoRaw := C.GoString(info)
	// log.Printf("CHROME TAB EVENT FROM GO: %s", tabInfo)

	idx := strings.IndexByte(tabInfoRaw, '|')
	if idx == -1 {
		log.Println("Can't parse chrome's tab info")
		return
	}

	url := tabInfoRaw[:idx]
	title := tabInfoRaw[idx+1:]
	tabInfo := datatypes.BrowserTabInfo{
		Title: title,
		URL:   url,
	}

	timekeeper.PushEvent(models.AppSwitchEvent{
		AppName:        constants.GoogleChrome,
		Time:           time.Now(),
		AdditionalData: tabInfo,
	})
}

func StartTabObserver(pid int, t *core.TimeKeeper) {
	if timekeeper == nil {
		timekeeper = t
	}

	log.Println("🚀 Listening for tab changes in Chrome...")
	C.startTabObserver(C.int(pid))
}
