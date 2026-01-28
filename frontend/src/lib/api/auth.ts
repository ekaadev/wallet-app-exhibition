// API functions untuk endpoint users (login, register, profile)

import { apiRequest, apiAuthRequest, type ApiResponse } from './api';

// Response types berdasarkan API spec
export interface UserResponse {
    id: number;
    username: string;
    token: string;
}

export interface UserProfileWalletInfo {
    id: number;
    balance: string;
}

export interface UserProfileResponse {
    id: number;
    username: string;
    role: 'super_admin' | 'admin' | 'user';
    wallet: UserProfileWalletInfo;
}

// Request types
export interface LoginRequest {
    username: string;
    password: string;
}

export interface RegisterRequest {
    username: string;
    password: string;
}

// Login user
export async function login(credentials: LoginRequest): Promise<ApiResponse<UserResponse>> {
    return apiRequest<UserResponse>('/users/login', {
        method: 'POST',
        body: JSON.stringify(credentials),
    });
}

// Register user baru
export async function register(credentials: RegisterRequest): Promise<ApiResponse<UserResponse>> {
    return apiRequest<UserResponse>('/users/register', {
        method: 'POST',
        body: JSON.stringify(credentials),
    });
}

// Get profile user yang sedang login
export async function getProfile(): Promise<ApiResponse<UserProfileResponse>> {
    return apiRequest<UserProfileResponse>('/users/me', {
        method: 'GET',
    });
}

// Response type untuk wallet berdasarkan API spec
export interface WalletResponse {
    id: number;
    user_id: number;
    balance: string;
}

// Get wallet user yang sedang login
export async function getWallet(): Promise<ApiResponse<WalletResponse>> {
    return apiRequest<WalletResponse>('/wallets/me', {
        method: 'GET',
    });
}
