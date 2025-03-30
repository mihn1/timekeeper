<script>
  import { trackingEnabled } from '../stores/timekeeper';
</script>

<header>
  <div class="logo">
    <img src="/logo.png" alt="TimeKeeper Logo" />
    <h1>TimeKeeper</h1>
  </div>
  
  <div class="status" class:active={$trackingEnabled}>
    {$trackingEnabled ? 'Tracking Active' : 'Tracking Paused'}
  </div>
</header>

<style>
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: var(--card-bg-color);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .logo img {
    width: 32px;
    height: 32px;
  }

  h1 {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-color);
    margin: 0;
  }

  .status {
    font-size: 0.9rem;
    padding: 0.35rem 0.75rem;
    border-radius: 1rem;
    background-color: var(--button-bg-color);
    color: var(--secondary-color);
  }

  .status.active {
    background-color: var(--success-color);
    color: white;
  }
</style>