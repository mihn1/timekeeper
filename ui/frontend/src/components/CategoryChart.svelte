<script>
  import { onMount } from 'svelte';
  import { Bar } from 'svelte-chartjs';
  import { formatTimeElapsed } from '../utils/formatters';
  import { refreshData } from '../stores/timekeeper';
  import { GetCategoryUsageData } from '../../wailsjs/go/main/App'; // Correct import path
  import { theme } from '../stores/theme';

  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale
  } from 'chart.js';

  ChartJS.register(
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale
  );

  export let date;

  let categoryData = [];
  let isLoading = true;
  let isDarkMode;

  $: isDarkMode = $theme === 'dark';

  $: if (date || $refreshData) {
    loadCategoryData();
  }

  onMount(() => {
    loadCategoryData();
  });

  async function loadCategoryData() {
    isLoading = true;
    try {
      categoryData = await GetCategoryUsageData(date); // Correct function call

      // Sort by time spent
      categoryData.sort((a, b) => b.TimeElapsed - a.TimeElapsed);
    } catch (err) {
      console.error('Error loading category data:', err);
      categoryData = [];
    } finally {
      isLoading = false;
    }
  }

  $: chartData = {
    labels: categoryData.map(cat => cat.Name),
    datasets: [{
      label: 'Time Spent (minutes)',
      data: categoryData.map(cat => cat.TimeElapsed / 60000), // Convert ms to mins
      backgroundColor: [
        '#4E79A7', '#F28E2C', '#E15759', '#76B7B2', '#59A14F',
        '#EDC949', '#AF7AA1', '#FF9DA7', '#9C755F', '#BAB0AB'
      ]
    }]
  };

  $: chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false
      },
      tooltip: {
        callbacks: {
          label: function(context) {
            const value = context.raw;
            const hours = Math.floor(value / 60);
            const mins = Math.floor(value % 60);
            return `${hours}h ${mins}m`;
          }
        },
        backgroundColor: isDarkMode ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.8)',
        titleColor: isDarkMode ? '#f0f0f0' : '#333333',
        bodyColor: isDarkMode ? '#f0f0f0' : '#333333',
        borderColor: isDarkMode ? '#444444' : '#e0e0e0',
        borderWidth: 1
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        grid: {
          color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'
        },
        ticks: {
          color: isDarkMode ? '#f0f0f0' : '#333333',
          callback: function(value) {
            const hours = Math.floor(value / 60);
            const mins = Math.floor(value % 60);
            return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
          }
        }
      },
      x: {
        grid: {
          color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'
        },
        ticks: {
          color: isDarkMode ? '#f0f0f0' : '#333333'
        }
      }
    }
  };
</script>

<div class="chart-wrapper">
  {#if isLoading}
    <div class="loading">Loading category data...</div>
  {:else if categoryData.length > 0}
    <Bar data={chartData} options={chartOptions} />
  {:else}
    <div class="no-data">No category data for this date</div>
  {/if}
</div>

<style>
  .chart-wrapper {
    height: 300px;
    position: relative;
    background-color: var(--card-bg-color);
    border-radius: 4px;
    padding: 1rem;
  }

  .loading, .no-data {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: var(--secondary-color);
    font-style: italic;
  }
</style>