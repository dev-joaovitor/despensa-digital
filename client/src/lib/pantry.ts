import { apiFetch } from '$lib/api';

export interface StockProduct {
	id: number;
	name: string;
	created_at: string;
	updated_at: string;
	brand: { id: number; name: string };
	category: { id: number; name: string };
	measurement: { id: number; size: number; acronym: string };
	stock: { initial: number; remaining: number };
}

export async function listStockProducts(search?: string): Promise<StockProduct[]> {
	const query = search ? `?search=${encodeURIComponent(search)}` : '';
	const { status, body } = await apiFetch<StockProduct[]>(`/api/v1/stock/products${query}`);
	return status === 200 && body.data ? body.data : [];
}

export async function getStockProduct(id: number): Promise<StockProduct | null> {
	const { status, body } = await apiFetch<StockProduct>(`/api/v1/stock/products/${id}`);
	return status === 200 && body.data ? body.data : null;
}
