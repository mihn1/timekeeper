<script>
  export let events = [];
  export let categoryUsageData = [];

  // Category color palette keyed by categoryId.
  const CAT_COLORS = {
    0: '#9ca3af', // Excluded — grey
    1: '#3b82f6', // Work — blue
    2: '#f59e0b', // Entertainment — amber
    3: '#22c55e', // Personal — green
    4: '#d1d5db', // Undefined — light grey
  };

  function catColor(categoryId) {
    return CAT_COLORS[categoryId] ?? `hsl(${(categoryId * 137) % 360}, 60%, 55%)`;
  }

  function timeToSecs(timeStr) {
    if (!timeStr || timeStr === '—') return null;
    const parts = timeStr.split(':').map(Number);
    return parts[0] * 3600 + parts[1] * 60 + (parts[2] ?? 0);
  }

  const DAY_SECS = 86400;

  $: rows = (() => {
    const appMap = new Map();
    for (const ev of events) {
      if (!appMap.has(ev.appName)) appMap.set(ev.appName, []);
      appMap.get(ev.appName).push(ev);
    }
    return Array.from(appMap.entries()).map(([appName, evs]) => ({ appName, events: evs }));
  })();

  const HOUR_MARKS = Array.from({ length: 25 }, (_, i) => i);

  function pct(secs) {
    return (secs / DAY_SECS) * 100;
  }

  let tooltip = null;
  let tooltipStyle = '';

  function showTooltip(ev, event, mouseEv) {
    ev.stopPropagation();
    const startSecs = timeToSecs(event.startTime);
    const endSecs   = event.endTime !== '—' ? timeToSecs(event.endTime) : null;
    const dur = event.durationSecs > 0
      ? `${Math.floor(event.durationSecs / 60)}m ${event.durationSecs % 60}s`
      : 'in progress';
    tooltip = {
      app:      event.appName,
      title:    event.urlOrTitle || '',
      start:    event.startTime,
      end:      event.endTime,
      duration: dur,
    };
    const x = mouseEv.clientX + 12;
    const y = mouseEv.clientY + 12;
    tooltipStyle = `left:${x}px;top:${y}px`;
  }

  function hideTooltip() { tooltip = null; }
</script>

<svelte:window on:click={hideTooltip} />

<div class="timeline-panel">
  <h2>Activity Timeline</h2>

  {#if rows.length === 0}
    <div class="empty">No event data for this date.</div>
  {:else}
    <div class="timeline-wrap">
      <!-- Hour ruler -->
      <div class="ruler-row">
        <div class="row-label"></div>
        <div class="ruler">
          {#each HOUR_MARKS as h}
            <span class="hour-mark" style="left:{pct(h * 3600)}%">{String(h).padStart(2,'0')}</span>
          {/each}
        </div>
      </div>

      <!-- App swimlanes -->
      {#each rows as row}
        <div class="swimlane-row">
          <div class="row-label" title={row.appName}>{row.appName}</div>
          <div class="swimlane">
            {#each row.events as event}
              {@const startSecs = timeToSecs(event.startTime)}
              {@const endSecs   = event.endTime !== '—' ? timeToSecs(event.endTime) : Math.min(DAY_SECS, (startSecs ?? 0) + event.durationSecs)}
              {#if startSecs !== null && endSecs !== null && endSecs > startSecs}
                <div
                  class="block"
                  style="left:{pct(startSecs)}%;width:{Math.max(0.2, pct(endSecs - startSecs))}%;background:{catColor(event.categoryId ?? 4)}"
                  on:click|stopPropagation={(e) => showTooltip(e, event, e)}
                  title="{event.appName}"
                ></div>
              {/if}
            {/each}
          </div>
        </div>
      {/each}
    </div>

    {#if tooltip}
      <div class="tooltip" style={tooltipStyle}>
        <div class="tt-app">{tooltip.app}</div>
        {#if tooltip.title}<div class="tt-title">{tooltip.title}</div>{/if}
        <div class="tt-time">{tooltip.start} – {tooltip.end} ({tooltip.duration})</div>
      </div>
    {/if}
  {/if}
</div>

<style>
  .timeline-panel {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: var(--card-shadow);
    overflow: hidden;
  }

  h2 {
    margin: 0 0 1rem;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .timeline-wrap {
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow-x: auto;
  }

  .ruler-row, .swimlane-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .row-label {
    width: 120px;
    min-width: 120px;
    font-size: 0.75rem;
    color: var(--secondary-color);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: right;
    padding-right: 0.5rem;
  }

  .ruler, .swimlane {
    position: relative;
    flex: 1;
    height: 18px;
  }

  .ruler { height: 20px; }

  .hour-mark {
    position: absolute;
    transform: translateX(-50%);
    font-size: 0.65rem;
    color: var(--secondary-color);
    user-select: none;
  }

  .swimlane {
    background-color: var(--table-border-color);
    border-radius: 3px;
    overflow: hidden;
  }

  .block {
    position: absolute;
    top: 1px;
    height: calc(100% - 2px);
    border-radius: 2px;
    opacity: 0.85;
    cursor: pointer;
    transition: opacity 0.1s;
  }

  .block:hover { opacity: 1; }

  .tooltip {
    position: fixed;
    background: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 6px;
    padding: 0.5rem 0.75rem;
    box-shadow: 0 4px 16px rgba(0,0,0,0.15);
    z-index: 9999;
    pointer-events: none;
    max-width: 300px;
  }

  .tt-app   { font-weight: 600; font-size: 0.85rem; color: var(--text-color); }
  .tt-title { font-size: 0.75rem; color: var(--secondary-color); margin-top: 2px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .tt-time  { font-size: 0.75rem; color: var(--secondary-color); margin-top: 4px; }

  .empty {
    color: var(--secondary-color);
    font-style: italic;
    font-size: 0.875rem;
    text-align: center;
    padding: 1rem 0;
  }
</style>
