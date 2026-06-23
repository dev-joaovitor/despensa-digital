<script lang="ts">
	import PrimaryButton from './PrimaryButton.svelte';
	import { formatCurrency, formatDate, type PriceObservation, type Product } from '$lib/price-observations';

	interface Props {
		observation: PriceObservation;
		onhistory: (product: Product) => void;
	}

	let { observation, onhistory }: Props = $props();
	const { product, current, average_observed_price, lowest } = $derived(observation);
</script>

<article class="card">
	<h3 class="name">{product.name}</h3>
	<p class="meta">{product.brand.name} · {product.measurement.size}{product.measurement.acronym}</p>

	<div class="prices">
		<div class="price">
			<span class="label">Menor</span>
			<span class="value" title="{lowest.establishment.name} · {formatDate(lowest.observed_at)}">
				{formatCurrency(lowest.observed_price)}
			</span>
		</div>
		<div class="price">
			<span class="label">Média</span>
			<span class="value">{formatCurrency(average_observed_price)}</span>
		</div>
		<div class="price">
			<span class="label">Atual</span>
			<span class="value" title="{current.establishment.name} · {formatDate(current.observed_at)}">
				{formatCurrency(current.observed_price)}
			</span>
		</div>
	</div>

	<div class="actions">
		<PrimaryButton onclick={() => onhistory(product)}>Histórico</PrimaryButton>
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
	}

	.name {
		font-size: 1.125rem;
		color: var(--color-text);
	}

	.meta {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.prices {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-md);
	}

	.price {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.price .label {
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.03em;
		color: var(--color-text-muted);
	}

	.price .value {
		font-size: 0.9375rem;
		color: var(--color-text);
		cursor: default;
	}

	.actions {
		display: flex;
		justify-content: flex-end;
	}
</style>
