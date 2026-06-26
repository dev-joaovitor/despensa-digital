<script lang="ts">
	import { onMount } from 'svelte';
	import logo from '$lib/assets/intellistock-logo.png';
	import { isAuthenticated } from '$lib/auth';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';

	let loggedIn = $state<boolean | null>(null);

	onMount(async () => {
		loggedIn = await isAuthenticated();
	});
</script>

<svelte:head>
	<title>Página não encontrada — IntelliStock</title>
</svelte:head>

<main>
	<section class="card">
		<img class="logo" src={logo} alt="IntelliStock" />

		<h1>404</h1>
		<p class="message">Página não encontrada.</p>

		{#if loggedIn === true}
			<PrimaryButton href="/">Início</PrimaryButton>
		{:else if loggedIn === false}
			<PrimaryButton href="/login">Entrar</PrimaryButton>
			<SecondaryButton href="/register">Cadastre-se</SecondaryButton>
		{/if}
	</section>
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

	.card {
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
		font-size: 2.5rem;
		color: var(--color-text);
	}

	.message {
		text-align: center;
		color: var(--color-text-muted);
		margin-bottom: var(--space-sm);
	}
</style>
