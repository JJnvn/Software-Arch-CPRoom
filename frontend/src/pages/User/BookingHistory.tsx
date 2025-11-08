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
        setBookings(data.bookings ?? []);
      } catch {
        setBookings([]);
      }
    })();
  }, []);

  return (
    <div className="page">
      <h1 className="page-title">Booking History</h1>
      <div className="space-y-3">
        {bookings.length === 0 && <div className="text-gray-600">No bookings found.</div>}
        {bookings.map((b) => (
          <BookingCard
            key={b.id}
            booking={{ id: b.id, room: b.room || '—', date: b.date || '—', time: b.time || '—', status: b.status || '—' }}
            onCancel={(id) => alert(`Cancel booking ${id}`)}
            onEdit={(id) => navigate(`/bookings/${id}/edit`)}
            onTransfer={(id) => navigate(`/bookings/${id}/transfer`)}
          />
        ))}
      </div>
    </div>
  );
}

