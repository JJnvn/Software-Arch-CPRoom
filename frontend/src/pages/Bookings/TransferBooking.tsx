import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import * as auth from '@/services/auth';

export default function TransferBooking() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [booking, setBooking] = useState<any>(null);
  const [newUserEmail, setNewUserEmail] = useState('');
  const [status, setStatus] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    (async () => {
      try {
        const history = await auth.getBookingHistory();
        const found = history.find((b: any) => b.booking_id === id);
        if (found) {
          setBooking(found);
        } else {
          setError('Booking not found');
        }
      } catch (error) {
        setError('Failed to load booking');
      }
    })();
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!id) {
      setError('Invalid booking ID');
      return;
    }

    if (!newUserEmail.trim()) {
      setError('Please enter a user email');
      return;
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(newUserEmail)) {
      setError('Please enter a valid email address');
      return;
    }

    setIsSubmitting(true);
    setError('');
    setStatus('Transferring booking...');
    
    try {
      await auth.transferBooking(id, { new_user_email: newUserEmail });
      setStatus('Booking transferred successfully!');
      setTimeout(() => navigate('/booking-history'), 2000);
    } catch (error: any) {
      setError(error.response?.data?.error || 'Failed to transfer booking');
      setStatus('');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!booking) {
    return (
      <div className="page">
        <div className="card max-w-2xl mx-auto">
          <div className="flex items-center justify-center gap-2">
            {error ? (
              <div className="alert-error">{error}</div>
            ) : (
              <>
                <div className="spinner"></div>
                <span>Loading booking details...</span>
              </>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="max-w-2xl mx-auto">
        <button 
          onClick={() => navigate('/booking-history')} 
          className="btn-secondary btn-sm mb-4"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to Bookings
        </button>

        <h1 className="page-title">Transfer Booking</h1>
        
        <div className="card mb-4 bg-blue-50 border border-blue-200">
          <h3 className="font-semibold text-blue-900 mb-2">Current Booking Details</h3>
          <div className="space-y-2 text-sm text-blue-800">
            <div className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
              <strong>Room:</strong> {booking.room_name || booking.room_id}
            </div>
            <div className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <strong>Date:</strong> {new Date(booking.start_time).toLocaleDateString()}
            </div>
            <div className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Time:</strong> {new Date(booking.start_time).toLocaleTimeString()} - {new Date(booking.end_time).toLocaleTimeString()}
            </div>
            <div className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Status:</strong> {booking.status}
            </div>
          </div>
        </div>
        
        <div className="card">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-semibold mb-2">
                New Owner Email Address
              </label>
              <input
                type="email"
                value={newUserEmail}
                onChange={(e) => {
                  setNewUserEmail(e.target.value);
                  setError(''); // Clear error on input
                }}
                className={error && !status ? 'input-error' : 'input'}
                placeholder="user@example.com"
                disabled={isSubmitting}
                required
              />
              <p className="text-xs text-gray-500 mt-2 flex items-start gap-1">
                <svg className="w-4 h-4 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Enter the email address of the person you want to transfer this booking to. They will receive a notification.
              </p>
            </div>

            {error && (
              <div className="alert-error flex items-start gap-2">
                <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {error}
              </div>
            )}

            {status && (
              <div className="alert-success flex items-start gap-2">
                <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {status}
              </div>
            )}

            <div className="flex gap-3 pt-2">
              <button 
                type="submit" 
                className="btn-primary flex-1"
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <div className="spinner"></div>
                    Transferring...
                  </>
                ) : (
                  <>
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                    </svg>
                    Transfer Booking
                  </>
                )}
              </button>
              <button 
                type="button" 
                onClick={() => navigate('/booking-history')} 
                className="btn-secondary"
                disabled={isSubmitting}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
