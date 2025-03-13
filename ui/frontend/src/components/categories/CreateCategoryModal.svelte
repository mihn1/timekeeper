<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { AddCategory } from '../../../wailsjs/go/main/App';
  import Modal from '../common/Modal.svelte';
  import { dtos } from '../../../wailsjs/go/models';
  
  export let show = false;
  
  const dispatch = createEventDispatcher();
  
  type FormErrors = {
    Name?: string;
    [key: string]: string;
  };
  
  let newCategory: dtos.CategoryCreate = { 
    name: '',
    description: '',
  };
  
  let formErrors: FormErrors = {};
  
  function validateForm() {
    formErrors = {};
    let isValid = true;
    
    if (!newCategory.name.trim()) {
      formErrors.Name = 'Category Name is required';
      isValid = false;
    }
    
    return isValid;
  }

  async function addCategory() {
    if (!validateForm()) return;
    
    try {
      await AddCategory(newCategory);
      resetForm();
      dispatch('categoryAdded');
    } catch (err) {
      console.error('Error adding category:', err);
    }
  }
  
  function resetForm() {
    newCategory = { 
      name: '',
      description: '',
    };
    formErrors = {};
  }
  
  function close() {
    resetForm();
    dispatch('close');
  }
</script>

<Modal {show} title="Create New Category" on:close={close}>
  <div class="p-4">
    <div class="grid md:grid-cols-2 gap-4">
      <div>
        <label for="category-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
        <input 
          id="category-name"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.Name ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={newCategory.name}
        />
        {#if formErrors.Name}
          <p class="text-red-500 text-xs mt-1">{formErrors.Name}</p>
        {/if}
      </div>
      
      <div>
        <label for="category-description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <input 
          id="category-description"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.Description ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={newCategory.description}
        />
        {#if formErrors.Description}
          <p class="text-red-500 text-xs mt-1">{formErrors.Description}</p>
        {/if}
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
        on:click={addCategory}
      >
        Add Category
      </button>
    </div>
  </div>
</Modal>