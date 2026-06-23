import { listResources } from '$lib/resources';
import { listPriceObservations, listProducts, listUnitMeasurements } from '$lib/price-observations';

export async function load({ url }) {
	const search = url.searchParams.get('search') ?? '';
	const [observations, products, establishments, brands, categories, measurements] =
		await Promise.all([
			listPriceObservations(search),
			listProducts(),
			listResources('establishments'),
			listResources('brands'),
			listResources('categories'),
			listUnitMeasurements()
		]);
	return { observations, products, establishments, brands, categories, measurements, search };
}
