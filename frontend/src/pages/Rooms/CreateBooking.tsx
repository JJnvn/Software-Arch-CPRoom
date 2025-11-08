import { FormEvent, useState } from 'react';
import * as rooms from '@/services/rooms';

export default function CreateBooking() {
  const [roomId, setRoomId] = useState('');
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [duration, setDuration] = useState<number | ''>('');
  const [status, setStatus] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    await rooms.createBooking({ roomId, date, time, duration: Number(duration) || 0 });
    setStatus('Booking created');
    setTimeout(() => setStatus(null), 2000);
  }

  return (
    <div className="page">
      <h1 className="page-title">Create Booking</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-green-700 text-sm">{status}</div>}
        <div>
          <label className="block text-sm mb-1">Room ID</label>
          <input value={roomId} onChange={e => setRoomId(e.target.value)} className="w-full border rounded px-3 py-2" placeholder="e.g., 101" required />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div>
            <label className="block text-sm mb-1">Date</label>
            <input value={date} onChange={e => setDate(e.target.value)} type="date" className="w-full border rounded px-3 py-2" required />
          </div>
          <div>
            <label className="block text-sm mb-1">Start Time</label>
            <input value={time} onChange={e => setTime(e.target.value)} type="time" className="w-full border rounded px-3 py-2" required />
          </div>
          <div>
            <label className="block text-sm mb-1">Duration (min)</label>
            <input value={duration} onChange={e => setDuration(e.target.value === '' ? '' : Number(e.target.value))} type="number" min={15} step={15} className="w-full border rounded px-3 py-2" />
          </div>
        </div>
        <button className="px-4 py-2 bg-green-600 text-white rounded">Create Booking</button>
      </form>
    </div>
  );
}

