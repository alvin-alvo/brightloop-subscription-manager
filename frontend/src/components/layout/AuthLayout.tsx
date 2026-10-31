import React from 'react';
import { Outlet, Link } from 'react-router-dom';
import { ShieldCheck } from 'lucide-react';

interface AuthLayoutProps {
  children?: React.ReactNode;
  title?: string;
  subtitle?: string;
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({
  children,
  title,
  subtitle,
}) => {
  return (
    <div className="min-h-screen bg-stripe-canvas-soft flex flex-col justify-between relative overflow-hidden font-sans selection:bg-indigo-100 selection:text-stripe-primary">
      {/* Stripe Signature Atmospheric Gradient Mesh Wash (Upper Third) */}
      <div
        className="absolute top-0 inset-x-0 h-96 sm:h-[420px] pointer-events-none overflow-hidden select-none"
        aria-hidden="true"
      >
        {/* Soft atmospheric gradient blobs */}
        <div className="absolute -top-32 -left-20 w-[600px] h-[500px] rounded-full bg-gradient-to-br from-[#f5e9d4] via-[#b9b9f9]/50 to-[#533afd]/20 blur-3xl opacity-60" />
        <div className="absolute -top-28 right-[-80px] w-[550px] h-[500px] rounded-full bg-gradient-to-bl from-[#f96bee]/30 via-[#ea2261]/20 to-[#533afd]/30 blur-3xl opacity-50" />
        <div className="absolute top-24 left-1/2 -translate-x-1/2 w-[720px] h-[280px] rounded-full bg-gradient-to-r from-[#b9b9f9]/30 via-[#533afd]/15 to-[#f5e9d4]/30 blur-2xl opacity-70" />
        <div className="absolute bottom-0 inset-x-0 h-24 bg-gradient-to-b from-transparent to-stripe-canvas-soft" />
      </div>

      {/* Main Centered Content Area */}
      <main className="relative z-10 flex-1 flex flex-col items-center justify-center px-4 py-12 sm:px-6 lg:px-8">
        {/* Brand Header */}
        <div className="mb-8 text-center flex flex-col items-center">
          <Link
            to="/"
            className="group inline-flex items-center gap-2.5 transition-transform duration-200 hover:scale-[1.02]"
            aria-label="Brightloop Home"
          >
            <div className="w-10 h-10 rounded-xl bg-stripe-primary flex items-center justify-center text-white font-bold text-lg shadow-stripe-sm group-hover:bg-stripe-primary-deep transition-colors">
              BL
            </div>
            <span className="font-semibold text-2xl tracking-tight text-stripe-ink">
              Brightloop
            </span>
          </Link>
          <p className="text-xs text-stripe-ink-mute mt-1.5 font-normal tracking-wide uppercase">
            B2B SaaS Subscription Infrastructure
          </p>
        </div>

        {/* Stripe Polished Surface Card */}
        <div className="w-full max-w-[440px] bg-white border border-stripe-hairline rounded-2xl shadow-stripe-md p-7 sm:p-9 transition-all">
          {(title || subtitle) && (
            <div className="mb-6 text-center sm:text-left">
              {title && (
                <h1 className="text-2xl sm:text-[26px] font-light text-stripe-ink tracking-tight">
                  {title}
                </h1>
              )}
              {subtitle && (
                <p className="text-sm text-stripe-ink-mute mt-1.5 font-normal leading-relaxed">
                  {subtitle}
                </p>
              )}
            </div>
          )}

          {/* Form Content / Route View */}
          {children || <Outlet />}
        </div>

        {/* Security & System Trust Signal */}
        <div className="mt-8 flex items-center gap-2 text-xs text-stripe-ink-mute">
          <ShieldCheck className="w-4 h-4 text-emerald-600" />
          <span>256-bit TLS Encrypted & SOC 2 Type II Certified</span>
        </div>
      </main>

      {/* Footer */}
      <footer className="relative z-10 py-6 text-center text-xs text-stripe-ink-mute border-t border-stripe-hairline/60 bg-white/40 backdrop-blur-sm">
        <div className="max-w-md mx-auto px-4 flex flex-wrap items-center justify-center gap-x-4 gap-y-1">
          <span>&copy; {new Date().getFullYear()} Brightloop Software Inc.</span>
          <span className="hidden sm:inline">&bull;</span>
          <a
            href="#terms"
            onClick={(e) => e.preventDefault()}
            className="hover:text-stripe-ink transition-colors"
          >
            Terms of Service
          </a>
          <span className="hidden sm:inline">&bull;</span>
          <a
            href="#privacy"
            onClick={(e) => e.preventDefault()}
            className="hover:text-stripe-ink transition-colors"
          >
            Privacy Policy
          </a>
          <span className="hidden sm:inline">&bull;</span>
          <a
            href="#security"
            onClick={(e) => e.preventDefault()}
            className="hover:text-stripe-ink transition-colors"
          >
            Security
          </a>
        </div>
      </footer>
    </div>
  );
};

export default AuthLayout;

