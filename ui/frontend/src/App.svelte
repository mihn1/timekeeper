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
  import EventLog from './components/events/EventLog.svelte';
  import Goals from './components/Goals.svelte';
  import Preferences from './components/Preferences.svelte';
  import Menu from './components/Menu.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import ToastContainer from './components/common/ToastContainer.svelte';

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
  
  <div class="content-area">
    {#if isInitialized}
      {#if $currentView === 'dashboard'}
        <Dashboard />
      {:else if $currentView === 'rules'}
        <Rules />
      {:else if $currentView === 'categories'}
        <Categories />
      {:else if $currentView === 'goals'}
        <Goals />
      {:else if $currentView === 'events'}
        <EventLog />
      {:else if $currentView === 'preferences'}
        <Preferences />
      {:else}
        <Dashboard />
      {/if}
    {:else}
      <div class="loading">
        <p>Initializing TimeKeeper...</p>
      </div>
    {/if}
  </div>

  <StatusBar
    isEnabled={$trackingEnabled}
    onToggle={toggleTracking}
  />

  <ToastContainer />
</main>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
    background-color: var(--background-color);
    color: var(--text-color);
  }

  .content-area {
    flex: 1;
    overflow-y: auto;
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
