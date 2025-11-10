import { FormEvent, useState } from 'react';
import * as rooms from '@/services/rooms';
import { useAuth } from '@/hooks/useAuth';

export default function AdminCreateRoom() {
  const { user } = useAuth();
  const role = (user?.role || '').toUpperCase();
  const isPrivileged = role === 'ADMIN' || role === 'STAFF';
  const [name, setName] = useState('');
  const [capacity, setCapacity] = useState<number | ''>('');
  const [features, setFeatures] = useState('');
  const [status, setStatus] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!isPrivileged) return;
    if (!name.trim()) return;
    try {
      const payload = {
        name: name.trim(),
        capacity: capacity === '' ? undefined : Number(capacity),
        features: features
          ? features.split(',').map((s) => s.trim()).filter(Boolean)
          : undefined,
      };
      const created = await rooms.createRoom(payload);
      setStatus(`Created room "${created.name}" (id ${created.id})`);
      setName('');
      setCapacity('');
      setFeatures('');
      setTimeout(() => setStatus(null), 2000);
    } catch (err: any) {
      const msg = err?.response?.data?.error || 'Failed to create room';
      setStatus(msg);
      setTimeout(() => setStatus(null), 2500);
    }
  }

  if (!isPrivileged) {
    return (
      <div className="page">
        <h1 className="page-title">Create Room</h1>
        <div className="card text-red-700">Not authorized</div>
      </div>
    );
  }

  return (
    <div className="page">
      <h1 className="page-title">Create Room</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-sm">{status}</div>}
        <div>
          <label className="block text-sm mb-1">Room Name</label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border rounded px-3 py-2"
            placeholder="Unique room name"
            required
          />
        </div>
        <div>
          <label className="block text-sm mb-1">Capacity</label>
          <input
            value={capacity}
            onChange={(e) => setCapacity(e.target.value === '' ? '' : Number(e.target.value))}
            type="number"
            min={0}
            className="w-full border rounded px-3 py-2"
            placeholder="Capacity (optional)"
          />
        </div>
        <div>
          <label className="block text-sm mb-1">Features</label>
          <input
            value={features}
            onChange={(e) => setFeatures(e.target.value)}
            className="w-full border rounded px-3 py-2"
            placeholder="Comma-separated features (optional)"
          />
        </div>
        <button className="px-4 py-2 bg-indigo-600 text-white rounded">Create</button>
      </form>
    </div>
  );
}
