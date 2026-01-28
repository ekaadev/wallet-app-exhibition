// Auth store untuk state management autentikasi

import { writable, derived, get } from 'svelte/store';
import { login as apiLogin, register as apiRegister, getProfile, apiRequest, type UserResponse, type UserProfileResponse } from '$lib/api';
import { connectWebSocket, disconnectWebSocket } from '$lib/api/websocket';

// Stores
export const token = writable<string | null>(null);
export const user = writable<AuthUser | null>(null);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Derived store untuk mengecek apakah user terautentikasi
export const isAuthenticated = derived(user, ($user) => !!$user);


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
            connectWebSocket();
            
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
            connectWebSocket();

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
    loading.set(true);
    error.set(null);

    try {
        const response = await getProfile();
        
        if (response.data) {
            const profile = response.data;
            user.set({
                id: profile.id,
                username: profile.username,
                role: profile.role,
                wallet: profile.wallet,
            });

            // Connect WebSocket jika session exist
            // Kita butuh trigger khusus untuk check cookie or just connect
            connectWebSocket();

            return true;
        }
        
        error.set('Gagal memuat profil');
        return false;
    } catch (err) {
        // Jika gagal, kemungkinan token expired
        error.set(err instanceof Error ? err.message : 'Gagal memuat profil');
        user.set(null); // Clear user if failed implies invalid session
        return false;
    } finally {
        loading.set(false);
    }
}

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

// Action: Logout
export async function logout(): Promise<void> {
    try {
        await apiRequest('/users/logout', { method: 'POST' });
    } catch (e) {
        console.error('Logout error:', e);
    }
    disconnectWebSocket(); // Disconnect WS
    token.set(null);
    user.set(null);
    error.set(null);
    // Force reload to clear any memory/state if needed, or just redirect
    if (typeof window !== 'undefined') {
        window.location.href = '/login';
    }
}

// Action: Clear error
export function clearError(): void {
    error.set(null);
}
