<script context="module">
  import { currentView } from '../stores/navigation';

  function navigate(view) {
    currentView.set(view);
  }
</script>

<script>
  import { theme, toggleTheme } from '../stores/theme';
</script>

<nav class="menu">
  <div class="menu-items">
    <button class="cursor-pointer" on:click={() => navigate('dashboard')}>Dashboard</button>
    <button class="cursor-pointer" on:click={() => navigate('rules')}>Rule Management</button>
    <button class="cursor-pointer" on:click={() => navigate('categories')}>Category Management</button>
  </div>
  <div class="theme-toggle">
    <button on:click={toggleTheme} aria-label="Toggle theme">
      {#if $theme === 'light'}
        🌙 <!-- Moon emoji for dark mode -->
      {:else}
        ☀️ <!-- Sun emoji for light mode -->
      {/if}
    </button>
  </div>
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

  .theme-toggle button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    background-color: var(--theme-toggle-bg);
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
</style>