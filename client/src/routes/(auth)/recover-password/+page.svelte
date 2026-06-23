<script lang="ts">
	import { goto } from '$app/navigation';
	import logo from '$lib/assets/intellistock-logo.png';
	import { apiFetch } from '$lib/api';
	import Input from '$lib/components/Input.svelte';
	import EmailInput from '$lib/components/EmailInput.svelte';
	import PasswordInput from '$lib/components/PasswordInput.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import TextButton from '$lib/components/TextButton.svelte';

	type Step = 'email' | 'code' | 'password';

	let step = $state<Step>('email');
	let email = $state('');
	let code = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');

	let loading = $state(false);
	let errorMsg = $state('');
	let infoMsg = $state('');

	let canSubmitPassword = $derived(
		newPassword.length >= 6 && confirmPassword !== '' && newPassword === confirmPassword
	);

	function backToEmail() {
		step = 'email';
		code = '';
		errorMsg = '';
		infoMsg = '';
	}

	async function sendCode(): Promise<boolean> {
		const { status, body } = await apiFetch('/api/v1/auth/send-recovery-code', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
		return status === 200 && !body.error;
	}

	async function handleSendCode() {
		if (loading || email.trim() === '') return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			if (await sendCode()) {
				step = 'code';
			} else {
				errorMsg = 'Não foi possível enviar o código. Tente novamente.';
			}
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	async function handleResend() {
		if (loading) return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			if (await sendCode()) {
				infoMsg = 'Novo código enviado.';
			} else {
				errorMsg = 'Não foi possível reenviar o código. Tente novamente.';
			}
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	async function handleVerify() {
		if (loading || code.trim().length !== 4) return;
		loading = true;
		errorMsg = '';
		infoMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/auth/verify-recovery-code', {
				method: 'POST',
				body: JSON.stringify({ email, code })
			});

			if (status === 200 && !body.error) {
				step = 'password';
				return;
			}

			errorMsg = body.message || 'O código está errado ou já expirou.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	async function handleChangePassword() {
		if (loading || !canSubmitPassword) return;
		loading = true;
		errorMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/auth/change-password', {
				method: 'POST',
				body: JSON.stringify({
					new_password: newPassword,
					new_password_confirmation: confirmPassword
				})
			});

			if (status === 200 && !body.error) {
				await goto('/');
				return;
			}

			errorMsg = body.message || 'Não foi possível alterar a senha.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	// Single submit handler dispatches by current step for Enter-key support.
	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (step === 'email') handleSendCode();
		else if (step === 'code') handleVerify();
		else handleChangePassword();
	}
</script>

<svelte:head>
	<title>Recuperar senha — IntelliStock</title>
</svelte:head>

<main>
	<form class="recover-card" onsubmit={handleSubmit}>
		<img class="logo" src={logo} alt="IntelliStock" />

		<h1>Recuperar senha</h1>

		{#if step === 'email'}
			<EmailInput bind:value={email} required />
			<PrimaryButton type="submit" disabled={email.trim() === ''} {loading}>
				Enviar código
			</PrimaryButton>
		{:else if step === 'code'}
			<p class="sent-to">Enviado para o email: {email}</p>
			<SecondaryButton type="button" onclick={backToEmail}>Alterar email</SecondaryButton>

			<Input
				label="Código de verificação"
				bind:value={code}
				inputmode="numeric"
				maxlength={4}
				autocomplete="one-time-code"
				required
			/>
			<div class="resend">
				<TextButton type="button" onclick={handleResend}>Reenviar código</TextButton>
			</div>

			{#if infoMsg}
				<p class="info" role="status">{infoMsg}</p>
			{/if}

			<PrimaryButton type="submit" disabled={code.trim().length !== 4} {loading}>
				Verificar
			</PrimaryButton>
		{:else}
			<PasswordInput
				label="Nova senha"
				bind:value={newPassword}
				autocomplete="new-password"
				required
			/>
			<PasswordInput
				label="Confirmar nova senha"
				bind:value={confirmPassword}
				autocomplete="new-password"
				required
			/>
			<PrimaryButton type="submit" disabled={!canSubmitPassword} {loading}>
				Alterar senha
			</PrimaryButton>
		{/if}

		{#if errorMsg}
			<p class="error" role="alert">{errorMsg}</p>
		{/if}

		<div class="login-link">
			<TextButton href="/login">Voltar ao login</TextButton>
		</div>
	</form>
</main>

<style>
	main {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-lg);
		background-color: var(--color-surface);
	}

	.recover-card {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
		width: 100%;
		max-width: 22rem;
		padding: var(--space-xl);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
	}

	.logo {
		display: block;
		width: 12rem;
		max-width: 100%;
		height: auto;
		margin: 0 auto var(--space-sm);
	}

	h1 {
		text-align: center;
		font-size: 1.5rem;
		color: var(--color-text);
	}

	.sent-to {
		font-size: 0.875rem;
		color: var(--color-text-muted);
		text-align: center;
	}

	.resend {
		display: flex;
		justify-content: flex-end;
		margin-top: calc(-1 * var(--space-sm));
	}

	.login-link {
		display: flex;
		justify-content: center;
	}

	.info {
		font-size: 0.875rem;
		color: var(--color-text-muted);
		text-align: center;
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
		text-align: center;
	}
</style>
