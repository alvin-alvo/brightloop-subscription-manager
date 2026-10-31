export interface OverdueInvoice {
  id: string;
  invoiceNumber: string;
  customerName: string;
  customerEmail: string;
  amountCents: number;
  dueDate: string;
  daysOverdue: number;
  status: 'overdue' | 'paid' | 'open';
  remindersSent: number;
  planName: string;
}

export interface ProductVariant {
  id: string;
  name: string;
  brand: string;
  priceAdjustmentCents: number;
}

export interface AdminProduct {
  id: string;
  name: string;
  type: 'Software' | 'Service' | 'Add-on' | 'Infrastructure';
  salesPriceCents: number;
  costPriceCents: number;
  notes: string;
  variants: ProductVariant[];
  active: boolean;
  createdAt: string;
}

export interface PlanBehavior {
  autoCloseAtEnd: boolean;
  allowEarlyClose: boolean;
  allowPause: boolean;
  allowRenew: boolean;
}

export interface AdminPlan {
  id: string;
  productId: string;
  productName: string;
  name: string;
  costCents: number;
  priceCents: number;
  billingFrequency: 'daily' | 'weekly' | 'monthly' | 'yearly';
  minQuantity: number;
  validFrom: string;
  validUntil: string;
  behaviors: PlanBehavior;
  subscribersCount: number;
  active: boolean;
  createdAt: string;
}

export interface AdminDiscount {
  id: string;
  code: string;
  name: string;
  type: 'percentage' | 'fixed_amount';
  value: number; // e.g. 20 for 20%, or 5000 cents ($50.00)
  minPurchaseAmountCents: number;
  minQuantity: number;
  startDate: string;
  endDate: string;
  usageLimit: number;
  redemptions: number;
  target: 'subscriptions' | 'products';
  targetProductIds?: string[];
  active: boolean;
  createdAt: string;
}

export interface AdminStaff {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'staff';
  department: string;
  status: 'active' | 'invited' | 'suspended';
  lastActive: string;
  createdAt: string;
}

export interface FinanceMetrics {
  activeSubscriptions: number;
  activeSubscriptionsGrowth: string;
  totalRevenueCents: number;
  totalRevenueGrowth: string;
  outstandingAmountCents: number;
  overdueInvoicesCount: number;
  overdueAmountCents: number;
  monthlyRecurringRevenueCents: number;
}

// Initial Seed Data (Realistic B2B SaaS Subscriptions & Invoices)
const INITIAL_INVOICES: OverdueInvoice[] = [
  {
    id: 'inv_1042',
    invoiceNumber: 'INV-2024-1042',
    customerName: 'Vertex Technologies',
    customerEmail: 'billing@vertextech.io',
    amountCents: 49900,
    dueDate: '2024-10-15',
    daysOverdue: 16,
    status: 'overdue',
    remindersSent: 2,
    planName: 'Growth Tier (Monthly)',
  },
  {
    id: 'inv_1039',
    invoiceNumber: 'INV-2024-1039',
    customerName: 'Helix Health Labs',
    customerEmail: 'accounts@helixhealth.com',
    amountCents: 125000,
    dueDate: '2024-10-10',
    daysOverdue: 21,
    status: 'overdue',
    remindersSent: 3,
    planName: 'Enterprise SLA Tier',
  },
  {
    id: 'inv_1033',
    invoiceNumber: 'INV-2024-1033',
    customerName: 'Nexus Cloud Infrastructure',
    customerEmail: 'finance@nexuscloud.net',
    amountCents: 345000,
    dueDate: '2024-09-28',
    daysOverdue: 33,
    status: 'overdue',
    remindersSent: 4,
    planName: 'Custom Enterprise Node',
  },
  {
    id: 'inv_1028',
    invoiceNumber: 'INV-2024-1028',
    customerName: 'Horizon Media Group',
    customerEmail: 'ap@horizonmedia.org',
    amountCents: 89000,
    dueDate: '2024-10-22',
    daysOverdue: 9,
    status: 'overdue',
    remindersSent: 1,
    planName: 'Growth Tier (Monthly)',
  },
  {
    id: 'inv_1025',
    invoiceNumber: 'INV-2024-1025',
    customerName: 'Apex Data Science',
    customerEmail: 'treasury@apexdata.ai',
    amountCents: 74900,
    dueDate: '2024-10-18',
    daysOverdue: 13,
    status: 'overdue',
    remindersSent: 2,
    planName: 'Pro Tier (Annual)',
  },
  {
    id: 'inv_1019',
    invoiceNumber: 'INV-2024-1019',
    customerName: 'Sterling Logistics Corp',
    customerEmail: 'billing@sterling-logistics.com',
    amountCents: 520000,
    dueDate: '2024-09-15',
    daysOverdue: 46,
    status: 'overdue',
    remindersSent: 5,
    planName: 'Enterprise Platform',
  },
  {
    id: 'inv_1012',
    invoiceNumber: 'INV-2024-1012',
    customerName: 'Cascade BioInformatics',
    customerEmail: 'ops@cascadebio.co',
    amountCents: 29900,
    dueDate: '2024-10-27',
    daysOverdue: 4,
    status: 'overdue',
    remindersSent: 1,
    planName: 'Starter Tier (Monthly)',
  },
];

