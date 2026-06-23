<script lang="ts">
	import { goto } from '$app/navigation';
	import logo from '$lib/assets/intellistock-logo.png';
	import { apiFetch } from '$lib/api';
	import EmailInput from '$lib/components/EmailInput.svelte';
	import PasswordInput from '$lib/components/PasswordInput.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import TextButton from '$lib/components/TextButton.svelte';

	let email = $state('');
	let password = $state('');
	let loading = $state(false);
	let errorMsg = $state('');

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading) return;

		loading = true;
		errorMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/auth/login', {
				method: 'POST',
				body: JSON.stringify({ email, password })
			});

			if (status === 200 && !body.error) {
				await goto('/');
				return;
			}

			errorMsg = body.message || 'Não foi possível entrar. Verifique seus dados.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login — IntelliStock</title>
</svelte:head>

<main>
	<form class="login-card" onsubmit={handleSubmit}>
		<img class="logo" src={logo} alt="IntelliStock" />

		<h1>Login</h1>

		<EmailInput bind:value={email} required />
		<PasswordInput bind:value={password} required />

		<div class="recovery">
			<TextButton href="/recover-password">Esqueci minha senha</TextButton>
		</div>

		{#if errorMsg}
			<p class="error" role="alert">{errorMsg}</p>
		{/if}

		<PrimaryButton type="submit" {loading}>Entrar</PrimaryButton>
		<SecondaryButton href="/register">Criar nova conta</SecondaryButton>
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

	.login-card {
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

	.recovery {
		display: flex;
		justify-content: flex-end;
		margin-top: calc(-1 * var(--space-sm));
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
		text-align: center;
	}
</style>
