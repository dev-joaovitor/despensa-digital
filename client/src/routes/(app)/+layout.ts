import { redirect } from '@sveltejs/kit';
import { getUser } from '$lib/auth';

export async function load() {
	const user = await getUser();
	if (!user) redirect(307, '/login');
	return { user };
}
