import { useState } from "react";
import { NavLink } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

interface NavItemProps {
  to: string;
  label: string;
  icon: React.ReactNode;
  collapsed: boolean;
}

function NavItem({ to, label, icon, collapsed }: NavItemProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `flex items-center ${collapsed ? 'justify-center' : 'gap-3'} px-3 py-2.5 rounded-lg transition-all duration-200 group relative ${
          isActive
            ? "bg-blue-50 text-blue-700 font-medium"
            : "text-gray-700 hover:bg-gray-100"
        }`
      }
      title={collapsed ? label : undefined}
    >
      <div className="w-5 h-5 flex-shrink-0">{icon}</div>
      <span className={`text-sm whitespace-nowrap transition-all duration-300 ${
        collapsed ? 'opacity-0 w-0 overflow-hidden' : 'opacity-100'
      }`}>
        {label}
      </span>
      {collapsed && (
        <div className="absolute left-full ml-2 px-3 py-1.5 bg-gray-900 text-white text-xs rounded-md opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap z-50 transition-opacity duration-200 shadow-lg">
          {label}
        </div>
      )}
    </NavLink>
  );
}

interface SectionHeaderProps {
  children: React.ReactNode;
  collapsed: boolean;
}

function SectionHeader({ children, collapsed }: SectionHeaderProps) {
  if (collapsed) {
    return <div className="border-t border-gray-200 my-2" />;
  }
  return (
    <div className="flex items-center gap-2 px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">
      {children}
    </div>
  );
}

export default function Sidebar() {
  const { user } = useAuth();
  const role = (user?.role || '').toUpperCase();
  const isPrivileged = role === 'ADMIN' || role === 'STAFF';
  const [collapsed, setCollapsed] = useState(false);

  return (
    <aside className={`${collapsed ? 'w-20' : 'w-64'} bg-white border-r border-gray-200 flex flex-col h-screen sticky top-0 transition-[width] duration-300 ease-in-out flex-shrink-0 overflow-hidden`}>
      {/* Toggle Button */}
      <div className="p-4 border-b border-gray-200 flex-shrink-0">
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="w-full flex items-center justify-center p-2 rounded-lg hover:bg-gray-100 active:bg-gray-200 transition-all duration-200 group"
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          <div className="relative w-full flex items-center justify-center">
            {/* Collapsed Icon */}
            <svg 
              className={`w-5 h-5 text-gray-600 group-hover:text-gray-900 transition-all duration-300 absolute ${
                collapsed ? 'opacity-100 rotate-0' : 'opacity-0 -rotate-90'
              }`}
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" />
            </svg>
            {/* Expanded Content */}
            <div className={`flex items-center justify-between w-full transition-all duration-300 ${
              collapsed ? 'opacity-0' : 'opacity-100'
            }`}>
              <span className={`text-sm font-semibold text-gray-700 transition-all duration-300 ${
                collapsed ? 'w-0 overflow-hidden' : ''
              }`}>Menu</span>
              <svg 
                className={`w-5 h-5 text-gray-600 group-hover:text-gray-900 transition-all duration-300 ${
                  collapsed ? 'opacity-0 rotate-90' : 'opacity-100 rotate-0'
                }`}
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
              </svg>
            </div>
          </div>
        </button>
      </div>

      <div className="flex-1 overflow-y-auto overflow-x-hidden p-4 space-y-6">
        {/* User Role Badge */}
        {user && (
          <div className={`transition-all duration-300 ${
            collapsed ? 'px-0' : 'px-3 py-2'
          } ${
            collapsed ? '' : 'bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg'
          }`}>
            <div className={`flex items-center ${collapsed ? 'justify-center' : 'gap-2'}`}>
              <div className={`${collapsed ? 'w-10 h-10' : 'w-8 h-8'} bg-blue-600 rounded-full flex items-center justify-center flex-shrink-0 transition-all duration-300`} 
                   title={collapsed ? `${user.name} (${role})` : undefined}>
                <span className={`${collapsed ? 'text-sm' : 'text-xs'} font-bold text-white transition-all duration-300`}>
                  {user.name?.charAt(0).toUpperCase() || 'U'}
                </span>
              </div>
              <div className={`flex-1 min-w-0 transition-all duration-300 ${
                collapsed ? 'opacity-0 w-0 overflow-hidden' : 'opacity-100'
              }`}>
                <p className="text-xs font-medium text-gray-900 truncate">{user.name}</p>
                <p className="text-xs text-blue-700 font-semibold">{role}</p>
              </div>
            </div>
          </div>
        )}

        {/* Rooms Section */}
        <div>
          <SectionHeader collapsed={collapsed}>
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
            Rooms & Booking
          </SectionHeader>
          <div className="space-y-1">
            <NavItem
              to="/rooms/search"
              label="Search & Book"
              collapsed={collapsed}
              icon={
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              }
            />
            <NavItem
              to="/booking-history"
              label="My Bookings"
              collapsed={collapsed}
              icon={
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              }
            />
          </div>
        </div>

        {/* User Section */}
        <div>
          <SectionHeader collapsed={collapsed}>
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Account
          </SectionHeader>
          <div className="space-y-1">
            <NavItem
              to="/profile"
              label="Profile"
              collapsed={collapsed}
              icon={
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              }
            />
            <NavItem
              to="/preferences"
              label="Preferences"
              collapsed={collapsed}
              icon={
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              }
            />
          </div>
        </div>

        {/* Notifications Section */}
        <div>
          <SectionHeader collapsed={collapsed}>
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
            Notifications
          </SectionHeader>
          <div className="space-y-1">
            <NavItem
              to="/notifications"
              label="History"
              collapsed={collapsed}
              icon={
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              }
            />
          </div>
        </div>

        {/* Staff/Admin Section */}
        {isPrivileged && (
          <div>
            <SectionHeader collapsed={collapsed}>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              {role} Access
            </SectionHeader>
            <div className="space-y-1">
              <NavItem
                to="/approvals/pending"
                label="Pending Approvals"
                collapsed={collapsed}
                icon={
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                }
              />
              <NavItem
                to="/approvals/approved"
                label="Approved Bookings"
                collapsed={collapsed}
                icon={
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                }
              />
              <NavItem
                to="/admin/rooms/create"
                label="Create Room"
                collapsed={collapsed}
                icon={
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                }
              />
            </div>
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="p-4 border-t border-gray-200 flex-shrink-0">
        <div className={`text-xs transition-all duration-300 ${
          collapsed ? 'text-center' : 'text-center'
        }`}>
          <p className={`text-gray-500 font-medium transition-all duration-300 ${
            collapsed ? 'opacity-0 h-0 overflow-hidden' : 'opacity-100'
          }`}>
            CPRoom Booking System
          </p>
          <p className={`text-gray-400 font-bold transition-all duration-300 ${
            collapsed ? 'text-base' : 'mt-1'
          }`}>
            {collapsed ? 'v1' : 'v1.0.0'}
          </p>
        </div>
      </div>
    </aside>
  );
}
