package macos

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/mihn1/timekeeper/constants"
	"github.com/mihn1/timekeeper/internal/models"
	"github.com/mihn1/timekeeper/macos/browsers"
	"github.com/progrium/darwinkit/macos"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
)

type Observer struct {
	browserListeners map[string]bool
	mu               sync.Mutex
	callback         func(models.AppSwitchEvent)
	logger           *slog.Logger
}

func NewObserver(callback func(models.AppSwitchEvent), logger *slog.Logger) *Observer {
	return &Observer{
		browserListeners: make(map[string]bool),
		mu:               sync.Mutex{},
		callback:         callback,
		logger:           logger,
	}
}

var (
	applicationKey foundation.String = foundation.String_StringWithString("NSWorkspaceApplicationKey")
)

func (o *Observer) Start() error {
	o.logger.Info("Starting MACOS browsers observer")
	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
		ws := appkit.Workspace_SharedWorkspace()
		notificationCenter := ws.NotificationCenter()

		// Register for launching a new app
		notificationCenter.AddObserverForNameObjectQueueUsingBlock(
			"NSWorkspaceDidLaunchApplicationNotification",
			nil,
			foundation.OperationQueue_MainQueue(),
			func(notification foundation.Notification) {
				event, pid := getEvent(notification)
				o.registerBrowserObserver(pid, event.AppName, o.callback)
			})

		// Register for activating an app
		notificationCenter.AddObserverForNameObjectQueueUsingBlock(
			"NSWorkspaceDidActivateApplicationNotification",
			nil,
			foundation.OperationQueue_MainQueue(),
			func(notification foundation.Notification) {
				event, pid := getEvent(notification)
				o.registerBrowserObserver(pid, event.AppName, o.callback)
				o.callback(event) // Push event to timekeeper in case of app activation
			})

		// Register for terminating an app
		notificationCenter.AddObserverForNameObjectQueueUsingBlock(
			"NSWorkspaceDidTerminateApplicationNotification",
			nil,
			foundation.OperationQueue_MainQueue(),
			func(notification foundation.Notification) {
				event, _ := getEvent(notification)
				o.stopBrowserObserver(event.AppName)
			})
	})

	return nil
}

func (o *Observer) Stop() error {
	o.logger.Info("Stopping MACOS browsers observer")
	for browserName := range o.browserListeners {
		o.stopBrowserObserver(browserName)
	}

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

func (o *Observer) registerBrowserObserver(pid int, browserName string, callback func(models.AppSwitchEvent)) {
	switch browserName {
	case constants.BRAVE, constants.GOOGLE_CHROME, constants.SAFARI:
		o.mu.Lock()
		defer o.mu.Unlock()

		if val, ok := o.browserListeners[browserName]; !ok || !val {
			log.Printf("Registering browser observer for %v", browserName)
			success := browsers.StartTabObserver(pid, browserName, callback)
			if !success {
				log.Printf("Failed to start observer for %v", browserName)
				return
			}

			o.browserListeners[browserName] = true
		}
	}
}

func (o *Observer) stopBrowserObserver(browserName string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, ok := o.browserListeners[browserName]; !ok {
		return
	}

	browsers.StopTabObserver(browserName)
	o.browserListeners[browserName] = false
}
