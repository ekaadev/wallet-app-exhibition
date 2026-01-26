<script>
    import { ArrowLeft, Send } from 'lucide-svelte';

    // State untuk input form (Svelte 5)
    let recipient = $state('');
    let amount = $state('');
    let note = $state('');
    let isLoading = $state(false);

    // Fungsi format rupiah sederhana untuk display (opsional)
    // @ts-ignore
    function formatRupiah(value) {
        if (!value) return '';
        return new Intl.NumberFormat('id-ID').format(value);
    }

    function handleSend() {
        if (!recipient || !amount) {
            alert("Mohon isi tujuan dan jumlah transfer");
            return;
        }

        isLoading = true;
        
        // Simulasi delay kirim
        setTimeout(() => {
            alert(`Berhasil kirim Rp${formatRupiah(amount)} ke ${recipient}`);
            isLoading = false;
            // Reset form atau redirect
        }, 1500);
    }
</script>

<div class="page-container">
    
    <header class="header">
        <a href="/wallet" class="back-btn" aria-label="Kembali">
            <ArrowLeft size={24} color="#007aff" />
        </a>
        <h1>Transfer</h1>
        <div class="spacer"></div> </header>

    <div class="ios-list-group">
        
        <div class="list-row">
            <label for="recipient" class="row-label">Ke</label>
            <input 
                id="recipient" 
                type="text" 
                class="row-input" 
                placeholder="Username / Wallet ID"
                bind:value={recipient}
            />
        </div>

        <div class="list-row">
            <label for="amount" class="row-label">Jumlah</label>
            <div class="amount-wrapper">
                <span class="currency">Rp</span>
                <input 
                    id="amount" 
                    type="number" 
                    class="row-input amount-input" 
                    placeholder="0"
                    bind:value={amount}
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
            />
        </div>
    </div>

    <div class="fixed-footer">
        <button 
            class="send-btn {isLoading ? 'loading' : ''}" 
            onclick={handleSend}
            disabled={isLoading}
        >
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
    /* Reset & Base */
    :global(body) {
        background-color: #ffffff;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    }

    .page-container {
        padding: 1rem 1.5rem;
        max-width: 480px;
        margin: 0 auto;
        min-height: 100vh;
        display: flex;
        flex-direction: column;
    }

    /* --- HEADER --- */
    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 2rem;
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

    .spacer { width: 24px; } /* Penyeimbang layout header */


    /* --- IOS GROUP STYLE (Sama seperti Profile) --- */
    .ios-list-group {
        background-color: #f2f4f6; /* Abu-abu muda */
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
        width: 80px; /* Lebar tetap untuk label agar input sejajar */
        flex-shrink: 0;
    }

    /* --- INPUT STYLES --- */
    .row-input {
        border: none;
        background: transparent;
        font-size: 0.95rem;
        color: #000; /* Input teks hitam */
        width: 100%;
        text-align: right; /* Gaya iOS: Value di kanan */
        font-family: inherit;
        outline: none;
    }

    .row-input::placeholder {
        color: #c7c7cc; /* Placeholder abu-abu */
    }

    /* Khusus Input Angka */
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
        color: #007aff; /* Warna biru untuk jumlah uang */
    }

    /* Hapus spinner di input number */
    input[type=number]::-webkit-inner-spin-button, 
    input[type=number]::-webkit-outer-spin-button { 
        -webkit-appearance: none; 
        margin: 0; 
    }


    .fixed-footer {
        position: fixed;
        bottom: 0;
        /* Trik agar footer tetap di tengah meski di layar desktop lebar */
        left: 50%;
        transform: translateX(-50%);
        width: 100%;
        max-width: 480px; /* Samakan dengan max-width container */
        
        background: rgba(255, 255, 255, 0.9); /* Latar belakang agak transparan */
        backdrop-filter: blur(10px); /* Efek blur glass */
        border-top: 1px solid rgba(0,0,0,0.05); /* Garis batas tipis */
        
        padding: 1rem 1.5rem 2rem 1.5rem; /* Padding bawah lebih besar untuk area HP swipe */
        z-index: 100;
    }

    .send-btn {
        background-color: #1c1c1e;
        color: white;
        width: 100%;
        padding: 1rem;
        border: none;
        border-radius: 50px; /* Tombol bulat penuh lebih modern untuk floating */
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        box-shadow: 0 4px 15px rgba(0,0,0,0.15); /* Shadow agar terlihat melayang */
    }

    .send-btn:active {
        transform: scale(0.98);
    }

    .send-btn:disabled {
        background-color: #8e8e93;
        cursor: not-allowed;
    }
</style>