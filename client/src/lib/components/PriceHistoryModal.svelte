<script lang="ts">
	import Modal from './Modal.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import type { NamedResource } from '$lib/resources';
	import {
		listPriceHistory,
		formatCurrency,
		formatDate,
		productLabel,
		type HistoryEntry,
		type Product
	} from '$lib/price-observations';

	interface Props {
		open?: boolean;
		product: Product | null;
		establishments: NamedResource[];
	}

	let { open = $bindable(false), product, establishments }: Props = $props();

	function toISO(d: Date): string {
		const y = d.getFullYear();
		const m = String(d.getMonth() + 1).padStart(2, '0');
		const day = String(d.getDate()).padStart(2, '0');
		return `${y}-${m}-${day}`;
	}

	function shift(iso: string, days: number): string {
		const d = new Date(`${iso}T00:00:00`);
		d.setDate(d.getDate() + days);
		return toISO(d);
	}

	let establishmentId = $state<number | null>(null);
	let from = $state('');
	let to = $state('');
	let entries = $state<HistoryEntry[]>([]);
	let loading = $state(false);

	// Initialize the window (last 7 days) and clear selection each time it opens.
	$effect(() => {
		if (open) {
			const today = toISO(new Date());
			from = shift(today, -7);
			to = today;
			establishmentId = null;
			entries = [];
		}
	});

	// Refetch whenever a chip is selected or the window changes.
	$effect(() => {
		if (!open || product == null || establishmentId == null) return;
		const pid = product.id;
		const eid = establishmentId;
		const f = from;
		const t = to;
		loading = true;
		listPriceHistory(pid, { establishment_id: eid, from: f, to: t })
			.then((data) => {
				entries = data;
			})
			.finally(() => {
				loading = false;
			});
	});

	const today = toISO(new Date());
	const canAdvance = $derived(to < today);

	function slideBack() {
		from = shift(from, -7);
		to = shift(to, -7);
	}

	function slideForward() {
		let nextFrom = shift(from, 7);
		let nextTo = shift(to, 7);
		// Don't slide past today; clamp the window so `to` lands on today.
		if (nextTo > today) {
			const overshoot = daysBetween(today, nextTo);
			nextFrom = shift(nextFrom, -overshoot);
			nextTo = today;
		}
		from = nextFrom;
		to = nextTo;
	}

	function daysBetween(a: string, b: string): number {
		const ms = new Date(`${b}T00:00:00`).getTime() - new Date(`${a}T00:00:00`).getTime();
		return Math.round(ms / 86400000);
	}
</script>

<Modal bind:open size="lg" title={product ? productLabel(product) : ''}>
	<div class="chips">
		{#each establishments as e (e.id)}
			<button
				type="button"
				class="chip"
				class:active={establishmentId === e.id}
				onclick={() => (establishmentId = e.id)}
			>
				{e.name}
			</button>
		{/each}
	</div>

	<div class="list">
		{#if establishmentId == null}
			<p class="empty">Selecione um estabelecimento para ver o histórico.</p>
		{:else if loading}
			<p class="empty">Carregando...</p>
		{:else if entries.length === 0}
			<p class="empty">Nenhuma observação no período.</p>
		{:else}
			{#each entries as entry (entry.id)}
				<div class="row">
					<span class="date">{formatDate(entry.observed_at)}</span>
					<span class="price">{formatCurrency(entry.observed_price)}</span>
				</div>
			{/each}
		{/if}
	</div>

	<div class="window-actions">
		<PrimaryButton onclick={slideBack}>- 7 dias</PrimaryButton>
		<PrimaryButton onclick={slideForward} disabled={!canAdvance}>+ 7 dias</PrimaryButton>
	</div>
</Modal>

<style>
	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: var(--space-xs);
		width: 100%;
	}

	.chip {
		font-family: inherit;
		font-size: 0.8125rem;
		color: var(--color-primary);
		background-color: transparent;
		border: 1px solid var(--color-primary);
		border-radius: 999px;
		padding: 0.25rem 0.75rem;
		cursor: pointer;
		transition:
			background-color 0.15s ease,
			color 0.15s ease;
	}

	.chip.active {
		background-color: var(--color-primary);
		color: var(--color-primary-contrast);
	}

	.list {
		display: flex;
		flex-direction: column;
		width: 100%;
		max-height: 50vh;
		overflow-y: auto;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
	}

	.empty {
		padding: var(--space-md);
		font-size: 0.875rem;
		color: var(--color-text-muted);
		text-align: center;
	}

	.row {
		display: flex;
		justify-content: space-between;
		padding: 0.625rem 0.75rem;
		border-bottom: 1px solid var(--color-border);
	}

	.row:last-child {
		border-bottom: none;
	}

	.date {
		color: var(--color-text-muted);
		font-size: 0.875rem;
	}

	.price {
		color: var(--color-text);
		font-size: 0.9375rem;
	}

	.window-actions {
		display: flex;
		justify-content: space-between;
		gap: var(--space-sm);
		width: 100%;
	}
</style>
