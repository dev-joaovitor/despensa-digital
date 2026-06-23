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

	type Mode = 'invite' | 'create';

	let mode = $state<Mode>('invite');
	let invitationCode = $state('');
	let householdName = $state('');
	let codeVerified = $state(false);

	let fullName = $state('');
	let email = $state('');
	let password = $state('');

	let loading = $state(false);
	let errorMsg = $state('');

	// Account fields appear once a household path is settled.
	let showAccountFields = $derived(mode === 'create' || codeVerified);

	let canCreate = $derived(
		showAccountFields &&
			fullName.trim() !== '' &&
			email.trim() !== '' &&
			password !== '' &&
			(mode === 'create' ? householdName.trim() !== '' : codeVerified)
	);

	function switchMode(next: Mode) {
		mode = next;
		errorMsg = '';
		codeVerified = false;
	}

	function editCode() {
		codeVerified = false;
		errorMsg = '';
	}

	async function verifyCode() {
		if (loading || invitationCode.trim() === '') return;

		loading = true;
		errorMsg = '';

		try {
			const { status, body } = await apiFetch('/api/v1/households/verify-invitation-code', {
				method: 'POST',
				body: JSON.stringify({ invitation_code: invitationCode })
			});

			if (status === 200 && !body.error) {
				codeVerified = true;
				return;
			}

			errorMsg = body.message || 'Código de convite inválido.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (loading || !canCreate) return;

		loading = true;
		errorMsg = '';

		const payload =
			mode === 'invite'
				? { full_name: fullName, email, password, invitation_code: invitationCode }
				: { full_name: fullName, email, password, household_name: householdName };

		try {
			const { status, body } = await apiFetch('/api/v1/users/', {
				method: 'POST',
				body: JSON.stringify(payload)
			});

			if (status === 201 && !body.error) {
				await goto('/');
				return;
			}

			errorMsg = body.message || 'Não foi possível criar a conta. Verifique seus dados.';
		} catch {
			errorMsg = 'Não foi possível conectar. Tente novamente.';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Criação de conta — IntelliStock</title>
</svelte:head>

<main>
	<form class="register-card" onsubmit={handleSubmit}>
		<img class="logo" src={logo} alt="IntelliStock" />

		<h1>Criação de conta</h1>

		{#if mode === 'invite'}
			<Input
				label="Código de convite da residência"
				bind:value={invitationCode}
				placeholder="Insira o código recebido"
				readonly={codeVerified}
				required
			/>

			{#if !codeVerified}
				<PrimaryButton
					type="button"
					disabled={invitationCode.trim() === ''}
					{loading}
					onclick={verifyCode}
				>
					Próximo
				</PrimaryButton>
				<SecondaryButton type="button" onclick={() => switchMode('create')}>
					Criar nova residência
				</SecondaryButton>
			{:else}
				<div class="recovery">
					<TextButton type="button" onclick={editCode}>Alterar código</TextButton>
				</div>
			{/if}
		{:else}
			<Input
				label="Nome da residência"
				bind:value={householdName}
				placeholder="Ex.: Casa da família"
				required
			/>
		{/if}

		{#if showAccountFields}
			<Input label="Nome" bind:value={fullName} autocomplete="name" required />
			<EmailInput label="Email" bind:value={email} required />
			<PasswordInput label="Senha" bind:value={password} autocomplete="new-password" required />
		{/if}

		{#if errorMsg}
			<p class="error" role="alert">{errorMsg}</p>
		{/if}

		{#if showAccountFields}
			<PrimaryButton type="submit" disabled={!canCreate} {loading}>Criar</PrimaryButton>
		{/if}

		{#if mode === 'create'}
			<SecondaryButton type="button" onclick={() => switchMode('invite')}>
				Já tenho um convite
			</SecondaryButton>
		{/if}

		<div class="login-link">
			<TextButton href="/login">Já tenho uma conta</TextButton>
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

	.register-card {
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

	.login-link {
		display: flex;
		justify-content: center;
	}

	.error {
		font-size: 0.875rem;
		color: var(--color-error);
		text-align: center;
	}
</style>
