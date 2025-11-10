import api from './api';

function fmtDateTimeFromSeconds(seconds?: number) {
  if (!seconds && seconds !== 0) return { date: '', time: '' };
  const d = new Date(seconds * 1000);
  const date = d.toISOString().slice(0, 10);
  const time = d.toISOString().slice(11, 16);
  return { date, time };
}

export async function getPendingApprovals() {
  const { data } = await api.get('/approvals/pending');
  const items = Array.isArray(data?.pending) ? data.pending : [];
  const requests = items.map((it: any) => {
    const start = fmtDateTimeFromSeconds(it?.start?.seconds);
    const roomLabel = it?.room_name || it?.roomId || it?.room_id;
    const userLabel = it?.user_name || it?.userId || it?.user_id;
    return {
      id: it.booking_id || it.bookingId,
      room: roomLabel,
      requester: userLabel,
      date: start.date,
      time: start.time,
      raw: it,
    };
  });
  return { requests };
}

export async function approveBookingRequest(requestId: string) {
  const { data } = await api.post(`/approvals/${requestId}/approve`);
  return data;
}

export async function denyBookingRequest(requestId: string, reason: string) {
  const { data } = await api.post(`/approvals/${requestId}/deny`, { reason });
  return data;
}

export async function getApprovedApprovals() {
  const { data } = await api.get('/approvals/approved');
  const items = Array.isArray(data?.approved) ? data.approved : [];
  const requests = items.map((it: any) => {
    const start = fmtDateTimeFromSeconds(it?.start?.seconds);
    const roomLabel = it?.room_name || it?.roomId || it?.room_id;
    const userLabel = it?.user_name || it?.userId || it?.user_id;
    return {
      id: it.booking_id || it.bookingId,
      room: roomLabel,
      requester: userLabel,
      date: start.date,
      time: start.time,
      raw: it,
    };
  });
  return { requests };
}

