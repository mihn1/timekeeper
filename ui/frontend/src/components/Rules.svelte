<script>
  import { onMount } from 'svelte';
  import { GetRules, GetRule, AddRule, DeleteRule } from '../../wailsjs/go/main/App';
  import { refreshData } from '../stores/timekeeper';

  let rules = [];
  let newRule = { RuleId: 0, CategoryId: '', AppName: '', AdditionalDataKey: '', Expression: '', IsRegex: false, Priority: 0 };
  let isLoading = true;

  $: if ($refreshData) {
    loadRules();
  }

  onMount(() => {
    loadRules();
  });

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

  async function addRule() {
    try {
      await AddRule(newRule);
      newRule = { RuleId: 0, CategoryId: '', AppName: '', AdditionalDataKey: '', Expression: '', IsRegex: false, Priority: 0 };
      loadRules();
    } catch (err) {
      console.error('Error adding rule:', err);
    }
  }

  async function deleteRule(ruleId) {
    try {
      await DeleteRule(ruleId);
      loadRules();
    } catch (err) {
      console.error('Error deleting rule:', err);
    }
  }
</script>

<div class="rule-management">
  <h2>Rule Management</h2>
  <div class="rule-form">
    <input type="text" placeholder="Category ID" bind:value={newRule.CategoryId} />
    <input type="text" placeholder="App Name" bind:value={newRule.AppName} />
    <input type="text" placeholder="Additional Data Key" bind:value={newRule.AdditionalDataKey} />
    <input type="text" placeholder="Expression" bind:value={newRule.Expression} />
    <input type="checkbox" bind:checked={newRule.IsRegex} /> Is Regex
    <input type="number" placeholder="Priority" bind:value={newRule.Priority} />
    <button on:click={addRule}>Add Rule</button>
  </div>
  {#if isLoading}
    <div class="loading">Loading rules...</div>
  {:else}
    <ul>
      {#each rules as rule}
        <li>
          {rule.AppName} - {rule.Expression}
          <button on:click={() => deleteRule(rule.RuleId)}>Delete</button>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .rule-management {
    padding: 1rem;
  }

  .rule-form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .loading {
    color: #666;
  }
</style>