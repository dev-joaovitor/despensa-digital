<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open?: boolean;
		title?: string;
		size?: 'sm' | 'lg';
		onclose?: () => void;
		children: Snippet;
	}

	let { open = $bindable(false), title, size = 'sm', onclose, children }: Props = $props();

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
		<div class="dialog {size}" role="dialog" aria-modal="true" aria-label={title}>
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
		gap: var(--space-md);
		width: 100%;
		max-height: calc(100vh - 2 * var(--space-md));
		padding: var(--space-lg);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
	}

	.dialog.sm {
		align-items: center;
		max-width: 28rem;
		text-align: center;
	}

	.dialog.lg {
		max-width: 40rem;
	}

	h2 {
		font-size: 1.25rem;
		color: var(--color-text);
	}
</style>
