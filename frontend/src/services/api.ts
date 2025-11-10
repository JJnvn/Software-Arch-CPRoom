import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  withCredentials: false, // disable cookies for header-based auth
});

// attach Authorization header automatically
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('AUTH_TOKEN'); // key must match login
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default api;
