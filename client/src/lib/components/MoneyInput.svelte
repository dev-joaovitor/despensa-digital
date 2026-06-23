<script lang="ts">
	interface Props {
		value?: number;
		label?: string;
	}

	let { value = $bindable(0), label }: Props = $props();

	const id = crypto.randomUUID();

	function formatCents(cents: number): string {
		const reais = Math.floor(cents / 100);
		const frac = String(cents % 100).padStart(2, '0');
		return `${reais.toLocaleString('pt-BR')},${frac}`;
	}

	// Display string derived from the numeric value (rounded to cents).
	let display = $state(formatCents(Math.round(value * 100)));

	function handleInput(event: Event) {
		const raw = (event.target as HTMLInputElement).value.replace(/\D/g, '');
		const cents = raw ? parseInt(raw, 10) : 0;
		value = cents / 100;
		display = formatCents(cents);
	}
</script>

<div class="field">
	{#if label}<label for={id}>{label}</label>{/if}
	<div class="control">
		<span class="prefix">R$</span>
		<input
			{id}
			type="text"
			inputmode="numeric"
			bind:value={display}
			oninput={handleInput}
		/>
	</div>
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

	.control {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		padding: 0.625rem 0.75rem;
		transition: border-color 0.15s ease;
	}

	.control:focus-within {
		border-color: var(--color-primary);
	}

	.prefix {
		color: var(--color-text-muted);
		font-size: 1rem;
	}

	input {
		flex: 1;
		font-family: inherit;
		font-size: 1rem;
		color: var(--color-text);
		background-color: transparent;
		border: none;
		padding: 0;
		width: 100%;
	}

	input:focus {
		outline: none;
	}
</style>
