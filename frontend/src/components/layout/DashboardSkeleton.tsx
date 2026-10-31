import React from 'react';

export const DashboardSkeleton: React.FC = () => {
  return (
    <div className="w-full space-y-6 animate-pulse" aria-busy="true" aria-label="Loading content">
      {/* Page Title & Actions Skeleton */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="space-y-2">
          <div className="h-4 w-32 bg-slate-200 rounded"></div>
          <div className="h-7 w-56 bg-slate-200 rounded-md"></div>
        </div>
        <div className="flex items-center gap-2.5">
          <div className="h-9 w-28 bg-slate-200 rounded-full"></div>
          <div className="h-9 w-32 bg-slate-200 rounded-full"></div>
        </div>
      </div>

      {/* Metric Cards Skeleton Grid (Stripe 4-up) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {[1, 2, 3, 4].map((i) => (
          <div
            key={i}
            className="p-5 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm space-y-3"
          >
            <div className="flex items-center justify-between">
              <div className="h-3.5 w-24 bg-slate-200 rounded"></div>
              <div className="h-4 w-4 bg-slate-200 rounded-full"></div>
            </div>
            <div className="h-8 w-36 bg-slate-200 rounded"></div>
            <div className="h-3 w-28 bg-slate-100 rounded"></div>
          </div>
        ))}
      </div>

      {/* Main Table / Content Skeleton */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm p-6 space-y-4">
        {/* Table Toolbar */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-2 border-b border-stripe-hairline">
          <div className="h-9 w-64 bg-slate-100 rounded-md"></div>
          <div className="flex items-center gap-2">
            <div className="h-9 w-24 bg-slate-100 rounded-full"></div>
            <div className="h-9 w-24 bg-slate-100 rounded-full"></div>
          </div>
        </div>

        {/* Table Rows Shimmer */}
        <div className="space-y-3 pt-2">
          {[1, 2, 3, 4, 5].map((row) => (
            <div
              key={row}
              className="flex items-center justify-between py-3 border-b border-slate-100 last:border-b-0"
            >
              <div className="flex items-center gap-3 w-1/3">
                <div className="w-8 h-8 rounded-full bg-slate-200 shrink-0"></div>
                <div className="space-y-1.5 w-full">
                  <div className="h-3.5 w-3/4 bg-slate-200 rounded"></div>
                  <div className="h-2.5 w-1/2 bg-slate-100 rounded"></div>
                </div>
              </div>
              <div className="h-4 w-20 bg-slate-200 rounded"></div>
              <div className="h-6 w-16 bg-slate-100 rounded-full"></div>
              <div className="h-4 w-24 bg-slate-200 rounded"></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default DashboardSkeleton;

