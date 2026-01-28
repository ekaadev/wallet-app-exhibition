<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { History, Banknote, Scan } from 'lucide-svelte';
	import { user, token, loadProfile, logout } from '$lib/stores/auth';

	// State untuk wallet data
	let loading = $state(true);
	let walletError = $state('');

	// Load profile/wallet data saat mount
	onMount(async () => {
		// Load profile dari API (termasuk wallet info)
		const success = await loadProfile();
		if (!success) {
			// Jika gagal, kemungkinan token expired
			walletError = 'Gagal memuat data wallet';
			await logout();
			return;
		}

		loading = false;
	});

	// Format saldo ke format Rupiah
	function formatBalance(balance?: string): string {
		if (!balance) return 'Rp0';
		const num = Number(balance);
		return `Rp${num.toLocaleString('id-ID')}`;
	}

	// Fungsi untuk navigasi/aksi
	function handleTransfer() {
		goto('/transaction');
	}

	function handleScan() {
		console.log('Buka Kamera untuk Scan QRIS');
	}
</script>

<div class="wallet-container">
	{#if loading}
		<div class="loading-state">
			<p>Memuat wallet...</p>
		</div>
	{:else if walletError}
		<div class="error-state">
			<p>{walletError}</p>
		</div>
	{:else if $user}
		<header class="header">
			<h1>Wallet</h1>
			<div class="actions">
				<a href="/history">
					<button class="icon-btn" aria-label="Riwayat"><History size={24} /></button>
				</a>
			</div>
		</header>

		<p class="user-label">{$user.username}</p>

		<div class="card-wrapper">
			<div class="jago-card-clean">
				<div class="card-top">
					<span class="visa-logo">EXHIBITION</span>
				</div>

				<div class="card-bottom">
					<div class="card-details">
						<span class="card-name">{$user.username.toUpperCase()}</span>
					</div>
				</div>

				<div class="card-notch"></div>
			</div>
		</div>

		<div class="balance-info">
			<div class="pocket-label">
				<span>Balance</span>
			</div>
			<span class="amount">{formatBalance($user.wallet?.balance)}</span>
		</div>

		<div class="quick-actions">
			<button class="action-btn" onclick={handleTransfer}>
				<div class="icon-wrapper transfer-icon">
					<Banknote size={24} />
				</div>
				<span class="action-text">Transfer & Bayar</span>
			</button>

			<button class="action-btn" onclick={handleScan}>
				<div class="icon-wrapper scan-icon">
					<Scan size={24} />
				</div>
				<span class="action-text">Scan QR</span>
			</button>
		</div>
	{/if}
</div>

<style>
	:global(body) {
		margin: 0;
		font-family: 'Inter', sans-serif;
		background-color: #ffffff;
	}

	.wallet-container {
		padding: 1.5rem;
		padding-bottom: 6rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.loading-state,
	.error-state {
		text-align: center;
		padding: 2rem;
		color: #666;
	}

	/* Header */
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	h1 {
		font-size: 1.5rem;
		font-weight: 800;
		margin: 0;
		color: #111;
	}

	.actions {
		display: flex;
		gap: 0.75rem;
	}

	.icon-btn {
		background: #fff8e6;
		border: none;
		width: 36px;
		height: 36px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #333;
		cursor: pointer;
	}

	.user-label {
		color: #888;
		font-weight: 600;
		font-size: 0.9rem;
		margin-bottom: 0.75rem;
		padding-left: 0.25rem;
	}

	.card-wrapper {
		display: flex;
		justify-content: center;
		margin-bottom: 1rem;
	}

	.jago-card-clean {
		width: 100%;
		aspect-ratio: 3.5 / 5;
		background-color: #fdb92e;
		border-radius: 20px;
		position: relative;
		padding: 24px;
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		box-shadow: 0 4px 12px rgba(253, 185, 46, 0.3);
	}

	.visa-logo {
		font-style: italic;
		font-weight: 900;
		font-size: 1.4rem;
		color: rgba(255, 255, 255, 0.9);
		letter-spacing: 1px;
	}

	.card-bottom {
		margin-bottom: 15px;
	}

	.card-name {
		display: block;
		color: white;
		font-weight: 600;
		font-size: 1rem;
		margin-bottom: 0.5rem;
		letter-spacing: 0.5px;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
	}

	.card-notch {
		position: absolute;
		bottom: -30px;
		left: 50%;
		transform: translateX(-50%);
		width: 60px;
		height: 60px;
		background-color: #fff;
		border-radius: 50%;
		z-index: 10;
	}

	.balance-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 0.5rem;
	}

	.pocket-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #666;
		font-weight: 500;
		font-size: 0.95rem;
	}

	.amount {
		font-size: 1.1rem;
		font-weight: 700;
		color: #111;
	}

	.quick-actions {
		display: flex;
		gap: 1rem;
		width: 100%;
		margin-top: 2rem;
	}

	.action-btn {
		flex: 1 1 0;
		background-color: #ffffff;
		border: 1px solid #f0f0f0;
		border-radius: 100px;
		padding: 1rem 0.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
		cursor: pointer;
		transition:
			transform 0.1s ease,
			box-shadow 0.2s ease;
		min-width: 0;
	}

	.action-text {
		font-weight: 700;
		font-size: 0.85rem;
		color: #1a1a1a;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.icon-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.transfer-icon {
		color: #f59e0b;
	}

	.scan-icon {
		color: #9333ea;
	}
</style>
