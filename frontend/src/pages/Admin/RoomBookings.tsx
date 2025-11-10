import { useState, useEffect } from 'react';
import * as auth from '@/services/auth';

export default function AdminRoomBookings() {
  const [rooms, setRooms] = useState<any[]>([]);
  const [selectedRoom, setSelectedRoom] = useState('');
  const [bookings, setBookings] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadRooms();
  }, []);

  const loadRooms = async () => {
    try {
      // This would typically come from a rooms API endpoint
      // For now, we'll use a placeholder
      setRooms([
        { id: '11111111-1111-1111-1111-111111111111', name: 'Conference Room A' },
        { id: '22222222-2222-2222-2222-222222222222', name: 'Meeting Room B' },
        { id: '33333333-3333-3333-3333-333333333333', name: 'Board Room' },
      ]);
    } catch (error) {
      console.error('Failed to load rooms:', error);
    }
  };

  const loadBookings = async (roomId: string) => {
    if (!roomId) return;
    
    setLoading(true);
    try {
      // Call admin endpoint to get all bookings for a room
      const response = await auth.getAdminRoomBookings(roomId);
      setBookings(response.bookings || []);
    } catch (error) {
      console.error('Failed to load bookings:', error);
      setBookings([]);
    } finally {
      setLoading(false);
    }
  };

  const handleRoomChange = (roomId: string) => {
    setSelectedRoom(roomId);
    loadBookings(roomId);
  };

  const formatDateTime = (isoString: string) => {
    const d = new Date(isoString);
    return {
      date: d.toLocaleDateString(),
      time: d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    };
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'confirmed': return 'bg-green-100 text-green-800';
      case 'pending': return 'bg-yellow-100 text-yellow-800';
      case 'cancelled': return 'bg-red-100 text-red-800';
      case 'completed': return 'bg-blue-100 text-blue-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="page">
      <h1 className="page-title">Admin - Room Bookings</h1>
      
      <div className="card mb-4">
        <label className="block text-sm font-medium mb-2">Select Room</label>
        <select
          value={selectedRoom}
          onChange={(e) => handleRoomChange(e.target.value)}
          className="input"
        >
          <option value="">-- Select a room --</option>
          {rooms.map((room) => (
            <option key={room.id} value={room.id}>
              {room.name}
            </option>
          ))}
        </select>
      </div>

      {selectedRoom && (
        <div className="card">
          <h2 className="font-semibold mb-3">
            Bookings for {rooms.find(r => r.id === selectedRoom)?.name}
          </h2>
          
          {loading ? (
            <div className="text-center py-4">Loading bookings...</div>
          ) : bookings.length === 0 ? (
            <div className="text-gray-600 py-4 text-center">
              No bookings found for this room
            </div>
          ) : (
            <div className="space-y-2">
              {bookings.map((booking: any) => {
                const start = formatDateTime(booking.start_time);
                const end = formatDateTime(booking.end_time);
                
                return (
                  <div key={booking.booking_id} className="border rounded p-3">
                    <div className="flex items-center justify-between mb-2">
                      <div className="font-medium">
                        {start.date} • {start.time} - {end.time}
                      </div>
                      <span className={`text-xs px-2 py-1 rounded ${getStatusColor(booking.status)}`}>
                        {booking.status}
                      </span>
                    </div>
                    <div className="text-sm text-gray-600 space-y-1">
                      <div>Booking ID: {booking.booking_id}</div>
                      <div>User: {booking.user_name || booking.user_id}</div>
                      <div>Created: {new Date(booking.created_at).toLocaleString()}</div>
                      {booking.updated_at !== booking.created_at && (
                        <div>Updated: {new Date(booking.updated_at).toLocaleString()}</div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
