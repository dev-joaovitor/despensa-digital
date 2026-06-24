import { listResources } from '$lib/resources';
import { listProducts, listUnitMeasurements } from '$lib/price-observations';
import { listStockProducts } from '$lib/pantry';

export async function load({ url }) {
	const search = url.searchParams.get('search') ?? '';
	const [stockProducts, products, establishments, brands, categories, measurements] =
		await Promise.all([
			listStockProducts(search),
			listProducts(),
			listResources('establishments'),
			listResources('brands'),
			listResources('categories'),
			listUnitMeasurements()
		]);
	return { stockProducts, products, establishments, brands, categories, measurements, search };
}
