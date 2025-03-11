<script>
  import { onMount } from 'svelte';
  import { GetCategories, GetCategory, AddCategory, DeleteCategory } from '../../wailsjs/go/main/App';
  import { refreshData } from '../stores/timekeeper';

  let categories = [];
  let newCategory = { Id: '', Name: '', Description: '', CategoryTypeId: 0 };
  let isLoading = true;

  $: if ($refreshData) {
    loadCategories();
  }

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

  async function addCategory() {
    try {
      await AddCategory(newCategory);
      newCategory = { Id: '', Name: '', Description: '', CategoryTypeId: 0 };
      loadCategories();
    } catch (err) {
      console.error('Error adding category:', err);
    }
  }

  async function deleteCategory(categoryId) {
    try {
      await DeleteCategory(categoryId);
      loadCategories();
    } catch (err) {
      console.error('Error deleting category:', err);
    }
  }
</script>

<div class="category-management">
  <h2>Category Management</h2>
  <div class="category-form">
    <input type="text" placeholder="ID" bind:value={newCategory.Id} />
    <input type="text" placeholder="Name" bind:value={newCategory.Name} />
    <input type="text" placeholder="Description" bind:value={newCategory.Description} />
    <input type="number" placeholder="Category Type ID" bind:value={newCategory.CategoryTypeId} />
    <button on:click={addCategory}>Add Category</button>
  </div>
  {#if isLoading}
    <div class="loading">Loading categories...</div>
  {:else}
    <ul>
      {#each categories as category}
        <li>
          {category.Name} - {category.Description}
          <button on:click={() => deleteCategory(category.Id)}>Delete</button>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .category-management {
    padding: 1rem;
  }

  .category-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .loading {
    color: #666;
  }
</style>