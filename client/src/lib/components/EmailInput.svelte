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

	let { label = 'E-mail', value = $bindable(''), ...rest }: Props = $props();
</script>

<Input
	{label}
	bind:value
	type="email"
	autocomplete="email"
	placeholder="seu@email.com"
	{...rest}
/>
