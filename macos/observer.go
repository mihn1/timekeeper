package macos

import (
	"fmt"
	"time"

	"github.com/mihn1/timekeeper/internal/core"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/macos/chrome"
	"github.com/progrium/darwinkit/macos"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
)

type Observer struct{}

func NewObserver() *Observer {
	return &Observer{}
}

var (
	applicationKey foundation.String = foundation.String_StringWithString("NSWorkspaceApplicationKey")
)

func (o *Observer) StartObserving(t *core.TimeKeeper) error {
	appListeners := make(map[string]bool) // App name -> isListening

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

				if event.AppName == "Google Chrome" {
					if !appListeners[event.AppName] {
						chrome.StartTabObserver(pid, t)
						appListeners[event.AppName] = true
					}
				}

				t.PushEvent(event)
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
		Time:           time.Now().UTC(),
		AdditionalData: desc,
	}

	return event, int(runningApp.ProcessIdentifier())
}
