<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Column } from '../../types/table';

  export let data = [];
  export let columns: Column[] = [];
  export let pageSize = 10;
  export let actionIcon = null; // For backward compatibility
  export let rowActions = []; // New way to define multiple actions
  export let emptyMessage = "No data found";
  export let noPagination = false;
  export let noHeader = false;
  
  const dispatch = createEventDispatcher();
  
  let currentPage = 1;
  let sortKey = null;
  let sortDirection = 'asc';
  
  $: totalPages = Math.ceil(data.length / pageSize);
  $: startIndex = (currentPage - 1) * pageSize;
  $: endIndex = Math.min(startIndex + pageSize, data.length);
  $: paginatedData = noPagination ? data : data.slice(startIndex, endIndex);
  
  $: sortedData = sortKey !== null 
    ? [...paginatedData].sort((a, b) => {
        if (a[sortKey] < b[sortKey]) return sortDirection === 'asc' ? -1 : 1;
        if (a[sortKey] > b[sortKey]) return sortDirection === 'asc' ? 1 : -1;
        return 0;
      })
    : paginatedData;
  
  function nextPage() {
    if (currentPage < totalPages) currentPage++;
  }
  
  function prevPage() {
    if (currentPage > 1) currentPage--;
  }
  
  function sort(column) {
    if (!column.sortable) return;
    
    if (sortKey === column.key) {
      if (sortDirection === 'asc') {
        sortDirection = 'desc';
      } else {
        // Reset sorting when clicked a third time
        sortKey = null;
        sortDirection = 'asc';
      }
    } else {
      sortKey = column.key;
      sortDirection = 'asc';
    }
  }
  
  function handleActionClick(row, actionHandler) {
    if (actionHandler) {
      actionHandler(row);
    } else {
      // For backward compatibility
      dispatch('rowAction', { row });
    }
  }
</script>

<div class="overflow-hidden">
  <div class="overflow-x-auto">
    <table class="min-w-full divide-y data-table">
      {#if !noHeader}
      <thead class="table-header">
        <tr>
          {#each columns as column}
            <th 
              class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider cursor-pointer"
              class:cursor-default={!column.sortable}
              on:click={() => sort(column)}
            >
              <div class="flex items-center space-x-1">
                <span>{column.title}</span>
                {#if column.sortable && sortKey === column.key}
                  <span class="ml-1">
                    {#if sortDirection === 'asc'}
                      ▲
                    {:else}
                      ▼
                    {/if}
                  </span>
                {/if}
              </div>
            </th>
          {/each}
          {#if rowActions.length > 0 || actionIcon}
            <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider">
              Actions
            </th>
          {/if}
        </tr>
      </thead>
      {/if}
      
      <tbody class="table-body divide-y">
        {#if sortedData.length === 0}
          <tr>
            <td colspan={columns.length + ((rowActions.length > 0 || actionIcon) ? 1 : 0)} class="px-6 py-4 text-center empty-message">
              {emptyMessage}
            </td>
          </tr>
        {:else}
          {#each sortedData as row}
            <tr class="hover-row">
              {#each columns as column}
                <td class="px-6 py-4 whitespace-nowrap">
                  {#if column.formatter}
                    {@html column.formatter(row[column.key])}
                  {:else}
                    {row[column.key]}
                  {/if}
                </td>
              {/each}
              {#if rowActions.length > 0 || actionIcon}
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div class="flex justify-end space-x-2">
                    {#if rowActions.length > 0}
                      {#each rowActions as action}
                        <button 
                          class={`focus:outline-none cursor-pointer action-btn ${action.icon === 'trash' ? 'delete-btn' : 'edit-btn'}`}
                          on:click={() => handleActionClick(row, action.handler)}
                          title={action.title || ''}
                        >
                          {#if action.icon === 'trash'}
                            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                            </svg>
                          {:else if action.icon === 'edit'}
                            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                            </svg>
                          {:else}
                            {action.icon}
                          {/if}
                        </button>
                      {/each}
                    {:else if actionIcon}
                      <!-- Backward compatibility -->
                      <button 
                        class="focus:outline-none cursor-pointer delete-btn"
                        on:click={() => handleActionClick(row, null)}
                      >
                        {#if actionIcon === 'trash'}
                          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                          </svg>
                        {:else if actionIcon === 'edit'}
                          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                          </svg>
                        {:else}
                          {actionIcon}
                        {/if}
                      </button>
                    {/if}
                  </div>
                </td>
              {/if}
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  
  {#if !noPagination && totalPages > 1}
    <div class="pagination">
      <div class="flex-1 flex justify-between items-center">
        <button
          class="pagination-btn"
          disabled={currentPage === 1}
          on:click={prevPage}
        >
          Previous
        </button>
        <span class="pagination-info">
          Page {currentPage} of {totalPages}
        </span>
        <button
          class="pagination-btn"
          disabled={currentPage === totalPages}
          on:click={nextPage}
        >
          Next
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .data-table {
    color: var(--text-color);
    border-color: var(--table-border-color);
  }
  
  .table-header {
    background-color: var(--table-header-bg);
  }
  
  .table-body {
    background-color: var(--card-bg-color);
    color: var(--text-color);
    border-color: var(--table-border-color);
  }
  
  .hover-row:hover {
    background-color: var(--table-row-hover);
  }
  
  .empty-message {
    color: var(--secondary-color);
  }
  
  .pagination {
    padding: 0.75rem 1.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-top: 1px solid var(--table-border-color);
    background-color: var(--card-bg-color);
  }
  
  .pagination-btn {
    padding: 0.5rem 1rem;
    border: 1px solid var(--table-border-color);
    border-radius: 0.25rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-color);
    background-color: var(--button-bg-color);
  }
  
  .pagination-btn:hover:not([disabled]) {
    background-color: var(--button-hover-bg-color);
  }
  
  .pagination-btn[disabled] {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .pagination-info {
    font-size: 0.875rem;
    color: var(--text-color);
  }
  
  .action-btn {
    color: var(--primary-color);
  }
  
  .edit-btn {
    color: var(--primary-color);
  }
  
  .edit-btn:hover {
    color: var(--info-color);
  }
  
  .delete-btn {
    color: var(--danger-color);
  }
  
  .delete-btn:hover {
    opacity: 0.8;
  }

  .divide-y > :not([hidden]) ~ :not([hidden]) {
    border-top-width: 1px;
    border-top-style: solid;
    border-top-color: var(--table-border-color);
  }
</style>
