---
name: Rules, categories, and exclusion logic
description: How rules match events, exclusion short-circuit, and validation fixes applied
type: project
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
## Rule matching

`internal/models/rule.go` — `IsMatch(event)`:
- If no `AdditionalDataKey`/`Expression`: match by app name only
- Else evaluate `event.AdditionalData[key]` against `Expression` (plain substring or regex if `IsRegex`)
- `Priority` determines ordering (higher first)
- `IsExclusion` rules short-circuit processing — matched event is dropped entirely

Resolution order in `core/timekeeper.go`:
1. App-specific rules
2. Exclusion check → `EventExcludedError`
3. Global rules (`ALL_APPS`)
4. Exclusion check
5. First matching non-exclusion rule → category
6. Fallback: `UNDEFINED`

## Hardcoded exclusions (no rule required)

- `constants.SYSTEM_PAUSED` — synthetic sleep/lock marker, excluded in `getRulesForEvent` before any rule lookup

## Validation fixes applied

**Frontend (CreateRuleModal.svelte):** category not required when `isExclusion` is ticked.
**Backend (api_rules.go `AddRule`/`UpdateRule`):** `CategoryID <= 0 && !IsExclusion` — same condition, both fixed.

## TestRuleMatch API

`api_rules.go TestRuleMatch(appName, additionalDataKey, value)`:
- Builds synthetic event
- Fetches app-specific + global rules, sorts by priority desc
- Uses `DefaultCategoryResolver` to resolve category
- Returns `RuleMatchResult{Matched, CategoryId, CategoryName, MatchedRule}`

Exposed as `RuleTestPanel.svelte` in the Rules view.

## Category IDs (seeded defaults)

| ID | Name |
|----|------|
| 0 | Excluded |
| 1 | Work |
| 2 | Entertainment |
| 3 | Personal |
| 4 | Undefined (fallback) |

**How to apply:** When debugging "event not tracked" or "wrong category", trace through this resolution order.
