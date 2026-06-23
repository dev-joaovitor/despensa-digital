<script lang="ts">
	import { apiFetch } from '$lib/api';
	import type { Household } from '$lib/auth';
	import Input from '$lib/components/Input.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	let { household }: { household: Household } = $props();

	let code = $state(household.invitation_code);

	let generating = $state(false);
	let errorMsg = $state('');
	let infoMsg = $state('');

	async function handleCopy() {
		errorMsg = '';
		infoMsg = '';
		try {
			await navigator.clipboard.writeText(code);
			infoMsg = 'Convite copiado para a área de transferência.';
		} catch {
			errorMsg = 'Não foi possível copiar. Copie o código manualmente.';
		}
	}

	async function handleGenerate() {
		if (generating) return;
		generating = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const { status, body } = await apiFetch<{ new_code: string }>('/api/v1/households/code', {
				method: 'POST'
			});

			if (status === 200 && !body.error && body.data) {
				code = body.data.new_code;
				infoMsg = 'Novo código gerado.';
				return;
			}

			errorMsg = body.message || 'Não foi possível gerar um novo código.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			generating = false;
		}
	}
</script>

<div class="card">
	<Input label="Código da residência" value={code} readonly />

	{#if infoMsg}
		<p class="info" role="status">{infoMsg}</p>
	{/if}
	{#if errorMsg}
		<p class="error" role="alert">{errorMsg}</p>
	{/if}

	<div class="actions">
		<PrimaryButton type="button" onclick={handleCopy}>Copiar convite</PrimaryButton>
		<PrimaryButton type="button" onclick={handleGenerate} loading={generating}>
			Gerar novo código
		</PrimaryButton>
	</div>
</div>

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

	.actions {
		display: flex;
		flex-wrap: wrap;
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
