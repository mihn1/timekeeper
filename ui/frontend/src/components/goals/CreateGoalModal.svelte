<script>
  import { createEventDispatcher } from 'svelte';
  import { AddGoal } from '../../../wailsjs/go/main/App';
  import Modal from '../common/Modal.svelte';

  export let show = false;
  export let categories = [];

  const dispatch = createEventDispatcher();

  const FREQ = [
    { value: 1, label: 'Daily' },
    { value: 2, label: 'Weekly' },
    { value: 3, label: 'Monthly' },
  ];

  function hmToMs(h, m) { return (Number(h) * 3600 + Number(m) * 60) * 1000; }

  let name = '';
  let freq = 1;
  let categoryIds = [];
  let hours = 1;
  let minutes = 0;
  let error = null;
  let saving = false;

  $: availableCategories = categories.filter(c => c.id > 0);

  function toggleCategory(id) {
    const numId = Number(id);
    if (categoryIds.includes(numId)) {
      categoryIds = categoryIds.filter(x => x !== numId);
    } else {
      categoryIds = [...categoryIds, numId];
    }
  }

  async function save() {
    if (!name.trim()) { error = 'Name is required.'; return; }
    if (!categoryIds.length) { error = 'Select at least one category.'; return; }
    const ms = hmToMs(hours, minutes);
    if (ms <= 0) { error = 'Target must be greater than 0.'; return; }
    error = null;
    saving = true;
    try {
      await AddGoal(name.trim(), categoryIds.map(Number), freq, ms);
      resetForm();
      dispatch('goalAdded');
    } catch (err) {
      console.error('Error adding goal:', err);
      error = 'Failed to save goal.';
    } finally {
      saving = false;
    }
  }

  function resetForm() {
    name = '';
    freq = 1;
    categoryIds = [];
    hours = 1;
    minutes = 0;
    error = null;
  }

  function close() {
    resetForm();
    dispatch('close');
  }
</script>

<Modal {show} title="Add Goal" on:close={close}>
  <div class="modal-body">
    <div class="field">
      <label class="field-label" for="goal-name">Name</label>
      <input
        id="goal-name"
        type="text"
        class="field-input"
        placeholder="e.g. Deep Work"
        bind:value={name}
      />
    </div>

    <div class="field">
      <label class="field-label" for="goal-freq">Frequency</label>
      <select id="goal-freq" class="field-select" bind:value={freq}>
        {#each FREQ as f}
          <option value={f.value}>{f.label}</option>
        {/each}
      </select>
    </div>

    <div class="field">
      <div class="field-label">Target</div>
      <div class="hm-row">
        <input type="number" min="0" max="999" step="1" class="hm-input" bind:value={hours} />
        <span class="unit">h</span>
        <input type="number" min="0" max="59" step="5" class="hm-input" bind:value={minutes} />
        <span class="unit">m</span>
      </div>
    </div>

    <div class="field">
      <div class="field-label">Categories</div>
      <div class="cat-grid">
        {#each availableCategories as cat}
          <label class="cat-check">
            <input
              type="checkbox"
              checked={categoryIds.includes(cat.id)}
              on:change={() => toggleCategory(cat.id)}
            />
            {cat.name}
          </label>
        {/each}
        {#if availableCategories.length === 0}
          <span class="empty-cats">No categories available.</span>
        {/if}
      </div>
    </div>

    {#if error}
      <p class="error-msg">{error}</p>
    {/if}

    <div class="modal-actions">
      <button class="btn-cancel" on:click={close}>Cancel</button>
      <button
        class="btn-save"
        on:click={save}
        disabled={saving || !name.trim() || !categoryIds.length}
      >
        {saving ? 'Saving…' : 'Add Goal'}
      </button>
    </div>
  </div>
</Modal>

<style>
  .modal-body {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1.25rem 1.5rem;
  }

  .field { display: flex; flex-direction: column; gap: 0.3rem; }

  .field-label {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--secondary-color);
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }

  .field-input, .field-select {
    padding: 0.4rem 0.6rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.875rem;
  }
  .field-input { width: 100%; box-sizing: border-box; }
  .field-select { min-width: 120px; }

  .hm-row { display: flex; align-items: center; gap: 0.4rem; }
  .hm-input {
    width: 4rem;
    padding: 0.35rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.875rem;
    text-align: right;
  }
  .unit { color: var(--secondary-color); font-size: 0.875rem; }

  .cat-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem 1rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--input-border-color);
    border-radius: 6px;
    background: var(--input-bg-color);
    min-height: 2.5rem;
  }

  .cat-check {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.875rem;
    color: var(--text-color);
    cursor: pointer;
  }

  .empty-cats {
    font-size: 0.8rem;
    color: var(--secondary-color);
    font-style: italic;
  }

  .error-msg {
    font-size: 0.8rem;
    color: #ef4444;
    margin: 0;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    padding-top: 0.25rem;
    border-top: 1px solid var(--card-border-color);
    margin-top: 0.25rem;
  }

  .btn-cancel {
    padding: 0.4rem 1rem;
    background-color: var(--button-bg-color);
    color: var(--text-color);
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
  }
  .btn-cancel:hover { background-color: var(--button-hover-bg-color); }

  .btn-save {
    padding: 0.4rem 1rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
  }
  .btn-save:hover:not([disabled]) { opacity: 0.85; }
  .btn-save[disabled] { opacity: 0.4; cursor: default; }
</style>
