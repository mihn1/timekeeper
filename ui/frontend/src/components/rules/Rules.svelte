<script lang="ts">
  import { onMount } from 'svelte';
  import { GetRules, DeleteRule, GetCategories } from '../../../wailsjs/go/main/App';
  import { refreshData } from '../../stores/timekeeper';
  import Modal from '../common/Modal.svelte';
  import DataTable from '../common/DataTable.svelte';
  import CreateRuleModal from './CreateRuleModal.svelte';
  import CreateCategoryModal from '../categories/CreateCategoryModal.svelte';
  import { dtos } from '../../../wailsjs/go/models';
  import type { Column } from '../../types/table';

  let rules: dtos.RuleListItem[] = [];
  let categories: dtos.CategoryListItem[] = [];
  let isLoading = true;
  let isLoadingCategories = true;
  let showDeleteModal = false;
  let showCreateRuleModal = false;
  let showCreateCategoryModal = false;
  let ruleToDelete: dtos.RuleListItem | null = null;
  let searchTerm = '';
  let pageSizes = [5, 10, 25, 50];
  let selectedPageSize = 10;
  let isGroupedView = true;

  // Reactive statement to force update when categories load
  $: categoriesLoaded = !isLoadingCategories && categories.length > 0;
  
  // Force component update when categories finish loading
  $: if (categoriesLoaded && rules.length > 0) {
    // Force a refresh of the component
    rules = [...rules];
  }

  $: if ($refreshData) {
    loadRules();
  }

  $: filteredRules = rules
    .filter(rule => 
      rule.appName.toLowerCase().includes(searchTerm.toLowerCase()) || 
      rule.expression.toLowerCase().includes(searchTerm.toLowerCase()) ||
      getCategoryName(rule.categoryId).toLowerCase().includes(searchTerm.toLowerCase())
    )
    .sort((a, b) => b.priority - a.priority); // Sort by priority (highest first)

  $: groupedRules = groupRulesByAppName(filteredRules);

  onMount(() => {
    // Load both in parallel for better performance
    loadCategories(),
    loadRules()
  });

  function getCategoryName(categoryId: number): string {
    if (isLoadingCategories) return "Loading...";
    
    // Simple find operation - perfectly fine for small arrays
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : `Unknown (${categoryId})`;
  }

  async function loadRules() {
    isLoading = true;
    try {
      rules = await GetRules();
    } catch (err) {
      console.error('Error loading rules:', err);
    } finally {
      isLoading = false;
    }
  }

  async function loadCategories() {
    isLoadingCategories = true;
    try {
      categories = await GetCategories();
    } catch (err) {
      console.error('Error loading categories:', err);
    } finally {
      isLoadingCategories = false;
    }
  }

  function confirmDelete(rule: dtos.RuleListItem) {
    ruleToDelete = rule;
    showDeleteModal = true;
  }

  async function executeDelete() {
    if (!ruleToDelete) return;
    
    try {
      await DeleteRule(ruleToDelete.id);
      showDeleteModal = false;
      ruleToDelete = null;
      loadRules();
    } catch (err) {
      console.error('Error deleting rule:', err);
    }
  }

  function cancelDelete() {
    showDeleteModal = false;
    ruleToDelete = null;
  }
  
  function openCreateRuleModal() {
    showCreateRuleModal = true;
  }
  
  function handleRuleAdded() {
    loadRules();
    showCreateRuleModal = false;
  }
  
  function handleCategoryAdded() {
    loadCategories();
    showCreateCategoryModal = false;
  }
  
  function handlePageSizeChange(event) {
    selectedPageSize = parseInt(event.target.value);
  }

  function groupRulesByAppName(rules: dtos.RuleListItem[]) {
    const grouped = {};
    
    rules.forEach(rule => {
      if (!grouped[rule.appName]) {
        grouped[rule.appName] = [];
      }
      grouped[rule.appName].push(rule);
    });
    
    // Convert to array format for rendering
    return Object.keys(grouped).map(appName => ({
      appName,
      rules: grouped[appName],
      _isGroupHeader: true
    }));
  }
  
  function toggleGroupView() {
    isGroupedView = !isGroupedView;
  }

  const tableColumns: Column[] = [
    { 
      key: 'categoryId', 
      title: 'Category', 
      sortable: true, 
      formatter: (value) => getCategoryName(value)
    },
    { key: 'appName', title: 'App Name', sortable: true },
    { key: 'additionalDataKey', title: 'Data Key', sortable: true },
    { key: 'expression', title: 'Expression', sortable: true },
    { key: 'isRegex', title: 'Regex', sortable: true, formatter: (value) => value ? 'Yes' : 'No' },
    { key: 'priority', title: 'Priority', sortable: true },
  ];
