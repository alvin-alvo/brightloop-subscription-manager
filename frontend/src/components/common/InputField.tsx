import React, { useState } from 'react';
import { Eye, EyeOff, AlertCircle } from 'lucide-react';

export interface InputFieldProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  helperText?: string;
  cornerAction?: React.ReactNode;
}

export const InputField: React.FC<InputFieldProps> = ({
  label,
  id,
  name,
  type = 'text',
  value,
  onChange,
  onBlur,
  error,
  helperText,
  placeholder,
  disabled,
  required,
  autoComplete,
  cornerAction,
  className = '',
  ...props
}) => {
  const [showPassword, setShowPassword] = useState(false);
  const inputId = id || name || `input_${label.toLowerCase().replace(/\s+/g, '_')}`;
  const isPassword = type === 'password';
  const effectiveType = isPassword ? (showPassword ? 'text' : 'password') : type;

  return (
    <div className="w-full">
      <div className="flex items-center justify-between mb-1.5">
        <label
          htmlFor={inputId}
          className="block text-xs font-medium text-stripe-ink-secondary tracking-normal"
        >
          {label}
          {required && <span className="text-stripe-ruby ml-0.5">*</span>}
        </label>
        {cornerAction && <div className="text-xs">{cornerAction}</div>}
      </div>

      <div className="relative flex items-center">
        <input
          id={inputId}
          name={name}
          type={effectiveType}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          disabled={disabled}
          placeholder={placeholder}
          autoComplete={autoComplete}
          aria-invalid={Boolean(error)}
          aria-describedby={error ? `${inputId}-error` : undefined}
          className={`w-full bg-white text-stripe-ink text-sm rounded-md px-3.5 py-2.5 transition-all duration-150 placeholder:text-slate-400 outline-none ${
            isPassword ? 'pr-10' : ''
          } ${
            error
              ? 'border border-stripe-ruby focus:border-stripe-ruby focus:ring-2 focus:ring-stripe-ruby/20'
              : 'border border-stripe-hairline-input focus:border-stripe-primary focus:ring-2 focus:ring-stripe-primary/20 hover:border-slate-400'
          } ${disabled ? 'bg-slate-50 text-slate-400 cursor-not-allowed' : ''} ${className}`}
          {...props}
        />

        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword((prev) => !prev)}
            disabled={disabled}
            tabIndex={-1}
            className="absolute right-2.5 top-1/2 -translate-y-1/2 p-1 text-slate-400 hover:text-stripe-ink transition-colors rounded focus:outline-none"
            aria-label={showPassword ? 'Hide password' : 'Show password'}
          >
            {showPassword ? (
              <EyeOff className="w-4 h-4" />
            ) : (
              <Eye className="w-4 h-4" />
            )}
          </button>
        )}
      </div>

      {error ? (
        <p
          id={`${inputId}-error`}
          className="mt-1.5 flex items-center gap-1.5 text-xs text-stripe-ruby font-normal animate-fadeIn"
          role="alert"
        >
          <AlertCircle className="w-3.5 h-3.5 shrink-0" />
          <span>{error}</span>
        </p>
      ) : helperText ? (
        <p className="mt-1.5 text-xs text-stripe-ink-mute">{helperText}</p>
      ) : null}
    </div>
  );
};

export default InputField;

