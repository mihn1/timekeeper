<script context="module">
  import { currentView } from '../stores/navigation';

  function navigate(view) {
    currentView.set(view);
  }
</script>

<nav class="menu">
  <div class="menu-items">
    <button
      class="cursor-pointer menu-item {$currentView === 'dashboard' ? 'active' : ''}"
      on:click={() => navigate('dashboard')}
    >Dashboard</button>
    <button
      class="cursor-pointer menu-item {$currentView === 'rules' ? 'active' : ''}"
      on:click={() => navigate('rules')}
    >Rules</button>
    <button
      class="cursor-pointer menu-item {$currentView === 'categories' ? 'active' : ''}"
      on:click={() => navigate('categories')}
    >Categories</button>
    <button
      class="cursor-pointer menu-item {$currentView === 'goals' ? 'active' : ''}"
      on:click={() => navigate('goals')}
    >Goals</button>
    <button
      class="cursor-pointer menu-item {$currentView === 'events' ? 'active' : ''}"
      on:click={() => navigate('events')}
    >Event Log</button>
  </div>

  <button
    class="cursor-pointer menu-item prefs-btn {$currentView === 'preferences' ? 'active' : ''}"
    on:click={() => navigate('preferences')}
    title="Preferences"
    aria-label="Preferences"
  >
    ⚙ Preferences
  </button>
</nav>

<style>
  .menu {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background-color: var(--menu-bg-color);
  }

  .menu-items {
    display: flex;
    gap: 1rem;
  }

  button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    cursor: pointer;
    transition: background-color 0.2s;
  }

  button:hover {
    background-color: var(--button-hover-bg-color);
  }

  .menu-item.active {
    background-color: var(--button-hover-bg-color);
    position: relative;
  }

  .menu-item.active::after {
    content: '';
    position: absolute;
    bottom: -4px;
    left: 0;
    width: 100%;
    height: 3px;
    background-color: var(--primary-color);
  }

  .prefs-btn {
    font-size: 0.875rem;
  }
</style>
