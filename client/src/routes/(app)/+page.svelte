<script lang="ts">
	import LinkCard from '$lib/components/LinkCard.svelte';

	const main = [
		{ href: '/pantry', title: 'Despensa', subtitle: 'Espelho digital da sua despensa' },
		{ href: '/shopping-list', title: 'Lista de compras', subtitle: 'Sua lista de compras atual' },
		{
			href: '/price-observations',
			title: 'Histórico de preços',
			subtitle: 'Observações e histórico de preços'
		}
	];

	const settings = [
		{
			href: '/account',
			title: 'Você e sua residência',
			subtitle: 'Gerencie seus dados e de sua residência'
		},
		{
			href: '/establishments',
			title: 'Estabelecimentos',
			subtitle: 'Estabelecimentos que você costuma comprar'
		},
		{ href: '/brands', title: 'Marcas', subtitle: 'Marcas preferidas' },
		{
			href: '/categories',
			title: 'Categorias de produtos',
			subtitle: 'As categorias de seus produtos'
		}
	];
</script>

<svelte:head>
	<title>Início — IntelliStock</title>
</svelte:head>

<section class="grid">
	{#each main as card (card.href)}
		<LinkCard {...card} />
	{/each}
</section>

<h2>Configurações</h2>

<section class="grid">
	{#each settings as card (card.href)}
		<LinkCard {...card} />
	{/each}
</section>

<style>
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(16rem, 1fr));
		gap: var(--space-lg);
	}

	h2 {
		font-size: 1.25rem;
		color: var(--color-text);
		margin: var(--space-xl) 0 var(--space-md);
	}
</style>
