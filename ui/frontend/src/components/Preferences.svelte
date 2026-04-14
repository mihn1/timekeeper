<script>
  import { preferences, timezone } from '../stores/preferences.js';
  import { theme } from '../stores/theme.js';

  // ── Timezone picker ────────────────────────────────────────────────────────

  function getAllTimezones() {
    try {
      return Intl.supportedValuesOf('timeZone');
    } catch {
      // Fallback for environments where supportedValuesOf is unavailable.
      return [
        'UTC',
        'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
        'America/Sao_Paulo', 'Europe/London', 'Europe/Paris', 'Europe/Berlin',
        'Europe/Moscow', 'Asia/Dubai', 'Asia/Kolkata', 'Asia/Bangkok',
        'Asia/Ho_Chi_Minh', 'Asia/Shanghai', 'Asia/Tokyo', 'Asia/Seoul',
        'Australia/Sydney', 'Pacific/Auckland',
      ];
    }
  }

  function getUtcOffset(tz) {
    try {
      const parts = new Intl.DateTimeFormat('en', {
        timeZone: tz,
        timeZoneName: 'shortOffset',
      }).formatToParts(new Date());
      return parts.find(p => p.type === 'timeZoneName')?.value ?? '';
    } catch {
      return '';
    }
  }

  // Build once at module init — ~400 zones, each needing one Intl call.
  const allTzOptions = getAllTimezones().map(tz => ({
    value: tz,
    label: `${tz} (${getUtcOffset(tz)})`,
  }));

  let tzSearch = '';
  let selectedTz = $timezone;
  let tzSaving = false;
  let tzSaveError = '';

  $: filteredOptions = tzSearch.trim()
    ? allTzOptions.filter(o => o.value.toLowerCase().includes(tzSearch.toLowerCase()))
    : allTzOptions;

  // Keep local selection in sync when the store updates (e.g. loaded from backend).
  $: selectedTz = $timezone;

  async function applyTimezone(e) {
    const tz = e.target.value;
    selectedTz = tz;
    tzSaving = true;
    tzSaveError = '';
    try {
      await preferences.setTimezone(tz);
    } catch {
      tzSaveError = 'Failed to save. Will retry next session.';
    } finally {
      tzSaving = false;
    }
  }

  function clearSearch() {
    tzSearch = '';
  }
</script>

