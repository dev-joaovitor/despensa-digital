<script lang="ts">
	interface Props {
		value?: string; // YYYY-MM-DD, '' while incomplete
		label?: string;
	}

	let { value = $bindable(''), label }: Props = $props();

	const id = crypto.randomUUID();

	// Parse an incoming value into the three parts (runs once on init).
	const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(value);
	let day = $state(match ? match[3] : '');
	let month = $state(match ? match[2] : '');
	let year = $state(match ? match[1] : '');

	function digits(raw: string, max: number): string {
		return raw.replace(/\D/g, '').slice(0, max);
	}

	// Emit a single YYYY-MM-DD value only when all parts are present.
	function sync() {
		if (day.length >= 1 && month.length >= 1 && year.length === 4) {
			value = `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
		} else {
			value = '';
		}
	}

	function handleDay() {
		day = digits(day, 2);
		sync();
	}

	function handleMonth() {
		month = digits(month, 2);
		sync();
	}

	function handleYear() {
		year = digits(year, 4);
		sync();
	}
</script>

<div class="field">
	{#if label}<label for={id}>{label}</label>{/if}
	<div class="parts">
		<input
			{id}
			class="dd"
			type="text"
			inputmode="numeric"
			placeholder="DD"
			bind:value={day}
			oninput={handleDay}
		/>
		<span class="sep">/</span>
		<input
			class="mm"
			type="text"
			inputmode="numeric"
			placeholder="MM"
			bind:value={month}
			oninput={handleMonth}
		/>
		<span class="sep">/</span>
		<input
			class="yyyy"
			type="text"
			inputmode="numeric"
			placeholder="YYYY"
			bind:value={year}
			oninput={handleYear}
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

	.parts {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		padding: 0.625rem 0.75rem;
		transition: border-color 0.15s ease;
	}

	.parts:focus-within {
		border-color: var(--color-primary);
	}

	.sep {
		color: var(--color-text-muted);
	}

	input {
		font-family: inherit;
		font-size: 1rem;
		color: var(--color-text);
		background-color: transparent;
		border: none;
		padding: 0;
		text-align: center;
	}

	input:focus {
		outline: none;
	}

	input::placeholder {
		color: var(--color-text-muted);
		opacity: 0.7;
	}

	.dd,
	.mm {
		width: 2.5rem;
	}

	.yyyy {
		width: 3.5rem;
	}
</style>
