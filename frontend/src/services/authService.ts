import { api } from './api';
import type { LoginRequest, RegisterRequest, AuthResponse } from '../types/auth';
import type { ApiResponse } from '../types/common';

export const authService = {
    async login(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>> {
        console.log('🔐 Sending login request:', credentials);
        try {
            const response = await api.post('/login', credentials);
            console.log('✅ Login response:', response.data);
            return response.data;
        } catch (error: any) {
            console.error('❌ Login error:', {
                status: error.response?.status,
                data: error.response?.data,
                message: error.message
            });
            throw error;
        }
    },

    async register(userData: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
        console.log('📝 Sending register request:', userData);
        try {
            const response = await api.post('/register', userData);
            console.log('✅ Register response:', response.data);
            return response.data;
        } catch (error: any) {
            console.error('❌ Register error:', {
                status: error.response?.status,
                data: error.response?.data,
                message: error.message
            });
            throw error;
        }
    },

    async logout(): Promise<void> {
        try {
            const response = await api.post('/logout');
            return response.data;
        } catch (error: any) {
            console.error('❌ Logout error:', error);
            throw error;
        }
    },

    async refreshToken(): Promise<ApiResponse<AuthResponse>> {
        try {
            const response = await api.post('/refresh');
            return response.data;
        } catch (error: any) {
            console.error('❌ Refresh token error:', error);
            throw error;
        }
    },

    async healthCheck(): Promise<ApiResponse<any>> {
        try {
            const response = await api.get('/health');
            return response.data;
        } catch (error: any) {
            console.error('❌ Health check error:', error);
            throw error;
        }
    }
};