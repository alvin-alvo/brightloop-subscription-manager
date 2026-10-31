import React from 'react';
import { CreditCard, Calendar, Check, ArrowRight } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { useAuth } from '../../hooks/useAuth';

export const CustomerSubscriptionsPage: React.FC = () => {
  const { user } = useAuth();
  const { addToast } = useToast();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
          My Active Subscriptions
        </h1>
        <p className="text-sm text-stripe-ink-mute mt-1">
          Review your contracted software licenses, billing cycles, and payment methods.
        </p>
      </div>

      {/* Active Subscription Highlight Card (Stripe Featured Tier Style) */}
      <div className="p-6 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
        <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
          <div>
            <div className="flex items-center gap-2">
              <span className="text-xs uppercase font-medium tracking-wider text-stripe-primary">
                Current Plan
              </span>
              <StatusBadge label="Active" variant="success" />
            </div>
            <h2 className="text-xl font-normal text-stripe-ink mt-1">
              Enterprise Annual Tier
            </h2>
            <p className="text-xs text-stripe-ink-secondary mt-0.5">
              Assigned to: {user?.companyName || 'Acme Enterprises Inc.'}
            </p>
          </div>
          <div className="text-right">
            <div className="text-3xl font-light text-stripe-ink tnum">$14,400.00</div>
            <p className="text-xs text-stripe-ink-mute mt-0.5">Billed Annually</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6 pt-6 border-t border-stripe-hairline">
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-slate-50 text-stripe-ink-secondary">
              <Calendar className="w-4 h-4" />
            </div>
            <div>
              <p className="text-xs text-stripe-ink-mute">Next Renewal Date</p>
              <p className="text-sm font-medium text-stripe-ink tnum">January 15, 2025</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-slate-50 text-stripe-ink-secondary">
              <CreditCard className="w-4 h-4" />
            </div>
            <div>
              <p className="text-xs text-stripe-ink-mute">Payment Method</p>
              <p className="text-sm font-medium text-stripe-ink tnum">Mastercard ending in 9821</p>
            </div>
          </div>

          <div className="flex items-center justify-end">
            <button
              onClick={() =>
                addToast(
                  'Upgrade Request Submitted',
                  'Your account manager will contact you within 2 hours.',
                  'info'
                )
              }
              className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-stripe-primary bg-indigo-50 hover:bg-indigo-100 rounded-full transition-colors"
            >
              Modify Plan / Add Seats
              <ArrowRight className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        <div className="mt-6 pt-4 border-t border-stripe-hairline">
          <h3 className="text-xs font-medium text-stripe-ink uppercase tracking-wider mb-3">
            Included Entitlements
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-2 text-xs text-stripe-ink-secondary">
            <div className="flex items-center gap-2">
              <Check className="w-3.5 h-3.5 text-emerald-600" />
              <span>Dedicated Account Manager</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-3.5 h-3.5 text-emerald-600" />
              <span>99.99% Financial-grade SLA</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="w-3.5 h-3.5 text-emerald-600" />
              <span>Unlimited API Event Integrations</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CustomerSubscriptionsPage;

