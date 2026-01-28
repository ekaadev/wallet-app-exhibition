import { redirect } from '@sveltejs/kit';

export function load({ cookies }) {
	const token = cookies.get('jwt');
	
	if (!token) {
		throw redirect(302, '/login');
	}
}
