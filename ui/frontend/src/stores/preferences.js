import { writable, derived } from 'svelte/store';
import { GetPreferences, SavePreferences } from '../../wailsjs/go/main/App';

const STORAGE_KEY = 'tk_preferences';

const defaultPrefs = {
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
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
    if (backendPrefs?.timezone) {
      _store.update(p => ({ ...p, timezone: backendPrefs.timezone }));
    }
  })
  .catch(() => { /* backend not ready yet — local cache is fine */ });

export const preferences = {
  subscribe: _store.subscribe,

  async setTimezone(tz) {
    _store.update(p => ({ ...p, timezone: tz }));
    try {
      await SavePreferences({ timezone: tz });
    } catch (err) {
      console.error('Failed to save preferences to backend:', err);
    }
  },
};

// Derived convenience store: IANA timezone string.
export const timezone = derived(_store, $p => $p.timezone);
