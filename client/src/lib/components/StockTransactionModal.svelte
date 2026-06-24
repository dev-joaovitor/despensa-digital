<script lang="ts">
	import Modal from './Modal.svelte';
	import Input from './Input.svelte';
	import SearchSelect from './SearchSelect.svelte';
	import MoneyInput from './MoneyInput.svelte';
	import DateInput from './DateInput.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import SecondaryButton from './SecondaryButton.svelte';
	import ResourceFormModal from './ResourceFormModal.svelte';
	import type { NamedResource } from '$lib/resources';
	import { productLabel, type Product } from '$lib/price-observations';
	import { createStockPurchase } from '$lib/stock';

	interface Props {
		open?: boolean;
		title?: string;
		product: Product | null;
		establishments: NamedResource[];
		onsuccess?: () => void;
		onskip?: () => void;
	}

	let {
		open = $bindable(false),
		title = 'Estoque inicial',
		product,
		establishments = $bindable(),
		onsuccess,
		onskip
	}: Props = $props();

	let establishmentId = $state<number | null>(null);
	let quantity = $state('');
	let price = $state(0);
	let expirationDate = $state('');
	let loading = $state(false);
	let errorMsg = $state('');

	let establishmentFormOpen = $state(false);
	let establishmentInitialName = $state('');

	function reset() {
		establishmentId = null;
		quantity = '';
		price = 0;
		expirationDate = '';
		errorMsg = '';
	}

	$effect(() => {
		if (open) {
			reset();
			loading = false;
		}
	});

	const canSubmit = $derived(
		product != null &&
			establishmentId != null &&
			Number.isInteger(Number(quantity)) &&
			Number(quantity) > 0 &&
			price > 0 &&
			/^\d{4}-\d{2}-\d{2}$/.test(expirationDate)
	);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		try {
			const result = await createStockPurchase({
				product_id: product!.id,
				establishment_id: establishmentId!,
				quantity: Number(quantity),
				unit_price: price,
				expiration_date: expirationDate
			});
			if (result.ok) {
				open = false;
				onsuccess?.();
				return;
			}
			errorMsg = result.message || 'Não foi possível salvar.';
		} catch {
			errorMsg = 'Erro de conexão. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	function skip() {
		open = false;
		onskip?.();
	}

	function openCreateEstablishment(text: string) {
		establishmentInitialName = text;
		establishmentFormOpen = true;
	}

	function handleEstablishmentCreated(item: NamedResource) {
		establishments = [...establishments, item];
		establishmentId = item.id;
	}
</script>

<Modal bind:open {title}>
	<form onsubmit={handleSubmit}>
		{#if product}<p class="product">{productLabel(product)}</p>{/if}

		<SearchSelect
			label="Estabelecimento"
			placeholder="Pesquise um estabelecimento"
			items={establishments}
			getLabel={(e) => e.name}
			bind:value={establishmentId}
			oncreate={openCreateEstablishment}
		/>
		<Input label="Quantidade" type="number" min="1" step="1" bind:value={quantity} />
		<MoneyInput label="Preço unitário" bind:value={price} />
		<DateInput label="Data de validade" bind:value={expirationDate} />

		{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

		<div class="actions">
			<SecondaryButton onclick={skip}>Agora não</SecondaryButton>
			<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Salvar</PrimaryButton>
		</div>
	</form>
</Modal>

<ResourceFormModal
	bind:open={establishmentFormOpen}
	kind="establishments"
	initialName={establishmentInitialName}
	onsuccess={handleEstablishmentCreated}
/>

<style>
	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		width: 100%;
	}

	.product {
		font-size: 0.9375rem;
		color: var(--color-text);
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-sm);
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
	}
</style>
