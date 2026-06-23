import { env } from '$env/dynamic/public';

const API_BASE_URL = env.PUBLIC_API_BASE_URL ?? '';

export interface ApiResponse<T = unknown> {
	error: boolean;
	message?: string;
	data?: T;
}

export async function apiFetch<T = unknown>(
	path: string,
	init: RequestInit = {}
): Promise<{ status: number; body: ApiResponse<T> }> {
	const response = await fetch(`${API_BASE_URL}${path}`, {
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...init.headers
		},
		...init
	});

	const body = (await response.json()) as ApiResponse<T>;
	return { status: response.status, body };
}