<div class="prefs-page">
  <div class="page-header">
    <h1>Preferences</h1>
    <p class="page-subtitle">Customize how TimeKeeper looks and behaves.</p>
  </div>

  <!-- ── Appearance ──────────────────────────────────────────────────────── -->
  <section class="prefs-section">
    <div class="section-header">
      <h2>Appearance</h2>
    </div>

    <div class="pref-row">
      <div class="pref-info">
        <span class="pref-label">Theme</span>
        <span class="pref-desc">Choose between light and dark interface.</span>
      </div>
      <div class="theme-control">
        <button
          class="theme-btn"
          class:active={$theme === 'light'}
          on:click={() => theme.set('light')}
        >
          ☀ Light
        </button>
        <button
          class="theme-btn"
          class:active={$theme === 'dark'}
          on:click={() => theme.set('dark')}
        >
          ☾ Dark
        </button>
      </div>
    </div>
  </section>

  <!-- ── Regional ───────────────────────────────────────────────────────── -->
  <section class="prefs-section">
    <div class="section-header">
      <h2>Regional</h2>
    </div>

    <div class="pref-row tz-row">
      <div class="pref-info">
        <span class="pref-label">Timezone</span>
        <span class="pref-desc">
          Dates and the "today" reference in the dashboard use this timezone.
          Data is always stored as UTC.
        </span>
        <span class="pref-current">
          Current: <strong>{selectedTz}</strong> ({getUtcOffset(selectedTz)})
        </span>
      </div>

      <div class="tz-control">
        <div class="tz-search-wrap">
          <input
            class="tz-search"
            type="text"
            placeholder="Filter timezones…"
            bind:value={tzSearch}
          />
          {#if tzSearch}
            <button class="clear-btn" on:click={clearSearch} aria-label="Clear search">✕</button>
          {/if}
        </div>
        <select
          class="tz-select"
          value={selectedTz}
          on:change={applyTimezone}
          size="6"
        >
          {#each filteredOptions as opt (opt.value)}
            <option value={opt.value}>{opt.label}</option>
          {/each}
        </select>
        <p class="tz-hint">
          {filteredOptions.length} timezone{filteredOptions.length !== 1 ? 's' : ''}
          {tzSearch ? 'matching' : 'available'}
          {#if tzSaving}<span class="tz-saving">Saving…</span>{/if}
          {#if tzSaveError}<span class="tz-error">{tzSaveError}</span>{/if}
        </p>
      </div>
    </div>
  </section>
</div>

<style>
  .prefs-page {
    max-width: 760px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .page-header h1 {
    margin: 0 0 0.25rem;
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--text-color);
  }

  .page-subtitle {
    margin: 0;
    font-size: 0.875rem;
    color: var(--secondary-color);
  }

  /* ── Section card ─────────────────────────────────────────────── */
  .prefs-section {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    box-shadow: var(--card-shadow);
    overflow: hidden;
  }

  .section-header {
    padding: 0.85rem 1.25rem;
    border-bottom: 1px solid var(--card-border-color);
    background-color: var(--table-header-bg);
  }

  .section-header h2 {
    margin: 0;
    font-size: 0.8rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--secondary-color);
  }

  /* ── Preference row ───────────────────────────────────────────── */
  .pref-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 2rem;
    padding: 1.1rem 1.25rem;
    border-bottom: 1px solid var(--card-border-color);
  }

  .pref-row:last-child { border-bottom: none; }

  .pref-info {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    flex: 1;
    min-width: 0;
  }

  .pref-label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .pref-desc {
    font-size: 0.8rem;
    color: var(--secondary-color);
    line-height: 1.45;
  }

  .pref-current {
    font-size: 0.78rem;
    color: var(--secondary-color);
    margin-top: 0.2rem;
  }

  .pref-current strong {
    color: var(--text-color);
  }

  /* ── Theme control ────────────────────────────────────────────── */
  .theme-control {
    display: flex;
    border: 1px solid var(--input-border-color);
    border-radius: 6px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .theme-btn {
    padding: 0.45rem 1rem;
    border: none;
    background: none;
    color: var(--secondary-color);
    font-size: 0.85rem;
    cursor: pointer;
    transition: background-color 0.15s, color 0.15s;
    white-space: nowrap;
  }

  .theme-btn:first-child {
    border-right: 1px solid var(--input-border-color);
  }

  .theme-btn.active {
    background-color: var(--primary-color);
    color: #fff;
    font-weight: 600;
  }

  .theme-btn:not(.active):hover {
    background-color: var(--button-hover-bg-color);
    color: var(--text-color);
  }

  /* ── Timezone control ─────────────────────────────────────────── */
  .tz-row {
    align-items: flex-start;
  }

  .tz-control {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    flex-shrink: 0;
    width: 280px;
  }

  .tz-search-wrap {
    position: relative;
    display: flex;
    align-items: center;
  }

  .tz-search {
    width: 100%;
    padding: 0.4rem 2rem 0.4rem 0.6rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
    box-sizing: border-box;
  }

  .clear-btn {
    position: absolute;
    right: 0.4rem;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--secondary-color);
    font-size: 0.75rem;
    padding: 0.1rem 0.2rem;
    line-height: 1;
  }

  .tz-select {
    width: 100%;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.82rem;
    padding: 0.2rem;
  }

  .tz-hint {
    margin: 0;
    font-size: 0.72rem;
    color: var(--secondary-color);
    text-align: right;
  }

  .tz-saving { color: var(--primary-color); margin-left: 0.4rem; }
  .tz-error  { color: var(--danger-color);  margin-left: 0.4rem; }
</style>