</script>

<!-- Modals -->
<CreateRuleModal 
  show={showCreateRuleModal} 
  categories={categories}
  on:close={() => showCreateRuleModal = false}
  on:ruleAdded={handleRuleAdded}
/>

<CreateCategoryModal 
  show={showCreateCategoryModal}
  on:close={() => showCreateCategoryModal = false}
  on:categoryAdded={handleCategoryAdded}
/>

<!-- Delete Confirmation Modal -->
<Modal 
  show={showDeleteModal} 
  title="Confirm Delete"
  on:close={() => showDeleteModal = false}
>
  <div class="p-4">
    <p class="mb-4">Are you sure you want to delete this rule for <strong>{ruleToDelete?.appName}</strong>?</p>
    <p class="mb-4 text-red-600">This action cannot be undone.</p>
    
    <div class="flex justify-end gap-2">
      <button class="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400 cursor-pointer" on:click={cancelDelete}>
        Cancel
      </button>
      <button class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 cursor-pointer" on:click={executeDelete}>
        Delete
      </button>
    </div>
  </div>
</Modal>

<div class="p-6 max-w-6xl mx-auto">
  <h1 class="text-2xl font-bold mb-6 text-gray-800">Rule Management</h1>
  
  <!-- Action Buttons -->
  <div class="flex justify-between mb-6">
    <div class="flex gap-4">
      <button 
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 flex items-center cursor-pointer"
        on:click={openCreateRuleModal}
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
        </svg>
        Add Rule
      </button>
    </div>
    
    <div class="flex items-center">
      <label class="inline-flex items-center cursor-pointer">
        <input type="checkbox" class="form-checkbox h-5 w-5 text-blue-600" bind:checked={isGroupedView}>
        <span class="ml-2 text-sm text-gray-700">Group by App Name</span>
      </label>
    </div>
  </div>
  
  <!-- Rule List Section -->
  <div class="bg-white rounded-lg shadow-md overflow-hidden">
    <div class="p-5 border-b border-gray-200">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-semibold text-gray-700">Rules</h2>
        
        <div class="flex gap-4 items-center">
          <div class="flex items-center">
            <label for="pageSize" class="mr-2 text-sm text-gray-700">Items per page:</label>
            <select 
              id="pageSize"
              class="p-1 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              bind:value={selectedPageSize}
              on:change={handlePageSizeChange}
            >
              {#each pageSizes as size}
                <option value={size}>{size}</option>
              {/each}
            </select>
          </div>
          
          <div class="relative">
            <input 
              type="text" 
              placeholder="Search rules..." 
              class="pl-10 p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              bind:value={searchTerm}
            />
            <div class="absolute left-3 top-2.5">
              <svg class="h-4 w-4 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    {#if isLoading}
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    {:else if isGroupedView}
      {#each groupedRules as group}
        <div class="border-b border-gray-200 last:border-b-0">
          <div class="bg-gray-100 p-3 font-medium">
            {group.appName} ({group.rules.length} {group.rules.length === 1 ? 'rule' : 'rules'})
          </div>
          <DataTable 
            data={group.rules} 
            columns={tableColumns}
            on:rowAction={(e) => confirmDelete(e.detail.row)}
            actionIcon="trash"
            emptyMessage="No rules found"
            pageSize={group.rules.length}
            noPagination={true}
            noHeader={false} 
            />
        </div>
      {:else}
        <div class="p-8 text-center text-gray-500">
          No rules found
        </div>
      {/each}
    {:else}
      <DataTable 
        data={filteredRules} 
        columns={tableColumns}
        on:rowAction={(e) => confirmDelete(e.detail.row)}
        actionIcon="trash"
        emptyMessage="No rules found"
        pageSize={selectedPageSize}
      />
    {/if}
  </div>
</div>