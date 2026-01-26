<script>
	import { Copy, ChevronRight, LogOut, ShieldCheck, User } from 'lucide-svelte';

	// Data Mockup
	const profileData = {
		username: 'Yohanes Eka Putra',
		email: 'ekaaputra223@gmail.com', // Saya tambah email biar header lebih cantik
		role: 'Super Admin',
		wallet_id: '1017 1913 7671',
		created_at: '20 Jan 2024',
		status: 'Aktif'
	};

	/**
	 * @param {string} text
	 */
	function copyToClipboard(text) {
		navigator.clipboard.writeText(text);
		// Bisa tambah toast notification disini
		alert('ID disalin!');
	}
</script>

<div class="profile-container">
	<div class="ios-list-group">
		<div class="list-row">
			<span class="row-label">Wallet ID</span>
			<button class="row-value copy-btn" onclick={() => copyToClipboard(profileData.wallet_id)}>
				{profileData.wallet_id}
				<Copy size={14} class="icon-copy" />
			</button>
		</div>

		<div class="list-row">
			<span class="row-label">Role</span>
			<span class="row-value">{profileData.role}</span>
		</div>

		<div class="list-row">
			<span class="row-label">Status Akun</span>
			<span class="row-value status-active">{profileData.status}</span>
		</div>

		<div class="list-row">
			<span class="row-label">Bergabung</span>
			<span class="row-value">{profileData.created_at}</span>
		</div>
	</div>

	<div class="ios-list-group action-group">
		<button class="list-row action">
			<div class="row-left">
				<LogOut size={18} class="action-icon logout-icon" />
				<span class="row-label logout-text">Keluar</span>
			</div>
		</button>
	</div>
</div>

<style>
	:global(body) {
		background-color: #ffffff;
		font-family:
			-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; /* Font System ala Apple */
	}

	.profile-container {
		padding: 2rem 1.5rem 6rem 1.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	/* --- IOS LIST GROUP STYLE --- */
	.ios-list-group {
		background-color: #f2f4f6; /* Warna latar kartu abu-abu muda ala referensi */
		border-radius: 12px;
		overflow: hidden; /* Agar border-radius memotong child elemen */
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

		/* Garis pemisah antar item */
		border-bottom: 1px solid #e5e5ea;
	}

	/* Menghilangkan garis di item terakhir */
	.list-row:last-child {
		border-bottom: none;
	}

	.row-label {
		font-size: 0.95rem;
		font-weight: 600;
		color: #1c1c1e; /* Hampir hitam */
	}

	.row-value {
		font-size: 0.95rem;
		color: #636366; /* Abu-abu teks sekunder */
		font-family: inherit; /* Untuk button */
	}

	/* --- SPECIFIC STYLES --- */
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
		color: #34c759; /* Apple Green */
		font-weight: 600;
	}

	/* Action Rows */
	.list-row.action {
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.list-row.action:active {
		background-color: #e5e5ea; /* Efek tekan */
	}

	.row-left {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.logout-text {
		color: #ff3b30; /* Apple Red */
	}
</style>
