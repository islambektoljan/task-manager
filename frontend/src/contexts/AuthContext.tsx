import React, { createContext, useContext, useReducer, useEffect } from 'react';
import type {AuthState, LoginRequest, RegisterRequest} from '../types/auth';
import type {User} from '../types/common';
import { authService } from '../services/authService';

type AuthAction =
    | { type: 'AUTH_START' }
    | { type: 'AUTH_SUCCESS'; payload: { user: User; token: string } }
    | { type: 'AUTH_FAILURE'; payload: string }
    | { type: 'LOGOUT' }
    | { type: 'CLEAR_ERROR' };

const initialState: AuthState = {
    user: null,
    token: localStorage.getItem('token'),
    isAuthenticated: !!localStorage.getItem('token'),
    loading: false,
    error: null,
};

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
    switch (action.type) {
        case 'AUTH_START':
            return { ...state, loading: true, error: null };
        case 'AUTH_SUCCESS':
            return {
                ...state,
                loading: false,
                isAuthenticated: true,
                user: action.payload.user,
                token: action.payload.token,
                error: null,
            };
        case 'AUTH_FAILURE':
            return { ...state, loading: false, error: action.payload, isAuthenticated: false };
        case 'LOGOUT':
            return { ...initialState, token: null, isAuthenticated: false };
        case 'CLEAR_ERROR':
            return { ...state, error: null };
        default:
            return state;
    }
};

interface AuthContextType extends AuthState {
    login: (credentials: LoginRequest) => Promise<void>;
    register: (userData: RegisterRequest) => Promise<void>;
    logout: () => Promise<void>;
    clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [state, dispatch] = useReducer(authReducer, initialState);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const user = localStorage.getItem('user');

        if (token && user) {
            try {
                const userData = JSON.parse(user);
                dispatch({
                    type: 'AUTH_SUCCESS',
                    payload: { user: userData, token }
                });
            } catch (error) {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
            }
        }
    }, []);

    const login = async (credentials: LoginRequest) => {
        dispatch({ type: 'AUTH_START' });
        try {
            const response = await authService.login(credentials);
            if (response.success && response.data) {
                const user = {
                    id: response.data.user_id,
                    email: response.data.email,
                    role: response.data.role,
                };

                localStorage.setItem('token', response.data.token);
                localStorage.setItem('user', JSON.stringify(user));

                dispatch({
                    type: 'AUTH_SUCCESS',
                    payload: { user, token: response.data.token },
                });
            } else {
                dispatch({ type: 'AUTH_FAILURE', payload: response.error || 'Login failed' });
            }
        } catch (error: any) {
            dispatch({
                type: 'AUTH_FAILURE',
                payload: error.response?.data?.error || 'Login failed',
            });
        }
    };

    const register = async (userData: RegisterRequest) => {
        dispatch({ type: 'AUTH_START' });
        try {
            const response = await authService.register(userData);
            if (response.success && response.data) {
                const user = {
                    id: response.data.user_id,
                    email: response.data.email,
                    role: response.data.role,
                };

                localStorage.setItem('token', response.data.token);
                localStorage.setItem('user', JSON.stringify(user));

                dispatch({
                    type: 'AUTH_SUCCESS',
                    payload: { user, token: response.data.token },
                });
            } else {
                dispatch({ type: 'AUTH_FAILURE', payload: response.error || 'Registration failed' });
            }
        } catch (error: any) {
            dispatch({
                type: 'AUTH_FAILURE',
                payload: error.response?.data?.error || 'Registration failed',
            });
        }
    };

    const logout = async () => {
        try {
            await authService.logout();
        } catch (error) {
            console.error('Logout error:', error);
        } finally {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            dispatch({ type: 'LOGOUT' });
        }
    };

    const clearError = () => {
        dispatch({ type: 'CLEAR_ERROR' });
    };

    return (
        <AuthContext.Provider
            value={{
                ...state,
                login,
                register,
                logout,
                clearError,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};