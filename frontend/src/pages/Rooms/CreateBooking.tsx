import { FormEvent, useState } from "react";
import * as rooms from "@/services/rooms";

export default function CreateBooking() {
    const [roomName, setRoomName] = useState("");
    const [date, setDate] = useState("");
    const [time, setTime] = useState("");
    const [duration, setDuration] = useState<number | "">("");
    const [status, setStatus] = useState<string | null>(null);

    async function onSubmit(e: FormEvent) {
        e.preventDefault();
        if (!roomName || !date || !time || !duration) return;

        const token = localStorage.getItem('AUTH_TOKEN');
        if (!token) {
            setStatus("Please log in to create a booking");
            setTimeout(() => setStatus(null), 2000);
            return;
        }

        // parse auth_user from localStorage
        const userJson = localStorage.getItem("auth_user");
        if (!userJson) {
            setStatus("User not logged in");
            return;
        }
        const user = JSON.parse(userJson);

        // combine date + time into a JS Date object
        const startDateTime = new Date(`${date}T${time}`);
        const endDateTime = new Date(
            startDateTime.getTime() + Number(duration) * 60000
        );

        // resolve room name to id
        let roomId: string | null = null;
        try {
            const allRooms = await rooms.listRooms();
            const match = allRooms.find(r => (r.name || '').toLowerCase() === roomName.trim().toLowerCase());
            if (match) roomId = match.id;
        } catch (e) {
            // ignore here; handled below
        }

        if (!roomId) {
            setStatus("Room not found");
            setTimeout(() => setStatus(null), 2000);
            return;
        }

        // create payload
        const payload = {
            user_id: user.id,
            room_id: roomId,
            start_time: startDateTime.toISOString(),
            end_time: endDateTime.toISOString(),
        };

        try {
            await rooms.createBooking(payload);
            setStatus("Booking created");
            setTimeout(() => setStatus(null), 2000);
        } catch (err) {
            console.error(err);
            setStatus("Failed to create booking");
            setTimeout(() => setStatus(null), 2000);
        }
    }

    return (
        <div className="page">
            <h1 className="page-title">Create Booking</h1>
            <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
                {status && (
                    <div className="text-green-700 text-sm">{status}</div>
                )}
                <div>
                    <label className="block text-sm mb-1">Room Name</label>
                    <input
                        value={roomName}
                        onChange={(e) => setRoomName(e.target.value)}
                        className="w-full border rounded px-3 py-2"
                        placeholder="e.g., Phoenix"
                        required
                    />
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                    <div>
                        <label className="block text-sm mb-1">Date</label>
                        <input
                            value={date}
                            onChange={(e) => setDate(e.target.value)}
                            type="date"
                            className="w-full border rounded px-3 py-2"
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm mb-1">Start Time</label>
                        <input
                            value={time}
                            onChange={(e) => setTime(e.target.value)}
                            type="time"
                            className="w-full border rounded px-3 py-2"
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm mb-1">
                            Duration (min)
                        </label>
                        <input
                            value={duration}
                            onChange={(e) =>
                                setDuration(
                                    e.target.value === ""
                                        ? ""
                                        : Number(e.target.value)
                                )
                            }
                            type="number"
                            min={15}
                            step={15}
                            className="w-full border rounded px-3 py-2"
                            required
                        />
                    </div>
                </div>
                <button className="px-4 py-2 bg-green-600 text-white rounded">
                    Create Booking
                </button>
            </form>
        </div>
    );
}
