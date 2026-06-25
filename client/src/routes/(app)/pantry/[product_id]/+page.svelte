<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import BatchCard from '$lib/components/BatchCard.svelte';
	import BatchTransactionModal from '$lib/components/BatchTransactionModal.svelte';
	import StockTransactionModal from '$lib/components/StockTransactionModal.svelte';
	import type { StockBatch, BatchTransactionType } from '$lib/stock';
	import type { Product } from '$lib/price-observations';

	let { data } = $props();

	let establishments = $state(data.establishments);

	$effect(() => {
		establishments = data.establishments;
	});

	let purchaseOpen = $state(false);

	let txOpen = $state(false);
	let txBatch = $state<StockBatch | null>(null);
	let txType = $state<BatchTransactionType>('consumption');

	function openTransaction(batch: StockBatch, type: BatchTransactionType) {
		txBatch = batch;
		txType = type;
		txOpen = true;
	}
</script>

<svelte:head>
	<title>{data.product ? data.product.name : 'Lotes'} — IntelliStock</title>
</svelte:head>

<div class="page">
	<a class="back" href="/pantry">Voltar à Despensa</a>

	{#if !data.product}
		<p class="empty">Produto não encontrado.</p>
	{:else}
		<header class="head">
			<h1>{data.product.name}</h1>
			<p class="subtitle">
				{data.product.brand.name} · {data.product.measurement.size}{data.product.measurement
					.acronym}
			</p>
		</header>

		<div class="toolbar">
			<PrimaryButton onclick={() => (purchaseOpen = true)}>Comprar</PrimaryButton>
		</div>

		{#if data.batches.length === 0}
			<p class="empty">Nenhum lote cadastrado.</p>
		{:else}
			<div class="list">
				{#each data.batches as batch (batch.id)}
					<BatchCard {batch} onaction={openTransaction} />
				{/each}
			</div>
		{/if}
	{/if}
</div>

{#if data.product}
	<StockTransactionModal
		bind:open={purchaseOpen}
		title="Novo lote"
		product={data.product as Product}
		bind:establishments
		onsuccess={invalidateAll}
	/>
{/if}

<BatchTransactionModal
	bind:open={txOpen}
	batch={txBatch}
	type={txType}
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
	}

	.back {
		font-size: 0.875rem;
		color: var(--color-primary);
		text-decoration: none;
	}

	.back:hover {
		text-decoration: underline;
	}

	.head {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-xs);
		text-align: center;
	}

	h1 {
		font-size: 1.5rem;
		color: var(--color-text);
	}

	.subtitle {
		font-size: 0.9375rem;
		color: var(--color-text-muted);
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
