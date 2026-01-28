import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const token = cookies.get('jwt');
	
	// Jika ada token, redirect ke wallet
	if (token) {
		throw redirect(302, '/wallet');
	}
	
	// Jika tidak ada token, redirect ke login
	throw redirect(302, '/login');
};
