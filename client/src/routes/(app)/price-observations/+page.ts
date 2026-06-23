import { listResources } from '$lib/resources';
import { listPriceObservations, listProducts } from '$lib/price-observations';

export async function load({ url }) {
	const search = url.searchParams.get('search') ?? '';
	const [observations, products, establishments] = await Promise.all([
		listPriceObservations(search),
		listProducts(),
		listResources('establishments')
	]);
	return { observations, products, establishments, search };
}
