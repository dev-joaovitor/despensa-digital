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
	return { ok: status === 200 && !body.error, message: body.message };
}

export interface StockBatch {
	id: number;
	unit_price: number;
	initial_quantity: number;
	remaining_quantity: number;
	expiration_date: string;
	created_at: string;
	updated_at: string;
	establishment: { id: number, name: string };
}

export async function listProductBatches(productId: number): Promise<StockBatch[]> {
	const { status, body } = await apiFetch<StockBatch[]>(
		`/api/v1/stock/products/${productId}/batches`
	);
	return status === 200 && body.data ? body.data : [];
}

export type BatchTransactionType = 'consumption' | 'waste' | 'correction';

export interface CorrectionInput {
	establishment_id?: number;
	unit_price?: number;
	expiration_date?: string; // YYYY-MM-DD
}

async function transact(payload: Record<string, unknown>): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch('/api/v1/stock/transact', {
		method: 'POST',
		body: JSON.stringify(payload)
	});
	return { ok: status === 200 && !body.error, message: body.message };
}

export function consumeBatch(batchId: number, quantity: number) {
	return transact({ type: 'consumption', batch_id: batchId, quantity });
}

export function wasteBatch(batchId: number, quantity: number) {
	return transact({ type: 'waste', batch_id: batchId, quantity });
}

export function correctBatch(batchId: number, quantity: number, extra: CorrectionInput = {}) {
	return transact({ type: 'correction', batch_id: batchId, quantity, ...extra });
}
