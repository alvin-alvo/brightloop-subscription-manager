import React, { lazy } from 'react';
import {
  TrendingUp,
  Layers,
  Percent,
  ShieldCheck,
  Users,
  FileSpreadsheet,
  Repeat,
  Receipt,
  CreditCard,
  FileText,
  LayoutDashboard,
  type LucideIcon,
} from 'lucide-react';
import type { UserRole } from '../types';

// Lazy-loaded route views for optimal performance & code splitting
const FinanceDashboardPage = lazy(() => import('../pages/admin/FinanceDashboardPage'));
const ProductsPlansPage = lazy(() => import('../pages/admin/ProductsPlansPage'));
const DiscountsPage = lazy(() => import('../pages/admin/DiscountsPage'));
const StaffManagementPage = lazy(() => import('../pages/admin/StaffManagementPage'));

const CrmCustomersPage = lazy(() => import('../pages/staff/CrmCustomersPage'));
const QuotationsPage = lazy(() => import('../pages/staff/QuotationsPage'));
const SubscriptionsPage = lazy(() => import('../pages/staff/SubscriptionsPage'));
const BillingOperationsPage = lazy(() => import('../pages/staff/BillingOperationsPage'));

const CustomerDashboardPage = lazy(() => import('../pages/customer/CustomerDashboardPage'));
const CustomerSubscriptionsPage = lazy(() => import('../pages/customer/CustomerSubscriptionsPage'));
const CustomerInvoicesPage = lazy(() => import('../pages/customer/CustomerInvoicesPage'));

export interface AppRouteConfig {
  id: string;
  path: string;
  label: string;
  icon: LucideIcon;
  allowedRoles: UserRole[];
  section: string;
  element: React.LazyExoticComponent<React.FC>;
}

// Dynamic routing array in TypeScript (no static JSON)
export const APP_ROUTES: AppRouteConfig[] = [
  // --- ADMIN ROUTES (Admin Hub Control Center) ---
  {
    id: 'admin_finance',
    path: '/admin/finance',
    label: 'Finance Hub',
    icon: TrendingUp,
    allowedRoles: ['admin'],
    section: 'Administration',
    element: FinanceDashboardPage,
  },
  {
    id: 'admin_catalog',
    path: '/admin/catalog',
    label: 'Products & Plans',
    icon: Layers,
    allowedRoles: ['admin'],
    section: 'Administration',
    element: ProductsPlansPage,
  },
  {
    id: 'admin_discounts',
    path: '/admin/discounts',
    label: 'Discounts',
    icon: Percent,
    allowedRoles: ['admin'],
    section: 'Administration',
    element: DiscountsPage,
  },
  {
    id: 'admin_staff',
    path: '/admin/staff',
    label: 'Staff Management',
    icon: ShieldCheck,
    allowedRoles: ['admin'],
    section: 'Administration',
    element: StaffManagementPage,
  },

  // --- STAFF ROUTES ---
  {
    id: 'staff_crm',
    path: '/staff/customers',
    label: 'CRM / Customers',
    icon: Users,
    allowedRoles: ['staff', 'admin'],
    section: 'Operations & Sales',
    element: CrmCustomersPage,
  },
  {
    id: 'staff_quotations',
    path: '/staff/quotations',
    label: 'Quotations',
    icon: FileSpreadsheet,
    allowedRoles: ['staff', 'admin'],
    section: 'Operations & Sales',
    element: QuotationsPage,
  },
  {
    id: 'staff_subscriptions',
    path: '/staff/subscriptions',
    label: 'Subscriptions',
    icon: Repeat,
    allowedRoles: ['staff', 'admin'],
    section: 'Operations & Sales',
    element: SubscriptionsPage,
  },
  {
    id: 'staff_billing',
    path: '/staff/billing',
    label: 'Billing',
    icon: Receipt,
    allowedRoles: ['staff', 'admin'],
    section: 'Operations & Sales',
    element: BillingOperationsPage,
  },

  // --- CUSTOMER PORTAL ROUTES ---
  {
    id: 'customer_dashboard',
    path: '/customer/dashboard',
    label: 'Portal Overview',
    icon: LayoutDashboard,
    allowedRoles: ['customer'],
    section: 'Customer Portal',
    element: CustomerDashboardPage,
  },
  {
    id: 'customer_subscriptions',
    path: '/customer/subscriptions',
    label: 'My Subscriptions',
    icon: CreditCard,
    allowedRoles: ['customer'],
    section: 'Customer Portal',
    element: CustomerSubscriptionsPage,
  },
  {
    id: 'customer_invoices',
    path: '/customer/invoices',
    label: 'Invoices & Receipts',
    icon: FileText,
    allowedRoles: ['customer'],
    section: 'Customer Portal',
    element: CustomerInvoicesPage,
  },
];

// Helper to get default redirect route for each role
export function getDefaultRouteForRole(role: UserRole): string {
  switch (role) {
    case 'admin':
      return '/admin/finance';
    case 'staff':
      return '/staff/customers';
    case 'customer':
      return '/customer/dashboard';
    default:
      return '/customer/dashboard';
  }
}
