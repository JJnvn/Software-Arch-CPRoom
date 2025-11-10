import { useEffect, useState } from 'react';
import * as staff from '@/services/staff';

export default function PendingApprovals() {
  const [requests, setRequests] = useState<any[]>([]);

  useEffect(() => {
    (async () => {
      const data = await staff.getPendingApprovals();
      setRequests(data.requests ?? []);
    })();
  }, []);

  async function approve(id: string) {
    try {
      await staff.approveBookingRequest(id);
      setRequests((r) => r.filter((x) => x.id !== id));
    } catch (e) {
      // noop; could show a toast
    }
  }
  async function deny(id: string) {
    const reason = window.prompt('Reason for denial?');
    if (!reason) return;
    try {
      await staff.denyBookingRequest(id, reason);
      setRequests((r) => r.filter((x) => x.id !== id));
    } catch (e) {
      // noop; could show a toast
    }
  }

  return (
    <div className="page">
      <h1 className="page-title">Pending Booking Approvals</h1>
      <div className="space-y-2">
        {requests.length === 0 && <div className="text-gray-600">No pending requests.</div>}
        {requests.map((r) => (
          <div key={r.id} className="card flex items-center justify-between">
            <div>
              <div className="font-medium">Room {r.room || '—'} • {r.date || '—'} {r.time || ''}</div>
              <div className="text-sm text-gray-600">Requested by: {r.requester || '—'}</div>
            </div>
            <div className="flex items-center gap-2">
              <button className="px-3 py-1 rounded bg-green-600 text-white" onClick={() => approve(r.id)}>Approve</button>
              <button className="px-3 py-1 rounded bg-red-600 text-white" onClick={() => deny(r.id)}>Deny</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

