import { getHousehold } from '$lib/auth';

export async function load({ parent }) {
	const { user } = await parent();
	const household = await getHousehold();
	const isCreator = !!household && household.creator_id === user.id;
	return { user, household, isCreator };
}
