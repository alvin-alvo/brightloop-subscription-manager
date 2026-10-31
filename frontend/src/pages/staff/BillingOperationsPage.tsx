import React from 'react';
import { RefreshCw } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const BillingOperationsPage: React.FC = () => {
  const { addToast } = useToast();

  const billingBatches = [
    {
      id: 'BATCH-2024-10-W4',
      date: 'Today, 00:00 UTC',
      itemsCount: 142,
      totalBilled: '$38,450.00',
      successRate: '98.6%',
      status: 'completed' as const,
    },
    {
      id: 'BATCH-2024-10-W3',
      date: 'Oct 24, 2024',
      itemsCount: 139,
      totalBilled: '$36,120.00',
      successRate: '99.2%',
      status: 'completed' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Billing & Invoicing Operations
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Scheduled recurring invoice generation, charge automation, and dunning retry management.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Billing Engine Triggered', 'Manual billing batch run started.', 'info')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <RefreshCw className="w-4 h-4" />
          Run Recurring Billing
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal uppercase text-stripe-ink-mute">
            Auto-Charge Status
          </span>
          <div className="flex items-center gap-2 mt-2">
            <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <span className="text-lg font-light text-stripe-ink">Autonomous Scheduler</span>
          </div>
          <p className="text-xs text-stripe-ink-mute mt-1">Runs daily at 00:00 UTC</p>
        </div>

        <div className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal uppercase text-stripe-ink-mute">
            Dunning Engine
          </span>
          <div className="text-lg font-light text-stripe-ink mt-2">
            Smart Retries (3 Attempts)
          </div>
          <p className="text-xs text-stripe-ink-mute mt-1">87% recovery rate on expired cards</p>
        </div>

        <div className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal uppercase text-stripe-ink-mute">
            Pending Invoices
          </span>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">12</div>
          <p className="text-xs text-stripe-ink-mute mt-1">Ready for next automated cycle</p>
        </div>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="p-4 border-b border-stripe-hairline">
          <h2 className="text-sm font-normal text-stripe-ink">Recurring Billing Runs</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Batch ID</th>
                <th className="py-3 px-5">Execution Date</th>
                <th className="py-3 px-5">Invoices Processed</th>
                <th className="py-3 px-5">Total Billed</th>
                <th className="py-3 px-5">Success Rate</th>
                <th className="py-3 px-5">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {billingBatches.map((b) => (
                <tr key={b.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs font-semibold text-stripe-ink-secondary">
                    {b.id}
                  </td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-mute tnum">{b.date}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink tnum">{b.itemsCount}</td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">
                    {b.totalBilled}
                  </td>
                  <td className="py-3 px-5 text-xs text-emerald-600 font-medium tnum">
                    {b.successRate}
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge label="Completed" variant="success" />
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

export default BillingOperationsPage;
