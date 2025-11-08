import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';

export default function RoomSchedule() {
  const { id } = useParams();
  const [date, setDate] = useState<string>('');
  const [schedule, setSchedule] = useState<any[]>([]);

  useEffect(() => {
    (async () => {
      if (!id) return;
      const data = await rooms.getRoomSchedule(id, { date: date || undefined });
      setSchedule(data.schedule ?? []);
    })();
  }, [id, date]);

  return (
    <div className="page">
      <h1 className="page-title">Room Schedule</h1>
      <div className="card mb-4 flex items-center gap-2">
        <span className="text-sm text-gray-600">Room ID:</span>
        <span className="font-medium">{id}</span>
        <input value={date} onChange={e => setDate(e.target.value)} type="date" className="border rounded px-3 py-2 ml-auto" />
      </div>
      <div className="space-y-2">
        {schedule.length === 0 && <div className="text-gray-600">No events for selected date.</div>}
        {schedule.map((s) => (
          <div key={s.id} className="card">
            <div className="font-medium">{s.title || 'Reserved'}</div>
            <div className="text-sm text-gray-600">{s.date || '—'} • {s.time || '—'}</div>
            <div className="text-xs text-gray-500">Owner: {s.owner || '—'}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

