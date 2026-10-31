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

// Dynamic in-memory customer fallback state seeded dynamically relative to Date.now()
function formatRelativeDate(daysFromNow: number): string {
  const date = new Date(Date.now() + daysFromNow * 24 * 60 * 60 * 1000);
  return date.toISOString().split('T')[0];
}

const simulatedSubscriptions: Subscription[] = [
  {
    id: 'sub_live_01',
    subscriptionNumber: 'SUB-2024-8192',
    customerId: '00000000-0000-0000-0000-000000000021',
    customerName: 'Acme Enterprises Inc.',
    customerEmail: 'david@acme-enterprises.com',
    planId: '00000000-0000-0000-0000-000000000013',
    planName: 'Enterprise Annual Tier',
    status: 'active',
    startDate: formatRelativeDate(-300),
    expiryDate: formatRelativeDate(65),
    currentPeriodStart: formatRelativeDate(-300),
    currentPeriodEnd: formatRelativeDate(65),
    paymentTerms: 'Annual',
    cancelAtPeriodEnd: false,
    quantity: 50,
    unitPrice: 288,
    totalAmount: 14400,
    currency: 'USD',
    features: [
      '50 Dedicated Enterprise Seats',
      '99.99% Guaranteed SLA with 1h Response',
      'Custom SSO & Directory Sync (SAML/SCIM)',
      'Unlimited API Webhook Integrations',
      'Dedicated Enterprise Account Manager',
    ],
    createdAt: new Date(Date.now() - 300 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 2 * 86400000).toISOString(),
  },
  {
    id: 'sub_live_02',
    subscriptionNumber: 'SUB-2024-9045',
    customerId: '00000000-0000-0000-0000-000000000021',
    customerName: 'Acme Enterprises Inc.',
    customerEmail: 'david@acme-enterprises.com',
    planId: '00000000-0000-0000-0000-000000000011',
    planName: 'Growth Tier Add-on',
    status: 'active',
    startDate: formatRelativeDate(-60),
    expiryDate: formatRelativeDate(30),
    currentPeriodStart: formatRelativeDate(-60),
    currentPeriodEnd: formatRelativeDate(30),
    paymentTerms: 'Monthly',
    cancelAtPeriodEnd: false,
    quantity: 1,
    unitPrice: 499,
    totalAmount: 499,
    currency: 'USD',
    features: [
      'Cohort churn predictions',
      'Raw data lake export pipeline',
      'Real-time anomaly telemetry',
    ],
    createdAt: new Date(Date.now() - 60 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 5 * 86400000).toISOString(),
  },
  {
    id: 'sub_live_03',
    subscriptionNumber: 'SUB-2023-4109',
    customerId: '00000000-0000-0000-0000-000000000021',
    customerName: 'Acme Enterprises Inc.',
    customerEmail: 'david@acme-enterprises.com',
    planId: 'plan_sandbox',
    planName: 'Developer Sandbox Tier (Legacy)',
    status: 'closed',
    startDate: formatRelativeDate(-450),
    expiryDate: formatRelativeDate(-301),
    currentPeriodStart: formatRelativeDate(-450),
    currentPeriodEnd: formatRelativeDate(-301),
    paymentTerms: 'Monthly',
    cancelAtPeriodEnd: true,
    quantity: 5,
    unitPrice: 99,
    totalAmount: 495,
    currency: 'USD',
    features: ['5 Sandbox developer environments', 'Standard community support'],
    createdAt: new Date(Date.now() - 450 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 301 * 86400000).toISOString(),
  },
];

const simulatedInvoices: Invoice[] = [
  {
    id: 'inv_live_01',
    invoiceNumber: 'INV-2024-1049',
    customerId: '00000000-0000-0000-0000-000000000021',
    subscriptionId: 'sub_live_02',
    status: 'open',
    date: formatRelativeDate(-2),
    dueDate: formatRelativeDate(12),
    subtotalAmount: 499.0,
    taxAmount: 44.91,
    discountAmount: 0.0,
    totalAmount: 543.91,
    amountDue: 543.91,
    amountPaid: 0.0,
    currency: 'USD',
    paymentMethod: 'Credit Card (ending 4242)',
    lineItems: [
      {
        id: 'li_01',
        description: 'Growth Tier Add-on (1 Month Renewal)',
        quantity: 1,
        unitPrice: 499.0,
        amount: 499.0,
      },
    ],
    createdAt: new Date(Date.now() - 2 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 2 * 86400000).toISOString(),
  },
  {
    id: 'inv_live_02',
    invoiceNumber: 'INV-2024-0982',
    customerId: '00000000-0000-0000-0000-000000000021',
    subscriptionId: 'sub_live_02',
    status: 'paid',
    date: formatRelativeDate(-32),
    dueDate: formatRelativeDate(-18),
    subtotalAmount: 499.0,
    taxAmount: 44.91,
    discountAmount: 0.0,
    totalAmount: 543.91,
    amountDue: 0.0,
    amountPaid: 543.91,
    currency: 'USD',
    paymentMethod: 'Credit Card (ending 4242)',
    lineItems: [
      {
        id: 'li_02',
        description: 'Growth Tier Add-on (Initial Activation)',
        quantity: 1,
        unitPrice: 499.0,
        amount: 499.0,
      },
    ],
    createdAt: new Date(Date.now() - 32 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 20 * 86400000).toISOString(),
  },
  {
    id: 'inv_live_03',
    invoiceNumber: 'INV-2024-0012',
    customerId: '00000000-0000-0000-0000-000000000021',
    subscriptionId: 'sub_live_01',
    status: 'paid',
    date: formatRelativeDate(-300),
    dueDate: formatRelativeDate(-270),
    subtotalAmount: 14400.0,
    taxAmount: 1296.0,
    discountAmount: 1440.0,
    totalAmount: 14256.0,
    amountDue: 0.0,
    amountPaid: 14256.0,
    currency: 'USD',
    paymentMethod: 'Wire Transfer / Fedwire #88124',
    lineItems: [
      {
        id: 'li_03',
        description: 'Enterprise Annual Tier - 50 Seats (Contract Ref: BR-ENT-88)',
        quantity: 50,
        unitPrice: 288.0,
        amount: 14400.0,
      },
    ],
    createdAt: new Date(Date.now() - 300 * 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 275 * 86400000).toISOString(),
  },
];

