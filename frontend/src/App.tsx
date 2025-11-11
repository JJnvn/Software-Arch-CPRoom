import { Navigate, Outlet, Route, Routes, useLocation } from "react-router-dom";
import Navbar from "./components/Navbar";
import Sidebar from "./components/Sidebar";
import Login from "./pages/Auth/Login";
import Register from "./pages/Auth/Register";
import OAuthCallback from "./pages/Auth/OAuthCallback";
import Profile from "./pages/User/Profile";
import ProfileEdit from "./pages/User/ProfileEdit";
import BookingHistory from "./pages/User/BookingHistory";
import PreferencesEdit from "./pages/User/PreferencesEdit";
import RescheduleBooking from "./pages/Bookings/RescheduleBooking";
import TransferBooking from "./pages/Bookings/TransferBooking";
import RoomSchedule from "./pages/Rooms/RoomSchedule";
import SearchAndBook from "./pages/Rooms/SearchAndBook";
import EditBooking from "./pages/Rooms/EditBooking";

import AdminCreateRoom from "./pages/Admin/CreateRoom";
import EditRoom from "./pages/Admin/EditRoom";
import PendingApprovals from "./pages/Staff/PendingApprovals";
import ApprovedBookings from "./pages/Staff/ApprovedBookings";
import AuditTrail from "./pages/Staff/AuditTrail";
import AdminRoomBookings from "./pages/Admin/RoomBookings";
import NotificationHistory from "./pages/Notifications/NotificationHistory";
import { useAuth } from "./hooks/useAuth";

function AppLayout() {
    return (
        <div className="flex h-screen overflow-hidden">
            <Sidebar />
            <div className="flex flex-col flex-1 overflow-hidden">
                <Navbar />
                <main className="flex-1 overflow-y-auto bg-gray-50">
                    <Outlet />
                </main>
            </div>
        </div>
    );
}

function ProtectedRoute({ children }: { children: JSX.Element }) {
    const { user, loading } = useAuth();
    const location = useLocation();
    if (loading) {
        return null;
    }
    if (!user) {
        return <Navigate to="/login" replace state={{ from: location }} />;
    }
    return children;
}

export default function App() {
    return (
        <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/auth/callback" element={<OAuthCallback />} />

            <Route
                element={
                    <ProtectedRoute>
                        <AppLayout />
                    </ProtectedRoute>
                }
            >
                <Route path="/" element={<Navigate to="/rooms/search" replace />} />
                
                {/* User Profile & Settings */}
                <Route path="/profile" element={<Profile />} />
                <Route path="/profile/edit" element={<ProfileEdit />} />
                <Route path="/preferences" element={<PreferencesEdit />} />

                {/* Rooms & Bookings */}
                <Route path="/rooms/search" element={<SearchAndBook />} />
                <Route path="/rooms/:id/schedule" element={<RoomSchedule />} />
                <Route path="/booking-history" element={<BookingHistory />} />
                <Route path="/bookings/:id/edit" element={<EditBooking />} />
                <Route path="/bookings/:id/transfer" element={<TransferBooking />} />

                {/* Notifications */}
                <Route path="/notifications" element={<NotificationHistory />} />

                {/* Staff & Admin */}
                <Route path="/approvals/pending" element={<PendingApprovals />} />
                <Route path="/approvals/approved" element={<ApprovedBookings />} />
                <Route path="/admin/rooms/create" element={<AdminCreateRoom />} />
                <Route path="/admin/rooms/:id/edit" element={<EditRoom />} />
                <Route path="/admin/room-bookings" element={<AdminRoomBookings />} />
                <Route path="/staff/audit-trail" element={<AuditTrail />} />
            </Route>
        </Routes>
    );
}
