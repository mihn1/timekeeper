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
    background-color: #ffffff;
    border-top: 1px solid #eaeaea;
  }

  .status {
    font-size: 0.9rem;
    color: #666;
  }

  .status span {
    font-weight: 600;
    color: #777;
  }

  .status span.active {
    color: #1E8E3E;
  }

  .toggle-button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    background-color: #f0f0f0;
    color: #666;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .toggle-button:hover {
    background-color: #e0e0e0;
  }

  .toggle-button.active {
    background-color: #f9aea4;
    color: #c53030;
  }

  .toggle-button.active:hover {
    background-color: #f79b8e;
  }
</style>