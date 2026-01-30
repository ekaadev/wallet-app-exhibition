<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { X, Camera, AlertCircle } from 'lucide-svelte';
	import { Html5Qrcode } from 'html5-qrcode';

	// Props
	let { show = $bindable(false), onClose = () => {} }: { show: boolean; onClose?: () => void } =
		$props();

	// State
	let scannerContainer: HTMLDivElement | null = $state(null);
	let scanner: Html5Qrcode | null = $state(null);
	let scanError = $state('');
	let isScanning = $state(false);
	let permissionDenied = $state(false);

	// Watch show state untuk start/stop scanner
	$effect(() => {
		if (show && scannerContainer) {
			startScanner();
		} else if (!show && scanner) {
			stopScanner();
		}
	});

	// Inisialisasi dan start scanner
	async function startScanner() {
		if (!scannerContainer || isScanning) return;

		scanError = '';
		permissionDenied = false;

		try {
			scanner = new Html5Qrcode('qr-reader');
			isScanning = true;

			await scanner.start(
				{ facingMode: 'environment' },
				{
					fps: 10,
					qrbox: { width: 250, height: 250 },
					aspectRatio: 1
				},
				onScanSuccess,
				onScanError
			);
		} catch (err) {
			isScanning = false;
			console.error('Scanner error:', err);

			// Handle permission denied
			if (err instanceof Error && err.message.includes('Permission')) {
				permissionDenied = true;
				scanError = 'Akses kamera ditolak. Izinkan akses kamera untuk scan QR.';
			} else {
				scanError = 'Gagal membuka kamera. Pastikan perangkat memiliki kamera.';
			}
		}
	}

	// Stop scanner
	async function stopScanner() {
		if (scanner && isScanning) {
			try {
				await scanner.stop();
			} catch (err) {
				console.error('Stop scanner error:', err);
			}
		}
		isScanning = false;
		scanner = null;
	}

	// Callback saat QR berhasil di-scan
	function onScanSuccess(decodedText: string) {
		try {
			// Parse QR data
			const data = JSON.parse(decodedText);

			if (data.userId) {
				// Stop scanner dan close modal
				stopScanner();
				handleClose();

				// Build URL params
				const params = new URLSearchParams();
				params.set('userId', String(data.userId));

				if (data.amount && data.amount > 0) {
					params.set('amount', String(data.amount));
				}

				// Redirect ke transaction page dengan params
				goto(`/transaction?${params.toString()}`);
			} else {
				scanError = 'QR Code tidak valid. Pastikan scan QR dari pengguna lain.';
			}
		} catch {
			scanError = 'Format QR Code tidak dikenali.';
		}
	}

	// Callback error scanning (dipanggil setiap frame yang tidak berhasil)
	function onScanError(_errorMessage: string) {
		// Abaikan error ini, karena dipanggil setiap frame tanpa QR
	}

	// Handle close modal
	function handleClose() {
		stopScanner();
		show = false;
		scanError = '';
		permissionDenied = false;
		onClose();
	}

	// Handle klik di luar modal
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
		stopScanner();
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
				<h2>Scan QR Code</h2>
				<button class="close-btn" onclick={handleClose} aria-label="Tutup">
					<X size={24} />
				</button>
			</div>

			{#if scanError}
				<div class="error-box">
					<AlertCircle size={20} />
					<span>{scanError}</span>
				</div>
			{/if}

			<div class="scanner-wrapper">
				{#if permissionDenied}
					<div class="permission-denied">
						<Camera size={48} />
						<p>Kamera tidak tersedia</p>
						<button class="retry-btn" onclick={startScanner}>Coba Lagi</button>
					</div>
				{:else}
					<div id="qr-reader" bind:this={scannerContainer}></div>
				{/if}
			</div>

			<p class="hint">Arahkan kamera ke QR Code pengguna lain untuk melakukan transfer</p>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		background-color: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal-content {
		background-color: #1c1c1e;
		border-radius: 20px;
		padding: 1.5rem;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
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
		margin-bottom: 1rem;
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.2rem;
		font-weight: 700;
		color: #ffffff;
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
		color: #ffffff;
	}

	.error-box {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background-color: rgba(220, 38, 38, 0.2);
		color: #fca5a5;
		padding: 0.75rem 1rem;
		border-radius: 10px;
		margin-bottom: 1rem;
		font-size: 0.85rem;
	}

	.scanner-wrapper {
		background-color: #000;
		border-radius: 12px;
		overflow: hidden;
		aspect-ratio: 1;
		position: relative;
	}

	.scanner-wrapper :global(#qr-reader) {
		width: 100%;
		height: 100%;
	}

	.scanner-wrapper :global(#qr-reader video) {
		width: 100% !important;
		height: 100% !important;
		object-fit: cover;
	}

	.scanner-wrapper :global(#qr-reader__scan_region) {
		min-height: unset !important;
	}

	.scanner-wrapper :global(#qr-reader__dashboard) {
		display: none !important;
	}

	.permission-denied {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		min-height: 300px;
		color: #8e8e93;
		gap: 1rem;
	}

	.permission-denied p {
		margin: 0;
		font-size: 0.95rem;
	}

	.retry-btn {
		background-color: #007aff;
		color: #fff;
		border: none;
		border-radius: 20px;
		padding: 0.6rem 1.5rem;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
	}

	.retry-btn:hover {
		background-color: #0056b3;
	}

	.hint {
		margin-top: 1rem;
		font-size: 0.85rem;
		color: #8e8e93;
		text-align: center;
	}
</style>
