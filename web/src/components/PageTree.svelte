<script lang="ts">
  import type { PageSummary } from "../types/Page";
  import PageTree from "./PageTree.svelte";

  interface Props {
    pages: PageSummary[];
    parentId: string | null;
    selectedId: string | null;
    onSelect: (id: string) => void;
  }

  let { pages, parentId, selectedId, onSelect }: Props = $props();

  let children = $derived(pages.filter((p) => p.parentId === parentId));
  let isNested = $derived(parentId !== null);
</script>

<ul class="tree" class:nested={isNested}>
  {#each children as page (page.id)}
    <li>
      <button
        class="node"
        class:selected={page.id === selectedId}
        onclick={() => onSelect(page.id)}
      >
        {page.title}
      </button>
      {#if pages.some((p) => p.parentId === page.id)}
        <PageTree {pages} parentId={page.id} {selectedId} {onSelect} />
      {/if}
    </li>
  {/each}
</ul>

<style>
  .tree {
    list-style: none;
    margin: 0;
    padding-left: 0;
  }

  .tree.nested {
    padding-left: 0.9rem;
    border-left: 1px solid var(--border);
    margin-left: 0.6rem;
  }

  .node {
    display: block;
    width: 100%;
    text-align: left;
    background: none;
    border: none;
    border-left: 2px solid transparent;
    color: var(--muted);
    padding: 0.3rem 0.5rem;
    border-radius: 0 4px 4px 0;
    cursor: pointer;
    font-size: 0.9rem;
  }

  .node:hover {
    background: var(--surface);
    color: var(--text);
  }

  .node.selected {
    background: var(--accent-tint);
    border-left-color: var(--accent);
    color: var(--text);
  }
</style>
