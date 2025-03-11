<script context="module">
  import { currentView } from '../stores/navigation';

  function navigate(view) {
    currentView.set(view);
  }
</script>

<nav class="menu">
  <button on:click={() => navigate('dashboard')}>Dashboard</button>
  <button on:click={() => navigate('rules')}>Rule Management</button>
  <button on:click={() => navigate('categories')}>Category Management</button>
</nav>

<style>
  .menu {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    background-color: #f0f0f0;
  }

  button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    background-color: #ffffff;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  button:hover {
    background-color: #e0e0e0;
  }
</style>