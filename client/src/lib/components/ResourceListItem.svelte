<script lang="ts">
	import type { NamedResource } from '$lib/resources';
	import PrimaryButton from './PrimaryButton.svelte';
	import SecondaryButton from './SecondaryButton.svelte';

	interface Props {
		item: NamedResource;
		onedit: (item: NamedResource) => void;
		onremove: (item: NamedResource) => void;
	}

	let { item, onedit, onremove }: Props = $props();
</script>

<li class="item">
	<span class="name">{item.name}</span>
	<div class="actions">
		<PrimaryButton onclick={() => onedit(item)}>Editar</PrimaryButton>
		<SecondaryButton onclick={() => onremove(item)}>Remover</SecondaryButton>
	</div>
</li>

<style>
	.item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-md);
		padding: var(--space-md) var(--space-lg);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
	}

	.name {
		color: var(--color-text);
		overflow-wrap: anywhere;
	}

	.actions {
		display: flex;
		gap: var(--space-sm);
		flex-shrink: 0;
	}
</style>
