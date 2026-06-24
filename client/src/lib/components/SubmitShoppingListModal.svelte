<script lang="ts">
	import Modal from './Modal.svelte';
	import Input from './Input.svelte';
	import DateInput from './DateInput.svelte';
	import SearchSelect from './SearchSelect.svelte';
	import MoneyInput from './MoneyInput.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import ResourceFormModal from './ResourceFormModal.svelte';
	import type { NamedResource } from '$lib/resources';
	import { productLabel } from '$lib/price-observations';
	import { submitShoppingList, type ShoppingItem } from '$lib/shopping-list';

	interface Props {
		open?: boolean;
		items: ShoppingItem[];
		establishments: NamedResource[];
		onsuccess: () => void;
	}

	let {
		open = $bindable(false),
		items,
		establishments = $bindable(),
		onsuccess
	}: Props = $props();

	interface DraftRow {
		item: ShoppingItem;
		price: number;
		quantity: string;
		expiration_date: string;
		establishment_id: number | null;
		// Whether the user picked an establishment for this row specifically.
		overridden: boolean;
	}

	let globalEstablishmentId = $state<number | null>(null);
	let rows = $state<DraftRow[]>([]);
	let loading = $state(false);
	let errorMsg = $state('');

	let establishmentFormOpen = $state(false);
	let establishmentInitialName = $state('');

	// Rebuild the local draft whenever the modal opens. Nothing here hits the API.
	$effect(() => {
		if (open) {
			globalEstablishmentId = null;
			errorMsg = '';
			loading = false;
			rows = items.map((item) => ({
				item,
				price: 0,
				quantity: String(item.quantity),
				expiration_date: '',
				establishment_id: null,
				overridden: false
			}));
		}
	});

	function applyGlobal(id: number) {
		for (const row of rows) {
			if (!row.overridden) row.establishment_id = id;
		}
	}

	function overrideRow(row: DraftRow) {
		row.overridden = true;
	}

	const canSubmit = $derived(
		rows.length > 0 &&
			rows.every(
				(r) =>
					r.price > 0 &&
					r.establishment_id != null &&
					Number.isInteger(Number(r.quantity)) &&
					Number(r.quantity) > 0 &&
					/^\d{4}-\d{2}-\d{2}$/.test(r.expiration_date)
			)
	);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		try {
			const result = await submitShoppingList(
				rows.map((r) => ({
					product_id: r.item.product.id,
					establishment_id: r.establishment_id!,
					expiration_date: r.expiration_date,
					price: r.price,
					quantity: Number(r.quantity)
				}))
			);
			if (result.ok) {
				open = false;
				onsuccess();
				return;
			}
			errorMsg = result.message || 'Não foi possível enviar a lista.';
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
		globalEstablishmentId = item.id;
		applyGlobal(item.id);
	}
</script>

<Modal bind:open title="Enviar lista" size="lg">
	<form onsubmit={handleSubmit}>
		<div class="global">
			<SearchSelect
				label="Estabelecimento"
				placeholder="Pesquise um estabelecimento"
				items={establishments}
				getLabel={(e) => e.name}
				bind:value={globalEstablishmentId}
				oncreate={openCreateEstablishment}
				onchange={applyGlobal}
			/>
		</div>

		{#if rows.length === 0}
			<p class="info">Nenhum item marcado para enviar.</p>
		{:else}
			<ul class="rows">
				{#each rows as row (row.item.id)}
					<li class="item">
						<span class="label">{productLabel(row.item.product)}</span>
						<div class="fields">
							<MoneyInput bind:value={row.price} />
							<Input type="number" min="1" step="1" bind:value={row.quantity} />
							<DateInput bind:value={row.expiration_date} />
							<SearchSelect
								placeholder="Estabelecimento"
								items={establishments}
								getLabel={(e) => e.name}
								bind:value={row.establishment_id}
								onchange={() => overrideRow(row)}
							/>
						</div>
					</li>
				{/each}
			</ul>
		{/if}

		{#if errorMsg}<p class="error" role="alert">{errorMsg}</p>{/if}

		<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Enviar lista</PrimaryButton>
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

	.global {
		max-width: 24rem;
	}

	.rows {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		margin: 0;
		padding: 0;
		list-style: none;
		max-height: 50vh;
		overflow-y: auto;
	}

	.item {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		padding: var(--space-sm);
		background-color: var(--color-surface);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
	}

	.label {
		font-size: 0.9375rem;
		color: var(--color-text);
	}

	.fields {
		display: grid;
		grid-template-columns: 1fr 5rem 1fr 1.5fr;
		gap: var(--space-sm);
		align-items: center;
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
