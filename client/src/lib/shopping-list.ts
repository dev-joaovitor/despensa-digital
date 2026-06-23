import { apiFetch } from '$lib/api';
import type { Product } from '$lib/price-observations';

export interface ShoppingItem {
	id: number;
	product: Product;
	quantity: number;
	is_checked: boolean;
}

export interface SubmitItem {
	product_id: number;
	establishment_id: number;
	expiration_date: string;
	price: number;
	quantity: number;
}

export async function listShoppingItems(): Promise<ShoppingItem[]> {
	const { status, body } = await apiFetch<ShoppingItem[]>('/api/v1/shopping-list');
	return status === 200 && body.data ? body.data : [];
}

export async function addShoppingItem(
	product_id: number,
	quantity: number
): Promise<{ ok: boolean; data?: { id: number }; message?: string }> {
	const { status, body } = await apiFetch<{ id: number }>('/api/v1/shopping-list', {
		method: 'POST',
		body: JSON.stringify({ product_id, quantity })
	});
	return { ok: status === 201 && !body.error, data: body.data, message: body.message };
}

export async function updateShoppingItemQuantity(
	id: number,
	quantity: number
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch(`/api/v1/shopping-list/${id}`, {
		method: 'PATCH',
		body: JSON.stringify({ quantity })
	});
	return { ok: status === 200 && !body.error, message: body.message };
}

export async function tickShoppingItem(id: number): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch(`/api/v1/shopping-list/${id}/tick`, { method: 'POST' });
	return { ok: status === 200 && !body.error, message: body.message };
}

export async function removeShoppingItem(id: number): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch(`/api/v1/shopping-list/${id}`, { method: 'DELETE' });
	return { ok: status === 200 && !body.error, message: body.message };
}

export async function submitShoppingList(
	items: SubmitItem[]
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch('/api/v1/shopping-list/submit', {
		method: 'POST',
		body: JSON.stringify({ items })
	});
	return { ok: status === 200 && !body.error, message: body.message };
}
