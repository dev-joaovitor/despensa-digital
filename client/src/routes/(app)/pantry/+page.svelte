<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import PantryProductCard from '$lib/components/PantryProductCard.svelte';
	import ProductFormModal from '$lib/components/ProductFormModal.svelte';
	import StockTransactionModal from '$lib/components/StockTransactionModal.svelte';
	import AddPriceObservationModal from '$lib/components/AddPriceObservationModal.svelte';
	import type { StockProduct } from '$lib/pantry';
	import type { Product } from '$lib/price-observations';

	let { data } = $props();

	let searchValue = $state(data.search);
	let products = $state(data.products);
	let establishments = $state(data.establishments);
	let brands = $state(data.brands);
	let categories = $state(data.categories);

	let editOpen = $state(false);
	let editProduct = $state<StockProduct | null>(null);
	let priceOpen = $state(false);

	let stockOpen = $state(false);
	let stockProduct = $state<Product | null>(null);

	// Keep state in sync when the loader reruns (search / back-forward navigation).
	$effect(() => {
		searchValue = data.search;
	});
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

	function handleSearch(query: string) {
		const url = new URL(page.url);
		if (query) url.searchParams.set('search', query);
		else url.searchParams.delete('search');
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}

	function openEdit(product: StockProduct) {
		editProduct = product;
		editOpen = true;
	}

	function openCreate() {
		editProduct = null;
		editOpen = true;
	}

	function handleProductSaved(product: Product) {
		invalidateAll();
		if (editProduct == null) {
			// Created (not edited) from pantry: offer to record initial stock.
			editOpen = false;
			stockProduct = product;
			stockOpen = true;
		}
	}
</script>

<svelte:head>
	<title>Despensa — IntelliStock</title>
</svelte:head>

<div class="page">
	<header class="head">
		<h1>Despensa</h1>
	</header>

	<div class="toolbar">
		<SearchInput bind:value={searchValue} onsearch={handleSearch} />
		<PrimaryButton onclick={openCreate}>Adicionar produto</PrimaryButton>
		<PrimaryButton onclick={() => (priceOpen = true)}>Registrar preço</PrimaryButton>
	</div>

	{#if data.stockProducts.length === 0}
		<p class="empty">Nenhum produto encontrado.</p>
	{:else}
		<div class="list">
			{#each data.stockProducts as product (product.id)}
				<PantryProductCard {product} onedit={openEdit} />
			{/each}
		</div>
	{/if}
</div>

<ProductFormModal
	bind:open={editOpen}
	product={editProduct}
	bind:brands
	bind:categories
	measurements={data.measurements}
	onsuccess={handleProductSaved}
/>

<StockTransactionModal
	bind:open={stockOpen}
	product={stockProduct}
	bind:establishments
	onsuccess={invalidateAll}
/>

<AddPriceObservationModal
	bind:open={priceOpen}
	bind:products
	bind:establishments
	bind:brands
	bind:categories
	measurements={data.measurements}
	onsuccess={invalidateAll}
/>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		width: 100%;
		max-width: 52rem;
	}

	.head {
		display: flex;
		justify-content: center;
	}

	h1 {
		font-size: 1.5rem;
		color: var(--color-text);
		text-align: center;
	}

	.toolbar {
		display: flex;
		align-items: center;
		gap: var(--space-md);
	}

	.toolbar :global(.btn) {
		white-space: nowrap;
	}

	.list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(15rem, 1fr));
		gap: var(--space-md);
	}

	.empty {
		color: var(--color-text-muted);
		text-align: center;
	}
</style>
