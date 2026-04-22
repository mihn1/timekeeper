<script>
  import { rerunStatus } from '../stores/rerun';

  export let isEnabled = false;
  export let onToggle = () => {};

  $: status = $rerunStatus;
  $: hasJob = status && status.state && status.state !== 'idle';
  $: progressPct = status?.totalEvents > 0
    ? Math.min(100, Math.round((status.processedEvents / status.totalEvents) * 100))
    : 0;
  $: rerunLabel = buildRerunLabel(status, progressPct);

  function buildRerunLabel(s, pct) {
    if (!s) return '';
    switch (s.state) {
      case 'running':
        return `Rerunning rules ${s.startDate} → ${s.endDate}: ${s.processedEvents}/${s.totalEvents} (${pct}%)`;
      case 'completed':
        return `Rerun complete ${s.startDate} → ${s.endDate}: ${s.totalEvents} events`;
      case 'failed':
        return `Rerun failed: ${s.errorMessage || 'unknown error'}`;
      default:
        return '';
    }
  }
</script>

<div class="status-bar">
  <div class="status-left">
    <div class="status">
      Status: <span class:active={isEnabled}>{isEnabled ? 'Active' : 'Paused'}</span>
    </div>

    {#if hasJob}
      <div class="rerun-status" class:running={status.state === 'running'} class:failed={status.state === 'failed'}>
        {#if status.state === 'running'}
          <div class="spinner"></div>
        {/if}
        <span class="rerun-label">{rerunLabel}</span>
        {#if status.state === 'running' && status.totalEvents > 0}
          <div class="progress-track">
            <div class="progress-bar" style="width: {progressPct}%"></div>
          </div>
        {/if}
      </div>
    {/if}
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
    padding: 0.5rem 1.5rem;
    background-color: var(--card-bg-color);
    border-top: 1px solid var(--card-border-color);
    gap: 1rem;
  }

  .status-left {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex: 1;
    min-width: 0;
  }

  .status {
    font-size: 0.9rem;
    color: var(--text-color);
    flex-shrink: 0;
  }

  .status span {
    font-weight: 600;
    color: var(--secondary-color);
  }

  .status span.active {
    color: var(--success-color);
  }

  .rerun-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
    color: var(--text-color);
    opacity: 0.85;
    min-width: 0;
    flex: 1;
  }

  .rerun-label {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .rerun-status.failed .rerun-label {
    color: var(--danger-color);
  }

  .progress-track {
    width: 120px;
    height: 4px;
    background-color: var(--card-border-color);
    border-radius: 2px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .progress-bar {
    height: 100%;
    background-color: var(--primary-color, #3b82f6);
    transition: width 0.2s ease;
  }

  .spinner {
    width: 12px;
    height: 12px;
    border: 2px solid var(--card-border-color);
    border-top-color: var(--primary-color, #3b82f6);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    flex-shrink: 0;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
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
    flex-shrink: 0;
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
