import api from './api';

// Booking-related functions
export async function cancelBooking(bookingId: string) {
  const { data } = await api.post(`/bookings/${bookingId}/cancel`);
  return data;
}

export async function rescheduleBooking(bookingId: string, payload: { start_time: string; end_time: string }) {
  const { data } = await api.put(`/bookings/${bookingId}`, payload);
  return data;
}

export async function transferBooking(bookingId: string, payload: { new_user_email: string }) {
  const { data } = await api.post(`/bookings/${bookingId}/transfer`, payload);
  return data;
}

export async function getRoomSchedule(roomId: string, date: string) {
  const { data } = await api.get(`/rooms/${roomId}/schedule`, { params: { date } });
  return data;
}

export async function getAdminRoomBookings(roomId: string) {
  const { data } = await api.get(`/admin/rooms/${roomId}/bookings`);
  return data;
}

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

export async function updateProfile(payload: Partial<{ name: string; email: string; password: string }>) {
  const { data } = await api.put('/auth/profile', payload);
  return data;
}


export async function getBookingHistory() {
  const { data } = await api.get('/bookings/mine');
  return data;
}

export async function updatePreferences(payload: { notification_type?: string; language?: string; enabled_channels?: string[] }) {
  const { data } = await api.put('/preferences', payload);
  return data;
}

export async function getPreferences() {
  const { data } = await api.get('/preferences');
  return data;
}

