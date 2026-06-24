import { apiFetch } from '$lib/api';

export interface Brand {
	name: string;
}

export interface Measurement {
	size: number;
	acronym: string;
}

export interface UnitMeasurement {
	id: number;
	name: string;
	acronym: string;
}

export interface Product {
	id: number;
	name: string;
	brand: Brand;
	category?: { name: string };
	measurement: Measurement;
}

export interface PriceSnapshot {
	observed_price: number;
	observed_at: string;
	establishment: { name: string };
}

export interface PriceObservation {
	product: Product;
	current: PriceSnapshot;
	average_observed_price: number;
	lowest: PriceSnapshot;
}

export interface HistoryEntry {
	id: number;
	product_id: number;
	establishment: { id: number; name: string };
	observed_price: number;
	observed_at: string;
}

export interface CreatePriceObservationInput {
	product_id: number;
	establishment_id: number;
	price: number;
}

export async function listPriceObservations(search?: string): Promise<PriceObservation[]> {
	const query = search ? `?search=${encodeURIComponent(search)}` : '';
	const { status, body } = await apiFetch<PriceObservation[]>(`/api/v1/price-observations${query}`);
	return status === 200 && body.data ? body.data : [];
}

export async function listProducts(): Promise<Product[]> {
	const { status, body } = await apiFetch<Product[]>('/api/v1/products');
	return status === 200 && body.data ? body.data : [];
}

export async function listUnitMeasurements(): Promise<UnitMeasurement[]> {
	const { status, body } = await apiFetch<UnitMeasurement[]>('/api/v1/unit-measurements');
	return status === 200 && body.data ? body.data : [];
}

export interface CreateProductInput {
	name: string;
	brand_id: number;
	measurement_id: number;
	category_id: number;
	unit_size: number;
}

export async function createProduct(
	input: CreateProductInput
): Promise<{ ok: boolean; data?: { id: number }; message?: string }> {
	const { status, body } = await apiFetch<{ id: number }>('/api/v1/products', {
		method: 'POST',
		body: JSON.stringify(input)
	});
	return { ok: status === 201 && !body.error, data: body.data, message: body.message };
}

export async function updateProduct(
	id: number,
	input: CreateProductInput
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch(`/api/v1/products/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(input)
	});
	return { ok: status === 200 && !body.error, message: body.message };
}

export async function createPriceObservation(
	input: CreatePriceObservationInput
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch('/api/v1/price-observations', {
		method: 'POST',
		body: JSON.stringify(input)
	});
	return { ok: status === 201 && !body.error, message: body.message };
}

export async function listPriceHistory(
	productId: number,
	params: { establishment_id: number; from: string; to: string }
): Promise<HistoryEntry[]> {
	const query = new URLSearchParams({
		establishment_id: String(params.establishment_id),
		from: params.from,
		to: params.to
	});
	const { status, body } = await apiFetch<HistoryEntry[]>(
		`/api/v1/price-observations/history/${productId}?${query}`
	);
	return status === 200 && body.data ? body.data : [];
}

const currencyFmt = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' });
const dateFmt = new Intl.DateTimeFormat('pt-BR');

export function formatCurrency(value: number): string {
	return currencyFmt.format(value);
}

export function formatDate(iso: string): string {
	return dateFmt.format(new Date(iso));
}

export function productLabel(p: Product): string {
	return `${p.name} ${p.brand.name} ${p.measurement.size}${p.measurement.acronym}`;
}
