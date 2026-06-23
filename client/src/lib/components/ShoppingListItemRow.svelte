<script lang="ts">
	import PrimaryButton from './PrimaryButton.svelte';
	import { productLabel } from '$lib/price-observations';
	import {
		tickShoppingItem,
		updateShoppingItemQuantity,
		removeShoppingItem,
		type ShoppingItem
	} from '$lib/shopping-list';

	interface Props {
		item: ShoppingItem;
		onchange: () => void;
	}

	let { item, onchange }: Props = $props();

	let loading = $state(false);
	let qtyText = $state(String(item.quantity));

	// Keep the local input in sync when the item changes externally (e.g. after refetch).
	$effect(() => {
		qtyText = String(item.quantity);
	});

	async function setQuantity(next: number) {
		if (loading || item.is_checked) return;
		const clamped = Math.min(9999, Math.max(1, Math.trunc(next)));
		if (clamped === item.quantity) {
			qtyText = String(item.quantity);
			return;
		}
		loading = true;
		try {
			const result = await updateShoppingItemQuantity(item.id, clamped);
			if (result.ok) onchange();
		} finally {
			loading = false;
		}
	}

	function commitInput() {
		const parsed = parseInt(qtyText, 10);
		if (Number.isNaN(parsed)) {
			qtyText = String(item.quantity);
			return;
		}
		setQuantity(parsed);
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			(event.target as HTMLInputElement).blur();
		}
	}

	async function toggleCheck() {
		if (loading) return;
		loading = true;
		try {
			const result = await tickShoppingItem(item.id);
			if (result.ok) onchange();
		} finally {
			loading = false;
		}
	}

	async function remove() {
		if (loading || item.is_checked) return;
		loading = true;
		try {
			const result = await removeShoppingItem(item.id);
			if (result.ok) onchange();
		} finally {
			loading = false;
		}
	}
</script>

<div class="row" class:checked={item.is_checked}>
	<button type="button" class="name" class:struck={item.is_checked} onclick={toggleCheck} disabled={loading}>
		{productLabel(item.product)}
	</button>

	<div class="qty">
		<button
			type="button"
			class="step"
			onclick={() => setQuantity(item.quantity - 1)}
			disabled={loading || item.is_checked}
			aria-label="Diminuir">−</button
		>
		<input
			type="text"
			inputmode="numeric"
			bind:value={qtyText}
			onblur={commitInput}
			onkeydown={handleKeydown}
			disabled={loading || item.is_checked}
		/>
		<button
			type="button"
			class="step"
			onclick={() => setQuantity(item.quantity + 1)}
			disabled={loading || item.is_checked}
			aria-label="Aumentar">+</button
		>
	</div>

	<PrimaryButton onclick={remove} disabled={loading || item.is_checked}>Remover</PrimaryButton>
</div>

<style>
	.row {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-sm) var(--space-md);
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
	}

	.name {
		flex: 1;
		text-align: left;
		font-family: inherit;
		font-size: 1rem;
		color: var(--color-text);
		background: transparent;
		border: none;
		padding: 0;
		cursor: pointer;
	}

	.name.struck {
		text-decoration: line-through;
		color: var(--color-text-muted);
	}

	.qty {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
	}

	.step {
		width: 2rem;
		height: 2rem;
		font-size: 1.125rem;
		line-height: 1;
		color: var(--color-text);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		cursor: pointer;
	}

	.step:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.qty input {
		width: 3.5rem;
		text-align: center;
		font-family: inherit;
		font-size: 1rem;
		color: var(--color-text);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		padding: 0.375rem 0.25rem;
	}

	.qty input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.qty input:disabled {
		opacity: 0.5;
	}
</style>
