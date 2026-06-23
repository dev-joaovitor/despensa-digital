import { apiFetch } from '$lib/api';

export interface SessionUser {
	id: number;
	full_name: string;
	email: string;
}

export interface Household {
	id: number;
	name: string;
	creator_id: number;
	// Empty string when the session user is not the creator.
	invitation_code: string;
	created_at: string;
	updated_at: string;
}

export async function getUser(): Promise<SessionUser | null> {
	try {
		const { status, body } = await apiFetch<SessionUser>('/api/v1/auth/me');
		return status === 200 && body.data ? body.data : null;
	} catch {
		return null;
	}
}

export async function getHousehold(): Promise<Household | null> {
	try {
		const { status, body } = await apiFetch<Household>('/api/v1/households/');
		return status === 200 && body.data ? body.data : null;
	} catch {
		return null;
	}
}

export async function isAuthenticated(): Promise<boolean> {
	return (await getUser()) !== null;
}
