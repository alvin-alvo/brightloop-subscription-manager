import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Mail, ArrowLeft, CheckCircle2, AlertCircle } from 'lucide-react';
import { InputField } from '../components/common/InputField';
import { Button } from '../components/common/Button';
import { isValidEmail } from '../utils/validations';
import { mockForgotPassword } from '../services/authService';
import { useToast } from '../context/ToastContext';

export const ForgotPasswordPage: React.FC = () => {
  const { addToast } = useToast();

  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [touched, setTouched] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  const validate = (value: string) => {
    if (!value.trim()) {
      return 'Email address is required';
    }
    if (!isValidEmail(value)) {
      return 'Please enter a valid email address';
    }
    return '';
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    setServerError(null);

    if (touched) {
      setError(validate(e.target.value));
    }
  };

  const handleBlur = () => {
    setTouched(true);
    setError(validate(email));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setServerError(null);
    setTouched(true);

    const validationError = validate(email);
    setError(validationError);

    if (validationError) {
      return;
    }

    setIsSubmitting(true);

    try {
      await mockForgotPassword(email);

      // Display required success toast
      addToast(
        'Password Reset Sent',
        'If an account exists, a reset link has been sent to your email.',
        'success'
      );

      setIsSuccess(true);
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Unable to process reset request.';
      setServerError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div>
      {/* Return to Sign In link */}
      <div className="mb-4">
        <Link
          to="/login"
          className="inline-flex items-center gap-1.5 text-xs font-medium text-stripe-ink-mute hover:text-stripe-ink transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Back to sign in</span>
        </Link>
      </div>

      {/* Page Heading */}
      <div className="mb-6">
        <h2 className="text-2xl font-light text-stripe-ink tracking-tight">
          Reset your password
        </h2>
        <p className="text-sm text-stripe-ink-mute mt-1">
          Enter your email address and we'll send you a password recovery link.
        </p>
      </div>

      {/* Global Server Error Banner */}
      {serverError && (
        <div
          className="mb-5 p-3.5 bg-rose-50 border border-rose-200 rounded-lg flex items-start gap-2.5 text-xs text-stripe-ruby animate-fadeIn"
          role="alert"
        >
          <AlertCircle className="w-4 h-4 shrink-0 mt-0.5" />
          <span className="leading-relaxed font-medium">{serverError}</span>
        </div>
      )}

      {isSuccess ? (
        /* Success State View */
        <div className="space-y-5 animate-fadeIn">
          <div className="p-4 bg-emerald-50/80 border border-emerald-200 rounded-xl flex items-start gap-3">
            <CheckCircle2 className="w-5 h-5 text-emerald-600 shrink-0 mt-0.5" />
            <div className="text-xs text-emerald-950 space-y-1">
              <p className="font-semibold text-sm text-emerald-900">
                Check your inbox
              </p>
              <p className="leading-relaxed">
                If an account exists, a reset link has been sent to your email.
              </p>
              <p className="text-[11px] text-emerald-700 pt-1 font-mono">
                {email}
              </p>
            </div>
          </div>

          <div className="pt-2 space-y-2.5">
            <Button
              type="button"
              variant="secondary"
              size="md"
              className="w-full"
              onClick={() => {
                setIsSuccess(false);
                setEmail('');
                setTouched(false);
                setError('');
              }}
            >
              Send another link
            </Button>

            <Link to="/login" className="block w-full">
              <Button
                type="button"
                variant="primary"
                size="md"
                className="w-full"
              >
                Return to sign in
              </Button>
            </Link>
          </div>
        </div>
      ) : (
        /* Forgot Password Form */
        <form onSubmit={handleSubmit} noValidate className="space-y-4">
          <InputField
            label="Account Email"
            id="email"
            name="email"
            type="email"
            value={email}
            onChange={handleChange}
            onBlur={handleBlur}
            error={touched ? error : undefined}
            placeholder="you@company.com"
            autoComplete="email"
            required
            disabled={isSubmitting}
          />

          <div className="pt-2">
            <Button
              type="submit"
              variant="primary"
              size="md"
              className="w-full shadow-stripe-sm"
              isLoading={isSubmitting}
              loadingText="Sending reset link..."
              leftIcon={<Mail className="w-4 h-4" />}
            >
              Send reset link
            </Button>
          </div>
        </form>
      )}

      {/* Footer Navigation */}
      <div className="mt-6 pt-5 border-t border-stripe-hairline text-center text-xs text-stripe-ink-secondary">
        <span>Need further assistance? </span>
        <a
          href="mailto:support@brightloop.io"
          className="text-stripe-primary font-medium hover:text-stripe-primary-deep transition-colors"
        >
          Contact Support
        </a>
      </div>
    </div>
  );
};

export default ForgotPasswordPage;
