import { apiClient } from './api';
import type { Subscription, Plan, Invoice, CustomerOverviewMetrics } from '../types';

// Backend response shapes (camelCase, money returned as both dollars and cents).
interface BackendPlan {
  id: string;
  productId: string;
  productName?: string;
  name: string;
  description: string;
  billingInterval: 'monthly' | 'quarterly' | 'yearly';
  price: number;
  currency: string;
  features: string[];
  trialPeriodDays: number;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

interface BackendSubscription {
  id: string;
  customerId: string;
  customerName?: string;
  customerEmail?: string;
  planId: string;
  plan?: BackendPlan;
  status: string;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  cancelAtPeriodEnd: boolean;
  quantity: number;
  unitPrice: number;
  totalAmount: number;
  currency: string;
  billingInterval: string;
  createdAt: string;
  updatedAt: string;
}

interface BackendInvoice {
  id: string;
  invoiceNumber: string;
  customerId: string;
  subscriptionId?: string;
  status: string;
  date: string;
  dueDate: string;
  subtotalAmount: number;
  taxAmount: number;
  discountAmount: number;
  totalAmount: number;
  amountDue: number;
  amountPaid: number;
  currency: string;
  paymentMethod?: string;
}

export function formatDisplayDate(isoDateStr: string): string {
  try {
    const d = new Date(isoDateStr);
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  } catch {
    return isoDateStr;
  }
}

function paymentTerms(interval: string): string {
  switch (interval) {
    case 'quarterly':
      return 'Quarterly';
    case 'yearly':
      return 'Annual';
    default:
      return 'Monthly';
  }
}

function subscriptionNumber(id: string): string {
  return `SUB-${id.replace(/-/g, '').slice(0, 8).toUpperCase()}`;
}

function toPlan(plan: BackendPlan): Plan {
  return {
    id: plan.id,
    productId: plan.productId,
    name: plan.name,
    description: plan.description,
    billingInterval: plan.billingInterval,
    price: plan.price,
    currency: plan.currency,
    features: plan.features || [],
    trialPeriodDays: plan.trialPeriodDays,
    active: plan.active,
    createdAt: plan.createdAt,
    updatedAt: plan.updatedAt,
  };
}

function toSubscription(sub: BackendSubscription): Subscription {
  return {
    id: sub.id,
    subscriptionNumber: subscriptionNumber(sub.id),
    customerId: sub.customerId,
    customerName: sub.customerName,
    customerEmail: sub.customerEmail,
    planId: sub.planId,
    planName: sub.plan?.name || '',
    plan: sub.plan ? toPlan(sub.plan) : undefined,
    status: (sub.status as Subscription['status']) || 'active',
    startDate: sub.currentPeriodStart,
    expiryDate: sub.currentPeriodEnd,
    currentPeriodStart: sub.currentPeriodStart,
    currentPeriodEnd: sub.currentPeriodEnd,
    paymentTerms: paymentTerms(sub.billingInterval),
    cancelAtPeriodEnd: sub.cancelAtPeriodEnd,
    quantity: sub.quantity,
    unitPrice: sub.unitPrice,
    totalAmount: sub.totalAmount,
    currency: sub.currency,
    features: sub.plan?.features,
    createdAt: sub.createdAt,
    updatedAt: sub.updatedAt,
  };
}

function toInvoice(inv: BackendInvoice): Invoice {
  return {
    id: inv.id,
    invoiceNumber: inv.invoiceNumber,
    customerId: inv.customerId,
    subscriptionId: inv.subscriptionId,
    status: (inv.status as Invoice['status']) || 'open',
    date: inv.date,
    dueDate: inv.dueDate,
    subtotalAmount: inv.subtotalAmount,
    taxAmount: inv.taxAmount,
    discountAmount: inv.discountAmount,
    totalAmount: inv.totalAmount,
    amountDue: inv.amountDue,
    amountPaid: inv.amountPaid,
    currency: inv.currency,
    paymentMethod: inv.paymentMethod,
    lineItems: [],
    createdAt: inv.date,
    updatedAt: inv.date,
  };
}

/** Real subscription-service integration (no mock fallback). */
export const subscriptionsService = {
  /** The current customer's subscriptions, scoped by the backend. */
  async getSubscriptions(): Promise<Subscription[]> {
    const raw = await apiClient<{ items: BackendSubscription[] }>('/v1/subscriptions');
    return (raw.items || []).map(toSubscription);
  },

  /** The current customer's invoices. */
  async getCustomerInvoices(): Promise<Invoice[]> {
    const raw = await apiClient<BackendInvoice[]>('/v1/invoices/customer');
    return (raw || []).map(toInvoice);
  },

  /** The current customer's dashboard overview. */
  async getCustomerOverview(): Promise<CustomerOverviewMetrics> {
    return apiClient<CustomerOverviewMetrics>('/v1/customers/me/overview');
  },

  /** Pay (settle) an open invoice. */
  async payInvoice(invoiceId: string, method = 'card'): Promise<Invoice> {
    const raw = await apiClient<BackendInvoice>(`/v1/invoices/${invoiceId}/pay`, {
      method: 'POST',
      body: JSON.stringify({ method }),
    });
    return toInvoice(raw);
  },

  async getPlans(): Promise<Plan[]> {
    const raw = await apiClient<BackendPlan[]>('/v1/subscriptions/plans');
    return (raw || []).map(toPlan);
  },

  async cancelSubscription(id: string): Promise<void> {
    await apiClient<unknown>(`/v1/subscriptions/${id}/cancel?immediate=true`, { method: 'POST' });
  },

  async pauseSubscription(id: string): Promise<void> {
    await apiClient<unknown>(`/v1/subscriptions/${id}/pause`, { method: 'POST' });
  },

  async resumeSubscription(id: string): Promise<void> {
    await apiClient<unknown>(`/v1/subscriptions/${id}/resume`, { method: 'POST' });
  },

  async renewSubscription(id: string): Promise<void> {
    await apiClient<unknown>(`/v1/subscriptions/${id}/renew`, { method: 'POST' });
  },
};

export type { BackendSubscription, BackendPlan, BackendInvoice };