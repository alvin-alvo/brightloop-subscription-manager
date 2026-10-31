import React from 'react';
import { NavLink } from 'react-router-dom';
import { APP_ROUTES } from '../../routes/routes.config';
import { useAuth } from '../../hooks/useAuth';
import { X } from 'lucide-react';

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const { role, user } = useAuth();

  // Filter routes strictly matching the current user role
  const visibleRoutes = APP_ROUTES.filter((route) =>
    route.allowedRoles.includes(role)
  );

  // Group routes by section
  const sections = Array.from(new Set(visibleRoutes.map((r) => r.section)));

  return (
    <>
      {/* Mobile Backdrop */}
      {isOpen && (
        <div
          onClick={onClose}
          className="fixed inset-0 z-40 bg-slate-900/30 backdrop-blur-sm lg:hidden transition-opacity"
          aria-hidden="true"
        />
      )}

      {/* Sidebar Container */}
      <aside
        className={`fixed top-0 bottom-0 left-0 z-50 w-64 bg-white border-r border-stripe-hairline flex flex-col transition-transform duration-200 ease-in-out lg:translate-x-0 ${isOpen ? 'translate-x-0' : '-translate-x-full'
          }`}
      >
        {/* Brand Header */}
        <div className="h-16 px-6 flex items-center justify-between border-b border-stripe-hairline shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-stripe-primary flex items-center justify-center text-white font-bold text-sm shadow-stripe-sm">
              BL
            </div>
            <div>
              <div className="flex items-center gap-1.5">
                <span className="font-semibold text-stripe-ink text-sm tracking-tight">
                  Brightloop
                </span>
                <span className="px-1.5 py-0.5 rounded text-[10px] font-medium bg-indigo-50 text-stripe-primary border border-indigo-100 uppercase">
                  {role}
                </span>
              </div>
              <p className="text-[11px] text-stripe-ink-mute truncate max-w-[130px]">
                {user?.companyName || 'Subscription HQ'}
              </p>
            </div>
          </div>

          <button
            onClick={onClose}
            className="p-1.5 rounded-md text-slate-400 hover:text-stripe-ink lg:hidden"
            aria-label="Close sidebar"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation Item Groups */}
        <div className="flex-1 overflow-y-auto px-4 py-5 space-y-6">
          {sections.map((sectionName) => {
            const routesInSection = visibleRoutes.filter(
              (r) => r.section === sectionName
            );

            return (
              <div key={sectionName} className="space-y-1">
                <div className="px-3 pb-1.5 text-[10px] font-semibold uppercase tracking-wider text-stripe-ink-mute">
                  {sectionName}
                </div>

                {routesInSection.map((route) => {
                  const Icon = route.icon;
                  return (
                    <NavLink
                      key={route.id}
                      to={route.path}
                      onClick={onClose}
                      className={({ isActive }) =>
                        `flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150 ${isActive
                          ? 'bg-indigo-50/80 text-stripe-primary font-medium shadow-sm ring-1 ring-stripe-primary/10'
                          : 'text-stripe-ink-secondary hover:bg-slate-50 hover:text-stripe-ink'
                        }`
                      }
                    >
                      {({ isActive }) => (
                        <>
                          <Icon
                            className={`w-4 h-4 shrink-0 transition-colors ${isActive ? 'text-stripe-primary' : 'text-slate-400'
                              }`}
                          />
                          <span className="truncate">{route.label}</span>
                          {isActive && (
                            <span className="ml-auto w-1.5 h-1.5 rounded-full bg-stripe-primary" />
                          )}
                        </>
                      )}
                    </NavLink>
                  );
                })}
              </div>
            );
          })}
        </div>

        {/* Sidebar Footer with Live Connection Status */}
        <div className="p-4 border-t border-stripe-hairline bg-slate-50/50 shrink-0">
          <div className="flex items-center justify-between text-xs text-stripe-ink-mute">
            <span className="flex items-center gap-1.5">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
              Live Gateway
            </span>
            <span className="font-mono text-[10px]">v1.0.0</span>
          </div>
        </div>
      </aside>
    </>
  );
};

export default Sidebar;