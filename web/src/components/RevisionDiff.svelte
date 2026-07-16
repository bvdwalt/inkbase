<script lang="ts">
  import { diffWords } from "diff";

  interface Props {
    oldText: string;
    newText: string;
  }

  let { oldText, newText }: Props = $props();

  let parts = $derived(diffWords(oldText, newText));
</script>

<div class="diff">
  {#if oldText === newText}
    <p class="no-change">No text changes.</p>
  {:else}
    {#each parts as part, i (i)}
      <span class:added={part.added} class:removed={part.removed}>{part.value}</span>
    {/each}
  {/if}
</div>

<style>
  .diff {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 0.75rem;
    margin: 0.5rem 0;
    white-space: pre-wrap;
    font-family: var(--font-mono);
    font-size: 0.82rem;
    line-height: 1.5;
  }

  .added {
    background: rgba(143, 209, 158, 0.14);
    color: #8fd19e;
  }

  .removed {
    background: rgba(209, 158, 143, 0.14);
    color: #d19e8f;
    text-decoration: line-through;
  }

  .no-change {
    color: var(--muted);
    margin: 0;
    font-family: var(--font-sans);
  }
</style>
