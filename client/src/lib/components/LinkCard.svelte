<script lang="ts">
	interface Props {
		href: string;
		title: string;
		subtitle: string;
	}

	let { href, title, subtitle }: Props = $props();
</script>

<a class="card" {href}>
	<h3>{title}</h3>
	<p>{subtitle}</p>
</a>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: var(--space-xs);
		padding: var(--space-lg);
		background-color: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		box-shadow: var(--shadow);
		text-decoration: none;
		transition:
			border-color 0.15s ease,
			transform 0.15s ease,
			box-shadow 0.15s ease;
	}

	.card:hover {
		border-color: var(--color-primary);
		transform: translateY(-2px);
	}

	h3 {
		font-size: 1.125rem;
		color: var(--color-text);
	}

	p {
		font-size: 0.9rem;
		color: var(--color-text-muted);
	}
</style>
