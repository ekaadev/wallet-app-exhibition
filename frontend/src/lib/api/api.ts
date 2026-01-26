// API base configuration dan helper functions untuk HTTP requests

const API_BASE_URL = 'http://localhost:3000';

// Tipe response standar dari API
export interface ApiResponse<T> {
    data?: T;
    message?: string;
}

// Tipe error response
export interface ApiError {
    message: string;
}

// Helper untuk membuat request tanpa auth header
export async function apiRequest<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<ApiResponse<T>> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const defaultHeaders: HeadersInit = {
        'Content-Type': 'application/json',
    };

    const response = await fetch(url, {
        ...options,
        headers: {
            ...defaultHeaders,
            ...options.headers,
        },
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.message || 'Terjadi kesalahan');
    }

    return data;
}

// Helper untuk membuat request dengan auth header
export async function apiAuthRequest<T>(
    endpoint: string,
    token: string,
    options: RequestInit = {}
): Promise<ApiResponse<T>> {
    return apiRequest<T>(endpoint, {
        ...options,
        headers: {
            ...options.headers,
            'Authorization': `Bearer ${token}`,
        },
    });
}
