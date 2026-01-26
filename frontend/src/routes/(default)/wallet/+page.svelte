<script>
	import { goto } from '$app/navigation';
	import { History, Banknote, Scan } from 'lucide-svelte';

	const userData = {
		name: 'Yohanes Eka Putra',
		balance: 'Rp235'
	};

	// Fungsi placeholder untuk navigasi/aksi
	function handleTransfer() {
        goto('/transaction');
	}

	function handleScan() {
		console.log('Buka Kamera untuk Scan QRIS');
	}
</script>

<div class="wallet-container">
	<header class="header">
		<h1>Wallet</h1>
		<div class="actions">
            <a href="/history">
			    <button class="icon-btn" aria-label="Riwayat"><History size={24} /></button>
            </a>
		</div>
	</header>

	<p class="user-label">{userData.name}</p>

	<div class="card-wrapper">
		<div class="jago-card-clean">
			<div class="card-top">
				<span class="visa-logo">EXHIBITION</span>
			</div>

			<div class="card-bottom">
				<div class="card-details">
					<span class="card-name">{userData.name.toUpperCase()}</span>
				</div>
			</div>

			<div class="card-notch"></div>
		</div>
	</div>

	<div class="balance-info">
		<div class="pocket-label">
			<span>Balance</span>
		</div>
		<span class="amount">{userData.balance}</span>
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
</div>

<style>
	:global(body) {
		margin: 0;
		font-family: 'Inter', sans-serif;
		background-color: #ffffff;
	}

	.wallet-container {
		padding: 1.5rem;
		padding-bottom: 6rem; /* Ruang untuk navbar */
		max-width: 480px;
		margin: 0 auto;
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
		background: #fff8e6; /* Kuning sangat muda */
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
		aspect-ratio: 3.5 / 5; /* Rasio Portrait Kartu */

		/* Warna Solid Tanpa Gradasi Aneh */
		background-color: #fdb92e;

		border-radius: 20px;
		position: relative;
		padding: 24px;
		box-sizing: border-box;

		/* Layout Flexbox */
		display: flex;
		flex-direction: column;
		justify-content: space-between;

		/* Shadow tipis agar tidak flat mati */
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
		margin-bottom: 15px; /* Jarak dari notch */
	}

	.card-name {
		display: block;
		color: white;
		font-weight: 600;
		font-size: 1rem;
		margin-bottom: 0.5rem;
		letter-spacing: 0.5px;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1); /* Sedikit shadow biar teks jelas */
	}

	/* (Lengkungan bawah) */
	.card-notch {
		position: absolute;
		bottom: -30px;
		left: 50%;
		transform: translateX(-50%);
		width: 60px;
		height: 60px;
		background-color: #fff; /* Harus sama dengan warna background body */
		border-radius: 50%;
		z-index: 10;
	}

	/* Info Saldo Bawah */
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
        gap: 1rem; /* Jarak antar tombol */
        width: 100%; /* Pastikan container mengambil lebar penuh */
        margin-top: 2rem;
    }

    .action-btn {
        flex: 1 1 0; 
        
        background-color: #ffffff;
        border: 1px solid #f0f0f0;
        border-radius: 100px;
        padding: 1rem 0.5rem; /* Padding kiri-kanan dikurangi sedikit jaga-jaga layar kecil */
        
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem; /* Gap icon dan text dirapatkan sedikit */
        
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
        cursor: pointer;
        transition: transform 0.1s ease, box-shadow 0.2s ease;
        
        /* Opsional: Mencegah tombol jadi gepeng kalau layar sangat kecil */
        min-width: 0; 
    }

    /* Text Style */
    .action-text {
        font-weight: 700;
        font-size: 0.85rem; /* Font dikecilkan sedikit agar muat 50:50 di HP kecil */
        color: #1a1a1a;
        white-space: nowrap; /* Mencegah teks turun baris */
        overflow: hidden;
        text-overflow: ellipsis; /* Jika layar kekecilan, teks akan jadi ... */
    }

    /* Icon Styling */
    .icon-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0; /* Mencegah icon tergencet */
    }
	/* Warna Icon sesuai referensi (Orange & Ungu) */
	.transfer-icon {
		color: #f59e0b; /* Orange-ish */
	}

	.scan-icon {
		color: #9333ea; /* Purple-ish */
	}
</style>
