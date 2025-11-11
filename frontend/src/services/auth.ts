import api from './api';

export async function login(payload: { email: string; password: string }) {
  const { data } = await api.post('/auth/login', payload);
  if (data.token) {
    localStorage.setItem('AUTH_TOKEN', data.token);
  }
  return data;
}

export async function register(payload: { name: string; email: string; password: string }) {
  const { data } = await api.post('/auth/register', payload);
  if (data?.token) {
    localStorage.setItem('AUTH_TOKEN', data.token);
  }
  return data;
}

export async function logout() {
  const { data } = await api.get('/auth/logout');
  return data;
}

export async function getProfile() {
  const { data } = await api.get('/auth/my-profile'); 
  return data;
}

export async function updateProfile(
  id: string,
  payload: Partial<{ name: string; email: string }>
) {
  const { data } = await api.put(`/auth/users/${id}`, payload);
  return data;
}


export async function getBookingHistory() {
  const { data } = await api.get('/bookings/mine');
  return data;
}

export async function updatePreferences(payload: { notificationType?: string; language?: string }) {
  // See notification-service: use /preferences/:userId instead
  return { success: false } as any;
}

