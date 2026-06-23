<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import type { Household } from '$lib/auth';
	import Input from '$lib/components/Input.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	let { household, isCreator }: { household: Household; isCreator: boolean } = $props();

	let name = $state(household.name);

	let loading = $state(false);
	let errorMsg = $state('');
	let infoMsg = $state('');

	let canSubmit = $derived(isCreator && name.trim() !== '');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/households/', {
				method: 'PATCH',
				body: JSON.stringify({ name })
			});

			if (status === 200 && !body.error) {
				infoMsg = 'Residência atualizada com sucesso.';
				await invalidateAll();
				return;
			}

			errorMsg = body.message || 'Não foi possível salvar.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}
</script>

<form class="card" onsubmit={handleSubmit}>
	<Input label="Nome" bind:value={name} readonly={!isCreator} required />

	{#if infoMsg}
		<p class="info" role="status">{infoMsg}</p>
	{/if}
	{#if errorMsg}
		<p class="error" role="alert">{errorMsg}</p>
	{/if}

	{#if isCreator}
		<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Salvar</PrimaryButton>
	{/if}
</form>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		padding: var(--space-lg);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
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
