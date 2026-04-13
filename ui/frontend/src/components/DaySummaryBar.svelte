<script>
  import { formatTimeElapsed } from '../utils/formatters';

  // IDs of categories counted as "productive" (Work = 1 by default).
  const PRODUCTIVE_IDS = [1];

  export let appUsageData = [];
  export let categoryUsageData = [];
  export let eventCount = 0;

  $: totalMs = categoryUsageData
    .filter(c => c.id !== 0) // exclude EXCLUDED category
    .reduce((sum, c) => sum + c.timeElapsed, 0);

  $: productiveMs = categoryUsageData
    .filter(c => PRODUCTIVE_IDS.includes(c.id))
    .reduce((sum, c) => sum + c.timeElapsed, 0);

  $: topApp = appUsageData.length > 0 ? appUsageData[0] : null;

  $: topCategory = categoryUsageData.length > 0
    ? [...categoryUsageData].sort((a, b) => b.timeElapsed - a.timeElapsed)[0]
    : null;

  $: topCategoryPct = totalMs > 0 && topCategory
    ? Math.round((topCategory.timeElapsed / totalMs) * 100)
    : 0;
</script>

{#if totalMs > 0}
  <div class="summary-bar">
    <div class="stat">
      <span class="label">Total</span>
      <span class="value">{formatTimeElapsed(totalMs)}</span>
    </div>
    <div class="divider"></div>
    <div class="stat">
      <span class="label">Productive</span>
      <span class="value productive">{formatTimeElapsed(productiveMs)}</span>
    </div>
    {#if topApp}
      <div class="divider"></div>
      <div class="stat">
        <span class="label">Top app</span>
        <span class="value">{topApp.appName} · {formatTimeElapsed(topApp.timeElapsed)}</span>
      </div>
    {/if}
    {#if topCategory}
      <div class="divider"></div>
      <div class="stat">
        <span class="label">Top category</span>
        <span class="value">{topCategory.name} · {topCategoryPct}%</span>
      </div>
    {/if}
    <div class="divider"></div>
    <div class="stat">
      <span class="label">Switches</span>
      <span class="value">{eventCount}</span>
    </div>
  </div>
{/if}

<style>
  .summary-bar {
    display: flex;
    align-items: center;
    gap: 0;
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 0.75rem 1.25rem;
    box-shadow: var(--card-shadow);
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    padding: 0 0.75rem;
  }

  .stat:first-child {
    padding-left: 0;
  }

  .label {
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--secondary-color);
    font-weight: 600;
  }

  .value {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text-color);
    white-space: nowrap;
  }

  .productive {
    color: #22c55e;
  }

  .divider {
    width: 1px;
    height: 2rem;
    background-color: var(--card-border-color);
    flex-shrink: 0;
  }
</style>
