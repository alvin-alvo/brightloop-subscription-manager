import React from 'react';
import { Plus } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const SubscriptionsPage: React.FC = () => {
  const { addToast } = useToast();

  const subscriptions = [
    {
      id: 'sub_9921',
      customer: 'Acme Enterprises Inc.',
      plan: 'Enterprise Annual',
      billingCycle: 'Yearly',
      nextBilling: 'Jan 15, 2025',
      mrr: '$1,200.00',
      status: 'active' as const,
    },
    {
      id: 'sub_9920',
      customer: 'Vertex Technologies',
      plan: 'Growth Tier',
      billingCycle: 'Monthly',
      nextBilling: 'Nov 14, 2024',
      mrr: '$499.00',
      status: 'active' as const,
    },
    {
      id: 'sub_9919',
      customer: 'Quantum AI Systems',
      plan: 'Scale Tier',
      billingCycle: 'Monthly',
      nextBilling: 'Nov 18, 2024',
      mrr: '$890.00',
      status: 'active' as const,
    },
    {
      id: 'sub_9918',
      customer: 'Helix Health Labs',
      plan: 'Pro Tier',
      billingCycle: 'Monthly',
      nextBilling: 'Oct 28, 2024',
      mrr: '$299.00',
      status: 'past_due' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Subscriptions Lifecycle Manager
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Active recurring contracts, renewal timelines, upgrades, and cancellations.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Subscription Form', 'Assign new subscription tier.', 'info')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          Assign Subscription
        </button>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Sub ID</th>
                <th className="py-3 px-5">Subscriber</th>
                <th className="py-3 px-5">Plan</th>
                <th className="py-3 px-5">Cycle</th>
                <th className="py-3 px-5">Monthly Value</th>
                <th className="py-3 px-5">Next Renewal</th>
                <th className="py-3 px-5">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {subscriptions.map((s) => (
                <tr key={s.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs font-semibold text-stripe-ink-secondary">
                    {s.id}
                  </td>
                  <td className="py-3 px-5 font-medium text-stripe-ink">{s.customer}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-secondary">{s.plan}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-mute">{s.billingCycle}</td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">{s.mrr}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-mute tnum">{s.nextBilling}</td>
                  <td className="py-3 px-5">
                    <StatusBadge
                      label={s.status === 'active' ? 'Active' : 'Past Due'}
                      variant={s.status === 'active' ? 'success' : 'danger'}
                    />
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

export default SubscriptionsPage;
