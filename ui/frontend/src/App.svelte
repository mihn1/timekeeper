<script>
  import logo from './assets/images/logo-universal.png'
  import {Greet} from '../wailsjs/go/main/App.js'
  import { onMount, onDestroy } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import { IsTrackingEnabled, EnableTracking, DisableTracking } from '../wailsjs/go/main/App';
  import Dashboard from './components/Dashboard.svelte';
  import Header from './components/Header.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import { trackingEnabled, refreshData } from './stores/timekeeper';

  let isInitialized = false;
  let unsubscribe;

  onMount(async () => {
    // Check if tracking is enabled
    const enabled = await IsTrackingEnabled();
    trackingEnabled.set(enabled);
    
    // Listen for real-time updates
    unsubscribe = EventsOn('timekeeper:data-updated', () => {
      refreshData.set(Date.now());
    });
    
    isInitialized = true;
  });

  onDestroy(() => {
    if (unsubscribe) unsubscribe();
  });

  // Toggle tracking functionality
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
  <Header />
  
  {#if isInitialized}
    <div class="content">
      <Dashboard />
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

  <!-- <img alt="Wails logo" id="logo" src="{logo}"> -->
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

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

</style>
