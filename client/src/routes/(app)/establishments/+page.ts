import { listResources } from '$lib/resources';

export async function load() {
	return { items: await listResources('establishments') };
}
