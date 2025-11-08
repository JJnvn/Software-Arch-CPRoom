import { FormEvent, useState } from 'react';
import { useParams } from 'react-router-dom';
import * as rooms from '@/services/rooms';

export default function TransferBooking() {
  const { id } = useParams();
  const [newOwnerEmail, setNewOwnerEmail] = useState('');
  const [status, setStatus] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!id) return;
    await rooms.transferBookingOwnership(id, { newOwnerEmail });
    setStatus('Ownership transferred');
    setTimeout(() => setStatus(null), 2000);
  }

  return (
    <div className="page">
      <h1 className="page-title">Transfer Booking</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-green-700 text-sm">{status}</div>}
        <div>
          <label className="block text-sm mb-1">New Owner Email</label>
          <input value={newOwnerEmail} onChange={e => setNewOwnerEmail(e.target.value)} type="email" className="w-full border rounded px-3 py-2" placeholder="new.owner@company.com" required />
        </div>
        <button className="px-4 py-2 bg-amber-600 text-white rounded">Transfer Ownership</button>
      </form>
    </div>
  );
}

