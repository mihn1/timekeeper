<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Bar } from 'svelte-chartjs';
  import { formatTimeElapsed } from '../utils/formatters';
  import { theme } from '../stores/theme';
  import type { dtos } from '../../wailsjs/go/models';

  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale
  } from 'chart.js';

  ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

  // Accept data as a prop (Dashboard centralises fetching).
  export let data: dtos.CategoryUsageItem[] = [];
  // Optionally highlight a selected category (from drill-down).
  export let selectedCategoryId: number | null = null;

  const dispatch = createEventDispatcher();

  let isDarkMode: boolean;
  $: isDarkMode = $theme === 'dark';

  const BASE_COLORS = [
    '#4E79A7', '#F28E2C', '#E15759', '#76B7B2', '#59A14F',
    '#EDC949', '#AF7AA1', '#FF9DA7', '#9C755F', '#BAB0AB'
  ];

  $: chartData = {
    labels: data.map(cat => cat.name),
    datasets: [{
      label: 'Time Spent (minutes)',
      data: data.map(cat => cat.timeElapsed / 60000),
      backgroundColor: data.map((cat, i) => {
        const base = BASE_COLORS[i % BASE_COLORS.length];
        if (selectedCategoryId !== null && cat.id !== selectedCategoryId) {
          return base + '55'; // dim non-selected
        }
        return base;
      }),
    }]
  };

  function handleClick(_event: any, elements: any[]) {
    if (elements.length > 0) {
      const idx = elements[0].index;
      const cat = data[idx];
      const newId = selectedCategoryId === cat.id ? null : cat.id;
      dispatch('categorySelected', { categoryId: newId });
    }
  }

  $: chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    onClick: handleClick,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label(context: any) {
            const value = context.raw;
            const hours = Math.floor(value / 60);
            const mins  = Math.floor(value % 60);
            return `${hours}h ${mins}m`;
          }
        },
        backgroundColor: isDarkMode ? 'rgba(0,0,0,0.8)'   : 'rgba(255,255,255,0.8)',
        titleColor:      isDarkMode ? '#f0f0f0' : '#333333',
        bodyColor:       isDarkMode ? '#f0f0f0' : '#333333',
        borderColor:     isDarkMode ? '#444444' : '#e0e0e0',
        borderWidth: 1
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        grid:  { color: isDarkMode ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)' },
        ticks: {
          color: isDarkMode ? '#f0f0f0' : '#333333',
          callback(value: any) {
            const hours = Math.floor(value / 60);
            const mins  = Math.floor(value % 60);
            return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
          }
        }
      },
      x: {
        grid:  { color: isDarkMode ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)' },
        ticks: { color: isDarkMode ? '#f0f0f0' : '#333333' }
      }
    }
  };
</script>

<div class="chart-wrapper">
  {#if data.length > 0}
    <Bar data={chartData} options={chartOptions} />
    {#if selectedCategoryId !== null}
      <button class="clear-filter" on:click={() => dispatch('categorySelected', { categoryId: null })}>
        Clear filter
      </button>
    {/if}
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

  .no-data {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: var(--secondary-color);
    font-style: italic;
  }

  .clear-filter {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    font-size: 0.72rem;
    padding: 0.2rem 0.5rem;
    background: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    opacity: 0.8;
  }

  .clear-filter:hover { opacity: 1; }
</style>
