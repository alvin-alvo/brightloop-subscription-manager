import React from 'react';
import { StatusBadge } from '../components/StatusBadge';

export const HomePage: React.FC = () => {
  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-center items-center p-6">
      <div className="max-w-md w-full bg-white rounded-xl shadow-sm border border-slate-200 p-8 text-center">
        <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-indigo-50 text-indigo-600 mb-4 font-bold text-xl">
          BL
        </div>
        <h1 className="text-2xl font-bold text-slate-900 mb-2">Brightloop</h1>
        <p className="text-sm text-slate-600 mb-6">
          B2B SaaS Subscription Management System
        </p>
        <div className="flex justify-center gap-2 mb-6">
          <StatusBadge label="Vite + React" variant="info" />
          <StatusBadge label="TypeScript" variant="success" />
          <StatusBadge label="Tailwind CSS" variant="neutral" />
        </div>
        <p className="text-xs text-slate-500">
          Frontend initialized. Ready for Authentication flow.
        </p>
      </div>
    </div>
  );
};
