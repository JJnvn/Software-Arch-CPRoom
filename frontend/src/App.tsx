import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import Navbar from './components/Navbar';
import Sidebar from './components/Sidebar';
import Login from './pages/Auth/Login';
import Register from './pages/Auth/Register';
import Profile from './pages/User/Profile';
import BookingHistory from './pages/User/BookingHistory';
import Preferences from './pages/User/Preferences';
import SearchRooms from './pages/Rooms/SearchRooms';
import RoomSchedule from './pages/Rooms/RoomSchedule';
import CreateBooking from './pages/Rooms/CreateBooking';
import EditBooking from './pages/Rooms/EditBooking';
import TransferBooking from './pages/Rooms/TransferBooking';
import AdminRoomBookings from './pages/Rooms/AdminRoomBookings';
import PendingApprovals from './pages/Staff/PendingApprovals';
import NotificationHistory from './pages/Notifications/NotificationHistory';
import { useAuth } from './hooks/useAuth';

function AppLayout() {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 bg-gray-50">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function ProtectedRoute({ children }: { children: JSX.Element }) {
  const { user } = useAuth();
  const location = useLocation();
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

      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/" element={<Navigate to="/rooms/search" replace />} />
        {/* User */}
        <Route path="/profile" element={<Profile />} />
        <Route path="/bookings" element={<BookingHistory />} />
        <Route path="/preferences" element={<Preferences />} />

        {/* Rooms */}
        <Route path="/rooms/search" element={<SearchRooms />} />
        <Route path="/rooms/:id/schedule" element={<RoomSchedule />} />
        <Route path="/bookings/create" element={<CreateBooking />} />
        <Route path="/bookings/:id/edit" element={<EditBooking />} />
        <Route path="/bookings/:id/transfer" element={<TransferBooking />} />
        <Route path="/admin/rooms/:id/bookings" element={<AdminRoomBookings />} />

        {/* Staff */}
        <Route path="/staff/pending-approvals" element={<PendingApprovals />} />

        {/* Notifications */}
        <Route path="/notifications" element={<NotificationHistory />} />
      </Route>
    </Routes>
  );
}

