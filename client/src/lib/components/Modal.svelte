<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open?: boolean;
		title?: string;
		onclose?: () => void;
		children: Snippet;
	}

	let { open = $bindable(false), title, onclose, children }: Props = $props();

	function close() {
		open = false;
		onclose?.();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') close();
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

{#if open}
	<div
		class="backdrop"
		role="presentation"
		onclick={(e) => {
			if (e.target === e.currentTarget) close();
		}}
	>
		<div class="dialog" role="dialog" aria-modal="true" aria-label={title}>
			{#if title}
				<h2>{title}</h2>
			{/if}
			{@render children()}
		</div>
	</div>
{/if}

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-md);
		background-color: rgba(43, 33, 26, 0.45);
		z-index: 100;
	}

	.dialog {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-md);
		width: 100%;
		max-width: 28rem;
		padding: var(--space-lg);
		text-align: center;
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
	}

	h2 {
		font-size: 1.25rem;
		color: var(--color-text);
	}
</style>