const INITIAL_PRODUCTS: AdminProduct[] = [
  {
    id: 'prod_brightloop_core',
    name: 'Brightloop Subscription Suite',
    type: 'Software',
    salesPriceCents: 29900,
    costPriceCents: 5000,
    notes: 'Core SaaS subscription engine with multi-tenant recurring billing and dunning automation.',
    variants: [
      {
        id: 'var_odoo_560',
        name: 'Odoo Enterprise Connector',
        brand: 'Odoo',
        priceAdjustmentCents: 56000, // +$560.00
      },
      {
        id: 'var_sap_connector',
        name: 'SAP S/4HANA Sync Module',
        brand: 'SAP',
        priceAdjustmentCents: 85000, // +$850.00
      },
    ],
    active: true,
    createdAt: '2024-01-10T10:00:00Z',
  },
  {
    id: 'prod_ai_dunning',
    name: 'Autonomous Revenue Recovery AI',
    type: 'Add-on',
    salesPriceCents: 14900,
    costPriceCents: 2000,
    notes: 'Machine learning dunning retries that analyze banking patterns to maximize recovery.',
    variants: [
      {
        id: 'var_priority_routing',
        name: 'High-Volume Card Updater',
        brand: 'Visa/MC',
        priceAdjustmentCents: 19900,
      },
    ],
    active: true,
    createdAt: '2024-02-15T12:00:00Z',
  },
  {
    id: 'prod_dedicated_gateway',
    name: 'Dedicated Private Edge Gateway',
    type: 'Infrastructure',
    salesPriceCents: 85000,
    costPriceCents: 22000,
    notes: 'Isolated tenant deployment with single-tenant VPC peering and HIPAA/SOC2 compliance.',
    variants: [
      {
        id: 'var_eu_residency',
        name: 'EU Frankfurt Data Residency',
        brand: 'AWS EU',
        priceAdjustmentCents: 25000,
      },
    ],
    active: true,
    createdAt: '2024-03-01T08:00:00Z',
  },
];

const INITIAL_PLANS: AdminPlan[] = [
  {
    id: 'plan_starter_m',
    productId: 'prod_brightloop_core',
    productName: 'Brightloop Subscription Suite',
    name: 'Starter Tier (Monthly)',
    costCents: 5000,
    priceCents: 9900,
    billingFrequency: 'monthly',
    minQuantity: 1,
    validFrom: '2024-01-01',
    validUntil: '2025-12-31',
    behaviors: {
      autoCloseAtEnd: false,
      allowEarlyClose: true,
      allowPause: true,
      allowRenew: true,
    },
    subscribersCount: 420,
    active: true,
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'plan_growth_m',
    productId: 'prod_brightloop_core',
    productName: 'Brightloop Subscription Suite',
    name: 'Growth Tier (Monthly)',
    costCents: 8000,
    priceCents: 29900,
    billingFrequency: 'monthly',
    minQuantity: 1,
    validFrom: '2024-01-01',
    validUntil: '2025-12-31',
    behaviors: {
      autoCloseAtEnd: false,
      allowEarlyClose: true,
      allowPause: true,
      allowRenew: true,
    },
    subscribersCount: 615,
    active: true,
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'plan_enterprise_y',
    productId: 'prod_brightloop_core',
    productName: 'Brightloop Subscription Suite',
    name: 'Enterprise Tier (Yearly)',
    costCents: 30000,
    priceCents: 120000,
    billingFrequency: 'yearly',
    minQuantity: 5,
    validFrom: '2024-01-01',
    validUntil: '2026-12-31',
    behaviors: {
      autoCloseAtEnd: false,
      allowEarlyClose: false,
      allowPause: true,
      allowRenew: true,
    },
    subscribersCount: 249,
    active: true,
    createdAt: '2024-01-01T00:00:00Z',
  },
];

