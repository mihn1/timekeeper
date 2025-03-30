<script>
  import { createEventDispatcher, onMount } from 'svelte';
  
  export let show = false;
  export let title = 'Modal Title';
  
  const dispatch = createEventDispatcher();
  
  function closeModal() {
    dispatch('close');
  }
  
  function handleKeydown(event) {
    // Close modal on Escape key
    if (event.key === 'Escape') {
      closeModal();
    }
  }
  
  // Set up keyboard event listener when modal is shown
  onMount(() => {
    return () => {
      document.body.classList.remove('overflow-hidden');
    };
  });
  
  // Prevent scrolling when modal is open
  $: if (show) {
    document.body.classList.add('overflow-hidden');
    // Add event listener when modal is shown
    window.addEventListener('keydown', handleKeydown);
  } else {
    document.body.classList.remove('overflow-hidden');
    // Remove event listener when modal is hidden
    window.removeEventListener('keydown', handleKeydown);
  }
</script>

{#if show}
  <!-- Use role="dialog" and proper aria attributes -->
  <div 
    class="modal-backdrop" 
    role="presentation"
    on:click|self={closeModal}
  >
    <div 
      class="modal-container" 
      role="dialog"
      aria-labelledby="modal-title" 
      aria-modal="true"
      tabindex="-1"
    >
      <div class="border-b px-4 py-3 flex justify-between items-center">
        <h3 id="modal-title" class="text-lg font-semibold">{title}</h3>
        <!-- Use button element for close button -->
        <button 
          type="button"
          class="close-button focus:outline-none cursor-pointer" 
          on:click={closeModal}
          aria-label="Close modal"
        >
          <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
      <slot></slot>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 50;
  }

  .modal-container {
    background-color: var(--card-bg-color);
    border-radius: 0.5rem;
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: var(--card-shadow);
    color: var(--text-color);
  }

  /* Focus outline for keyboard navigation */
  .modal-container:focus {
    outline: 2px solid var(--primary-color);
    outline-offset: 2px;
  }
  
  .border-b {
    border-bottom: 1px solid var(--card-border-color);
  }
  
  .text-lg {
    font-size: 1.125rem;
  }
  
  .font-semibold {
    font-weight: 600;
  }
  
  .close-button {
    color: var(--text-color);
  }
  
  .close-button:hover {
    color: var(--primary-color);
  }

  @media (max-width: 640px) {
    .modal-container {
      max-width: 90%;
    }
  }
</style>
