import { useEffect, useState } from 'react';
import BookingCard from '@/components/BookingCard';
import * as auth from '@/services/auth';
import { useNavigate } from 'react-router-dom';

export default function BookingHistory() {
  const [bookings, setBookings] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [alert, setAlert] = useState<{ type: 'success' | 'error'; message: string } | null>(null);
  const [cancellingId, setCancellingId] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadBookings();
  }, []);

  async function loadBookings() {
    setLoading(true);
    try {
      const data = await auth.getBookingHistory();
      setBookings(data ?? []);
    } catch {
      setBookings([]);
      showAlert('error', 'Failed to load booking history');
    } finally {
      setLoading(false);
    }
  }

  function showAlert(type: 'success' | 'error', message: string) {
    setAlert({ type, message });
    setTimeout(() => setAlert(null), 5000);
  }

  const formatDateTime = (isoString?: string) => {
    if (!isoString) return { date: '—', time: '—' };
    const d = new Date(isoString);
    const date = d.toISOString().slice(0, 10);
    const time = d.toISOString().slice(11, 16);
    return { date, time };
  };

  const handleCancel = async (id: string, roomName: string) => {
    if (!confirm(`Are you sure you want to cancel the booking for ${roomName}?\n\nThis action cannot be undone.`)) return;
    
    setCancellingId(id);
    try {
      await auth.cancelBooking(id);
      setBookings(bookings.filter(booking => booking.booking_id !== id));
      showAlert('success', 'Booking cancelled successfully');
    } catch (error: any) {
      showAlert('error', error.response?.data?.error || 'Failed to cancel booking');
    } finally {
      setCancellingId(null);
    }
  };

  if (loading) {
    return (
      <div className="page">
        <div className="card max-w-2xl mx-auto">
          <div className="flex items-center justify-center gap-2 py-8">
            <div className="spinner"></div>
            <span>Loading your bookings...</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-4">
          <h1 className="page-title mb-0">My Bookings</h1>
          <button 
            onClick={() => navigate('/rooms/search')}
            className="btn-primary btn-sm"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            New Booking
          </button>
        </div>

        {alert && (
          <div className={alert.type === 'success' ? 'alert-success mb-4' : 'alert-error mb-4'}>
            {alert.message}
          </div>
        )}
        
        <div className="space-y-3">
          {bookings.length === 0 && (
            <div className="card text-center py-12">
              <svg className="w-20 h-20 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <p className="text-gray-600 text-lg mb-2">No bookings yet</p>
              <p className="text-gray-500 text-sm mb-4">Start by searching for available rooms</p>
              <button 
                onClick={() => navigate('/rooms/search')}
                className="btn-primary"
              >
                Search Rooms
              </button>
            </div>
          )}
          
          {bookings.map((b) => {
            const { date, time} = formatDateTime(b.start_time);
            const isProcessing = cancellingId === b.booking_id;
            
            return (
              <div key={b.booking_id} className={isProcessing ? 'opacity-50 pointer-events-none' : ''}>
                <BookingCard
                  booking={{ 
                    id: b.booking_id, 
                    room: b.room_name || b.room_id || '—', 
                    date, 
                    time, 
                    status: b.status || '—' 
                  }}
                  onCancel={(id) => handleCancel(id, b.room_name || b.room_id || 'this room')}
                  onEdit={(id) => navigate(`/bookings/${id}/edit`)}
                  onTransfer={(id) => navigate(`/bookings/${id}/transfer`)}
                />
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

