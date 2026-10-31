import React from 'react';
import { DollarSign, TrendingUp, Users, CreditCard, ArrowUpRight } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const FinanceDashboardPage: React.FC = () => {
  const { addToast } = useToast();

  const metrics = [
    {
      title: 'Monthly Recurring Revenue',
      value: '$48,250.00',
      change: '+14.2%',
      trend: 'up',
      subtitle: 'vs. $42,240.00 last month',
      icon: DollarSign,
    },
    {
      title: 'Active Subscriptions',
      value: '1,284',
      change: '+8.1%',
      trend: 'up',
      subtitle: '94 net new this month',
      icon: CreditCard,
    },
    {
      title: 'Net Retention Rate',
      value: '108.4%',
      change: '+2.1%',
      trend: 'up',
      subtitle: 'Target: >105%',
      icon: TrendingUp,
    },
    {
      title: 'Total Customers',
      value: '842',
      change: '+5.4%',
      trend: 'up',
      subtitle: 'Enterprise & Pro accounts',
      icon: Users,
    },
  ];

  const recentTransactions = [
    {
      id: 'tx_98124',
      customer: 'Acme Enterprises Inc.',
      plan: 'Enterprise Annual',
      amount: '$14,400.00',
      status: 'success' as const,
      date: 'Today, 10:14 AM',
    },
    {
      id: 'tx_98123',
      customer: 'Vertex Technologies',
      plan: 'Growth Tier (Monthly)',
      amount: '$499.00',
      status: 'success' as const,
      date: 'Today, 09:42 AM',
    },
    {
      id: 'tx_98122',
      customer: 'Quantum AI Systems',
      plan: 'Scale Tier (Annual)',
      amount: '$8,900.00',
      status: 'success' as const,
      date: 'Yesterday',
    },
    {
      id: 'tx_98121',
      customer: 'Helix Health',
      plan: 'Pro Monthly',
      amount: '$299.00',
      status: 'warning' as const,
      date: 'Oct 29, 2024',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Page Title & Action Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Finance & Revenue Dashboard
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Real-time financial performance, recurring metrics, and settlement streams.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() =>
              addToast(
                'Financial Report Exported',
                'Q4 MRR and cash reconciliation report generated.',
                'success'
              )
            }
            className="px-4 py-2 text-sm font-normal text-stripe-ink bg-white border border-stripe-hairline hover:bg-slate-50 rounded-full shadow-stripe-sm transition-colors"
          >
            Export Ledger
          </button>
          <button
            onClick={() =>
              addToast(
                'Payout Scheduled',
                'Next automatic payout scheduled for Nov 1, 2024.',
                'info'
              )
            }
            className="px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
          >
            Create Payout
          </button>
        </div>
      </div>

      {/* Metric Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {metrics.map((metric) => {
          const Icon = metric.icon;
          return (
            <div
              key={metric.title}
              className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm hover:shadow-stripe-md transition-shadow"
            >
              <div className="flex items-center justify-between text-stripe-ink-mute mb-2">
                <span className="text-xs font-normal uppercase tracking-wider">
                  {metric.title}
                </span>
                <div className="p-1.5 rounded-md bg-slate-50 text-stripe-ink-secondary">
                  <Icon className="w-4 h-4" />
                </div>
              </div>
              <div className="flex items-baseline gap-2">
                <span className="text-2xl font-light text-stripe-ink tracking-tight tnum">
                  {metric.value}
                </span>
                <span className="inline-flex items-center text-xs font-medium text-emerald-600 tnum">
                  <ArrowUpRight className="w-3 h-3 mr-0.5" />
                  {metric.change}
                </span>
              </div>
              <p className="text-xs text-stripe-ink-mute mt-1.5 tnum">
                {metric.subtitle}
              </p>
            </div>
          );
        })}
      </div>

      {/* Live Transaction Ledger Card */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="p-5 border-b border-stripe-hairline flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h2 className="text-base font-normal text-stripe-ink">
              Recent Settlement Stream
            </h2>
            <p className="text-xs text-stripe-ink-mute mt-0.5">
              Live transactions processed across active subscription contracts.
            </p>
          </div>
          <span className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-normal bg-indigo-50 text-stripe-primary border border-indigo-100">
            Realtime Sync Active
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Transaction ID</th>
                <th className="py-3 px-5">Customer</th>
                <th className="py-3 px-5">Plan</th>
                <th className="py-3 px-5">Amount</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Date</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {recentTransactions.map((tx) => (
                <tr key={tx.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs text-stripe-ink-secondary">
                    {tx.id}
                  </td>
                  <td className="py-3 px-5 font-medium text-stripe-ink">
                    {tx.customer}
                  </td>
                  <td className="py-3 px-5 text-stripe-ink-secondary text-xs">
                    {tx.plan}
                  </td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">
                    {tx.amount}
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge
                      label={tx.status === 'success' ? 'Settled' : 'Pending'}
                      variant={tx.status === 'success' ? 'success' : 'warning'}
                    />
                  </td>
                  <td className="py-3 px-5 text-right text-xs text-stripe-ink-mute tnum">
                    {tx.date}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default FinanceDashboardPage;

