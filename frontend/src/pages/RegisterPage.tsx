import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { UserPlus, AlertCircle } from 'lucide-react';
import { InputField } from '../components/common/InputField';
import { Button } from '../components/common/Button';
import { PasswordRequirementsChecklist } from '../components/auth/PasswordRequirementsChecklist';
import { checkPasswordRequirements, isValidEmail } from '../utils/validations';
import { mockRegister } from '../services/authService';
import { useAuth } from '../hooks/useAuth';
import { useToast } from '../context/ToastContext';

export const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const { addToast } = useToast();

  // Controlled form state
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  });

  // Inline error state
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

      case 'password': {
        if (!value) {
          return 'Password is required';
        }
        const { allValid } = checkPasswordRequirements(value);
        if (!allValid) {
          return 'Password must satisfy all security requirements below';
        }
        return '';
      }

      case 'confirmPassword':
        if (!value) {
          return 'Please confirm your password';
        }
        if (value !== formData.password) {
          return 'Passwords do not match';
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

    // If already touched, validate live
    if (touched[name]) {
      const fieldError = validateField(name, value);
      setErrors((prev) => ({ ...prev, [name]: fieldError }));
    }

    // Live confirm password check if password changes
    if (name === 'password' && touched.confirmPassword && formData.confirmPassword) {
      if (formData.confirmPassword !== value) {
        setErrors((prev) => ({ ...prev, confirmPassword: 'Passwords do not match' }));
      } else {
        setErrors((prev) => ({ ...prev, confirmPassword: '' }));
      }
    }
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setTouched((prev) => ({ ...prev, [name]: true }));
    const fieldError = validateField(name, value);
    setErrors((prev) => ({ ...prev, [name]: fieldError }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setServerError(null);

    // Mark all as touched
    setTouched({
      email: true,
      password: true,
      confirmPassword: true,
    });

    const emailErr = validateField('email', formData.email);
    const passErr = validateField('password', formData.password);
    const confirmErr = validateField('confirmPassword', formData.confirmPassword);

    const newErrors = {
      email: emailErr,
      password: passErr,
      confirmPassword: confirmErr,
    };

    setErrors(newErrors);

    if (emailErr || passErr || confirmErr) {
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await mockRegister({
        email: formData.email,
        password: formData.password,
        name: formData.name,
      });

      // Update auth session
      login(response.user, response.token);

      addToast(
        'Account created successfully!',
        `Welcome to Brightloop, ${response.user.name}.`,
        'success'
      );

      // Redirect to dashboard as required
      navigate('/dashboard', { replace: true });
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Registration failed. Please try again.';
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
          Create your account
        </h2>
        <p className="text-sm text-stripe-ink-mute mt-1">
          Start managing subscriptions and billing seamlessly.
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

      {/* Registration Form */}
      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        {/* Optional Name field */}
        <InputField
          label="Full Name (optional)"
          id="name"
          name="name"
          type="text"
          value={formData.name}
          onChange={handleChange}
          placeholder="Jane Doe"
          autoComplete="name"
          disabled={isSubmitting}
        />

        {/* Email Field */}
        <InputField
          label="Work Email"
          id="email"
          name="email"
          type="email"
          value={formData.email}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.email ? errors.email : undefined}
          placeholder="you@company.com"
          autoComplete="email"
          required
          disabled={isSubmitting}
        />

        {/* Password Field */}
        <div>
          <InputField
            label="Password"
            id="password"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleChange}
            onBlur={handleBlur}
            error={touched.password ? errors.password : undefined}
            placeholder="Create a strong password"
            autoComplete="new-password"
            required
            disabled={isSubmitting}
          />

          {/* STRICT Real-Time Visual Feedback Checklist */}
          <PasswordRequirementsChecklist password={formData.password} />
        </div>

        {/* Confirm Password Field */}
        <InputField
          label="Confirm Password"
          id="confirmPassword"
          name="confirmPassword"
          type="password"
          value={formData.confirmPassword}
          onChange={handleChange}
          onBlur={handleBlur}
          error={touched.confirmPassword ? errors.confirmPassword : undefined}
          placeholder="Re-enter your password"
          autoComplete="new-password"
          required
          disabled={isSubmitting}
        />

        {/* Submit Button */}
        <div className="pt-2">
          <Button
            type="submit"
            variant="primary"
            size="md"
            className="w-full shadow-stripe-sm"
            isLoading={isSubmitting}
            loadingText="Creating account..."
            leftIcon={<UserPlus className="w-4 h-4" />}
          >
            Create account
          </Button>
        </div>
      </form>

      {/* Switch to Login Link */}
      <div className="mt-6 pt-5 border-t border-stripe-hairline text-center text-xs text-stripe-ink-secondary">
        <span>Already have an account? </span>
        <Link
          to="/login"
          className="text-stripe-primary font-medium hover:text-stripe-primary-deep transition-colors ml-1 inline-flex items-center"
        >
          Sign in &rarr;
        </Link>
      </div>
    </div>
  );
};

export default RegisterPage;

