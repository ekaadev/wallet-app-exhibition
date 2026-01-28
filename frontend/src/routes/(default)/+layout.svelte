<script lang="ts">
	import NavMenu from '$lib/components/nav-menu.svelte';
	import { user, token } from '$lib/stores/auth';
	import { connectWebSocket } from '$lib/api/websocket';
	import { browser } from '$app/environment';

	let { children, data } = $props();

	$effect(() => {
		if (data.user) {
			user.set(data.user);
			// Also set token if it was available in cookies,
			// though layout.server.ts reads it.
			// Ideally we should pass token too if we want full store sync,
			// but user object is what UI needs.

			// Auto connect WS if user exists
			if (browser) {
				connectWebSocket();
			}
		}
	});
</script>

<svelte:head>
	<title>Wallet Exhibition App</title>
</svelte:head>

<div>
	<NavMenu />
	{@render children()}
</div>
