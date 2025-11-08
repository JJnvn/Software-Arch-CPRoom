import { FormEvent, useState } from 'react';
import * as rooms from '@/services/rooms';
import { Link } from 'react-router-dom';

export default function SearchRooms() {
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [capacity, setCapacity] = useState<number | ''>('');
  const [features, setFeatures] = useState('');
  const [results, setResults] = useState<any[]>([]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const params = { date, time, capacity: Number(capacity) || undefined, features: features ? features.split(',').map(s => s.trim()) : undefined };
    const data = await rooms.searchRooms(params);
    setResults(data.rooms ?? []);
  }

  return (
    <div className="page">
      <h1 className="page-title">Search Rooms</h1>
      <form onSubmit={onSubmit} className="card grid grid-cols-1 md:grid-cols-5 gap-3 mb-4">
        <input value={date} onChange={e => setDate(e.target.value)} type="date" className="border rounded px-3 py-2" />
        <input value={time} onChange={e => setTime(e.target.value)} type="time" className="border rounded px-3 py-2" />
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
              <Link className="px-3 py-1 rounded bg-green-600 text-white" to="/bookings/create">Book</Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

