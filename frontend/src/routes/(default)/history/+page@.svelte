<script>
	// @ts-nocheck

	import {
		ArrowUpRight,
		ArrowDownLeft,
		User,
		Smartphone,
		Zap,
		QrCode,
		Wallet,
		ArrowLeft
	} from 'lucide-svelte';

	// --- MOCKUP DATA TRANSAKSI WALLET ---
	const transactions = [
		// BULAN INI (JANUARI 2026)
		{
			id: 1,
			title: 'Budi Santoso',
			desc: 'Bayar makan siang',
			date: '28 Jan 2026',
			type: 'expense',
			amount: 25000,
			category: 'transfer_out'
		},
		{
			id: 2,
			title: 'Top Up BCA',
			desc: 'Via m-Banking',
			date: '27 Jan 2026',
			type: 'income',
			amount: 500000,
			category: 'topup'
		},
		{
			id: 3,
			title: 'Token Listrik',
			desc: 'ID: 1234567890',
			date: '25 Jan 2026',
			type: 'expense',
			amount: 102500,
			category: 'bill'
		},
		{
			id: 4,
			title: 'Siti Aminah',
			desc: 'Patungan kado',
			date: '24 Jan 2026',
			type: 'income',
			amount: 75000,
			category: 'transfer_in'
		},
		{
			id: 5,
			title: 'Kopi Kenangan',
			desc: 'QRIS Payment',
			date: '20 Jan 2026',
			type: 'expense',
			amount: 22000,
			category: 'qris'
		},

		// BULAN LALU (DESEMBER 2025)
		{
			id: 6,
			title: 'Tarik Tunai ATM',
			desc: 'ATM Bersama',
			date: '30 Des 2025',
			type: 'expense',
			amount: 100000,
			category: 'withdraw'
		},
		{
			id: 7,
			title: 'Rizky Febian',
			desc: 'Ganti uang bensin',
			date: '28 Des 2025',
			type: 'expense',
			amount: 50000,
			category: 'transfer_out'
		},
		{
			id: 8,
			title: 'Top Up Mandiri',
			desc: 'Via ATM',
			date: '20 Des 2025',
			type: 'income',
			amount: 1000000,
			category: 'topup'
		}
	];

	// Helper Format Rupiah
	// @ts-ignore
	function formatRupiah(value) {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			maximumFractionDigits: 0
		}).format(value);
	}

	// 1. Hitung Summary Otomatis
	let totalIncome = $derived(
		transactions.filter((t) => t.type === 'income').reduce((sum, t) => sum + t.amount, 0)
	);

	let totalExpense = $derived(
		transactions.filter((t) => t.type === 'expense').reduce((sum, t) => sum + t.amount, 0)
	);

	// 2. Grouping Data per Bulan
	// @ts-ignore
	function groupTransactions(data) {
		const groups = {};
		// @ts-ignore
		data.forEach((t) => {
			const dateParts = t.date.split(' ');
			const monthYear = `${dateParts[1]} ${dateParts[2]}`;
			// @ts-ignore
			if (!groups[monthYear]) groups[monthYear] = [];
			// @ts-ignore
			groups[monthYear].push(t);
		});
		return groups;
	}

	const groupedData = groupTransactions(transactions);
	const months = Object.keys(groupedData);

	// 3. Logic Icon & Warna berdasarkan Kategori Wallet
	// @ts-ignore
	function getIcon(category) {
		switch (category) {
			case 'transfer_out':
				return User; // Kirim ke orang
			case 'transfer_in':
				return User; // Terima dari orang
			case 'topup':
				return Wallet; // Isi Saldo
			case 'bill':
				return Zap; // Bayar Tagihan
			case 'qris':
				return QrCode; // Scan QR
			case 'withdraw':
				return Smartphone; // Tarik Tunai/Pulsa
			default:
				return User;
		}
	}

	// @ts-ignore
	function getColor(category) {
		switch (category) {
			case 'transfer_in':
				return '#34c759'; // Hijau (Masuk)
			case 'topup':
				return '#34c759'; // Hijau (Masuk)
			case 'transfer_out':
				return '#007aff'; // Biru (Transfer)
			case 'bill':
				return '#ff9500'; // Orange (Tagihan)
			case 'qris':
				return '#af52de'; // Ungu (QRIS)
			default:
				return '#8e8e93'; // Abu-abu
		}
	}
</script>

