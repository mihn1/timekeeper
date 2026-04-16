import { writable, derived } from 'svelte/store';
import { GetPreferences, SavePreferences } from '../../wailsjs/go/main/App';

const STORAGE_KEY = 'tk_preferences';

const defaultPrefs = {
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  minEventDurationMs: 1000,
};

function loadLocalPrefs() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return { ...defaultPrefs, ...JSON.parse(raw) };
  } catch {}
  return { ...defaultPrefs };
}

const _store = writable(loadLocalPrefs());

// Mirror to localStorage for instant reads before backend responds.
_store.subscribe(val => {
  try { localStorage.setItem(STORAGE_KEY, JSON.stringify(val)); } catch {}
});

// Sync from backend on startup — backend is the source of truth.
GetPreferences()
  .then(backendPrefs => {
    if (!backendPrefs) return;
    _store.update(p => ({
      ...p,
      ...(backendPrefs.timezone           ? { timezone: backendPrefs.timezone }                       : {}),
      ...(backendPrefs.minEventDurationMs != null ? { minEventDurationMs: backendPrefs.minEventDurationMs } : {}),
    }));
  })
  .catch(() => { /* backend not ready yet — local cache is fine */ });

export const preferences = {
  subscribe: _store.subscribe,

  async setTimezone(tz) {
    _store.update(p => ({ ...p, timezone: tz }));
    try {
      let current;
      _store.subscribe(v => { current = v; })();
      await SavePreferences({ timezone: tz, minEventDurationMs: current.minEventDurationMs });
    } catch (err) {
      console.error('Failed to save preferences to backend:', err);
    }
  },

  async setMinEventDurationMs(ms) {
    const clamped = Math.max(0, Math.round(ms));
    _store.update(p => ({ ...p, minEventDurationMs: clamped }));
    try {
      let current;
      _store.subscribe(v => { current = v; })();
      await SavePreferences({ timezone: current.timezone, minEventDurationMs: clamped });
    } catch (err) {
      console.error('Failed to save preferences to backend:', err);
    }
  },
};

// Derived convenience stores.
export const timezone = derived(_store, $p => $p.timezone);
export const minEventDurationMs = derived(_store, $p => $p.minEventDurationMs);
