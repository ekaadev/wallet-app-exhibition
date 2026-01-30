<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { X } from 'lucide-svelte';
	import QRCode from 'qrcode';
	import { user } from '$lib/stores/auth';

	// Props
	let { show = $bindable(false), onClose = () => {} }: { show: boolean; onClose?: () => void } =
		$props();

	// State
	let amount = $state('');
	let qrCanvas: HTMLCanvasElement | null = $state(null);
	let qrError = $state('');

	// Reactive: generate QR saat amount berubah atau modal dibuka
	$effect(() => {
		if (show && qrCanvas && $user) {
			generateQR();
		}
	});

	// Helper format rupiah input
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
		generateQR();
	}

	// Generate QR Code dengan data user
	async function generateQR() {
		if (!qrCanvas || !$user) return;

		qrError = '';

		try {
			// Clean amount (remove dots)
			const cleanAmount = amount ? amount.replace(/\./g, '') : '';

			// Data untuk QR Code berisi userId dan amount (jika ada)
			const qrData = {
				userId: $user.id,
				amount: cleanAmount ? parseFloat(cleanAmount) : null
			};

			await QRCode.toCanvas(qrCanvas, JSON.stringify(qrData), {
				width: 250,
				margin: 2,
				color: {
					dark: '#1c1c1e',
					light: '#ffffff'
				}
			});
		} catch (err) {
			qrError = 'Gagal generate QR Code';
			console.error('QR Generation error:', err);
		}
	}

	// Handle close modal
	function handleClose() {
		show = false;
		amount = '';
		onClose();
	}

	// Handle klik di luar modal untuk close
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			handleClose();
		}
	}

	// Handle keyboard escape
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			handleClose();
		}
	}

	onMount(() => {
		if (typeof window !== 'undefined') {
			window.addEventListener('keydown', handleKeydown);
		}
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			window.removeEventListener('keydown', handleKeydown);
		}
	});
</script>

{#if show}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="modal-backdrop" onclick={handleBackdropClick}>
		<div class="modal-content">
			<div class="modal-header">
				<h2>QR Code Saya</h2>
				<button class="close-btn" onclick={handleClose} aria-label="Tutup">
					<X size={24} />
				</button>
			</div>

			<div class="qr-container">
				{#if qrError}
					<p class="error">{qrError}</p>
				{/if}
				<canvas bind:this={qrCanvas}></canvas>
				<p class="user-info">
					{$user?.username ?? 'User'} (ID: {$user?.id ?? '-'})
				</p>
			</div>

			<div class="amount-section">
				<label for="qr-amount">Jumlah yang diminta (opsional)</label>
				<div class="amount-input-wrapper">
					<span class="currency">Rp</span>
					<input
						id="qr-amount"
						type="text"
						inputmode="numeric"
						placeholder="0"
						value={amount}
						oninput={handleAmountInput}
					/>
				</div>
				<p class="hint">Kosongkan jika ingin penerima memasukkan jumlah sendiri</p>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal-content {
		background-color: #ffffff;
		border-radius: 20px;
		padding: 1.5rem;
		width: 100%;
		max-width: 360px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
		animation: slideUp 0.2s ease-out;
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.2rem;
		font-weight: 700;
		color: #1c1c1e;
	}

	.close-btn {
		background: none;
		border: none;
		cursor: pointer;
		color: #8e8e93;
		padding: 0;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.close-btn:hover {
		color: #1c1c1e;
	}

	.qr-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 1rem;
		background-color: #f8f8f8;
		border-radius: 12px;
		margin-bottom: 1.5rem;
	}

	.qr-container canvas {
		border-radius: 8px;
	}

	.user-info {
		margin-top: 0.75rem;
		font-size: 0.9rem;
		color: #666;
		font-weight: 500;
	}

	.error {
		color: #dc2626;
		font-size: 0.85rem;
		margin-bottom: 0.5rem;
	}

	.amount-section {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.amount-section label {
		font-size: 0.9rem;
		font-weight: 600;
		color: #1c1c1e;
	}

	.amount-input-wrapper {
		display: flex;
		align-items: center;
		background-color: #f2f4f6;
		border-radius: 10px;
		padding: 0.75rem 1rem;
		gap: 4px;
	}

	.currency {
		font-weight: 600;
		color: #1c1c1e;
	}

	.amount-input-wrapper input {
		border: none;
		background: transparent;
		font-size: 1rem;
		font-weight: 600;
		color: #007aff;
		width: 100%;
		outline: none;
		font-family: inherit;
	}

	.amount-input-wrapper input::placeholder {
		color: #c7c7cc;
	}

	/* input[type='number'] removed */

	.hint {
		font-size: 0.8rem;
		color: #8e8e93;
		margin: 0;
	}
</style>
