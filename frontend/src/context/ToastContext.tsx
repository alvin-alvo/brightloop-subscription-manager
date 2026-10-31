import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { CheckCircle2, AlertTriangle, Info, XCircle, X } from 'lucide-react';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastItem {
  id: string;
  title: string;
  message?: string;
  type: ToastType;
  duration?: number;
}

interface ToastContextType {
  toasts: ToastItem[];
  addToast: (
    title: string,
    message?: string,
    type?: ToastType,
    duration?: number
  ) => string;
  removeToast: (id: string) => void;
  simulateRealtimeEvent: () => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  }, []);

  const addToast = useCallback(
    (
      title: string,
      message?: string,
      type: ToastType = 'info',
      duration = 4500
    ) => {
      const id = `toast_${Date.now()}_${Math.random().toString(36).substring(2, 7)}`;
      const newToast: ToastItem = { id, title, message, type, duration };

      setToasts((prev) => [...prev, newToast]);

      if (duration > 0) {
        setTimeout(() => {
          removeToast(id);
        }, duration);
      }

      return id;
    },
    [removeToast]
  );

  const simulateRealtimeEvent = useCallback(() => {
    const events: Array<{ title: string; message: string; type: ToastType }> = [
      {
        title: 'Payment Received',
        message: 'Acme Corp paid $1,450.00 for Annual Enterprise Plan.',
        type: 'success',
      },
      {
        title: 'New Subscription Active',
        message: 'Starlight Labs upgraded to Growth Tier ($299/mo).',
        type: 'info',
      },
      {
        title: 'Invoice Past Due',
        message: 'Invoice #INV-2024-089 ($420.00) requires follow-up.',
        type: 'warning',
      },
      {
        title: 'Payment Attempt Failed',
        message: 'Card ending in •••• 4242 failed for customer Nexus Tech.',
        type: 'error',
      },
    ];

    const randomEvent = events[Math.floor(Math.random() * events.length)];
    addToast(randomEvent.title, randomEvent.message, randomEvent.type);
  }, [addToast]);

  // Background simulation engine: actively generates mock payments,
  // subscription state changes, and dunning updates without page refresh.
  useEffect(() => {
    const initialTimer = setTimeout(() => {
      simulateRealtimeEvent();
    }, 12000);

    const recurringTimer = setInterval(() => {
      simulateRealtimeEvent();
    }, 28000);

    return () => {
      clearTimeout(initialTimer);
      clearInterval(recurringTimer);
    };
  }, [simulateRealtimeEvent]);

  return (
    <ToastContext.Provider
      value={{
        toasts,
        addToast,
        removeToast,
        simulateRealtimeEvent,
      }}
    >
      {children}
      {/* Toast viewport fixed top-right */}
      <div
        className="fixed top-4 right-4 z-50 flex flex-col gap-2.5 max-w-sm w-full pointer-events-none"
        aria-live="polite"
      >
        {toasts.map((toast) => (
          <div
            key={toast.id}
            className="pointer-events-auto flex items-start gap-3 p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-md transition-all duration-200 transform translate-y-0"
          >
            <div className="shrink-0 mt-0.5">
              {toast.type === 'success' && (
                <CheckCircle2 className="w-5 h-5 text-emerald-600" />
              )}
              {toast.type === 'error' && (
                <XCircle className="w-5 h-5 text-stripe-ruby" />
              )}
              {toast.type === 'warning' && (
                <AlertTriangle className="w-5 h-5 text-amber-500" />
              )}
              {toast.type === 'info' && (
                <Info className="w-5 h-5 text-stripe-primary" />
              )}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-stripe-ink">{toast.title}</p>
              {toast.message && (
                <p className="text-xs text-stripe-ink-mute mt-0.5 leading-relaxed">
                  {toast.message}
                </p>
              )}
            </div>
            <button
              onClick={() => removeToast(toast.id)}
              className="shrink-0 text-slate-400 hover:text-stripe-ink p-1 rounded-md transition-colors"
              aria-label="Dismiss notification"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
};

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
}

