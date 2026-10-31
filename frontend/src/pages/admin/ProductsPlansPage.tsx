import React, { useState } from 'react';
import {
  Plus,
  Layers,
  Repeat,
  Trash2,
  Calendar,
  Sparkles,
  X,
} from 'lucide-react';
import {
  mockAdminService,
  formatCents,
  type AdminProduct,
  type AdminPlan,
  type ProductVariant,
} from '../../services/mockAdminData';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { Button } from '../../components/common/Button';
import { InputField } from '../../components/common/InputField';

export const ProductsPlansPage: React.FC = () => {
  const { addToast } = useToast();

  const [products, setProducts] = useState<AdminProduct[]>(() =>
    mockAdminService.getProducts()
  );
  const [plans, setPlans] = useState<AdminPlan[]>(() =>
    mockAdminService.getPlans()
  );

  const [activeTab, setActiveTab] = useState<'catalog' | 'plans'>('catalog');
  const [showProductModal, setShowProductModal] = useState(false);
  const [showPlanModal, setShowPlanModal] = useState(false);

  // --- NEW PRODUCT FORM STATE ---
  const [prodForm, setProdForm] = useState({
    name: '',
    type: 'Software' as AdminProduct['type'],
    salesPrice: '',
    costPrice: '',
    notes: '',
  });

  // Dynamic variants list for the new product
  const [variants, setVariants] = useState<ProductVariant[]>([
    {
      id: `var_${Date.now()}`,
      name: 'Odoo Enterprise Connector',
      brand: 'Odoo',
      priceAdjustmentCents: 56000, // +$560.00
    },
  ]);

  // Temporary row for adding another variant
  const [newVar, setNewVar] = useState({
    name: '',
    brand: '',
    priceAdjustment: '',
  });

  // --- NEW PLAN FORM STATE ---
  const [planForm, setPlanForm] = useState({
    productId: products[0]?.id || '',
    name: '',
    price: '',
    cost: '',
    billingFrequency: 'monthly' as AdminPlan['billingFrequency'],
    minQuantity: 1,
    validFrom: new Date().toISOString().split('T')[0],
    validUntil: '2025-12-31',
    autoCloseAtEnd: false,
    allowEarlyClose: true,
    allowPause: true,
    allowRenew: true,
  });

  // Add a variant dynamically
  const handleAddVariant = () => {
    if (!newVar.name.trim() || !newVar.priceAdjustment) {
      addToast('Validation', 'Please provide a variant name and price adjustment.', 'warning');
      return;
    }
    const cents = Math.round(parseFloat(newVar.priceAdjustment) * 100);
    setVariants((prev) => [
      ...prev,
      {
        id: `var_${Date.now()}_${Math.random().toString(36).substring(2, 5)}`,
        name: newVar.name.trim(),
        brand: newVar.brand.trim() || 'Standard',
        priceAdjustmentCents: cents,
      },
    ]);
    setNewVar({ name: '', brand: '', priceAdjustment: '' });
  };

  const handleRemoveVariant = (id: string) => {
    setVariants((prev) => prev.filter((v) => v.id !== id));
  };

  // Submit New Product
  const handleCreateProduct = (e: React.FormEvent) => {
    e.preventDefault();
    if (!prodForm.name.trim() || !prodForm.salesPrice) {
      addToast('Validation', 'Product name and sales price are required.', 'warning');
      return;
    }

    const salesPriceCents = Math.round(parseFloat(prodForm.salesPrice) * 100);
    const costPriceCents = prodForm.costPrice ? Math.round(parseFloat(prodForm.costPrice) * 100) : 0;

    const created = mockAdminService.addProduct({
      name: prodForm.name.trim(),
      type: prodForm.type,
      salesPriceCents,
      costPriceCents,
      notes: prodForm.notes.trim(),
      variants,
      active: true,
    });

    setProducts(mockAdminService.getProducts());
    setShowProductModal(false);
    setProdForm({ name: '', type: 'Software', salesPrice: '', costPrice: '', notes: '' });
    setVariants([]);

    addToast(
      'Product Configured',
      `"${created.name}" created with ${created.variants.length} dynamic variants.`,
      'success'
    );
  };

  // Submit New Plan
  const handleCreatePlan = (e: React.FormEvent) => {
    e.preventDefault();
    if (!planForm.name.trim() || !planForm.price || !planForm.productId) {
      addToast('Validation', 'Please specify plan name, linked product, and price.', 'warning');
      return;
    }

    const selectedProduct = products.find((p) => p.id === planForm.productId);
    const priceCents = Math.round(parseFloat(planForm.price) * 100);
    const costCents = planForm.cost ? Math.round(parseFloat(planForm.cost) * 100) : 0;

    const created = mockAdminService.addPlan({
      productId: planForm.productId,
      productName: selectedProduct?.name || 'Selected Product',
      name: planForm.name.trim(),
      priceCents,
      costCents,
      billingFrequency: planForm.billingFrequency,
      minQuantity: planForm.minQuantity || 1,
      validFrom: planForm.validFrom,
      validUntil: planForm.validUntil,
      behaviors: {
        autoCloseAtEnd: planForm.autoCloseAtEnd,
        allowEarlyClose: planForm.allowEarlyClose,
        allowPause: planForm.allowPause,
        allowRenew: planForm.allowRenew,
      },
      active: true,
    });

    setPlans(mockAdminService.getPlans());
    setShowPlanModal(false);
    setPlanForm({
      productId: products[0]?.id || '',
      name: '',
      price: '',
      cost: '',
      billingFrequency: 'monthly',
      minQuantity: 1,
      validFrom: new Date().toISOString().split('T')[0],
      validUntil: '2025-12-31',
      autoCloseAtEnd: false,
      allowEarlyClose: true,
      allowPause: true,
      allowRenew: true,
    });

    addToast(
      'Plan Created',
      `"${created.name}" linked to ${created.productName} (${created.billingFrequency}).`,
      'success'
    );
  };

  const handleTogglePlan = (id: string) => {
    const updated = mockAdminService.togglePlanStatus(id);
    setPlans(mockAdminService.getPlans());
    addToast(
      'Plan Status Updated',
      `"${updated.name}" is now ${updated.active ? 'active' : 'paused'}.`,
      'info'
    );
  };

  return (
    <div className="space-y-6 animate-fadeIn">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Product & Plan Configurator
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Author SaaS products, price dynamic brand variants, and configure recurring billing plans.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Button
            type="button"
            variant="secondary"
            size="sm"
            onClick={() => setShowProductModal(true)}
            leftIcon={<Plus className="w-3.5 h-3.5" />}
          >
            New Product
          </Button>

          <Button
            type="button"
            variant="primary"
            size="sm"
            onClick={() => setShowPlanModal(true)}
            leftIcon={<Plus className="w-3.5 h-3.5" />}
          >
            Create Plan
          </Button>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="border-b border-stripe-hairline flex items-center gap-6">
        <button
          onClick={() => setActiveTab('catalog')}
          className={`pb-3 text-sm font-medium transition-colors relative flex items-center gap-2 ${
            activeTab === 'catalog'
              ? 'text-stripe-primary font-semibold'
              : 'text-stripe-ink-secondary hover:text-stripe-ink'
          }`}
        >
          <Layers className="w-4 h-4" />
          <span>Products & Brand Variants ({products.length})</span>
          {activeTab === 'catalog' && (
            <span className="absolute bottom-0 inset-x-0 h-0.5 bg-stripe-primary" />
          )}
        </button>

        <button
          onClick={() => setActiveTab('plans')}
          className={`pb-3 text-sm font-medium transition-colors relative flex items-center gap-2 ${
            activeTab === 'plans'
              ? 'text-stripe-primary font-semibold'
              : 'text-stripe-ink-secondary hover:text-stripe-ink'
          }`}
        >
          <Repeat className="w-4 h-4" />
          <span>Configured Plans ({plans.length})</span>
          {activeTab === 'plans' && (
            <span className="absolute bottom-0 inset-x-0 h-0.5 bg-stripe-primary" />
          )}
        </button>
      </div>

      {/* --- TAB 1: PRODUCTS & VARIANTS CATALOG --- */}
      {activeTab === 'catalog' && (
        <div className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
            {products.map((product) => (
              <div
                key={product.id}
                className="bg-white border border-stripe-hairline rounded-xl p-5 shadow-stripe-sm hover:shadow-stripe-md transition-shadow flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="flex items-center gap-2">
                        <h2 className="text-base font-normal text-stripe-ink">
                          {product.name}
                        </h2>
                        <span className="px-2 py-0.5 rounded text-[10px] font-medium bg-slate-100 text-stripe-ink-secondary border border-slate-200 uppercase">
                          {product.type}
                        </span>
                      </div>
                      <p className="text-xs text-stripe-ink-mute mt-1 leading-relaxed">
                        {product.notes || 'No description provided.'}
                      </p>
                    </div>

                    <div className="text-right shrink-0">
                      <span className="text-lg font-light text-stripe-ink tnum block">
                        {formatCents(product.salesPriceCents)}
                      </span>
                      <span className="text-[11px] text-stripe-ink-mute tnum">
                        Cost: {formatCents(product.costPriceCents)}
                      </span>
                    </div>
                  </div>

                  {/* Dynamic Variants Section */}
                  <div className="mt-4 pt-3.5 border-t border-slate-100">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-[11px] font-semibold text-stripe-ink-mute uppercase tracking-wider flex items-center gap-1">
                        <Sparkles className="w-3 h-3 text-stripe-primary" />
                        Dynamic Brand Variants ({product.variants.length})
                      </span>
                    </div>

                    {product.variants.length === 0 ? (
                      <p className="text-xs text-slate-400 italic">No variants configured.</p>
                    ) : (
                      <div className="space-y-2">
                        {product.variants.map((v) => (
                          <div
                            key={v.id}
                            className="flex items-center justify-between p-2.5 bg-slate-50/80 rounded-lg border border-slate-200/60 text-xs"
                          >
                            <div className="flex items-center gap-2">
                              <span className="px-1.5 py-0.5 rounded text-[10px] font-semibold bg-indigo-50 text-stripe-primary border border-indigo-100">
                                {v.brand}
                              </span>
                              <span className="font-medium text-stripe-ink">{v.name}</span>
                            </div>
                            <span className="font-medium text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200 tnum">
                              +{formatCents(v.priceAdjustmentCents)}
                            </span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                <div className="mt-4 pt-3 border-t border-stripe-hairline flex items-center justify-between text-xs text-stripe-ink-mute">
                  <span className="font-mono text-[11px]">{product.id}</span>
                  <button
                    onClick={() => {
                      mockAdminService.deleteProduct(product.id);
                      setProducts(mockAdminService.getProducts());
                      addToast('Product Deleted', `Removed ${product.name}`, 'info');
                    }}
                    className="text-rose-600 hover:text-rose-700 font-medium inline-flex items-center gap-1"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* --- TAB 2: LINKED PLANS TABLE --- */}
      {activeTab === 'plans' && (
        <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                  <th className="py-3 px-5">Plan Name</th>
                  <th className="py-3 px-5">Linked Product</th>
                  <th className="py-3 px-5">Frequency</th>
                  <th className="py-3 px-5">Price (Cost)</th>
                  <th className="py-3 px-5">Plan Behaviors</th>
                  <th className="py-3 px-5">Subscribers</th>
                  <th className="py-3 px-5 text-right">Status / Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-stripe-hairline text-sm">
                {plans.map((p) => (
                  <tr key={p.id} className="hover:bg-slate-50/70 transition-colors">
                    <td className="py-3.5 px-5 font-medium text-stripe-ink">
                      <div>{p.name}</div>
                      <div className="text-[11px] text-stripe-ink-mute font-mono">{p.id}</div>
                    </td>

                    <td className="py-3.5 px-5 text-xs text-stripe-ink-secondary">
                      {p.productName}
                    </td>

                    <td className="py-3.5 px-5">
                      <span className="capitalize px-2 py-0.5 rounded text-xs bg-slate-100 text-stripe-ink font-medium">
                        {p.billingFrequency}
                      </span>
                    </td>

                    <td className="py-3.5 px-5">
                      <div className="font-medium text-stripe-ink tnum">
                        {formatCents(p.priceCents)}
                      </div>
                      <div className="text-[11px] text-stripe-ink-mute tnum">
                        Cost: {formatCents(p.costCents)}
                      </div>
                    </td>

                    {/* Plan Behaviors Indicators */}
                    <td className="py-3.5 px-5">
                      <div className="flex flex-wrap gap-1 text-[10px]">
                        {p.behaviors.autoCloseAtEnd && (
                          <span className="px-1.5 py-0.5 bg-amber-50 text-amber-700 border border-amber-200 rounded">
                            Auto-Close
                          </span>
                        )}
                        {p.behaviors.allowEarlyClose && (
                          <span className="px-1.5 py-0.5 bg-slate-100 text-slate-700 rounded">
                            Early-Close
                          </span>
                        )}
                        {p.behaviors.allowPause && (
                          <span className="px-1.5 py-0.5 bg-indigo-50 text-stripe-primary rounded">
                            Pauseable
                          </span>
                        )}
                        {p.behaviors.allowRenew && (
                          <span className="px-1.5 py-0.5 bg-emerald-50 text-emerald-700 rounded">
                            Renewable
                          </span>
                        )}
                      </div>
                    </td>

                    <td className="py-3.5 px-5 text-xs font-normal text-stripe-ink tnum">
                      {p.subscribersCount} users
                    </td>

                    <td className="py-3.5 px-5 text-right">
                      <button
                        onClick={() => handleTogglePlan(p.id)}
                        className="inline-flex items-center gap-1.5 text-xs font-medium"
                      >
                        <StatusBadge
                          label={p.active ? 'Active' : 'Paused'}
                          variant={p.active ? 'success' : 'neutral'}
                        />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* --- MODAL 1: ADD PRODUCT WITH DYNAMIC VARIANTS --- */}
      {showProductModal && (
        <div className="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white border border-stripe-hairline rounded-2xl shadow-stripe-lg max-w-lg w-full max-h-[90vh] overflow-y-auto p-6 animate-fadeIn">
            <div className="flex items-center justify-between pb-3 border-b border-stripe-hairline">
              <h2 className="text-lg font-light text-stripe-ink">Add New Product</h2>
              <button
                onClick={() => setShowProductModal(false)}
                className="p-1 text-slate-400 hover:text-stripe-ink rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreateProduct} className="space-y-4 pt-4">
              <InputField
                label="Product Name"
                value={prodForm.name}
                onChange={(e) => setProdForm({ ...prodForm, name: e.target.value })}
                placeholder="e.g., Enterprise Billing Connector"
                required
              />

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Product Type
                  </label>
                  <select
                    value={prodForm.type}
                    onChange={(e) =>
                      setProdForm({ ...prodForm, type: e.target.value as AdminProduct['type'] })
                    }
                    className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  >
                    <option value="Software">Software</option>
                    <option value="Service">Service</option>
                    <option value="Add-on">Add-on</option>
                    <option value="Infrastructure">Infrastructure</option>
                  </select>
                </div>

                <InputField
                  label="Sales Price ($ USD)"
                  type="number"
                  step="0.01"
                  value={prodForm.salesPrice}
                  onChange={(e) => setProdForm({ ...prodForm, salesPrice: e.target.value })}
                  placeholder="299.00"
                  required
                />
              </div>

              <InputField
                label="Cost Price ($ USD)"
                type="number"
                step="0.01"
                value={prodForm.costPrice}
                onChange={(e) => setProdForm({ ...prodForm, costPrice: e.target.value })}
                placeholder="50.00"
              />

              <div>
                <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                  Notes & Internal Description
                </label>
                <textarea
                  rows={2}
                  value={prodForm.notes}
                  onChange={(e) => setProdForm({ ...prodForm, notes: e.target.value })}
                  placeholder="Technical specs, target segment, or operational requirements..."
                  className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                />
              </div>

              {/* Dynamic Variants Section */}
              <div className="p-3.5 bg-slate-50 rounded-xl border border-stripe-hairline space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-semibold uppercase tracking-wider text-stripe-ink">
                    Dynamic Brand Variants
                  </span>
                  <span className="text-[11px] text-stripe-ink-mute">
                    e.g., Odoo brand adds +560
                  </span>
                </div>

                {/* List of current variants */}
                <div className="space-y-1.5">
                  {variants.map((v) => (
                    <div
                      key={v.id}
                      className="flex items-center justify-between p-2 bg-white rounded border border-stripe-hairline text-xs"
                    >
                      <div>
                        <span className="font-semibold text-stripe-primary mr-1">[{v.brand}]</span>
                        <span>{v.name}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-emerald-600 font-medium">
                          +{formatCents(v.priceAdjustmentCents)}
                        </span>
                        <button
                          type="button"
                          onClick={() => handleRemoveVariant(v.id)}
                          className="text-slate-400 hover:text-rose-600 p-0.5"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Add Variant Row */}
                <div className="grid grid-cols-3 gap-2 pt-1 border-t border-slate-200">
                  <input
                    type="text"
                    value={newVar.brand}
                    onChange={(e) => setNewVar({ ...newVar, brand: e.target.value })}
                    placeholder="Brand (e.g. Odoo)"
                    className="text-xs px-2.5 py-1.5 bg-white border border-slate-300 rounded"
                  />
                  <input
                    type="text"
                    value={newVar.name}
                    onChange={(e) => setNewVar({ ...newVar, name: e.target.value })}
                    placeholder="Option name"
                    className="text-xs px-2.5 py-1.5 bg-white border border-slate-300 rounded"
                  />
                  <div className="flex gap-1.5">
                    <input
                      type="number"
                      step="1"
                      value={newVar.priceAdjustment}
                      onChange={(e) => setNewVar({ ...newVar, priceAdjustment: e.target.value })}
                      placeholder="+Price"
                      className="text-xs px-2 py-1.5 bg-white border border-slate-300 rounded w-full"
                    />
                    <button
                      type="button"
                      onClick={handleAddVariant}
                      className="px-2.5 py-1.5 text-xs bg-slate-800 text-white rounded hover:bg-black font-medium shrink-0"
                    >
                      Add
                    </button>
                  </div>
                </div>
              </div>

              <div className="pt-2 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="md"
                  onClick={() => setShowProductModal(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" variant="primary" size="md">
                  Save Product
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* --- MODAL 2: ADD PLAN LINKED TO PRODUCT --- */}
      {showPlanModal && (
        <div className="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white border border-stripe-hairline rounded-2xl shadow-stripe-lg max-w-lg w-full max-h-[90vh] overflow-y-auto p-6 animate-fadeIn">
            <div className="flex items-center justify-between pb-3 border-b border-stripe-hairline">
              <h2 className="text-lg font-light text-stripe-ink">Configure Linked Recurring Plan</h2>
              <button
                onClick={() => setShowPlanModal(false)}
                className="p-1 text-slate-400 hover:text-stripe-ink rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleCreatePlan} className="space-y-4 pt-4">
              {/* Linked Product Selector */}
              <div>
                <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                  Linked Product
                </label>
                <select
                  value={planForm.productId}
                  onChange={(e) => setPlanForm({ ...planForm, productId: e.target.value })}
                  className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  required
                >
                  {products.map((prod) => (
                    <option key={prod.id} value={prod.id}>
                      {prod.name} ({formatCents(prod.salesPriceCents)})
                    </option>
                  ))}
                </select>
              </div>

              <InputField
                label="Plan Name"
                value={planForm.name}
                onChange={(e) => setPlanForm({ ...planForm, name: e.target.value })}
                placeholder="e.g., Growth Tier Annual"
                required
              />

              <div className="grid grid-cols-2 gap-3">
                <InputField
                  label="Plan Price ($ USD)"
                  type="number"
                  step="0.01"
                  value={planForm.price}
                  onChange={(e) => setPlanForm({ ...planForm, price: e.target.value })}
                  placeholder="299.00"
                  required
                />

                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Billing Frequency
                  </label>
                  <select
                    value={planForm.billingFrequency}
                    onChange={(e) =>
                      setPlanForm({
                        ...planForm,
                        billingFrequency: e.target.value as AdminPlan['billingFrequency'],
                      })
                    }
                    className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary capitalize"
                  >
                    <option value="daily">Daily</option>
                    <option value="weekly">Weekly</option>
                    <option value="monthly">Monthly</option>
                    <option value="yearly">Yearly</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <InputField
                  label="Cost Price ($)"
                  type="number"
                  step="0.01"
                  value={planForm.cost}
                  onChange={(e) => setPlanForm({ ...planForm, cost: e.target.value })}
                  placeholder="50.00"
                />

                <InputField
                  label="Min Quantity"
                  type="number"
                  min="1"
                  value={planForm.minQuantity}
                  onChange={(e) =>
                    setPlanForm({ ...planForm, minQuantity: parseInt(e.target.value) || 1 })
                  }
                  required
                />

                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Valid From
                  </label>
                  <input
                    type="date"
                    value={planForm.validFrom}
                    onChange={(e) => setPlanForm({ ...planForm, validFrom: e.target.value })}
                    className="w-full bg-white text-stripe-ink text-xs rounded-md px-2.5 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                  Valid Until Date
                </label>
                <div className="relative">
                  <Calendar className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type="date"
                    value={planForm.validUntil}
                    onChange={(e) => setPlanForm({ ...planForm, validUntil: e.target.value })}
                    className="w-full bg-white text-stripe-ink text-xs rounded-md pl-9 pr-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  />
                </div>
              </div>

              {/* Toggle Switches for Plan Behaviors */}
              <div className="p-4 bg-slate-50 rounded-xl border border-stripe-hairline space-y-3">
                <span className="text-xs font-semibold uppercase tracking-wider text-stripe-ink block">
                  Plan Behavior Switches
                </span>

                {/* Auto-close at end */}
                <label className="flex items-center justify-between cursor-pointer">
                  <div>
                    <div className="text-xs font-medium text-stripe-ink">Auto-close at end</div>
                    <div className="text-[11px] text-stripe-ink-mute">
                      Automatically terminate contract when validity expires
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={planForm.autoCloseAtEnd}
                    onChange={(e) =>
                      setPlanForm({ ...planForm, autoCloseAtEnd: e.target.checked })
                    }
                    className="w-4 h-4 text-stripe-primary rounded border-slate-300 focus:ring-stripe-primary"
                  />
                </label>

                {/* Allow early close */}
                <label className="flex items-center justify-between cursor-pointer">
                  <div>
                    <div className="text-xs font-medium text-stripe-ink">Allow early close</div>
                    <div className="text-[11px] text-stripe-ink-mute">
                      Permit customer or staff to cancel prior to maturity
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={planForm.allowEarlyClose}
                    onChange={(e) =>
                      setPlanForm({ ...planForm, allowEarlyClose: e.target.checked })
                    }
                    className="w-4 h-4 text-stripe-primary rounded border-slate-300 focus:ring-stripe-primary"
                  />
                </label>

                {/* Allow pause */}
                <label className="flex items-center justify-between cursor-pointer">
                  <div>
                    <div className="text-xs font-medium text-stripe-ink">Allow pause</div>
                    <div className="text-[11px] text-stripe-ink-mute">
                      Allow subscription pausing without voiding the agreement
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={planForm.allowPause}
                    onChange={(e) => setPlanForm({ ...planForm, allowPause: e.target.checked })}
                    className="w-4 h-4 text-stripe-primary rounded border-slate-300 focus:ring-stripe-primary"
                  />
                </label>

                {/* Allow renew */}
                <label className="flex items-center justify-between cursor-pointer">
                  <div>
                    <div className="text-xs font-medium text-stripe-ink">Allow renew</div>
                    <div className="text-[11px] text-stripe-ink-mute">
                      Permit auto-renewal or customer contract extension
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={planForm.allowRenew}
                    onChange={(e) => setPlanForm({ ...planForm, allowRenew: e.target.checked })}
                    className="w-4 h-4 text-stripe-primary rounded border-slate-300 focus:ring-stripe-primary"
                  />
                </label>
              </div>

              <div className="pt-2 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="md"
                  onClick={() => setShowPlanModal(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" variant="primary" size="md">
                  Create Plan
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProductsPlansPage;
