---
name: Svelte reactive declarations not available at variable init time
description: Don't initialize let variables from $: reactive labels — use store get or $store directly
type: feedback
originSessionId: 24a106a6-d389-4c1f-b68e-f7633a0973cd
---
In Svelte, `$: foo = expr` is a reactive label, not an eagerly-evaluated const. When another `let bar = foo` runs at component init, `foo` is still `undefined`.

**Caused:** blank page crash when `let selectedDate = today` and `$: today = todayInTz($timezone)` — `selectedDate` was `undefined`, then `selectedDate.split('-')` threw.

**Fix pattern:**
```js
// WRONG
$: today = todayInTz($timezone);
let selectedDate = today; // undefined!

// CORRECT
$: today = todayInTz($timezone);
let selectedDate = todayInTz($timezone); // $store IS available at init
```

`$store` (auto-subscription shorthand) IS available synchronously at component initialization because the store subscription runs first.

**How to apply:** Never derive initial `let` values from `$:` labels. Use `$store` directly or `get(store)` from `svelte/store`.