const INITIAL_DISCOUNTS: AdminDiscount[] = [
  {
    id: 'disc_annual_20',
    code: 'ANNUAL20',
    name: 'Annual Contract Incentive',
    type: 'percentage',
    value: 20,
    minPurchaseAmountCents: 100000,
    minQuantity: 1,
    startDate: '2024-01-01',
    endDate: '2025-12-31',
    usageLimit: 500,
    redemptions: 142,
    target: 'subscriptions',
    active: true,
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'disc_startup_50',
    code: 'STARTUP50',
    name: 'Early Stage Accelerator Program',
    type: 'percentage',
    value: 50,
    minPurchaseAmountCents: 29900,
    minQuantity: 1,
    startDate: '2024-06-01',
    endDate: '2025-06-01',
    usageLimit: 100,
    redemptions: 89,
    target: 'subscriptions',
    active: true,
    createdAt: '2024-06-01T00:00:00Z',
  },
  {
    id: 'disc_core_100',
    code: 'ODOCONNECT100',
    name: 'Odoo Connector Promotional Credit',
    type: 'fixed_amount',
    value: 10000, // $100.00 off
    minPurchaseAmountCents: 56000,
    minQuantity: 1,
    startDate: '2024-08-01',
    endDate: '2025-08-01',
    usageLimit: 200,
    redemptions: 64,
    target: 'products',
    targetProductIds: ['prod_brightloop_core'],
    active: true,
    createdAt: '2024-08-01T00:00:00Z',
  },
];

const INITIAL_STAFF: AdminStaff[] = [
  {
    id: 'stf_01',
    name: 'Alex Vance',
    email: 'alex.admin@brightloop.io',
    role: 'admin',
    department: 'Executive Leadership',
    status: 'active',
    lastActive: 'Just now',
    createdAt: '2024-01-15T08:00:00Z',
  },
  {
    id: 'stf_02',
    name: 'Sarah Miller',
    email: 'sarah.miller@brightloop.io',
    role: 'staff',
    department: 'Revenue & Operations',
    status: 'active',
    lastActive: '12m ago',
    createdAt: '2024-02-01T10:00:00Z',
  },
  {
    id: 'stf_03',
    name: 'Marcus Brody',
    email: 'marcus.brody@brightloop.io',
    role: 'staff',
    department: 'Customer Billing Support',
    status: 'active',
    lastActive: '2h ago',
    createdAt: '2024-03-12T14:30:00Z',
  },
  {
    id: 'stf_04',
    name: 'Elena Rostova',
    email: 'elena.rostova@brightloop.io',
    role: 'staff',
    department: 'Enterprise Sales',
    status: 'active',
    lastActive: 'Yesterday',
    createdAt: '2024-05-18T09:15:00Z',
  },
];

// Storage keys for reactive session persistence
const KEYS = {
  INVOICES: 'brightloop_admin_invoices',
  PRODUCTS: 'brightloop_admin_products',
  PLANS: 'brightloop_admin_plans',
  DISCOUNTS: 'brightloop_admin_discounts',
  STAFF: 'brightloop_admin_staff',
};

function readStore<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key);
    if (raw) return JSON.parse(raw);
  } catch {
    // Ignore fallback
  }
  return fallback;
}

function writeStore<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch {
    // Ignore storage quota
  }
}

// Simulated query delay strictly resolving in under 2 seconds (e.g. 800ms)
export const simulateQueryDelay = (ms = 800): Promise<void> =>
  new Promise((res) => setTimeout(res, ms));

// Format cents to USD currency string
export function formatCents(cents: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(cents / 100);
}

