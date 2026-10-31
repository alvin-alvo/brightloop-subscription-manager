import React from 'react';
import { Plus, Send } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const QuotationsPage: React.FC = () => {
  const { addToast } = useToast();

  const quotations = [
    {
      id: 'QT-2024-0041',
      customer: 'Acme Enterprises Inc.',
      scope: 'Annual Enterprise + 20 Dedicated Seats',
      amount: '$16,500.00',
      validUntil: 'Nov 15, 2024',
      status: 'pending' as const,
    },
    {
      id: 'QT-2024-0040',
      customer: 'Quantum AI Systems',
      scope: 'Scale Tier Migration + Custom SLA',
      amount: '$9,200.00',
      validUntil: 'Nov 10, 2024',
      status: 'approved' as const,
    },
    {
      id: 'QT-2024-0039',
      customer: 'Pinnacle Logistics',
      scope: 'Starter Tier Trial Extension',
      amount: '$1,188.00',
      validUntil: 'Oct 30, 2024',
      status: 'expired' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Quotations & Custom Proposals
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Generate customized commercial proposals, terms, and convert quotes directly into live subscriptions.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Draft Initialized', 'New commercial proposal created.', 'info')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          Create Quotation
        </button>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Quote ID</th>
                <th className="py-3 px-5">Customer</th>
                <th className="py-3 px-5">Proposal Scope</th>
                <th className="py-3 px-5">Proposed Total</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {quotations.map((q) => (
                <tr key={q.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs font-medium text-stripe-primary">
                    {q.id}
                  </td>
                  <td className="py-3 px-5 font-medium text-stripe-ink">{q.customer}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-secondary">{q.scope}</td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">{q.amount}</td>
                  <td className="py-3 px-5">
                    <StatusBadge
                      label={
                        q.status === 'approved'
                          ? 'Approved'
                          : q.status === 'pending'
                          ? 'Pending Signature'
                          : 'Expired'
                      }
                      variant={
                        q.status === 'approved'
                          ? 'success'
                          : q.status === 'pending'
                          ? 'warning'
                          : 'neutral'
                      }
                    />
                  </td>
                  <td className="py-3 px-5 text-right">
                    <button
                      onClick={() =>
                        addToast('Proposal Dispatched', `Quote ${q.id} sent to client.`, 'success')
                      }
                      className="inline-flex items-center gap-1 text-xs text-stripe-primary hover:underline"
                    >
                      <Send className="w-3 h-3" />
                      Send to Sign
                    </button>
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

export default QuotationsPage;
