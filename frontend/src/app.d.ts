// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				id: number;
				username: string;
				role: 'super_admin' | 'admin' | 'user';
				wallet: {
					id: number;
					balance: string;
				};
			} | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
