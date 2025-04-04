<script lang="ts">
  import { onMount } from 'svelte';
  import { GetCategories, DeleteCategory } from '../../../wailsjs/go/main/App';
  import { refreshData } from '../../stores/timekeeper';
  import Modal from '../common/Modal.svelte';
  import DataTable from '../common/DataTable.svelte';
  import CreateCategoryModal from './CreateCategoryModal.svelte';
  import EditCategoryModal from './EditCategoryModal.svelte';
  import { dtos } from '../../../wailsjs/go/models';
  import type { Column } from '../../types/table'; // Import the shared type

  let categories: dtos.CategoryListItem[] = [];
  let isLoading = true;
  let showDeleteModal = false;
  let showCreateCategoryModal = false;
  let showEditCategoryModal = false;
  let categoryToDelete = null;
  let categoryToEdit = null;
  let searchTerm = '';
  let pageSizes = [5, 10, 25, 50];
  let selectedPageSize = 10;

  $: if ($refreshData) {
    loadCategories();
  }

  $: filteredCategories = categories.filter(category => 
    category.name?.toLowerCase().includes(searchTerm.toLowerCase()) || 
    category.description?.toLowerCase().includes(searchTerm.toLowerCase())
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

  function confirmDelete(category: dtos.CategoryListItem) {
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

  function editCategory(category) {
    console.log('Editing category:', category);
    categoryToEdit = { ...category }; // Create a copy to avoid reference issues
    showEditCategoryModal = true;
  }

  function handleCategoryEdited() {
    loadCategories();
    showEditCategoryModal = false;
    categoryToEdit = null;
  }

  const tableColumns: Column[] = [
    { key: 'id', title: 'ID', sortable: true },
    { key: 'name', title: 'Name', sortable: true },
    { key: 'description', title: 'Description', sortable: true }
  ];

  const rowActions = [
    { 
      icon: 'edit', 
      handler: editCategory,
      title: 'Edit category'
    },
    {
      icon: 'trash', 
      handler: confirmDelete,
      title: 'Delete category'
    }
  ];
</script>

<!-- Modals -->
<CreateCategoryModal 
  show={showCreateCategoryModal}
  on:close={() => showCreateCategoryModal = false}
  on:categoryAdded={handleCategoryAdded}
/>

<EditCategoryModal
  show={showEditCategoryModal}
  category={categoryToEdit}
  on:close={() => { showEditCategoryModal = false; categoryToEdit = null; }}
  on:categoryEdited={handleCategoryEdited}
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
  <h1 class="text-2xl font-bold mb-6 category-title">Category Management</h1>
  
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
  <div class="category-container rounded-lg shadow-md overflow-hidden">
    <div class="p-5 border-b table-header-wrapper">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-semibold table-header-title">Categories</h2>
        
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
        rowActions={rowActions}
        emptyMessage="No categories found"
        pageSize={selectedPageSize}
      />
    {/if}
  </div>
</div>

<style>
  .category-title {
    color: var(--text-color);
  }

  .category-container {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
  }

  .table-header-wrapper {
    background-color: var(--table-header-bg);
    border-color: var(--table-border-color);
  }

  .table-header-title {
    color: var(--text-color);
  }
</style>