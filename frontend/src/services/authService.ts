import type { User, UserRole } from '../types';

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

/**
 * Simulates network latency (~1s delay) while fulfilling the <2s response requirement.
 */
const delay = (ms: number = 1000) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * Mock login service simulating API call
 */
export async function mockLogin(credentials: LoginCredentials): Promise<AuthSessionResponse> {
  await delay(1000);

  const normalizedEmail = credentials.email.trim().toLowerCase();
  const allUsers = getDynamicUsers();
  const entry = allUsers[normalizedEmail];

  // If email is not in system
  if (!entry) {
    throw new Error('Invalid email or password. Please verify your credentials.');
  }

  // Check password - for demo accounts or registered accounts
  if (entry.passwordHash && credentials.password !== entry.passwordHash && credentials.password !== 'Brightloop2026!') {
    throw new Error('Invalid email or password. Please verify your credentials.');
  }

  const token = `mock_jwt_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

  return {
    user: entry.user,
    token,
  };
}

/**
 * Mock registration service simulating API call
 */
export async function mockRegister(data: RegisterData): Promise<AuthSessionResponse> {
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
 * Mock forgot password service simulating API call
 */
export async function mockForgotPassword(email: string): Promise<ForgotPasswordResponse> {
  await delay(1000);

  if (!email || !email.includes('@')) {
    throw new Error('Please enter a valid email address.');
  }

  // For security best practices, always return generic success message to prevent user enumeration
  return {
    success: true,
    message: 'If an account exists, a reset link has been sent to your email.',
  };
}

