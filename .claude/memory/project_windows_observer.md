---
name: Windows observer implementation details
description: WinEvent hooks, sleep/lock detection, own-process filtering, lock screen exclusion
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
## WinEvent hooks

`platforms/windows/observer.go` installs two hooks:
- `EVENT_SYSTEM_FOREGROUND` — fires when active window changes
- `EVENT_OBJECT_NAMECHANGE` — fires when foreground window title changes (catches browser tab switches)

Both use `WINEVENT_OUTOFCONTEXT` which requires a Win32 message pump (`GetMessageW`/`DispatchMessageW`) on the hook thread.

## Own-process / shell filtering in `collectWindowInfo`

Two filters applied before emitting events:
1. **Class name filter**: Skips "Progman" and "WorkerW" (Desktop Shell windows that briefly become foreground at startup — caused "explorer" to be recorded on launch).
2. **Own-PID filter**: Skips windows belonging to `ownPid` (`os.Getpid()`). `WINEVENT_SKIPOWNPROCESS` only filters the hook thread's own events, not WebView2 renderer events.

## Sleep / lock detection (system:paused marker)

`handleWindowChange` detects machine idle and emits a `constants.SYSTEM_PAUSED` synthetic event:

- **Lock screen** (`lockapp`, `logonui`, `winlogon` executables detected by `isLockScreenApp`): mapped to `SYSTEM_PAUSED` immediately.
- **Sleep/hibernate without lock screen**: Gap between `lastEmitTime` and `now` > `maxIdleGap` (5 minutes) triggers a `SYSTEM_PAUSED` injection at `lastEmitTime + 1ms`.

`SYSTEM_PAUSED` is excluded in `core/timekeeper.go` `getRulesForEvent` — never aggregated, never stored. This closes the previous app's event at the right time and drops the sleep duration.

Observer tracks `isPaused bool` and `lastEmitTime time.Time` to avoid duplicate pause markers.

## App name normalization

`normalizeAppName` maps `.exe` base names → display names:
- `chrome.exe` → `Google Chrome`
- `brave.exe` → `Brave Browser`
- `msedge.exe` → `Microsoft Edge`
- `firefox.exe` → `Firefox`
- others → trimmed lowercase stem

**How to apply:** When debugging "wrong app name" or "sleep time accumulating" issues, start here.
