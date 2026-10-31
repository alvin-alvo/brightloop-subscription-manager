import React, { useState, useRef, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Menu,
  ChevronRight,
  Bell,
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

  // Close dropdown on outside click
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsDropdownOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

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
        {/* Real-time simulation event button */}
        <button
          onClick={simulateRealtimeEvent}
          className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-stripe-primary bg-indigo-50/80 hover:bg-indigo-100/80 border border-indigo-100 rounded-full transition-colors"
          title="Simulate incoming WebSocket webhook"
        >
          <Bell className="w-3.5 h-3.5 animate-bounce" />
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

