<script lang="ts" generics="T extends { id: number }">
	import TextButton from './TextButton.svelte';

	interface Props {
		items: T[];
		getLabel: (item: T) => string;
		value?: number | null;
		label?: string;
		placeholder?: string;
		oncreate?: (text: string) => void;
	}

	let {
		items,
		getLabel,
		value = $bindable(null),
		label,
		placeholder,
		oncreate
	}: Props = $props();

	let text = $state('');
	let open = $state(false);
	const id = crypto.randomUUID();

	const selected = $derived(value != null ? items.find((i) => i.id === value) : undefined);

	// Keep the input text in sync with the selected item (e.g. when set externally).
	$effect(() => {
		if (selected) text = getLabel(selected);
	});

	const matches = $derived(
		text.trim()
			? items.filter((i) => getLabel(i).toLowerCase().includes(text.trim().toLowerCase()))
			: items
	);

	function select(item: T) {
		value = item.id;
		text = getLabel(item);
		open = false;
	}

	function handleInput() {
		open = true;
		// Typing invalidates a previous selection until the user picks again.
		if (selected && text !== getLabel(selected)) value = null;
	}
</script>

<div class="search-select">
	{#if label}<label for={id}>{label}</label>{/if}
	<div class="row">
		<div class="control">
			<input
				{id}
				type="text"
				bind:value={text}
				{placeholder}
				autocomplete="off"
				oninput={handleInput}
				onfocus={() => (open = true)}
				onblur={() => setTimeout(() => (open = false), 150)}
			/>
			{#if open && matches.length}
				<ul class="options" role="listbox">
					{#each matches as item (item.id)}
						<li>
							<button type="button" onmousedown={() => select(item)}>{getLabel(item)}</button>
						</li>
					{/each}
				</ul>
			{/if}
		</div>
		{#if oncreate}
			<TextButton onclick={() => oncreate(text.trim())}>Criar</TextButton>
		{/if}
	</div>
</div>

<style>
	.search-select {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		width: 100%;
	}

	label {
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}

	.row {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}

	.control {
		position: relative;
		flex: 1;
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

	input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.options {
		position: absolute;
		top: calc(100% + 2px);
		left: 0;
		right: 0;
		z-index: 10;
		max-height: 12rem;
		overflow-y: auto;
		margin: 0;
		padding: 0;
		list-style: none;
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		box-shadow: var(--shadow);
	}

	.options button {
		display: block;
		width: 100%;
		text-align: left;
		font-family: inherit;
		font-size: 0.9375rem;
		color: var(--color-text);
		background-color: transparent;
		border: none;
		padding: 0.5rem 0.75rem;
		cursor: pointer;
	}

	.options button:hover {
		background-color: var(--color-surface);
	}
</style>
