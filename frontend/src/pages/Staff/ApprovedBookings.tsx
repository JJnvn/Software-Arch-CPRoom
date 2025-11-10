import { useEffect, useState } from 'react';
import * as staff from '@/services/staff';

export default function ApprovedBookings() {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    (async () => {
      const data = await staff.getApprovedApprovals();
      setItems(data.requests ?? []);
    })();
  }, []);

  return (
    <div className="page">
      <h1 className="page-title">Approved Bookings</h1>
      <div className="space-y-2">
        {items.length === 0 && <div className="text-gray-600">No approved bookings.</div>}
        {items.map((r) => (
          <div key={r.id} className="card flex items-center justify-between">
            <div>
              <div className="font-medium">Room {r.room || '—'} • {r.date || '—'} {r.time || ''}</div>
              <div className="text-sm text-gray-600">User: {r.requester || '—'}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

