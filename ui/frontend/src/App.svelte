<script>
  import { onMount, onDestroy } from 'svelte';
  import { currentView } from './stores/navigation';
  import { trackingEnabled, refreshData } from './stores/timekeeper';
  import { theme } from './stores/theme';
  import { IsTrackingEnabled, EnableTracking, DisableTracking } from '../wailsjs/go/main/App';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import Dashboard from './components/Dashboard.svelte';
  import Rules from './components/rules/Rules.svelte';
  import Categories from './components/categories/Categories.svelte';
  import Menu from './components/Menu.svelte';
  import StatusBar from './components/StatusBar.svelte';

  let isInitialized = false;
  let unsubscribe;

  $: currentTheme = $theme;

  onMount(async () => {
    const enabled = await IsTrackingEnabled();
    trackingEnabled.set(enabled);

    unsubscribe = EventsOn('timekeeper:data-updated', () => {
      refreshData.set(Date.now());
    });

    isInitialized = true;
  });

  onDestroy(() => {
    if (unsubscribe) unsubscribe();
  });

  async function toggleTracking() {
    const enabled = $trackingEnabled;

    if (enabled) {
      await DisableTracking();
      trackingEnabled.set(false);
    } else {
      await EnableTracking();
      trackingEnabled.set(true);
    }
  }
</script>

<main class="app-container">
  <Menu />
  
  {#if isInitialized}
    <div class="content-area">
      {#if $currentView === 'dashboard'}
        <Dashboard />
      {:else if $currentView === 'rules'}
        <Rules />
      {:else if $currentView === 'categories'}
        <Categories />
      {:else}
        <Dashboard />
      {/if}
    </div>
    
    <StatusBar 
      isEnabled={$trackingEnabled} 
      onToggle={toggleTracking} 
    />
  {:else}
    <div class="loading">
      <p>Initializing TimeKeeper...</p>
    </div>
  {/if}
</main>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    background-color: var(--background-color);
    color: var(--text-color);
  }
  
  .content-area {
    flex: 1;
    padding: 1rem;
  }

  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    font-size: 1.2rem;
    color: #555;
  }
</style>
