import React, { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Check, CreditCard, X } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { fetchStaffWorkspace, type StaffInvoice } from '../../utils/mockStaffData';

type InvoiceFilter = 'Outstanding' | 'Paid';
interface PaymentRecord {
  invoiceNumber: string;
  amountCents: number;
  method: string;
  date: string;
  note: string;
}

const currency = (cents: number): string =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100);

export const BillingOperationsPage: React.FC = () => {
  const { addToast } = useToast();
  const { data, isLoading } = useQuery({ queryKey: ['staff-workspace'], queryFn: fetchStaffWorkspace });
  const [invoices, setInvoices] = useState<StaffInvoice[]>([]);
  const [filter, setFilter] = useState<InvoiceFilter>('Outstanding');
  const [selectedInvoice, setSelectedInvoice] = useState<StaffInvoice | null>(null);
  const [payments, setPayments] = useState<PaymentRecord[]>([]);
  const [method, setMethod] = useState('Card');
  const [amount, setAmount] = useState('');
  const [paymentDate, setPaymentDate] = useState(new Date().toISOString().slice(0, 10));
  const [note, setNote] = useState('');

  useEffect(() => {
    if (data) setInvoices(data.invoices);
  }, [data]);

  const openPayment = (invoice: StaffInvoice) => {
    setSelectedInvoice(invoice);
    setAmount((invoice.amountCents / 100).toFixed(2));
    setPaymentDate(new Date().toISOString().slice(0, 10));
    setNote('');
  };

  const recordPayment = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selectedInvoice) return;
    const amountCents = Math.round(Number(amount) * 100);
    if (!Number.isFinite(amountCents) || amountCents <= 0 || amountCents > selectedInvoice.amountCents) return;
    setPayments((current) => [...current, { invoiceNumber: selectedInvoice.invoiceNumber, amountCents, method, date: paymentDate, note }]);
    setInvoices((current) => current.map((invoice) => {
      if (invoice.id !== selectedInvoice.id) return invoice;
      const remaining = invoice.amountCents - amountCents;
      return { ...invoice, amountCents: remaining, status: remaining === 0 ? 'Paid' : 'Outstanding', method: remaining === 0 ? method : undefined };
    }));
    addToast('Payment recorded', `${currency(amountCents)} applied to ${selectedInvoice.invoiceNumber}.`, 'success');
    setSelectedInvoice(null);
  };

  const visibleInvoices = invoices.filter((invoice) => invoice.status === filter);
  const outstandingTotal = invoices.filter((invoice) => invoice.status === 'Outstanding').reduce((total, invoice) => total + invoice.amountCents, 0);
  const paidTotal = invoices.filter((invoice) => invoice.status === 'Paid').reduce((total, invoice) => total + invoice.amountCents, 0);

  return (
    <div className="space-y-6">
      <div><p className="text-xs font-medium uppercase text-stripe-ink-mute">Accounts receivable</p><h1 className="mt-1 text-2xl font-light text-stripe-ink">Billing & payments</h1><p className="mt-1 text-sm text-stripe-ink-mute">Track invoice status and record payments against the correct account.</p></div>

      <div className="grid gap-4 sm:grid-cols-2"><section className="rounded-lg border border-stripe-hairline bg-white p-5 shadow-stripe-sm"><p className="text-xs uppercase text-stripe-ink-mute">Outstanding balance</p><p className="mt-2 text-2xl font-light tabular-nums text-stripe-ink">{currency(outstandingTotal)}</p><p className="mt-1 text-xs text-stripe-ink-mute">{invoices.filter((invoice) => invoice.status === 'Outstanding').length} invoices to follow up</p></section><section className="rounded-lg border border-stripe-hairline bg-white p-5 shadow-stripe-sm"><p className="text-xs uppercase text-stripe-ink-mute">Paid invoices</p><p className="mt-2 text-2xl font-light tabular-nums text-stripe-ink">{currency(paidTotal)}</p><p className="mt-1 text-xs text-stripe-ink-mute">{invoices.filter((invoice) => invoice.status === 'Paid').length} invoices settled</p></section></div>

      <section className="overflow-hidden rounded-lg border border-stripe-hairline bg-white shadow-stripe-sm">
        <div className="flex flex-col gap-3 border-b border-stripe-hairline p-4 sm:flex-row sm:items-center sm:justify-between"><div><h2 className="text-sm font-medium text-stripe-ink">Invoices</h2><p className="mt-1 text-xs text-stripe-ink-mute">{payments.length} payment{payments.length === 1 ? '' : 's'} recorded this session</p></div><div className="inline-flex w-fit rounded-md border border-stripe-hairline p-1" role="group" aria-label="Filter invoices">{(['Outstanding', 'Paid'] as const).map((status) => <button key={status} type="button" onClick={() => setFilter(status)} aria-pressed={filter === status} className={`rounded px-3 py-1.5 text-xs font-medium ${filter === status ? 'bg-stripe-primary text-white' : 'text-stripe-ink-secondary hover:bg-slate-50'}`}>{status}</button>)}</div></div>
        <div className="overflow-x-auto"><table className="w-full min-w-[760px] text-left text-sm"><thead className="bg-stripe-canvas-soft text-[11px] uppercase text-stripe-ink-mute"><tr><th className="px-5 py-3 font-medium">Invoice</th><th className="px-5 py-3 font-medium">Customer</th><th className="px-5 py-3 font-medium">Issued</th><th className="px-5 py-3 font-medium">Due date</th><th className="px-5 py-3 font-medium">Amount due</th><th className="px-5 py-3 font-medium">Status</th><th className="px-5 py-3 text-right font-medium">Action</th></tr></thead><tbody className="divide-y divide-stripe-hairline">{isLoading ? <tr><td colSpan={7} className="px-5 py-12 text-center text-sm text-stripe-ink-mute">Loading invoices...</td></tr> : visibleInvoices.map((invoice) => {const overdue = invoice.status === 'Outstanding' && invoice.dueDate < new Date().toISOString().slice(0, 10); return <tr key={invoice.id} className="hover:bg-stripe-canvas-soft/70"><td className="px-5 py-4 font-mono text-xs font-medium text-stripe-primary">{invoice.invoiceNumber}</td><td className="px-5 py-4 font-medium text-stripe-ink">{invoice.customer}</td><td className="px-5 py-4 text-stripe-ink-secondary">{new Date(`${invoice.issuedAt}T00:00:00`).toLocaleDateString()}</td><td className={`px-5 py-4 ${overdue ? 'font-medium text-rose-700' : 'text-stripe-ink-secondary'}`}>{new Date(`${invoice.dueDate}T00:00:00`).toLocaleDateString()}{overdue && <span className="ml-2 text-[10px] uppercase">Overdue</span>}</td><td className="px-5 py-4 font-medium tabular-nums text-stripe-ink">{currency(invoice.amountCents)}</td><td className="px-5 py-4"><StatusBadge label={invoice.status} variant={invoice.status === 'Paid' ? 'success' : overdue ? 'danger' : 'warning'} /></td><td className="px-5 py-4 text-right">{invoice.status === 'Outstanding' ? <button type="button" onClick={() => openPayment(invoice)} className="inline-flex items-center gap-1.5 rounded-full border border-stripe-hairline px-3 py-1.5 text-xs font-medium text-stripe-primary hover:bg-indigo-50"><CreditCard className="h-3.5 w-3.5" /> Record payment</button> : <span className="inline-flex items-center gap-1 text-xs text-emerald-700"><Check className="h-3.5 w-3.5" />{invoice.method ?? 'Paid'}</span>}</td></tr>;})}{!isLoading && visibleInvoices.length === 0 && <tr><td colSpan={7} className="px-5 py-12 text-center text-sm text-stripe-ink-mute">No {filter.toLowerCase()} invoices.</td></tr>}</tbody></table></div>
      </section>

      {selectedInvoice && <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4" onMouseDown={(event) => { if (event.target === event.currentTarget) setSelectedInvoice(null); }}><section role="dialog" aria-modal="true" aria-labelledby="payment-modal-title" className="w-full max-w-lg rounded-lg border border-stripe-hairline bg-white shadow-stripe-lg"><div className="flex items-start justify-between border-b border-stripe-hairline p-5"><div><p className="font-mono text-xs text-stripe-primary">{selectedInvoice.invoiceNumber}</p><h2 id="payment-modal-title" className="mt-1 text-lg font-medium text-stripe-ink">Record payment</h2><p className="mt-1 text-xs text-stripe-ink-mute">{selectedInvoice.customer} · balance {currency(selectedInvoice.amountCents)}</p></div><button type="button" onClick={() => setSelectedInvoice(null)} aria-label="Close payment dialog" className="rounded p-1.5 text-stripe-ink-mute hover:bg-slate-100"><X className="h-4 w-4" /></button></div><form onSubmit={recordPayment} className="space-y-4 p-5"><label className="block text-xs font-medium text-stripe-ink-secondary">Payment method<select value={method} onChange={(event) => setMethod(event.target.value)} className="mt-1.5 w-full rounded-md border border-stripe-hairline-input bg-white px-3 py-2.5 text-sm text-stripe-ink"><option>Card</option><option>Bank transfer</option><option>Direct debit</option><option>Cash</option></select></label><div className="grid gap-4 sm:grid-cols-2"><label className="block text-xs font-medium text-stripe-ink-secondary">Amount<input required type="number" min="0.01" max={(selectedInvoice.amountCents / 100).toFixed(2)} step="0.01" value={amount} onChange={(event) => setAmount(event.target.value)} className="mt-1.5 w-full rounded-md border border-stripe-hairline-input px-3 py-2.5 text-sm text-stripe-ink" /></label><label className="block text-xs font-medium text-stripe-ink-secondary">Payment date<input required type="date" value={paymentDate} onChange={(event) => setPaymentDate(event.target.value)} className="mt-1.5 w-full rounded-md border border-stripe-hairline-input px-3 py-2.5 text-sm text-stripe-ink" /></label></div><label className="block text-xs font-medium text-stripe-ink-secondary">Note <span className="font-normal text-stripe-ink-mute">(optional)</span><textarea rows={3} value={note} onChange={(event) => setNote(event.target.value)} placeholder="Reference or follow-up note" className="mt-1.5 w-full resize-y rounded-md border border-stripe-hairline-input px-3 py-2.5 text-sm text-stripe-ink placeholder:text-slate-400" /></label><div className="flex justify-end gap-2 border-t border-stripe-hairline pt-4"><button type="button" onClick={() => setSelectedInvoice(null)} className="rounded-full px-4 py-2 text-sm text-stripe-ink-secondary hover:bg-slate-50">Cancel</button><button type="submit" className="inline-flex items-center gap-2 rounded-full bg-stripe-primary px-4 py-2 text-sm font-medium text-white hover:bg-stripe-primary-deep"><Check className="h-4 w-4" /> Save payment</button></div></form></section></div>}
    </div>
  );
};

export default BillingOperationsPage;
