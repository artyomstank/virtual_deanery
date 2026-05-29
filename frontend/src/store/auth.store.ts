import { create } from 'zustand';
import { User } from '../types/api';
import { apiClient } from '../lib/api';

interface AuthState {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isLoading: false,
  error: null,

  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.post('/auth/login', { email, password });
      const { access_token, user } = response;

      apiClient.setToken(access_token);
      set({
        user,
        token: access_token,
        isLoading: false,
        error: null,
      });
    } catch (error: any) {
      set({
        isLoading: false,
        error: error.message || 'Login failed',
      });
      throw error;
    }
  },

  logout: () => {
    apiClient.clearToken();
    set({
      user: null,
      token: null,
      error: null,
    });
  },

  setUser: (user) => {
    set({ user });
  },

  setToken: (token) => {
    set({ token });
    if (token) {
      apiClient.setToken(token);
    } else {
      apiClient.clearToken();
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
