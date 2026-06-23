import { apiFetch } from '$lib/api';

export async function isAuthenticated(): Promise<boolean> {
	try {
		const { status } = await apiFetch('/api/v1/auth/me');
		return status === 200;
	} catch {
		return false;
	}
}
