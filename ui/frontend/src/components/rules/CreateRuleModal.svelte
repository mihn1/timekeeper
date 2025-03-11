<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { AddRule } from '../../../wailsjs/go/main/App';
  import Modal from '../common/Modal.svelte';
  
  export let show = false;
  export let categories = [];
  
  const dispatch = createEventDispatcher();
  
  type Rule = {
    RuleId: number;
    CategoryId: string;
    AppName: string;
    AdditionalDataKey: string;
    Expression: string;
    IsRegex: boolean;
    Priority: number;
  };
  
  type FormErrors = {
    CategoryId?: string;
    AppName?: string;
    Expression?: string;
    [key: string]: string;
  };
  
  let newRule: Rule = { 
    RuleId: 0, 
    CategoryId: '', 
    AppName: '', 
    AdditionalDataKey: '', 
    Expression: '', 
    IsRegex: false, 
    Priority: 0 
  };
  
  let formErrors: FormErrors = {};
  
  function validateForm() {
    formErrors = {};
    let isValid = true;
    
    if (!newRule.CategoryId.trim()) {
      formErrors.CategoryId = 'Category ID is required';
      isValid = false;
    }
    
    if (!newRule.AppName.trim()) {
      formErrors.AppName = 'App Name is required';
      isValid = false;
    }
    
    if (!newRule.Expression.trim()) {
      formErrors.Expression = 'Expression is required';
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
      RuleId: 0, 
      CategoryId: '', 
      AppName: '', 
      AdditionalDataKey: '', 
      Expression: '', 
      IsRegex: false, 
      Priority: 0 
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
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.CategoryId ? 'border-red-500' : 'border-gray-300'}"
          bind:value={newRule.CategoryId}
        >
          <option value="">Select a category...</option>
          {#each categories as category}
            <option value={category.Id}>{category.Name}</option>
          {/each}
        </select>
        {#if formErrors.CategoryId}
          <p class="text-red-500 text-xs mt-1">{formErrors.CategoryId}</p>
        {/if}
      </div>
      
      <div>
        <label for="app-name" class="block text-sm font-medium text-gray-700 mb-1">App Name</label>
        <input 
          id="app-name"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.AppName ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={newRule.AppName}
        />
        {#if formErrors.AppName}
          <p class="text-red-500 text-xs mt-1">{formErrors.AppName}</p>
        {/if}
      </div>
      
      <div>
        <label for="additional-data-key" class="block text-sm font-medium text-gray-700 mb-1">Additional Data Key</label>
        <input 
          id="additional-data-key"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={newRule.AdditionalDataKey}
        />
      </div>
      
      <div>
        <label for="expression" class="block text-sm font-medium text-gray-700 mb-1">Expression</label>
        <input 
          id="expression"
          type="text" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 {formErrors.Expression ? 'border-red-500' : 'border-gray-300'}" 
          bind:value={newRule.Expression}
        />
        {#if formErrors.Expression}
          <p class="text-red-500 text-xs mt-1">{formErrors.Expression}</p>
        {/if}
      </div>
      
      <div>
        <label for="priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
        <input 
          id="priority"
          type="number" 
          class="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500 border-gray-300" 
          bind:value={newRule.Priority}
        />
      </div>
      
      <div class="flex items-center">
        <label class="inline-flex items-center mt-4">
          <input type="checkbox" class="form-checkbox h-5 w-5 text-blue-600" bind:checked={newRule.IsRegex} />
          <span class="ml-2 text-sm text-gray-700">Is Regex</span>
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