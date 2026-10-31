import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import {
  CreditCard,
  Calendar,
  AlertCircle,
  CheckCircle2,
  ArrowRight,
  ShieldCheck,
  Building,
  ExternalLink,
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { DashboardSkeleton } from '../../components/layout/DashboardSkeleton';
import {
  fetchCustomerOverview,
  formatDisplayDate,
} from '../../services/mockCustomerData';

interface CustomerDashboardPageProps {
  customerId?: string;
}

export const CustomerDashboardPage: React.FC<CustomerDashboardPageProps> = ({
  customerId: propCustomerId,
}) => {
  const { user } = useAuth();
  // Strict data isolation: prioritize passed prop, fallback to logged-in user ID
  const effectiveCustomerId = propCustomerId || user?.id || 'usr_cust_03';

  const { data: metrics, isLoading } = useQuery({
    queryKey: ['customerOverview', effectiveCustomerId],
    queryFn: () => fetchCustomerOverview(effectiveCustomerId),
    staleTime: 1000 * 60 * 2,
  });

  if (isLoading || !metrics) {
    return <DashboardSkeleton />;
  }

  const hasOutstanding = metrics.totalOutstandingAmount > 0;

  return (
    <div className="space-y-6">
      {/* Welcome Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2 border-b border-stripe-hairline">
        <div>
          <div className="inline-flex items-center gap-2 px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200 mb-2">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
            Account in Good Standing
          </div>
          <h1 className="text-2xl sm:text-3xl font-light text-stripe-ink tracking-tight">
            Welcome back, {user?.name || 'Customer'}
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1 flex items-center gap-1.5">
            <Building className="w-4 h-4 text-slate-400" />
            <span>{user?.companyName || 'Acme Enterprises Inc.'}</span>
            <span>•</span>
            <span>Customer ID: <code className="font-mono text-xs text-stripe-ink-secondary">{effectiveCustomerId}</code></span>
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Link
            to="/customer/invoices"
            className="px-4 py-2 text-xs sm:text-sm font-medium text-stripe-ink bg-white border border-stripe-hairline hover:bg-slate-50 rounded-full shadow-stripe-sm transition-colors"
          >
            Billing Statements
          </Link>
          <Link
            to="/customer/subscriptions"
            className="px-4 py-2 text-xs sm:text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
          >
            My Subscriptions
          </Link>
        </div>
      </div>

      {/* 3 Top-Level Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
        {/* Card 1: Active Subscriptions */}
        <div className="p-6 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm hover:shadow-stripe-md transition-shadow">
          <div className="flex items-center justify-between text-stripe-ink-mute mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider">
              Active Subscriptions
            </span>
            <div className="p-2 rounded-lg bg-indigo-50 text-stripe-primary">
              <CreditCard className="w-4 h-4" />
            </div>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-light text-stripe-ink tnum">
              {metrics.activeSubscriptionsCount}
            </span>
            <span className="text-xs text-stripe-ink-mute">contracted tiers</span>
          </div>
          <p className="text-xs text-stripe-ink-secondary mt-2 tnum">
            Estimated ~${Math.round(metrics.totalMonthlyValue).toLocaleString()}/mo run-rate
          </p>
          <div className="mt-4 pt-3 border-t border-stripe-hairline">
            <Link
              to="/customer/subscriptions"
              className="inline-flex items-center gap-1 text-xs font-medium text-stripe-primary hover:text-stripe-primary-deep"
            >
              Inspect Active Licenses
              <ArrowRight className="w-3.5 h-3.5" />
            </Link>
          </div>
        </div>

        {/* Card 2: Next Payment Due */}
        <div className="p-6 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm hover:shadow-stripe-md transition-shadow">
          <div className="flex items-center justify-between text-stripe-ink-mute mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider">
              Next Payment Due
            </span>
            <div className="p-2 rounded-lg bg-slate-50 text-stripe-ink-secondary">
              <Calendar className="w-4 h-4" />
            </div>
          </div>
          {metrics.nextPaymentDue ? (
            <>
              <div className="text-3xl font-light text-stripe-ink tnum">
                ${metrics.nextPaymentDue.amount.toLocaleString(undefined, {
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                })}
              </div>
              <p className="text-xs text-stripe-ink-secondary mt-2">
                Renewal on <span className="font-medium text-stripe-ink tnum">{formatDisplayDate(metrics.nextPaymentDue.date)}</span>
              </p>
              <div className="mt-4 pt-3 border-t border-stripe-hairline flex items-center justify-between text-xs text-stripe-ink-mute">
                <span className="truncate max-w-[170px]">{metrics.nextPaymentDue.planName}</span>
                <span className="font-mono text-[11px]">{metrics.nextPaymentDue.subscriptionNumber}</span>
              </div>
            </>
          ) : (
            <p className="text-sm text-stripe-ink-mute mt-2">No upcoming payments scheduled.</p>
          )}
        </div>

        {/* Card 3: Total Outstanding */}
        <div
          className={`p-6 rounded-xl border transition-shadow ${
            hasOutstanding
              ? 'bg-amber-50/40 border-amber-200 shadow-stripe-sm'
              : 'bg-white border-stripe-hairline shadow-stripe-sm'
          }`}
        >
          <div className="flex items-center justify-between mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider text-stripe-ink-mute">
              Total Outstanding
            </span>
            <div
              className={`p-2 rounded-lg ${
                hasOutstanding
                  ? 'bg-amber-100 text-amber-700'
                  : 'bg-emerald-50 text-emerald-600'
              }`}
            >
              {hasOutstanding ? (
                <AlertCircle className="w-4 h-4" />
              ) : (
                <CheckCircle2 className="w-4 h-4" />
              )}
            </div>
          </div>

          <div className="text-3xl font-light text-stripe-ink tnum">
            ${metrics.totalOutstandingAmount.toLocaleString(undefined, {
              minimumFractionDigits: 2,
              maximumFractionDigits: 2,
            })}
          </div>

          <p className="text-xs mt-2 text-stripe-ink-secondary">
            {hasOutstanding
              ? `${metrics.outstandingInvoicesCount} invoice requires payment confirmation`
              : 'All statements are paid in full'}
          </p>

          <div className="mt-4 pt-3 border-t border-stripe-hairline">
            {hasOutstanding ? (
              <Link
                to="/customer/invoices"
                className="inline-flex items-center gap-1.5 text-xs font-medium text-amber-800 hover:text-amber-900 underline underline-offset-2"
              >
                Resolve Pending Invoices
                <ArrowRight className="w-3.5 h-3.5" />
              </Link>
            ) : (
              <span className="inline-flex items-center gap-1 text-xs text-emerald-700 font-medium">
                <ShieldCheck className="w-3.5 h-3.5" />
                Zero Balance Due
              </span>
            )}
          </div>
        </div>
      </div>

      {/* Account Overview & Self-Service Shortcuts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile Card */}
        <div className="p-6 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm lg:col-span-1 space-y-4">
          <h2 className="text-sm font-semibold text-stripe-ink uppercase tracking-wider">
            Commercial Account Profile
          </h2>
          <div className="space-y-3 text-xs text-stripe-ink-secondary">
            <div>
              <span className="text-stripe-ink-mute block">Company Organization</span>
              <span className="font-medium text-stripe-ink text-sm">
                {user?.companyName || 'Acme Enterprises Inc.'}
              </span>
            </div>
            <div>
              <span className="text-stripe-ink-mute block">Contract Billing Representative</span>
              <span className="font-medium text-stripe-ink">{user?.name} ({user?.email})</span>
            </div>
            <div>
              <span className="text-stripe-ink-mute block">Standard Payment Terms</span>
              <span className="font-medium text-stripe-ink">Net 30 via Automated ACH / Wire</span>
            </div>
            <div>
              <span className="text-stripe-ink-mute block">Settlement Currency</span>
              <span className="font-medium text-stripe-ink">USD ($)</span>
            </div>
          </div>
        </div>

        {/* Quick Self-Service Portal Guide */}
        <div className="p-6 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm lg:col-span-2 flex flex-col justify-between">
          <div>
            <div className="flex items-center justify-between">
              <h2 className="text-sm font-semibold text-stripe-ink uppercase tracking-wider">
                Self-Service Capabilities
              </h2>
              <span className="text-xs text-stripe-primary font-medium">Customer Portal</span>
            </div>
            <p className="text-xs text-stripe-ink-mute mt-2 leading-relaxed">
              Use your portal to review live software licenses, verify seat allocations, inspect historical tax invoices with discount line items, and instantly settle open statements without support tickets.
            </p>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mt-4">
              <Link
                to="/customer/subscriptions"
                className="p-3 rounded-lg border border-stripe-hairline bg-slate-50 hover:bg-slate-100/80 transition-colors flex items-center justify-between"
              >
                <div>
                  <div className="font-medium text-stripe-ink text-xs">Contract Subscriptions</div>
                  <div className="text-[11px] text-stripe-ink-mute">Inspect term dates & seats</div>
                </div>
                <ExternalLink className="w-3.5 h-3.5 text-slate-400" />
              </Link>

              <Link
                to="/customer/invoices"
                className="p-3 rounded-lg border border-stripe-hairline bg-slate-50 hover:bg-slate-100/80 transition-colors flex items-center justify-between"
              >
                <div>
                  <div className="font-medium text-stripe-ink text-xs">PDF Invoices & Tax Receipts</div>
                  <div className="text-[11px] text-stripe-ink-mute">Download statements</div>
                </div>
                <ExternalLink className="w-3.5 h-3.5 text-slate-400" />
              </Link>
            </div>
          </div>

          <div className="pt-4 mt-4 border-t border-stripe-hairline text-xs text-stripe-ink-mute flex items-center justify-between">
            <span>Need commercial contract amendments?</span>
            <span className="text-stripe-primary hover:underline cursor-pointer">
              Contact Enterprise Account Manager
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CustomerDashboardPage;

