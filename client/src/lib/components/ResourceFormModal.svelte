<script lang="ts">
	import {
		createResource,
		updateResource,
		RESOURCE_CONFIG,
		type NamedResource,
		type ResourceKind
	} from '$lib/resources';
	import Modal from './Modal.svelte';
	import Input from './Input.svelte';
	import PrimaryButton from './PrimaryButton.svelte';

	interface Props {
		open?: boolean;
		kind: ResourceKind;
		item?: NamedResource | null;
		onsuccess: (item: NamedResource) => void;
	}

	let { open = $bindable(false), kind, item = null, onsuccess }: Props = $props();

	const config = $derived(RESOURCE_CONFIG[kind]);
	const singular = $derived(config.singular);
	const isEdit = $derived(!!item);
	const novo = $derived(config.gender === 'm' ? 'Novo' : 'Nova');
	const added_ = $derived(config.gender === 'm' ? 'adicionado' : 'adicionada');
	const title = $derived(isEdit ? `Editar ${singular}` : `${novo} ${singular}`);
	const submitLabel = $derived(isEdit ? 'Salvar' : 'Adicionar');

	let name = $state('');
	let loading = $state(false);
	let errorMsg = $state('');
	let added = $state(false);

	// Reset local state whenever the modal opens for a (possibly different) item.
	$effect(() => {
		if (open) {
			name = item?.name ?? '';
			loading = false;
			errorMsg = '';
			added = false;
		}
	});

	const canSubmit = $derived(name.trim().length >= 4);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';

		try {
			if (item) {
				// Update returns no body, so reconstruct the saved item locally.
				const result = await updateResource(kind, item.id, name.trim());
				if (result.ok) {
					onsuccess({ ...item, name: name.trim() });
					open = false;
					return;
				}
				errorMsg = result.message || 'Não foi possível salvar.';
				return;
			}

			const result = await createResource(kind, name.trim());
			if (result.ok && result.data) {
				onsuccess(result.data);
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
		name = '';
		added = false;
		errorMsg = '';
	}
</script>

<Modal bind:open {title}>
	{#if added}
		<p class="info" role="status">
			{singular[0].toUpperCase() + singular.slice(1)} {added_} com sucesso.
		</p>
		<div class="actions">
			<PrimaryButton onclick={addMore}>Adicionar mais</PrimaryButton>
			<PrimaryButton onclick={() => (open = false)}>Voltar a lista</PrimaryButton>
		</div>
	{:else}
		<form onsubmit={handleSubmit}>
			<Input label="Nome" bind:value={name} required />

			{#if errorMsg}
				<p class="error" role="alert">{errorMsg}</p>
			{/if}

			<PrimaryButton type="submit" disabled={!canSubmit} {loading}>{submitLabel}</PrimaryButton>
		</form>
	{/if}
</Modal>

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
