import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { LogIn, AlertCircle, Sparkles } from 'lucide-react';
import { InputField } from '../components/common/InputField';
import { Button } from '../components/common/Button';
import { isValidEmail } from '../utils/validations';
import { mockLogin, DEFAULT_DEMO_USERS } from '../services/authService';
import { useAuth } from '../hooks/useAuth';
import { useToast } from '../context/ToastContext';

export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const { addToast } = useToast();

  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  const validateField = (field: string, value: string) => {
    switch (field) {
      case 'email':
        if (!value.trim()) {
          return 'Email is required';
        }
        if (!isValidEmail(value)) {
          return 'Please enter a valid email address';
        }
        return '';

      case 'password':
        if (!value) {
          return 'Password is required';
        }
        return '';

      default:
        return '';
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setServerError(null);

    if (touched[name]) {
      const fieldError = validateField(name, value);
      setErrors((prev) => ({ ...prev, [name]: fieldError }));
    }
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setTouched((prev) => ({ ...prev, [name]: true }));
    const fieldError = validateField(name, value);
    setErrors((prev) => ({ ...prev, [name]: fieldError }));
  };

  const handleQuickDemoFill = (email: string) => {
    const demo = DEFAULT_DEMO_USERS[email];
    if (demo) {
      setFormData({
        email: demo.user.email,
        password: demo.passwordHash,
      });
      setErrors({});
      setServerError(null);
      setTouched({ email: true, password: true });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setServerError(null);

    setTouched({
      email: true,
      password: true,
    });

    const emailErr = validateField('email', formData.email);
    const passErr = validateField('password', formData.password);

    setErrors({
      email: emailErr,
      password: passErr,
    });

    if (emailErr || passErr) {
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await mockLogin({
        email: formData.email,
        password: formData.password,
      });

      // Update auth session
      login(response.user, response.token);

      addToast(
        'Signed in successfully',
        `Welcome back, ${response.user.name}.`,
        'success'
      );

      // Determine destination: prior target from ProtectedRoute or default to /dashboard
      const fromPath = (location.state as { from?: { pathname?: string } })?.from?.pathname;
      const targetDestination = fromPath && fromPath !== '/login' ? fromPath : '/dashboard';

      navigate(targetDestination, { replace: true });
    } catch (err: unknown) {
      const message =
        err instanceof Error
          ? err.message
          : 'Failed to sign in. Please check your credentials.';
      setServerError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div>
      {/* Page Heading */}
      <div className="mb-6">
        <h2 className="text-2xl font-light text-stripe-ink tracking-tight">
          Sign in to Brightloop
        </h2>
        <p className="text-sm text-stripe-ink-mute mt-1">
          Access your billing dashboard and subscription controls.
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

      {/* Login Form */}
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {/* Email Field */}
        <InputField
          label="Email address"
          id="email"
          name="email"
          type="email"
          value={formData.email}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.email ? errors.email : undefined}
          placeholder="alex.admin@brightloop.io"
          autoComplete="email"
          required
          disabled={isSubmitting}
        />

        {/* Password Field with Forgot Link */}
        <InputField
          label="Password"
          id="password"
          name="password"
          type="password"
          value={formData.password}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.password ? errors.password : undefined}
          placeholder="Enter your password"
          autoComplete="current-password"
          required
          disabled={isSubmitting}
          cornerAction={
            <Link
              to="/forgot-password"
              className="text-stripe-primary hover:text-stripe-primary-deep text-xs font-medium transition-colors"
            >
              Forgot your password?
            </Link>
          }
        />

        {/* Submit Button */}
        <div className="pt-2">
          <Button
            type="submit"
            variant="primary"
            size="md"
            className="w-full shadow-stripe-sm"
            isLoading={isSubmitting}
            loadingText="Signing in..."
            leftIcon={<LogIn className="w-4 h-4" />}
          >
            Sign in
          </Button>
        </div>
      </form>

      {/* Quick Demo Credentials Helper */}
      <div className="mt-5 p-3 bg-slate-50/80 rounded-xl border border-stripe-hairline">
        <div className="flex items-center gap-1.5 text-[11px] font-semibold text-stripe-ink-mute uppercase tracking-wider mb-2">
          <Sparkles className="w-3 h-3 text-stripe-primary" />
          <span>Quick Demo Logins</span>
        </div>
        <div className="flex flex-wrap gap-1.5">
          <button
            type="button"
            onClick={() => handleQuickDemoFill('alex.admin@brightloop.io')}
            disabled={isSubmitting}
            className="px-2.5 py-1 text-xs bg-white hover:bg-slate-100 border border-stripe-hairline rounded-md text-stripe-ink-secondary hover:text-stripe-primary transition-colors"
          >
            Admin
          </button>
          <button
            type="button"
            onClick={() => handleQuickDemoFill('sarah.miller@brightloop.io')}
            disabled={isSubmitting}
            className="px-2.5 py-1 text-xs bg-white hover:bg-slate-100 border border-stripe-hairline rounded-md text-stripe-ink-secondary hover:text-stripe-primary transition-colors"
          >
            Staff
          </button>
          <button
            type="button"
            onClick={() => handleQuickDemoFill('david@acme-enterprises.com')}
            disabled={isSubmitting}
            className="px-2.5 py-1 text-xs bg-white hover:bg-slate-100 border border-stripe-hairline rounded-md text-stripe-ink-secondary hover:text-stripe-primary transition-colors"
          >
            Customer
          </button>
        </div>
      </div>

      {/* Switch to Registration Link */}
      <div className="mt-6 pt-5 border-t border-stripe-hairline text-center text-xs text-stripe-ink-secondary">
        <span>Don't have an account? </span>
        <Link
          to="/register"
          className="text-stripe-primary font-medium hover:text-stripe-primary-deep transition-colors ml-1 inline-flex items-center"
        >
          Sign up &rarr;
        </Link>
      </div>
    </div>
  );
};

export default LoginPage;

