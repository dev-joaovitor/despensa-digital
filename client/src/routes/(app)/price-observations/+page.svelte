<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import PriceObservationCard from '$lib/components/PriceObservationCard.svelte';
	import AddPriceObservationModal from '$lib/components/AddPriceObservationModal.svelte';
	import PriceHistoryModal from '$lib/components/PriceHistoryModal.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import { invalidateAll } from '$app/navigation';
	import type { Product } from '$lib/price-observations';

	let { data } = $props();

	let searchValue = $state(data.search);
	let establishments = $state(data.establishments);
	let products = $state(data.products);
	let brands = $state(data.brands);
	let categories = $state(data.categories);

	let addOpen = $state(false);
	let historyOpen = $state(false);
	let historyProduct = $state<Product | null>(null);

	// Keep the input in sync when the page reloads (e.g. back/forward navigation).
	$effect(() => {
		searchValue = data.search;
	});

	$effect(() => {
		establishments = data.establishments;
	});

	$effect(() => {
		products = data.products;
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

	function openHistory(product: Product) {
		historyProduct = product;
		historyOpen = true;
	}
</script>

<svelte:head>
	<title>Observações de preço — IntelliStock</title>
</svelte:head>

<div class="page">
	<header class="head">
		<h1>Observações de preço</h1>
	</header>

	<div class="toolbar">
		<SearchInput bind:value={searchValue} onsearch={handleSearch} />
		<PrimaryButton onclick={() => (addOpen = true)}>Adicionar observação de preço</PrimaryButton>
	</div>

	{#if data.observations.length === 0}
		<p class="empty">Nenhuma observação de preço encontrada.</p>
	{:else}
		<div class="list">
			{#each data.observations as observation (observation.product.id)}
				<PriceObservationCard {observation} onhistory={openHistory} />
			{/each}
		</div>
	{/if}
</div>

<AddPriceObservationModal
	bind:open={addOpen}
	bind:products
	bind:establishments
	bind:brands
	bind:categories
	measurements={data.measurements}
	onsuccess={invalidateAll}
/>

<PriceHistoryModal bind:open={historyOpen} product={historyProduct} {establishments} />

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
