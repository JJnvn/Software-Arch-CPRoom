import api from './api';

export async function getNotificationHistory() {
  const { data } = await api.get('/notifications');
  return data;
}

export async function updateNotificationPreferences(payload: { notificationType?: string; language?: string }) {
  const { data } = await api.put('/users/me/notification-preferences', payload);
  return data;
}

