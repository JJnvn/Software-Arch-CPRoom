import { FormEvent, useState } from 'react';
import * as rooms from '@/services/rooms';
import { Link } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export default function SearchRooms() {
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [capacity, setCapacity] = useState<number | ''>('');
  const [features, setFeatures] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [duration, setDuration] = useState<number | ''>('');
  const [status, setStatus] = useState<string | null>(null);
  const { user } = useAuth();

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!date || !time || duration === '' || Number(duration) <= 0) {
      setResults([]);
      return;
    }

    const start = new Date(`${date}T${time}`);
    const end = new Date(start.getTime() + Number(duration) * 60000);

    const params = {
      start: start.toISOString(),
      end: end.toISOString(),
      capacity: Number(capacity) || undefined,
      features: features ? features.split(',').map(s => s.trim()) : undefined,
    };
    const data = await rooms.searchRooms(params);
    setResults(data.rooms ?? []);
  }

  async function handleBook(room: any) {
    if (!user) {
      setStatus('Please log in to book');
      setTimeout(() => setStatus(null), 2000);
      return;
    }

    const token = localStorage.getItem('AUTH_TOKEN');
    if (!token) {
      setStatus('Please log in to book');
      setTimeout(() => setStatus(null), 2000);
      return;
    }

    if (!date || !time || duration === '' || Number(duration) <= 0) {
      setStatus('Provide date, time and duration');
      setTimeout(() => setStatus(null), 2000);
      return;
    }

    const start = new Date(`${date}T${time}`);
    const end = new Date(start.getTime() + Number(duration) * 60000);

    const payload = {
      user_id: user.id,
      room_id: room.id,
      start_time: start.toISOString(),
      end_time: end.toISOString(),
    };

    try {
      await rooms.createBooking(payload);
      setStatus(`Booked room "${room.name || room.id}" successfully`);
      setTimeout(() => setStatus(null), 2500);
    } catch (err) {
      setStatus('Failed to create booking');
      setTimeout(() => setStatus(null), 2500);
    }
  }

  return (
    <div className="page">
      <h1 className="page-title">Search Rooms</h1>
      {status && <div className="mb-3 text-sm text-green-700">{status}</div>}
      <form onSubmit={onSubmit} className="card grid grid-cols-1 md:grid-cols-6 gap-3 mb-4">
        <input value={date} onChange={e => setDate(e.target.value)} type="date" className="border rounded px-3 py-2" />
        <input value={time} onChange={e => setTime(e.target.value)} type="time" className="border rounded px-3 py-2" />
        <input value={duration} onChange={e => setDuration(e.target.value === '' ? '' : Number(e.target.value))} type="number" min={15} step={15} className="border rounded px-3 py-2" placeholder="Duration (min)" />
        <input value={capacity} onChange={e => setCapacity(e.target.value === '' ? '' : Number(e.target.value))} type="number" min={1} className="border rounded px-3 py-2" placeholder="Capacity" />
        <input value={features} onChange={e => setFeatures(e.target.value)} className="border rounded px-3 py-2" placeholder="Features (comma-separated)" />
        <button className="bg-blue-600 text-white rounded px-4">Search</button>
      </form>

      <div className="space-y-2">
        {results.length === 0 && <div className="text-gray-600">No rooms found. Try a different search.</div>}
        {results.map((r) => (
          <div key={r.id} className="card flex items-center justify-between">
            <div>
              <div className="font-medium">Room {r.name || r.id}</div>
              <div className="text-sm text-gray-600">Capacity: {r.capacity ?? '—'}</div>
            </div>
            <div className="flex items-center gap-2">
              <Link className="px-3 py-1 rounded bg-indigo-600 text-white" to={`/rooms/${r.id}/schedule`}>View Schedule</Link>
              <button type="button" onClick={() => handleBook(r)} className="px-3 py-1 rounded bg-green-600 text-white">Book Now</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

