<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { AddRule } from '../../../wailsjs/go/main/App';
  import Modal from '../common/Modal.svelte';
  import type { dtos } from '../../../wailsjs/go/models';
  
  export let show = false;
  export let categories = [];
  
  const dispatch = createEventDispatcher();
  
  type FormErrors = {
    categoryId?: string;
    appName?: string;
    [key: string]: string;
  };
  
  let newRule: dtos.RuleCreate = { 
    categoryId: 0, 
    appName: '', 
    additionalDataKey: '', 
    expression: '', 
    isRegex: false, 
    priority: 0,
    isExclusion: false
  };
  
  let formErrors: FormErrors = {};
  
  function validateForm() {
    formErrors = {};
    let isValid = true;
    
    if (newRule.categoryId == 0) {
      formErrors.categoryId = 'Category ID is required';
      isValid = false;
    }
    
    if (!newRule.appName.trim()) {
      formErrors.appName = 'App Name is required';
      isValid = false;
    }
    
    return isValid;
  }

  async function addRule() {
    if (!validateForm()) return;
    
    try {
      await AddRule(newRule);
      resetForm();
      dispatch('ruleAdded');
    } catch (err) {
      console.error('Error adding rule:', err);
    }
  }
  
  function resetForm() {
    newRule = { 
      categoryId: 0, 
      appName: '', 
      additionalDataKey: '', 
      expression: '', 
      isRegex: false, 
      priority: 0,
      isExclusion: false
    };
    formErrors = {};
  }
  
  function close() {
    resetForm();
    dispatch('close');
  }
</script>

<Modal {show} title="Create New Rule" on:close={close}>
  <div class="p-4">
    <div class="grid md:grid-cols-2 gap-4">
      <div>
        <label for="category-select" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
        <select 
          id="category-select"
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.categoryId ? 'border-red-500' : 'border-gray-300'}"
          bind:value={newRule.categoryId}
        >
          <option value="">Select a category...</option>
          {#each categories as category}
            <option value={category.id}>{category.name}</option>
          {/each}
        </select>
        {#if formErrors.categoryId}
          <p class="text-red-500 text-xs mt-1">{formErrors.categoryId}</p>
        {/if}
      </div>
      
      <div>
        <label for="app-name" class="block text-sm font-medium text-gray-700 mb-1">App Name</label>
        <input 
          id="app-name"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.appName ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={newRule.appName}
        />
        {#if formErrors.appName}
          <p class="text-red-500 text-xs mt-1">{formErrors.appName}</p>
        {/if}
      </div>
      
      <div>
        <label for="additional-data-key" class="block text-sm font-medium text-gray-700 mb-1">Additional Data Key</label>
        <input 
          id="additional-data-key"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={newRule.additionalDataKey}
        />
      </div>
      
      <div>
        <label for="expression" class="block text-sm font-medium text-gray-700 mb-1">Expression</label>
        <input 
          id="expression"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={newRule.expression}
        />
      </div>
      
      <div>
        <label for="priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
        <input 
          id="priority"
          type="number" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={newRule.priority}
        />
      </div>
      
      <div class="flex items-center">
        <label class="inline-flex items-center mt-4">
          <input type="checkbox" class="form-checkbox h-5 w-5 text-blue-600" bind:checked={newRule.isRegex} />
          <span class="ml-2 text-sm text-gray-700">Is Regex</span>
        </label>
      </div>
      
      <div class="flex items-center">
        <label class="inline-flex items-center mt-4">
          <input type="checkbox" class="form-checkbox h-5 w-5 text-blue-600" bind:checked={newRule.isExclusion} />
          <span class="ml-2 text-sm text-gray-700">Is Exclusion</span>
        </label>
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
        on:click={addRule}
      >
        Add Rule
      </button>
    </div>
  </div>
</Modal>