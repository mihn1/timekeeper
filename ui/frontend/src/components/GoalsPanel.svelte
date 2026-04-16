<script>
  import { GetGoals, GetCategories, AddGoal, UpdateGoal, DeleteGoal, GetCategoryUsageTotals } from '../../wailsjs/go/main/App';
  import { formatTimeElapsed } from '../utils/formatters';
  import { shiftDateStr } from '../utils/dateUtils';
  import { onMount } from 'svelte';

  export let selectedDate = '';
  export let refreshTick = 0;

  const FREQ = [
    { value: 1, label: 'Daily' },
    { value: 2, label: 'Weekly' },
    { value: 3, label: 'Monthly' },
  ];
  function freqLabel(f) { return FREQ.find(x => x.value === f)?.label ?? 'Daily'; }
  function msToHm(ms) { return { h: Math.floor(ms / 3600000), m: Math.floor((ms % 3600000) / 60000) }; }
  function hmToMs(h, m) { return (Number(h) * 3600 + Number(m) * 60) * 1000; }

  function calendarWeekStart(dateStr) {
    const d = new Date(dateStr + 'T12:00:00');
    const dow = d.getDay();
    const toMon = (dow === 0) ? -6 : 1 - dow;
    return shiftDateStr(dateStr, toMon);
  }
  function calendarMonthStart(dateStr) { return dateStr.slice(0, 8) + '01'; }

  let goals = [];
  let categories = [];
  let dailyUsage = [];
  let weeklyUsage = [];
  let monthlyUsage = [];
  let showAddForm = false;
  let isLoading = true;

  // Add form state
  let newName = '';
  let newFreq = 1;
  let newCategoryIds = [];
  let newHours = 1;
  let newMinutes = 0;
  let addError = null;

  $: weekStart  = selectedDate ? calendarWeekStart(selectedDate)  : '';
  $: monthStart = selectedDate ? calendarMonthStart(selectedDate) : '';

  $: activeGoals = goals.filter(g => g.isActive);
  $: availableCategories = categories.filter(c => c.id > 0);

  $: dailyGoals   = activeGoals.filter(g => g.frequency === 1);
  $: weeklyGoals  = activeGoals.filter(g => g.frequency === 2);
  $: monthlyGoals = activeGoals.filter(g => g.frequency === 3);

  const SECTIONS = [
    { label: 'Daily',   goals: () => dailyGoals,   sublabel: () => selectedDate },
    { label: 'Weekly',  goals: () => weeklyGoals,  sublabel: () => `${weekStart} – ${selectedDate}` },
    { label: 'Monthly', goals: () => monthlyGoals, sublabel: () => `${monthStart} – ${selectedDate}` },
  ];

  onMount(async () => {
    await loadGoals();
    await loadPeriodData();
  });

  // Reload when date changes (navigation) or dashboard refreshes.
  $: if (selectedDate) { console.log('[GoalsPanel] selectedDate changed:', selectedDate); loadPeriodData(); }
  $: if (refreshTick > 0) { console.log('[GoalsPanel] refreshTick changed:', refreshTick); loadPeriodData(); }

  async function loadGoals() {
    isLoading = true;
    try {
      [goals, categories] = await Promise.all([GetGoals(), GetCategories()]);
    } catch (err) {
      console.error('Error loading goals:', err);
    } finally {
      isLoading = false;
    }
  }

  async function loadPeriodData() {
    if (!selectedDate) return;
    const date = selectedDate; // snapshot to detect stale results
    console.log('[GoalsPanel] loadPeriodData start, date=', date);
    try {
      const [d, w, m] = await Promise.all([
        GetCategoryUsageTotals(date, date),
        GetCategoryUsageTotals(calendarWeekStart(date), date),
        GetCategoryUsageTotals(calendarMonthStart(date), date),
      ]);
      if (date !== selectedDate) {
        console.log('[GoalsPanel] loadPeriodData stale, discarding. date=', date, 'selectedDate=', selectedDate);
        return;
      }
      console.log('[GoalsPanel] loadPeriodData done, date=', date, 'd=', d);
      dailyUsage   = d ?? [];
      weeklyUsage  = w ?? [];
      monthlyUsage = m ?? [];
    } catch (err) {
      console.error('[GoalsPanel] Error loading period usage:', err);
    }
  }

  function usageForGoal(goal) {
    const usageData = goal.frequency === 2 ? weeklyUsage
                    : goal.frequency === 3 ? monthlyUsage
                    : dailyUsage;
    return (goal.categoryIds ?? []).reduce((sum, cid) => {
      const c = (usageData ?? []).find(x => x.id === cid);
      return sum + (c ? (c.timeElapsed ?? 0) : 0);
    }, 0);
  }

  function progressPct(actual, target) {
    if (!target || target <= 0) return 0;
    return Math.min(100, Math.round((actual / target) * 100));
  }

  function progressClass(pct) {
    if (pct >= 100) return 'met';
    if (pct >= 50)  return 'close';
    return 'low';
  }

  async function removeGoal(goal) {
    try {
      await DeleteGoal(goal.id);
      await loadGoals();
    } catch (err) {
      console.error('Error deleting goal:', err);
    }
  }

  async function addGoal() {
    if (!newName.trim()) { addError = 'Name required.'; return; }
    if (!newCategoryIds.length) { addError = 'Select at least one category.'; return; }
    const ms = hmToMs(newHours, newMinutes);
    if (ms <= 0) { addError = 'Target must be > 0.'; return; }
    addError = null;
    try {
      await AddGoal(newName.trim(), newCategoryIds.map(Number), newFreq, ms);
      newName = '';
      newFreq = 1;
      newCategoryIds = [];
      newHours = 1;
      newMinutes = 0;
      showAddForm = false;
      await loadGoals();
      await loadPeriodData();
    } catch (err) {
      console.error('Error adding goal:', err);
      addError = 'Failed to save.';
    }
  }

  function toggleCategory(id) {
    const numId = Number(id);
    if (newCategoryIds.includes(numId)) {
      newCategoryIds = newCategoryIds.filter(x => x !== numId);
    } else {
      newCategoryIds = [...newCategoryIds, numId];
    }
  }
