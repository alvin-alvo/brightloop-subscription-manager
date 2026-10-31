export type UserRole = 'admin' | 'staff' | 'customer';

export type UserStatus = 'active' | 'inactive' | 'suspended';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  companyName?: string;
  phone?: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Product {
  id: string;
  name: string;
  description: string;
  active: boolean;
  features: string[];
  createdAt: string;
  updatedAt: string;
}

export type BillingInterval = 'monthly' | 'quarterly' | 'yearly';

export interface Plan {
  id: string;
  productId: string;
  name: string;
  description: string;
  billingInterval: BillingInterval;
  price: number;
  currency: string;
  features: string[];
  trialPeriodDays?: number;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export type SubscriptionStatus =
  | 'active'
  | 'closed'
  | 'trialing'
  | 'past_due'
  | 'canceled'
  | 'paused'
  | 'unpaid';

export interface Subscription {
  id: string;
  subscriptionNumber: string;
  customerId: string;
  customerName?: string;
  customerEmail?: string;
  planId: string;
  planName: string;
  plan?: Plan;
  status: SubscriptionStatus;
  startDate: string;
  expiryDate: string;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  paymentTerms: string;
  cancelAtPeriodEnd: boolean;
  quantity: number;
  unitPrice: number;
  totalAmount: number;
  currency: string;
  features?: string[];
  createdAt: string;
  updatedAt: string;
}

export type InvoiceStatus = 'draft' | 'open' | 'confirmed' | 'paid' | 'uncollectible' | 'void';

export interface InvoiceLineItem {
  id: string;
  description: string;
  quantity: number;
  unitPrice: number;
  amount: number;
}

export interface Invoice {
  id: string;
  invoiceNumber: string;
  customerId: string;
  customerName?: string;
  customerEmail?: string;
  subscriptionId?: string;
  status: InvoiceStatus;
  date: string;
  dueDate: string;
  subtotalAmount: number;
  taxAmount: number;
  discountAmount: number;
  totalAmount: number;
  amountDue: number;
  amountPaid: number;
  currency: string;
  paidAt?: string;
  paymentMethod?: string;
  lineItems: InvoiceLineItem[];
  createdAt: string;
  updatedAt: string;
}

export interface CustomerOverviewMetrics {
  activeSubscriptionsCount: number;
  totalMonthlyValue: number;
  nextPaymentDue: {
    amount: number;
    date: string;
    subscriptionNumber: string;
    planName: string;
  } | null;
  totalOutstandingAmount: number;
  outstandingInvoicesCount: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface AuthResponse {
  user: User;
  token: string;
}
