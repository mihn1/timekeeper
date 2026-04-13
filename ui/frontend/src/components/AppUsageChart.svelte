<script>
  import { Doughnut } from 'svelte-chartjs';
  import { theme } from '../stores/theme';

  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    CategoryScale
  } from 'chart.js';

  ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale);

  // Accepts AppUsageItem[] (camelCase fields from the new DTO).
  export let data = [];

  let isDarkMode;
  $: isDarkMode = $theme === 'dark';

  $: if (data) prepareChartData();

  let chartData = { labels: [], datasets: [] };

  function prepareChartData() {
    const topApps = data.slice(0, 10);

    let otherTime = 0;
    for (let i = 10; i < data.length; i++) {
      otherTime += data[i].timeElapsed ?? 0;
    }

    const labels = topApps.map(app => app.appName);
    const values = topApps.map(app => app.timeElapsed / 60000);

    if (otherTime > 0) {
      labels.push('Other');
      values.push(otherTime / 60000);
    }

    const colors = generateColors(labels.length);
    chartData = {
      labels,
      datasets: [{ data: values, backgroundColor: colors, hoverOffset: 4 }]
    };
  }

  function generateColors(count) {
    const colors = [
      '#4E79A7', '#F28E2C', '#E15759', '#76B7B2', '#59A14F',
      '#EDC949', '#AF7AA1', '#FF9DA7', '#9C755F', '#BAB0AB'
    ];
    for (let i = colors.length; i < count; i++) {
      const h = (i * 137.5) % 360;
      colors.push(`hsl(${h}, 70%, 60%)`);
    }
    return colors.slice(0, count);
  }

  $: chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'right',
        labels: {
          boxWidth: 15,
          padding: 15,
          font: { size: 12 },
          color: isDarkMode ? '#f0f0f0' : '#333333'
        }
      },
      tooltip: {
        callbacks: {
          label(context) {
            const value = context.raw;
            const hours = Math.floor(value / 60);
            const mins  = Math.floor(value % 60);
            return `${context.label}: ${hours}h ${mins}m`;
          }
        },
        backgroundColor: isDarkMode ? 'rgba(0,0,0,0.8)'   : 'rgba(255,255,255,0.8)',
        titleColor:      isDarkMode ? '#f0f0f0' : '#333333',
        bodyColor:       isDarkMode ? '#f0f0f0' : '#333333',
        borderColor:     isDarkMode ? '#444444' : '#e0e0e0',
        borderWidth: 1
      }
    }
  };
</script>

<div class="chart-wrapper">
  {#if data.length > 0}
    <Doughnut data={chartData} options={chartOptions} />
  {:else}
    <div class="no-data">No application usage data for this date</div>
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
</style>