<div class="page-container">
	<header class="header">
		<a href="/wallet" class="back-btn" aria-label="Kembali">
			<ArrowLeft size={24} color="#007aff" />
		</a>

		<h1>Riwayat</h1>

		<div class="spacer"></div>
	</header>

	<div class="summary-card">
		<div class="summary-item">
			<div class="icon-circle income-bg">
				<ArrowDownLeft size={20} color="#34c759" />
			</div>
			<div class="summary-text">
				<span class="label">Uang Masuk</span>
				<span class="value income">{formatRupiah(totalIncome)}</span>
			</div>
		</div>

		<div class="divider"></div>

		<div class="summary-item">
			<div class="icon-circle expense-bg">
				<ArrowUpRight size={20} color="#ff3b30" />
			</div>
			<div class="summary-text">
				<span class="label">Uang Keluar</span>
				<span class="value expense">{formatRupiah(totalExpense)}</span>
			</div>
		</div>
	</div>

	<div class="history-list">
		{#each months as month}
			<div class="month-header">{month.toUpperCase()}</div>

			<div class="ios-list-group">
				{#each groupedData[month] as item}
					{@const ItemIcon = getIcon(item.category)}
					{@const iconColor = getColor(item.category)}

					<div class="list-row">
						<div class="row-left">
							<div class="category-icon" style="background-color: {iconColor}15;">
								<ItemIcon size={18} color={iconColor} />
							</div>
							<div class="item-details">
								<span class="item-title">{item.title}</span>
								<span class="item-desc">{item.desc}</span>
							</div>
						</div>

						<div class="row-right">
							<span class="amount {item.type}">
								{item.type === 'income' ? '+' : '-'}{formatRupiah(item.amount)}
							</span>
							<span class="date-label">{item.date.split(' ')[0]} {item.date.split(' ')[1]}</span>
						</div>
					</div>
				{/each}
			</div>
		{/each}
	</div>
</div>

<style>
	:global(body) {
		background-color: #ffffff; /* Background Putih */
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}

	.page-container {
		padding: 1rem 1.5rem 6rem 1.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	/* HEADER */
	.header {
		margin-bottom: 1.5rem;
		padding-top: 1rem;

		/* Flexbox untuk menata 3 elemen (Btn - Judul - Spacer) */
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.header h1 {
		font-size: 1.2rem; /* Ukuran font disesuaikan agar rapi dgn tombol */
		font-weight: 700;
		margin: 0;
		color: #111;
	}

	/* Style Tombol Back */
	.back-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		display: flex;
		align-items: center;
		justify-content: flex-start;
		width: 40px; /* Lebar area klik */
	}

	/* Style Spacer Penyeimbang */
	.spacer {
		width: 40px; /* Harus sama dengan lebar back-btn */
	}

	/* SUMMARY CARD */
	.summary-card {
		background-color: #ffffff;
		border-radius: 16px;
		padding: 1.25rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		border: 1px solid #f0f0f0;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
		margin-bottom: 2rem;
	}

	.summary-item {
		display: flex;
		align-items: center;
		gap: 12px;
		flex: 1;
	}
	.icon-circle {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.income-bg {
		background-color: #f0fdf4;
	}
	.expense-bg {
		background-color: #fef2f2;
	}

	.summary-text {
		display: flex;
		flex-direction: column;
	}
	.summary-text .label {
		font-size: 0.75rem;
		color: #8e8e93;
		font-weight: 600;
	}
	.summary-text .value {
		font-size: 0.95rem;
		font-weight: 700;
	}
	.value.income {
		color: #34c759;
	}
	.value.expense {
		color: #1c1c1e;
	}
	.divider {
		width: 1px;
		height: 40px;
		background-color: #f0f0f0;
		margin: 0 1rem;
	}

	/* HISTORY LIST */
	.month-header {
		font-size: 0.85rem;
		font-weight: 700;
		color: #9ca3af;
		margin-bottom: 0.75rem;
		padding-left: 0.5rem;
		letter-spacing: 0.5px;
		margin-top: 1.5rem;
	}
	.history-list > .month-header:first-child {
		margin-top: 0;
	}

	.ios-list-group {
		background-color: #ffffff;
		border-radius: 14px;
		overflow: hidden;
		border: 1px solid #f0f0f0;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.02);
	}

	.list-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-bottom: 1px solid #f8f8f8;
	}
	.list-row:last-child {
		border-bottom: none;
	}

	.row-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.category-icon {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.item-details {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}
	.item-title {
		font-size: 0.95rem;
		font-weight: 600;
		color: #1a1a1a;
	}
	.item-desc {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.row-right {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 2px;
	}
	.amount {
		font-size: 0.95rem;
		font-weight: 700;
	}
	.amount.income {
		color: #34c759;
	}
	.amount.expense {
		color: #1a1a1a;
	}
	.date-label {
		font-size: 0.7rem;
		color: #d1d5db;
		font-weight: 500;
	}
</style>
