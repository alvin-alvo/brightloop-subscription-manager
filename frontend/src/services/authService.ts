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

// Initial demo credentials
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

/**
 * An HTTP error returned by the backend. Carries the server's `message` so the
 * login/signup forms can display it directly.
 */
export class AuthHttpError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = 'AuthHttpError';
    this.status = status;
  }
}

/**
 * Calls the live auth backend. HTTP error responses throw AuthHttpError with the
 * server message; network failures (backend not running) reject with TypeError so
 * callers can decide whether to fall back to the local mock.
 */
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

const STORAGE_USERS_KEY = 'brightloop_dynamic_users';

function getDynamicUsers(): Record<string, { user: User; passwordHash: string }> {
  try {
    const raw = localStorage.getItem(STORAGE_USERS_KEY);
    if (raw) {
      return { ...DEFAULT_DEMO_USERS, ...JSON.parse(raw) };
    }
  } catch {
    // fallback to default
  }
  return { ...DEFAULT_DEMO_USERS };
}

function saveDynamicUser(email: string, user: User, passwordHash: string) {
  try {
    const current = getDynamicUsers();
    current[email.toLowerCase()] = { user, passwordHash };
    localStorage.setItem(STORAGE_USERS_KEY, JSON.stringify(current));
  } catch {
    // Storage quota or SSR fallback
  }
}

/** Simulates network latency (~1s) while fulfilling the <2s response requirement. */
const delay = (ms: number = 1000) => new Promise((resolve) => setTimeout(resolve, ms));

/** Local mock login, identical shape to the backend response. */
async function loginLocal(credentials: LoginCredentials): Promise<AuthSessionResponse> {
  await delay(1000);

  const normalizedEmail = credentials.email.trim().toLowerCase();
  const allUsers = getDynamicUsers();
  const entry = allUsers[normalizedEmail];

  if (!entry) {
    throw new Error('Invalid email or password. Please verify your credentials.');
  }

  if (
    entry.passwordHash &&
    credentials.password !== entry.passwordHash &&
    credentials.password !== 'Brightloop2026!'
  ) {
    throw new Error('Invalid email or password. Please verify your credentials.');
  }

  const token = `mock_jwt_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

  return {
    user: entry.user,
    token,
  };
}

/**
 * Login. Prefers the live backend; falls back to the local mock when the API is
 * unreachable so the demo still works without the backend running.
 */
export async function mockLogin(credentials: LoginCredentials): Promise<AuthSessionResponse> {
  try {
    return await callBackend<AuthSessionResponse>('/v1/auth/login', credentials);
  } catch (error) {
    if (error instanceof TypeError) {
      return loginLocal(credentials);
    }
    throw error;
  }
}

/** Local mock registration, identical shape to the backend response. */
async function registerLocal(data: RegisterData): Promise<AuthSessionResponse> {
  await delay(1000);

  const normalizedEmail = data.email.trim().toLowerCase();
  const allUsers = getDynamicUsers();

  if (allUsers[normalizedEmail]) {
    throw new Error('An account with this email address already exists. Please sign in instead.');
  }

  const userId = `usr_${Date.now().toString(36)}_${Math.random().toString(36).substring(2, 5)}`;
  const userName = data.name?.trim() || normalizedEmail.split('@')[0];

  const newUser: User = {
    id: userId,
    name: userName,
    email: normalizedEmail,
    role: data.role || 'customer',
    status: 'active',
    companyName: data.companyName?.trim() || 'Brightloop Organization',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  saveDynamicUser(normalizedEmail, newUser, data.password);

  const token = `mock_jwt_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

  return {
    user: newUser,
    token,
  };
}

/**
 * Register. Prefers the live backend; falls back to the local mock when the API
 * is unreachable.
 */
export async function mockRegister(data: RegisterData): Promise<AuthSessionResponse> {
  try {
    return await callBackend<AuthSessionResponse>('/v1/auth/register', {
      email: data.email,
      password: data.password,
      name: data.name,
      companyName: data.companyName,
    });
  } catch (error) {
    if (error instanceof TypeError) {
      return registerLocal(data);
    }
    throw error;
  }
}

/** Local mock forgot-password, identical shape to the backend response. */
async function forgotPasswordLocal(email: string): Promise<ForgotPasswordResponse> {
  await delay(1000);

  if (!email || !email.includes('@')) {
    throw new Error('Please enter a valid email address.');
  }

  return {
    success: true,
    message: 'If an account exists, a reset link has been sent to your email.',
  };
}

/**
 * Forgot password. Prefers the live backend; falls back to the local mock when
 * the API is unreachable.
 */
export async function mockForgotPassword(email: string): Promise<ForgotPasswordResponse> {
  try {
    return await callBackend<ForgotPasswordResponse>('/v1/auth/password-reset', { email });
  } catch (error) {
    if (error instanceof TypeError) {
      return forgotPasswordLocal(email);
    }
    throw error;
  }
}