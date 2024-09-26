import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { QrCode } from '$lib/types';


export const load: PageServerLoad<{ code: QrCode | null }> = async ({ params }: { params: { id: string } }) => {
	try {
		const resp = await fetch(`http://localhost:8080/api/get?id=${params.id}`,
			{
				method: 'GET',
				headers: {
					Authorization: 'Bearer 1234567890'
				}
			});

		if (!resp.ok) {
			if (resp.status === 404) error(404, 'Uh oh! Could not find that code');
			error(500, 'Uh oh! Something went wrong');
		}

		const code = await resp.json() as QrCode;
		return {
			code
		}
	} catch (err) {
		error(err.status, err.body && err.body.message);
	}
};
