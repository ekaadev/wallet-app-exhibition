<script lang="ts">
	import NavMenu from '$lib/components/nav-menu.svelte';
	let { children } = $props();
</script>

<svelte:head>
	<title>Wallet Exhibition App</title>
</svelte:head>

<div>
	<NavMenu />
	{@render children()}
</div>
