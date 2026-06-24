<script lang="ts">
	import SecondaryButton from './SecondaryButton.svelte';
	import type { StockProduct } from '$lib/pantry';

	interface Props {
		product: StockProduct;
		onedit: (product: StockProduct) => void;
	}

	let { product, onedit }: Props = $props();

	const { initial, remaining } = $derived(product.stock);
	const hasStock = $derived(initial > 0);
	const fillPct = $derived(hasStock ? Math.round((remaining / initial) * 100) : 0);
	const meterLabel = $derived(hasStock ? `${remaining}/${initial} em estoque` : 'Sem estoque');
</script>

<article class="card">
	<header class="top">
		<h3 class="name">{product.name}</h3>
		<span class="chip">{product.category.name}</span>
	</header>

	<p class="brand">{product.brand.name}</p>

	<div class="bottom">
		<SecondaryButton onclick={() => onedit(product)}>Editar</SecondaryButton>
		<span class="size">{product.measurement.size}{product.measurement.acronym}</span>
	</div>

	<div class="meter" class:empty={!hasStock}>
		<div class="meter-fill" style="width: {fillPct}%"></div>
		<span class="meter-label">{meterLabel}</span>
	</div>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		padding: var(--space-md);
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		overflow: hidden;
	}

	.top {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-sm);
	}

	.name {
		font-size: 1.125rem;
		color: var(--color-text);
	}

	.chip {
		flex-shrink: 0;
		font-size: 0.75rem;
		color: var(--color-primary);
		background-color: transparent;
		border: 1px solid var(--color-primary);
		border-radius: 999px;
		padding: 0.125rem 0.625rem;
	}

	.brand {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.bottom {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-sm);
	}

	.size {
		font-size: 0.8125rem;
		color: var(--color-text-muted);
	}

	/* Flush-attached to the card's bottom edge: bleed over the card padding. */
	.meter {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		height: 1.5rem;
		margin: var(--space-sm) calc(-1 * var(--space-md)) calc(-1 * var(--space-md));
		background-color: var(--color-border);
	}

	.meter-fill {
		position: absolute;
		inset: 0 auto 0 0;
		background-color: var(--color-primary-light);
	}

	.meter.empty .meter-fill {
		background-color: transparent;
	}

	.meter-label {
		position: relative;
		font-size: 0.75rem;
		color: var(--color-primary-contrast);
	}
</style>
