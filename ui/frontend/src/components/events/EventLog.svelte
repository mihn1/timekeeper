<script lang="ts">
  import { onMount } from 'svelte';
  import { GetEventLog, GetCategories, DeleteEvent } from '../../../wailsjs/go/main/App';
  import { refreshData } from '../../stores/timekeeper';
  import { timezone } from '../../stores/preferences.js';
  import { todayInTz } from '../../utils/dateUtils.js';
  import DataTable from '../common/DataTable.svelte';
  import { dtos } from '../../../wailsjs/go/models';
  import type { Column } from '../../types/table';

  let selectedDate: string = todayInTz($timezone);
  let events: dtos.EventLogItem[] = [];
  let categories: dtos.CategoryListItem[] = [];
  let isLoading = true;
  let loadError: string | null = null;
  let searchTerm = '';
  let showIgnored = false;

  $: if ($refreshData || selectedDate) { loadEvents(); }

  $: filteredEvents = events
    .filter(e => showIgnored || e.categoryId !== 0)
    .filter(e => e.appName.toLowerCase().includes(searchTerm.toLowerCase()));

  onMount(() => {
    loadCategories();
    // loadEvents() is triggered by the reactive $: block above on mount
  });

  async function loadEvents() {
    isLoading = true;
    loadError = null;
    try {
      const result = await GetEventLog(selectedDate);
      events = result ?? [];
    } catch (err) {
      loadError = 'Failed to load events. Please try again.';
      events = [];
    } finally {
      isLoading = false;
    }
  }

  async function loadCategories() {
    try {
      const result = await GetCategories();
      if (result) categories = result;
    } catch (_) {}
  }

  function getCategoryName(id: number): string {
    return categories.find(c => c.id === id)?.name ?? (id === 0 ? 'Ignored' : `#${id}`);
  }

  function formatDuration(secs: number): string {
    if (secs <= 0) return '—';
    if (secs < 60) return `${secs}s`;
    const m = Math.floor(secs / 60);
    const s = secs % 60;
    if (m < 60) return `${m}m ${s}s`;
    return `${Math.floor(m / 60)}h ${m % 60}m`;
  }

  const columns: Column[] = [
    { key: 'startTime',    title: 'Start',    sortable: true },
    { key: 'endTime',      title: 'End',      sortable: false },
    { key: 'appName',      title: 'App',      sortable: true },
    { key: 'durationSecs', title: 'Duration', sortable: true,
      formatter: (v: number) => formatDuration(v) },
    { key: 'categoryId',   title: 'Category', sortable: true,
      formatter: (v: number) => getCategoryName(v) },
    { key: 'urlOrTitle',   title: 'URL / Title', sortable: false,
      formatter: (v: string) => v && v.length > 60 ? v.slice(0, 57) + '…' : (v ?? '') },
  ];

  async function deleteEvent(row: dtos.EventLogItem) {
    if (!confirm(`Delete this event?\n${row.appName} — ${row.startTime}`)) return;
    try {
      await DeleteEvent(row.id);
      events = events.filter(e => e.id !== row.id);
    } catch (err) {
      alert('Failed to delete event.');
    }
  }

  const rowActions = [
    { icon: 'trash', title: 'Delete event', handler: deleteEvent },
  ];
</script>

<div class="p-6 max-w-6xl mx-auto">
  <h1 class="text-2xl font-bold mb-6 title">Event Log</h1>

  <div class="flex gap-4 items-center mb-6 flex-wrap">
    <label class="font-medium" for="date-picker">Date:</label>
    <input
      id="date-picker"
      type="date"
      class="p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      bind:value={selectedDate}
    />
    <label class="flex items-center gap-2 cursor-pointer select-none show-ignored">
      <input type="checkbox" bind:checked={showIgnored} />
      <span>Show ignored events</span>
    </label>
    <button
      class="refresh-button"
      on:click={loadEvents}
      aria-label="Refresh events"
      title="Refresh events"
    >
      <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
      </svg>
    </button>
  </div>

  <div class="log-container rounded-lg shadow-md overflow-hidden">
    <div class="p-5 border-b table-header-wrapper">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-semibold table-header-title">
          Events ({filteredEvents.length})
        </h2>
        <input
          type="text"
          placeholder="Filter by app..."
          class="pl-3 p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          bind:value={searchTerm}
        />
      </div>
    </div>

    {#if isLoading}
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    {:else if loadError}
      <div class="flex items-center gap-2 p-4 m-4 text-red-700 bg-red-50 border border-red-200 rounded">
        <svg class="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <span>{loadError}</span>
        <button class="ml-auto text-sm underline cursor-pointer" on:click={loadEvents}>Retry</button>
      </div>
    {:else if filteredEvents.length === 0}
      <div class="flex flex-col items-center justify-center py-12 gap-2 empty-state">
        <p class="text-lg">No events recorded for this date.</p>
        <p class="text-sm">Event recording requires SQLite mode. In-memory mode does not persist events.</p>
      </div>
    {:else}
      <DataTable
        data={filteredEvents}
        columns={columns}
        rowActions={rowActions}
        emptyMessage="No events found"
        pageSize={25}
        initialSortKey="startTime"
        initialSortDirection="desc"
      />
    {/if}
  </div>
</div>

<style>
  .title { color: var(--text-color); }
  .log-container {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
  }
  .table-header-wrapper {
    background-color: var(--table-header-bg);
    border-color: var(--table-border-color);
  }
  .table-header-title { color: var(--text-color); }
  .empty-state { color: var(--secondary-color); font-style: italic; }
  .show-ignored {
    font-size: 0.875rem;
    color: var(--text-color);
  }
  .refresh-button {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.45rem;
    border-radius: 50%;
    border: none;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    cursor: pointer;
    margin-left: auto;
  }
  .refresh-button:hover { background-color: var(--button-hover-bg-color); }
  .icon { width: 1.1rem; height: 1.1rem; }
</style>
