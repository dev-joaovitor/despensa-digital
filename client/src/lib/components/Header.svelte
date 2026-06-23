<script lang="ts">
	import { page } from '$app/state';
	import type { SessionUser } from '$lib/auth';
	import { getInitials } from '$lib/initials';
	import logo from '$lib/assets/intellistock-logo-minimalist.png';

	interface Props {
		user: SessionUser;
		nav?: boolean;
	}

	let { user, nav = true }: Props = $props();

	const links = [
		{ href: '/', label: 'Início' },
		{ href: '/pantry', label: 'Despensa' },
		{ href: '/shopping-list', label: 'Lista de compras' },
		{ href: '/price-observations', label: 'Observação de preços' }
	];

	const accountHref: string = '/account';
	let initials = $derived(getInitials(user.full_name));
</script>

<header>
	<a class="brand" href="/" aria-label="IntelliStock">
		<img src={logo} alt="IntelliStock" />
	</a>

	{#if nav}
		<nav aria-label="Navegação principal">
			{#each links as link (link.href)}
				<a
					href={link.href}
					class="link"
					class:active={page.url.pathname === link.href}
					aria-current={page.url.pathname === link.href ? 'page' : undefined}
				>
					{link.label}
				</a>
			{/each}

			<a
				class="avatar"
				href={accountHref}
				class:active={page.url.pathname === accountHref}
				aria-label="Sua conta"
				title={user.full_name}
			>
				{initials}
			</a>
		</nav>
	{/if}
</header>

<style>
	header {
		position: sticky;
		top: 0;
		z-index: 10;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-md);
		padding: var(--space-sm) var(--space-lg);
		background-color: var(--color-bg);
		border-bottom: 1px solid var(--color-border);
		box-shadow: var(--shadow);
	}

	.brand {
		display: inline-flex;
		align-items: center;
	}

	.brand img {
		height: 2.25rem;
		width: auto;
		display: block;
	}

	nav {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		flex-wrap: wrap;
		justify-content: flex-end;
	}

	.link {
		color: var(--color-text-muted);
		text-decoration: none;
		font-size: 0.95rem;
		padding: var(--space-xs) 0;
		border-bottom: 2px solid transparent;
		transition: color 0.15s ease, border-color 0.15s ease;
	}

	.link:hover {
		color: var(--color-text);
	}

	.link.active {
		color: var(--color-primary);
		border-bottom-color: var(--color-primary);
	}

	.avatar {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 2.25rem;
		height: 2.25rem;
		border-radius: 50%;
		background-color: var(--color-primary);
		color: var(--color-primary-contrast);
		font-size: 0.85rem;
		font-weight: 600;
		text-decoration: none;
		flex-shrink: 0;
		transition: background-color 0.15s ease;
	}

	.avatar:hover,
	.avatar.active {
		background-color: var(--color-primary-dark);
	}
</style>
