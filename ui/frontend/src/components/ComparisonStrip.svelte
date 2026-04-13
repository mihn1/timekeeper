<script>
  import { formatTimeElapsed } from '../utils/formatters';

  export let todayMs = 0;
  export let yesterdayMs = 0;
  export let weekAvgMs = 0;

  $: vYesterdayMs = todayMs - yesterdayMs;
  $: vWeekMs = todayMs - weekAvgMs;

  function arrow(delta) {
    if (delta > 0) return '▲';
    if (delta < 0) return '▼';
    return '–';
  }

  function cls(delta) {
    if (delta > 0) return 'positive';
    if (delta < 0) return 'negative';
    return 'neutral';
  }
</script>

{#if yesterdayMs > 0 || weekAvgMs > 0}
  <div class="comparison-strip">
    {#if yesterdayMs > 0}
      <span class="cmp">
        vs yesterday:
        <span class={cls(vYesterdayMs)}>
          {arrow(vYesterdayMs)} {formatTimeElapsed(Math.abs(vYesterdayMs))}
        </span>
      </span>
    {/if}
    {#if weekAvgMs > 0}
      <span class="sep">·</span>
      <span class="cmp">
        vs 7-day avg:
        <span class={cls(vWeekMs)}>
          {arrow(vWeekMs)} {formatTimeElapsed(Math.abs(vWeekMs))}
        </span>
      </span>
    {/if}
  </div>
{/if}

<style>
  .comparison-strip {
    font-size: 0.8rem;
    color: var(--secondary-color);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.25rem 0.25rem;
  }

  .cmp { display: inline-flex; align-items: center; gap: 0.3rem; }
  .sep { opacity: 0.4; }

  .positive { color: #22c55e; font-weight: 600; }
  .negative { color: #ef4444; font-weight: 600; }
  .neutral  { color: var(--secondary-color); }
</style>
