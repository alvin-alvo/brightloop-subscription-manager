import type { Subscription, Invoice, CustomerOverviewMetrics } from '../types';

/**
 * Dynamic customer data simulator for Brightloop.
 * Uses relative dates anchored to Date.now() to fulfill dynamic data requirements
 * without static JSON files.
 */

// Format helper: YYYY-MM-DD
function formatRelativeDate(daysFromNow: number): string {
  const date = new Date(Date.now() + daysFromNow * 24 * 60 * 60 * 1000);
  return date.toISOString().split('T')[0];
}

// Format friendly date: "Nov 14, 2026"
export function formatDisplayDate(isoDateStr: string): string {
  try {
    const d = new Date(isoDateStr);
    return d.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  } catch {
    return isoDateStr;
  }
}

// In-memory state storage initialized dynamically per session
let dynamicSubscriptionsStore: Subscription[] = [];
let dynamicInvoicesStore: Invoice[] = [];

export function seedCustomerData(): void {
  const now = Date.now();

  // Subscriptions seeded relative to current timestamp
  dynamicSubscriptionsStore = [
    {
      id: 'sub_live_01',
      subscriptionNumber: 'SUB-2024-8192',
      customerId: 'usr_cust_03', // David Chen / Acme Enterprises Inc.
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      planId: 'plan_enterprise',
      planName: 'Enterprise Annual Tier',
      status: 'active',
      startDate: formatRelativeDate(-300), // ~10 months ago
      expiryDate: formatRelativeDate(65),  // Renewal in ~65 days
      currentPeriodStart: formatRelativeDate(-300),
      currentPeriodEnd: formatRelativeDate(65),
      paymentTerms: 'Billed Annually (Net 30)',
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
      createdAt: new Date(now - 300 * 86400000).toISOString(),
      updatedAt: new Date(now - 2 * 86400000).toISOString(),
    },
    {
      id: 'sub_live_02',
      subscriptionNumber: 'SUB-2024-9045',
      customerId: 'usr_cust_03',
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      planId: 'plan_analytics_addon',
      planName: 'Advanced Revenue Intelligence Add-on',
      status: 'active',
      startDate: formatRelativeDate(-60),
      expiryDate: formatRelativeDate(30),
      currentPeriodStart: formatRelativeDate(-60),
      currentPeriodEnd: formatRelativeDate(30),
      paymentTerms: 'Billed Monthly (Auto-Debit)',
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
      createdAt: new Date(now - 60 * 86400000).toISOString(),
      updatedAt: new Date(now - 5 * 86400000).toISOString(),
    },
    {
      id: 'sub_live_03',
      subscriptionNumber: 'SUB-2023-4109',
      customerId: 'usr_cust_03',
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      planId: 'plan_starter_legacy',
      planName: 'Developer Sandbox Tier (Legacy)',
      status: 'closed',
      startDate: formatRelativeDate(-450),
      expiryDate: formatRelativeDate(-301),
      currentPeriodStart: formatRelativeDate(-450),
      currentPeriodEnd: formatRelativeDate(-301),
      paymentTerms: 'Billed Monthly',
      cancelAtPeriodEnd: true,
      quantity: 5,
      unitPrice: 99,
      totalAmount: 495,
      currency: 'USD',
      features: ['5 Sandbox developer environments', 'Standard community support'],
      createdAt: new Date(now - 450 * 86400000).toISOString(),
      updatedAt: new Date(now - 301 * 86400000).toISOString(),
    },
    // Subscriptions for another customer (to demonstrate strict data isolation)
    {
      id: 'sub_live_other_01',
      subscriptionNumber: 'SUB-2024-5512',
      customerId: 'usr_cust_other_99',
      customerName: 'Global Logistics Corp',
      customerEmail: 'billing@globallogistics.com',
      planId: 'plan_growth',
      planName: 'Growth Tier',
      status: 'active',
      startDate: formatRelativeDate(-90),
      expiryDate: formatRelativeDate(275),
      currentPeriodStart: formatRelativeDate(-90),
      currentPeriodEnd: formatRelativeDate(275),
      paymentTerms: 'Billed Annually',
      cancelAtPeriodEnd: false,
      quantity: 1,
      unitPrice: 3588,
      totalAmount: 3588,
      currency: 'USD',
      features: ['25 team seats', 'Priority support'],
      createdAt: new Date(now - 90 * 86400000).toISOString(),
      updatedAt: new Date(now - 90 * 86400000).toISOString(),
    },
  ];

  // Invoices seeded relative to current timestamp
  dynamicInvoicesStore = [
    {
      id: 'inv_live_01',
      invoiceNumber: 'INV-2024-1004',
      customerId: 'usr_cust_03',
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      subscriptionId: 'sub_live_02',
      status: 'confirmed', // Unpaid/Confirmed awaiting payment
      date: formatRelativeDate(-4),
      dueDate: formatRelativeDate(10), // Due in 10 days
      subtotalAmount: 499.0,
      taxAmount: 44.91,
      discountAmount: 0.0,
      totalAmount: 543.91,
      amountDue: 543.91,
      amountPaid: 0.0,
      currency: 'USD',
      paymentMethod: 'Corporate ACH (Pending)',
      lineItems: [
        {
          id: 'li_01',
          description: 'Advanced Revenue Intelligence Add-on (Nov Billing)',
          quantity: 1,
          unitPrice: 499.0,
          amount: 499.0,
        },
      ],
      createdAt: new Date(now - 4 * 86400000).toISOString(),
      updatedAt: new Date(now - 4 * 86400000).toISOString(),
    },
    {
      id: 'inv_live_02',
      invoiceNumber: 'INV-2024-0902',
      customerId: 'usr_cust_03',
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      subscriptionId: 'sub_live_02',
      status: 'paid',
      date: formatRelativeDate(-34),
      dueDate: formatRelativeDate(-20),
      subtotalAmount: 499.0,
      taxAmount: 44.91,
      discountAmount: 50.0, // $50 promotional coupon applied
      totalAmount: 493.91,
      amountDue: 0.0,
      amountPaid: 493.91,
      currency: 'USD',
      paidAt: formatRelativeDate(-22),
      paymentMethod: 'Mastercard ending in 9821',
      lineItems: [
        {
          id: 'li_02',
          description: 'Advanced Revenue Intelligence Add-on (Oct Billing)',
          quantity: 1,
          unitPrice: 499.0,
          amount: 499.0,
        },
      ],
      createdAt: new Date(now - 34 * 86400000).toISOString(),
      updatedAt: new Date(now - 22 * 86400000).toISOString(),
    },
    {
      id: 'inv_live_03',
      invoiceNumber: 'INV-2024-0101',
      customerId: 'usr_cust_03',
      customerName: 'Acme Enterprises Inc.',
      customerEmail: 'david@acme-enterprises.com',
      subscriptionId: 'sub_live_01',
      status: 'paid',
      date: formatRelativeDate(-300),
      dueDate: formatRelativeDate(-270),
      subtotalAmount: 14400.0,
      taxAmount: 1296.0,
      discountAmount: 1440.0, // 10% annual prepayment discount
      totalAmount: 14256.0,
      amountDue: 0.0,
      amountPaid: 14256.0,
      currency: 'USD',
      paidAt: formatRelativeDate(-275),
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
      createdAt: new Date(now - 300 * 86400000).toISOString(),
      updatedAt: new Date(now - 275 * 86400000).toISOString(),
    },
    // Another customer's invoice (must be filtered out for customerId 'usr_cust_03')
    {
      id: 'inv_live_other_01',
      invoiceNumber: 'INV-2024-9981',
      customerId: 'usr_cust_other_99',
      customerName: 'Global Logistics Corp',
      customerEmail: 'billing@globallogistics.com',
      subscriptionId: 'sub_live_other_01',
      status: 'paid',
      date: formatRelativeDate(-15),
      dueDate: formatRelativeDate(-1),
      subtotalAmount: 3588.0,
      taxAmount: 322.92,
      discountAmount: 0.0,
      totalAmount: 3910.92,
      amountDue: 0.0,
      amountPaid: 3910.92,
      currency: 'USD',
      paidAt: formatRelativeDate(-2),
      lineItems: [],
      createdAt: new Date(now - 15 * 86400000).toISOString(),
      updatedAt: new Date(now - 2 * 86400000).toISOString(),
    },
  ];
}

