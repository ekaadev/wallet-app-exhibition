import { redirect } from '@sveltejs/kit';
import { PUBLIC_API_BASE_URL } from '$env/static/public';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('jwt');
	const API_BASE_URL = PUBLIC_API_BASE_URL || 'http://localhost:3000';

	if (!token) {
		throw redirect(302, '/login');
	}

	const res = await fetch(`${API_BASE_URL}/users/me`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${token}`,
			Cookie: `jwt=${token}`
		},
		credentials: 'include'
	});

	if (!res.ok) {
		// Jika token invalid atau expired, redirect ke login
		throw redirect(302, '/login');
	}

	const responseData = await res.json();

	return {
		user: responseData.data
	};
};
