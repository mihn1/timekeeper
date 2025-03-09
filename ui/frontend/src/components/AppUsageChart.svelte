<script>
  import { onMount } from 'svelte';
  import { Doughnut } from 'svelte-chartjs';
  import { formatTimeElapsed } from '../utils/formatters';
  
  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    ArcElement,
    CategoryScale
  } from 'chart.js';
  
  ChartJS.register(
    Title,
    Tooltip,
    Legend,
    ArcElement,
    CategoryScale
  );
  
  export let data = [];
  
  let chartData = {
    labels: [],
    datasets: []
  };
  
  $: if (data) {
    prepareChartData();
  }
  
  function prepareChartData() {
    // Take top 10 apps by usage
    const topApps = data.slice(0, 10);
    
    // Calculate "Other" if there are more than 10 apps
    let otherTime = 0;
    if (data.length > 10) {
      for (let i = 10; i < data.length; i++) {
        otherTime += data[i].TimeElapsed || 0;
      }
    }
    
    // Prepare chart labels and data
    const labels = topApps.map(app => app.AppName);
    const values = topApps.map(app => app.TimeElapsed / 60000); // Convert to minutes
    
    if (otherTime > 0) {
      labels.push('Other');
      values.push(otherTime / 60000);
    }
    
    // Generate nice colors
    const colors = generateColors(labels.length);
    
    chartData = {
      labels: labels,
      datasets: [
        {
          data: values,
          backgroundColor: colors,
          hoverOffset: 4
        }
      ]
    };
  }
  
  function generateColors(count) {
    const colors = [
      '#4E79A7', '#F28E2C', '#E15759', '#76B7B2', '#59A14F',
      '#EDC949', '#AF7AA1', '#FF9DA7', '#9C755F', '#BAB0AB'
    ];
    
    // If we need more colors than in our set, we'll generate some
    if (count > colors.length) {
      for (let i = colors.length; i < count; i++) {
        const h = (i * 137.5) % 360;
        colors.push(`hsl(${h}, 70%, 60%)`);
      }
    }
    
    return colors.slice(0, count);
  }
  
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'right',
        labels: {
          boxWidth: 15,
          padding: 15,
          font: {
            size: 12
          }
        }
      },
      tooltip: {
        callbacks: {
          label: function(context) {
            const value = context.raw;
            const hours = Math.floor(value / 60);
            const mins = Math.floor(value % 60);
            return `${context.label}: ${hours}h ${mins}m`;
          }
        }
      }
    }
  };
</script>

<div class="chart-wrapper">
  {#if data.length > 0}
    <Doughnut data={chartData} />
  {:else}
    <div class="no-data">No application usage data for this date</div>
  {/if}
</div>

<style>
  .chart-wrapper {
    height: 300px;
    position: relative;
  }
  
  .no-data {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: #777;
    font-style: italic;
  }
</style>