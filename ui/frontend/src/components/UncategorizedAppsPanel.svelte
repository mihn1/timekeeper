<script>
  import { createEventDispatcher } from 'svelte';

  export let apps = [];

  const dispatch = createEventDispatcher();
</script>

{#if apps.length > 0}
  <div class="panel">
    <div class="panel-header">
      <h3>Uncategorized Apps</h3>
      <span class="badge">{apps.length}</span>
    </div>
    <p class="hint">These apps had no matching rule. Create a rule to categorize them.</p>
    <ul class="app-list">
      {#each apps as appName}
        <li class="app-item">
          <span class="app-name">{appName}</span>
          <button
            class="create-btn"
            on:click={() => dispatch('createRule', { appName })}
          >
            + Rule
          </button>
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

  .panel-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.25rem;
  }

  h3 {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-color);
  }

  .badge {
    background-color: #f59e0b;
    color: #fff;
    border-radius: 99px;
    font-size: 0.7rem;
    font-weight: 700;
    padding: 0.1rem 0.45rem;
    line-height: 1.4;
  }

  .hint {
    font-size: 0.8rem;
    color: var(--secondary-color);
    margin: 0 0 0.75rem;
  }

  .app-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .app-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.3rem 0;
    border-bottom: 1px solid var(--table-border-color);
  }

  .app-item:last-child { border-bottom: none; }

  .app-name {
    font-size: 0.85rem;
    color: var(--text-color);
  }

  .create-btn {
    font-size: 0.75rem;
    padding: 0.2rem 0.6rem;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    white-space: nowrap;
  }

  .create-btn:hover { background-color: var(--button-hover-bg-color); }
</style>
