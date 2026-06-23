<script lang="ts">
	import Modal from './Modal.svelte';
	import PrimaryButton from './PrimaryButton.svelte';
	import SecondaryButton from './SecondaryButton.svelte';

	interface Props {
		open?: boolean;
		title: string;
		subtitle?: string;
		confirmLabel?: string;
		loading?: boolean;
		onconfirm: () => void;
	}

	let {
		open = $bindable(false),
		title,
		subtitle,
		confirmLabel = 'Excluir',
		loading = false,
		onconfirm
	}: Props = $props();
</script>

<Modal bind:open {title}>
	{#if subtitle}
		<p class="subtitle">{subtitle}</p>
	{/if}
	<div class="actions">
		<PrimaryButton onclick={onconfirm} {loading}>{confirmLabel}</PrimaryButton>
		<SecondaryButton onclick={() => (open = false)} disabled={loading}>Cancelar</SecondaryButton>
	</div>
</Modal>

<style>
	.subtitle {
		color: var(--color-text-muted);
	}

	.actions {
		display: flex;
		justify-content: center;
		gap: var(--space-sm);
	}
</style>
