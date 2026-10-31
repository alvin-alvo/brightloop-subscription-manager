import React, { useState, useRef, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Menu,
  ChevronRight,
  Bell,
  BellRing,
  CheckCheck,
  Zap,
  LogOut,
  Shield,
  UserCheck,
  User,
  Sparkles,
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useToast } from '../../context/ToastContext';
import { APP_ROUTES, getDefaultRouteForRole } from '../../routes/routes.config';
import type { UserRole } from '../../types';
import {
  notificationsService,
  type AppNotification,
} from '../../services/notificationsService';

interface HeaderProps {
  onOpenSidebar: () => void;
}

export const Header: React.FC<HeaderProps> = ({ onOpenSidebar }) => {
  const { user, role, switchRole, logout } = useAuth();
  const { simulateRealtimeEvent } = useToast();
  const location = useLocation();
  const navigate = useNavigate();

  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const [isNotifOpen, setIsNotifOpen] = useState(false);
  const notifRef = useRef<HTMLDivElement>(null);
  const [notifications, setNotifications] = useState<AppNotification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [notifLoading, setNotifLoading] = useState(true);

  const refreshNotifications = async () => {
    try {
      const result = await notificationsService.getNotifications(10);
      setNotifications(result.items);
      setUnreadCount(result.unreadCount);
    } catch {
      // The service falls back to local data; keep current state otherwise.
    } finally {
      setNotifLoading(false);
    }
  };

  useEffect(() => {
    refreshNotifications();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Close dropdowns on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      const node = event.target as Node;
      if (dropdownRef.current && !dropdownRef.current.contains(node)) {
        setIsDropdownOpen(false);
      }
      if (notifRef.current && !notifRef.current.contains(node)) {
        setIsNotifOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleMarkAllRead = async () => {
    await notificationsService.markAllRead();
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    setUnreadCount(0);
  };

  const handleMarkRead = async (id: string) => {
    await notificationsService.markRead(id);
    setNotifications((prev) =>
      prev.map((n) => (n.id === id ? { ...n, read: true } : n))
    );
    setUnreadCount((count) => Math.max(0, count - 1));
  };

  const handleBellToggle = async () => {
    setIsNotifOpen((prev) => {
      const open = !prev;
      if (open && notifications.length === 0) refreshNotifications();
      return open;
    });
  };

  // Compute active breadcrumb item from route configuration
  const currentRoute = APP_ROUTES.find((r) => r.path === location.pathname);

  const handleRoleChange = (newRole: UserRole) => {
    switchRole(newRole);
    setIsDropdownOpen(false);
    // Navigate to the role's default dashboard immediately
    const nextRoute = getDefaultRouteForRole(newRole);
    navigate(nextRoute);
  };

  const roleColors: Record<UserRole, string> = {
    admin: 'bg-indigo-50 text-stripe-primary border-indigo-200',
    staff: 'bg-slate-100 text-slate-800 border-slate-300',
    customer: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  };

  return (
    <header className="sticky top-0 z-30 h-16 bg-white/95 backdrop-blur-md border-b border-stripe-hairline px-4 sm:px-6 lg:px-8 flex items-center justify-between">
      {/* Left Section: Mobile Menu + Breadcrumbs */}
      <div className="flex items-center gap-3">
        <button
          onClick={onOpenSidebar}
          className="p-1.5 -ml-1 text-slate-500 hover:text-stripe-ink rounded-lg hover:bg-slate-50 lg:hidden"
          aria-label="Open sidebar"
        >
          <Menu className="w-5 h-5" />
        </button>

        {/* Dynamic Breadcrumbs */}
        <nav aria-label="Breadcrumb" className="flex items-center gap-1.5 text-xs">
          <span className="text-stripe-ink-mute hidden sm:inline">Brightloop</span>
          <ChevronRight className="w-3.5 h-3.5 text-slate-300 hidden sm:inline" />
          <span className="font-normal text-stripe-ink-mute capitalize">
            {role === 'customer' ? 'Customer Portal' : role}
          </span>
          <ChevronRight className="w-3.5 h-3.5 text-slate-300" />
          <span className="font-medium text-stripe-ink">
            {currentRoute?.label || 'Dashboard'}
          </span>
        </nav>
      </div>

      {/* Right Section: Real-time trigger & User Menu */}
      <div className="flex items-center gap-2.5 sm:gap-3">
        {/* Notifications dropdown */}
        <div className="relative" ref={notifRef}>
          <button
            onClick={handleBellToggle}
            className="relative p-2 text-slate-500 hover:text-stripe-ink rounded-full hover:bg-slate-100 transition-colors"
            aria-label="Notifications"
            title="Notifications"
          >
            <Bell className="w-5 h-5" />
            {unreadCount > 0 && (
              <span className="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 rounded-full bg-rose-500 text-white text-[10px] font-semibold flex items-center justify-center">
                {unreadCount > 9 ? '9+' : unreadCount}
              </span>
            )}
          </button>

          {isNotifOpen && (
            <div className="absolute right-0 mt-2 w-80 bg-white border border-stripe-hairline rounded-xl shadow-stripe-md overflow-hidden z-50 text-xs">
              <div className="px-4 py-3 border-b border-stripe-hairline flex items-center justify-between">
                <p className="font-medium text-stripe-ink text-sm flex items-center gap-2">
                  <BellRing className="w-4 h-4 text-stripe-primary" />
                  Notifications
                  {unreadCount > 0 && (
                    <span className="text-[10px] font-semibold text-stripe-primary bg-indigo-50 px-1.5 py-0.5 rounded-full">
                      {unreadCount} new
                    </span>
                  )}
                </p>
                {unreadCount > 0 && (
                  <button
                    onClick={handleMarkAllRead}
                    className="inline-flex items-center gap-1 text-stripe-primary hover:text-stripe-primary-deep font-medium"
                  >
                    <CheckCheck className="w-3.5 h-3.5" />
                    Mark all read
                  </button>
                )}
              </div>

              <div className="max-h-80 overflow-y-auto divide-y divide-stripe-hairline">
                {notifLoading && (
                  <div className="px-4 py-6 text-center text-stripe-ink-mute">
                    Loading notifications...
                  </div>
                )}
                {!notifLoading && notifications.length === 0 && (
                  <div className="px-4 py-8 text-center text-stripe-ink-mute">
                    You're all caught up.
                  </div>
                )}
                {notifications.map((notification) => (
                  <button
                    key={notification.id}
                    onClick={() => handleMarkRead(notification.id)}
                    className={`w-full text-left px-4 py-3 transition-colors ${
                      notification.read
                        ? 'hover:bg-slate-50'
                        : 'bg-indigo-50/50 hover:bg-indigo-50'
                    }`}
                  >
                    <div className="flex items-start gap-2.5">
                      {!notification.read && (
                        <span className="w-2 h-2 mt-1.5 shrink-0 rounded-full bg-stripe-primary" />
                      )}
                      <span className="flex-1 min-w-0">
                        <span className="block font-medium text-stripe-ink pr-6">
                          {notification.title}
                        </span>
                        <span className="block text-stripe-ink-mute mt-0.5 leading-snug">
                          {notification.body}
                        </span>
                        <span className="block text-[11px] text-stripe-ink-secondary mt-1">
                          {notification.displayTime}
                        </span>
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Real-time simulation event button */}
        <button
          onClick={simulateRealtimeEvent}
          className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-stripe-primary bg-indigo-50/80 hover:bg-indigo-100/80 border border-indigo-100 rounded-full transition-colors"
          title="Simulate incoming WebSocket webhook"
        >
          <Zap className="w-3.5 h-3.5 animate-bounce" />
          <span className="hidden sm:inline">Simulate Event</span>
        </button>

        {/* User Profile & Role Switcher Dropdown */}
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setIsDropdownOpen((prev) => !prev)}
            className="flex items-center gap-2 p-1 pl-2 bg-slate-50 hover:bg-slate-100/80 border border-stripe-hairline rounded-full transition-colors"
            aria-expanded={isDropdownOpen}
            aria-haspopup="true"
          >
            <span
              className={`px-2 py-0.5 rounded-full text-[10px] font-semibold uppercase border ${roleColors[role]}`}
            >
              {role}
            </span>
            <div className="w-7 h-7 rounded-full bg-stripe-primary text-white flex items-center justify-center font-medium text-xs shadow-stripe-sm">
              {user?.name
                ?.split(' ')
                .map((n) => n[0])
                .join('')
                .slice(0, 2) || 'U'}
            </div>
          </button>

          {/* Dropdown Menu (Stripe card style) */}
          {isDropdownOpen && (
            <div className="absolute right-0 mt-2 w-64 bg-white border border-stripe-hairline rounded-xl shadow-stripe-md py-2 z-50 text-xs">
              {/* Profile Details */}
              <div className="px-4 py-2.5 border-b border-stripe-hairline">
                <p className="font-medium text-stripe-ink text-sm truncate">
                  {user?.name}
                </p>
                <p className="text-stripe-ink-mute truncate">{user?.email}</p>
                <p className="text-[11px] text-stripe-ink-secondary mt-1 font-medium truncate">
                  {user?.companyName}
                </p>
              </div>

              {/* Role Switcher (For Instant Testing) */}
              <div className="px-3 py-2 border-b border-stripe-hairline">
                <p className="px-2 pb-1.5 text-[10px] font-semibold text-stripe-ink-mute uppercase tracking-wider flex items-center gap-1">
                  <Sparkles className="w-3 h-3 text-stripe-primary" />
                  Quick Role Switcher
                </p>
                <div className="space-y-1">
                  <button
                    onClick={() => handleRoleChange('admin')}
                    className={`w-full flex items-center justify-between px-2.5 py-1.5 rounded-md text-left transition-colors ${
                      role === 'admin'
                        ? 'bg-indigo-50 font-medium text-stripe-primary'
                        : 'text-stripe-ink-secondary hover:bg-slate-50'
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <Shield className="w-3.5 h-3.5" />
                      Admin (Full System)
                    </span>
                    {role === 'admin' && (
                      <span className="w-1.5 h-1.5 rounded-full bg-stripe-primary" />
                    )}
                  </button>

                  <button
                    onClick={() => handleRoleChange('staff')}
                    className={`w-full flex items-center justify-between px-2.5 py-1.5 rounded-md text-left transition-colors ${
                      role === 'staff'
                        ? 'bg-indigo-50 font-medium text-stripe-primary'
                        : 'text-stripe-ink-secondary hover:bg-slate-50'
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <UserCheck className="w-3.5 h-3.5" />
                      Staff (Operations & Sales)
                    </span>
                    {role === 'staff' && (
                      <span className="w-1.5 h-1.5 rounded-full bg-stripe-primary" />
                    )}
                  </button>

                  <button
                    onClick={() => handleRoleChange('customer')}
                    className={`w-full flex items-center justify-between px-2.5 py-1.5 rounded-md text-left transition-colors ${
                      role === 'customer'
                        ? 'bg-indigo-50 font-medium text-stripe-primary'
                        : 'text-stripe-ink-secondary hover:bg-slate-50'
                    }`}
                  >
                    <span className="flex items-center gap-2">
                      <User className="w-3.5 h-3.5" />
                      Customer (Self Portal)
                    </span>
                    {role === 'customer' && (
                      <span className="w-1.5 h-1.5 rounded-full bg-stripe-primary" />
                    )}
                  </button>
                </div>
              </div>

              {/* Logout Action */}
              <div className="pt-1 px-1">
                <button
                  onClick={() => {
                    logout();
                    setIsDropdownOpen(false);
                    navigate('/login');
                  }}
                  className="w-full flex items-center gap-2 px-3 py-2 text-rose-600 hover:bg-rose-50 rounded-md transition-colors"
                >
                  <LogOut className="w-3.5 h-3.5" />
                  Sign Out
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;

