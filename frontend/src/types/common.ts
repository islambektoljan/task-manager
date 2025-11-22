export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    error?: string;
    code?: number;
}

export interface User {
    id: string;
    email: string;
    role: string;
}