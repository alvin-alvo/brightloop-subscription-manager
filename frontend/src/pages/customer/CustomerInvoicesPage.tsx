import React from 'react';
import { Download } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const CustomerInvoicesPage: React.FC = () => {
  const { addToast } = useToast();

  const invoices = [
    {
      id: 'INV-2024-001',
      date: 'Jan 15, 2024',
      period: 'Jan 15, 2024 - Jan 14, 2025',
      amount: '$14,400.00',
      status: 'paid' as const,
    },
    {
      id: 'INV-2023-012',
      date: 'Jan 15, 2023',
      period: 'Jan 15, 2023 - Jan 14, 2024',
      amount: '$12,000.00',
      status: 'paid' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Billing History & Tax Invoices
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Download compliant PDF tax invoices and inspect historical statement records.
          </p>
        </div>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Invoice Reference</th>
                <th className="py-3 px-5">Billing Period</th>
                <th className="py-3 px-5">Amount Paid</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">PDF Receipt</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {invoices.map((inv) => (
                <tr key={inv.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs font-semibold text-stripe-ink">
                    {inv.id}
                  </td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-secondary tnum">
                    {inv.period}
                  </td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">
                    {inv.amount}
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge label="Paid in Full" variant="success" />
                  </td>
                  <td className="py-3 px-5 text-right">
                    <button
                      onClick={() =>
                        addToast('PDF Downloaded', `Invoice ${inv.id} receipt generated.`, 'success')
                      }
                      className="inline-flex items-center gap-1.5 px-3 py-1 text-xs font-medium text-stripe-ink hover:text-stripe-primary bg-slate-50 hover:bg-slate-100 rounded-md transition-colors"
                    >
                      <Download className="w-3.5 h-3.5" />
                      Download
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

export default CustomerInvoicesPage;
