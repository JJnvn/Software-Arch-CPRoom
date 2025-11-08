import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';

export default function AdminRoomBookings() {
  const { id } = useParams();
  const [bookings, setBookings] = useState<any[]>([]);

  useEffect(() => {
    (async () => {
      if (!id) return;
      const data = await rooms.getAdminRoomBookings(id);
      setBookings(data.bookings ?? []);
    })();
  }, [id]);

  return (
    <div className="page">
      <h1 className="page-title">Admin Room Bookings</h1>
      <div className="card mb-3">Room ID: <span className="font-medium">{id}</span></div>
      <div className="space-y-2">
        {bookings.length === 0 && <div className="text-gray-600">No bookings for this room.</div>}
        {bookings.map((b) => (
          <div key={b.id} className="card">
            <div className="font-medium">{b.date} • {b.time}</div>
            <div className="text-sm text-gray-600">Owner: {b.owner || '—'}</div>
            <div className="text-xs text-gray-500">Status: {b.status || '—'}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

