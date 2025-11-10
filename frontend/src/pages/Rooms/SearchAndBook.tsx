import { FormEvent, useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import * as rooms from '@/services/rooms';
import { useAuth } from '@/hooks/useAuth';

interface Room {
  id: string;
  name: string;
  capacity?: number;
  features?: string[];
}

export default function SearchAndBook() {
  const { user } = useAuth();
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [duration, setDuration] = useState<number>(60);
  const [capacity, setCapacity] = useState<number | ''>('');
  const [features, setFeatures] = useState('');
  const [results, setResults] = useState<Room[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
  const [isBooking, setIsBooking] = useState(false);
  const [notification, setNotification] = useState<{ type: 'success' | 'error' | 'info'; message: string } | null>(null);

  // Set default date and time to now + 1 hour
  useEffect(() => {
    const now = new Date();
    now.setHours(now.getHours() + 1, 0, 0, 0);
    const dateStr = now.toISOString().split('T')[0];
    const timeStr = now.toTimeString().substring(0, 5);
    setDate(dateStr);
    setTime(timeStr);
  }, []);

  const showNotification = (type: 'success' | 'error' | 'info', message: string) => {
    setNotification({ type, message });
    setTimeout(() => setNotification(null), 5000);
  };

  const handleSearch = async (e: FormEvent) => {
    e.preventDefault();

    if (!date || !time || !duration || duration <= 0) {
      showNotification('error', 'Please fill in date, time, and duration');
      return;
    }

    setIsSearching(true);
    setSelectedRoom(null);

    try {
      const start = new Date(`${date}T${time}`);
      const end = new Date(start.getTime() + duration * 60000);

      // Validate booking is in the future
      if (start < new Date()) {
        showNotification('error', 'Booking time must be in the future');
        setIsSearching(false);
        return;
      }

      const params = {
        start: start.toISOString(),
        end: end.toISOString(),
        capacity: capacity ? Number(capacity) : undefined,
        features: features ? features.split(',').map(s => s.trim()).filter(Boolean) : undefined,
      };

      const data = await rooms.searchRooms(params);
      setResults(data.rooms || []);

      if (data.rooms.length === 0) {
        showNotification('info', 'No rooms available for the selected time slot. Try different dates or filters.');
      }
    } catch (err) {
      console.error('Search error:', err);
      showNotification('error', 'Failed to search rooms. Please try again.');
      setResults([]);
    } finally {
      setIsSearching(false);
    }
  };

  const handleBookRoom = async (room: Room) => {
    if (!user) {
      showNotification('error', 'Please log in to book a room');
      return;
    }

    if (!date || !time || !duration) {
      showNotification('error', 'Please select date, time, and duration first');
      return;
    }

    setIsBooking(true);
    setSelectedRoom(room);

    try {
      const start = new Date(`${date}T${time}`);
      const end = new Date(start.getTime() + duration * 60000);

      const payload = {
        user_id: user.id,
        room_id: room.id,
        start_time: start.toISOString(),
        end_time: end.toISOString(),
      };

      await rooms.createBooking(payload);
      showNotification('success', `Successfully booked ${room.name}! Check your booking history.`);
      
      // Refresh search results to show updated availability
      setTimeout(() => {
        const fakeEvent = { preventDefault: () => {} } as FormEvent;
        handleSearch(fakeEvent);
      }, 1000);
    } catch (err: any) {
      console.error('Booking error:', err);
      const errorMsg = err.response?.data?.error || 'Failed to create booking. The room may no longer be available.';
      showNotification('error', errorMsg);
    } finally {
      setIsBooking(false);
      setSelectedRoom(null);
    }
  };

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours === 0) return `${mins}m`;
    if (mins === 0) return `${hours}h`;
    return `${hours}h ${mins}m`;
  };

  const getStartEndTime = () => {
    if (!date || !time || !duration) return null;
    const start = new Date(`${date}T${time}`);
    const end = new Date(start.getTime() + duration * 60000);
    return {
      start: start.toLocaleString('en-US', { 
        weekday: 'short', 
        month: 'short', 
        day: 'numeric', 
        hour: '2-digit', 
        minute: '2-digit' 
      }),
      end: end.toLocaleTimeString('en-US', { 
        hour: '2-digit', 
        minute: '2-digit' 
      }),
    };
  };

  const timeInfo = getStartEndTime();

  return (
    <div className="page">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Search & Book Rooms</h1>
          <p className="text-gray-600">Find available rooms and book instantly</p>
        </div>

        {/* Notification */}
        {notification && (
          <div className={`mb-6 rounded-lg p-4 flex items-start gap-3 ${
            notification.type === 'success' ? 'bg-green-50 border border-green-200' :
            notification.type === 'error' ? 'bg-red-50 border border-red-200' :
            'bg-blue-50 border border-blue-200'
          }`}>
            <svg className={`w-5 h-5 flex-shrink-0 ${
              notification.type === 'success' ? 'text-green-600' :
              notification.type === 'error' ? 'text-red-600' :
              'text-blue-600'
            }`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {notification.type === 'success' ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              ) : notification.type === 'error' ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              )}
            </svg>
            <p className={`text-sm font-medium ${
              notification.type === 'success' ? 'text-green-800' :
              notification.type === 'error' ? 'text-red-800' :
              'text-blue-800'
            }`}>{notification.message}</p>
          </div>
        )}

        {/* Search Form */}
        <form onSubmit={handleSearch} className="card mb-6">
          <div className="flex items-center gap-2 mb-4">
            <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <h2 className="text-lg font-semibold text-gray-900">Search Criteria</h2>
          </div>

          {/* Time Selection */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
              <input
                type="date"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Start Time</label>
              <input
                type="time"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Duration</label>
              <select
                value={duration}
                onChange={(e) => setDuration(Number(e.target.value))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              >
                <option value={30}>30 minutes</option>
                <option value={60}>1 hour</option>
                <option value={90}>1.5 hours</option>
                <option value={120}>2 hours</option>
                <option value={180}>3 hours</option>
                <option value={240}>4 hours</option>
                <option value={480}>8 hours (Full day)</option>
              </select>
            </div>
          </div>

          {/* Filters */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Minimum Capacity <span className="text-gray-500">(optional)</span>
              </label>
              <input
                type="number"
                value={capacity}
                onChange={(e) => setCapacity(e.target.value === '' ? '' : Number(e.target.value))}
                min={1}
                placeholder="e.g., 10"
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Required Features <span className="text-gray-500">(comma-separated)</span>
              </label>
              <input
                type="text"
                value={features}
                onChange={(e) => setFeatures(e.target.value)}
                placeholder="e.g., Projector, Whiteboard, Video Conference"
                className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* Booking Summary */}
          {timeInfo && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-4">
              <div className="flex items-center gap-2 text-sm text-blue-800">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="font-medium">Booking for:</span>
                <span>{timeInfo.start} - {timeInfo.end}</span>
                <span className="text-blue-600">({formatDuration(duration)})</span>
              </div>
            </div>
          )}

          <button
            type="submit"
            disabled={isSearching}
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-3 rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            {isSearching ? (
              <>
                <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Searching...
              </>
            ) : (
              <>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                Search Available Rooms
              </>
            )}
          </button>
        </form>

        {/* Results */}
        {results.length > 0 && (
          <div>
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">
                Available Rooms <span className="text-gray-500 font-normal">({results.length})</span>
              </h2>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {results.map((room) => (
                <div
                  key={room.id}
                  className="card hover:shadow-lg transition-shadow border-2 border-transparent hover:border-blue-200"
                >
                  {/* Room Header */}
                  <div className="flex items-start justify-between mb-3">
                    <div className="flex items-center gap-2">
                      <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                        <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                        </svg>
                      </div>
                      <div>
                        <h3 className="font-semibold text-gray-900">{room.name}</h3>
                        <p className="text-xs text-gray-500">ID: {room.id.substring(0, 8)}</p>
                      </div>
                    </div>
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Available
                    </span>
                  </div>

                  {/* Room Details */}
                  <div className="space-y-2 mb-4">
                    {room.capacity && (
                      <div className="flex items-center gap-2 text-sm text-gray-600">
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                        <span>Capacity: <span className="font-medium">{room.capacity} people</span></span>
                      </div>
                    )}
                    {room.features && room.features.length > 0 && (
                      <div>
                        <div className="flex items-center gap-2 text-sm text-gray-600 mb-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span className="font-medium">Features:</span>
                        </div>
                        <div className="flex flex-wrap gap-1 ml-6">
                          {room.features.slice(0, 4).map((feature, idx) => (
                            <span
                              key={idx}
                              className="inline-block px-2 py-0.5 bg-gray-100 text-gray-700 text-xs rounded"
                            >
                              {feature}
                            </span>
                          ))}
                          {room.features.length > 4 && (
                            <span className="inline-block px-2 py-0.5 bg-gray-100 text-gray-500 text-xs rounded">
                              +{room.features.length - 4} more
                            </span>
                          )}
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2 pt-3 border-t border-gray-200">
                    <Link
                      to={`/rooms/${room.id}/schedule`}
                      className="flex-1 text-center px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                    >
                      View Schedule
                    </Link>
                    <button
                      onClick={() => handleBookRoom(room)}
                      disabled={isBooking && selectedRoom?.id === room.id}
                      className="flex-1 px-3 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white rounded-lg text-sm font-medium transition-colors flex items-center justify-center gap-1"
                    >
                      {isBooking && selectedRoom?.id === room.id ? (
                        <>
                          <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                          </svg>
                          Booking...
                        </>
                      ) : (
                        <>
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                          </svg>
                          Book Now
                        </>
                      )}
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {results.length === 0 && !isSearching && date && time && (
          <div className="text-center py-12">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-1">No rooms found</h3>
            <p className="text-gray-600 mb-4">Try adjusting your search criteria or selecting a different time</p>
          </div>
        )}
      </div>
    </div>
  );
}
