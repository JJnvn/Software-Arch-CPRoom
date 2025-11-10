import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';

interface Booking {
  booking_id: string;
  user_id: string;
  start_time: string;
  end_time: string;
  status: string;
}

interface RoomScheduleData {
  room_id: string;
  room_name: string;
  date: string;
  bookings: Booking[];
}

export default function RoomSchedule() {
  const { id } = useParams();
  const [date, setDate] = useState<string>(new Date().toISOString().split('T')[0]);
  const [scheduleData, setScheduleData] = useState<RoomScheduleData | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    (async () => {
      if (!id) return;
      setLoading(true);
      try {
        const data = await rooms.getRoomSchedule(id, { date: date || undefined });
        setScheduleData(data);
      } catch (error) {
        console.error('Failed to load schedule:', error);
      } finally {
        setLoading(false);
      }
    })();
  }, [id, date]);

  const formatTime = (isoString: string) => {
    const dateObj = new Date(isoString);
    return dateObj.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false });
  };

  const formatDate = (isoString: string) => {
    const dateObj = new Date(isoString);
    return dateObj.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'confirmed':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="page">
      <h1 className="page-title">Room Schedule</h1>
      
      <div className="card mb-4">
        <div className="flex items-center justify-between gap-4">
          <div>
            <div className="text-sm text-gray-600">Room</div>
            <div className="font-semibold text-lg">{scheduleData?.room_name || 'Loading...'}</div>
            <div className="text-xs text-gray-500">{scheduleData?.room_id || id}</div>
          </div>
          <div className="flex items-center gap-2">
            <label className="text-sm text-gray-600">Date:</label>
            <input 
              value={date} 
              onChange={e => setDate(e.target.value)} 
              type="date" 
              className="border rounded px-3 py-2" 
            />
          </div>
        </div>
      </div>

      {loading && (
        <div className="text-center text-gray-600 py-8">Loading schedule...</div>
      )}

      {!loading && scheduleData && (
        <div className="space-y-3">
          <div className="text-sm text-gray-600">
            {scheduleData.bookings.length} booking(s) on {formatDate(scheduleData.date)}
          </div>
          
          {scheduleData.bookings.length === 0 && (
            <div className="card text-center text-gray-600 py-8">
              No bookings for this date. Room is available!
            </div>
          )}
          
          {scheduleData.bookings.map((booking) => (
            <div key={booking.booking_id} className="card hover:shadow-md transition-shadow">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(booking.status)}`}>
                      {booking.status.toUpperCase()}
                    </span>
                  </div>
                  <div className="text-lg font-semibold text-gray-800 mb-1">
                    {formatTime(booking.start_time)} - {formatTime(booking.end_time)}
                  </div>
                  <div className="text-sm text-gray-600">
                    Duration: {Math.round((new Date(booking.end_time).getTime() - new Date(booking.start_time).getTime()) / (1000 * 60))} minutes
                  </div>
                </div>
                <div className="text-right text-sm">
                  <div className="text-gray-500 mb-1">Booking ID</div>
                  <div className="text-xs font-mono text-gray-700">{booking.booking_id.substring(0, 8)}...</div>
                  <div className="text-gray-500 mt-2 mb-1">User ID</div>
                  <div className="text-xs font-mono text-gray-700">{booking.user_id.substring(0, 8)}...</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

