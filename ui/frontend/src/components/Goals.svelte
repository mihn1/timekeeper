<script>
  import { onMount } from 'svelte';
  import { GetGoals, GetCategories, UpdateGoal, DeleteGoal } from '../../wailsjs/go/main/App';
  import { refreshData } from '../stores/timekeeper';
  import CreateGoalModal from './goals/CreateGoalModal.svelte';

  const FREQ = [
    { value: 1, label: 'Daily' },
    { value: 2, label: 'Weekly' },
    { value: 3, label: 'Monthly' },
  ];

  function freqLabel(f) { return FREQ.find(x => x.value === f)?.label ?? 'Daily'; }
  function msToHm(ms) { return { h: Math.floor(ms / 3600000), m: Math.floor((ms % 3600000) / 60000) }; }
  function hmToMs(h, m) { return (Number(h) * 3600 + Number(m) * 60) * 1000; }
  function formatTarget(ms) {
    const { h, m } = msToHm(ms);
    if (h === 0) return `${m}m`;
    if (m === 0) return `${h}h`;
    return `${h}h ${m}m`;
  }

  let goals = [];
  let categories = [];
  let isLoading = true;
  let loadError = null;
  let showAddModal = false;

  // Edit state — keyed by goal.id
  let editingId = null;
  let editName = '';
  let editFreq = 1;
  let editCategoryIds = [];
  let editHours = 0;
  let editMinutes = 0;
  let editError = null;

  $: availableCategories = categories.filter(c => c.id > 0);
  $: goalsByFreq = FREQ.map(f => ({
    freq: f.value,
    label: f.label,
    goals: goals.filter(g => g.frequency === f.value),
  })).filter(g => g.goals.length > 0);

  $: if ($refreshData) loadGoals();

  onMount(loadGoals);

  async function loadGoals() {
    isLoading = true;
    loadError = null;
    try {
      [goals, categories] = await Promise.all([GetGoals(), GetCategories()]);
    } catch (err) {
      console.error('Error loading goals:', err);
      loadError = 'Failed to load goals.';
    } finally {
      isLoading = false;
    }
  }

  function toggleEditCategory(id) {
    const numId = Number(id);
    if (editCategoryIds.includes(numId)) {
      editCategoryIds = editCategoryIds.filter(x => x !== numId);
    } else {
      editCategoryIds = [...editCategoryIds, numId];
    }
  }

  function startEdit(goal) {
    const { h, m } = msToHm(goal.targetMs);
    editingId = goal.id;
    editName = goal.name;
    editFreq = goal.frequency;
    editCategoryIds = [...(goal.categoryIds ?? [])];
    editHours = h;
    editMinutes = m;
    editError = null;
  }

  function cancelEdit() {
    editingId = null;
    editError = null;
  }

  async function saveEdit(goal) {
    if (!editName.trim()) { editError = 'Name is required.'; return; }
    if (!editCategoryIds.length) { editError = 'Select at least one category.'; return; }
    const ms = hmToMs(editHours, editMinutes);
    if (ms <= 0) { editError = 'Target must be greater than 0.'; return; }
    editError = null;
    try {
      await UpdateGoal(goal.id, editName.trim(), editCategoryIds.map(Number), editFreq, ms, goal.isActive);
      editingId = null;
      await loadGoals();
    } catch (err) {
      console.error('Error updating goal:', err);
      editError = 'Failed to update.';
    }
  }

  async function toggleActive(goal) {
    try {
      await UpdateGoal(goal.id, goal.name, goal.categoryIds ?? [], goal.frequency, goal.targetMs, !goal.isActive);
      await loadGoals();
    } catch (err) {
      console.error('Error toggling goal:', err);
    }
  }

  async function removeGoal(goal) {
    try {
      await DeleteGoal(goal.id);
      await loadGoals();
    } catch (err) {
      console.error('Error deleting goal:', err);
    }
  }
</script>

<CreateGoalModal
  show={showAddModal}
  {categories}
  on:goalAdded={async () => { showAddModal = false; await loadGoals(); }}
  on:close={() => { showAddModal = false; }}
/>

