import { apiFetch } from '$lib/api';

export interface StockPurchaseInput {
	product_id: number;
	establishment_id: number;
	quantity: number;
	unit_price: number;
	expiration_date: string; // YYYY-MM-DD
}

export async function createStockPurchase(
	input: StockPurchaseInput
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch('/api/v1/stock/transact', {
		method: 'POST',
		body: JSON.stringify({ type: 'purchase', ...input })
	});
	return { ok: status === 201 && !body.error, message: body.message };
}
