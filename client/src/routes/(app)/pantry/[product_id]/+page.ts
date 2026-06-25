import { listProductBatches } from '$lib/stock';
import { listStockProducts } from '$lib/pantry';
import { listResources } from '$lib/resources';

export async function load({ params }) {
	const productId = Number(params.product_id);
	const [batches, products, establishments] = await Promise.all([
		listProductBatches(productId),
		listStockProducts(),
		listResources('establishments')
	]);
	const product = products.find((p) => p.id === productId) ?? null;
	return { productId, product, batches, establishments };
}