<div class="p-6 max-w-3xl mx-auto">
  <div class="flex items-center justify-between mb-6">
    <h1 class="text-2xl font-bold page-title">Goals</h1>
    <button class="add-goal-btn" on:click={() => { showAddModal = true; }}>+ Add Goal</button>
  </div>

  <!-- Existing goals -->
  <div class="section-card">
    <div class="mb-4">
      <h2 class="text-lg font-semibold">Time Goals</h2>
      <p class="text-sm section-desc mt-1">Daily, weekly, and monthly targets across one or more categories. Progress shown on the Dashboard.</p>
    </div>

    {#if isLoading}
      <div class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-6 w-6 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    {:else if loadError}
      <div class="error-banner p-3 rounded text-sm">
        {loadError}
        <button class="underline ml-2 cursor-pointer" on:click={loadGoals}>Retry</button>
      </div>
    {:else if goals.length === 0}
      <p class="empty-msg text-sm italic py-4 text-center">No goals set yet.</p>
    {:else}
      {#each goalsByFreq as group}
        <div class="goal-group mb-5">
          <div class="group-label text-xs font-semibold uppercase tracking-wider mb-2">{group.label}</div>
          <table class="goals-table w-full text-sm">
            <thead>
              <tr class="table-head-row">
                <th class="text-left py-2 px-3 font-medium">Name</th>
                <th class="text-left py-2 px-3 font-medium">Categories</th>
                <th class="text-left py-2 px-3 font-medium">Target</th>
                <th class="py-2 px-3 font-medium text-center" style="width:4rem">Active</th>
                <th class="py-2 px-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each group.goals as goal (goal.id)}
                <tr class="goal-row">
                  {#if editingId === goal.id}
                    <td class="py-2 px-3" colspan="5">
                      <div class="edit-form">
                        <div class="edit-row">
                          <input type="text" class="edit-input flex-1" placeholder="Goal name" bind:value={editName} />
                          <select class="edit-select" bind:value={editFreq}>
                            {#each FREQ as f}
                              <option value={f.value}>{f.label}</option>
                            {/each}
                          </select>
                        </div>
                        <div class="edit-cats">
                          {#each availableCategories as cat}
                            <label class="cat-check">
                              <input type="checkbox" checked={editCategoryIds.includes(cat.id)} on:change={() => toggleEditCategory(cat.id)} />
                              {cat.name}
                            </label>
                          {/each}
                        </div>
                        <div class="edit-row">
                          <input type="number" min="0" max="999" step="1" class="hm-input" bind:value={editHours} />
                          <span class="unit-label">h</span>
                          <input type="number" min="0" max="59" step="5" class="hm-input" bind:value={editMinutes} />
                          <span class="unit-label">m</span>
                          {#if editError}<span class="text-red-500 text-xs ml-2">{editError}</span>{/if}
                          <div class="flex gap-2 ml-auto">
                            <button class="action-btn save-btn" on:click={() => saveEdit(goal)}>Save</button>
                            <button class="action-btn cancel-btn" on:click={cancelEdit}>Cancel</button>
                          </div>
                        </div>
                      </div>
                    </td>
                  {:else}
                    <td class="py-2 px-3 font-medium">{goal.name}</td>
                    <td class="py-2 px-3 secondary-text">{(goal.categoryNames ?? []).join(', ')}</td>
                    <td class="py-2 px-3">
                      <span class="target-value">{formatTarget(goal.targetMs)}</span>
                    </td>
                    <td class="py-2 px-3 text-center">
                      <button
                        class="active-toggle"
                        class:active={goal.isActive}
                        on:click={() => toggleActive(goal)}
                        title={goal.isActive ? 'Active — click to deactivate' : 'Inactive — click to activate'}
                      >
                        {goal.isActive ? '●' : '○'}
                      </button>
                    </td>
                    <td class="py-2 px-3 text-right">
                      <div class="flex justify-end gap-2">
                        <button class="icon-btn edit-icon" on:click={() => startEdit(goal)} title="Edit">
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                          </svg>
                        </button>
                        <button class="icon-btn delete-icon" on:click={() => removeGoal(goal)} title="Delete">
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                          </svg>
                        </button>
                      </div>
                    </td>
                  {/if}
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .page-title { color: var(--text-color); }

  .section-card {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1.25rem 1.5rem;
    box-shadow: var(--card-shadow);
  }

  .section-desc { color: var(--secondary-color); }

  .error-banner {
    background-color: #fef2f2;
    border: 1px solid #fecaca;
    color: #b91c1c;
  }

  .empty-msg { color: var(--secondary-color); }
  .group-label { color: var(--secondary-color); }
  .secondary-text { color: var(--secondary-color); font-size: 0.8rem; }

  .goals-table { border-collapse: collapse; }

  .table-head-row {
    border-bottom: 1px solid var(--table-border-color);
    color: var(--secondary-color);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .goal-row {
    border-bottom: 1px solid var(--table-border-color);
    color: var(--text-color);
  }
  .goal-row:last-child { border-bottom: none; }
  .goal-row:hover { background-color: var(--table-row-hover); }

  .target-value { font-variant-numeric: tabular-nums; font-weight: 500; }

  .active-toggle {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    color: var(--secondary-color);
    padding: 0;
    line-height: 1;
  }
  .active-toggle.active { color: #22c55e; }

  /* Edit form */
  .edit-form { display: flex; flex-direction: column; gap: 0.5rem; }
  .edit-row { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
  .edit-input {
    padding: 0.25rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
    min-width: 120px;
  }
  .edit-select {
    padding: 0.25rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
  }
  .edit-cats { display: flex; flex-wrap: wrap; gap: 0.4rem; }

  .cat-check {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.8rem;
    color: var(--text-color);
    cursor: pointer;
  }
  .hm-input {
    width: 3.5rem;
    padding: 0.2rem 0.4rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
    text-align: right;
  }

  .unit-label { color: var(--secondary-color); font-size: 0.85rem; }

  .icon-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.2rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
  }
  .edit-icon { color: var(--primary-color); }
  .edit-icon:hover { color: var(--info-color); }
  .delete-icon { color: var(--danger-color); }
  .delete-icon:hover { opacity: 0.75; }

  .action-btn {
    padding: 0.2rem 0.6rem;
    border: none;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }
  .save-btn { background-color: var(--primary-color); color: white; }
  .save-btn:hover { opacity: 0.85; }
  .cancel-btn { background-color: var(--button-bg-color); color: var(--text-color); }
  .cancel-btn:hover { background-color: var(--button-hover-bg-color); }

  .add-goal-btn {
    padding: 0.4rem 1rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
  }
  .add-goal-btn:hover:not([disabled]) { opacity: 0.85; }
  .add-goal-btn[disabled] { opacity: 0.4; cursor: default; }
</style>
