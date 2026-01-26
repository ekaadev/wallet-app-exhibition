// Auth store untuk state management autentikasi

import { writable, derived, get } from 'svelte/store';
import { login as apiLogin, register as apiRegister, getProfile, type UserResponse, type UserProfileResponse } from '$lib/api';
import { connectWebSocket, disconnectWebSocket } from '$lib/api/websocket';

// Key untuk localStorage
const TOKEN_KEY = 'wallet_auth_token';
const USER_KEY = 'wallet_auth_user';

// Interface untuk user data yang disimpan di store
export interface AuthUser {
    id: number;
    username: string;
    role?: 'super_admin' | 'admin' | 'user';
    wallet?: {
        id: number;
        balance: string;
    };
}

// Fungsi helper untuk membaca dari localStorage (browser only)
function getStoredToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
}

function getStoredUser(): AuthUser | null {
    if (typeof window === 'undefined') return null;
    const stored = localStorage.getItem(USER_KEY);
    return stored ? JSON.parse(stored) : null;
}

// Stores
export const token = writable<string | null>(getStoredToken());
export const user = writable<AuthUser | null>(getStoredUser());
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Derived store untuk mengecek apakah user terautentikasi
export const isAuthenticated = derived(token, ($token) => !!$token);

// Subscribe ke token untuk menyimpan ke localStorage
token.subscribe((value) => {
    if (typeof window === 'undefined') return;
    if (value) {
        localStorage.setItem(TOKEN_KEY, value);
    } else {
        localStorage.removeItem(TOKEN_KEY);
    }
});

// Subscribe ke user untuk menyimpan ke localStorage
user.subscribe((value) => {
    if (typeof window === 'undefined') return;
    if (value) {
        localStorage.setItem(USER_KEY, JSON.stringify(value));
    } else {
        localStorage.removeItem(USER_KEY);
    }
});

// Action: Login
export async function login(username: string, password: string): Promise<boolean> {
    loading.set(true);
    error.set(null);

    try {
        const response = await apiLogin({ username, password });
        
        if (response.data) {
            const userData = response.data;
            token.set(userData.token);
            user.set({
                id: userData.id,
                username: userData.username,
            });
            
            // Connect WebSocket setalah login berhasil
            connectWebSocket(userData.token);
            
            return true;
        }
        
        error.set('Login gagal');
        return false;
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'Login gagal');
        return false;
    } finally {
        loading.set(false);
    }
}

// Action: Register
export async function register(username: string, password: string): Promise<boolean> {
    loading.set(true);
    error.set(null);

    try {
        const response = await apiRegister({ username, password });
        
        if (response.data) {
            const userData = response.data;
            token.set(userData.token);
            user.set({
                id: userData.id,
                username: userData.username,
            });
            
            // Connect WebSocket setelah register berhasil
            connectWebSocket(userData.token);

            return true;
        }
        
        error.set('Registrasi gagal');
        return false;
    } catch (err) {
        error.set(err instanceof Error ? err.message : 'Registrasi gagal');
        return false;
    } finally {
        loading.set(false);
    }
}

// Action: Load profile dari API
export async function loadProfile(): Promise<boolean> {
    const currentToken = get(token);
    if (!currentToken) {
        error.set('Tidak terautentikasi');
        return false;
    }

    loading.set(true);
    error.set(null);

    try {
        const response = await getProfile(currentToken);
        
        if (response.data) {
            const profile = response.data;
            user.set({
                id: profile.id,
                username: profile.username,
                role: profile.role,
                wallet: profile.wallet,
            });

            // Connect WebSocket jika session exist
            connectWebSocket(currentToken);

            return true;
        }
        
        error.set('Gagal memuat profil');
        return false;
    } catch (err) {
        // Jika gagal, kemungkinan token expired
        error.set(err instanceof Error ? err.message : 'Gagal memuat profil');
        return false;
    } finally {
        loading.set(false);
    }
}

// Action: Logout
export function logout(): void {
    disconnectWebSocket(); // Disconnect WS
    token.set(null);
    user.set(null);
    error.set(null);
}

// Action: Clear error
export function clearError(): void {
    error.set(null);
}
