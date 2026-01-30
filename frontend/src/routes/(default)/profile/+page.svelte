<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Copy, LogOut } from 'lucide-svelte';
	import {
		user,
		loading,
		error,
		loadProfile,
		logout,
		isAuthenticated,
		token
	} from '$lib/stores/auth';

	// FORCE reload profile setiap kali page mount
	// Jangan rely pada cached data dari layout
	onMount(async () => {
		// Always load fresh profile from API
		const success = await loadProfile();
		if (!success) {
			// Jika gagal load profile, kemungkinan token expired
			logout();
			goto('/login');
		}
	});

	// Fungsi untuk copy ke clipboard
	function copyToClipboard(text: string) {
		navigator.clipboard.writeText(text);
		alert('ID disalin!');
	}

	// Format tanggal bergabung (dari created_at jika ada, atau gunakan default)
	function formatDate(date?: string): string {
		if (!date) return '-';
		return new Date(date).toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}

	// Handle logout
	async function handleLogout() {
		await logout();
		goto('/login');
	}

	// Format role untuk display
	function formatRole(role?: string): string {
		if (!role) return '-';
		return role.replace('_', ' ').replace(/\b\w/g, (c) => c.toUpperCase());
	}
</script>

<div class="profile-container">
	{#if $loading}
		<div class="loading-state">
			<p>Memuat profil...</p>
		</div>
	{:else if $error}
		<div class="error-state">
			<p>{$error}</p>
			<button onclick={() => goto('/login')}>Kembali ke Login</button>
		</div>
	{:else if $user}
		<div class="ios-list-group">
			<div class="list-row">
				<span class="row-label">Username</span>
				<span class="row-value">{$user.username}</span>
			</div>

			<div class="list-row">
				<span class="row-label">User ID</span>
				<button class="row-value copy-btn" onclick={() => copyToClipboard(String($user.id || ''))}>
					{$user.id}
					<Copy size={14} class="icon-copy" />
				</button>
			</div>

			{#if $user.wallet}
				<div class="list-row">
					<span class="row-label">Wallet ID</span>
					<button
						class="row-value copy-btn"
						onclick={() => copyToClipboard(String($user?.wallet?.id || ''))}
					>
						{$user.wallet.id}
						<Copy size={14} class="icon-copy" />
					</button>
				</div>

				<div class="list-row">
					<span class="row-label">Saldo</span>
					<span class="row-value">Rp {Number($user.wallet.balance).toLocaleString('id-ID')}</span>
				</div>
			{/if}

			<div class="list-row">
				<span class="row-label">Role</span>
				<span class="row-value">{formatRole($user.role)}</span>
			</div>

			<div class="list-row">
				<span class="row-label">Status Akun</span>
				<span class="row-value status-active">Aktif</span>
			</div>
		</div>

		<div class="ios-list-group action-group">
			<button class="list-row action" onclick={handleLogout}>
				<div class="row-left">
					<LogOut size={18} class="action-icon logout-icon" />
					<span class="row-label logout-text">Keluar</span>
				</div>
			</button>
		</div>
	{/if}
</div>

<style>
	:global(body) {
		background-color: #ffffff;
		font-family:
			-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
	}

	.profile-container {
		padding: 2rem 1.5rem 6rem 1.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.loading-state,
	.error-state {
		text-align: center;
		padding: 2rem;
	}

	.error-state button {
		margin-top: 1rem;
		padding: 0.5rem 1rem;
		background: #007aff;
		color: white;
		border: none;
		border-radius: 8px;
		cursor: pointer;
	}

	.ios-list-group {
		background-color: #f2f4f6;
		border-radius: 12px;
		overflow: hidden;
		margin-bottom: 1.5rem;
	}

	.list-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1rem;
		background: none;
		border: none;
		width: 100%;
		text-align: left;
		border-bottom: 1px solid #e5e5ea;
	}

	.list-row:last-child {
		border-bottom: none;
	}

	.row-label {
		font-size: 0.95rem;
		font-weight: 600;
		color: #1c1c1e;
	}

	.row-value {
		font-size: 0.95rem;
		color: #636366;
		font-family: inherit;
	}

	.copy-btn {
		background: none;
		border: none;
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 0;
	}

	.status-active {
		color: #34c759;
		font-weight: 600;
	}

	.list-row.action {
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.list-row.action:active {
		background-color: #e5e5ea;
	}

	.row-left {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.logout-text {
		color: #ff3b30;
	}
</style>
