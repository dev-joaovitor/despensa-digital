<script lang="ts">
	import ResourcePage from '$lib/components/ResourcePage.svelte';

	let { data } = $props();
</script>

<ResourcePage kind="establishments" items={data.items} />
