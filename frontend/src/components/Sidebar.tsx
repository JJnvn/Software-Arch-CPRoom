import { NavLink } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

function Item({ to, label }: { to: string; label: string }) {
    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                `block px-4 py-2 rounded hover:bg-gray-100 ${
                    isActive ? "bg-gray-200 font-medium" : ""
                }`
            }
        >
            {label}
        </NavLink>
    );
}

export default function Sidebar() {
    const { user } = useAuth();
    const role = (user?.role || '').toUpperCase();
    const isPrivileged = role === 'ADMIN' || role === 'STAFF';
    return (
        <aside className="w-60 bg-white border-r p-4 space-y-2">
            <div className="text-xs uppercase text-gray-500 px-2">User</div>
            <Item to="/profile" label="Profile" />
            <Item to="/bookings" label="Booking History" />
            <Item to="/preferences" label="Preferences" />
            <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                Rooms
            </div>
            <Item to="/rooms/search" label="Search Rooms" />
            <Item to="/bookings/create" label="Create Booking" />
            {isPrivileged && (
                <>
                    <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                        Staff
                    </div>
                    <Item to="/approvals/pending" label="Pending Approvals" />
                    <Item to="/admin/rooms/create" label="Create Room" />
                    <Item to="/admin/rooms/101/bookings" label="Admin Room Bookings" />
                </>
            )}
            <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                Notifications
            </div>
            <Item to="/notifications" label="Notification History" />
        </aside>
    );
}
