import React, { useEffect, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, ArrowRight, Plus, Trash2 } from 'lucide-react';
import { Button } from '../../components/common/Button';
import { useToast } from '../../context/ToastContext';
import {
  fetchStaffWorkspace,
  type QuoteLineTemplate,
  type StaffSubscription,
  type StaffWorkspaceData,
} from '../../utils/mockStaffData';

interface EditableQuoteLine extends QuoteLineTemplate {
  rowId: string;
}

const currency = (cents: number): string =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100);

export const QuoteBuilderPage: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { addToast } = useToast();
  const { data, isLoading } = useQuery({
    queryKey: ['staff-workspace'],
    queryFn: fetchStaffWorkspace,
  });
  const [templateId, setTemplateId] = useState('');
  const [customerId, setCustomerId] = useState('');
  const [validUntil, setValidUntil] = useState('');
  const [lines, setLines] = useState<EditableQuoteLine[]>([]);
  const selectedTemplate = data?.templates.find((template) => template.id === templateId);

  useEffect(() => {
    if (data && !templateId) setTemplateId(data.templates[0]?.id ?? '');
    if (data && !customerId) setCustomerId(data.customers[0]?.id ?? '');
  }, [data, templateId, customerId]);

  useEffect(() => {
    if (!selectedTemplate) return;
    setLines(selectedTemplate.lines.map((line, index) => ({
      ...line,
      rowId: `${selectedTemplate.id}-${index}`,
    })));
    const expires = new Date();
    expires.setDate(expires.getDate() + selectedTemplate.validityDays);
    setValidUntil(expires.toISOString().slice(0, 10));
  }, [selectedTemplate]);

  const subtotalCents = lines.reduce((total, line) => total + line.quantity * line.unitPriceCents, 0);
  const taxCents = lines.reduce(
    (total, line) => total + Math.round(line.quantity * line.unitPriceCents * line.taxRateBasisPoints / 10000),
    0
  );

  const updateLine = (rowId: string, changes: Partial<EditableQuoteLine>) => {
    setLines((current) => current.map((line) => line.rowId === rowId ? { ...line, ...changes } : line));
  };

  const createQuote = (convertToSubscription: boolean) => {
    const customer = data?.customers.find((item) => item.id === customerId);
    if (!customer || !selectedTemplate || !lines.length || !validUntil) {
      addToast('Complete the quotation', 'Choose a customer, validity date, and at least one product line.', 'warning');
      return;
    }

    const now = Date.now();
    const quoteId = `QT-${new Date().getFullYear()}-${String(now).slice(-5)}`;
    const quote = {
      id: quoteId,
      customerId: customer.id,
      customer: customer.name,
      plan: selectedTemplate.plan,
      amountCents: subtotalCents + taxCents,
      validUntil,
      status: convertToSubscription ? 'Approved' as const : 'Pending' as const,
    };
    const subscription: StaffSubscription = {
      id: `sub_${now}`,
      customerId: customer.id,
      customer: customer.name,
      plan: selectedTemplate.plan,
      monthlyValueCents: subtotalCents + taxCents,
      stage: 'Confirmed',
      updatedAt: new Date().toISOString().slice(0, 10),
    };

    queryClient.setQueryData<StaffWorkspaceData>(['staff-workspace'], (current) => current && ({
      ...current,
      quotes: [quote, ...current.quotes],
      subscriptions: convertToSubscription ? [subscription, ...current.subscriptions] : current.subscriptions,
    }));
    addToast(
      convertToSubscription ? 'Subscription confirmed' : 'Quotation created',
      convertToSubscription
        ? `${customer.name} was added to the Confirmed pipeline stage.`
        : `Proposal for ${customer.name} is ready to share.`,
      'success'
    );
    navigate(convertToSubscription ? '/staff/subscriptions' : '/staff/quotations');
  };

  if (isLoading || !data) {
    return <p className="py-12 text-center text-sm text-stripe-ink-mute">Loading quote catalog...</p>;
  }

  return (
    <form onSubmit={(event) => event.preventDefault()} className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <button type="button" onClick={() => navigate('/staff/quotations')} className="mb-3 inline-flex items-center gap-2 text-sm text-stripe-ink-mute hover:text-stripe-primary">
            <ArrowLeft className="h-4 w-4" /> Back to quotations
          </button>
          <h1 className="text-2xl font-light text-stripe-ink">Build a quotation</h1>
          <p className="mt-1 text-sm text-stripe-ink-mute">Configure the offer, review the totals, then create the proposal.</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button type="button" variant="secondary" onClick={() => createQuote(false)}>Save quotation</Button>
          <Button type="button" onClick={() => createQuote(true)} rightIcon={<ArrowRight className="h-4 w-4" />}>Create subscription</Button>
        </div>
      </div>

      <section className="rounded-lg border border-stripe-hairline bg-white p-5 shadow-stripe-sm sm:p-6">
        <div className="grid gap-5 md:grid-cols-3">
          <label className="block text-xs font-medium text-stripe-ink-secondary">
            Select template
            <select value={templateId} onChange={(event) => setTemplateId(event.target.value)} className="mt-2 w-full rounded-md border border-stripe-hairline-input bg-white px-3 py-2.5 text-sm text-stripe-ink focus:border-stripe-primary focus:outline-none">
              {data.templates.map((template) => <option key={template.id} value={template.id}>{template.name}</option>)}
            </select>
          </label>
          <label className="block text-xs font-medium text-stripe-ink-secondary">
            Customer
            <select value={customerId} onChange={(event) => setCustomerId(event.target.value)} className="mt-2 w-full rounded-md border border-stripe-hairline-input bg-white px-3 py-2.5 text-sm text-stripe-ink focus:border-stripe-primary focus:outline-none">
              {data.customers.map((customer) => <option key={customer.id} value={customer.id}>{customer.name}</option>)}
            </select>
          </label>
          <label className="block text-xs font-medium text-stripe-ink-secondary">
            Valid until
            <input type="date" value={validUntil} onChange={(event) => setValidUntil(event.target.value)} className="mt-2 w-full rounded-md border border-stripe-hairline-input bg-white px-3 py-2.5 text-sm text-stripe-ink focus:border-stripe-primary focus:outline-none" />
          </label>
        </div>
        {selectedTemplate && <p className="mt-4 text-xs text-stripe-ink-mute">Template plan: <span className="font-medium text-stripe-ink">{selectedTemplate.plan}</span> · {selectedTemplate.validityDays}-day validity</p>}
      </section>

      <section className="overflow-hidden rounded-lg border border-stripe-hairline bg-white shadow-stripe-sm">
        <div className="flex items-center justify-between border-b border-stripe-hairline px-5 py-4">
          <div>
            <h2 className="text-sm font-medium text-stripe-ink">Product lines</h2>
            <p className="mt-1 text-xs text-stripe-ink-mute">Prices are in USD. Tax is calculated per line.</p>
          </div>
          <button type="button" onClick={() => setLines((current) => [...current, { rowId: `custom-${Date.now()}`, productId: 'prod-service', product: 'Additional services', quantity: 1, unitPriceCents: 0, taxRateBasisPoints: 850 }])} className="inline-flex items-center gap-2 text-sm font-medium text-stripe-primary hover:text-stripe-primary-deep">
            <Plus className="h-4 w-4" /> Add item
          </button>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[720px] text-left text-sm">
            <thead className="bg-stripe-canvas-soft text-[11px] uppercase text-stripe-ink-mute">
              <tr><th className="px-5 py-3 font-medium">Product</th><th className="px-3 py-3 font-medium">Qty</th><th className="px-3 py-3 font-medium">Unit price</th><th className="px-3 py-3 font-medium">Tax</th><th className="px-3 py-3 text-right font-medium">Total</th><th className="w-12 px-3 py-3" /></tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline">
              {lines.map((line) => {
                const lineSubtotal = line.quantity * line.unitPriceCents;
                const lineTotal = lineSubtotal + Math.round(lineSubtotal * line.taxRateBasisPoints / 10000);
                return <tr key={line.rowId}>
                  <td className="px-5 py-3"><input aria-label="Product name" value={line.product} onChange={(event) => updateLine(line.rowId, { product: event.target.value })} className="w-full min-w-48 rounded border border-transparent bg-transparent px-2 py-1 text-sm text-stripe-ink hover:border-stripe-hairline-input focus:border-stripe-primary focus:outline-none" /></td>
                  <td className="px-3 py-3"><input aria-label={`Quantity for ${line.product}`} type="number" min="1" value={line.quantity} onChange={(event) => updateLine(line.rowId, { quantity: Math.max(1, Number(event.target.value)) })} className="w-16 rounded-md border border-stripe-hairline-input px-2 py-1.5 text-sm" /></td>
                  <td className="px-3 py-3"><div className="flex items-center rounded-md border border-stripe-hairline-input px-2"><span className="text-xs text-stripe-ink-mute">$</span><input aria-label={`Unit price for ${line.product}`} type="number" min="0" step="0.01" value={(line.unitPriceCents / 100).toFixed(2)} onChange={(event) => updateLine(line.rowId, { unitPriceCents: Math.round(Math.max(0, Number(event.target.value)) * 100) })} className="w-24 px-1 py-1.5 text-sm outline-none" /></div></td>
                  <td className="px-3 py-3"><select aria-label={`Tax for ${line.product}`} value={line.taxRateBasisPoints} onChange={(event) => updateLine(line.rowId, { taxRateBasisPoints: Number(event.target.value) })} className="rounded-md border border-stripe-hairline-input bg-white px-2 py-1.5 text-sm"><option value={0}>0%</option><option value={850}>8.5%</option><option value={1000}>10%</option></select></td>
                  <td className="px-3 py-3 text-right font-medium tabular-nums text-stripe-ink">{currency(lineTotal)}</td>
                  <td className="px-3 py-3"><button type="button" onClick={() => setLines((current) => current.filter((item) => item.rowId !== line.rowId))} aria-label={`Remove ${line.product}`} className="rounded p-1.5 text-stripe-ink-mute hover:bg-rose-50 hover:text-rose-700"><Trash2 className="h-4 w-4" /></button></td>
                </tr>;
              })}
            </tbody>
          </table>
        </div>
        <div className="ml-auto max-w-sm space-y-2 border-t border-stripe-hairline px-5 py-4 text-sm">
          <div className="flex justify-between text-stripe-ink-secondary"><span>Subtotal</span><span className="tabular-nums">{currency(subtotalCents)}</span></div>
          <div className="flex justify-between text-stripe-ink-secondary"><span>Tax</span><span className="tabular-nums">{currency(taxCents)}</span></div>
          <div className="flex justify-between border-t border-stripe-hairline pt-3 font-semibold text-stripe-ink"><span>Total due</span><span className="tabular-nums">{currency(subtotalCents + taxCents)}</span></div>
        </div>
      </section>
    </form>
  );
};

export default QuoteBuilderPage;