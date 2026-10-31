import React, { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { ArrowRight, CalendarDays, Plus } from 'lucide-react';
import { useToast } from '../../context/ToastContext';
import { fetchStaffWorkspace, type PipelineStage, type StaffSubscription } from '../../utils/mockStaffData';

const stages: PipelineStage[] = ['Draft', 'Quotation', 'Confirmed', 'Active', 'Closed'];
const currency = (cents: number): string =>
  new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(cents / 100);

export const SubscriptionsPage: React.FC = () => {
  const { addToast } = useToast();
  const { data, isLoading } = useQuery({ queryKey: ['staff-workspace'], queryFn: fetchStaffWorkspace });
  const [subscriptions, setSubscriptions] = useState<StaffSubscription[]>([]);

  useEffect(() => {
    if (data) setSubscriptions(data.subscriptions);
  }, [data]);

  const moveTo = (id: string, stage: PipelineStage) => {
    setSubscriptions((current) => current.map((item) => item.id === id ? { ...item, stage, updatedAt: new Date().toISOString().slice(0, 10) } : item));
    const subscription = subscriptions.find((item) => item.id === id);
    if (subscription) addToast('Pipeline updated', `${subscription.customer} moved to ${stage}.`, 'success');
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div><p className="text-xs font-medium uppercase text-stripe-ink-mute">Revenue operations</p><h1 className="mt-1 text-2xl font-light text-stripe-ink">Subscription pipeline</h1><p className="mt-1 text-sm text-stripe-ink-mute">A shared view of every customer agreement and its current stage.</p></div>
        <button type="button" onClick={() => addToast('Subscription draft started', 'New subscription added to the Draft stage.', 'info')} className="inline-flex items-center justify-center gap-2 rounded-full bg-stripe-primary px-4 py-2.5 text-sm font-medium text-white shadow-stripe-sm hover:bg-stripe-primary-deep"><Plus className="h-4 w-4" /> New subscription</button>
      </div>

      {isLoading ? <div className="rounded-lg border border-stripe-hairline bg-white p-12 text-center text-sm text-stripe-ink-mute">Loading subscription pipeline...</div> : <div className="flex gap-4 overflow-x-auto pb-3">
        {stages.map((stage, index) => {
          const stageItems = subscriptions.filter((item) => item.stage === stage);
          return <section key={stage} aria-label={`${stage} subscriptions`} className="w-[min(82vw,300px)] shrink-0 rounded-lg border border-stripe-hairline bg-slate-50/80 p-3 sm:w-[280px]">
            <div className="mb-3 flex items-center justify-between px-1"><div className="flex items-center gap-2"><h2 className="text-sm font-semibold text-stripe-ink">{stage}</h2><span className="rounded-full bg-white px-2 py-0.5 text-xs tabular-nums text-stripe-ink-mute">{stageItems.length}</span></div><span className="text-[10px] uppercase text-stripe-ink-mute">{String(index + 1).padStart(2, '0')}</span></div>
            <div className="min-h-28 space-y-3">
              {stageItems.map((subscription) => <article key={subscription.id} className="rounded-md border border-stripe-hairline bg-white p-4 shadow-stripe-sm">
                <div className="flex items-start justify-between gap-3"><div><p className="font-medium text-stripe-ink">{subscription.customer}</p><p className="mt-1 font-mono text-[11px] text-stripe-ink-mute">{subscription.id}</p></div><span className="rounded bg-indigo-50 px-2 py-1 text-[10px] font-medium text-stripe-primary">{subscription.plan.split(' ')[0]}</span></div>
                <div className="mt-4 flex items-center justify-between border-t border-stripe-hairline pt-3"><span className="text-xs text-stripe-ink-mute">Monthly value</span><span className="text-sm font-medium tabular-nums text-stripe-ink">{currency(subscription.monthlyValueCents)}</span></div>
                <p className="mt-2 flex items-center gap-1.5 text-[11px] text-stripe-ink-mute"><CalendarDays className="h-3.5 w-3.5" /> Updated {new Date(`${subscription.updatedAt}T00:00:00`).toLocaleDateString()}</p>
                <label className="mt-3 flex items-center gap-2 border-t border-stripe-hairline pt-3 text-xs text-stripe-ink-secondary">Move to <select aria-label={`Move ${subscription.customer} to stage`} value={subscription.stage} onChange={(event) => moveTo(subscription.id, event.target.value as PipelineStage)} className="min-w-0 flex-1 rounded border border-stripe-hairline-input bg-white px-2 py-1.5 text-xs text-stripe-ink focus:border-stripe-primary focus:outline-none">{stages.map((nextStage) => <option key={nextStage} value={nextStage}>{nextStage}</option>)}</select><ArrowRight className="h-3.5 w-3.5 text-stripe-ink-mute" /></label>
              </article>)}
              {stageItems.length === 0 && <p className="rounded-md border border-dashed border-slate-300 px-3 py-6 text-center text-xs text-stripe-ink-mute">No subscriptions in this stage</p>}
            </div>
          </section>;
        })}
      </div>}
    </div>
  );
};

export default SubscriptionsPage;
