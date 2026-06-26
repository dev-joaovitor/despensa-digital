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
	import {
		consumeBatch,
		wasteBatch,
		correctBatch,
		type StockBatch,
		type BatchTransactionType,
		type CorrectionInput
	} from '$lib/stock';

	interface Props {
		open?: boolean;
		batch: StockBatch | null;
		type: BatchTransactionType;
		establishments: NamedResource[];
		onsuccess?: () => void;
	}

	let {
		open = $bindable(false),
		batch,
		type,
		establishments = $bindable(),
		onsuccess
	}: Props = $props();

	const TITLES: Record<BatchTransactionType, string> = {
		consumption: 'Consumir',
		waste: 'Descartar',
		correction: 'Corrigir lote'
	};

	let quantity = $state('');
	let establishmentId = $state<number | null>(null);
	let price = $state(0);
	let expirationDate = $state('');
	let loading = $state(false);
	let errorMsg = $state('');

	let establishmentFormOpen = $state(false);
	let establishmentInitialName = $state('');

	function dateOnly(value: string): string {
		const match = /^(\d{4}-\d{2}-\d{2})/.exec(value);
		return match ? match[1] : '';
	}

	// Reset the form to the batch's current values each time the modal opens.
	$effect(() => {
		if (open) {
			quantity = batch ? String(batch.remaining_quantity) : '';
			establishmentId = batch ? batch.establishment.id : null;
			price = batch ? batch.unit_price : 0;
			expirationDate = batch ? dateOnly(batch.expiration_date) : '';
			errorMsg = '';
			loading = false;
		}
	});

	const canSubmit = $derived(
		batch != null &&
			Number.isInteger(Number(quantity)) &&
			Number(quantity) > 0
	);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		try {
			const qty = Number(quantity);
			let result: { ok: boolean; message?: string };
			if (type === 'consumption') {
				result = await consumeBatch(batch!.id, qty);
			} else if (type === 'waste') {
				result = await wasteBatch(batch!.id, qty);
			} else {
				const extra: CorrectionInput = {};
				if (establishmentId != null) extra.establishment_id = establishmentId;
				if (price > 0) extra.unit_price = price;
				if (/^\d{4}-\d{2}-\d{2}$/.test(expirationDate)) extra.expiration_date = expirationDate;
				result = await correctBatch(batch!.id, qty, extra);
			}
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

	function openCreateEstablishment(text: string) {
		establishmentInitialName = text;
		establishmentFormOpen = true;
	}

	function handleEstablishmentCreated(item: NamedResource) {
		establishments = [...establishments, item];
		establishmentId = item.id;
	}
</script>

<Modal bind:open title={TITLES[type]}>
	<form onsubmit={handleSubmit}>
		<Input label="Quantidade" type="number" min="1" step="1" bind:value={quantity} />

		{#if type === 'correction'}
			<SearchSelect
				label="Estabelecimento"
				placeholder="Pesquise um estabelecimento"
				items={establishments}
				getLabel={(e) => e.name}
				bind:value={establishmentId}
				oncreate={openCreateEstablishment}
			/>
			<MoneyInput label="Preço unitário" bind:value={price} />
			<DateInput label="Data de validade" bind:value={expirationDate} />
		{/if}

		{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

		<div class="form-actions">
			<SecondaryButton onclick={() => (open = false)}>Cancelar</SecondaryButton>
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

	.form-actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-sm);
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
	}
</style>
