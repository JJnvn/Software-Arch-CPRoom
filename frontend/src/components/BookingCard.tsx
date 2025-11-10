type Props = {
  booking: {
    id: string;
    room: string;
    date: string;
    time: string;
    status: string;
  };
  onCancel?: (id: string) => void;
  onEdit?: (id: string) => void;
  onTransfer?: (id: string) => void;
};

function getStatusBadgeClass(status: string) {
  const normalized = status.toLowerCase();
  switch (normalized) {
    case 'pending':
      return 'badge-pending';
    case 'confirmed':
      return 'badge-confirmed';
    case 'cancelled':
      return 'badge-cancelled';
    case 'expired':
      return 'badge-expired';
    case 'denied':
      return 'badge-denied';
    default:
      return 'badge bg-gray-100 text-gray-800';
  }
}

export default function BookingCard({ booking, onCancel, onEdit, onTransfer }: Props) {
  const canModify = booking.status.toLowerCase() === 'pending' || booking.status.toLowerCase() === 'confirmed';
  
  return (
    <div className="card hover:shadow-lg transition-shadow duration-200">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <div className="font-semibold text-lg">Room {booking.room}</div>
            <span className={getStatusBadgeClass(booking.status)}>
              {booking.status}
            </span>
          </div>
          <div className="text-sm text-gray-600 space-y-0.5">
            <div className="flex items-center gap-1">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span>{booking.date}</span>
            </div>
            <div className="flex items-center gap-1">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{booking.time}</span>
            </div>
          </div>
        </div>
        <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
          {onEdit && canModify && (
            <button 
              className="btn-sm btn-primary" 
              onClick={() => onEdit(booking.id)}
              title="Reschedule booking"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Reschedule
            </button>
          )}
          {onTransfer && canModify && (
            <button 
              className="btn-sm btn-warning" 
              onClick={() => onTransfer(booking.id)}
              title="Transfer booking to another user"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
              Transfer
            </button>
          )}
          {onCancel && canModify && (
            <button 
              className="btn-sm btn-danger" 
              onClick={() => onCancel(booking.id)}
              title="Cancel this booking"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
              Cancel
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

