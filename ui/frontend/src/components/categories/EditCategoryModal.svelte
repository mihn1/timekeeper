<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { UpdateCategory } from '../../../wailsjs/go/main/App';
  import Modal from '../common/Modal.svelte';
  import { dtos } from '../../../wailsjs/go/models';
  
  export let show = false;
  export let category = null;
  
  const dispatch = createEventDispatcher();
  
  type FormErrors = {
    name?: string;
    [key: string]: string;
  };
  
  let editedCategory: dtos.CategoryUpdate = null;
  
  let formErrors: FormErrors = {};
  let initialized = false;
  
  // Update local state when category prop changes or show becomes true
  $: if (category && show && !initialized) {
    editedCategory = {
      id: category.id !== undefined ? category.id : (category.Id || ''),
      name: category.name !== undefined ? category.name : (category.Name || ''),
      description: category.description !== undefined ? category.description : (category.Description || '')
    };
    initialized = true;
  }

  // Reset the initialization flag when the modal is closed
  $: if (!show) {
    editedCategory = null;
    initialized = false;
  }
  
  function validateForm() {
    formErrors = {};
    let isValid = true;
    
    if (!editedCategory.name || !editedCategory.name.trim()) {
      formErrors.name = 'Category Name is required';
      isValid = false;
    }
    
    return isValid;
  }

  async function saveCategory() {
    if (!validateForm()) return;
    
    try {
      console.log('Saving category:', editedCategory);
      await UpdateCategory(editedCategory.id, editedCategory);
      dispatch('categoryEdited');
    } catch (err) {
      console.error('Error updating category:', err);
    }
  }
  
  function close() {
    dispatch('close');
  }
</script>

<Modal {show} title="Edit Category" on:close={close}>
  <div class="p-4">
    <div class="grid md:grid-cols-2 gap-4">
      <div>
        <label for="category-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
        <input 
          id="category-name"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.name ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={editedCategory.name}
        />
        {#if formErrors.name}
          <p class="text-red-500 text-xs mt-1">{formErrors.name}</p>
        {/if}
      </div>
      
      <div>
        <label for="category-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <textarea 
          id="category-description"
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={editedCategory.description}
          rows="3"
        ></textarea>
      </div>
    </div>
    
    <div class="mt-6 flex justify-end gap-2">
      <button 
        class="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400 focus:outline-none focus:ring-2 cursor-pointer"
        on:click={close}
      >
        Cancel
      </button>
      <button 
        class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer"
        on:click={saveCategory}
      >
        Save Changes
      </button>
    </div>
  </div>
</Modal>
