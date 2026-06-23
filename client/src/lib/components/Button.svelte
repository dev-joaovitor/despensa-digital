<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from 'svelte/elements';

	type Variant = 'primary' | 'secondary' | 'text';

	interface Props {
		variant?: Variant;
		href?: string;
		type?: 'button' | 'submit' | 'reset';
		disabled?: boolean;
		loading?: boolean;
		onclick?: (event: MouseEvent) => void;
		children: Snippet;
		rest?: HTMLButtonAttributes & HTMLAnchorAttributes;
	}

	let {
		variant = 'primary',
		href,
		type = 'button',
		disabled = false,
		loading = false,
		onclick,
		children,
		...rest
	}: Props = $props();

	let isDisabled = $derived(disabled || loading);
</script>

{#if href}
	<a
		{href}
		class="btn {variant}"
		class:disabled={isDisabled}
		aria-disabled={isDisabled}
		tabindex={isDisabled ? -1 : undefined}
		{...rest}
	>
		{@render children()}
	</a>
{:else}
	<button class="btn {variant}" {type} disabled={isDisabled} {onclick} {...rest}>
		{#if loading}<span class="spinner" aria-hidden="true"></span>{/if}
		{@render children()}
	</button>
{/if}

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-sm);
		font-family: inherit;
		font-size: 1rem;
		line-height: 1;
		padding: 0.75rem 1.25rem;
		border-radius: var(--radius);
		border: 1px solid transparent;
		cursor: pointer;
		text-decoration: none;
		text-align: center;
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
	}

	.btn:disabled,
	.btn.disabled {
		opacity: 0.6;
		cursor: not-allowed;
		pointer-events: none;
	}

	.primary {
		background-color: var(--color-primary);
		color: var(--color-primary-contrast);
	}
	.primary:hover {
		background-color: var(--color-primary-dark);
	}

	.secondary {
		background-color: transparent;
		color: var(--color-primary);
		border-color: var(--color-primary);
	}
	.secondary:hover {
		background-color: var(--color-surface);
	}

	.text {
		background-color: transparent;
		color: var(--color-primary);
		padding: var(--space-xs) var(--space-sm);
	}
	.text:hover {
		text-decoration: underline;
	}

	.spinner {
		width: 1em;
		height: 1em;
		border: 2px solid currentColor;
		border-right-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
