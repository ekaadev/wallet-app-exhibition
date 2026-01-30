import { redirect, type Handle } from '@sveltejs/kit';
import { PUBLIC_API_BASE_URL } from '$env/static/public';

const API_BASE_URL = PUBLIC_API_BASE_URL || 'http://localhost:3000';

// Routes yang memerlukan authentication
const protectedRoutes = ['/wallet', '/history', '/profile', '/transaction'];

// Routes yang tidak boleh diakses jika sudah login
const authRoutes = ['/login', '/signup'];

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('jwt');
	const pathname = event.url.pathname;

	// Always clear locals.user first
	// This prevents caching old user data
	event.locals.user = null;

	// Check if route is protected
	const isProtectedRoute = protectedRoutes.some(route => pathname.startsWith(route));
	const isAuthRoute = authRoutes.some(route => pathname.startsWith(route));

	// If accessing protected route without token
	if (isProtectedRoute && !token) {
		throw redirect(302, '/login');
	}

	// If accessing auth route with valid token, redirect to wallet
	if (isAuthRoute && token) {
		// Verify token is still valid
		try {
			const res = await fetch(`${API_BASE_URL}/users/me`, {
				headers: {
					'Authorization': `Bearer ${token}`,
				},
				credentials: 'include'
			});

			if (res.ok) {
				throw redirect(302, '/wallet');
			}
		} catch (error) {
			// If verification fails, continue to auth route
			if (error instanceof Response && error.status === 302) {
				throw error;
			}
			// If API error, clear invalid cookie
			event.cookies.delete('jwt', { path: '/' });
		}
	}

	// Add user to locals if token exists
	if (token) {
		try {
			const res = await fetch(`${API_BASE_URL}/users/me`, {
				headers: {
					'Authorization': `Bearer ${token}`,
				},
				credentials: 'include'
			});

			if (res.ok) {
				const data = await res.json();
				// Set fresh user data
				event.locals.user = data.data;
			} else {
				// Invalid token, clear cookie AND locals
				event.cookies.delete('jwt', { path: '/' });
				event.locals.user = null;
			}
		} catch (error) {
			console.error('Error validating token:', error);
			// On error, clear everything
			event.cookies.delete('jwt', { path: '/' });
			event.locals.user = null;
		}
	}

	// Add cache control headers to prevent stale data
	const response = await resolve(event);
	
	// Don't cache authenticated pages
	if (isProtectedRoute) {
		response.headers.set('Cache-Control', 'no-store, no-cache, must-revalidate');
	}
	
	return response;
};