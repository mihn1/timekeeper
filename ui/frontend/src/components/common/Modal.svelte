<script>
  import { createEventDispatcher } from 'svelte';
  
  export let show = false;
  export let title = "";
  
  const dispatch = createEventDispatcher();
  
  function closeModal() {
    dispatch('close');
  }
</script>

{#if show}
  <div class="fixed inset-0 z-50 overflow-auto bg-black bg-opacity-50 flex items-center justify-center">
    <div class="bg-white rounded-lg shadow-lg w-full max-w-md mx-4 md:mx-0">
      <div class="border-b px-4 py-3 flex justify-between items-center">
        <h3 class="text-lg font-semibold text-gray-800">{title}</h3>
        <button 
          class="text-gray-500 hover:text-gray-700 focus:outline-none cursor-pointer" 
          on:click={closeModal}
        >
          <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
      <div>
        <slot></slot>
      </div>
    </div>
  </div>
{/if}
