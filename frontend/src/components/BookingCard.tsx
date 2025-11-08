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

export default function BookingCard({ booking, onCancel, onEdit, onTransfer }: Props) {
  return (
    <div className="card flex items-center justify-between">
      <div>
        <div className="font-medium">Room {booking.room}</div>
        <div className="text-sm text-gray-600">{booking.date} • {booking.time}</div>
        <div className="text-xs text-gray-500">Status: {booking.status}</div>
      </div>
      <div className="flex items-center gap-2">
        {onEdit && (
          <button className="px-3 py-1 text-sm rounded bg-blue-600 text-white" onClick={() => onEdit(booking.id)}>Edit</button>
        )}
        {onTransfer && (
          <button className="px-3 py-1 text-sm rounded bg-amber-600 text-white" onClick={() => onTransfer(booking.id)}>Transfer</button>
        )}
        {onCancel && (
          <button className="px-3 py-1 text-sm rounded bg-red-600 text-white" onClick={() => onCancel(booking.id)}>Cancel</button>
        )}
      </div>
    </div>
  );
}

