<script>
	import { page } from '$app/stores';
	import { WalletCards, User } from 'lucide-svelte';

	const menuItems = [
		{ name: 'Wallet', icon: WalletCards, link: '/wallet' },
		{ name: 'Profile', icon: User, link: '/profile' }
	];
</script>

<nav class="glass-navbar">
	{#each menuItems as item}
		{@const Icon = item.icon}
		{@const isActive = $page.url.pathname.startsWith(item.link)}

		<a href={item.link} class="nav-item {isActive ? 'active' : ''}" aria-label={item.name}>
			<Icon size={20} />
			<span class="label">{item.name}</span>
		</a>
	{/each}
</nav>

<style>
	/* ... Global body style biarkan sama ... */
	:global(body) {
		background-color: #f8f9fa;
		color: #333;
		margin: 0;
		min-height: 100vh;
	}

	.glass-navbar {
		position: fixed;
		/* Turunkan sedikit posisinya biar tidak terlalu makan tempat di layar */
		bottom: 1.5rem;
		left: 50%;
		transform: translateX(-50%);
		z-index: 1000;

		display: flex;
		align-items: center;
		justify-content: center; /* Center content */

		gap: 1.5rem; /* Jarak antar menu lebih rapat */
		padding: 0.6rem 1.5rem; /* Padding atas/bawah & kiri/kanan diperkecil */

		/* Glassmorphism Styles (Tetap) */
		background: rgba(255, 255, 255, 0.75); /* Opacity dinaikkan dikit biar tulisan jelas */
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(0, 0, 0, 0.05);
		box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05); /* Shadow lebih halus */
		border-radius: 50px;
		color: #333;
	}

	.nav-item {
		background: none;
		border: none;
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;
		/* Jarak antara icon dan text didempetkan */
		gap: 0.2rem;

		text-decoration: none; /* Penting untuk tag <a> */
		color: #94a3b8;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

		/* Padding item diperkecil */
		padding: 0.25rem 0.5rem;
	}

	.label {
		/* Ukuran font diperkecil agar pas dengan icon size 20 */
		font-size: 0.65rem;
		font-weight: 600;
		letter-spacing: 0.01em;
	}

	/* Hover */
	.nav-item:hover {
		color: #0f172a;
		transform: translateY(-2px);
	}

	/* Active State */
	.nav-item.active {
		color: #000000;
		transform: translateY(-2px);
	}

	.nav-item.active :global(svg) {
		filter: drop-shadow(0 2px 3px rgba(0, 0, 0, 0.15));
		stroke-width: 2.5px; /* Sedikit ditebalkan icon-nya saat aktif */
	}
</style>
