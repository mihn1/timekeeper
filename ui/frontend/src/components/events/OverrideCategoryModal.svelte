<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import Modal from '../common/Modal.svelte';
  import { OverrideEventCategory } from '../../../wailsjs/go/main/App';
  import { dtos } from '../../../wailsjs/go/models';

  export let event: dtos.EventLogItem;
  export let categories: dtos.CategoryListItem[] = [];

  const dispatch = createEventDispatcher();

  let selectedCategoryId: number = event.categoryId;
  let submitting = false;
  let submitError = '';

  function close() {
    submitError = '';
    dispatch('close');
  }

  async function submit() {
    if (submitting) return;
    if (selectedCategoryId === event.categoryId) {
      close();
      return;
    }
    submitting = true;
    submitError = '';
    try {
      await OverrideEventCategory(event.id, selectedCategoryId);
      dispatch('saved', { eventId: event.id, categoryId: selectedCategoryId });
      close();
    } catch (err) {
      submitError = err?.message || String(err) || 'Failed to override category';
    } finally {
      submitting = false;
    }
  }
</script>

<Modal show={true} title="Override Event Category" on:close={close}>
  <div class="p-4 space-y-4">
      <div class="text-sm event-info">
        <div><strong>App:</strong> {event.appName}</div>
        <div><strong>Time:</strong> {event.startTime} – {event.endTime}</div>
        {#if event.urlOrTitle}
          <div class="truncate"><strong>URL/Title:</strong> {event.urlOrTitle}</div>
        {/if}
      </div>

      <label class="flex flex-col text-sm">
        <span class="mb-1">Category</span>
        <select
          class="p-2 border border-gray-300 rounded category-select"
          bind:value={selectedCategoryId}
          disabled={submitting}
        >
          <option value={0}>Ignored</option>
          {#each categories as cat}
            <option value={cat.id}>{cat.name}</option>
          {/each}
        </select>
      </label>

      {#if submitError}
        <p class="text-sm text-red-600">{submitError}</p>
      {/if}

      <p class="text-xs hint">
        Category aggregation for this date will be adjusted. App-level totals are unchanged.
      </p>

      <div class="flex justify-end gap-2 pt-2">
        <button
          type="button"
          class="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400 cursor-pointer"
          on:click={close}
          disabled={submitting}
        >Cancel</button>
        <button
          type="button"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          on:click={submit}
          disabled={submitting}
        >{submitting ? 'Saving…' : 'Save'}</button>
      </div>
  </div>
</Modal>

<style>
  .event-info {
    color: var(--text-color);
    opacity: 0.9;
    line-height: 1.5;
  }
  .category-select {
    background-color: var(--input-bg-color, var(--card-bg-color));
    color: var(--text-color);
  }
  .hint {
    color: var(--secondary-color);
    font-style: italic;
  }
</style>
