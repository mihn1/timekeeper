<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Column } from '../../types/table';

  export let data = [];
  export let columns: Column[] = [];
  export let pageSize = 10;
  export let actionIcon = null;
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
  
  $: sortedData = sortKey 
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
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortKey = column.key;
      sortDirection = 'asc';
    }
  }
  
  function handleActionClick(row) {
    dispatch('rowAction', { row });
  }
</script>

<div class="overflow-hidden">
  <div class="overflow-x-auto">
    <table class="min-w-full divide-y divide-gray-200">
      {#if !noHeader}
      <thead class="bg-gray-50">
        <tr>
          {#each columns as column}
            <th 
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
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
          {#if actionIcon}
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          {/if}
        </tr>
      </thead>
      {/if}
      
      <tbody class="bg-white divide-y divide-gray-200">
        {#if sortedData.length === 0}
          <tr>
            <td colspan={columns.length + (actionIcon ? 1 : 0)} class="px-6 py-4 text-center text-gray-500">
              {emptyMessage}
            </td>
          </tr>
        {:else}
          {#each sortedData as row}
            <tr class="hover:bg-gray-50">
              {#each columns as column}
                <td class="px-6 py-4 whitespace-nowrap">
                  {#if column.formatter}
                    {@html column.formatter(row[column.key])}
                  {:else}
                    {row[column.key]}
                  {/if}
                </td>
              {/each}
              {#if actionIcon}
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button 
                    class="text-red-600 hover:text-red-900 focus:outline-none cursor-pointer"
                    on:click={() => handleActionClick(row)}
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
                </td>
              {/if}
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
  
  {#if !noPagination && totalPages > 1}
    <div class="px-6 py-3 flex items-center justify-between border-t border-gray-200">
      <div class="flex-1 flex justify-between items-center">
        <button
          class="px-4 py-2 border rounded text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={currentPage === 1}
          on:click={prevPage}
        >
          Previous
        </button>
        <span class="text-sm text-gray-700">
          Page {currentPage} of {totalPages}
        </span>
        <button
          class="px-4 py-2 border rounded text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={currentPage === totalPages}
          on:click={nextPage}
        >
          Next
        </button>
      </div>
    </div>
  {/if}
</div>
