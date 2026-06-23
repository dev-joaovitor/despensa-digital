<script lang="ts">
	interface Props {
		value?: string;
		placeholder?: string;
		delay?: number;
		onsearch: (query: string) => void;
	}

	let { value = $bindable(''), placeholder = 'Pesquise', delay = 2000, onsearch }: Props = $props();

	let timer: ReturnType<typeof setTimeout> | undefined;

	function handleInput() {
		clearTimeout(timer);
		timer = setTimeout(() => onsearch(value.trim()), delay);
	}
</script>

<div class="search">
	<input
		type="search"
		bind:value
		{placeholder}
		oninput={handleInput}
		aria-label={placeholder}
	/>
</div>

<style>
	.search {
		width: 100%;
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
</style>
