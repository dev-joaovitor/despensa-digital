<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import BasicInfoForm from './BasicInfoForm.svelte';
	import PasswordForm from './PasswordForm.svelte';
	import HouseholdNameForm from './HouseholdNameForm.svelte';
	import HouseholdCodeForm from './HouseholdCodeForm.svelte';

	let { data } = $props();

	let loggingOut = $state(false);

	async function handleLogout() {
		if (loggingOut) return;
		loggingOut = true;
		try {
			// Session is cleared server-side; redirect regardless of the response body.
			await apiFetch('/api/v1/auth/logout', { method: 'POST' });
		} catch {
			// Ignore network errors and still leave the page.
		}
		await goto('/login');
	}
</script>

<svelte:head>
	<title>Você e sua residência — IntelliStock</title>
</svelte:head>

<div class="account">
	<section class="section">
		<div class="section-head">
			<h2>Você</h2>
			<SecondaryButton onclick={handleLogout} loading={loggingOut}>Sair</SecondaryButton>
		</div>
		<div class="forms">
			<BasicInfoForm user={data.user} />
			<PasswordForm user={data.user} />
		</div>
	</section>

	{#if data.household}
		<section class="section">
			<h2>Residência</h2>
			<div class="forms">
				<HouseholdNameForm household={data.household} isCreator={data.isCreator} />
				{#if data.isCreator}
					<HouseholdCodeForm household={data.household} />
				{/if}
			</div>
		</section>
	{/if}
</div>

<style>
	.account {
		display: flex;
		flex-direction: column;
		gap: var(--space-xl);
		width: 100%;
		max-width: 52rem;
	}

	.section {
		display: flex;
		flex-direction: column;
		gap: var(--space-md);
	}

	.section-head {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: var(--space-md);
	}

	.forms {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(18rem, 1fr));
		gap: var(--space-lg);
	}

	h2 {
		font-size: 1.25rem;
		color: var(--color-text);
	}
</style>
