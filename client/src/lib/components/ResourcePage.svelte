<script lang="ts">
	import {
		deleteResource,
		RESOURCE_CONFIG,
		type NamedResource,
		type ResourceKind
	} from '$lib/resources';
	import PrimaryButton from './PrimaryButton.svelte';
	import ResourceList from './ResourceList.svelte';
	import ResourceFormModal from './ResourceFormModal.svelte';
	import ConfirmModal from './ConfirmModal.svelte';

	interface Props {
		kind: ResourceKind;
		items: NamedResource[];
	}

	let { kind, items: initialItems }: Props = $props();

	const config = $derived(RESOURCE_CONFIG[kind]);

	let items = $state(initialItems);

	let formOpen = $state(false);
	let editingItem = $state<NamedResource | null>(null);

	let confirmOpen = $state(false);
	let removingItem = $state<NamedResource | null>(null);
	let removing = $state(false);

	function openCreate() {
		editingItem = null;
		formOpen = true;
	}

	function openEdit(item: NamedResource) {
		editingItem = item;
		formOpen = true;
	}

	function handleSuccess(saved: NamedResource) {
		const index = items.findIndex((i) => i.id === saved.id);
		if (index === -1) {
			items = [...items, saved];
		} else {
			items[index] = saved;
		}
	}

	function openRemove(item: NamedResource) {
		removingItem = item;
		confirmOpen = true;
	}

	async function confirmRemove() {
		if (!removingItem || removing) return;
		removing = true;
		const { ok } = await deleteResource(kind, removingItem.id);
		if (ok) {
			items = items.filter((i) => i.id !== removingItem!.id);
			confirmOpen = false;
		}
		removing = false;
	}
</script>

<svelte:head>
	<title>{config.title} — IntelliStock</title>
</svelte:head>

<div class="page">
	<header class="head">
		<h1>{config.title}</h1>
		<PrimaryButton onclick={openCreate}>Adicionar {config.singular}</PrimaryButton>
	</header>

	<ResourceList {items} onedit={openEdit} onremove={openRemove} />
</div>

<ResourceFormModal bind:open={formOpen} {kind} item={editingItem} onsuccess={handleSuccess} />

<ConfirmModal
	bind:open={confirmOpen}
	title={`Deseja realmente excluir ${config.gender === 'm' ? 'um' : 'uma'} ${config.singular}?`}
	subtitle="Esta exclusão irá impactar diretamente seus produtos e histórico de preços"
	loading={removing}
	onconfirm={confirmRemove}
/>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-lg);
		width: 100%;
		max-width: 52rem;
	}

	.head {
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		align-items: center;
		gap: var(--space-md);
	}

	h1 {
		grid-column: 2;
		font-size: 1.5rem;
		color: var(--color-text);
		text-align: center;
	}

	.head :global(.btn) {
		grid-column: 3;
		justify-self: end;
	}
</style>
