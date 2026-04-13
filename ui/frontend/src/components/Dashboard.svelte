<script>
  import {
    GetAppUsageData,
    GetCategoryUsageData,
    GetEventLog,
    GetUncategorizedApps,
    GetCategoryUsageRange,
    GetActivityCalendar,
    GetCategories,
  } from '../../wailsjs/go/main/App';
  import { formatTimeElapsed } from '../utils/formatters';
  import { refreshData } from '../stores/timekeeper';
  import AppUsageChart from './AppUsageChart.svelte';
  import CategoryChart from './CategoryChart.svelte';
  import DaySummaryBar from './DaySummaryBar.svelte';
  import ComparisonStrip from './ComparisonStrip.svelte';
  import GoalsPanel from './GoalsPanel.svelte';
  import ActivityTimeline from './ActivityTimeline.svelte';
  import TopDomainsPanel from './TopDomainsPanel.svelte';
  import UncategorizedAppsPanel from './UncategorizedAppsPanel.svelte';
  import TrendChart from './TrendChart.svelte';
  import HeatmapCalendar from './HeatmapCalendar.svelte';
  import CreateRuleModal from './rules/CreateRuleModal.svelte';

  function localDateStr(d = new Date()) {
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${y}-${m}-${day}`;
  }

  const today = localDateStr();

  // ── State ─────────────────────────────────────────────────────────────────
  let selectedDate = today;
  let viewMode = 'day'; // 'day' | '7d' | '14d' | '30d'
  let minDurationMs = 0;
  let selectedCategoryFilter = null;
  let showDateInput = false;
  let showCreateRuleModal = false;
  let prefillAppName = '';

  // ── Data ──────────────────────────────────────────────────────────────────
  let appUsageData = [];
  let categoryUsageData = [];
  let eventLog = [];
  let uncategorizedApps = [];
  let trendData = [];
  let calendarData = [];
  let categories = [];
  let yesterdayMs = 0;
  let weekAvgMs = 0;

  let isLoading = true;
  let loadError = null;

  // ── Reactive ──────────────────────────────────────────────────────────────
  $: isToday = selectedDate === today;
  $: selectedYear = parseInt(selectedDate.split('-')[0]);

  $: filteredAppUsageData = appUsageData
    .filter(app => app.timeElapsed >= minDurationMs)
    .filter(app => selectedCategoryFilter === null || app.categoryId === selectedCategoryFilter);

  $: if ($refreshData || selectedDate) {
    loadData();
  }

  $: if (viewMode !== 'day') {
    loadRangeData();
  }

  // ── Data loading ──────────────────────────────────────────────────────────
  async function loadData() {
    isLoading = true;
    loadError = null;
    try {
      const [apps, cats, events, uncats, allCats] = await Promise.all([
        GetAppUsageData(selectedDate),
        GetCategoryUsageData(selectedDate),
        GetEventLog(selectedDate),
        GetUncategorizedApps(selectedDate),
        GetCategories(),
      ]);

      appUsageData = (apps ?? []).sort((a, b) => b.timeElapsed - a.timeElapsed);
      categoryUsageData = (cats ?? []).sort((a, b) => b.timeElapsed - a.timeElapsed);
      eventLog = events ?? [];
      uncategorizedApps = uncats ?? [];
      categories = allCats ?? [];

      // Load comparison data (yesterday + 7-day window) in background.
      loadComparisonData();

      // Load calendar for the current year.
      GetActivityCalendar(selectedYear).then(r => { calendarData = r ?? []; }).catch(() => {});
    } catch (err) {
      console.error('Error loading data:', err);
      loadError = 'Failed to load data. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  async function loadComparisonData() {
    try {
      const yesterday = shiftDateStr(selectedDate, -1);
      const weekStart = shiftDateStr(selectedDate, -7);

      const [yData, weekData] = await Promise.all([
        GetCategoryUsageData(yesterday),
        GetCategoryUsageRange(weekStart, shiftDateStr(selectedDate, -1)),
      ]);

      yesterdayMs = (yData ?? []).filter(c => c.id !== 0).reduce((s, c) => s + c.timeElapsed, 0);

      // Weekly average: sum all days / number of days with data.
      const dayTotals = new Map();
      for (const row of (weekData ?? [])) {
        if (row.categoryId !== 0) {
          dayTotals.set(row.date, (dayTotals.get(row.date) ?? 0) + row.timeElapsed);
        }
      }
      const days = dayTotals.size;
      weekAvgMs = days > 0
        ? Array.from(dayTotals.values()).reduce((s, v) => s + v, 0) / days
        : 0;
    } catch {}
  }

  async function loadRangeData() {
    const days = viewMode === '7d' ? 7 : viewMode === '14d' ? 14 : 30;
    const startDate = shiftDateStr(selectedDate, -(days - 1));
    try {
      trendData = await GetCategoryUsageRange(startDate, selectedDate) ?? [];
    } catch (err) {
      console.error('Error loading range data:', err);
      trendData = [];
    }
  }

  // ── Date navigation ───────────────────────────────────────────────────────
  function shiftDateStr(dateStr, days) {
    const d = new Date(dateStr + 'T00:00:00');
    d.setDate(d.getDate() + days);
    return localDateStr(d);
  }

  function shiftDate(days) { selectedDate = shiftDateStr(selectedDate, days); }
  function goToday() { selectedDate = today; }

  function formatDisplayDate(dateStr) {
    const d = new Date(dateStr + 'T00:00:00');
    return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' });
  }

  // ── UI handlers ───────────────────────────────────────────────────────────
  function handleCategoryFilter(e) {
    selectedCategoryFilter = e.detail.categoryId;
  }

  function handleCreateRule(e) {
    prefillAppName = e.detail.appName;
    showCreateRuleModal = true;
  }

  function handleDaySelected(e) {
    selectedDate = e.detail.date;
    viewMode = 'day';
  }
</script>

<div class="dashboard">
  <!-- ── Controls bar ─────────────────────────────────────────────────────── -->
  <div class="controls">
    <div class="date-nav">
      <button class="nav-btn" on:click={() => shiftDate(-1)} title="Previous day">&#8249;</button>

      {#if showDateInput}
        <input
          type="date"
          class="date-input-popup"
          bind:value={selectedDate}
          on:blur={() => showDateInput = false}
          on:change={() => showDateInput = false}
        />
      {:else}
        <button class="date-display" on:click={() => showDateInput = true} title="Click to pick a date">
          {formatDisplayDate(selectedDate)}
        </button>
      {/if}

      <button class="nav-btn" on:click={() => shiftDate(1)} title="Next day" disabled={isToday}>&#8250;</button>
      <button class="today-btn" on:click={goToday} disabled={isToday}>Today</button>
    </div>

    <div class="controls-right">
      <!-- View mode toggle -->
      <div class="view-toggle">
        {#each ['day','7d','14d','30d'] as mode}
          <button
            class="mode-btn"
            class:active={viewMode === mode}
            on:click={() => { viewMode = mode; if (mode !== 'day') loadRangeData(); }}
          >{mode}</button>
        {/each}
      </div>

      <!-- Min duration filter (day mode only) -->
      {#if viewMode === 'day'}
        <div class="filter-group">
          <label for="min-dur">Show:</label>
          <select id="min-dur" bind:value={minDurationMs}>
            <option value={0}>All</option>
            <option value={60000}>&ge; 1m</option>
            <option value={300000}>&ge; 5m</option>
            <option value={900000}>&ge; 15m</option>
          </select>
        </div>
      {/if}

      <button class="refresh-button" on:click={loadData} aria-label="Refresh data" title="Refresh data">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
        </svg>
      </button>
    </div>
  </div>

  <!-- ── Error state ───────────────────────────────────────────────────────── -->
  {#if loadError}
    <div class="error-box">
      <span>{loadError}</span>
      <button on:click={loadData}>Retry</button>
    </div>
  {/if}

  <!-- ── Loading overlay ───────────────────────────────────────────────────── -->
  {#if isLoading}
    <div class="loading">Loading data…</div>
  {:else}

    <!-- ── DAY MODE ────────────────────────────────────────────────────────── -->
    {#if viewMode === 'day'}

      <!-- Summary bar -->
      <DaySummaryBar
        appUsageData={appUsageData}
        categoryUsageData={categoryUsageData}
        eventCount={eventLog.length}
      />

      <!-- Comparison strip -->
      <ComparisonStrip
        todayMs={categoryUsageData.filter(c => c.id !== 0).reduce((s,c) => s + c.timeElapsed, 0)}
        {yesterdayMs}
        {weekAvgMs}
      />

      <!-- Goals panel -->
      <GoalsPanel categoryUsageData={categoryUsageData} />

      <!-- Charts row -->
      <div class="chart-container">
        <div class="chart-box">
          <h2>Application Usage</h2>
          <AppUsageChart data={filteredAppUsageData} />
        </div>

        <div class="chart-box">
          <h2>Categories
            {#if selectedCategoryFilter !== null}
              <span class="filter-tag">
                Filtered
                <button class="clear-tag" on:click={() => selectedCategoryFilter = null}>✕</button>
              </span>
            {/if}
          </h2>
          <CategoryChart
            data={categoryUsageData}
            selectedCategoryId={selectedCategoryFilter}
            on:categorySelected={handleCategoryFilter}
          />
        </div>
      </div>

      <!-- Activity timeline -->
      <ActivityTimeline events={eventLog} categoryUsageData={categoryUsageData} />

      <!-- App details table -->
      <div class="data-table">
        <h2>Application Details
          {#if minDurationMs > 0 || selectedCategoryFilter !== null}
            <span class="filter-count">({filteredAppUsageData.length} of {appUsageData.length})</span>
          {/if}
        </h2>
        <table>
          <thead>
            <tr>
              <th>Application</th>
              <th>Category</th>
              <th>Time Spent</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredAppUsageData as app}
              <tr>
                <td>{app.appName}</td>
                <td>
                  <span class="cat-badge" on:click={() => selectedCategoryFilter = selectedCategoryFilter === app.categoryId ? null : app.categoryId} title="Filter by this category">
                    {app.categoryName}
                  </span>
                </td>
                <td>{formatTimeElapsed(app.timeElapsed)}</td>
              </tr>
            {/each}
            {#if filteredAppUsageData.length === 0}
              <tr><td colspan="3" class="empty-cell">No apps match the current filter.</td></tr>
            {/if}
          </tbody>
        </table>
      </div>

      <!-- Bottom panels row -->
      <div class="bottom-panels">
        <TopDomainsPanel events={eventLog} />
        <UncategorizedAppsPanel apps={uncategorizedApps} on:createRule={handleCreateRule} />
      </div>

    <!-- ── RANGE MODE ──────────────────────────────────────────────────────── -->
    {:else}
      <div class="chart-box">
        <h2>Category Trends — last {viewMode}</h2>
        <TrendChart data={trendData} />
      </div>
    {/if}

    <!-- Heatmap calendar — always visible -->
    <div class="chart-box">
      <h2>Activity Calendar — {selectedYear}</h2>
      <HeatmapCalendar data={calendarData} year={selectedYear} on:daySelected={handleDaySelected} />
    </div>

  {/if}
</div>

<!-- Create Rule modal (launched from Uncategorized panel) -->
{#if showCreateRuleModal}
  <CreateRuleModal
    show={showCreateRuleModal}
    {categories}
    prefillAppName={prefillAppName}
    on:close={() => { showCreateRuleModal = false; prefillAppName = ''; }}
    on:ruleAdded={() => { showCreateRuleModal = false; loadData(); }}
  />
{/if}

<style>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  /* ── Controls ──────────────────────────────────────────────── */
  .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .date-nav {
    display: flex;
    align-items: center;
    gap: 0.35rem;
  }

  .nav-btn {
    background: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 50%;
    width: 2rem;
    height: 2rem;
    font-size: 1.3rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
  }
  .nav-btn:hover:not(:disabled) { background: var(--button-hover-bg-color); }
  .nav-btn:disabled { opacity: 0.35; cursor: default; }

  .date-display {
    background: none;
    border: 1px solid var(--input-border-color);
    border-radius: 6px;
    padding: 0.35rem 0.75rem;
    font-size: 0.9rem;
    color: var(--text-color);
    cursor: pointer;
    font-weight: 500;
    white-space: nowrap;
  }
  .date-display:hover { border-color: var(--button-bg-color); }

  .date-input-popup {
    padding: 0.35rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 6px;
    font-size: 0.9rem;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
  }

  .today-btn {
    background: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 6px;
    padding: 0.3rem 0.65rem;
    font-size: 0.8rem;
    cursor: pointer;
    margin-left: 0.2rem;
  }
  .today-btn:hover:not(:disabled) { background: var(--button-hover-bg-color); }
  .today-btn:disabled { opacity: 0.35; cursor: default; }

  .controls-right {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    flex-wrap: wrap;
  }

  .view-toggle {
    display: flex;
    border: 1px solid var(--input-border-color);
    border-radius: 6px;
    overflow: hidden;
  }
  .mode-btn {
    background: none;
    border: none;
    border-right: 1px solid var(--input-border-color);
    padding: 0.3rem 0.65rem;
    font-size: 0.8rem;
    color: var(--text-color);
    cursor: pointer;
  }
  .mode-btn:last-child { border-right: none; }
  .mode-btn.active { background: var(--button-bg-color); color: var(--button-text-color); }
  .mode-btn:hover:not(.active) { background: var(--table-row-hover); }

  .filter-group {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    font-size: 0.85rem;
    color: var(--text-color);
  }

  .filter-group select {
    padding: 0.3rem 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    font-size: 0.85rem;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
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
  }
  .refresh-button:hover { background-color: var(--button-hover-bg-color); }
  .icon { width: 1.1rem; height: 1.1rem; }

  /* ── Feedback ──────────────────────────────────────────────── */
  .loading {
    display: flex;
    justify-content: center;
    padding: 2rem;
    color: var(--text-color);
  }

  .error-box {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: rgba(239,68,68,0.08);
    border: 1px solid #ef444430;
    border-radius: 6px;
    color: var(--text-color);
    font-size: 0.875rem;
  }
  .error-box button {
    margin-left: auto;
    font-size: 0.8rem;
    text-decoration: underline;
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
  }

  /* ── Charts ────────────────────────────────────────────────── */
  .chart-container {
    display: flex;
    gap: 1.25rem;
  }

  .chart-box {
    flex: 1;
    background-color: var(--card-bg-color);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--card-border-color);
  }

  h2 {
    margin: 0 0 1rem;
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-color);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .filter-tag {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    font-size: 0.72rem;
    font-weight: 600;
    background: #3b82f620;
    color: #3b82f6;
    border-radius: 99px;
    padding: 0.1rem 0.5rem;
  }

  .clear-tag {
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    padding: 0;
    font-size: 0.8rem;
    line-height: 1;
  }

  /* ── App table ─────────────────────────────────────────────── */
  .data-table {
    background-color: var(--card-bg-color);
    border-radius: 8px;
    padding: 1.25rem;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--card-border-color);
  }

  .filter-count {
    font-size: 0.8rem;
    font-weight: 400;
    color: var(--secondary-color);
  }

  table { width: 100%; border-collapse: collapse; }

  th, td {
    padding: 0.65rem 1rem;
    text-align: left;
    border-bottom: 1px solid var(--table-border-color);
    color: var(--text-color);
    font-size: 0.875rem;
  }

  th {
    font-weight: 600;
    background-color: var(--table-header-bg);
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.03em;
    color: var(--secondary-color);
  }

  tr:hover { background-color: var(--table-row-hover); }

  .cat-badge {
    display: inline-block;
    padding: 0.15rem 0.5rem;
    border-radius: 99px;
    background: var(--table-border-color);
    font-size: 0.75rem;
    cursor: pointer;
    user-select: none;
  }
  .cat-badge:hover { background: var(--button-bg-color); color: var(--button-text-color); }

  .empty-cell {
    text-align: center;
    color: var(--secondary-color);
    font-style: italic;
  }

  /* ── Bottom panels ─────────────────────────────────────────── */
  .bottom-panels {
    display: flex;
    gap: 1.25rem;
    align-items: flex-start;
  }

  .bottom-panels > :global(*) { flex: 1; }

  /* ── Responsive ────────────────────────────────────────────── */
  @media (max-width: 768px) {
    .chart-container, .bottom-panels { flex-direction: column; }
  }
</style>
