<script>
  import { formatTimeElapsed } from '../utils/formatters';

  export let events = [];

  $: domains = (() => {
    const map = new Map();
    for (const ev of events) {
      let hostname = null;
      try {
        const url = ev.urlOrTitle;
        if (url && (url.startsWith('http://') || url.startsWith('https://'))) {
          hostname = new URL(url).hostname.replace(/^www\./, '');
        }
      } catch {}
      if (!hostname) continue;
      const prev = map.get(hostname) ?? 0;
      map.set(hostname, prev + ev.durationSecs * 1000);
    }
    return Array.from(map.entries())
      .map(([domain, ms]) => ({ domain, ms }))
      .sort((a, b) => b.ms - a.ms)
      .slice(0, 10);
  })();
</script>

{#if domains.length > 0}
  <div class="panel">
    <h3>Top Domains</h3>
    <ul class="domain-list">
      {#each domains as { domain, ms }, i}
        <li class="domain-item">
          <span class="rank">{i + 1}</span>
          <span class="domain">{domain}</span>
          <span class="duration">{formatTimeElapsed(ms)}</span>
        </li>
      {/each}
    </ul>
  </div>
{/if}

<style>
  .panel {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    box-shadow: var(--card-shadow);
  }

  h3 {
    margin: 0 0 0.75rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .domain-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .domain-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.3rem 0;
    border-bottom: 1px solid var(--table-border-color);
    font-size: 0.85rem;
  }

  .domain-item:last-child { border-bottom: none; }

  .rank {
    width: 1.2rem;
    text-align: right;
    color: var(--secondary-color);
    font-size: 0.75rem;
    flex-shrink: 0;
  }

  .domain { flex: 1; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .duration {
    color: var(--secondary-color);
    white-space: nowrap;
    font-size: 0.8rem;
  }
</style>
