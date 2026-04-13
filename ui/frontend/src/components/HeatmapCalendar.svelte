<script>
  import { createEventDispatcher } from 'svelte';

  export let data = []; // DayActivity[]
  export let year = new Date().getFullYear();

  const dispatch = createEventDispatcher();

  const CAT_COLORS = {
    1: '#3b82f6',
    2: '#f59e0b',
    3: '#22c55e',
    4: '#9ca3af',
  };

  // Build a lookup map: dateStr → DayActivity
  $: activityMap = new Map(data.map(d => [d.date, d]));

  // Max totalMs for intensity scaling
  $: maxMs = data.length > 0 ? Math.max(...data.map(d => d.totalMs)) : 1;

  function intensity(ms) {
    if (!ms || ms === 0) return 0;
    const ratio = ms / maxMs;
    if (ratio < 0.2) return 1;
    if (ratio < 0.4) return 2;
    if (ratio < 0.65) return 3;
    if (ratio < 0.85) return 4;
    return 5;
  }

  // Build weeks array: each week is 7 day slots (Sun–Sat), day = { date, dateStr } or null
  $: weeks = (() => {
    const jan1 = new Date(year, 0, 1);
    const dec31 = new Date(year, 11, 31);

    // Start from the Sunday on or before Jan 1
    const start = new Date(jan1);
    start.setDate(start.getDate() - start.getDay());

    // End at the Saturday on or after Dec 31
    const end = new Date(dec31);
    end.setDate(end.getDate() + (6 - end.getDay()));

    const allWeeks = [];
    let cur = new Date(start);
    while (cur <= end) {
      const week = [];
      for (let d = 0; d < 7; d++) {
        const dt = new Date(cur);
        const inYear = dt.getFullYear() === year;
        const dateStr = dt.toISOString().split('T')[0];
        week.push({ date: dt, dateStr, inYear });
        cur.setDate(cur.getDate() + 1);
      }
      allWeeks.push(week);
    }
    return allWeeks;
  })();

  const MONTHS = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];

  // Month label positions: find the week index where each month starts
  $: monthLabels = (() => {
    const labels = [];
    let lastMonth = -1;
    weeks.forEach((week, wi) => {
      const firstInYear = week.find(d => d.inYear);
      if (!firstInYear) return;
      const m = firstInYear.date.getMonth();
      if (m !== lastMonth) {
        labels.push({ month: MONTHS[m], weekIndex: wi });
        lastMonth = m;
      }
    });
    return labels;
  })();

  const DAYS = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];

  function cellColor(day) {
    if (!day.inYear) return 'transparent';
    const act = activityMap.get(day.dateStr);
    if (!act) return 'var(--table-border-color)';
    const lvl = intensity(act.totalMs);
    if (lvl === 0) return 'var(--table-border-color)';
    const base = CAT_COLORS[act.topCategoryId] ?? '#3b82f6';
    const alpha = 0.2 + lvl * 0.16;
    return base + Math.round(alpha * 255).toString(16).padStart(2, '0');
  }

  function formatTooltip(day) {
    if (!day.inYear) return '';
    const act = activityMap.get(day.dateStr);
    if (!act) return day.dateStr + ': no data';
    const h = Math.floor(act.totalMs / 3600000);
    const m = Math.floor((act.totalMs % 3600000) / 60000);
    return `${day.dateStr}: ${h}h ${m}m`;
  }

  function handleClick(day) {
    if (!day.inYear) return;
    dispatch('daySelected', { date: day.dateStr });
  }
</script>

<div class="heatmap-panel">
  <h3>Activity — {year}</h3>

  <div class="heatmap-scroll">
    <div class="heatmap">
      <!-- Month labels row -->
      <div class="month-row">
        <div class="day-label-spacer"></div>
        {#each weeks as _, wi}
          {@const label = monthLabels.find(l => l.weekIndex === wi)}
          <div class="month-cell">{label ? label.month : ''}</div>
        {/each}
      </div>

      <!-- Day rows (Sun=0 … Sat=6) -->
      {#each [0,1,2,3,4,5,6] as dayIdx}
        <div class="day-row">
          <div class="day-label">{dayIdx % 2 === 1 ? DAYS[dayIdx] : ''}</div>
          {#each weeks as week}
            {@const day = week[dayIdx]}
            <div
              class="cell"
              style="background:{cellColor(day)}"
              title={formatTooltip(day)}
              on:click={() => handleClick(day)}
              role={day.inYear ? 'button' : 'presentation'}
              tabindex={day.inYear ? 0 : -1}
              on:keydown={(e) => e.key === 'Enter' && handleClick(day)}
            ></div>
          {/each}
        </div>
      {/each}
    </div>
  </div>

  <div class="legend">
    <span class="legend-label">Less</span>
    {#each [1,2,3,4,5] as lvl}
      <div class="legend-cell" style="background: #3b82f6{Math.round((0.2 + lvl * 0.16) * 255).toString(16).padStart(2,'0')}"></div>
    {/each}
    <span class="legend-label">More</span>
  </div>
</div>

<style>
  .heatmap-panel {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    box-shadow: var(--card-shadow);
  }

  h3 {
    margin: 0 0 0.75rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .heatmap-scroll { overflow-x: auto; }

  .heatmap {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: max-content;
  }

  .month-row, .day-row {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .day-label-spacer { width: 28px; flex-shrink: 0; }
  .day-label { width: 28px; flex-shrink: 0; font-size: 0.6rem; color: var(--secondary-color); text-align: right; padding-right: 4px; }

  .month-cell { width: 12px; font-size: 0.6rem; color: var(--secondary-color); text-align: left; }

  .cell {
    width: 12px;
    height: 12px;
    border-radius: 2px;
    flex-shrink: 0;
    cursor: default;
    transition: transform 0.1s;
  }

  .cell[role="button"] { cursor: pointer; }
  .cell[role="button"]:hover { transform: scale(1.3); }

  .legend {
    display: flex;
    align-items: center;
    gap: 3px;
    margin-top: 0.5rem;
  }

  .legend-label {
    font-size: 0.65rem;
    color: var(--secondary-color);
    margin: 0 2px;
  }

  .legend-cell {
    width: 10px;
    height: 10px;
    border-radius: 2px;
  }
</style>
