import { useState, useEffect } from 'react';
import * as staff from '@/services/staff';

export default function ApprovalAuditTrail() {
  const [auditLogs, setAuditLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadAuditTrail();
  }, []);

  const loadAuditTrail = async () => {
    setLoading(true);
    try {
      // Combine both approved and denied bookings to show audit trail
      const approved = await staff.getApprovedApprovals();
      const pending = await staff.getPendingApprovals();
      
      // In a real system, you'd have a dedicated audit endpoint
      // For now, we'll show approved bookings as audit records
      const combined = [
        ...approved.requests.map((r: any) => ({ ...r, action: 'approved' })),
        ...pending.requests.map((r: any) => ({ ...r, action: 'pending' }))
      ];
      
      setAuditLogs(combined);
    } catch (error) {
      console.error('Failed to load audit trail:', error);
    } finally {
      setLoading(false);
    }
  };

  const getActionColor = (action: string) => {
    switch (action) {
      case 'approved': return 'bg-green-100 text-green-800';
      case 'denied': return 'bg-red-100 text-red-800';
      case 'pending': return 'bg-yellow-100 text-yellow-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="page">
      <h1 className="page-title">Approval Audit Trail</h1>
      
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold">Recent Approval Actions</h2>
          <button onClick={loadAuditTrail} className="btn-secondary text-sm">
            Refresh
          </button>
        </div>

        {loading ? (
          <div className="text-center py-4">Loading audit trail...</div>
        ) : auditLogs.length === 0 ? (
          <div className="text-gray-600 py-4 text-center">
            No audit records found
          </div>
        ) : (
          <div className="space-y-2">
            {auditLogs.map((log: any, idx: number) => (
              <div key={idx} className="border rounded p-3">
                <div className="flex items-center justify-between mb-2">
                  <div className="font-medium">
                    Booking #{log.id?.slice(0, 8)}
                  </div>
                  <span className={`text-xs px-2 py-1 rounded ${getActionColor(log.action)}`}>
                    {log.action}
                  </span>
                </div>
                <div className="text-sm text-gray-600 space-y-1">
                  <div>Room: {log.room}</div>
                  <div>Requester: {log.requester}</div>
                  <div>Date: {log.date} at {log.time}</div>
                  {log.raw?.approved_by && (
                    <div>Approved by: {log.raw.approved_by}</div>
                  )}
                  {log.raw?.approved_at && (
                    <div>Action taken: {new Date(log.raw.approved_at).toLocaleString()}</div>
                  )}
                  {log.raw?.denial_reason && (
                    <div className="text-red-600">Reason: {log.raw.denial_reason}</div>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
