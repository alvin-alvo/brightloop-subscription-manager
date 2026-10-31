import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import type { User, UserRole } from '../types';

interface AuthContextType {
  user: User | null;
  role: UserRole;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
  switchRole: (role: UserRole) => void;
}

export const DEMO_USERS: Record<UserRole, User> = {
  admin: {
    id: 'usr_admin_01',
    name: 'Alex Vance',
    email: 'alex.admin@brightloop.io',
    role: 'admin',
    status: 'active',
    companyName: 'Brightloop HQ',
    createdAt: '2024-01-15T08:00:00Z',
    updatedAt: '2024-10-31T09:00:00Z',
  },
  staff: {
    id: 'usr_staff_02',
    name: 'Sarah Miller',
    email: 'sarah.miller@brightloop.io',
    role: 'staff',
    status: 'active',
    companyName: 'Brightloop Operations',
    createdAt: '2024-02-01T10:00:00Z',
    updatedAt: '2024-10-31T09:00:00Z',
  },
  customer: {
    id: 'usr_cust_03',
    name: 'David Chen',
    email: 'david@acme-enterprises.com',
    role: 'customer',
    status: 'active',
    companyName: 'Acme Enterprises Inc.',
    createdAt: '2024-03-12T14:30:00Z',
    updatedAt: '2024-10-31T09:00:00Z',
  },
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const STORAGE_KEY_USER = 'brightloop_session_user';
const STORAGE_KEY_TOKEN = 'brightloop_token';

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY_USER);
      if (saved) {
        return JSON.parse(saved);
      }
    } catch {
      // Fallback
    }
    return null;
  });

  const [token, setToken] = useState<string | null>(() => {
    return localStorage.getItem(STORAGE_KEY_TOKEN) || null;
  });

  useEffect(() => {
    if (user) {
      localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(user));
    } else {
      localStorage.removeItem(STORAGE_KEY_USER);
    }
  }, [user]);

  const login = useCallback((newUser: User, newToken: string) => {
    setUser(newUser);
    setToken(newToken);
    localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(newUser));
    localStorage.setItem(STORAGE_KEY_TOKEN, newToken);
  }, []);

  const logout = useCallback(() => {
    setUser(null);
    setToken(null);
    localStorage.removeItem(STORAGE_KEY_USER);
    localStorage.removeItem(STORAGE_KEY_TOKEN);
  }, []);

  const switchRole = useCallback((newRole: UserRole) => {
    const demoUser = DEMO_USERS[newRole];
    setUser(demoUser);
    setToken(`mock_token_${newRole}`);
    localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(demoUser));
    localStorage.setItem(STORAGE_KEY_TOKEN, `mock_token_${newRole}`);
  }, []);

  const role: UserRole = user?.role || 'customer';
  const isAuthenticated = Boolean(user && token);

  return (
    <AuthContext.Provider
      value={{
        user,
        role,
        isAuthenticated,
        login,
        logout,
        switchRole,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuthContext(): AuthContextType {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used within an AuthProvider');
  }
  return context;
}

export default AuthContext;
