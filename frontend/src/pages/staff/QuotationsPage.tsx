import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Plus, Send } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { fetchStaffWorkspace, type StaffQuote } from '../../utils/mockStaffData';

const currency = (cents: number): string =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100);

const statusVariant: Record<StaffQuote['status'], 'warning' | 'success' | 'neutral'> = {
  Pending: 'warning', Approved: 'success', Expired: 'neutral',
};

export const QuotationsPage: React.FC = () => {
  const { addToast } = useToast();
  const { data, isLoading } = useQuery({ queryKey: ['staff-workspace'], queryFn: fetchStaffWorkspace });

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div><p className="text-xs font-medium uppercase text-stripe-ink-mute">Sales operations</p><h1 className="mt-1 text-2xl font-light text-stripe-ink">Quotations</h1><p className="mt-1 text-sm text-stripe-ink-mute">Review proposals and turn approved terms into subscriptions.</p></div>
        <Link to="/staff/quotations/new" className="inline-flex items-center justify-center gap-2 rounded-full bg-stripe-primary px-4 py-2.5 text-sm font-medium text-white shadow-stripe-sm hover:bg-stripe-primary-deep"><Plus className="h-4 w-4" /> Create quotation</Link>
      </div>

      <section className="overflow-hidden rounded-lg border border-stripe-hairline bg-white shadow-stripe-sm">
        <div className="border-b border-stripe-hairline px-5 py-4"><h2 className="text-sm font-medium text-stripe-ink">Recent proposals</h2><p className="mt-1 text-xs text-stripe-ink-mute">Amounts include applicable tax.</p></div>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[780px] text-left text-sm">
            <thead className="bg-stripe-canvas-soft text-[11px] uppercase text-stripe-ink-mute"><tr><th className="px-5 py-3 font-medium">Quote ID</th><th className="px-5 py-3 font-medium">Customer</th><th className="px-5 py-3 font-medium">Plan</th><th className="px-5 py-3 font-medium">Total</th><th className="px-5 py-3 font-medium">Valid until</th><th className="px-5 py-3 font-medium">Status</th><th className="px-5 py-3 text-right font-medium">Action</th></tr></thead>
            <tbody className="divide-y divide-stripe-hairline">
              {isLoading ? <tr><td colSpan={7} className="px-5 py-12 text-center text-sm text-stripe-ink-mute">Loading proposals...</td></tr> : data?.quotes.map((quote) => <tr key={quote.id} className="hover:bg-stripe-canvas-soft/70"><td className="px-5 py-4 font-mono text-xs font-medium text-stripe-primary">{quote.id}</td><td className="px-5 py-4 font-medium text-stripe-ink">{quote.customer}</td><td className="px-5 py-4 text-stripe-ink-secondary">{quote.plan}</td><td className="px-5 py-4 font-medium tabular-nums text-stripe-ink">{currency(quote.amountCents)}</td><td className="px-5 py-4 text-stripe-ink-secondary">{new Date(`${quote.validUntil}T00:00:00`).toLocaleDateString()}</td><td className="px-5 py-4"><StatusBadge label={quote.status} variant={statusVariant[quote.status]} /></td><td className="px-5 py-4 text-right"><button type="button" onClick={() => addToast('Proposal ready to send', `${quote.id} for ${quote.customer} is queued for sharing.`, 'success')} className="inline-flex items-center gap-1.5 text-xs font-medium text-stripe-primary hover:underline"><Send className="h-3.5 w-3.5" /> Send proposal</button></td></tr>)}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
};

export default QuotationsPage;
