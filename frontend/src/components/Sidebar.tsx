import { NavLink } from "react-router-dom";

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
            <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                Staff
            </div>
            <Item to="/approvals/pending" label="Pending Approvals" />
            <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                Admin
            </div>
            <Item to="/admin/rooms/101/bookings" label="Admin Room Bookings" />
            <div className="text-xs uppercase text-gray-500 px-2 mt-4">
                Notifications
            </div>
            <Item to="/notifications" label="Notification History" />
        </aside>
    );
}
