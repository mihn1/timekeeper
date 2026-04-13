<script>
  import { Line } from 'svelte-chartjs';
  import { theme } from '../stores/theme';

  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    LineElement,
    PointElement,
    LinearScale,
    CategoryScale,
  } from 'chart.js';

  ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, LinearScale, CategoryScale);

  // Array of DailyCategorySummary: { date, categoryId, categoryName, timeElapsed }
  export let data = [];

  let isDarkMode;
  $: isDarkMode = $theme === 'dark';

  const CAT_COLORS = {
    1: '#3b82f6',
    2: '#f59e0b',
    3: '#22c55e',
    4: '#9ca3af',
  };

  function catColor(catId) {
    return CAT_COLORS[catId] ?? `hsl(${(catId * 137) % 360}, 60%, 55%)`;
  }

  $: chartData = (() => {
    if (!data || data.length === 0) return { labels: [], datasets: [] };

    // Collect unique dates and categories.
    const dateSet = new Set(data.map(d => d.date));
    const dates   = Array.from(dateSet).sort();

    const catMap = new Map();
    for (const row of data) {
      if (!catMap.has(row.categoryId)) {
        catMap.set(row.categoryId, { name: row.categoryName, id: row.categoryId });
      }
    }

    // Pivot into per-category series.
    const datasets = Array.from(catMap.values()).map(cat => {
      const byDate = new Map(
        data.filter(d => d.categoryId === cat.id).map(d => [d.date, d.timeElapsed / 3600000])
      );
      return {
        label: cat.name,
        data: dates.map(d => byDate.get(d) ?? 0),
        borderColor: catColor(cat.id),
        backgroundColor: catColor(cat.id) + '33',
        tension: 0.3,
        pointRadius: 3,
        fill: false,
      };
    });

    const labels = dates.map(d => {
      const dt = new Date(d + 'T00:00:00');
      return dt.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
    });

    return { labels, datasets };
  })();

  $: chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
        labels: { color: isDarkMode ? '#f0f0f0' : '#333333', boxWidth: 12 }
      },
      tooltip: {
        callbacks: {
          label(ctx) {
            const h = Math.floor(ctx.raw);
            const m = Math.round((ctx.raw - h) * 60);
            return `${ctx.dataset.label}: ${h}h ${m}m`;
          }
        },
        backgroundColor: isDarkMode ? 'rgba(0,0,0,0.85)' : 'rgba(255,255,255,0.9)',
        titleColor: isDarkMode ? '#f0f0f0' : '#111',
        bodyColor:  isDarkMode ? '#f0f0f0' : '#333',
        borderColor: isDarkMode ? '#444' : '#ddd',
        borderWidth: 1,
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          color: isDarkMode ? '#f0f0f0' : '#333',
          callback: v => v === 0 ? '0' : `${v.toFixed(1)}h`,
        },
        grid: { color: isDarkMode ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)' }
      },
      x: {
        ticks: { color: isDarkMode ? '#f0f0f0' : '#333' },
        grid: { color: isDarkMode ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)' }
      }
    }
  };
</script>

<div class="chart-wrapper">
  {#if data.length > 0}
    <Line data={chartData} options={chartOptions} />
  {:else}
    <div class="no-data">No data for the selected range.</div>
  {/if}
</div>

<style>
  .chart-wrapper {
    height: 320px;
    position: relative;
    padding: 0.5rem;
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