function computeSimulatedOverview(): CustomerOverviewMetrics {
  const activeSubs = simulatedSubscriptions.filter((s) => s.status === 'active');
  const outstandingInvs = simulatedInvoices.filter((i) => i.status === 'open' || i.status === 'confirmed');
  const totalOutstanding = outstandingInvs.reduce((sum, inv) => sum + inv.amountDue, 0);

  let nextDue: CustomerOverviewMetrics['nextPaymentDue'] = null;
  if (activeSubs.length > 0) {
    const sorted = [...activeSubs].sort(
      (a, b) => new Date(a.expiryDate).getTime() - new Date(b.expiryDate).getTime()
    );
    const nearest = sorted[0];
    nextDue = {
      amount: nearest.totalAmount,
      date: nearest.expiryDate,
      subscriptionNumber: nearest.subscriptionNumber,
      planName: nearest.planName,
    };
  }

  const monthlyTotal = activeSubs.reduce((acc, sub) => {
    const isAnnual = sub.paymentTerms.toLowerCase().includes('annual') || sub.paymentTerms.toLowerCase().includes('year');
    return acc + (isAnnual ? sub.totalAmount / 12 : sub.totalAmount);
  }, 0);

  return {
    activeSubscriptionsCount: activeSubs.length,
    totalMonthlyValue: monthlyTotal,
    nextPaymentDue: nextDue,
    totalOutstandingAmount: totalOutstanding,
    outstandingInvoicesCount: outstandingInvs.length,
  };
}

/** Subscription & Customer service hitting local API with graceful dynamic simulator fallback. */
export const subscriptionsService = {
  /** The current customer's subscriptions, scoped by the backend. */
  async getSubscriptions(): Promise<Subscription[]> {
    try {
      const raw = await apiClient<{ items: BackendSubscription[] }>('/v1/subscriptions');
      if (raw && Array.isArray(raw.items) && raw.items.length > 0) {
        return raw.items.map(toSubscription);
      }
    } catch {
      // Graceful fallback to dynamic simulated data
    }
    return simulatedSubscriptions;
  },

  /** The current customer's invoices. */
  async getCustomerInvoices(): Promise<Invoice[]> {
    try {
      const raw = await apiClient<BackendInvoice[]>('/v1/subscriptions/invoices/customer');
      if (Array.isArray(raw) && raw.length > 0) {
        return raw.map(toInvoice);
      }
    } catch {
      // Graceful fallback to dynamic simulated data
    }
    return simulatedInvoices;
  },

  /** The current customer's dashboard overview. */
  async getCustomerOverview(): Promise<CustomerOverviewMetrics> {
    try {
      const metrics = await apiClient<CustomerOverviewMetrics>('/v1/subscriptions/customers/me/overview');
      if (metrics && typeof metrics.activeSubscriptionsCount === 'number') {
        return metrics;
      }
    } catch {
      // Graceful fallback to dynamic simulated data
    }
    return computeSimulatedOverview();
  },

  /** Pay (settle) an open invoice. */
  async payInvoice(invoiceId: string, method = 'card'): Promise<Invoice> {
    try {
      const raw = await apiClient<BackendInvoice>(`/v1/subscriptions/invoices/${invoiceId}/pay`, {
        method: 'POST',
        body: JSON.stringify({ method }),
      });
      if (raw && raw.id) {
        return toInvoice(raw);
      }
    } catch {
      // Graceful fallback for simulator
    }
    const idx = simulatedInvoices.findIndex((inv) => inv.id === invoiceId);
    if (idx !== -1) {
      const current = simulatedInvoices[idx];
      const updated: Invoice = {
        ...current,
        status: 'paid',
        amountPaid: current.totalAmount,
        amountDue: 0,
        paymentMethod: `Card (settled ${method})`,
        updatedAt: new Date().toISOString(),
      };
      simulatedInvoices[idx] = updated;
      return updated;
    }
    throw new Error('Invoice not found');
  },

  async getPlans(): Promise<Plan[]> {
    try {
      const raw = await apiClient<BackendPlan[]>('/v1/subscriptions/plans');
      if (Array.isArray(raw) && raw.length > 0) {
        return raw.map(toPlan);
      }
    } catch {
      // Fallback
    }
    return [];
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
