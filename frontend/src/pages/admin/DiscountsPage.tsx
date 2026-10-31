import React from 'react';
import { Plus } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const DiscountsPage: React.FC = () => {
  const { addToast } = useToast();

  const discounts = [
    {
      code: 'ANNUAL20',
      type: 'Percentage',
      value: '20% Off',
      duration: 'Once (1st Year)',
      redemptions: 142,
      maxRedemptions: 500,
      status: 'active' as const,
    },
    {
      code: 'STARTUP50',
      type: 'Percentage',
      value: '50% Off',
      duration: '6 Months',
      redemptions: 89,
      maxRedemptions: 100,
      status: 'active' as const,
    },
    {
      code: 'WELCOME100',
      type: 'Fixed Amount',
      value: '$100.00 Off',
      duration: 'Once',
      redemptions: 320,
      maxRedemptions: 320,
      status: 'neutral' as const,
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Discounts & Promotion Coupons
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Manage promotional discount codes, coupon redemptions, and referral incentives.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Discount Created', 'New discount code coupon draft initialized.', 'success')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          Create Coupon
        </button>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Promo Code</th>
                <th className="py-3 px-5">Discount Value</th>
                <th className="py-3 px-5">Duration</th>
                <th className="py-3 px-5">Redemptions</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {discounts.map((d) => (
                <tr key={d.code} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5 font-mono text-xs font-semibold text-stripe-primary">
                    {d.code}
                  </td>
                  <td className="py-3 px-5 text-stripe-ink font-medium tnum">{d.value}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink-secondary">{d.duration}</td>
                  <td className="py-3 px-5 text-xs text-stripe-ink tnum">
                    {d.redemptions} / {d.maxRedemptions}
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge
                      label={d.status === 'active' ? 'Active' : 'Depleted'}
                      variant={d.status === 'active' ? 'success' : 'neutral'}
                    />
                  </td>
                  <td className="py-3 px-5 text-right">
                    <button
                      onClick={() => addToast('Coupon Copied', `Code ${d.code} copied to clipboard`, 'info')}
                      className="text-xs text-stripe-primary hover:underline"
                    >
                      Copy Link
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

export default DiscountsPage;
