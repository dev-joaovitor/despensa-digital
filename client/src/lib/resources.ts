import { apiFetch } from '$lib/api';

export interface NamedResource {
	id: number;
	household_id?: number;
	name: string;
	created_at: string;
	updated_at: string;
}

export type ResourceKind = 'establishments' | 'brands' | 'categories';

interface ResourceConfig {
	title: string;
	singular: string;
	gender: 'm' | 'f';
}

export const RESOURCE_CONFIG: Record<ResourceKind, ResourceConfig> = {
	establishments: { title: 'Estabelecimentos', singular: 'estabelecimento', gender: 'm' },
	brands: { title: 'Marcas', singular: 'marca', gender: 'f' },
	categories: { title: 'Categorias', singular: 'categoria', gender: 'f' }
};

interface MutationResult {
	ok: boolean;
	data?: NamedResource;
	message?: string;
}

export async function listResources(kind: ResourceKind): Promise<NamedResource[]> {
	const { status, body } = await apiFetch<NamedResource[]>(`/api/v1/${kind}`);
	return status === 200 && body.data ? body.data : [];
}

export async function createResource(kind: ResourceKind, name: string): Promise<MutationResult> {
	const { status, body } = await apiFetch<NamedResource>(`/api/v1/${kind}`, {
		method: 'POST',
		body: JSON.stringify({ name })
	});
	const ok = status === 201 && !body.error;
	return { ok, data: body.data, message: body.message };
}

export async function updateResource(
	kind: ResourceKind,
	id: number,
	name: string
): Promise<MutationResult> {
	const { status, body } = await apiFetch<NamedResource>(`/api/v1/${kind}/${id}`, {
		method: 'PATCH',
		body: JSON.stringify({ name })
	});
	// 304 = nothing changed; still a success from the user's view.
	const ok = (status === 200 || status === 304) && !body.error;
	return { ok, data: body.data, message: body.message };
}

export async function deleteResource(
	kind: ResourceKind,
	id: number
): Promise<{ ok: boolean; message?: string }> {
	const { status, body } = await apiFetch(`/api/v1/${kind}/${id}`, { method: 'DELETE' });
	const ok = status === 200 && !body.error;
	return { ok, message: body.message };
}
