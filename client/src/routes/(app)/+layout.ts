import { redirect } from '@sveltejs/kit';
import { isAuthenticated } from '$lib/auth';

export async function load() {
	if (!(await isAuthenticated())) redirect(307, '/login');
}
