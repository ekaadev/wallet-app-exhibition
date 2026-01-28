// API functions untuk endpoint wallets dan wallet mutations

import { apiRequest, type ApiResponse } from './api';

// Response types berdasarkan API spec
export interface WalletMutationResponse {
    id: number;
    wallet_id: number;
    transaction_id: number;
    type: 'credit' | 'debit';
    amount: string;
    balance_before: string;
    balance_after: string;
    created_at: string;
}

export interface WalletMutationListResponse {
    mutations: WalletMutationResponse[];
    total: number;
    page: number;
    limit: number;
}

// Get daftar mutasi wallet
// Get daftar mutasi wallet
export async function getWalletMutations(
    page: number = 1,
    limit: number = 20
): Promise<ApiResponse<WalletMutationListResponse>> {
    return apiRequest<WalletMutationListResponse>(
        `/wallet-mutations?page=${page}&limit=${limit}`,
        { method: 'GET' }
    );
}
