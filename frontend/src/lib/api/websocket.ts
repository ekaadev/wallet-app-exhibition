// WebSocket service untuk komunikasi real-time

import { user, type AuthUser } from '$lib/stores/auth';
import { get } from 'svelte/store';
import { toast } from 'svelte-sonner';

const WS_URL = import.meta.env.PUBLIC_WS_URL;
let ws: WebSocket | null = null;
let reconnectValid = true;
let reconnectInterval = 5000;
let beforeUnloadHandler: ((e: BeforeUnloadEvent) => void) | null = null;

// Tipe payload dari backend
interface WebsocketMessage {
    type: 'transaction' | 'wallet_update';
    payload: any;
}

interface WalletUpdatePayload {
    wallet_id: number;
    new_balance: string;
    mutation_type: 'credit' | 'debit';
    amount: string;
}

interface TransactionPayload {
    transaction_id: number;
    transaction_type: 'top_up' | 'transfer';
    amount: string;
    description: string;
    from_user_id?: number;
    to_user_id: number;
}

export function connectWebSocket() {
    if (ws) {
        ws.onclose = null;
        ws.onerror = null;
        ws.close();
    }

    reconnectValid = true;
    // Browser automatically sends cookies with WebSocket connection
    const url = WS_URL;
    
    try {
        ws = new WebSocket(url);

        // Handle page unload to prevent reconnection loops
        beforeUnloadHandler = () => {
            reconnectValid = false;
            if (ws) {
                ws.onclose = null;
                ws.onerror = null;
                ws.close();
            }
        };
        window.addEventListener('beforeunload', beforeUnloadHandler);

        ws.onopen = () => {
            console.log('WebSocket connected');
        };

        ws.onmessage = (event) => {
            try {
                const message: WebsocketMessage = JSON.parse(event.data);
                handleMessage(message);
            } catch (err) {
                console.error('Error parsing websocket message:', err);
            }
        };

        ws.onclose = (event) => {
            console.log(`WebSocket disconnected. Code: ${event.code}, Reason: ${event.reason || 'No reason'}`);
            if (reconnectValid) {
                console.log('Reconnecting in 5s...');
                setTimeout(() => {
                    if (reconnectValid) connectWebSocket();
                }, reconnectInterval);
            }
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            ws?.close();
        };

    } catch (err) {
        console.error('Failed to connect websocket:', err);
    }
}

export function disconnectWebSocket() {
    reconnectValid = false;
    
    if (beforeUnloadHandler) {
        window.removeEventListener('beforeunload', beforeUnloadHandler);
        beforeUnloadHandler = null;
    }

    if (ws) {
        ws.onclose = null;
        ws.onerror = null;
        ws.close();
        ws = null;
    }
}

function handleMessage(message: WebsocketMessage) {
    console.log('Received WS message:', message);

    switch (message.type) {
        case 'wallet_update':
            handleWalletUpdate(message.payload);
            break;
        case 'transaction':
            handleTransaction(message.payload);
            break;
    }
}

function handleWalletUpdate(payload: WalletUpdatePayload) {
    // Update store saldo user secara real-time
    const currentUser = get(user);
    
    if (currentUser && currentUser.wallet && currentUser.wallet.id === payload.wallet_id) {
        user.update((u) => {
            if (u && u.wallet) {
                return {
                    ...u,
                    wallet: {
                        ...u.wallet,
                        balance: payload.new_balance
                    }
                };
            }
            return u;
        });

        // Tampilkan toast
        const type = payload.mutation_type === 'credit' ? 'masuk' : 'keluar';
        const amount = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(Number(payload.amount));
        
        toast.info(`Saldo terupdate: Uang ${type} ${amount}`, {
            amount: payload.new_balance
        } as any);
        console.log(`Saldo terupdate: ${payload.new_balance}`);
    }
}

function handleTransaction(payload: TransactionPayload) {
    // Bisa digunakan untuk menampilkan notifikasi popup transaksi baru
    const type = payload.transaction_type === 'top_up' ? 'Top Up' : 'Transfer';
    const amount = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(Number(payload.amount));
    
    // Tampilkan notifikasi browser native jika didukung/diizinkan
    if ('Notification' in window && Notification.permission === 'granted') {
        new Notification('Transaksi Baru', {
            body: `${type} berhasil sebesar ${amount}`
        });
    } else {
        console.log(`Notifikasi Transaksi: ${type} ${amount}`);
        // Tampilkan toast juga untuk transaksi
        toast.success(`Transaksi ${type} berhasil: ${amount}`);
    }
}
