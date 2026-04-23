<script>
  import { GetGoals, GetCategories, DeleteGoal, GetCategoryUsageTotals } from '../../wailsjs/go/main/App';
  import { formatTimeElapsed } from '../utils/formatters';
  import { shiftDateStr } from '../utils/dateUtils';
  import { onMount } from 'svelte';
  import { refreshData } from '../stores/timekeeper';
  import CreateGoalModal from './goals/CreateGoalModal.svelte';

  export let selectedDate = '';

  const FREQ = [
    { value: 1, label: 'Daily' },
    { value: 2, label: 'Weekly' },
    { value: 3, label: 'Monthly' },
  ];
  function freqLabel(f) { return FREQ.find(x => x.value === f)?.label ?? 'Daily'; }

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
  let showAddModal = false;
  let isLoading = true;

  $: weekStart  = selectedDate ? calendarWeekStart(selectedDate)  : '';
  $: monthStart = selectedDate ? calendarMonthStart(selectedDate) : '';

  $: activeGoals = goals.filter(g => g.isActive);
  $: availableCategories = categories.filter(c => c.id > 0);

  $: dailyGoals   = activeGoals.filter(g => g.frequency === 1);
  $: weeklyGoals  = activeGoals.filter(g => g.frequency === 2);
  $: monthlyGoals = activeGoals.filter(g => g.frequency === 3);

  // Rebuild reactively so Svelte tracks dailyGoals/weeklyGoals/monthlyGoals
  // and the date-range vars. Closures hide deps from Svelte's compile-time
  // tracking, so {@const} on section.goals() would read stale values.
  $: sections = [
    { label: 'Daily',   goals: dailyGoals,   sublabel: selectedDate },
    { label: 'Weekly',  goals: weeklyGoals,  sublabel: `${weekStart} – ${selectedDate}` },
    { label: 'Monthly', goals: monthlyGoals, sublabel: `${monthStart} – ${selectedDate}` },
  ];

  // Precompute usage per goal reactively so the template picks up changes to
  // dailyUsage/weeklyUsage/monthlyUsage after a date change or refresh.
  $: usageByGoalId = buildUsageByGoalId(activeGoals, dailyUsage, weeklyUsage, monthlyUsage);

  function buildUsageByGoalId(goalsList, daily, weekly, monthly) {
    const map = new Map();
    for (const g of goalsList) {
      const usage = g.frequency === 2 ? weekly : g.frequency === 3 ? monthly : daily;
      const total = (g.categoryIds ?? []).reduce((sum, cid) => {
        const c = (usage ?? []).find(x => x.id === cid);
        return sum + (c ? (c.timeElapsed ?? 0) : 0);
      }, 0);
      map.set(g.id, total);
    }
    return map;
  }

  onMount(async () => {
    await loadGoals();
    await loadPeriodData();
  });

  // Reload when date changes (navigation) or when global refresh fires.
  $: if (selectedDate) loadPeriodData();
  $: if ($refreshData) loadPeriodData();

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
    try {
      const [d, w, m] = await Promise.all([
        GetCategoryUsageTotals(date, date),
        GetCategoryUsageTotals(calendarWeekStart(date), date),
        GetCategoryUsageTotals(calendarMonthStart(date), date),
      ]);
      if (date !== selectedDate) return; // stale, a newer date is in flight
      dailyUsage   = d ?? [];
      weeklyUsage  = w ?? [];
      monthlyUsage = m ?? [];
    } catch (err) {
      console.error('[GoalsPanel] Error loading period usage:', err);
    }
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

  async function onGoalAdded() {
    showAddModal = false;
    await loadGoals();
    await loadPeriodData();
  }
</script>

<CreateGoalModal
  show={showAddModal}
  {categories}
  on:goalAdded={onGoalAdded}
  on:close={() => { showAddModal = false; }}
/>

<div class="goals-panel">
  <div class="panel-header">
    <h3>Goals</h3>
    <button class="add-btn" on:click={() => { showAddModal = true; }} title="Add goal">+</button>
  </div>

  {#if isLoading}
    <div class="empty">Loading…</div>
  {:else if activeGoals.length === 0}
    <div class="empty">No active goals. Click + to add one.</div>
  {:else}
    {#each sections as section}
      {#if section.goals.length > 0}
        <div class="section">
          <div class="section-label">
            {section.label}
            <span class="sublabel">({section.sublabel})</span>
          </div>
          <ul class="goal-list">
            {#each section.goals as goal}
              {@const actualMs = usageByGoalId.get(goal.id) ?? 0}
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
