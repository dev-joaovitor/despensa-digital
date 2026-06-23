<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import ShoppingListItemRow from '$lib/components/ShoppingListItemRow.svelte';
	import AddShoppingItemModal from '$lib/components/AddShoppingItemModal.svelte';
	import SubmitShoppingListModal from '$lib/components/SubmitShoppingListModal.svelte';

	let { data } = $props();

	let products = $state(data.products);
	let establishments = $state(data.establishments);
	let brands = $state(data.brands);
	let categories = $state(data.categories);

	let addOpen = $state(false);
	let submitOpen = $state(false);

	$effect(() => {
		products = data.products;
	});
	$effect(() => {
		establishments = data.establishments;
	});
	$effect(() => {
		brands = data.brands;
	});
	$effect(() => {
		categories = data.categories;
	});

	const checkedItems = $derived(data.items.filter((i) => i.is_checked));
</script>

<svelte:head>
	<title>Lista de compras — IntelliStock</title>
</svelte:head>

<div class="page">
	<header class="head">
		<h1>Lista de compras</h1>
		<PrimaryButton onclick={() => (addOpen = true)}>Adicionar item</PrimaryButton>
	</header>

	{#if data.items.length === 0}
		<p class="empty">Nenhum item na lista.</p>
	{:else}
		<div class="list">
			{#each data.items as item (item.id)}
				<ShoppingListItemRow {item} onchange={invalidateAll} />
			{/each}
		</div>
	{/if}

	<footer class="foot">
		<PrimaryButton onclick={() => (submitOpen = true)} disabled={checkedItems.length === 0}>
			Enviar lista
		</PrimaryButton>
	</footer>
</div>

<AddShoppingItemModal
	bind:open={addOpen}
	bind:products
	bind:brands
	bind:categories
	measurements={data.measurements}
	onsuccess={invalidateAll}
/>

<SubmitShoppingListModal
	bind:open={submitOpen}
	items={checkedItems}
	bind:establishments
	onsuccess={invalidateAll}
/>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		width: 100%;
		max-width: 52rem;
		max-height: calc(100vh - 8rem);
	}

	.head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-md);
	}

	.head :global(.btn) {
		white-space: nowrap;
	}

	h1 {
		font-size: 1.5rem;
		color: var(--color-text);
	}

	.list {
		flex: 1;
		min-height: 0;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding-right: var(--space-xs);
	}

	.foot {
		display: flex;
	}

	.foot :global(.btn) {
		flex: 1;
	}

	.empty {
		flex: 1;
		color: var(--color-text-muted);
		text-align: center;
	}
</style>
