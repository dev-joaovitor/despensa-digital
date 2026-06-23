<script lang="ts">
	import Input from './Input.svelte';

	interface Props {
		label?: string;
		value?: string;
		error?: string;
		name?: string;
		required?: boolean;
		disabled?: boolean;
	}

	let { label = 'Senha', value = $bindable(''), ...rest }: Props = $props();
</script>

<Input
	{label}
	bind:value
	type="password"
	autocomplete="current-password"
	placeholder="••••••••"
	{...rest}
/>
