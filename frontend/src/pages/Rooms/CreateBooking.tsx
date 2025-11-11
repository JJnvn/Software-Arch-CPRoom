import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import * as rooms from "@/services/rooms";

export default function CreateBooking() {
    const navigate = useNavigate();
    const [roomName, setRoomName] = useState("");
    const [date, setDate] = useState("");
    const [time, setTime] = useState("");
    const [duration, setDuration] = useState<number | "">("");
    const [status, setStatus] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    async function onSubmit(e: FormEvent) {
        e.preventDefault();
        
        setError(null);
        setStatus(null);

        if (!roomName || !date || !time || !duration) {
            setError("Please fill in all fields");
            return;
        }

        const token = localStorage.getItem('AUTH_TOKEN');
        if (!token) {
            setError("Please log in to create a booking");
            return;
        }

        const userJson = localStorage.getItem("auth_user");
        if (!userJson) {
            setError("User not logged in");
            return;
        }

        const user = JSON.parse(userJson);
        const startDateTime = new Date(`${date}T${time}`);
        const endDateTime = new Date(
            startDateTime.getTime() + Number(duration) * 60000
        );

        // Validate booking is in the future
        if (startDateTime < new Date()) {
            setError("Booking time must be in the future");
            return;
        }

        setIsSubmitting(true);
        setStatus("Creating booking...");

        try {
            // Resolve room name to id
            const allRooms = await rooms.listRooms();
            const match = allRooms.find(r => (r.name || '').toLowerCase() === roomName.trim().toLowerCase());
            
            if (!match) {
                setError("Room not found. Please check the room name.");
                setStatus(null);
                setIsSubmitting(false);
                return;
            }

            const payload = {
                user_id: user.id,
                room_id: match.id,
                start_time: startDateTime.toISOString(),
                end_time: endDateTime.toISOString(),
            };

            await rooms.createBooking(payload);
            setStatus("Booking created successfully!");
            setError(null);
            
            // Clear form
            setRoomName("");
            setDate("");
            setTime("");
            setDuration("");
            
            // Redirect after success
            setTimeout(() => navigate('/booking-history'), 2000);
        } catch (err: any) {
            console.error(err);
            setError(err.response?.data?.error || "Failed to create booking");
            setStatus(null);
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <div className="page">
            <div className="max-w-2xl mx-auto">
                <button 
                    onClick={() => navigate('/rooms/search')} 
                    className="btn-secondary btn-sm mb-4"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                    Back to Search
                </button>

                <h1 className="page-title">Create New Booking</h1>
                
                <form onSubmit={onSubmit} className="card space-y-5">
                    {error && (
                        <div className="alert-error flex items-start gap-2">
                            <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            {error}
                        </div>
                    )}
                    
                    {status && (
                        <div className="alert-success flex items-start gap-2">
                            <svg className="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            {status}
                        </div>
                    )}

                    <div>
                        <label className="block text-sm font-semibold mb-2">
                            <span className="flex items-center gap-1">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                                </svg>
                                Room Name
                            </span>
                        </label>
                        <input
                            value={roomName}
                            onChange={(e) => {
                                setRoomName(e.target.value);
                                setError(null);
                            }}
                            className="input"
                            placeholder="e.g., Conference Room A"
                            disabled={isSubmitting}
                            required
                        />
                        <p className="text-xs text-gray-500 mt-1">Enter the exact room name</p>
                    </div>

                    <div>
                        <label className="block text-sm font-semibold mb-2">
                            <span className="flex items-center gap-1">
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                </svg>
                                Date
                            </span>
                        </label>
                        <input
                            value={date}
                            onChange={(e) => {
                                setDate(e.target.value);
                                setError(null);
                            }}
                            type="date"
                            className="input"
                            min={new Date().toISOString().slice(0, 10)}
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-semibold mb-2">
                                <span className="flex items-center gap-1">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    Start Time
                                </span>
                            </label>
                            <input
                                value={time}
                                onChange={(e) => {
                                    setTime(e.target.value);
                                    setError(null);
                                }}
                                type="time"
                                className="input"
                                disabled={isSubmitting}
                                required
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-semibold mb-2">
                                <span className="flex items-center gap-1">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    Duration (minutes)
                                </span>
                            </label>
                            <input
                                value={duration}
                                onChange={(e) => {
                                    setDuration(
                                        e.target.value === ""
                                            ? ""
                                            : Number(e.target.value)
                                    );
                                    setError(null);
                                }}
                                type="number"
                                min={15}
                                step={15}
                                className="input"
                                placeholder="60"
                                disabled={isSubmitting}
                                required
                            />
                            <p className="text-xs text-gray-500 mt-1">Minimum 15 minutes</p>
                        </div>
                    </div>

                    <div className="flex gap-3 pt-2">
                        <button 
                            type="submit" 
                            className="btn-primary flex-1"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? (
                                <>
                                    <div className="spinner"></div>
                                    Creating Booking...
                                </>
                            ) : (
                                <>
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                    </svg>
                                    Create Booking
                                </>
                            )}
                        </button>
                        <button 
                            type="button"
                            onClick={() => navigate('/rooms/search')}
                            className="btn-secondary"
                            disabled={isSubmitting}
                        >
                            Cancel
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
