<script>
  import { onMount, onDestroy } from 'svelte';
  import { currentView } from './stores/navigation';
  import { trackingEnabled, refreshData } from './stores/timekeeper';
  import { IsTrackingEnabled, EnableTracking, DisableTracking } from '../wailsjs/go/main/App';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import Dashboard from './components/Dashboard.svelte';
  import Rules from './components/Rules.svelte';
  import Categories from './components/Categories.svelte';
  import Menu from './components/Menu.svelte';
  import StatusBar from './components/StatusBar.svelte';

  let isInitialized = false;
  let unsubscribe;

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

<main>
  <Menu />
  
  {#if isInitialized}
    <div class="content">
      {#if $currentView === 'dashboard'}
        <Dashboard />
      {:else if $currentView === 'rules'}
        <Rules />
      {:else if $currentView === 'categories'}
        <Categories />
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
  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background-color: #f5f5f7;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    color: #333;
  }

  .content {
    flex: 1;
    padding: 1rem;
    overflow-y: auto;
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
