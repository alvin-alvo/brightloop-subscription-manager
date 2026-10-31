import { useState, useCallback } from 'react';
import type { User } from '../types';

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem('brightloop_token')
  );

  const logout = useCallback(() => {
    localStorage.removeItem('brightloop_token');
    setUser(null);
    setToken(null);
  }, []);

  return {
    user,
    token,
    isAuthenticated: Boolean(token),
    setUser,
    setToken,
    logout,
  };
}
