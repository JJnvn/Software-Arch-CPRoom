import { useEffect, useState } from 'react';
import * as notifications from '@/services/notifications';
import { useAuth } from '@/hooks/useAuth';

export default function NotificationHistory() {
  const [items, setItems] = useState<any[]>([]);
  const { user } = useAuth();

  useEffect(() => {
    (async () => {
      if (!user?.id) return;
      const data = await notifications.getNotificationHistory(user.id);
      setItems(data.notifications ?? []);
    })();
  }, [user?.id]);

  return (
    <div className="page">
      <h1 className="page-title">Notification History</h1>
      <div className="space-y-2">
        {items.length === 0 && <div className="text-gray-600">No notifications found.</div>}
        {items.map((n) => (
          <div key={n.id} className="card">
            <div className="font-medium">{n.title || 'Notification'}</div>
            <div className="text-sm text-gray-600">{n.message || 'No details'}</div>
            <div className="text-xs text-gray-500">{n.createdAt || ''}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

