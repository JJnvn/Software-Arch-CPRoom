import api from './api';

export async function searchRooms(params: { start: string; end: string; capacity?: number; features?: string[] }) {
  const { data } = await api.get('/rooms/search', { params });
  const rooms = Array.isArray(data?.rooms)
    ? data.rooms.map((r: any) => ({
        id: r.room_id ?? r.id,
        name: r.name,
        capacity: r.capacity,
        features: r.features,
      }))
    : [];
  return { rooms };
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

export async function createRoom(payload: { name: string; capacity?: number; features?: string[] }) {
  const { data } = await api.post('/rooms', payload);
  return data;
}

export async function listRooms() {
  const { data } = await api.get('/rooms');
  return data as Array<{ id: string; name: string; capacity?: number; features?: string[] }>;
}

