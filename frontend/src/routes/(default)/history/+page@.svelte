<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { ArrowUpRight, ArrowDownLeft, User, Wallet, ArrowLeft } from 'lucide-svelte';
	import { token, user } from '$lib/stores/auth';
	import { getWalletMutations, type WalletMutationResponse } from '$lib/api/wallet';

	// State
	let mutations: WalletMutationResponse[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	// Load data saat mount
	onMount(async () => {
		if (!$token) {
			goto('/login');
			return;
		}

		try {
			const response = await getWalletMutations($token, 1, 50); // Get 50 transaksi terakhir
			if (response.data) {
				mutations = response.data.mutations;
			}
		} catch (err) {
			error = 'Gagal memuat riwayat transaksi';
			console.error(err);
		} finally {
			loading = false;
		}
	});

	// Helper Format Rupiah
	function formatRupiah(value: number | string): string {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			maximumFractionDigits: 0
		}).format(Number(value));
	}

	// Helper Format Tanggal
	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}

	// Menghitung Summary
	let totalIncome = $derived(
		mutations.filter((m) => m.type === 'credit').reduce((sum, m) => sum + Number(m.amount), 0)
	);

	let totalExpense = $derived(
		mutations.filter((m) => m.type === 'debit').reduce((sum, m) => sum + Number(m.amount), 0)
	);

	// Grouping Data per Bulan
	function groupMutations(data: WalletMutationResponse[]) {
		const groups: Record<string, WalletMutationResponse[]> = {};

		data.forEach((m) => {
			const date = new Date(m.created_at);
			const monthYear = date.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' });

			if (!groups[monthYear]) groups[monthYear] = [];
			groups[monthYear].push(m);
		});

		return groups;
	}

	let groupedData = $derived(groupMutations(mutations));
	let months = $derived(Object.keys(groupedData));

	// Logic Icon
	function getIcon(type: 'credit' | 'debit') {
		return type === 'credit' ? Wallet : User;
	}

	// Logic Warna
	function getColor(type: 'credit' | 'debit') {
		return type === 'credit' ? '#34c759' : '#1a1a1a'; // Hijau : Hitam
	}

	// Logic Title
	function getTitle(type: 'credit' | 'debit'): string {
		return type === 'credit' ? 'Uang Masuk' : 'Uang Keluar';
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

	{#if loading}
		<div class="loading-state">Memuat riwayat...</div>
	{:else if error}
		<div class="error-state">{error}</div>
	{:else}
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

		{#if mutations.length === 0}
			<div class="empty-state">
				<p>Belum ada transaksi</p>
			</div>
		{:else}
			<div class="history-list">
				{#each months as month}
					<div class="month-header">{month.toUpperCase()}</div>

					<div class="ios-list-group">
						{#each groupedData[month] as item}
							{@const ItemIcon = getIcon(item.type)}
							{@const iconColor = getColor(item.type)}

							<div class="list-row">
								<div class="row-left">
									<div
										class="category-icon"
										style="background-color: {item.type === 'credit' ? '#f0fdf4' : '#f3f4f6'};"
									>
										<ItemIcon size={18} color={iconColor} />
									</div>
									<div class="item-details">
										<span class="item-title">{getTitle(item.type)}</span>
										<span class="item-desc"
											>Ref: {item.transaction_id} • Sisa: {formatRupiah(item.balance_after)}</span
										>
									</div>
								</div>

								<div class="row-right">
									<span class="amount {item.type}">
										{item.type === 'credit' ? '+' : '-'}{formatRupiah(item.amount)}
									</span>
									<span class="date-label">{formatDate(item.created_at)}</span>
								</div>
							</div>
						{/each}
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>

<style>
	:global(body) {
		background-color: #ffffff;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}

	.page-container {
		padding: 1rem 1.5rem 6rem 1.5rem;
		max-width: 480px;
		margin: 0 auto;
	}

	.loading-state,
	.error-state,
	.empty-state {
		text-align: center;
		padding: 2rem;
		color: #8e8e93;
	}

	/* HEADER */
	.header {
		margin-bottom: 1.5rem;
		padding-top: 1rem;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.header h1 {
		font-size: 1.2rem;
		font-weight: 700;
		margin: 0;
		color: #111;
	}

	.back-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		display: flex;
		align-items: center;
		justify-content: flex-start;
		width: 40px;
	}

	.spacer {
		width: 40px;
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
	.amount.credit {
		color: #34c759;
	}
	.amount.debit {
		color: #1a1a1a;
	}
	.date-label {
		font-size: 0.7rem;
		color: #d1d5db;
		font-weight: 500;
	}
</style>
