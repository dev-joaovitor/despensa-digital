<script lang="ts">
	import { page } from '$app/state';
	import Header from '$lib/components/Header.svelte';

	let { data, children } = $props();

	let showNav = $derived(page.url.pathname !== '/');
</script>

<Header user={data.user} nav={showNav} />

<main>
	{@render children()}
</main>

<style>
	main {
		max-width: 64rem;
		margin: 0 auto;
		padding: var(--space-xl) var(--space-lg);
	}
</style>
