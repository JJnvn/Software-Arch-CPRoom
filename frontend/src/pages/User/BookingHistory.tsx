import { useEffect, useState } from 'react';
import BookingCard from '@/components/BookingCard';
import * as auth from '@/services/auth';
import { useNavigate } from 'react-router-dom';

export default function BookingHistory() {
  const [bookings, setBookings] = useState<any[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        const data = await auth.getBookingHistory();
        setBookings(data ?? []);
      } catch {
        setBookings([]);
      }
    })();
  }, []);

  const formatDateTime = (isoString?: string) => {
    if (!isoString) return { date: '—', time: '—' };
    const d = new Date(isoString);
    const date = d.toISOString().slice(0, 10);
    const time = d.toISOString().slice(11, 16);
    return { date, time };
  };

  return (
    <div className="page">
      <h1 className="page-title">Booking History</h1>
      <div className="space-y-3">
        {bookings.length === 0 && <div className="text-gray-600">No bookings found.</div>}
        {bookings.map((b) => {
          const { date, time } = formatDateTime(b.start_time);
          return (
            <BookingCard
              key={b.booking_id}
              booking={{ 
                id: b.booking_id, 
                room: b.room_name || b.room_id || '—', 
                date, 
                time, 
                status: b.status || '—' 
              }}
              onCancel={(id) => alert(`Cancel booking ${id}`)}
              onEdit={(id) => navigate(`/bookings/${id}/edit`)}
              onTransfer={(id) => navigate(`/bookings/${id}/transfer`)}
            />
          );
        })}
      </div>
    </div>
  );
}

