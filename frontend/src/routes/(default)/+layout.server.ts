import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// User sudah di-load di hooks.server.ts
	// Tidak perlu validasi ulang di sini
	return {
		user: locals.user || null
	};
};
