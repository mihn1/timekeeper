package constants

// Key names for regular additional data
const (
	KEY_BROWSER_URL   = "url"
	KEY_BROWSER_TITLE = "title"
	KEY_APP_DESC      = "description"
)

const (
	ALL_APPS = "All_Apps"

	// SYSTEM_PAUSED is a synthetic app name emitted when the machine is locked or
	// sleeping. Events with this name are always excluded from aggregations.
	SYSTEM_PAUSED = "system:paused"

	// AppNameChrome is the name of the Chrome application
	GOOGLE_CHROME  = "Google Chrome"
	MICROSOFT_EDGE = "Microsoft Edge"
	FIREFOX        = "Firefox"
	BRAVE          = "Brave Browser"
	SAFARI         = "Safari"
)
