import { useEffect, useState } from 'react';
import * as staff from '@/services/staff';

export default function PendingApprovals() {
  const [requests, setRequests] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [processingId, setProcessingId] = useState<string | null>(null);
  const [denyingId, setDenyingId] = useState<string | null>(null);
  const [denyReason, setDenyReason] = useState('');
  const [alert, setAlert] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  useEffect(() => {
    loadRequests();
  }, []);

  async function loadRequests() {
    setLoading(true);
    try {
      const data = await staff.getPendingApprovals();
      setRequests(data.requests ?? []);
    } catch (error) {
      showAlert('error', 'Failed to load pending requests');
    } finally {
      setLoading(false);
    }
  }

  function showAlert(type: 'success' | 'error', message: string) {
    setAlert({ type, message });
    setTimeout(() => setAlert(null), 5000);
  }

  async function approve(id: string, roomInfo: string) {
    if (!confirm(`Approve booking for ${roomInfo}?`)) return;
    
    setProcessingId(id);
    try {
      await staff.approveBookingRequest(id);
      setRequests((r) => r.filter((x) => x.id !== id));
      showAlert('success', 'Booking approved successfully');
    } catch (e) {
      showAlert('error', 'Failed to approve booking');
    } finally {
      setProcessingId(null);
    }
  }

  function startDeny(id: string) {
    setDenyingId(id);
    setDenyReason('');
  }

  function cancelDeny() {
    setDenyingId(null);
    setDenyReason('');
  }

  async function confirmDeny(id: string, roomInfo: string) {
    if (!denyReason.trim()) {
      showAlert('error', 'Please provide a reason for denial');
      return;
    }

    setProcessingId(id);
    try {
      await staff.denyBookingRequest(id, denyReason);
      setRequests((r) => r.filter((x) => x.id !== id));
      showAlert('success', 'Booking denied');
      setDenyingId(null);
      setDenyReason('');
    } catch (e) {
      showAlert('error', 'Failed to deny booking');
    } finally {
      setProcessingId(null);
    }
  }

  if (loading) {
    return (
      <div className="page">
        <div className="card max-w-2xl mx-auto">
          <div className="flex items-center justify-center gap-2">
            <div className="spinner"></div>
            <span>Loading pending approvals...</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="max-w-4xl mx-auto">
        <h1 className="page-title">Pending Booking Approvals</h1>
        
        {alert && (
          <div className={alert.type === 'success' ? 'alert-success mb-4' : 'alert-error mb-4'}>
            {alert.message}
          </div>
        )}
        
        <div className="space-y-3">
          {requests.length === 0 && (
            <div className="card text-center py-8">
              <svg className="w-16 h-16 mx-auto text-gray-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <p className="text-gray-600 text-lg">No pending approval requests</p>
              <p className="text-gray-500 text-sm mt-1">New booking requests will appear here</p>
            </div>
          )}
          
          {requests.map((r) => {
            const roomInfo = `Room ${r.room || '—'}`;
            const isProcessing = processingId === r.id;
            const isDenying = denyingId === r.id;

            return (
              <div key={r.id} className="card hover:shadow-lg transition-shadow duration-200">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                      </svg>
                      <span className="font-semibold text-lg">{roomInfo}</span>
                    </div>
                    <div className="space-y-1 text-sm text-gray-600">
                      <div className="flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>{r.date || '—'}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>{r.time || '—'}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                        <span>Requested by: {r.requester || '—'}</span>
                      </div>
                    </div>
                  </div>
                  
                  {!isDenying && (
                    <div className="flex items-center gap-2">
                      <button 
                        className="btn-sm btn-success"
                        onClick={() => approve(r.id, roomInfo)}
                        disabled={isProcessing}
                      >
                        {isProcessing ? (
                          <div className="spinner"></div>
                        ) : (
                          <>
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                            Approve
                          </>
                        )}
                      </button>
                      <button 
                        className="btn-sm btn-danger"
                        onClick={() => startDeny(r.id)}
                        disabled={isProcessing}
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                        Deny
                      </button>
                    </div>
                  )}
                </div>

                {isDenying && (
                  <div className="mt-4 pt-4 border-t space-y-3">
                    <div>
                      <label className="block text-sm font-semibold mb-2">Reason for Denial</label>
                      <textarea
                        value={denyReason}
                        onChange={(e) => setDenyReason(e.target.value)}
                        className="input resize-none"
                        rows={3}
                        placeholder="Please provide a reason for denying this booking request..."
                        autoFocus
                      />
                    </div>
                    <div className="flex gap-2">
                      <button
                        className="btn-sm btn-danger flex-1"
                        onClick={() => confirmDeny(r.id, roomInfo)}
                        disabled={isProcessing || !denyReason.trim()}
                      >
                        {isProcessing ? (
                          <>
                            <div className="spinner"></div>
                            Denying...
                          </>
                        ) : (
                          'Confirm Denial'
                        )}
                      </button>
                      <button
                        className="btn-sm btn-secondary"
                        onClick={cancelDeny}
                        disabled={isProcessing}
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

