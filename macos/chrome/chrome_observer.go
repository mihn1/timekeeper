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
	"time"
	"unsafe"

	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/internal/models"
)

var timekeeper *core.TimeKeeper

//export goTabChangeCallback
func goTabChangeCallback(info *C.char) {
	tabInfo := C.GoString(info)
	// log.Printf("CHROME TAB EVENT FROM GO: %s", tabInfo)
	defer C.free(unsafe.Pointer(info)) // Free allocated memory

	timekeeper.PushEvent(models.AppSwitchEvent{
		AppName:        "Google Chrome",
		Time:           time.Now(),
		AdditionalData: tabInfo,
	})
}

func StartTabObserver(pid int, t *core.TimeKeeper) {
	if timekeeper == nil {
		timekeeper = t
	}

	// Start listening for tab changes
	log.Println("🚀 Listening for tab changes in Chrome...")
	C.startTabObserver(C.int(pid))
}
