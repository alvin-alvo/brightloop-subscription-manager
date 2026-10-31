import React from 'react';
import { Search, Building2, Mail, ExternalLink, Plus } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const CrmCustomersPage: React.FC = () => {
  const { addToast } = useToast();

  const customers = [
    {
      id: 'cust_001',
      name: 'Acme Enterprises Inc.',
      contact: 'David Chen',
      email: 'david@acme-enterprises.com',
      mrr: '$1,200.00',
      activeSubscriptions: 2,
      tier: 'Enterprise Annual',
      status: 'active' as const,
    },
    {
      id: 'cust_002',
      name: 'Vertex Technologies',
      contact: 'Elena Rostova',
      email: 'elena@vertex.tech',
      mrr: '$499.00',
      activeSubscriptions: 1,
      tier: 'Growth Tier',
      status: 'active' as const,
    },
    {
      id: 'cust_003',
      name: 'Quantum AI Systems',
      contact: 'Liam O’Connor',
      email: 'liam@quantumai.dev',
      mrr: '$890.00',
      activeSubscriptions: 3,
      tier: 'Scale Tier',
      status: 'active' as const,
    },
    {
      id: 'cust_004',
      name: 'Helix Health Labs',
      contact: 'Dr. Maya Lin',
      email: 'clinics@helixhealth.io',
      mrr: '$299.00',
      activeSubscriptions: 1,
      tier: 'Pro Monthly',
      status: 'warning' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Customer Relationship Management
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Accounts, contracted contacts, active billing profiles, and lifetime values.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Customer Created', 'New corporate account onboarded.', 'success')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          Onboard Customer
        </button>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        {/* Search & Filter Toolbar */}
        <div className="p-4 border-b border-stripe-hairline flex items-center justify-between gap-3">
          <div className="relative w-72">
            <Search className="w-4 h-4 text-stripe-ink-mute absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              placeholder="Search companies or emails..."
              className="w-full pl-9 pr-3 py-1.5 text-xs bg-slate-50 border border-stripe-hairline-input rounded-md focus:bg-white focus:outline-none focus:border-stripe-primary"
            />
          </div>
          <span className="text-xs text-stripe-ink-mute tnum">4 total accounts shown</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Customer Organization</th>
                <th className="py-3 px-5">Primary Contact</th>
                <th className="py-3 px-5">Active Tier</th>
                <th className="py-3 px-5">Monthly Value</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Details</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {customers.map((c) => (
                <tr key={c.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5">
                    <div className="flex items-center gap-2.5">
                      <div className="p-1.5 rounded-md bg-slate-100 text-stripe-ink-secondary">
                        <Building2 className="w-4 h-4" />
                      </div>
                      <div>
                        <div className="font-medium text-stripe-ink">{c.name}</div>
                        <div className="text-[11px] font-mono text-stripe-ink-mute">
                          {c.id}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 px-5">
                    <div className="text-stripe-ink text-xs">{c.contact}</div>
                    <div className="text-[11px] text-stripe-ink-mute flex items-center gap-1">
                      <Mail className="w-3 h-3" />
                      {c.email}
                    </div>
                  </td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-secondary">
                    {c.tier}
                  </td>
                  <td className="py-3 px-5 font-normal text-stripe-ink tnum">
                    {c.mrr}
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge
                      label={c.status === 'active' ? 'Good Standing' : 'Past Due Notice'}
                      variant={c.status === 'active' ? 'success' : 'warning'}
                    />
                  </td>
                  <td className="py-3 px-5 text-right">
                    <button
                      onClick={() =>
                        addToast('Customer View', `Opening ledger for ${c.name}`, 'info')
                      }
                      className="p-1 rounded text-slate-400 hover:text-stripe-primary hover:bg-slate-50 transition-colors"
                      title="Open Profile"
                    >
                      <ExternalLink className="w-4 h-4" />
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

export default CrmCustomersPage;

