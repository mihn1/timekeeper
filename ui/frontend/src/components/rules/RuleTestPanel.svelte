<script lang="ts">
  import { TestRuleMatch } from '../../../wailsjs/go/main/App';
  import type { dtos } from '../../../wailsjs/go/models';

  let appName = '';
  let additionalDataKey = '';
  let value = '';
  let result: dtos.RuleMatchResult | null = null;
  let isLoading = false;
  let error = '';

  async function runTest() {
    if (!appName.trim()) return;
    isLoading = true;
    error = '';
    result = null;
    try {
      result = await TestRuleMatch(appName.trim(), additionalDataKey.trim(), value.trim());
    } catch (err) {
      error = String(err);
    } finally {
      isLoading = false;
    }
  }

  function reset() {
    appName = '';
    additionalDataKey = '';
    value = '';
    result = null;
    error = '';
  }
</script>

<div class="test-panel">
  <h3>Rule Tester</h3>
  <p class="hint">Simulate how an app event resolves to a category given the current rules.</p>

  <div class="form-row">
    <div class="field">
      <label>App Name</label>
      <input
        type="text"
        placeholder="e.g. Google Chrome"
        bind:value={appName}
        on:keydown={(e) => e.key === 'Enter' && runTest()}
      />
    </div>
    <div class="field">
      <label>Data Key <span class="optional">(optional)</span></label>
      <input
        type="text"
        placeholder="url / title"
        bind:value={additionalDataKey}
        on:keydown={(e) => e.key === 'Enter' && runTest()}
      />
    </div>
    <div class="field flex-2">
      <label for="test-value">Value <span class="optional">(optional)</span></label>
      <input
        id="test-value"
        type="text"
        placeholder="https://github.com/…"
        bind:value={value}
        on:keydown={(e) => e.key === 'Enter' && runTest()}
      />
    </div>
    <div class="actions">
      <button class="test-btn" on:click={runTest} disabled={!appName.trim() || isLoading}>
        {isLoading ? 'Testing…' : 'Test'}
      </button>
      {#if result || error}
        <button class="reset-btn" on:click={reset}>Clear</button>
      {/if}
    </div>
  </div>

  {#if error}
    <div class="result error">{error}</div>
  {:else if result}
    <div class="result {result.matched ? 'match' : 'no-match'}">
      {#if result.matched}
        <span class="badge match-badge">Matched</span>
        Category: <strong>{result.categoryName}</strong>
        {#if result.matchedRule}
          · Rule #{result.matchedRule.id}
          ({result.matchedRule.appName}{result.matchedRule.expression ? ' / ' + result.matchedRule.expression : ''})
        {/if}
      {:else}
        <span class="badge no-badge">No match</span>
        Resolves to: <strong>{result.categoryName}</strong> (default)
      {/if}
    </div>
  {/if}
</div>

<style>
  .test-panel {
    background-color: var(--card-bg-color);
    border: 1px solid var(--card-border-color);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    box-shadow: var(--card-shadow);
    margin-bottom: 1.25rem;
  }

  h3 { margin: 0 0 0.2rem; font-size: 1rem; font-weight: 600; color: var(--text-color); }
  .hint { margin: 0 0 0.75rem; font-size: 0.8rem; color: var(--secondary-color); }

  .form-row {
    display: flex;
    align-items: flex-end;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    flex: 1;
    min-width: 140px;
  }

  .field.flex-2 { flex: 2; }

  label { font-size: 0.75rem; font-weight: 500; color: var(--text-color); }
  .optional { font-weight: 400; color: var(--secondary-color); }

  input {
    padding: 0.4rem 0.6rem;
    border: 1px solid var(--input-border-color);
    border-radius: 4px;
    background-color: var(--input-bg-color);
    color: var(--input-text-color);
    font-size: 0.85rem;
  }

  .actions { display: flex; align-items: flex-end; gap: 0.4rem; }

  .test-btn {
    padding: 0.4rem 1rem;
    background-color: var(--button-bg-color);
    color: var(--button-text-color);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
    white-space: nowrap;
  }

  .test-btn:disabled { opacity: 0.4; cursor: default; }
  .test-btn:not(:disabled):hover { background-color: var(--button-hover-bg-color); }

  .reset-btn {
    padding: 0.4rem 0.75rem;
    background: none;
    border: 1px solid var(--card-border-color);
    border-radius: 4px;
    cursor: pointer;
    color: var(--secondary-color);
    font-size: 0.85rem;
  }

  .result {
    margin-top: 0.75rem;
    padding: 0.5rem 0.75rem;
    border-radius: 6px;
    font-size: 0.85rem;
    color: var(--text-color);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .result.match    { background: rgba(34,197,94,0.1); border: 1px solid #22c55e44; }
  .result.no-match { background: rgba(239,68,68,0.08); border: 1px solid #ef444430; }
  .result.error    { background: rgba(239,68,68,0.08); border: 1px solid #ef444430; }

  .badge {
    display: inline-block;
    padding: 0.1rem 0.5rem;
    border-radius: 99px;
    font-size: 0.72rem;
    font-weight: 700;
  }

  .match-badge { background: #22c55e; color: #fff; }
  .no-badge    { background: #ef4444; color: #fff; }
</style>
