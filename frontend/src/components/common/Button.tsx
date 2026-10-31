import React from 'react';
import { Loader2 } from 'lucide-react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  loadingText?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'md',
  isLoading = false,
  loadingText,
  leftIcon,
  rightIcon,
  disabled,
  className = '',
  ...props
}) => {
  const baseClasses =
    'inline-flex items-center justify-center font-medium transition-all duration-150 focus:outline-none focus:ring-2 focus:ring-offset-1 disabled:opacity-60 disabled:cursor-not-allowed select-none';

  // Stripe pill button geometry
  const sizeClasses = {
    sm: 'text-xs py-1.5 px-3.5 rounded-full gap-1.5 min-h-[32px]',
    md: 'text-sm py-2.5 px-5 rounded-full gap-2 min-h-[40px]',
    lg: 'text-base py-3 px-6 rounded-full gap-2.5 min-h-[46px]',
  };

  // Stripe color tokens from DESIGN.md
  const variantClasses = {
    primary:
      'bg-stripe-primary hover:bg-stripe-primary-deep active:bg-stripe-primary-press text-white shadow-stripe-sm focus:ring-stripe-primary/40',
    secondary:
      'bg-white text-stripe-primary border border-stripe-hairline hover:border-stripe-hairline-input hover:bg-slate-50 active:bg-slate-100 shadow-stripe-sm focus:ring-stripe-primary/30',
    ghost:
      'bg-transparent text-stripe-ink-secondary hover:text-stripe-ink hover:bg-slate-100/60 active:bg-slate-200/60 focus:ring-slate-300',
  };

  const isDisabled = disabled || isLoading;

  return (
    <button
      disabled={isDisabled}
      className={`${baseClasses} ${sizeClasses[size]} ${variantClasses[variant]} ${className}`}
      {...props}
    >
      {isLoading ? (
        <>
          <Loader2 className="w-4 h-4 animate-spin shrink-0" />
          <span>{loadingText || children}</span>
        </>
      ) : (
        <>
          {leftIcon && <span className="shrink-0">{leftIcon}</span>}
          <span>{children}</span>
          {rightIcon && <span className="shrink-0">{rightIcon}</span>}
        </>
      )}
    </button>
  );
};

export default Button;

