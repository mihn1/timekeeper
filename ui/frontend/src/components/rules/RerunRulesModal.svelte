<script>
  import { createEventDispatcher } from 'svelte';
  import Modal from '../common/Modal.svelte';
  import { StartRerunRules } from '../../../wailsjs/go/main/App';
  import { rerunStatus } from '../../stores/rerun';
  import { timezone } from '../../stores/preferences';
  import { todayInTz, shiftDateStr } from '../../utils/dateUtils';

  export let show = false;

  const dispatch = createEventDispatcher();

  $: maxDays = $rerunStatus?.maxRangeDays || 7;
  $: isRunning = $rerunStatus?.state === 'running';

  let startDate = '';
  let endDate = '';
  let submitError = '';
  let submitting = false;

  $: if (show && !startDate) {
    const today = todayInTz($timezone);
    endDate = today;
    startDate = shiftDateStr(today, -(maxDays - 1));
  }

  $: rangeDays = calcRangeDays(startDate, endDate);
  $: rangeInvalid = !startDate || !endDate || rangeDays < 1 || rangeDays > maxDays;

  function calcRangeDays(s, e) {
    if (!s || !e) return 0;
    const a = new Date(s + 'T00:00:00');
    const b = new Date(e + 'T00:00:00');
    return Math.floor((b - a) / 86400000) + 1;
  }

  function close() {
    submitError = '';
    dispatch('close');
  }

  async function submit() {
    if (rangeInvalid || submitting || isRunning) return;
    submitting = true;
    submitError = '';
    try {
      await StartRerunRules(startDate, endDate);
      dispatch('started');
      close();
    } catch (err) {
      submitError = err?.message || String(err) || 'Failed to start rerun';
    } finally {
      submitting = false;
    }
  }
</script>

<Modal show={show} title="Rerun Rules" on:close={close}>
  <div class="p-4 space-y-4">
    <p class="text-sm rerun-info">
      Re-apply current category rules to every event in the selected date range and rebuild
      aggregations for those days. Maximum range: <strong>{maxDays}</strong> days.
    </p>

    <div class="grid grid-cols-2 gap-3">
      <label class="flex flex-col text-sm">
        <span class="mb-1">Start date</span>
        <input
          type="date"
          class="p-2 border border-gray-300 rounded date-input"
          bind:value={startDate}
          disabled={submitting || isRunning}
        />
      </label>

      <label class="flex flex-col text-sm">
        <span class="mb-1">End date</span>
        <input
          type="date"
          class="p-2 border border-gray-300 rounded date-input"
          bind:value={endDate}
          disabled={submitting || isRunning}
        />
      </label>
    </div>

    {#if startDate && endDate}
      <p class="text-sm rerun-info">
        Selected: <strong>{rangeDays}</strong> {rangeDays === 1 ? 'day' : 'days'}
        {#if rangeDays > maxDays}
          <span class="text-red-600 ml-2">exceeds max of {maxDays}</span>
        {:else if rangeDays < 1}
          <span class="text-red-600 ml-2">end date must be on or after start date</span>
        {/if}
      </p>
    {/if}

    {#if isRunning}
      <p class="text-sm text-amber-600">A rerun job is already in progress.</p>
    {/if}

    {#if submitError}
      <p class="text-sm text-red-600">{submitError}</p>
    {/if}

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
        disabled={rangeInvalid || submitting || isRunning}
      >{submitting ? 'Starting…' : 'Run'}</button>
    </div>
  </div>
</Modal>

<style>
  .rerun-info {
    color: var(--text-color);
    opacity: 0.85;
  }
  .date-input {
    background-color: var(--input-bg-color, var(--card-bg-color));
    color: var(--text-color);
  }
</style>
