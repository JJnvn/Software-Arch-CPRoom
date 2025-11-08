import { Link } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export default function Navbar() {
  const { user, logout } = useAuth();
  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link to="/" className="font-semibold text-lg">CPRoom</Link>
        <nav className="flex items-center gap-4">
          {user && (
            <>
              <Link className="text-sm text-gray-700 hover:text-gray-900" to="/rooms/search">Rooms</Link>
              <Link className="text-sm text-gray-700 hover:text-gray-900" to="/bookings">My Bookings</Link>
              <Link className="text-sm text-gray-700 hover:text-gray-900" to="/notifications">Notifications</Link>
            </>
          )}
          <div className="flex items-center gap-3">
            {user ? (
              <>
                <span className="text-sm text-gray-600">{user.name}</span>
                <button onClick={logout} className="text-sm text-red-600 hover:underline">Logout</button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm hover:underline">Login</Link>
                <Link to="/register" className="text-sm hover:underline">Register</Link>
              </>
            )}
          </div>
        </nav>
      </div>
    </header>
  );
}

