<script lang="ts">
	import type { NamedResource } from '$lib/resources';
	import ResourceListItem from './ResourceListItem.svelte';

	interface Props {
		items: NamedResource[];
		emptyText?: string;
		onedit: (item: NamedResource) => void;
		onremove: (item: NamedResource) => void;
	}

	let { items, emptyText = 'Nenhum item cadastrado.', onedit, onremove }: Props = $props();
</script>

{#if items.length === 0}
	<p class="empty">{emptyText}</p>
{:else}
	<ul class="list">
		{#each items as item (item.id)}
			<ResourceListItem {item} {onedit} {onremove} />
		{/each}
	</ul>
{/if}

<style>
	.list {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.empty {
		color: var(--color-text-muted);
	}
</style>
