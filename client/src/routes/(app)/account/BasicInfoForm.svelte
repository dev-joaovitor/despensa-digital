<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import type { SessionUser } from '$lib/auth';
	import Input from '$lib/components/Input.svelte';
	import EmailInput from '$lib/components/EmailInput.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	let { user }: { user: SessionUser } = $props();

	let name = $state(user.full_name);
	let email = $state(user.email);

	let loading = $state(false);
	let errorMsg = $state('');
	let infoMsg = $state('');

	let canSubmit = $derived(name.trim() !== '' && email.trim() !== '');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/users/', {
				method: 'PATCH',
				body: JSON.stringify({ full_name: name, email })
			});

			if (status === 200 && !body.error) {
				infoMsg = 'Dados atualizados com sucesso.';
				await invalidateAll();
				return;
			}

			errorMsg = body.message || 'Não foi possível salvar seus dados.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}
</script>

<form class="card" onsubmit={handleSubmit}>
	<Input label="Nome" bind:value={name} autocomplete="name" required />
	<EmailInput bind:value={email} required />

	{#if infoMsg}
		<p class="info" role="status">{infoMsg}</p>
	{/if}
	{#if errorMsg}
		<p class="error" role="alert">{errorMsg}</p>
	{/if}

	<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Salvar</PrimaryButton>
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