// Core Aggregations & Query Functions
export const mockAdminService = {
  // --- INVOICES & METRICS ---
  getFinanceMetrics(): FinanceMetrics {
    const invoices = readStore<OverdueInvoice[]>(KEYS.INVOICES, INITIAL_INVOICES);
    const plans = readStore<AdminPlan[]>(KEYS.PLANS, INITIAL_PLANS);

    const activeSubscriptions = plans.reduce((acc, p) => acc + (p.active ? p.subscribersCount : 0), 0);
    const overdueInvoices = invoices.filter((i) => i.status === 'overdue');
    const overdueInvoicesCount = overdueInvoices.length;
    const overdueAmountCents = overdueInvoices.reduce((acc, i) => acc + i.amountCents, 0);

    // Total Revenue (In): Simulated lifetime collected revenue
    const totalCollectedCents = 28455000 + invoices.filter((i) => i.status === 'paid').reduce((acc, i) => acc + i.amountCents, 0);

    // Outstanding Amount (Owed): all unpaid invoices
    const outstandingAmountCents = invoices
      .filter((i) => i.status === 'overdue' || i.status === 'open')
      .reduce((acc, i) => acc + i.amountCents, 0);

    // Monthly Recurring Revenue normalized
    const monthlyRecurringRevenueCents = plans.reduce((acc, p) => {
      if (!p.active) return acc;
      const monthlyEquivalent =
        p.billingFrequency === 'yearly'
          ? Math.round(p.priceCents / 12)
          : p.billingFrequency === 'weekly'
          ? p.priceCents * 4
          : p.billingFrequency === 'daily'
          ? p.priceCents * 30
          : p.priceCents;
      return acc + monthlyEquivalent * p.subscribersCount;
    }, 0);

    return {
      activeSubscriptions: activeSubscriptions || 1284,
      activeSubscriptionsGrowth: '+8.4%',
      totalRevenueCents: totalCollectedCents,
      totalRevenueGrowth: '+14.2%',
      outstandingAmountCents: outstandingAmountCents,
      overdueInvoicesCount: overdueInvoicesCount,
      overdueAmountCents: overdueAmountCents,
      monthlyRecurringRevenueCents: monthlyRecurringRevenueCents,
    };
  },

  getInvoices(filters?: {
    search?: string;
    daysRange?: string;
    amountRange?: string;
    status?: string;
  }): OverdueInvoice[] {
    let list = readStore<OverdueInvoice[]>(KEYS.INVOICES, INITIAL_INVOICES);

    if (filters) {
      if (filters.search) {
        const query = filters.search.toLowerCase();
        list = list.filter(
          (i) =>
            i.customerName.toLowerCase().includes(query) ||
            i.invoiceNumber.toLowerCase().includes(query) ||
            i.customerEmail.toLowerCase().includes(query)
        );
      }
      if (filters.status && filters.status !== 'all') {
        list = list.filter((i) => i.status === filters.status);
      }
      if (filters.daysRange && filters.daysRange !== 'all') {
        if (filters.daysRange === '1-15') {
          list = list.filter((i) => i.daysOverdue >= 1 && i.daysOverdue <= 15);
        } else if (filters.daysRange === '16-30') {
          list = list.filter((i) => i.daysOverdue >= 16 && i.daysOverdue <= 30);
        } else if (filters.daysRange === '30+') {
          list = list.filter((i) => i.daysOverdue > 30);
        }
      }
      if (filters.amountRange && filters.amountRange !== 'all') {
        if (filters.amountRange === 'under-1k') {
          list = list.filter((i) => i.amountCents < 100000);
        } else if (filters.amountRange === '1k-3k') {
          list = list.filter((i) => i.amountCents >= 100000 && i.amountCents <= 300000);
        } else if (filters.amountRange === 'over-3k') {
          list = list.filter((i) => i.amountCents > 300000);
        }
      }
    }

    return list;
  },

  markInvoicePaid(id: string): OverdueInvoice {
    const list = readStore<OverdueInvoice[]>(KEYS.INVOICES, INITIAL_INVOICES);
    const updated = list.map((i) => (i.id === id ? { ...i, status: 'paid' as const } : i));
    writeStore(KEYS.INVOICES, updated);
    return updated.find((i) => i.id === id)!;
  },

  sendDunningReminder(id: string): OverdueInvoice {
    const list = readStore<OverdueInvoice[]>(KEYS.INVOICES, INITIAL_INVOICES);
    const updated = list.map((i) =>
      i.id === id ? { ...i, remindersSent: i.remindersSent + 1 } : i
    );
    writeStore(KEYS.INVOICES, updated);
    return updated.find((i) => i.id === id)!;
  },

  // --- PRODUCTS & VARIANTS ---
  getProducts(): AdminProduct[] {
    return readStore<AdminProduct[]>(KEYS.PRODUCTS, INITIAL_PRODUCTS);
  },

  addProduct(productData: Omit<AdminProduct, 'id' | 'createdAt'>): AdminProduct {
    const list = readStore<AdminProduct[]>(KEYS.PRODUCTS, INITIAL_PRODUCTS);
    const newProduct: AdminProduct = {
      ...productData,
      id: `prod_${Date.now().toString(36)}_${Math.random().toString(36).substring(2, 5)}`,
      createdAt: new Date().toISOString(),
    };
    const updated = [newProduct, ...list];
    writeStore(KEYS.PRODUCTS, updated);
    return newProduct;
  },

  deleteProduct(id: string): void {
    const list = readStore<AdminProduct[]>(KEYS.PRODUCTS, INITIAL_PRODUCTS);
    const updated = list.filter((p) => p.id !== id);
    writeStore(KEYS.PRODUCTS, updated);
  },

  // --- PLANS ---
  getPlans(): AdminPlan[] {
    return readStore<AdminPlan[]>(KEYS.PLANS, INITIAL_PLANS);
  },

  addPlan(planData: Omit<AdminPlan, 'id' | 'createdAt' | 'subscribersCount'>): AdminPlan {
    const list = readStore<AdminPlan[]>(KEYS.PLANS, INITIAL_PLANS);
    const newPlan: AdminPlan = {
      ...planData,
      id: `plan_${Date.now().toString(36)}_${Math.random().toString(36).substring(2, 5)}`,
      subscribersCount: 0,
      createdAt: new Date().toISOString(),
    };
    const updated = [newPlan, ...list];
    writeStore(KEYS.PLANS, updated);
    return newPlan;
  },

  togglePlanStatus(id: string): AdminPlan {
    const list = readStore<AdminPlan[]>(KEYS.PLANS, INITIAL_PLANS);
    const updated = list.map((p) => (p.id === id ? { ...p, active: !p.active } : p));
    writeStore(KEYS.PLANS, updated);
    return updated.find((p) => p.id === id)!;
  },

  // --- DISCOUNTS ---
  getDiscounts(): AdminDiscount[] {
    return readStore<AdminDiscount[]>(KEYS.DISCOUNTS, INITIAL_DISCOUNTS);
  },

  addDiscount(discountData: Omit<AdminDiscount, 'id' | 'createdAt' | 'redemptions'>): AdminDiscount {
    const list = readStore<AdminDiscount[]>(KEYS.DISCOUNTS, INITIAL_DISCOUNTS);
    const newDiscount: AdminDiscount = {
      ...discountData,
      code: discountData.code.trim().toUpperCase(),
      id: `disc_${Date.now().toString(36)}_${Math.random().toString(36).substring(2, 5)}`,
      redemptions: 0,
      createdAt: new Date().toISOString(),
    };
    const updated = [newDiscount, ...list];
    writeStore(KEYS.DISCOUNTS, updated);
    return newDiscount;
  },

  toggleDiscountStatus(id: string): AdminDiscount {
    const list = readStore<AdminDiscount[]>(KEYS.DISCOUNTS, INITIAL_DISCOUNTS);
    const updated = list.map((d) => (d.id === id ? { ...d, active: !d.active } : d));
    writeStore(KEYS.DISCOUNTS, updated);
    return updated.find((d) => d.id === id)!;
  },

  // --- STAFF MANAGEMENT ---
  getStaffMembers(): AdminStaff[] {
    return readStore<AdminStaff[]>(KEYS.STAFF, INITIAL_STAFF);
  },

  addStaffMember(staffData: Omit<AdminStaff, 'id' | 'createdAt' | 'lastActive'>): AdminStaff {
    const list = readStore<AdminStaff[]>(KEYS.STAFF, INITIAL_STAFF);
    const newStaff: AdminStaff = {
      ...staffData,
      id: `stf_${Date.now().toString(36)}_${Math.random().toString(36).substring(2, 5)}`,
      lastActive: 'Just invited',
      createdAt: new Date().toISOString(),
    };
    const updated = [...list, newStaff];
    writeStore(KEYS.STAFF, updated);
    return newStaff;
  },

  updateStaffRole(id: string, role: 'admin' | 'staff'): AdminStaff {
    const list = readStore<AdminStaff[]>(KEYS.STAFF, INITIAL_STAFF);
    const updated = list.map((s) => (s.id === id ? { ...s, role } : s));
    writeStore(KEYS.STAFF, updated);
    return updated.find((s) => s.id === id)!;
  },
};

