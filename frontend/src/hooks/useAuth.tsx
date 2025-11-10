import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import * as authApi from '@/services/auth';

type User = {
  id: string;
  name: string;
  email: string;
  role?: 'user' | 'staff' | 'admin';
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Restore session: prefer persisted user; else try token -> profile
    const raw = localStorage.getItem('auth_user');
    if (raw) {
      setUser(JSON.parse(raw));
      setLoading(false);
      return;
    }
    const token = localStorage.getItem('AUTH_TOKEN');
    if (token) {
      (async () => {
        try {
          const prof = await authApi.getProfile();
          setUser({ id: prof.id, name: prof.name, email: prof.email, role: prof.role });
        } catch (e) {
          localStorage.removeItem('AUTH_TOKEN');
        } finally {
          setLoading(false);
        }
      })();
    } else {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user) localStorage.setItem('auth_user', JSON.stringify(user));
    else localStorage.removeItem('auth_user');
  }, [user]);

  const value = useMemo(() => ({ user, loading,
    login: async (email: string, password: string) => {
      setLoading(true);
      try {
        const res = await authApi.login({ email, password });
        setUser(res.user ?? { id: 'me', name: 'User', email });
      } finally {
        setLoading(false);
      }
    },
    register: async (name: string, email: string, password: string) => {
      setLoading(true);
      try {
        const res = await authApi.register({ name, email, password });
        setUser(res.user ?? { id: 'me', name, email });
      } finally {
        setLoading(false);
      }
    },
    logout: async () => {
      try { await authApi.logout(); } catch {}
      localStorage.removeItem('AUTH_TOKEN');
      setUser(null);
    },
    refreshUser: async () => {
      try {
        const prof = await authApi.getProfile();
        setUser({ id: prof.id, name: prof.name, email: prof.email, role: prof.role });
      } catch (e) {
        console.error('Failed to refresh user profile:', e);
      }
    }
  }), [user, loading]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

