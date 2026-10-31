import { z } from 'zod';

export const passwordValidation = z
  .string()
  .refine((val) => val.length > 8, {
    message: 'Password must have more than 8 characters',
  })
  .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
  .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
  .regex(
    /[!@#$%^&*(),.?":{}|<>_\-+=[\]\\\/~`]/,
    'Password must contain at least one special character'
  );

export const registerSchema = z
  .object({
    name: z
      .string()
      .trim()
      .min(2, 'Name must be at least 2 characters long')
      .max(100, 'Name must not exceed 100 characters')
      .optional()
      .or(z.literal('')),
    email: z
      .string()
      .trim()
      .min(1, 'Email is required')
      .email('Please enter a valid email address')
      .toLowerCase(),
    companyName: z
      .string()
      .trim()
      .min(2, 'Company name must be at least 2 characters')
      .optional()
      .or(z.literal('')),
    role: z.enum(['admin', 'staff', 'customer']).default('customer'),
    password: passwordValidation,
    confirmPassword: z.string().min(1, 'Please confirm your password'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

export type RegisterFormData = z.infer<typeof registerSchema>;

export const loginSchema = z.object({
  email: z
    .string()
    .trim()
    .min(1, 'Email is required')
    .email('Please enter a valid email address')
    .toLowerCase(),
  password: z.string().min(1, 'Password is required'),
});

export type LoginFormData = z.infer<typeof loginSchema>;

export const forgotPasswordSchema = z.object({
  email: z
    .string()
    .trim()
    .min(1, 'Email is required')
    .email('Please enter a valid email address')
    .toLowerCase(),
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export interface PasswordRequirementStatus {
  id: string;
  label: string;
  isValid: boolean;
}

export interface PasswordValidationResult {
  requirements: PasswordRequirementStatus[];
  allValid: boolean;
  score: number;
}

/**
 * STRICT Real-time password requirement evaluator
 * Checks:
 * 1. More than 8 characters (length > 8)
 * 2. At least one uppercase letter ([A-Z])
 * 3. At least one lowercase letter ([a-z])
 * 4. At least one special character
 */
export function checkPasswordRequirements(password: string): PasswordValidationResult {
  const isMoreThan8 = password.length > 8;
  const hasUppercase = /[A-Z]/.test(password);
  const hasLowercase = /[a-z]/.test(password);
  const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>_\-+=[\]\\\/~`]/.test(password);

  const requirements: PasswordRequirementStatus[] = [
    {
      id: 'length',
      label: 'More than 8 characters',
      isValid: isMoreThan8,
    },
    {
      id: 'uppercase',
      label: 'At least one uppercase letter (A-Z)',
      isValid: hasUppercase,
    },
    {
      id: 'lowercase',
      label: 'At least one lowercase letter (a-z)',
      isValid: hasLowercase,
    },
    {
      id: 'special',
      label: 'At least one special character (!@#$%^&*)',
      isValid: hasSpecialChar,
    },
  ];

  const score = requirements.filter((r) => r.isValid).length;
  const allValid = requirements.every((r) => r.isValid);

  return {
    requirements,
    allValid,
    score,
  };
}

export function isValidEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim());
}
