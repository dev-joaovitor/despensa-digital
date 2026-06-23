<script lang="ts">
	import Modal from './Modal.svelte';
	import Input from './Input.svelte';
	import SearchSelect from './SearchSelect.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import ResourceFormModal from './ResourceFormModal.svelte';
	import type { NamedResource } from '$lib/resources';
	import { createProduct, type Product, type UnitMeasurement } from '$lib/price-observations';

	interface Props {
		open?: boolean;
		initialName?: string;
		brands: NamedResource[];
		categories: NamedResource[];
		measurements: UnitMeasurement[];
		onsuccess: (product: Product) => void;
	}

	let {
		open = $bindable(false),
		initialName = '',
		brands = $bindable(),
		categories = $bindable(),
		measurements,
		onsuccess
	}: Props = $props();

	let name = $state('');
	let brandId = $state<number | null>(null);
	let unitSize = $state('');
	let measurementId = $state<number | null>(null);
	let categoryId = $state<number | null>(null);
	let loading = $state(false);
	let errorMsg = $state('');
	let added = $state(false);

	let brandFormOpen = $state(false);
	let brandInitialName = $state('');
	let categoryFormOpen = $state(false);
	let categoryInitialName = $state('');

	function reset() {
		name = initialName;
		brandId = null;
		unitSize = '';
		measurementId = null;
		categoryId = null;
		errorMsg = '';
		added = false;
	}

	// Reset local state whenever the modal opens.
	$effect(() => {
		if (open) {
			reset();
			loading = false;
		}
	});

	const size = $derived(Number(unitSize));
	const canSubmit = $derived(
		name.trim().length >= 4 &&
			brandId != null &&
			measurementId != null &&
			categoryId != null &&
			Number.isInteger(size) &&
			size > 0
	);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';

		try {
			const result = await createProduct({
				name: name.trim(),
				brand_id: brandId!,
				measurement_id: measurementId!,
				category_id: categoryId!,
				unit_size: size
			});

			if (result.ok && result.data) {
				// The create response only returns ids, so build the enriched product locally.
				const brand = brands.find((b) => b.id === brandId);
				const category = categories.find((c) => c.id === categoryId);
				const measurement = measurements.find((m) => m.id === measurementId);
				onsuccess({
					id: result.data.id,
					name: name.trim(),
					brand: { name: brand?.name ?? '' },
					category: { name: category?.name ?? '' },
					measurement: { size, acronym: measurement?.acronym ?? '' }
				});
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

	function openCreateBrand(text: string) {
		brandInitialName = text;
		brandFormOpen = true;
	}

	function handleBrandCreated(item: NamedResource) {
		brands = [...brands, item];
		brandId = item.id;
	}

	function openCreateCategory(text: string) {
		categoryInitialName = text;
		categoryFormOpen = true;
	}

	function handleCategoryCreated(item: NamedResource) {
		categories = [...categories, item];
		categoryId = item.id;
	}
</script>

<Modal bind:open title="Criar produto">
	{#if added}
		<p class="info" role="status">Produto criado com sucesso.</p>
		<div class="actions">
			<PrimaryButton onclick={addMore}>Adicionar mais</PrimaryButton>
			<PrimaryButton onclick={() => (open = false)}>Voltar a lista</PrimaryButton>
		</div>
	{:else}
		<form onsubmit={handleSubmit}>
			<Input label="Nome" bind:value={name} required />

			<SearchSelect
				label="Marca"
				placeholder="Pesquise uma marca"
				items={brands}
				getLabel={(b) => b.name}
				bind:value={brandId}
				oncreate={openCreateBrand}
			/>

			<Input label="Tamanho unidade" type="number" min="1" step="1" bind:value={unitSize} />

			<SearchSelect
				label="Unidade de medida"
				placeholder="Pesquise uma unidade"
				items={measurements}
				getLabel={(m) => m.acronym}
				bind:value={measurementId}
			/>

			<SearchSelect
				label="Categoria"
				placeholder="Pesquise uma categoria"
				items={categories}
				getLabel={(c) => c.name}
				bind:value={categoryId}
				oncreate={openCreateCategory}
			/>

			{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

			<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Criar</PrimaryButton>
		</form>
	{/if}
</Modal>

<ResourceFormModal
	bind:open={brandFormOpen}
	kind="brands"
	initialName={brandInitialName}
	onsuccess={handleBrandCreated}
/>

<ResourceFormModal
	bind:open={categoryFormOpen}
	kind="categories"
	initialName={categoryInitialName}
	onsuccess={handleCategoryCreated}
/>

<style>
	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		width: 100%;
	}

	.actions {
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
