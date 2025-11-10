import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';
import { useAuth } from '@/hooks/useAuth';

export default function EditRoom() {
  const { id } = useParams();
  const { user } = useAuth();
  const role = (user?.role || '').toUpperCase();
  const isPrivileged = role === 'ADMIN' || role === 'STAFF';

  const [name, setName] = useState('');
  const [capacity, setCapacity] = useState<number | ''>('');
  const [features, setFeatures] = useState('');
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      if (!id) return;
      try {
        const r = await rooms.getRoom(id);
        setName(r.name || '');
        setCapacity(typeof r.capacity === 'number' ? r.capacity : '');
        setFeatures(Array.isArray(r.features) ? r.features.join(', ') : '');
      } catch (e) {
        setStatus('Failed to load room');
        setTimeout(() => setStatus(null), 2000);
      }
    })();
  }, [id]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!isPrivileged || !id) return;
    try {
      const payload: any = { name: name.trim() };
      if (capacity !== '') payload.capacity = Number(capacity);
      if (features.trim()) payload.features = features.split(',').map(s => s.trim()).filter(Boolean);
      await rooms.updateRoom(id, payload);
      setStatus('Room updated');
      setTimeout(() => setStatus(null), 2000);
    } catch (e) {
      setStatus('Failed to update room');
      setTimeout(() => setStatus(null), 2000);
    }
  }

  if (!isPrivileged) {
    return (
      <div className="page">
        <h1 className="page-title">Edit Room</h1>
        <div className="card text-red-700">Not authorized</div>
      </div>
    );
  }

  return (
    <div className="page">
      <h1 className="page-title">Edit Room</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-sm">{status}</div>}
        <div>
          <label className="block text-sm mb-1">Room Name</label>
          <input value={name} onChange={(e) => setName(e.target.value)} className="w-full border rounded px-3 py-2" required />
        </div>
        <div>
          <label className="block text-sm mb-1">Capacity</label>
          <input value={capacity} onChange={e => setCapacity(e.target.value === '' ? '' : Number(e.target.value))} type="number" min={0} className="w-full border rounded px-3 py-2" />
        </div>
        <div>
          <label className="block text-sm mb-1">Features</label>
          <input value={features} onChange={e => setFeatures(e.target.value)} className="w-full border rounded px-3 py-2" placeholder="Comma-separated" />
        </div>
        <button className="px-4 py-2 bg-indigo-600 text-white rounded">Save</button>
      </form>
    </div>
  );
}

