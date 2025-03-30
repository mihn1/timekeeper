<script>
  import { onMount } from 'svelte';
  import { GetAppUsageData } from '../../wailsjs/go/main/App';
  import { formatTimeElapsed } from '../utils/formatters';
  import AppUsageChart from './AppUsageChart.svelte';
  import CategoryChart from './CategoryChart.svelte';
  
  let selectedDate = new Date().toISOString().split('T')[0];
  let appUsageData = [];
  let isLoading = true;

  // Subscribe to refresh events
  $: if ($refreshData && selectedDate) {
    loadData();
  }

  onMount(() => {
    loadData();
  });

  async function loadData() {
    isLoading = true;
    try {
      appUsageData = await GetAppUsageData(selectedDate);
      // Sort by time elapsed (descending)
      appUsageData.sort((a, b) => b.TimeElapsed - a.TimeElapsed);
    } catch (err) {
      console.error('Error loading data:', err);
    } finally {
      isLoading = false;
    }
  }

  function handleDateChange(e) {
    selectedDate = e.target.value;
    loadData();
  }

  function refreshData() {
    loadData();
  }
</script>

<div class="dashboard">
  <div class="controls">
    <div class="date-picker">
      <label for="date-select">Select Date:</label>
      <input 
        type="date" 
        id="date-select"
        bind:value={selectedDate} 
        on:change={handleDateChange}
      />
    </div>
    
    <button 
      class="refresh-button" 
      on:click={refreshData}
      aria-label="Refresh data"
      title="Refresh data"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
      </svg>
    </button>
  </div>
  
  {#if isLoading}
    <div class="loading">Loading data...</div>
  {:else}
    <div class="chart-container">
      <div class="chart-box">
        <h2>Application Usage</h2>
        <AppUsageChart data={appUsageData} />
      </div>
      
      <div class="chart-box">
        <h2>Categories</h2>
        <CategoryChart date={selectedDate}  />
      </div>
    </div>
    
    <div class="data-table">
      <h2>Application Details</h2>
      <table>
        <thead>
          <tr>
            <th>Application</th>
            <th>Time Spent</th>
          </tr>
        </thead>
        <tbody>
          {#each appUsageData as app}
            <tr>
              <td>{app.AppName}</td>
              <td>{formatTimeElapsed(app.TimeElapsed)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .dashboard {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .date-picker {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-color);
  }

  input[type="date"] {
    padding: 0.5rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    font-size: 0.9rem;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
  }

  .refresh-button {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem;
    border-radius: 50%;
    border: none;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    cursor: pointer;
  }

  .refresh-button:hover {
    background-color: var(--button-hover-bg-color);
  }

  .chart-container {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .chart-box {
    flex: 1;
    background-color: var(--card-bg-color);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--card-border-color);
  }

  h2 {
    margin-top: 0;
    margin-bottom: 1rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .data-table {
    background-color: var(--card-bg-color);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--card-border-color);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid var(--table-border-color);
    color: var(--text-color);
  }

  th {
    font-weight: 600;
    background-color: var(--table-header-bg);
  }

  tr:hover {
    background-color: var(--table-row-hover);
  }

  .loading {
    display: flex;
    justify-content: center;
    padding: 2rem;
    color: var(--text-color);
  }

  @media (max-width: 768px) {
    .chart-container {
      flex-direction: column;
    }
  }
</style>