import api from './api';

export async function getPendingApprovals() {
  const { data } = await api.get('/staff/pending-approvals');
  return data;
}

export async function approveBookingRequest(requestId: string) {
  const { data } = await api.post(`/staff/requests/${requestId}/approve`);
  return data;
}

export async function denyBookingRequest(requestId: string) {
  const { data } = await api.post(`/staff/requests/${requestId}/deny`);
  return data;
}

