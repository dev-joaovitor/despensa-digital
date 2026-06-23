<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import Input from './Input.svelte';

	interface Props {
		label?: string;
		value?: string;
		error?: string;
		name?: string;
		required?: boolean;
		disabled?: boolean;
		autocomplete?: HTMLInputAttributes['autocomplete'];
	}

	let {
		label = 'Senha',
		value = $bindable(''),
		autocomplete = 'current-password',
		...rest
	}: Props = $props();
</script>

<Input
	{label}
	bind:value
	type="password"
	{autocomplete}
	placeholder="••••••••"
	{...rest}
/>
