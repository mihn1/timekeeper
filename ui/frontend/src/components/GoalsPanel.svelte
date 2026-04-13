<script>
  import { formatTimeElapsed } from '../utils/formatters';
  import { GetGoals, SetGoal, DeleteGoal, GetCategories } from '../../wailsjs/go/main/App';
  import { onMount } from 'svelte';

  export let categoryUsageData = [];

  let goals = [];
  let categories = [];
  let showAddForm = false;
  let newCategoryId = '';
  let newTargetHours = 1;
  let isLoading = true;

  onMount(async () => {
    await loadGoals();
  });

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

  async function addGoal() {
    if (!newCategoryId) return;
    try {
      await SetGoal(parseInt(newCategoryId), Math.round(newTargetHours * 3600000));
      newCategoryId = '';
      newTargetHours = 1;
      showAddForm = false;
      await loadGoals();
    } catch (err) {
      console.error('Error setting goal:', err);
    }
  }

  async function removeGoal(categoryId) {
    try {
      await DeleteGoal(categoryId);
      await loadGoals();
    } catch (err) {
      console.error('Error deleting goal:', err);
    }
  }

  function getActualMs(categoryId) {
    const cat = categoryUsageData.find(c => c.id === categoryId);
    return cat ? cat.timeElapsed : 0;
  }

  function progressPct(actualMs, targetMs) {
    if (targetMs <= 0) return 0;
    return Math.min(100, Math.round((actualMs / targetMs) * 100));
  }

  function progressClass(pct) {
    if (pct >= 100) return 'met';
    if (pct >= 50)  return 'close';
    return 'low';
  }

  $: availableCategories = categories.filter(
    c => c.id > 0 && !goals.some(g => g.categoryId === c.id)
  );
</script>

<div class="goals-panel">
  <div class="panel-header">
    <h3>Daily Goals</h3>
    <button class="add-btn" on:click={() => showAddForm = !showAddForm} title="Add goal">
      {showAddForm ? '✕' : '+'}
    </button>
  </div>

  {#if showAddForm}
    <div class="add-form">
      <select bind:value={newCategoryId}>
        <option value="">Category…</option>
        {#each availableCategories as cat}
          <option value={cat.id}>{cat.name}</option>
        {/each}
      </select>
      <input type="number" min="0.5" max="24" step="0.5" bind:value={newTargetHours} />
      <span class="unit">h</span>
      <button class="save-btn" on:click={addGoal} disabled={!newCategoryId}>Save</button>
    </div>
  {/if}

  {#if isLoading}
    <div class="empty">Loading…</div>
  {:else if goals.length === 0 && !showAddForm}
    <div class="empty">No goals set. Click + to add one.</div>
  {:else}
    <ul class="goal-list">
      {#each goals as goal}
        {@const actualMs = getActualMs(goal.categoryId)}
        {@const pct = progressPct(actualMs, goal.dailyTargetMs)}
        <li class="goal-item">
          <div class="goal-header">
            <span class="goal-name">{goal.categoryName}</span>
            <span class="goal-values">
              {formatTimeElapsed(actualMs)} / {formatTimeElapsed(goal.dailyTargetMs)}
              <span class="pct {progressClass(pct)}">{pct}%</span>
            </span>
            <button class="remove-btn" on:click={() => removeGoal(goal.categoryId)} title="Remove goal">✕</button>
          </div>
          <div class="progress-track">
            <div class="progress-fill {progressClass(pct)}" style="width: {pct}%"></div>
          </div>
        </li>
      {/each}
    </ul>
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

  h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-color);
  }

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
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
    flex-wrap: wrap;
  }

  .add-form select, .add-form input {
    padding: 0.3rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
  }

  .add-form input { width: 4rem; }
  .unit { color: var(--secondary-color); font-size: 0.85rem; }

  .save-btn {
    padding: 0.3rem 0.75rem;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
  }
  .save-btn:disabled { opacity: 0.4; cursor: default; }

  .goal-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }

  .goal-item { display: flex; flex-direction: column; gap: 0.25rem; }

  .goal-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
  }

  .goal-name { flex: 1; font-weight: 500; color: var(--text-color); }

  .goal-values {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--secondary-color);
    white-space: nowrap;
  }

  .pct { font-weight: 600; }
  .pct.met   { color: #22c55e; }
  .pct.close { color: #f59e0b; }
  .pct.low   { color: #ef4444; }

  .progress-track {
    height: 5px;
    background-color: var(--table-border-color);
    border-radius: 99px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    border-radius: 99px;
    transition: width 0.3s ease;
  }
  .progress-fill.met   { background: #22c55e; }
  .progress-fill.close { background: #f59e0b; }
  .progress-fill.low   { background: #ef4444; }

  .empty {
    font-size: 0.85rem;
    color: var(--secondary-color);
    font-style: italic;
    text-align: center;
    padding: 0.5rem 0;
  }
</style>
