import axios from 'axios';
import useAuthStore from '@/hooks/useAuth';
import Cookies from 'js-cookie';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Mutex state for refresh token handling
let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

const subscribeTokenRefresh = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb);
};

const onRefreshed = (token: string) => {
  refreshSubscribers.map((cb) => cb(token));
};

api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().accessToken;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const { config, response } = error;
    const originalRequest = config;

    // We only retry on 401 Unauthorized
    if (response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      if (!isRefreshing) {
        isRefreshing = true;

        try {
          const refreshToken = Cookies.get('refresh_token');
          if (!refreshToken) {
            useAuthStore.getState().logout();
            return Promise.reject(error);
          }

          const { data } = await axios.post(`${API_URL}/api/v1/auth/refresh`, {
            refresh_token: refreshToken,
          });

          const newAccessToken = data.access_token;
          const newRefreshToken = data.refresh_token;

          useAuthStore.getState().setTokens(newAccessToken, newRefreshToken);

          isRefreshing = false;
          onRefreshed(newAccessToken);
          refreshSubscribers = [];

          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return api(originalRequest);
        } catch (refreshError) {
          isRefreshing = false;
          refreshSubscribers = [];
          useAuthStore.getState().logout();
          return Promise.reject(refreshError);
        }
      } else {
        // If already refreshing, wait for the new token
        return new Promise((resolve) => {
          subscribeTokenRefresh((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(api(originalRequest));
          });
        });
      }
    }

    return Promise.reject(error);
  }
);