// Initial seed
seedCustomerData();

/**
 * Service API: Fetch Subscriptions with Data Isolation
 */
export async function fetchCustomerSubscriptions(
  customerId: string
): Promise<Subscription[]> {
  // Simulate network delay (< 500ms for high perceived performance)
  await new Promise((resolve) => setTimeout(resolve, 350));

  // Strict isolation filter
  return dynamicSubscriptionsStore.filter(
    (sub) => sub.customerId === customerId
  );
}

/**
 * Service API: Fetch Invoices with Data Isolation
 */
export async function fetchCustomerInvoices(
  customerId: string
): Promise<Invoice[]> {
  await new Promise((resolve) => setTimeout(resolve, 350));

  // Strict isolation filter
  return dynamicInvoicesStore.filter(
    (inv) => inv.customerId === customerId
  );
}

/**
 * Service API: Fetch Customer Dashboard Overview Metrics
 */
export async function fetchCustomerOverview(
  customerId: string
): Promise<CustomerOverviewMetrics> {
  await new Promise((resolve) => setTimeout(resolve, 400));

  const customerSubs = dynamicSubscriptionsStore.filter(
    (s) => s.customerId === customerId && s.status === 'active'
  );

  const customerInvoices = dynamicInvoicesStore.filter(
    (i) => i.customerId === customerId
  );

  const outstandingInvoices = customerInvoices.filter(
    (i) => i.status === 'confirmed' || i.status === 'open'
  );

  const totalOutstanding = outstandingInvoices.reduce(
    (acc, curr) => acc + curr.amountDue,
    0
  );

  // Find nearest upcoming renewal date
  let nextPaymentDue: CustomerOverviewMetrics['nextPaymentDue'] = null;
  if (customerSubs.length > 0) {
    const sorted = [...customerSubs].sort(
      (a, b) => new Date(a.expiryDate).getTime() - new Date(b.expiryDate).getTime()
    );
    const nearest = sorted[0];
    nextPaymentDue = {
      amount: nearest.totalAmount,
      date: nearest.expiryDate,
      subscriptionNumber: nearest.subscriptionNumber,
      planName: nearest.planName,
    };
  }

  const totalMonthlyValue = customerSubs.reduce((acc, curr) => {
    // If annual, approximate monthly equivalent
    const monthlyRate =
      curr.paymentTerms.toLowerCase().includes('annual')
        ? curr.totalAmount / 12
        : curr.totalAmount;
    return acc + monthlyRate;
  }, 0);

  return {
    activeSubscriptionsCount: customerSubs.length,
    totalMonthlyValue,
    nextPaymentDue,
    totalOutstandingAmount: totalOutstanding,
    outstandingInvoicesCount: outstandingInvoices.length,
  };
}

/**
 * Action API: Pay an Outstanding Invoice
 */
export async function payCustomerInvoice(
  invoiceId: string,
  customerId: string
): Promise<Invoice> {
  await new Promise((resolve) => setTimeout(resolve, 600));

  const index = dynamicInvoicesStore.findIndex(
    (inv) => inv.id === invoiceId && inv.customerId === customerId
  );

  if (index === -1) {
    throw new Error('Invoice not found or unauthorized');
  }

  const target = dynamicInvoicesStore[index];
  const updated: Invoice = {
    ...target,
    status: 'paid',
    amountPaid: target.totalAmount,
    amountDue: 0.0,
    paidAt: new Date().toISOString(),
    paymentMethod: 'Instant Corporate Card (Verified)',
    updatedAt: new Date().toISOString(),
  };

  dynamicInvoicesStore[index] = updated;
  return updated;
}

