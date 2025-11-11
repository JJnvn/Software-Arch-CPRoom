import api from './api';

// Get notification history for authenticated user (JWT-based)
export async function getNotificationHistory(page: number = 1, pageSize: number = 20) {
  const { data } = await api.get('/notifications/history', {
    params: { page, page_size: pageSize }
  });
  return data;
}

// Mark notification as read (if needed in the future)
export async function markNotificationRead(notificationId: string) {
  const { data } = await api.put(`/notifications/${notificationId}/read`);
  return data;
}

