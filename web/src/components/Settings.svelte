<script lang="ts">
  import * as api from "../lib/api";

  interface Props {
    onClose: () => void;
  }

  let { onClose }: Props = $props();

  let apiKey = $state<string | null>(null);
  let revealed = $state(false);
  let copied = $state(false);
  let regenerating = $state(false);
  let error = $state<string | null>(null);

  async function load() {
    error = null;
    try {
      const settings = await api.getSettings();
      apiKey = settings.apiKey;
    } catch (e) {
      error = (e as Error).message;
    }
  }

  async function copyKey() {
    if (!apiKey) return;
    await navigator.clipboard.writeText(apiKey);
    copied = true;
    setTimeout(() => (copied = false), 1500);
  }

  async function regenerate() {
    if (!confirm("Regenerate the API key? Anything using the current key will stop working.")) {
      return;
    }
    regenerating = true;
    error = null;
    try {
      const settings = await api.regenerateApiKey();
      apiKey = settings.apiKey;
      revealed = true;
    } catch (e) {
      error = (e as Error).message;
    } finally {
      regenerating = false;
    }
  }

  load();
</script>

<div class="settings">
  <div class="settings-header">
    <h2>Settings</h2>
    <button onclick={onClose} aria-label="Close settings">Close</button>
  </div>

  {#if error}
    <p class="error">{error}</p>
  {/if}

  <section>
    <h3>API key</h3>
    <p class="hint">
      Send this as an <code>X-Api-Key</code> header to call Palimpsest's API from
      other apps. It's not access control - anyone who can load this page can
      see it, so put an authenticating proxy in front if you need real privacy.
    </p>

    {#if apiKey}
      <div class="key-row">
        <code class="key-value">{revealed ? apiKey : "•".repeat(24)}</code>
        <button onclick={() => (revealed = !revealed)}>{revealed ? "Hide" : "Show"}</button>
        <button onclick={copyKey}>{copied ? "Copied!" : "Copy"}</button>
      </div>
      <button class="danger" onclick={regenerate} disabled={regenerating}>
        {regenerating ? "Regenerating..." : "Regenerate key"}
      </button>
    {:else if !error}
      <p class="hint">Loading...</p>
    {/if}
  </section>
</div>

<style>
  .settings {
    max-width: 40rem;
  }

  .settings-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .settings-header h2 {
    margin: 0;
  }

  section {
    margin-top: 1.5rem;
  }

  h3 {
    margin-bottom: 0.4rem;
  }

  .hint {
    color: var(--muted);
    font-size: 0.85rem;
    line-height: 1.5;
    max-width: 32rem;
  }

  .hint code {
    font-family: var(--font-mono);
  }

  .key-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 0.75rem 0;
  }

  .key-value {
    flex: 1;
    min-width: 0;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 0.5rem 0.7rem;
    font-family: var(--font-mono);
    font-size: 0.85rem;
    overflow-x: auto;
    white-space: nowrap;
  }

  .error {
    color: var(--danger);
  }

  button {
    background: var(--surface);
    border: 1px solid var(--border);
    color: inherit;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    transition:
      background-color 120ms ease,
      border-color 120ms ease;
  }

  button:hover:not(:disabled) {
    border-color: var(--accent);
  }

  button:disabled {
    opacity: 0.5;
    cursor: default;
  }

  button.danger {
    background: var(--danger-tint);
    border-color: transparent;
    color: var(--danger);
    margin-top: 0.5rem;
  }

  button.danger:hover:not(:disabled) {
    border-color: var(--danger);
  }
</style>
