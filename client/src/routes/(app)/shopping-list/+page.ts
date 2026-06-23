import { listResources } from '$lib/resources';
import { listProducts, listUnitMeasurements } from '$lib/price-observations';
import { listShoppingItems } from '$lib/shopping-list';

export async function load() {
	const [items, products, establishments, brands, categories, measurements] = await Promise.all([
		listShoppingItems(),
		listProducts(),
		listResources('establishments'),
		listResources('brands'),
		listResources('categories'),
		listUnitMeasurements()
	]);
	return { items, products, establishments, brands, categories, measurements };
}
