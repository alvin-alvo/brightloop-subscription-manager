import type { User, UserRole } from '../types';
import { BASE_URL } from './api';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  name?: string;
  companyName?: string;
  role?: UserRole;
}

export interface AuthSessionResponse {
  user: User;
  token: string;
}

export interface ForgotPasswordResponse {
  success: boolean;
  message: string;
}

// Demo credentials shown on the login page's "Quick Demo Logins" helper.
// These are real accounts seeded in the backend database.
export const DEFAULT_DEMO_USERS: Record<string, { user: User; passwordHash: string }> = {
  'alex.admin@brightloop.io': {
    user: {
      id: 'usr_admin_01',
      name: 'Alex Vance',
      email: 'alex.admin@brightloop.io',
      role: 'admin',
      status: 'active',
      companyName: 'Brightloop HQ',
      createdAt: '2024-01-15T08:00:00Z',
      updatedAt: '2024-10-31T09:00:00Z',
    },
    passwordHash: 'Brightloop2026!',
  },
  'sarah.miller@brightloop.io': {
    user: {
      id: 'usr_staff_02',
      name: 'Sarah Miller',
      email: 'sarah.miller@brightloop.io',
      role: 'staff',
      status: 'active',
      companyName: 'Brightloop Operations',
      createdAt: '2024-02-01T10:00:00Z',
      updatedAt: '2024-10-31T09:00:00Z',
    },
    passwordHash: 'Brightloop2026!',
  },
  'david@acme-enterprises.com': {
    user: {
      id: 'usr_cust_03',
      name: 'David Chen',
      email: 'david@acme-enterprises.com',
      role: 'customer',
      status: 'active',
      companyName: 'Acme Enterprises Inc.',
      createdAt: '2024-03-12T14:30:00Z',
      updatedAt: '2024-10-31T09:00:00Z',
    },
    passwordHash: 'Brightloop2026!',
  },
};

/** An HTTP error returned by the backend, carrying the server message. */
export class AuthHttpError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = 'AuthHttpError';
    this.status = status;
  }
}

/** Calls the live auth backend. Errors surface the server's message. */
async function callBackend<T>(endpoint: string, body: unknown): Promise<T> {
  const response = await fetch(`${BASE_URL}${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    const data = await response.json().catch(() => ({}));
    const message =
      (data as { message?: string })?.message || `Request failed with status ${response.status}`;
    throw new AuthHttpError(message, response.status);
  }

  return response.json() as Promise<T>;
}

/** Login against the real backend. */
export async function mockLogin(credentials: LoginCredentials): Promise<AuthSessionResponse> {
  return callBackend<AuthSessionResponse>('/v1/auth/login', credentials);
}

/** Register a customer account against the real backend. */
export async function mockRegister(data: RegisterData): Promise<AuthSessionResponse> {
  return callBackend<AuthSessionResponse>('/v1/auth/register', {
    email: data.email,
    password: data.password,
    name: data.name,
    companyName: data.companyName,
  });
}

/** Send a password-reset request to the real backend. */
export async function mockForgotPassword(email: string): Promise<ForgotPasswordResponse> {
  return callBackend<ForgotPasswordResponse>('/v1/auth/password-reset', { email });
}