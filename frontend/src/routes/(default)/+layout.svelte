<script lang="ts">
	import NavMenu from '$lib/components/nav-menu.svelte';
	import { user } from '$lib/stores/auth';
	import { connectWebSocket, disconnectWebSocket } from '$lib/api/websocket';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';

	let { children, data } = $props();

	// Sync user data dari server ke client store
	$effect(() => {
		if (data.user) {
			user.set(data.user);
		}
	});

	// WebSocket connection management
	onMount(() => {
		if (browser && data.user) {
			// Connect WebSocket saat component mount
			connectWebSocket();

			// Cleanup saat component unmount
			return () => {
				disconnectWebSocket();
			};
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
