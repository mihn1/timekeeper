<script>
  import { toasts, dismissToast } from '../../stores/toast.js';
</script>

<div class="toast-container" aria-live="polite" aria-atomic="true">
  {#each $toasts as toast (toast.id)}
    <div class="toast toast-{toast.type}" role="status">
      <span class="toast-msg">{toast.message}</span>
      <button class="toast-close" on:click={() => dismissToast(toast.id)} aria-label="Dismiss">×</button>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    bottom: 3rem;
    right: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    z-index: 100;
    pointer-events: none;
  }
  .toast {
    pointer-events: auto;
    min-width: 260px;
    max-width: 380px;
    padding: 0.6rem 0.9rem;
    border-radius: 6px;
    box-shadow: var(--card-shadow, 0 4px 12px rgba(0,0,0,0.18));
    color: #fff;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.9rem;
    animation: slide-in 0.18s ease-out;
  }
  .toast-success { background-color: #16a34a; }
  .toast-error   { background-color: #dc2626; }
  .toast-info    { background-color: #2563eb; }
  .toast-msg { flex: 1; }
  .toast-close {
    background: transparent;
    border: none;
    color: inherit;
    font-size: 1.1rem;
    line-height: 1;
    cursor: pointer;
    opacity: 0.85;
  }
  .toast-close:hover { opacity: 1; }
  @keyframes slide-in {
    from { transform: translateY(6px); opacity: 0; }
    to   { transform: translateY(0);   opacity: 1; }
  }
</style>
