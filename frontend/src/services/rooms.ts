import api from './api';

export async function searchRooms(params: { date?: string; time?: string; capacity?: number; features?: string[] }) {
  const { data } = await api.get('/rooms/search', { params });
  return data;
}

export async function getRoomSchedule(roomId: string, params: { date?: string }) {
  const { data } = await api.get(`/rooms/${roomId}/schedule`, { params });
  return data;
}

export async function createBooking(payload: any) {
  const { data } = await api.post('/bookings', payload);
  return data;
}

export async function cancelBooking(bookingId: string) {
  const { data } = await api.post(`/bookings/${bookingId}/cancel`);
  return data;
}

export async function rescheduleBooking(bookingId: string, payload: any) {
  const { data } = await api.post(`/bookings/${bookingId}/reschedule`, payload);
  return data;
}

export async function transferBookingOwnership(bookingId: string, payload: { newOwnerEmail: string }) {
  const { data } = await api.post(`/bookings/${bookingId}/transfer`, payload);
  return data;
}

export async function getAdminRoomBookings(roomId: string) {
  const { data } = await api.get(`/admin/rooms/${roomId}/bookings`);
  return data;
}

