import api from './api';

export async function getNotificationHistory() {
  const { data } = await api.get('/notifications');
  return data;
}

export async function updateNotificationPreferences(payload: { notificationType?: string; language?: string }) {
  // See notification-service: should call /preferences/:userId
  return { success: false } as any;
}

