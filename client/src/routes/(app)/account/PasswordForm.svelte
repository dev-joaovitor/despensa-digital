<script lang="ts">
	import { apiFetch } from '$lib/api';
	import type { SessionUser } from '$lib/auth';
	import Input from '$lib/components/Input.svelte';
	import PasswordInput from '$lib/components/PasswordInput.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import TextButton from '$lib/components/TextButton.svelte';

	type Mode = 'current' | 'code';

	let { user }: { user: SessionUser } = $props();

	let mode = $state<Mode>('current');
	let currentPassword = $state('');
	let code = $state('');
	let newPassword = $state('');
	let confirmation = $state('');

	let loading = $state(false);
	let sendingCode = $state(false);
	let errorMsg = $state('');
	let infoMsg = $state('');

	let firstField = $derived(mode === 'current' ? currentPassword : code);
	let canSubmit = $derived(
		firstField.trim() !== '' && newPassword !== '' && confirmation !== ''
	);

	async function handleForgotPassword() {
		if (sendingCode) return;
		sendingCode = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/auth/send-recovery-code', {
				method: 'POST',
				body: JSON.stringify({ email: user.email })
			});

			if (status === 200 && !body.error) {
				mode = 'code';
				currentPassword = '';
				infoMsg = 'Código enviado para seu email.';
				return;
			}

			errorMsg = body.message || 'Não foi possível enviar o código. Tente novamente.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			sendingCode = false;
		}
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canSubmit) return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const payload: Record<string, string> = {
				new_password: newPassword,
				new_password_confirmation: confirmation
			};
			if (mode === 'current') payload.old_password = currentPassword;
			else payload.code = code;

			const { status, body } = await apiFetch('/api/v1/users/', {
				method: 'PATCH',
				body: JSON.stringify(payload)
			});

			if (status === 200 && !body.error) {
				infoMsg = 'Senha alterada com sucesso.';
				mode = 'current';
				currentPassword = '';
				code = '';
				newPassword = '';
				confirmation = '';
				return;
			}

			errorMsg = body.message || 'Não foi possível alterar a senha.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}
</script>

<form class="card" onsubmit={handleSubmit}>
	{#if mode === 'current'}
		<PasswordInput label="Senha atual" bind:value={currentPassword} required />
	{:else}
		<Input
			label="Código de recuperação"
			bind:value={code}
			inputmode="numeric"
			maxlength={4}
			autocomplete="one-time-code"
			required
		/>
	{/if}

	<div class="forgot">
		<TextButton type="button" onclick={handleForgotPassword} loading={sendingCode}>
			Esqueci minha senha
		</TextButton>
	</div>

	<PasswordInput
		label="Senha nova"
		bind:value={newPassword}
		autocomplete="new-password"
		required
	/>
	<PasswordInput
		label="Confirmação"
		bind:value={confirmation}
		autocomplete="new-password"
		required
	/>

	{#if infoMsg}
		<p class="info" role="status">{infoMsg}</p>
	{/if}
	{#if errorMsg}
		<p class="error" role="alert">{errorMsg}</p>
	{/if}

	<PrimaryButton type="submit" disabled={!canSubmit} {loading}>Alterar</PrimaryButton>
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

	.forgot {
		display: flex;
		justify-content: flex-end;
		margin-top: calc(-1 * var(--space-sm));
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
