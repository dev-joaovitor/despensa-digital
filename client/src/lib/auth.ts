import { apiFetch } from '$lib/api';

export interface SessionUser {
	id: number;
	full_name: string;
	email: string;
}

export async function getUser(): Promise<SessionUser | null> {
	try {
		const { status, body } = await apiFetch<SessionUser>('/api/v1/auth/me');
		return status === 200 && body.data ? body.data : null;
	} catch {
		return null;
	}
}

export async function isAuthenticated(): Promise<boolean> {
	return (await getUser()) !== null;
}
