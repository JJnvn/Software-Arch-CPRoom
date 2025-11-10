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
  return data;
}

export async function logout() {
  const { data } = await api.get('/auth/logout');
  return data;
}

export async function getProfile() {
  const { data } = await api.get('/users/me'); 
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
  const { data } = await api.put('/users/me/preferences', payload);
  return data;
}

