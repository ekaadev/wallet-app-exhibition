// API functions untuk endpoint transactions (transfer, topup, list)

import { apiRequest, type ApiResponse } from './api';

// Response type untuk transaction
export interface TransactionResponse {
    id: number;
    type: 'top_up' | 'transfer';
    amount: string;
    from_wallet_id: number | null;
    to_wallet_id: number;
    performed_by_user_id: number;
    status: 'pending' | 'completed' | 'failed';
    description: string | null;
    created_at: string;
}

// Request type untuk transfer
export interface TransferRequest {
    to_user_id: number;
    amount: number;
    description?: string;
}

// Request type untuk topup (super admin only)
export interface TopUpRequest {
    to_user_id: number;
    amount: number;
    description?: string;
}

// Response type untuk list transactions
export interface TransactionListResponse {
    transactions: TransactionResponse[];
    total: number;
    page: number;
    limit: number;
}

// Transfer saldo ke user lain
export async function transfer(
    request: TransferRequest
): Promise<ApiResponse<TransactionResponse>> {
    return apiRequest<TransactionResponse>('/transactions/transfer', {
        method: 'POST',
        body: JSON.stringify(request),
    });
}

// Top-up saldo ke user (super admin only)
export async function topUp(
    request: TopUpRequest
): Promise<ApiResponse<TransactionResponse>> {
    return apiRequest<TransactionResponse>('/transactions/topup', {
        method: 'POST',
        body: JSON.stringify(request),
    });
}

// Get list transactions dengan pagination
export async function getTransactions(
    page: number = 1, 
    limit: number = 10
): Promise<ApiResponse<TransactionListResponse>> {
    return apiRequest<TransactionListResponse>(
        `/transactions?page=${page}&limit=${limit}`, 
        { method: 'GET' }
    );
}
