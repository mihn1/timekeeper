package macos

import (
	"fmt"
	"time"

	"github.com/mihn1/timekeeper/internal/core"
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

func (o *Observer) StartObserving(ch chan core.AppSwitchEvent) error {
	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
		fmt.Println("Starting")

		ws := appkit.Workspace_SharedWorkspace()
		notificationCenter := ws.NotificationCenter()
		notificationCenter.AddObserverForNameObjectQueueUsingBlock(
			"NSWorkspaceDidActivateApplicationNotification",
			nil,
			foundation.OperationQueue_MainQueue(),
			func(notification foundation.Notification) {
				event := getEvent(notification)
				ch <- event
			},
		)
	})

	return nil
}

func getEvent(notification foundation.Notification) core.AppSwitchEvent {
	userInfo := notification.UserInfo()
	runningApp := appkit.RunningApplicationFrom(userInfo.ObjectForKey(applicationKey).Ptr())
	appName := runningApp.LocalizedName()
	desc := runningApp.Description()

	return core.AppSwitchEvent{
		AppName:        appName,
		Time:           time.Now().UTC(),
		AdditionalData: desc,
	}
}
