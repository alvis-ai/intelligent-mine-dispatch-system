import { create } from 'zustand';
import apiClient from '../api/client';

interface User {
  id: number;
  username: string;
  realName: string;
  role: number;
  mineId: number;
}

interface AuthState {
  token: string | null;
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  validate: () => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: null,
  loading: false,

  login: async (username, password) => {
    set({ loading: true });
    try {
      const res = await apiClient.post('/api/v1/auth/login', { username, password });
      if (res.data.code === 0) {
        localStorage.setItem('token', res.data.data.token);
        set({ token: res.data.data.token, user: res.data.data.user, loading: false });
      } else {
        throw new Error(res.data.message);
      }
    } catch (err) {
      set({ loading: false });
      throw err;
    }
  },

  validate: async () => {
    const token = localStorage.getItem('token');
    if (!token) return;
    try {
      const res = await apiClient.post('/api/v1/auth/validate', { token });
      if (res.data.code === 0) {
        set({ user: res.data.data });
      }
    } catch {
      localStorage.removeItem('token');
      set({ token: null, user: null });
    }
  },

  logout: () => {
    localStorage.removeItem('token');
    set({ token: null, user: null });
  },
}));
