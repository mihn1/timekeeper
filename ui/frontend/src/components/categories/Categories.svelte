<script>
  import { onMount } from 'svelte';
  import { GetCategories, GetCategory, DeleteCategory } from '../../../wailsjs/go/main/App';
  import { refreshData } from '../../stores/timekeeper';
  import Modal from '../common/Modal.svelte';
  import DataTable from '../common/DataTable.svelte';
  import CreateCategoryModal from './CreateCategoryModal.svelte';

  let categories = [];
  let isLoading = true;
  let showDeleteModal = false;
  let showCreateCategoryModal = false;
  let categoryToDelete = null;
  let searchTerm = '';
  let pageSizes = [5, 10, 25, 50];
  let selectedPageSize = 10;

  $: if ($refreshData) {
    loadCategories();
  }

  $: filteredCategories = categories.filter(category => 
    category.Name?.toLowerCase().includes(searchTerm.toLowerCase()) || 
    category.Description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    category.Id?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  onMount(() => {
    loadCategories();
  });

  async function loadCategories() {
    isLoading = true;
    try {
      categories = await GetCategories();
    } catch (err) {
      console.error('Error loading categories:', err);
    } finally {
      isLoading = false;
    }
  }

  function confirmDelete(category) {
    categoryToDelete = category;
    showDeleteModal = true;
  }

  async function executeDelete() {
    if (!categoryToDelete) return;
    
    try {
      await DeleteCategory(categoryToDelete.Id);
      showDeleteModal = false;
      categoryToDelete = null;
      loadCategories();
    } catch (err) {
      console.error('Error deleting category:', err);
    }
  }

  function cancelDelete() {
    showDeleteModal = false;
    categoryToDelete = null;
  }
  
  function openCreateCategoryModal() {
    showCreateCategoryModal = true;
  }
  
  function handleCategoryAdded() {
    loadCategories();
    showCreateCategoryModal = false;
  }
  
  function handlePageSizeChange(event) {
    selectedPageSize = parseInt(event.target.value);
  }

  const tableColumns = [
    { key: 'Id', title: 'ID', sortable: true },
    { key: 'Name', title: 'Name', sortable: true },
    { key: 'Description', title: 'Description', sortable: true },
    { key: 'CategoryTypeId', title: 'Type', sortable: true },
    { key: 'actions', title: 'Actions', sortable: false }
  ];
</script>

<!-- Modals -->
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
    <p class="mb-4">Are you sure you want to delete the category: <strong>{categoryToDelete?.Name}</strong>?</p>
    <p class="mb-4 text-red-600">This action cannot be undone. All rules associated with this category will also be deleted.</p>
    
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
  <h1 class="text-2xl font-bold mb-6 text-gray-800">Category Management</h1>
  
  <!-- Action Buttons -->
  <div class="flex gap-4 mb-6">
    <button 
      class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 flex items-center cursor-pointer"
      on:click={openCreateCategoryModal}
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
      </svg>
      Add Category
    </button>
  </div>
  
  <!-- Category List Section -->
  <div class="bg-white rounded-lg shadow-md overflow-hidden">
    <div class="p-5 border-b border-gray-200">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-semibold text-gray-700">Categories</h2>
        
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
              placeholder="Search categories..." 
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
    {:else}
      <DataTable 
        data={filteredCategories} 
        columns={tableColumns}
        on:rowAction={(e) => confirmDelete(e.detail.row)}
        actionIcon="trash"
        emptyMessage="No categories found"
        pageSize={selectedPageSize}
      />
    {/if}
  </div>
</div>

<svelte:head>
  <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</svelte:head>