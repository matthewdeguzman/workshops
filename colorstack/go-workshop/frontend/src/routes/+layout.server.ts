import type { LayoutServerLoad } from './$types';
import type { QrCode } from '$lib/types';


export const load: LayoutServerLoad<{ codes: QrCode[] | null }> = async () => {
	let codes = null;
	try {
		const resp = await fetch(`http://localhost:8080/api/get?all=true`,
			{
				method: 'GET',
				headers: {
					Authorization: 'Bearer 1234567890'
				}
			});

		if (!resp.ok) {
			throw new Error(resp.statusText);
		}

		codes = (await resp.json()) as QrCode[];
	} catch (err) {
		console.error(err);
	}
	return {
		codes
	}
};
