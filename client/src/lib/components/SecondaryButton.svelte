<script lang="ts">
	import type { Snippet } from 'svelte';
	import Button from './Button.svelte';

	interface Props {
		href?: string;
		type?: 'button' | 'submit' | 'reset';
		disabled?: boolean;
		loading?: boolean;
		onclick?: (event: MouseEvent) => void;
		children: Snippet;
	}

	let { children, ...rest }: Props = $props();
</script>

<Button variant="secondary" {...rest}>
	{@render children()}
</Button>
