import React, { useState, useEffect } from 'react';
import {
  DollarSign,
  Users,
  AlertTriangle,
  Clock,
  ArrowUpRight,
  RefreshCw,
  Search,
  Mail,
  CheckCircle2,
  Download,
  Filter,
} from 'lucide-react';
import {
  mockAdminService,
  formatCents,
  simulateQueryDelay,
  type FinanceMetrics,
  type OverdueInvoice,
} from '../../services/mockAdminData';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { DashboardSkeleton } from '../../components/layout/DashboardSkeleton';

export const FinanceDashboardPage: React.FC = () => {
  const { addToast } = useToast();

  const [isLoading, setIsLoading] = useState(true);
  const [metrics, setMetrics] = useState<FinanceMetrics | null>(null);
  const [invoices, setInvoices] = useState<OverdueInvoice[]>([]);

  // Filtering state
  const [search, setSearch] = useState('');
  const [daysRange, setDaysRange] = useState('all');
  const [amountRange, setAmountRange] = useState('all');
  const [statusFilter, setStatusFilter] = useState<'overdue' | 'all'>('overdue');

  // Load finance data simulating query across thousands of records in < 2 seconds
  const loadData = async (showSkeleton = true) => {
    if (showSkeleton) setIsLoading(true);
    await simulateQueryDelay(750); // Simulates fast optimized aggregation in ~750ms

    const newMetrics = mockAdminService.getFinanceMetrics();
    const newInvoices = mockAdminService.getInvoices({
      search,
      daysRange,
      amountRange,
      status: statusFilter === 'all' ? undefined : 'overdue',
    });

    setMetrics(newMetrics);
    setInvoices(newInvoices);
    setIsLoading(false);
  };

  useEffect(() => {
    loadData(true);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [daysRange, amountRange, statusFilter]);

  // Handle text search filter
  useEffect(() => {
    const timer = setTimeout(() => {
      const filtered = mockAdminService.getInvoices({
        search,
        daysRange,
        amountRange,
        status: statusFilter === 'all' ? undefined : 'overdue',
      });
      setInvoices(filtered);
    }, 200);
    return () => clearTimeout(timer);
  }, [search, daysRange, amountRange, statusFilter]);

  const handleMarkPaid = (id: string, number: string, amountCents: number) => {
    mockAdminService.markInvoicePaid(id);
    addToast(
      'Invoice Settled',
      `${number} marked as paid (${formatCents(amountCents)}). Financial ledgers updated.`,
      'success'
    );
    loadData(false);
  };

  const handleSendReminder = (id: string, customer: string) => {
    const updated = mockAdminService.sendDunningReminder(id);
    addToast(
      'Dunning Reminder Sent',
      `Payment notice #${updated.remindersSent} dispatched to ${customer}.`,
      'info'
    );
    loadData(false);
  };

  if (isLoading || !metrics) {
    return <DashboardSkeleton />;
  }

  const metricCards = [
    {
      title: 'Active Subscriptions',
      value: metrics.activeSubscriptions.toLocaleString(),
      change: metrics.activeSubscriptionsGrowth,
      subtitle: `${formatCents(metrics.monthlyRecurringRevenueCents)} normalized MRR`,
      icon: Users,
      trend: 'up',
      tone: 'emerald',
    },
    {
      title: 'Total Revenue (In)',
      value: formatCents(metrics.totalRevenueCents),
      change: metrics.totalRevenueGrowth,
      subtitle: 'Net lifetime collected payments',
      icon: DollarSign,
      trend: 'up',
      tone: 'emerald',
    },
    {
      title: 'Outstanding Amount (Owed)',
      value: formatCents(metrics.outstandingAmountCents),
      change: 'Due now',
      subtitle: 'Open & aging accounts receivable',
      icon: Clock,
      trend: 'warning',
      tone: 'amber',
    },
    {
      title: 'Overdue Invoices',
      value: `${metrics.overdueInvoicesCount} Invoices`,
      change: formatCents(metrics.overdueAmountCents),
      subtitle: 'Requiring dunning follow-up',
      icon: AlertTriangle,
      trend: 'danger',
      tone: 'rose',
    },
  ];

  return (
    <div className="space-y-6 animate-fadeIn">
      {/* Page Title & Operational Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
              Finance & Revenue Control Center
            </h1>
            <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 uppercase tracking-wider">
              Live
            </span>
          </div>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Executive financial telemetry, active subscriptions, revenue realization, and dunning workflows.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            onClick={() => loadData(true)}
            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-medium text-stripe-ink-secondary bg-white hover:bg-slate-50 border border-stripe-hairline rounded-full shadow-stripe-sm transition-colors"
            title="Simulate querying thousands of subscriptions"
          >
            <RefreshCw className="w-3.5 h-3.5 text-slate-500" />
            <span>Query Subscriptions</span>
          </button>

          <button
            onClick={() =>
              addToast(
                'Financial Statement Downloaded',
                'Comprehensive AR aging and MRR report exported to CSV.',
                'success'
              )
            }
            className="inline-flex items-center gap-1.5 px-4 py-2 text-xs font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
          >
            <Download className="w-3.5 h-3.5" />
            <span>Export AR Ledger</span>
          </button>
        </div>
      </div>

      {/* Top-Level Metric Cards Grid (Stripe 4-up) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {metricCards.map((card) => {
          const Icon = card.icon;
          return (
            <div
              key={card.title}
              className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm hover:shadow-stripe-md transition-shadow relative overflow-hidden"
            >
              <div className="flex items-center justify-between mb-2">
                <span className="text-[11px] font-medium uppercase tracking-wider text-stripe-ink-mute">
                  {card.title}
                </span>
                <div className="p-1.5 rounded-lg bg-slate-50 border border-slate-100 text-stripe-ink-secondary">
                  <Icon className="w-4 h-4 text-stripe-primary" />
                </div>
              </div>

              <div className="flex items-baseline gap-2">
                <span className="text-2xl font-light text-stripe-ink tracking-tight tnum">
                  {card.value}
                </span>
                {card.trend === 'up' && (
                  <span className="inline-flex items-center text-xs font-medium text-emerald-600 tnum">
                    <ArrowUpRight className="w-3.5 h-3.5 mr-0.5" />
                    {card.change}
                  </span>
                )}
                {card.trend === 'danger' && (
                  <span className="inline-flex items-center text-xs font-medium text-stripe-ruby tnum">
                    {card.change}
                  </span>
                )}
              </div>

              <p className="text-xs text-stripe-ink-mute mt-1.5 truncate">
                {card.subtitle}
              </p>
            </div>
          );
        })}
      </div>

      {/* Overdue Invoices Table Section */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        {/* Table Header Bar */}
        <div className="p-5 border-b border-stripe-hairline flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2">
              <h2 className="text-base font-normal text-stripe-ink">
                Overdue Invoices & Aging Receivables
              </h2>
              <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-rose-50 text-stripe-ruby border border-rose-200 tnum">
                {invoices.filter((i) => i.status === 'overdue').length} Pending Action
              </span>
            </div>
            <p className="text-xs text-stripe-ink-mute mt-0.5">
              Actionable overdue billing runs requiring customer dunning, card retries, or manual reconciliation.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => {
                const overdueList = invoices.filter((i) => i.status === 'overdue');
                overdueList.forEach((inv) => mockAdminService.sendDunningReminder(inv.id));
                loadData(false);
                addToast(
                  'Batch Dunning Run Dispatched',
                  `Payment notices delivered to ${overdueList.length} overdue customer accounts.`,
                  'info'
                );
              }}
              className="px-3 py-1.5 text-xs font-medium text-stripe-primary bg-indigo-50/80 hover:bg-indigo-100/80 border border-indigo-100 rounded-full transition-colors flex items-center gap-1.5"
            >
              <Mail className="w-3.5 h-3.5" />
              <span>Remind All Overdue</span>
            </button>
          </div>
        </div>

        {/* Dense Filtering Toolbar */}
        <div className="p-4 bg-slate-50/60 border-b border-stripe-hairline flex flex-col lg:flex-row lg:items-center gap-3 justify-between">
          {/* Search Field */}
          <div className="relative flex-1 max-w-sm">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search customer, invoice #, or email..."
              className="w-full bg-white text-stripe-ink text-xs rounded-md pl-9 pr-3 py-2 border border-stripe-hairline-input placeholder:text-slate-400 focus:outline-none focus:border-stripe-primary focus:ring-1 focus:ring-stripe-primary"
            />
          </div>

          {/* Filter Dropdowns */}
          <div className="flex flex-wrap items-center gap-2.5">
            {/* Days Overdue Filter */}
            <div className="flex items-center gap-1.5 text-xs">
              <span className="text-stripe-ink-mute hidden sm:inline flex items-center gap-1">
                <Filter className="w-3 h-3" /> Aging:
              </span>
              <select
                value={daysRange}
                onChange={(e) => setDaysRange(e.target.value)}
                className="bg-white border border-stripe-hairline-input text-stripe-ink-secondary text-xs rounded-md px-2.5 py-1.5 focus:outline-none focus:border-stripe-primary"
              >
                <option value="all">All Days</option>
                <option value="1-15">1–15 Days Overdue</option>
                <option value="16-30">16–30 Days Overdue</option>
                <option value="30+">30+ Days Overdue (Critical)</option>
              </select>
            </div>

            {/* Amount Range Filter */}
            <div className="flex items-center gap-1.5 text-xs">
              <select
                value={amountRange}
                onChange={(e) => setAmountRange(e.target.value)}
                className="bg-white border border-stripe-hairline-input text-stripe-ink-secondary text-xs rounded-md px-2.5 py-1.5 focus:outline-none focus:border-stripe-primary"
              >
                <option value="all">All Amounts</option>
                <option value="under-1k">Under $1,000</option>
                <option value="1k-3k">$1,000 – $3,000</option>
                <option value="over-3k">Over $3,000</option>
              </select>
            </div>

            {/* Status Switcher */}
            <div className="flex items-center gap-1.5 text-xs">
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value as 'overdue' | 'all')}
                className="bg-white border border-stripe-hairline-input text-stripe-ink-secondary text-xs rounded-md px-2.5 py-1.5 focus:outline-none focus:border-stripe-primary"
              >
                <option value="overdue">Overdue Only</option>
                <option value="all">All Invoices</option>
              </select>
            </div>
          </div>
        </div>

        {/* High-Density Data Table */}
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Invoice #</th>
                <th className="py-3 px-5">Customer & Organization</th>
                <th className="py-3 px-5">Subscription Plan</th>
                <th className="py-3 px-5">Due Date</th>
                <th className="py-3 px-5">Aging Status</th>
                <th className="py-3 px-5">Amount Due</th>
                <th className="py-3 px-5 text-right">Dunning Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {invoices.length === 0 ? (
                <tr>
                  <td colSpan={7} className="py-12 text-center text-stripe-ink-mute text-sm">
                    No overdue invoices match the specified criteria.
                  </td>
                </tr>
              ) : (
                invoices.map((inv) => (
                  <tr key={inv.id} className="hover:bg-slate-50/70 transition-colors">
                    {/* Invoice ID */}
                    <td className="py-3 px-5 font-mono text-xs font-semibold text-stripe-ink">
                      {inv.invoiceNumber}
                    </td>

                    {/* Customer Info */}
                    <td className="py-3 px-5">
                      <div className="font-medium text-stripe-ink">{inv.customerName}</div>
                      <div className="text-xs text-stripe-ink-mute">{inv.customerEmail}</div>
                    </td>

                    {/* Plan */}
                    <td className="py-3 px-5 text-xs text-stripe-ink-secondary">
                      {inv.planName}
                    </td>

                    {/* Due Date */}
                    <td className="py-3 px-5 text-xs text-stripe-ink-mute tnum">
                      {inv.dueDate}
                    </td>

                    {/* Overdue Badge */}
                    <td className="py-3 px-5">
                      {inv.status === 'paid' ? (
                        <StatusBadge label="Paid / Settled" variant="success" />
                      ) : inv.daysOverdue > 30 ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-rose-50 text-stripe-ruby border border-rose-200 tnum">
                          {inv.daysOverdue}d overdue (Escalated)
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-amber-50 text-amber-700 border border-amber-200 tnum">
                          {inv.daysOverdue}d overdue
                        </span>
                      )}
                    </td>

                    {/* Amount */}
                    <td className="py-3 px-5 font-medium text-stripe-ink tnum">
                      {formatCents(inv.amountCents)}
                    </td>

                    {/* Operational Actions */}
                    <td className="py-3 px-5 text-right">
                      {inv.status === 'paid' ? (
                        <span className="inline-flex items-center gap-1 text-xs text-emerald-600 font-medium">
                          <CheckCircle2 className="w-3.5 h-3.5" /> Reconciled
                        </span>
                      ) : (
                        <div className="inline-flex items-center gap-2">
                          <button
                            onClick={() => handleSendReminder(inv.id, inv.customerName)}
                            className="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-normal text-stripe-ink-secondary bg-white hover:bg-slate-50 border border-stripe-hairline rounded-md shadow-sm transition-colors"
                            title={`Reminders sent: ${inv.remindersSent}`}
                          >
                            <Mail className="w-3 h-3 text-slate-400" />
                            <span>Remind ({inv.remindersSent})</span>
                          </button>

                          <button
                            onClick={() => handleMarkPaid(inv.id, inv.invoiceNumber, inv.amountCents)}
                            className="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-md shadow-sm transition-colors"
                          >
                            <CheckCircle2 className="w-3 h-3" />
                            <span>Mark Paid</span>
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default FinanceDashboardPage;