</script>

<div class="goals-panel">
  <div class="panel-header">
    <h3>Goals</h3>
    <button class="add-btn" on:click={() => { showAddForm = !showAddForm; addError = null; }} title="Add goal">
      {showAddForm ? '✕' : '+'}
    </button>
  </div>

  {#if showAddForm}
    <div class="add-form">
      <input class="name-input" type="text" placeholder="Goal name…" bind:value={newName} />
      <select class="freq-select" bind:value={newFreq}>
        {#each FREQ as f}
          <option value={f.value}>{f.label}</option>
        {/each}
      </select>
      <div class="cat-checkboxes">
        {#each availableCategories as cat}
          <label class="cat-check">
            <input type="checkbox" checked={newCategoryIds.includes(cat.id)} on:change={() => toggleCategory(cat.id)} />
            {cat.name}
          </label>
        {/each}
      </div>
      <div class="target-row">
        <input type="number" min="0" max="999" step="1" class="hm-input" bind:value={newHours} />
        <span class="unit">h</span>
        <input type="number" min="0" max="59" step="5" class="hm-input" bind:value={newMinutes} />
        <span class="unit">m</span>
        <button class="save-btn" on:click={addGoal} disabled={!newName.trim() || !newCategoryIds.length}>Save</button>
      </div>
      {#if addError}<div class="add-error">{addError}</div>{/if}
    </div>
  {/if}

  {#if isLoading}
    <div class="empty">Loading…</div>
  {:else if activeGoals.length === 0 && !showAddForm}
    <div class="empty">No active goals. Click + to add one.</div>
  {:else}
    {#each SECTIONS as section}
      {@const sectionGoals = section.goals()}
      {#if sectionGoals.length > 0}
        <div class="section">
          <div class="section-label">
            {section.label}
            <span class="sublabel">({section.sublabel()})</span>
          </div>
          <ul class="goal-list">
            {#each sectionGoals as goal}
              {@const actualMs = usageForGoal(goal)}
              {@const target   = goal.targetMs ?? 0}
              {@const pct      = progressPct(actualMs, target)}
              <li class="goal-item">
                <div class="goal-header">
                  <span class="goal-name">{goal.name} <span class="freq-badge">{freqLabel(goal.frequency)}</span></span>
                  <span class="goal-values">
                    {formatTimeElapsed(actualMs)} / {formatTimeElapsed(target)}
                    <span class="pct {progressClass(pct)}">{pct}%</span>
                  </span>
                  <button class="remove-btn" on:click={() => removeGoal(goal)} title="Remove">✕</button>
                </div>
                <div class="progress-track">
                  <div class="progress-fill {progressClass(pct)}" style="width: {pct}%"></div>
                </div>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    {/each}
  {/if}
</div>

<style>
  .goals-panel {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    box-shadow: var(--card-shadow);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }

  h3 { margin: 0; font-size: 1rem; font-weight: 600; color: var(--text-color); }

  .add-btn, .remove-btn {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--secondary-color);
    font-size: 1rem;
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
  }
  .add-btn:hover, .remove-btn:hover { background: var(--table-row-hover); }

  .add-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
    padding: 0.75rem;
    background: var(--table-row-hover);
    border-radius: 6px;
  }

  .name-input, .freq-select {
    padding: 0.3rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
  }
  .name-input { width: 100%; box-sizing: border-box; }

  .cat-checkboxes {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
  }
  .cat-check {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.8rem;
    color: var(--text-color);
    cursor: pointer;
  }

  .target-row {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .hm-input {
    width: 3.5rem;
    padding: 0.25rem 0.4rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
    text-align: right;
  }
  .unit { color: var(--secondary-color); font-size: 0.85rem; }

  .save-btn {
    padding: 0.3rem 0.75rem;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
    margin-left: auto;
  }
  .save-btn:disabled { opacity: 0.4; cursor: default; }

  .add-error { font-size: 0.75rem; color: #ef4444; }

  .section { margin-bottom: 0.75rem; }
  .section:last-child { margin-bottom: 0; }

  .section-label {
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--secondary-color);
    margin-bottom: 0.4rem;
  }
  .sublabel { font-weight: 400; text-transform: none; letter-spacing: 0; }

  .goal-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 0.5rem; }
  .goal-item { display: flex; flex-direction: column; gap: 0.2rem; }
  .goal-header { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; }
  .goal-name { flex: 1; font-weight: 500; color: var(--text-color); display: flex; align-items: center; gap: 0.35rem; }

  .freq-badge {
    font-size: 0.68rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--secondary-color);
    background: var(--table-border-color);
    border-radius: 99px;
    padding: 0.1rem 0.4rem;
  }

  .goal-values {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--secondary-color);
    white-space: nowrap;
    font-size: 0.8rem;
  }

  .pct { font-weight: 600; }
  .pct.met   { color: #22c55e; }
  .pct.close { color: #f59e0b; }
  .pct.low   { color: #ef4444; }

  .progress-track {
    height: 4px;
    background-color: var(--table-border-color);
    border-radius: 99px;
    overflow: hidden;
  }
  .progress-fill { height: 100%; border-radius: 99px; transition: width 0.3s ease; }
  .progress-fill.met   { background: #22c55e; }
  .progress-fill.close { background: #f59e0b; }
  .progress-fill.low   { background: #ef4444; }

  .empty { font-size: 0.85rem; color: var(--secondary-color); font-style: italic; text-align: center; padding: 0.5rem 0; }
</style>
