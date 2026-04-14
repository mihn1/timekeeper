---
name: Never use toISOString() for local date strings
description: toISOString() returns UTC date which is wrong for UTC+ timezones — use localDateStr or todayInTz
type: feedback
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
`new Date().toISOString().split('T')[0]` returns the **UTC date**, not the local date.

In UTC+ timezones (e.g. UTC+7), local midnight = previous UTC day. So:
- `new Date('2026-04-14T00:00:00').toISOString()` → `"2026-04-13T17:00:00Z"` → `"2026-04-13"` (WRONG)
- Forward date shift appeared to do nothing (same UTC date returned)

**Why:** User is in a UTC+ timezone. Forward navigation didn't move the date. Backward appeared to skip days.

**Fix:** Use `utils/dateUtils.js`:
- `todayInTz(tz)` — uses `Intl.DateTimeFormat('en-CA', {timeZone})` which returns YYYY-MM-DD in correct timezone
- `shiftDateStr(dateStr, days)` — shifts using local `getDate()`/`setDate()`, returns via `localDateStr()` (not toISOString)
- `formatDateDisplay(dateStr, tz)` — formats for display using `T12:00:00` noon to avoid DST edges

**How to apply:** Any time you write a date string for display or API calls, use `utils/dateUtils.js` helpers, never `toISOString().split('T')[0]`.
