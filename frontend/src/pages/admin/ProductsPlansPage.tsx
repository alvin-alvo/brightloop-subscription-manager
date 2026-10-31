import React from 'react';
import { Plus, Check } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const ProductsPlansPage: React.FC = () => {
  const { addToast } = useToast();

  const plans = [
    {
      id: 'plan_starter',
      name: 'Starter Tier',
      interval: 'Monthly',
      price: '$99',
      subscribers: 420,
      status: 'active' as const,
      features: ['Up to 5 team seats', 'Standard SLA (24h)', '5,000 API calls/mo'],
    },
    {
      id: 'plan_growth',
      name: 'Growth Tier',
      interval: 'Monthly',
      price: '$299',
      subscribers: 615,
      status: 'active' as const,
      featured: true,
      features: ['Up to 25 team seats', 'Priority support (4h SLA)', '50,000 API calls/mo', 'Custom domains'],
    },
    {
      id: 'plan_enterprise',
      name: 'Enterprise Tier',
      interval: 'Annual',
      price: '$1,200',
      subscribers: 249,
      status: 'active' as const,
      features: ['Unlimited team seats', 'Dedicated account manager', 'Unlimited API calls', 'Custom integrations', '99.99% SLA'],
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Products & Subscription Plans
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Configure tiered SaaS catalog, billing cadences, and feature entitlements.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Plan Creator Opened', 'Configure new SaaS product or recurring tier.', 'info')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          Create New Plan
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {plans.map((plan) => (
          <div
            key={plan.id}
            className={`p-6 rounded-xl border transition-all ${
              plan.featured
                ? 'border-indigo-300 bg-white shadow-stripe-md relative ring-1 ring-stripe-primary/20'
                : 'border-stripe-hairline bg-white shadow-stripe-sm'
            }`}
          >
            {plan.featured && (
              <span className="absolute -top-3 left-6 px-3 py-0.5 text-[11px] font-medium bg-stripe-primary text-white rounded-full">
                Most Popular
              </span>
            )}
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-normal text-stripe-ink">{plan.name}</h2>
              <StatusBadge label="Active" variant="success" />
            </div>

            <div className="mt-4 flex items-baseline gap-1">
              <span className="text-3xl font-light text-stripe-ink tnum">{plan.price}</span>
              <span className="text-xs text-stripe-ink-mute">/{plan.interval.toLowerCase()}</span>
            </div>
            <p className="text-xs text-stripe-ink-secondary mt-1 tnum">
              {plan.subscribers} active paying accounts
            </p>

            <ul className="mt-6 space-y-2.5 text-xs text-stripe-ink-secondary border-t border-stripe-hairline pt-4">
              {plan.features.map((feature, idx) => (
                <li key={idx} className="flex items-center gap-2">
                  <Check className="w-3.5 h-3.5 text-emerald-600 shrink-0" />
                  <span>{feature}</span>
                </li>
              ))}
            </ul>

            <div className="mt-6 pt-4 border-t border-stripe-hairline flex items-center justify-between">
              <button
                onClick={() => addToast('Plan Config', `Editing parameters for ${plan.name}`, 'info')}
                className="text-xs font-normal text-stripe-primary hover:text-stripe-primary-deep"
              >
                Edit Tier
              </button>
              <span className="text-xs text-stripe-ink-mute font-mono">{plan.id}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ProductsPlansPage;
