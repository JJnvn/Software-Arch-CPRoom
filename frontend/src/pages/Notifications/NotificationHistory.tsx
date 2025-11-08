import { useEffect, useState } from 'react';
import * as notifications from '@/services/notifications';

export default function NotificationHistory() {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    (async () => {
      const data = await notifications.getNotificationHistory();
      setItems(data.notifications ?? []);
    })();
  }, []);

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

