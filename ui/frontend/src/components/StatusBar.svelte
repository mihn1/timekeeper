<script>
  export let isEnabled = false;
  export let onToggle = () => {};
</script>

<div class="status-bar">
  <div class="status">
    Status: <span class:active={isEnabled}>{isEnabled ? 'Active' : 'Paused'}</span>
  </div>
  
  <button class="toggle-button cursor-pointer" class:active={isEnabled} on:click={onToggle}>
    {isEnabled ? 'Pause Tracking' : 'Start Tracking'}
  </button>
</div>

<style>
  .status-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1.5rem;
    background-color: var(--card-bg-color);
    border-top: 1px solid var(--card-border-color);
  }

  .status {
    font-size: 0.9rem;
    color: var(--text-color);
  }

  .status span {
    font-weight: 600;
    color: var(--secondary-color);
  }

  .status span.active {
    color: var(--success-color);
  }

  .toggle-button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .toggle-button:hover {
    background-color: var(--button-hover-bg-color);
  }

  .toggle-button.active {
    background-color: var(--danger-color);
    color: white;
  }

  .toggle-button.active:hover {
    background-color: var(--danger-color);
    opacity: 0.9;
  }
</style>