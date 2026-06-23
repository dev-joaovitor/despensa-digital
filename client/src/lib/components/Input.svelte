<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends HTMLInputAttributes {
		label?: string;
		value?: string;
		error?: string;
		id?: string;
	}

	let {
		label,
		value = $bindable(''),
		error,
		id = crypto.randomUUID(),
		type = 'text',
		...rest
	}: Props = $props();
</script>

<div class="field">
	{#if label}
		<label for={id}>{label}</label>
	{/if}
	<input
		{id}
		{type}
		bind:value
		class:error={!!error}
		aria-invalid={error ? 'true' : undefined}
		{...rest}
	/>
	{#if error}
		<span class="error-message">{error}</span>
	{/if}
</div>

<style>
	.field {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		width: 100%;
	}

	label {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	input {
		font-family: inherit;
		font-size: 1rem;
		color: var(--color-text);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		padding: 0.625rem 0.75rem;
		width: 100%;
		transition: border-color 0.15s ease;
	}

	input::placeholder {
		color: var(--color-text-muted);
		opacity: 0.7;
	}

	input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	input.error {
		border-color: var(--color-error);
	}

	.error-message {
		font-size: 0.8125rem;
		color: var(--color-error);
	}
</style>
