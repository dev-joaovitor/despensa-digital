<script lang="ts">
	import Modal from './Modal.svelte';
	import SearchSelect from './SearchSelect.svelte';
	import MoneyInput from './MoneyInput.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import ResourceFormModal from './ResourceFormModal.svelte';
	import ProductFormModal from './ProductFormModal.svelte';
	import type { NamedResource } from '$lib/resources';
	import {
		createPriceObservation,
		productLabel,
		type Product,
		type UnitMeasurement
	} from '$lib/price-observations';

	interface Props {
		open?: boolean;
		products: Product[];
		establishments: NamedResource[];
		brands: NamedResource[];
		categories: NamedResource[];
		measurements: UnitMeasurement[];
		onsuccess: () => void;
	}

	let {
		open = $bindable(false),
		products = $bindable(),
		establishments = $bindable(),
		brands = $bindable(),
		categories = $bindable(),
		measurements,
		onsuccess
	}: Props = $props();

	let productId = $state<number | null>(null);
	let establishmentId = $state<number | null>(null);
	let price = $state(0);
	let loading = $state(false);
	let errorMsg = $state('');
	let added = $state(false);

	let establishmentFormOpen = $state(false);
	let establishmentInitialName = $state('');

	let productFormOpen = $state(false);
	let productInitialName = $state('');

	function reset() {
		productId = null;
		establishmentId = null;
		price = 0;
		errorMsg = '';
		added = false;
	}

	$effect(() => {
		if (open) {
			reset();
			loading = false;
		}
	});

	const canSubmit = $derived(productId != null && establishmentId != null && price > 0);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		try {
			const result = await createPriceObservation({
				product_id: productId!,
				establishment_id: establishmentId!,
				price
			});
			if (result.ok) {
				onsuccess();
				added = true;
				return;
			}
			errorMsg = result.message || 'Não foi possível salvar.';
		} catch {
			errorMsg = 'Erro de conexão. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	function addMore() {
		reset();
	}

	function openCreateEstablishment(text: string) {
		establishmentInitialName = text;
		establishmentFormOpen = true;
	}

	function handleEstablishmentCreated(item: NamedResource) {
		establishments = [...establishments, item];
		establishmentId = item.id;
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

<Modal bind:open title="Observação de preço">
	{#if added}
		<p class="info" role="status">Observação de preço registrada com sucesso.</p>
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
			<SearchSelect
				label="Estabelecimento"
				placeholder="Pesquise um estabelecimento"
				items={establishments}
				getLabel={(e) => e.name}
				bind:value={establishmentId}
				oncreate={openCreateEstablishment}
			/>
			<MoneyInput label="Preço" bind:value={price} />

			{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

			<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Salvar</PrimaryButton>
		</form>
	{/if}
</Modal>

<ResourceFormModal
	bind:open={establishmentFormOpen}
	kind="establishments"
	initialName={establishmentInitialName}
	onsuccess={handleEstablishmentCreated}
/>

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
