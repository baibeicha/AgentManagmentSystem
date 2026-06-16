import { create } from 'zustand';
import Cookies from 'js-cookie';
import { v4 as uuidv4 } from 'uuid';
import { api } from '@/lib/api';

interface UserProfile {
  user_id: string;
  login: string;
  tenant_id: string;
  role: string;
}

interface AuthState {
  accessToken: string | null;
  profile: UserProfile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  deviceId: string;

  // Actions
  setTokens: (accessToken: string, refreshToken: string) => void;
  setProfile: (profile: UserProfile) => void;
  initialize: () => Promise<void>;
  logout: () => Promise<void>;
  getDeviceId: () => string;
}

const getOrCreateDeviceId = (): string => {
  if (typeof window === 'undefined') return '';
  let deviceId = localStorage.getItem('aegis_web_device_id');
  if (!deviceId) {
    deviceId = uuidv4();
    localStorage.setItem('aegis_web_device_id', deviceId);
  }
  return deviceId;
};

const useAuthStore = create<AuthState>((set, get) => ({
  accessToken: null,
  profile: null,
  isAuthenticated: false,
  isLoading: true,
  deviceId: '',

  setTokens: (accessToken: string, refreshToken: string) => {
    Cookies.set('refresh_token', refreshToken, { secure: true, sameSite: 'strict' });
    set({ accessToken, isAuthenticated: true });
  },

  setProfile: (profile: UserProfile) => {
    set({ profile });
  },

  getDeviceId: () => {
    const id = getOrCreateDeviceId();
    if (get().deviceId !== id) {
      set({ deviceId: id });
    }
    return id;
  },

  initialize: async () => {
    set({ isLoading: true });
    get().getDeviceId();

    // We don't have access token in local storage, so it will be empty initially.
    // However, if we have a refresh token, we can attempt to fetch the profile.
    // The axios interceptor will handle the 401 on /me by using the refresh token.

    const refreshToken = Cookies.get('refresh_token');
    if (refreshToken) {
      try {
        const { data } = await api.get('/api/v1/auth/me');
        set({ profile: data, isAuthenticated: true, isLoading: false });
      } catch (error) {
        // If it fails, logout will handle cleanup
        set({ accessToken: null, profile: null, isAuthenticated: false, isLoading: false });
      }
    } else {
      set({ isLoading: false });
    }
  },

  logout: async () => {
    const refreshToken = Cookies.get('refresh_token');
    if (refreshToken) {
      try {
        // We do a fire-and-forget to the backend
        await api.post('/api/v1/auth/logout', { refresh_token: refreshToken });
      } catch (err) {
        console.error('Logout API call failed:', err);
      }
    }

    Cookies.remove('refresh_token');
    set({ accessToken: null, profile: null, isAuthenticated: false });

    if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
       window.location.href = '/login';
    }
  },
}));

export default useAuthStore;
