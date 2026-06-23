<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		text: string;
		children: Snippet;
	}

	let { text, children }: Props = $props();
</script>

<span class="tooltip" tabindex="0" role="button">
	{@render children()}
	<span class="bubble" role="note">{text}</span>
</span>

<style>
	.tooltip {
		position: relative;
		display: inline-flex;
		cursor: help;
		outline: none;
	}

	.bubble {
		position: absolute;
		bottom: calc(100% + 0.5rem);
		left: 50%;
		transform: translateX(-50%);
		z-index: 20;
		padding: var(--space-xs) var(--space-sm);
		font-size: 0.75rem;
		line-height: 1.2;
		white-space: nowrap;
		color: var(--color-bg);
		background-color: var(--color-text);
		border-radius: var(--radius-sm);
		box-shadow: var(--shadow);
		opacity: 0;
		visibility: hidden;
		pointer-events: none;
		transition: opacity 0.08s ease;
	}

	.bubble::after {
		content: '';
		position: absolute;
		top: 100%;
		left: 50%;
		transform: translateX(-50%);
		border: 5px solid transparent;
		border-top-color: var(--color-text);
	}

	.tooltip:hover .bubble,
	.tooltip:focus-visible .bubble {
		opacity: 1;
		visibility: visible;
	}
</style>
