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
  | 'trialing'
  | 'past_due'
  | 'canceled'
  | 'paused'
  | 'unpaid';

export interface Subscription {
  id: string;
  customerId: string;
  customerName?: string;
  customerEmail?: string;
  planId: string;
  plan?: Plan;
  status: SubscriptionStatus;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  cancelAtPeriodEnd: boolean;
  quantity: number;
  unitPrice: number;
  totalAmount: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export type InvoiceStatus = 'draft' | 'open' | 'paid' | 'uncollectible' | 'void';

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
  amountDue: number;
  amountPaid: number;
  currency: string;
  dueDate: string;
  paidAt?: string;
  lineItems: InvoiceLineItem[];
  createdAt: string;
  updatedAt: string;
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
