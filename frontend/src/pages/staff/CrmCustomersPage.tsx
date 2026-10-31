import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Building2, ExternalLink, Mail, Plus, Search } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { fetchStaffWorkspace } from '../../utils/mockStaffData';

const currency = (cents: number): string =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100);

export const CrmCustomersPage: React.FC = () => {
  const { addToast } = useToast();
  const [search, setSearch] = useState('');
  const { data, isLoading } = useQuery({ queryKey: ['staff-workspace'], queryFn: fetchStaffWorkspace });
  const customers = (data?.customers ?? []).filter((customer) =>
    `${customer.name} ${customer.contact} ${customer.email}`.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs font-medium uppercase text-stripe-ink-mute">Accounts · {data?.customers.length ?? 0}</p>
          <h1 className="mt-1 text-2xl font-light text-stripe-ink">Customer relationships</h1>
          <p className="mt-1 text-sm text-stripe-ink-mute">Account owners, active plans, and monthly recurring value.</p>
        </div>
        <button onClick={() => addToast('Customer intake started', 'A new account profile is ready to complete.', 'info')} className="inline-flex items-center justify-center gap-2 rounded-full bg-stripe-primary px-4 py-2.5 text-sm font-medium text-white shadow-stripe-sm hover:bg-stripe-primary-deep">
          <Plus className="h-4 w-4" /> Add customer
        </button>
      </div>

      <section className="overflow-hidden rounded-lg border border-stripe-hairline bg-white shadow-stripe-sm">
        <div className="flex flex-col gap-3 border-b border-stripe-hairline p-4 sm:flex-row sm:items-center sm:justify-between">
          <label className="relative block w-full sm:max-w-sm">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-stripe-ink-mute" />
            <input value={search} onChange={(event) => setSearch(event.target.value)} type="search" placeholder="Search customers or contacts" className="w-full rounded-md border border-stripe-hairline-input bg-white py-2 pl-9 pr-3 text-sm text-stripe-ink outline-none focus:border-stripe-primary" />
          </label>
          <span className="text-xs text-stripe-ink-mute">{customers.length} accounts</span>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[760px] text-left text-sm">
            <thead className="bg-stripe-canvas-soft text-[11px] uppercase text-stripe-ink-mute"><tr><th className="px-5 py-3 font-medium">Customer</th><th className="px-5 py-3 font-medium">Primary contact</th><th className="px-5 py-3 font-medium">Current plan</th><th className="px-5 py-3 font-medium">Monthly value</th><th className="px-5 py-3 font-medium">Status</th><th className="px-5 py-3 text-right font-medium">Subscriptions</th></tr></thead>
            <tbody className="divide-y divide-stripe-hairline">
              {isLoading ? <tr><td colSpan={6} className="px-5 py-12 text-center text-sm text-stripe-ink-mute">Loading customer accounts...</td></tr> : customers.map((customer) => (
                <tr key={customer.id} className="hover:bg-stripe-canvas-soft/70">
                  <td className="px-5 py-4"><div className="flex items-center gap-3"><span className="flex h-9 w-9 items-center justify-center rounded-md bg-slate-100 text-stripe-ink-secondary"><Building2 className="h-4 w-4" /></span><span><span className="block font-medium text-stripe-ink">{customer.name}</span><span className="font-mono text-[11px] text-stripe-ink-mute">{customer.id}</span></span></div></td>
                  <td className="px-5 py-4"><span className="block text-stripe-ink">{customer.contact}</span><span className="mt-1 flex items-center gap-1 text-xs text-stripe-ink-mute"><Mail className="h-3 w-3" />{customer.email}</span></td>
                  <td className="px-5 py-4 text-stripe-ink-secondary">{customer.tier}</td>
                  <td className="px-5 py-4 font-medium tabular-nums text-stripe-ink">{currency(customer.monthlyValueCents)}</td>
                  <td className="px-5 py-4"><StatusBadge label={customer.status === 'active' ? 'Good standing' : 'Needs attention'} variant={customer.status === 'active' ? 'success' : 'warning'} /></td>
                  <td className="px-5 py-4 text-right"><span className="mr-3 tabular-nums text-stripe-ink-secondary">{customer.activeSubscriptions}</span><button type="button" title="Open customer profile" aria-label={`Open ${customer.name} profile`} onClick={() => addToast('Customer profile', `Opening the account view for ${customer.name}.`, 'info')} className="rounded p-1.5 text-stripe-ink-mute hover:bg-indigo-50 hover:text-stripe-primary"><ExternalLink className="h-4 w-4" /></button></td>
                </tr>
              ))}
              {!isLoading && customers.length === 0 && <tr><td colSpan={6} className="px-5 py-12 text-center text-sm text-stripe-ink-mute">No customers match that search.</td></tr>}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
};

export default CrmCustomersPage;
