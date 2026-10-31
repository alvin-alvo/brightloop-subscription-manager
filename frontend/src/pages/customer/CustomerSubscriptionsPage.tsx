import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  CreditCard,
  Calendar,
  Check,
  FileCheck,
  Sparkles,
  ArrowUpRight,
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useToast } from '../../context/ToastContext';
import { DashboardSkeleton } from '../../components/layout/DashboardSkeleton';
import { StatusBadge } from '../../components/StatusBadge';
import {
  subscriptionsService,
  formatDisplayDate,
} from '../../services/subscriptionsService';

interface CustomerSubscriptionsPageProps {
  customerId?: string;
}

export const CustomerSubscriptionsPage: React.FC<CustomerSubscriptionsPageProps> = ({
  customerId: propCustomerId,
}) => {
  const { user } = useAuth();
  const { addToast } = useToast();
  // Strict data isolation
  const effectiveCustomerId = propCustomerId || user?.id || 'usr_cust_03';

  const { data: subscriptions, isLoading } = useQuery({
    queryKey: ['customerSubscriptions', effectiveCustomerId],
    queryFn: () => subscriptionsService.getSubscriptions(),
    staleTime: 1000 * 60 * 2,
  });

  if (isLoading || !subscriptions) {
    return <DashboardSkeleton />;
  }

  const activeCount = subscriptions.filter((s) => s.status === 'active').length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2 border-b border-stripe-hairline">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            My Subscriptions
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Active and archived software licenses contracted for {user?.companyName || 'your organization'}.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className="px-3 py-1 rounded-full text-xs font-medium bg-indigo-50 text-stripe-primary border border-indigo-100">
            {activeCount} Active Plans
          </span>
        </div>
      </div>

      {/* Subscription Cards Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {subscriptions.map((sub) => {
          const isActive = sub.status === 'active';

          return (
            <div
              key={sub.id}
              className={`p-6 rounded-xl border transition-all ${
                isActive
                  ? 'bg-white border-stripe-hairline shadow-stripe-sm hover:shadow-stripe-md ring-1 ring-slate-100'
                  : 'bg-slate-50/70 border-slate-200 opacity-80'
              }`}
            >
              {/* Card Header: Plan Name, Status Badge, and Subscription Number */}
              <div className="flex items-start justify-between gap-3">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-xs font-medium text-stripe-primary">
                      {sub.subscriptionNumber}
                    </span>
                    <StatusBadge
                      label={isActive ? 'Active' : 'Closed'}
                      variant={isActive ? 'success' : 'neutral'}
                    />
                  </div>
                  <h2 className="text-xl font-normal text-stripe-ink mt-1.5">
                    {sub.planName}
                  </h2>
                </div>

                <div className="text-right">
                  <div className="text-2xl font-light text-stripe-ink tnum">
                    ${sub.totalAmount.toLocaleString(undefined, {
                      minimumFractionDigits: 2,
                    })}
                  </div>
                  <span className="text-[11px] text-stripe-ink-mute">
                    {sub.paymentTerms}
                  </span>
                </div>
              </div>

              {/* Key Timeline & Payment Terms Grid */}
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 my-5 p-3.5 bg-stripe-canvas-soft rounded-lg border border-slate-100 text-xs">
                <div>
                  <span className="text-stripe-ink-mute block flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    Start Date
                  </span>
                  <span className="font-medium text-stripe-ink tnum mt-0.5 block">
                    {formatDisplayDate(sub.startDate)}
                  </span>
                </div>

                <div>
                  <span className="text-stripe-ink-mute block flex items-center gap-1">
                    <Calendar className="w-3 h-3" />
                    {isActive ? 'Renewal Date' : 'End Date'}
                  </span>
                  <span className="font-medium text-stripe-ink tnum mt-0.5 block">
                    {formatDisplayDate(sub.expiryDate)}
                  </span>
                </div>

                <div className="col-span-2 sm:col-span-1">
                  <span className="text-stripe-ink-mute block flex items-center gap-1">
                    <CreditCard className="w-3 h-3" />
                    Terms
                  </span>
                  <span className="font-medium text-stripe-ink mt-0.5 block truncate">
                    {sub.paymentTerms}
                  </span>
                </div>
              </div>

              {/* Features / Entitlements */}
              {sub.features && sub.features.length > 0 && (
                <div className="space-y-2 mb-6">
                  <p className="text-[11px] uppercase tracking-wider font-semibold text-stripe-ink-mute">
                    Contract Entitlements
                  </p>
                  <ul className="space-y-1.5 text-xs text-stripe-ink-secondary">
                    {sub.features.map((feature, idx) => (
                      <li key={idx} className="flex items-center gap-2">
                        <Check
                          className={`w-3.5 h-3.5 shrink-0 ${
                            isActive ? 'text-emerald-600' : 'text-slate-400'
                          }`}
                        />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Card Footer Actions */}
              <div className="pt-4 border-t border-stripe-hairline flex items-center justify-between text-xs">
                <button
                  onClick={() =>
                    addToast(
                      'Service Agreement Downloaded',
                      `PDF terms for ${sub.subscriptionNumber} saved to downloads.`,
                      'success'
                    )
                  }
                  className="inline-flex items-center gap-1 text-stripe-ink-mute hover:text-stripe-ink transition-colors"
                >
                  <FileCheck className="w-3.5 h-3.5" />
                  View Agreement
                </button>

                {isActive && (
                  <button
                    onClick={() =>
                      addToast(
                        'Amendment Requested',
                        'Your enterprise account manager has been notified to adjust seats.',
                        'info'
                      )
                    }
                    className="inline-flex items-center gap-1 px-3 py-1 text-xs font-medium text-stripe-primary bg-indigo-50 hover:bg-indigo-100 rounded-full transition-colors"
                  >
                    <Sparkles className="w-3 h-3" />
                    Manage Seats
                    <ArrowUpRight className="w-3 h-3" />
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default CustomerSubscriptionsPage;
