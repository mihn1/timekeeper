package macos

import (
	"fmt"
	"sync"
	"time"

	"github.com/mihn1/timekeeper/internal/constants"
	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/macos/chromium"
	"github.com/progrium/darwinkit/macos"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
)

type Observer struct {
	timekeeper       *core.TimeKeeper
	browserListeners map[string]bool
	mu               sync.Mutex
}

func NewObserver(t *core.TimeKeeper) *Observer {
	return &Observer{
		timekeeper:       t,
		browserListeners: make(map[string]bool),
		mu:               sync.Mutex{},
	}
}

var (
	applicationKey foundation.String = foundation.String_StringWithString("NSWorkspaceApplicationKey")
)

func (o *Observer) StartObserving() error {
	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
		fmt.Println("Starting")

		ws := appkit.Workspace_SharedWorkspace()
		notificationCenter := ws.NotificationCenter()
		notificationCenter.AddObserverForNameObjectQueueUsingBlock(
			"NSWorkspaceDidActivateApplicationNotification",
			nil,
			foundation.OperationQueue_MainQueue(),
			func(notification foundation.Notification) {
				event, pid := getEvent(notification)

				if event.AppName == constants.GOOGLE_CHROME {
					o.registerChromiumObserver(pid, constants.GOOGLE_CHROME)
				} else if event.AppName == constants.BRAVE {
					o.registerChromiumObserver(pid, constants.BRAVE)
				}

				o.timekeeper.PushEvent(event)
			},
		)
	})

	return nil
}

func getEvent(notification foundation.Notification) (models.AppSwitchEvent, int) {
	userInfo := notification.UserInfo()
	runningApp := appkit.RunningApplicationFrom(userInfo.ObjectForKey(applicationKey).Ptr())
	appName := runningApp.LocalizedName()
	desc := runningApp.Description()

	event := models.AppSwitchEvent{
		AppName:        appName,
		StartTime:      time.Now().UTC(),
		AdditionalData: map[string]string{constants.KEY_APP_DESC: desc},
	}

	return event, int(runningApp.ProcessIdentifier())
}

func (o *Observer) registerChromiumObserver(pid int, browserName string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.browserListeners[browserName] {
		chromium.StartTabObserver(pid, browserName, o.timekeeper)
		o.browserListeners[browserName] = true
	}
}
