---
name: UI dashboard and preferences features
description: Components, stores, routes, and dashboard state as of current implementation
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
## Navigation routes (Menu.svelte → App.svelte)

`dashboard` | `rules` | `categories` | `events` | `preferences`

Menu labels: Dashboard, Rules, Categories, Event Log, ⚙ Preferences

## Dashboard features

**Date navigation**: prev/next/today buttons + click-to-edit date display. Forward button disabled when `selectedDate === today`. Date helpers in `utils/dateUtils.js` (timezone-aware, no UTC-offset bug).

**View modes**: Day / 7d / 14d / 30d toggle.

**Min duration filter**: All / ≥1m / ≥5m / ≥15m dropdown.

**Components rendered in day mode:**
- `DaySummaryBar` — total tracked, productive time, top app, app switch count
- `ComparisonStrip` — delta vs yesterday and 7-day average
- `GoalsPanel` — progress bars per category goal, inline add/delete
- Chart row: `AppUsageChart` + `CategoryChart` (click-to-filter by category)
- `ActivityTimeline` — swimlane chart, 00:00–24:00 ruler, hover tooltip
- App details table (filterable, sortable)
- `TopDomainsPanel` + `UncategorizedAppsPanel` (triggers CreateRuleModal with prefill)
- `HeatmapCalendar` — always visible, click cell → navigate to that day

**Range mode (7d/14d/30d):** shows `TrendChart` (line chart per category) + calendar.

**Category drill-down:** clicking a category in CategoryChart sets `selectedCategoryFilter`, dims other bars, filters app table.

**CreateRuleModal:** launched from UncategorizedAppsPanel with `prefillAppName`. Categories passed in from Dashboard's own `GetCategories()` fetch.

## Preferences page features

Sections (extensible card layout):
1. **Appearance** — Light/Dark theme segmented toggle
2. **Regional** — Timezone search + scrollable select (size=6), shows UTC offset, saves to backend on change

## Stores

| Store | File | Purpose |
|-------|------|---------|
| `currentView` | `stores/navigation.js` | Active route |
| `trackingEnabled` | `stores/timekeeper.js` | Observer on/off |
| `refreshData` | `stores/timekeeper.js` | Triggers data reload |
| `theme` | `stores/theme.js` | 'light' \| 'dark', persisted to localStorage |
| `preferences` | `stores/preferences.js` | `{timezone}`, synced with backend |
| `timezone` | `stores/preferences.js` (derived) | Current IANA string |

**How to apply:** Check this when locating which component owns which feature or adding a new dashboard panel.
