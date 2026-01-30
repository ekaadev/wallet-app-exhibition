<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { ArrowLeft, Send } from 'lucide-svelte';
	import { token, user, loadProfile } from '$lib/stores/auth';
	import { transfer } from '$lib/api';

	// State untuk input form
	let recipientId = $state('');
	let amount = $state('');
	let note = $state('');
	let isLoading = $state(false);
	let errorMessage = $state('');
	let successMessage = $state('');

	// Check auth saat mount dan baca URL params
	onMount(async () => {
		// Load profile untuk mendapatkan data wallet terbaru
		await loadProfile();

		// Baca URL params dari scan QR
		const urlParams = $page.url.searchParams;
		const userIdParam = urlParams.get('userId');
		const amountParam = urlParams.get('amount');

		// Pre-fill form jika ada params
		if (userIdParam) {
			recipientId = userIdParam;
		}
		if (amountParam) {
			amount = formatRupiahInput(amountParam);
		}
	});

	// Fungsi format rupiah untuk display
	function formatRupiah(value: string | number): string {
		if (!value) return '';
		return new Intl.NumberFormat('id-ID').format(Number(value));
	}

	function formatRupiahInput(value: string): string {
		const numberString = value.replace(/[^,\d]/g, '').toString();
		const split = numberString.split(',');
		let sisa = split[0].length % 3;
		let rupiah = split[0].substr(0, sisa);
		const ribuan = split[0].substr(sisa).match(/\d{3}/gi);

		if (ribuan) {
			const separator = sisa ? '.' : '';
			rupiah += separator + ribuan.join('.');
		}

		rupiah = split[1] != undefined ? rupiah + ',' + split[1] : rupiah;
		return rupiah;
	}

	function handleAmountInput(e: Event) {
		const target = e.target as HTMLInputElement;
		amount = formatRupiahInput(target.value);
	}

	// Handle submit transfer
	async function handleSend() {
		errorMessage = '';
		successMessage = '';

		// Parse amount dari format rupiah (hapus titik)
		const cleanAmount = amount.replace(/\./g, '');

		// Validasi input
		if (!recipientId || !cleanAmount) {
			errorMessage = 'Mohon isi ID penerima dan jumlah transfer';
			return;
		}

		const userId = parseInt(recipientId);
		if (isNaN(userId) || userId <= 0) {
			errorMessage = 'ID penerima harus berupa angka yang valid';
			return;
		}

		const transferAmount = parseFloat(cleanAmount);
		if (isNaN(transferAmount) || transferAmount <= 0) {
			errorMessage = 'Jumlah transfer harus lebih dari 0';
			return;
		}

		// Cek saldo cukup
		if ($user?.wallet) {
			const currentBalance = parseFloat($user.wallet.balance);
			if (transferAmount > currentBalance) {
				errorMessage = 'Saldo tidak cukup';
				return;
			}
		}

		isLoading = true;

		try {
			const response = await transfer({
				to_user_id: userId,
				amount: transferAmount,
				description: note || undefined
			});

			if (response.data) {
				successMessage = `Berhasil kirim Rp${formatRupiah(cleanAmount)} ke User ID ${recipientId}`;
				// Reset form
				recipientId = '';
				amount = '';
				note = '';
				// Reload profile untuk update saldo
				await loadProfile();
			}
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : 'Transfer gagal';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="page-container">
	<header class="header">
		<a href="/wallet" class="back-btn" aria-label="Kembali">
			<ArrowLeft size={24} color="#007aff" />
		</a>
		<h1>Transfer</h1>
		<div class="spacer"></div>
	</header>

	{#if errorMessage}
		<div class="message error-message">
			{errorMessage}
		</div>
	{/if}

	{#if successMessage}
		<div class="message success-message">
			{successMessage}
		</div>
	{/if}

	{#if $user?.wallet}
		<div class="balance-display">
			<span class="balance-label">Saldo Anda:</span>
			<span class="balance-value">Rp{formatRupiah($user.wallet.balance)}</span>
		</div>
	{/if}

	<div class="ios-list-group">
		<div class="list-row">
			<label for="recipient" class="row-label">Ke (User ID)</label>
			<input
				id="recipient"
				type="number"
				class="row-input"
				placeholder="Masukkan User ID"
				bind:value={recipientId}
				disabled={isLoading}
			/>
		</div>

		<div class="list-row">
			<label for="amount" class="row-label">Jumlah</label>
			<div class="amount-wrapper">
				<span class="currency">Rp</span>
				<input
					id="amount"
					type="text"
					inputmode="numeric"
					class="row-input amount-input"
					placeholder="0"
					value={amount}
					oninput={handleAmountInput}
					disabled={isLoading}
				/>
			</div>
		</div>
	</div>

	<div class="ios-list-group">
		<div class="list-row">
			<label for="note" class="row-label">Catatan</label>
			<input
				id="note"
				type="text"
				class="row-input"
				placeholder="Pembayaran makan siang..."
				bind:value={note}
				disabled={isLoading}
			/>
		</div>
	</div>

	<div class="fixed-footer">
		<button class="send-btn {isLoading ? 'loading' : ''}" onclick={handleSend} disabled={isLoading}>
			{#if isLoading}
				Memproses...
			{:else}
				<span class="btn-text">Kirim Sekarang</span>
				<Send size={18} />
			{/if}
		</button>
	</div>
</div>

<style>
	:global(body) {
		background-color: #ffffff;
		font-family:
			-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
	}

	.page-container {
		padding: 1rem 1.5rem;
		max-width: 480px;
		margin: 0 auto;
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1.5rem;
		padding-top: 1rem;
	}

	.header h1 {
		font-size: 1.1rem;
		font-weight: 700;
		margin: 0;
		color: #000;
	}

	.back-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		display: flex;
		align-items: center;
	}

	.spacer {
		width: 24px;
	}

	.message {
		padding: 0.75rem 1rem;
		border-radius: 8px;
		margin-bottom: 1rem;
		font-size: 0.9rem;
	}

	.error-message {
		background-color: #fee2e2;
		color: #dc2626;
	}

	.success-message {
		background-color: #dcfce7;
		color: #16a34a;
	}

	.balance-display {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		background-color: #f0f9ff;
		border-radius: 8px;
		margin-bottom: 1rem;
	}

	.balance-label {
		font-size: 0.9rem;
		color: #666;
	}

	.balance-value {
		font-size: 1rem;
		font-weight: 600;
		color: #007aff;
	}

	.ios-list-group {
		background-color: #f2f4f6;
		border-radius: 12px;
		overflow: hidden;
		margin-bottom: 1.5rem;
	}

	.list-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 1rem;
		border-bottom: 1px solid #e5e5ea;
		background: transparent;
	}

	.list-row:last-child {
		border-bottom: none;
	}

	.row-label {
		font-size: 0.95rem;
		font-weight: 600;
		color: #1c1c1e;
		width: 100px;
		flex-shrink: 0;
	}

	.row-input {
		border: none;
		background: transparent;
		font-size: 0.95rem;
		color: #000;
		width: 100%;
		text-align: right;
		font-family: inherit;
		outline: none;
	}

	.row-input::placeholder {
		color: #c7c7cc;
	}

	.row-input:disabled {
		opacity: 0.6;
	}

	.amount-wrapper {
		display: flex;
		align-items: center;
		justify-content: flex-end;
		width: 100%;
		gap: 2px;
	}

	.currency {
		font-weight: 600;
		color: #1c1c1e;
	}

	.amount-input {
		font-weight: 600;
		color: #007aff;
	}

	/* input[type='number'] removed */

	.fixed-footer {
		position: fixed;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		width: 100%;
		max-width: 480px;
		background: rgba(255, 255, 255, 0.9);
		backdrop-filter: blur(10px);
		border-top: 1px solid rgba(0, 0, 0, 0.05);
		padding: 1rem 1.5rem 2rem 1.5rem;
		z-index: 100;
	}

	.send-btn {
		background-color: #1c1c1e;
		color: white;
		width: 100%;
		padding: 1rem;
		border: none;
		border-radius: 50px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
	}

	.send-btn:active {
		transform: scale(0.98);
	}

	.send-btn:disabled {
		background-color: #8e8e93;
		cursor: not-allowed;
	}
</style>
