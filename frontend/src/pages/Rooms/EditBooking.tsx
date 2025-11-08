import { FormEvent, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';

export default function EditBooking() {
  const { id } = useParams();
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    // Could load booking details for id
  }, [id]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!id) return;
    await rooms.rescheduleBooking(id, { date, time });
    setStatus('Booking updated');
    setTimeout(() => setStatus(null), 2000);
  }

  return (
    <div className="page">
      <h1 className="page-title">Edit Booking</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-green-700 text-sm">{status}</div>}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div>
            <label className="block text-sm mb-1">New Date</label>
            <input value={date} onChange={e => setDate(e.target.value)} type="date" className="w-full border rounded px-3 py-2" />
          </div>
          <div>
            <label className="block text-sm mb-1">New Time</label>
            <input value={time} onChange={e => setTime(e.target.value)} type="time" className="w-full border rounded px-3 py-2" />
          </div>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded">Save Changes</button>
      </form>
    </div>
  );
}

