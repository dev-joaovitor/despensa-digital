import { listProductBatches } from '$lib/stock';
import { getStockProduct } from '$lib/pantry';
import { listResources } from '$lib/resources';

export async function load({ params }) {
	const productId = Number(params.product_id);
	const [batches, product, establishments] = await Promise.all([
		listProductBatches(productId),
		getStockProduct(productId),
		listResources('establishments')
	]);
	return { productId, product, batches, establishments };
}
