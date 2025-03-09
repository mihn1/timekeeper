<script>
  import { onMount } from 'svelte';
  import { GetAppUsageData } from '../../wailsjs/go/main/App';
  import { refreshData } from '../stores/timekeeper';
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
  }

  input[type="date"] {
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 0.9rem;
  }

  .chart-container {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .chart-box {
    flex: 1;
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  h2 {
    margin-top: 0;
    margin-bottom: 1rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: #333;
  }

  .data-table {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid #eee;
  }

  th {
    font-weight: 600;
    color: #555;
  }

  .loading {
    display: flex;
    justify-content: center;
    padding: 2rem;
    color: #666;
  }

  @media (max-width: 768px) {
    .chart-container {
      flex-direction: column;
    }
  }
</style>