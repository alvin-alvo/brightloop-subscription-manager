import React from 'react';
import { Check, Circle } from 'lucide-react';
import { checkPasswordRequirements } from '../../utils/validations';

interface PasswordRequirementsChecklistProps {
  password: string;
}

export const PasswordRequirementsChecklist: React.FC<PasswordRequirementsChecklistProps> = ({
  password,
}) => {
  const { requirements, score } = checkPasswordRequirements(password);

  // Strength label & color calculations
  const getStrengthMeta = (currentScore: number) => {
    if (password.length === 0) {
      return { label: '', color: 'bg-slate-200', textClass: 'text-slate-400' };
    }
    switch (currentScore) {
      case 1:
        return { label: 'Weak', color: 'bg-stripe-ruby', textClass: 'text-stripe-ruby' };
      case 2:
        return { label: 'Fair', color: 'bg-amber-500', textClass: 'text-amber-600' };
      case 3:
        return { label: 'Good', color: 'bg-indigo-500', textClass: 'text-indigo-600' };
      case 4:
        return { label: 'Strong', color: 'bg-emerald-600', textClass: 'text-emerald-600' };
      default:
        return { label: 'Weak', color: 'bg-stripe-ruby', textClass: 'text-stripe-ruby' };
    }
  };

  const strength = getStrengthMeta(score);

  return (
    <div className="mt-3 p-3.5 bg-slate-50/80 rounded-lg border border-stripe-hairline space-y-3">
      {/* Strength indicator bar */}
      <div className="space-y-1.5">
        <div className="flex items-center justify-between text-xs">
          <span className="text-stripe-ink-secondary font-medium">Password strength</span>
          {strength.label && (
            <span className={`font-semibold tracking-tight ${strength.textClass}`}>
              {strength.label}
            </span>
          )}
        </div>
        <div className="grid grid-cols-4 gap-1.5 h-1.5">
          {[1, 2, 3, 4].map((step) => {
            const isFilled = password.length > 0 && score >= step;
            return (
              <div
                key={step}
                className={`h-full rounded-full transition-all duration-300 ${
                  isFilled ? strength.color : 'bg-slate-200'
                }`}
              />
            );
          })}
        </div>
      </div>

      {/* Real-time Checklist */}
      <div className="space-y-2 pt-1 border-t border-slate-200/60">
        <p className="text-[11px] font-medium uppercase tracking-wider text-stripe-ink-mute">
          Requirements
        </p>
        <ul className="space-y-1.5" aria-label="Password requirements">
          {requirements.map((req) => (
            <li
              key={req.id}
              className={`flex items-center gap-2 text-xs transition-colors duration-200 ${
                req.isValid ? 'text-emerald-700 font-medium' : 'text-slate-500'
              }`}
            >
              <div
                className={`w-4 h-4 rounded-full flex items-center justify-center shrink-0 transition-colors duration-200 ${
                  req.isValid
                    ? 'bg-emerald-100 text-emerald-600'
                    : 'bg-slate-200/70 text-slate-400'
                }`}
              >
                {req.isValid ? (
                  <Check className="w-2.5 h-2.5 stroke-[3]" />
                ) : (
                  <Circle className="w-1.5 h-1.5 fill-current" />
                )}
              </div>
              <span className={req.isValid ? 'text-stripe-ink font-normal' : ''}>
                {req.label}
              </span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default PasswordRequirementsChecklist;

