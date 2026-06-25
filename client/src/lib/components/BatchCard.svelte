<script lang="ts">
	import SecondaryButton from './SecondaryButton.svelte';
	import type { StockBatch, BatchTransactionType } from '$lib/stock';

	interface Props {
		batch: StockBatch;
		onaction: (batch: StockBatch, type: BatchTransactionType) => void;
	}

	let { batch, onaction }: Props = $props();

	const initial = $derived(batch.initial_quantity);
	const remaining = $derived(batch.remaining_quantity);
	const hasStock = $derived(initial > 0);
	const fillPct = $derived(hasStock ? Math.round((remaining / initial) * 100) : 0);
	const meterLabel = $derived(hasStock ? `${remaining}/${initial} restante` : 'Sem estoque');

	// Show only the date portion of the ISO timestamp, formatted as DD/MM/YYYY.
	const expiration = $derived.by(() => {
		const match = /^(\d{4})-(\d{2})-(\d{2})/.exec(batch.expiration_date);
		return match ? `${match[3]}/${match[2]}/${match[1]}` : '—';
	});
</script>

<article class="card">
	<header class="top">
		<span class="label">Validade</span>
		<span class="date">{expiration}</span>
	</header>

	<p class="establishment">{batch.establishment.name}</p>

	<div class="actions">
		<SecondaryButton onclick={() => onaction(batch, 'consumption')}>Consumir</SecondaryButton>
		<SecondaryButton onclick={() => onaction(batch, 'waste')}>Descartar</SecondaryButton>
		<SecondaryButton onclick={() => onaction(batch, 'correction')}>Corrigir</SecondaryButton>
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
		align-items: baseline;
		justify-content: space-between;
		gap: var(--space-sm);
	}

	.label {
		font-size: 0.8125rem;
		color: var(--color-text-muted);
	}

	.date {
		font-size: 1.125rem;
		color: var(--color-text);
	}

	.establishment {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.actions {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-sm);
	}

	.actions :global(.btn) {
		flex: 1;
		white-space: nowrap;
	}

	/* Flush-attached to the card's bottom edge: bleed over the card padding. */
	.meter {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		height: 1.5rem;
		margin: auto calc(-1 * var(--space-md)) calc(-1 * var(--space-md));
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
