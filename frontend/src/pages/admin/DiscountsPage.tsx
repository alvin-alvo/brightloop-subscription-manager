import React, { useState } from 'react';
import {
  Plus,
  Percent,
  Copy,
  Calendar,
  X,
} from 'lucide-react';
import {
  mockAdminService,
  formatCents,
  type AdminDiscount,
} from '../../services/mockAdminData';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { Button } from '../../components/common/Button';
import { InputField } from '../../components/common/InputField';

export const DiscountsPage: React.FC = () => {
  const { addToast } = useToast();

  const [discounts, setDiscounts] = useState<AdminDiscount[]>(() =>
    mockAdminService.getDiscounts()
  );
  const products = mockAdminService.getProducts();

  const [showModal, setShowModal] = useState(false);

  // Form State
  const [form, setForm] = useState({
    code: '',
    name: '',
    type: 'percentage' as AdminDiscount['type'],
    value: '',
    minPurchaseAmount: '',
    minQuantity: 1,
    startDate: new Date().toISOString().split('T')[0],
    endDate: '2025-12-31',
    usageLimit: 100,
    target: 'subscriptions' as AdminDiscount['target'],
    selectedProductIds: [] as string[],
  });

  const handleToggleProductSelection = (prodId: string) => {
    setForm((prev) => {
      const exists = prev.selectedProductIds.includes(prodId);
      return {
        ...prev,
        selectedProductIds: exists
          ? prev.selectedProductIds.filter((id) => id !== prodId)
          : [...prev.selectedProductIds, prodId],
      };
    });
  };

  const handleCreateDiscount = (e: React.FormEvent) => {
    e.preventDefault();

    if (!form.code.trim() || !form.name.trim() || !form.value) {
      addToast('Validation', 'Please provide a discount code, name, and value.', 'warning');
      return;
    }

    if (form.target === 'products' && form.selectedProductIds.length === 0) {
      addToast('Validation', 'Please select at least one product for this discount.', 'warning');
      return;
    }

    const numericVal = parseFloat(form.value);
    const value = form.type === 'fixed_amount' ? Math.round(numericVal * 100) : numericVal;
    const minPurchaseAmountCents = form.minPurchaseAmount
      ? Math.round(parseFloat(form.minPurchaseAmount) * 100)
      : 0;

    const created = mockAdminService.addDiscount({
      code: form.code.trim().toUpperCase(),
      name: form.name.trim(),
      type: form.type,
      value: value,
      minPurchaseAmountCents,
      minQuantity: form.minQuantity || 1,
      startDate: form.startDate,
      endDate: form.endDate,
      usageLimit: form.usageLimit || 100,
      target: form.target,
      targetProductIds: form.target === 'products' ? form.selectedProductIds : undefined,
      active: true,
    });

    setDiscounts(mockAdminService.getDiscounts());
    setShowModal(false);

    // Reset Form
    setForm({
      code: '',
      name: '',
      type: 'percentage',
      value: '',
      minPurchaseAmount: '',
      minQuantity: 1,
      startDate: new Date().toISOString().split('T')[0],
      endDate: '2025-12-31',
      usageLimit: 100,
      target: 'subscriptions',
      selectedProductIds: [],
    });

    addToast(
      'Discount Rule Created',
      `Promo code ${created.code} successfully registered for ${
        created.target === 'subscriptions' ? 'whole subscriptions' : 'specific products'
      }.`,
      'success'
    );
  };

  const handleToggleStatus = (id: string) => {
    const updated = mockAdminService.toggleDiscountStatus(id);
    setDiscounts(mockAdminService.getDiscounts());
    addToast(
      'Discount Status Updated',
      `Code ${updated.code} is now ${updated.active ? 'active' : 'paused'}.`,
      'info'
    );
  };

  const handleCopyCode = (code: string) => {
    navigator.clipboard?.writeText(code);
    addToast('Coupon Copied', `Promo code "${code}" copied to clipboard.`, 'info');
  };

  return (
    <div className="space-y-6 animate-fadeIn">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
              Discount & Promotion Engine
            </h1>
            <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-indigo-50 text-stripe-primary border border-indigo-200 uppercase">
              Admin Exclusive
            </span>
          </div>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Author discount rules, govern subscription coupons, and control product-level price reductions.
          </p>
        </div>

        <Button
          type="button"
          variant="primary"
          size="md"
          onClick={() => setShowModal(true)}
          leftIcon={<Plus className="w-4 h-4" />}
        >
          Create Discount Rule
        </Button>
      </div>

      {/* Summary Metrics */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal text-stripe-ink-mute uppercase tracking-wider">
            Active Rules
          </span>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">
            {discounts.filter((d) => d.active).length}
          </div>
          <span className="text-xs text-stripe-ink-mute">Configured in system</span>
        </div>

        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal text-stripe-ink-mute uppercase tracking-wider">
            Total Redemptions
          </span>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">
            {discounts.reduce((acc, d) => acc + d.redemptions, 0)}
          </div>
          <span className="text-xs text-emerald-600 font-medium">Applied to billing cycles</span>
        </div>

        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <span className="text-xs font-normal text-stripe-ink-mute uppercase tracking-wider">
            Scope Distribution
          </span>
          <div className="text-sm font-medium text-stripe-ink mt-2 flex items-center gap-3">
            <span className="flex items-center gap-1.5">
              <span className="w-2 h-2 rounded-full bg-stripe-primary"></span>
              {discounts.filter((d) => d.target === 'subscriptions').length} Subscriptions
            </span>
            <span className="flex items-center gap-1.5">
              <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
              {discounts.filter((d) => d.target === 'products').length} Products
            </span>
          </div>
        </div>
      </div>

      {/* Discounts Table */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Promo Code</th>
                <th className="py-3 px-5">Discount Value</th>
                <th className="py-3 px-5">Target Scope</th>
                <th className="py-3 px-5">Eligibility Constraints</th>
                <th className="py-3 px-5">Validity Period</th>
                <th className="py-3 px-5">Redemptions</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {discounts.map((d) => {
                const isDepleted = d.redemptions >= d.usageLimit;
                return (
                  <tr key={d.id} className="hover:bg-slate-50/70 transition-colors">
                    <td className="py-3.5 px-5">
                      <div className="flex items-center gap-1.5 font-mono text-xs font-semibold text-stripe-primary">
                        <span>{d.code}</span>
                        <button
                          onClick={() => handleCopyCode(d.code)}
                          className="p-1 text-slate-400 hover:text-stripe-ink rounded transition-colors"
                          title="Copy promo code"
                        >
                          <Copy className="w-3.5 h-3.5" />
                        </button>
                      </div>
                      <div className="text-[11px] text-stripe-ink-mute mt-0.5">{d.name}</div>
                    </td>

                    <td className="py-3.5 px-5 font-medium text-stripe-ink tnum">
                      {d.type === 'percentage'
                        ? `${d.value}% Off`
                        : `${formatCents(d.value)} Off`}
                    </td>

                    <td className="py-3.5 px-5">
                      {d.target === 'subscriptions' ? (
                        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs bg-indigo-50 text-stripe-primary border border-indigo-100">
                          Whole Subscriptions
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs bg-emerald-50 text-emerald-700 border border-emerald-200">
                          Specific Products ({d.targetProductIds?.length || 1})
                        </span>
                      )}
                    </td>

                    <td className="py-3.5 px-5 text-xs text-stripe-ink-secondary">
                      <div>Min: {formatCents(d.minPurchaseAmountCents)}</div>
                      <div className="text-[11px] text-stripe-ink-mute">
                        Min Qty: {d.minQuantity} {d.minQuantity === 1 ? 'seat' : 'seats'}
                      </div>
                    </td>

                    <td className="py-3.5 px-5 text-xs text-stripe-ink-mute tnum">
                      <div>{d.startDate}</div>
                      <div className="text-[11px] text-slate-400">until {d.endDate}</div>
                    </td>

                    <td className="py-3.5 px-5 text-xs font-normal text-stripe-ink tnum">
                      <div className="flex items-center gap-2">
                        <span>
                          {d.redemptions} / {d.usageLimit}
                        </span>
                        <div className="w-16 h-1.5 bg-slate-100 rounded-full overflow-hidden hidden sm:block">
                          <div
                            className={`h-full ${
                              isDepleted ? 'bg-stripe-ruby' : 'bg-stripe-primary'
                            }`}
                            style={{
                              width: `${Math.min(
                                100,
                                Math.round((d.redemptions / d.usageLimit) * 100)
                              )}%`,
                            }}
                          />
                        </div>
                      </div>
                    </td>

                    <td className="py-3.5 px-5">
                      <StatusBadge
                        label={isDepleted ? 'Depleted' : d.active ? 'Active' : 'Paused'}
                        variant={isDepleted ? 'neutral' : d.active ? 'success' : 'warning'}
                      />
                    </td>

                    <td className="py-3.5 px-5 text-right">
                      <button
                        onClick={() => handleToggleStatus(d.id)}
                        className="text-xs text-stripe-primary hover:text-stripe-primary-deep font-medium"
                      >
                        {d.active ? 'Pause Rule' : 'Resume'}
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* --- CREATE DISCOUNT MODAL --- */}
      {showModal && (
        <div className="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white border border-stripe-hairline rounded-2xl shadow-stripe-lg max-w-lg w-full max-h-[90vh] overflow-y-auto p-6 animate-fadeIn">
            <div className="flex items-center justify-between pb-3 border-b border-stripe-hairline">
              <div className="flex items-center gap-2">
                <Percent className="w-5 h-5 text-stripe-primary" />
                <h2 className="text-lg font-light text-stripe-ink">Author Discount Rule</h2>
              </div>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 text-slate-400 hover:text-stripe-ink rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreateDiscount} className="space-y-4 pt-4">
              <div className="grid grid-cols-2 gap-3">
                <InputField
                  label="Promo Code"
                  value={form.code}
                  onChange={(e) => setForm({ ...form, code: e.target.value.toUpperCase() })}
                  placeholder="e.g., WINTER2026"
                  required
                />

                <InputField
                  label="Rule Name"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="e.g., Q1 Enterprise Promo"
                  required
                />
              </div>

              {/* Type and Value */}
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Discount Type
                  </label>
                  <select
                    value={form.type}
                    onChange={(e) =>
                      setForm({ ...form, type: e.target.value as AdminDiscount['type'] })
                    }
                    className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  >
                    <option value="percentage">Percentage (%)</option>
                    <option value="fixed_amount">Fixed Amount ($)</option>
                  </select>
                </div>

                <InputField
                  label={form.type === 'percentage' ? 'Percentage (e.g. 20)' : 'Amount in USD ($)'}
                  type="number"
                  step={form.type === 'percentage' ? '1' : '0.01'}
                  value={form.value}
                  onChange={(e) => setForm({ ...form, value: e.target.value })}
                  placeholder={form.type === 'percentage' ? '20' : '50.00'}
                  required
                />
              </div>

              {/* Target Selector: Products or Whole Subscriptions */}
              <div className="p-3.5 bg-slate-50 rounded-xl border border-stripe-hairline space-y-2.5">
                <label className="block text-xs font-semibold uppercase tracking-wider text-stripe-ink">
                  Application Target Selector
                </label>
                <div className="grid grid-cols-2 gap-2">
                  <label
                    className={`flex items-center gap-2 p-2.5 rounded-lg border text-xs cursor-pointer transition-all ${
                      form.target === 'subscriptions'
                        ? 'border-stripe-primary bg-white text-stripe-ink shadow-sm'
                        : 'border-slate-200 text-slate-600 bg-white/60'
                    }`}
                  >
                    <input
                      type="radio"
                      name="target"
                      value="subscriptions"
                      checked={form.target === 'subscriptions'}
                      onChange={() => setForm({ ...form, target: 'subscriptions' })}
                      className="text-stripe-primary focus:ring-stripe-primary"
                    />
                    <span className="font-medium">Whole Subscriptions</span>
                  </label>

                  <label
                    className={`flex items-center gap-2 p-2.5 rounded-lg border text-xs cursor-pointer transition-all ${
                      form.target === 'products'
                        ? 'border-stripe-primary bg-white text-stripe-ink shadow-sm'
                        : 'border-slate-200 text-slate-600 bg-white/60'
                    }`}
                  >
                    <input
                      type="radio"
                      name="target"
                      value="products"
                      checked={form.target === 'products'}
                      onChange={() => setForm({ ...form, target: 'products' })}
                      className="text-stripe-primary focus:ring-stripe-primary"
                    />
                    <span className="font-medium">Specific Products</span>
                  </label>
                </div>

                {/* Specific Products Picker */}
                {form.target === 'products' && (
                  <div className="pt-2 border-t border-slate-200 space-y-1.5 animate-fadeIn">
                    <span className="text-[11px] text-stripe-ink-mute block">
                      Choose which products this discount applies to:
                    </span>
                    <div className="space-y-1 max-h-32 overflow-y-auto">
                      {products.map((p) => {
                        const isChecked = form.selectedProductIds.includes(p.id);
                        return (
                          <label
                            key={p.id}
                            className="flex items-center justify-between p-2 rounded bg-white border border-slate-200 text-xs cursor-pointer hover:bg-slate-50"
                          >
                            <span className="font-medium text-stripe-ink">{p.name}</span>
                            <input
                              type="checkbox"
                              checked={isChecked}
                              onChange={() => handleToggleProductSelection(p.id)}
                              className="w-4 h-4 text-stripe-primary rounded border-slate-300 focus:ring-stripe-primary"
                            />
                          </label>
                        );
                      })}
                    </div>
                  </div>
                )}
              </div>

              {/* Min Purchase & Quantity */}
              <div className="grid grid-cols-2 gap-3">
                <InputField
                  label="Min Purchase Amount ($)"
                  type="number"
                  step="0.01"
                  value={form.minPurchaseAmount}
                  onChange={(e) => setForm({ ...form, minPurchaseAmount: e.target.value })}
                  placeholder="0.00"
                />

                <InputField
                  label="Min Quantity (Seats)"
                  type="number"
                  min="1"
                  value={form.minQuantity}
                  onChange={(e) =>
                    setForm({ ...form, minQuantity: parseInt(e.target.value) || 1 })
                  }
                  required
                />
              </div>

              {/* Start & End Dates */}
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={form.startDate}
                    onChange={(e) => setForm({ ...form, startDate: e.target.value })}
                    className="w-full bg-white text-stripe-ink text-xs rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    End Date
                  </label>
                  <div className="relative">
                    <Calendar className="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                    <input
                      type="date"
                      value={form.endDate}
                      onChange={(e) => setForm({ ...form, endDate: e.target.value })}
                      className="w-full bg-white text-stripe-ink text-xs rounded-md pl-9 pr-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                    />
                  </div>
                </div>
              </div>

              {/* Usage Limit */}
              <InputField
                label="Maximum Usage Limit (Redemptions)"
                type="number"
                min="1"
                value={form.usageLimit}
                onChange={(e) =>
                  setForm({ ...form, usageLimit: parseInt(e.target.value) || 1 })
                }
                helperText="Rule automatically deactivates once redemption limit is met."
                required
              />

              <div className="pt-2 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="md"
                  onClick={() => setShowModal(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" variant="primary" size="md">
                  Save Discount Rule
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default DiscountsPage;
