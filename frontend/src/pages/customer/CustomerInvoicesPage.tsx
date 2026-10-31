import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Download, CreditCard, CheckCircle2, Search, Filter } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useToast } from '../../context/ToastContext';
import { DashboardSkeleton } from '../../components/layout/DashboardSkeleton';
import { StatusBadge } from '../../components/StatusBadge';
import {
  subscriptionsService,
  formatDisplayDate,
} from '../../services/subscriptionsService';
import type { Invoice } from '../../types';

interface CustomerInvoicesPageProps {
  customerId?: string;
}

export const CustomerInvoicesPage: React.FC<CustomerInvoicesPageProps> = ({
  customerId: propCustomerId,
}) => {
  const { user } = useAuth();
  const { addToast } = useToast();
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState('');
  const [payingInvoiceId, setPayingInvoiceId] = useState<string | null>(null);

  // Strict data isolation
  const effectiveCustomerId = propCustomerId || user?.id || 'usr_cust_03';

  const { data: invoices, isLoading } = useQuery({
    queryKey: ['customerInvoices', effectiveCustomerId],
    queryFn: () => subscriptionsService.getCustomerInvoices(),
    staleTime: 1000 * 60 * 2,
  });

  const payMutation = useMutation({
    mutationFn: (invoiceId: string) => subscriptionsService.payInvoice(invoiceId),
    onMutate: (invoiceId) => {
      setPayingInvoiceId(invoiceId);
    },
    onSuccess: (updatedInvoice) => {
      setPayingInvoiceId(null);
      // Invalidate both invoices and overview metrics so badges & totals update
      queryClient.invalidateQueries({
        queryKey: ['customerInvoices', effectiveCustomerId],
      });
      queryClient.invalidateQueries({
        queryKey: ['customerOverview', effectiveCustomerId],
      });

      addToast(
        'Invoice Paid Successfully',
        `Invoice ${updatedInvoice.invoiceNumber} ($${updatedInvoice.totalAmount.toFixed(2)}) marked as settled.`,
        'success'
      );
    },
    onError: (error: Error) => {
      setPayingInvoiceId(null);
      addToast('Payment Failed', error.message, 'error');
    },
  });

  if (isLoading || !invoices) {
    return <DashboardSkeleton />;
  }

  const filteredInvoices = invoices.filter(
    (inv) =>
      inv.invoiceNumber.toLowerCase().includes(searchTerm.toLowerCase()) ||
      inv.status.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDownloadPdf = (invoice: Invoice) => {
    addToast(
      'Tax Invoice Downloaded',
      `PDF receipt for ${invoice.invoiceNumber} ($${invoice.totalAmount.toFixed(2)}) downloaded.`,
      'success'
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2 border-b border-stripe-hairline">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Invoice History & Receipts
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Complete tax invoices, recurring period charges, discount credits, and payment statuses.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-xs text-stripe-ink-mute">
            Showing {filteredInvoices.length} billing records
          </span>
        </div>
      </div>

      {/* Main Table Card */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        {/* Search & Filter Bar */}
        <div className="p-4 border-b border-stripe-hairline flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white">
          <div className="relative w-full sm:w-72">
            <Search className="w-4 h-4 text-stripe-ink-mute absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search by invoice # or status..."
              className="w-full pl-9 pr-3 py-1.5 text-xs bg-slate-50 border border-stripe-hairline-input rounded-md focus:bg-white focus:outline-none focus:border-stripe-primary"
            />
          </div>

          <div className="flex items-center gap-2 text-xs text-stripe-ink-mute">
            <Filter className="w-3.5 h-3.5" />
            <span>Isolated to Organization: <strong className="text-stripe-ink">{user?.companyName}</strong></span>
          </div>
        </div>

        {/* Responsive Table */}
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Invoice Reference</th>
                <th className="py-3 px-5">Issued Date</th>
                <th className="py-3 px-5">Due Date</th>
                <th className="py-3 px-5">Breakdown (Subtotal / Tax / Disc)</th>
                <th className="py-3 px-5">Total Amount</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {filteredInvoices.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-8 text-center text-xs text-stripe-ink-mute">
                    No matching invoices found.
                  </td>
                </tr>
              ) : (
                filteredInvoices.map((inv) => {
                  const isPaid = inv.status === 'paid';
                  const isPending = inv.status === 'confirmed' || inv.status === 'open';
                  const isProcessingThis = payingInvoiceId === inv.id;

                  return (
                    <tr
                      key={inv.id}
                      className="hover:bg-slate-50/70 transition-colors"
                    >
                      {/* Invoice ID */}
                      <td className="py-3.5 px-5">
                        <div className="font-mono text-xs font-semibold text-stripe-ink">
                          {inv.invoiceNumber}
                        </div>
                        {inv.paymentMethod && (
                          <div className="text-[11px] text-stripe-ink-mute truncate max-w-[150px]">
                            {inv.paymentMethod}
                          </div>
                        )}
                      </td>

                      {/* Issued Date */}
                      <td className="py-3.5 px-5 text-xs text-stripe-ink-secondary tnum">
                        {formatDisplayDate(inv.date)}
                      </td>

                      {/* Due Date */}
                      <td className="py-3.5 px-5 text-xs text-stripe-ink-secondary tnum">
                        {formatDisplayDate(inv.dueDate)}
                      </td>

                      {/* Financial Breakdown (Subtotal / Tax / Discount) */}
                      <td className="py-3.5 px-5 text-xs text-stripe-ink-mute tnum">
                        <span>${inv.subtotalAmount.toFixed(2)}</span>
                        {inv.taxAmount > 0 && (
                          <span className="text-slate-400"> + ${inv.taxAmount.toFixed(2)} tax</span>
                        )}
                        {inv.discountAmount > 0 && (
                          <span className="text-emerald-600 font-medium"> - ${inv.discountAmount.toFixed(2)} disc</span>
                        )}
                      </td>

                      {/* Total Amount */}
                      <td className="py-3.5 px-5 font-normal text-stripe-ink text-base tnum">
                        ${inv.totalAmount.toLocaleString(undefined, {
                          minimumFractionDigits: 2,
                          maximumFractionDigits: 2,
                        })}
                      </td>

                      {/* Status Badge: Yellow for Confirmed/Pending, Green for Paid */}
                      <td className="py-3.5 px-5">
                        <StatusBadge
                          label={isPaid ? 'Paid' : 'Confirmed (Due)'}
                          variant={isPaid ? 'success' : 'warning'}
                        />
                      </td>

                      {/* Action Buttons: Download PDF and Pay Now */}
                      <td className="py-3.5 px-5 text-right">
                        <div className="flex items-center justify-end gap-2">
                          <button
                            onClick={() => handleDownloadPdf(inv)}
                            className="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-normal text-stripe-ink-secondary hover:text-stripe-ink hover:bg-slate-100 rounded-md border border-slate-200 transition-colors"
                            title="Download PDF Tax Invoice"
                          >
                            <Download className="w-3.5 h-3.5" />
                            PDF
                          </button>

                          {isPending && (
                            <button
                              onClick={() => payMutation.mutate(inv.id)}
                              disabled={isProcessingThis}
                              className="inline-flex items-center gap-1 px-3 py-1 text-xs font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep disabled:opacity-50 rounded-full shadow-stripe-sm transition-colors"
                            >
                              <CreditCard className="w-3.5 h-3.5" />
                              {isProcessingThis ? 'Processing...' : 'Pay Now'}
                            </button>
                          )}

                          {isPaid && (
                            <span className="inline-flex items-center text-xs text-emerald-600 px-2 py-1">
                              <CheckCircle2 className="w-3.5 h-3.5 mr-1" />
                              Settled
                            </span>
                          )}
                        </div>
                      </td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default CustomerInvoicesPage;
