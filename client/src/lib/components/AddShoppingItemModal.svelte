<script lang="ts">
	import Modal from './Modal.svelte';
	import Input from './Input.svelte';
	import SearchSelect from './SearchSelect.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import ProductFormModal from './ProductFormModal.svelte';
	import type { NamedResource } from '$lib/resources';
	import { productLabel, type Product, type UnitMeasurement } from '$lib/price-observations';
	import { addShoppingItem } from '$lib/shopping-list';

	interface Props {
		open?: boolean;
		products: Product[];
		brands: NamedResource[];
		categories: NamedResource[];
		measurements: UnitMeasurement[];
		onsuccess: () => void;
	}

	let {
		open = $bindable(false),
		products = $bindable(),
		brands = $bindable(),
		categories = $bindable(),
		measurements,
		onsuccess
	}: Props = $props();

	let productId = $state<number | null>(null);
	let quantity = $state('1');
	let loading = $state(false);
	let errorMsg = $state('');
	let added = $state(false);

	let productFormOpen = $state(false);
	let productInitialName = $state('');

	function reset() {
		productId = null;
		quantity = '1';
		errorMsg = '';
		added = false;
	}

	$effect(() => {
		if (open) {
			reset();
			loading = false;
		}
	});

	const qty = $derived(Number(quantity));
	const canSubmit = $derived(productId != null && Number.isInteger(qty) && qty > 0);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		try {
			const result = await addShoppingItem(productId!, qty);
			if (result.ok) {
				onsuccess();
				added = true;
				return;
			}
			errorMsg = result.message || 'Não foi possível adicionar.';
		} catch {
			errorMsg = 'Erro de conexão. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	function addMore() {
		reset();
	}

	function openCreateProduct(text: string) {
		productInitialName = text;
		productFormOpen = true;
	}

	function handleProductCreated(product: Product) {
		products = [...products, product];
		productId = product.id;
	}
</script>

<Modal bind:open title="Adicionar item">
	{#if added}
		<p class="info" role="status">Item adicionado com sucesso.</p>
		<div class="success-actions">
			<PrimaryButton onclick={addMore}>Adicionar mais</PrimaryButton>
			<PrimaryButton onclick={() => (open = false)}>Voltar a lista</PrimaryButton>
		</div>
	{:else}
		<form onsubmit={handleSubmit}>
			<SearchSelect
				label="Produto"
				placeholder="Pesquise um produto"
				items={products}
				getLabel={productLabel}
				bind:value={productId}
				oncreate={openCreateProduct}
			/>
			<Input label="Quantidade" type="number" min="1" step="1" bind:value={quantity} />

			{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

			<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Adicionar</PrimaryButton>
		</form>
	{/if}
</Modal>

<ProductFormModal
	bind:open={productFormOpen}
	initialName={productInitialName}
	bind:brands
	bind:categories
	{measurements}
	onsuccess={handleProductCreated}
/>

<style>
	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		width: 100%;
	}

	.success-actions {
		display: flex;
		justify-content: center;
		gap: var(--space-sm);
	}

	.info {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
	}
</style>
