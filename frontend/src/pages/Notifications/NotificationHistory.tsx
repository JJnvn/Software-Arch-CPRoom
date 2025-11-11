import { useEffect, useState } from 'react';
import * as notifications from '@/services/notifications';
import { useAuth } from '@/hooks/useAuth';

interface NotificationItem {
  id: string;
  type: string;
  message: string;
  channel: string;
  sent_at: string;
  status: string;
  metadata?: Record<string, any>;
}

export default function NotificationHistory() {
  const [items, setItems] = useState<NotificationItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);
  const [hasMore, setHasMore] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    loadHistory();
  }, [user, page]);

  async function loadHistory() {
    if (!user) {
      setIsLoading(false);
      return;
    }

    try {
      setErrorMessage(null);
      const data = await notifications.getNotificationHistory(page, pageSize);
      const history = data.history || [];
      
      setItems(history);
      setHasMore(history.length === pageSize);
    } catch (error: any) {
      console.error('Failed to load notification history:', error);
      setErrorMessage(error.response?.data?.error || 'Failed to load notification history');
    } finally {
      setIsLoading(false);
    }
  }

  function getNotificationIcon(type: string) {
    switch (type.toLowerCase()) {
      case 'booking_created':
      case 'booking_confirmation':
        return (
          <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
      case 'booking_approved':
        return (
          <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
      case 'booking_denied':
        return (
          <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
      case 'booking_cancelled':
        return (
          <svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
      case 'booking_reminder':
        return (
          <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        );
      case 'booking_transferred':
        return (
          <svg className="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
          </svg>
        );
      default:
        return (
          <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
    }
  }

  function getChannelBadge(channel: string) {
    const channelLower = channel.toLowerCase();
    if (channelLower === 'email') {
      return <span className="badge-info">📧 Email</span>;
    } else if (channelLower === 'push') {
      return <span className="badge-warning">📱 Push</span>;
    } else if (channelLower === 'sms') {
      return <span className="badge-success">💬 SMS</span>;
    }
    return <span className="badge-default">{channel}</span>;
  }

  function getStatusBadge(status: string) {
    switch (status.toLowerCase()) {
      case 'sent':
      case 'delivered':
        return <span className="badge-success">✓ Sent</span>;
      case 'pending':
        return <span className="badge-warning">⏳ Pending</span>;
      case 'failed':
        return <span className="badge-danger">✗ Failed</span>;
      default:
        return <span className="badge-default">{status}</span>;
    }
  }

  function formatDate(dateStr: string) {
    try {
      const date = new Date(dateStr);
      const now = new Date();
      const diff = now.getTime() - date.getTime();
      const seconds = Math.floor(diff / 1000);
      const minutes = Math.floor(seconds / 60);
      const hours = Math.floor(minutes / 60);
      const days = Math.floor(hours / 24);

      if (days > 7) {
        return date.toLocaleDateString('en-US', { 
          month: 'short', 
          day: 'numeric', 
          year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined 
        });
      } else if (days > 0) {
        return `${days} day${days > 1 ? 's' : ''} ago`;
      } else if (hours > 0) {
        return `${hours} hour${hours > 1 ? 's' : ''} ago`;
      } else if (minutes > 0) {
        return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
      } else {
        return 'Just now';
      }
    } catch {
      return dateStr;
    }
  }

  function getNotificationTitle(type: string) {
    return type
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  }

  if (isLoading && items.length === 0) {
    return (
      <div className="page">
        <div className="flex items-center justify-center py-12">
          <div className="spinner"></div>
          <span className="ml-3 text-gray-600">Loading notification history...</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="page">
        <div className="card text-center py-8">
          <svg className="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          <p className="text-gray-600">Please log in to view your notification history.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="page-title mb-1">Notification History</h1>
            <p className="text-gray-600">View all notifications sent to you</p>
          </div>
          <button
            onClick={() => {
              setPage(1);
              loadHistory();
            }}
            className="btn-secondary"
            disabled={isLoading}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Refresh</span>
          </button>
        </div>

        {errorMessage && (
          <div className="alert-error mb-4">
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <span>{errorMessage}</span>
          </div>
        )}

        <div className="space-y-3">
          {items.length === 0 ? (
            <div className="card text-center py-12">
              <svg className="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
              <h3 className="text-lg font-medium text-gray-900 mb-2">No notifications yet</h3>
              <p className="text-gray-600">You'll see your notification history here once you receive notifications.</p>
            </div>
          ) : (
            items.map((item) => (
              <div key={item.id} className="card hover:shadow-md transition-shadow">
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0 mt-1">
                    {getNotificationIcon(item.type)}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-4 mb-2">
                      <div>
                        <h3 className="font-semibold text-gray-900">
                          {getNotificationTitle(item.type)}
                        </h3>
                        <p className="text-sm text-gray-600 mt-1">{item.message}</p>
                      </div>
                      <div className="flex-shrink-0 text-right">
                        <div className="text-xs text-gray-500 mb-1">
                          {formatDate(item.sent_at)}
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center gap-2 flex-wrap">
                      {getChannelBadge(item.channel)}
                      {getStatusBadge(item.status)}
                      
                      {item.metadata && Object.keys(item.metadata).length > 0 && (
                        <details className="text-xs">
                          <summary className="cursor-pointer text-blue-600 hover:text-blue-700">
                            View details
                          </summary>
                          <div className="mt-2 p-2 bg-gray-50 rounded text-gray-700">
                            {Object.entries(item.metadata).map(([key, value]) => (
                              <div key={key} className="flex gap-2">
                                <span className="font-medium">{key}:</span>
                                <span>{String(value)}</span>
                              </div>
                            ))}
                          </div>
                        </details>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Pagination */}
        {items.length > 0 && (
          <div className="flex items-center justify-between mt-6 pt-4 border-t">
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1 || isLoading}
              className="btn-secondary"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              <span>Previous</span>
            </button>
            
            <span className="text-sm text-gray-600">
              Page {page}
            </span>
            
            <button
              onClick={() => setPage(p => p + 1)}
              disabled={!hasMore || isLoading}
              className="btn-secondary"
            >
              <span>Next</span>
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

