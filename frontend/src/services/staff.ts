import api from './api';

export async function getPendingApprovals() {
  const { data } = await api.get('/approvals/pending');
  return data;
}

export async function approveBookingRequest(requestId: string) {
  const { data } = await api.post(`/approvals/${requestId}/approve`);
  return data;
}

export async function denyBookingRequest(requestId: string) {
  const { data } = await api.post(`/approvals/${requestId}/deny`);
  return data;
}

