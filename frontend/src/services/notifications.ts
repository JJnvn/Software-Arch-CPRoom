import api from './api';

export async function getNotificationHistory(userId: string) {
  const { data } = await api.get(`/notifications/history/${encodeURIComponent(userId)}`);
  return data;
}

export async function updateNotificationPreferences(payload: { notificationType?: string; language?: string }) {
  // See notification-service: should call /preferences/:userId
  return { success: false } as any;
}

