'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api, User, AuthOutput } from '@/lib/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    if (storedToken && storedUser) {
      try {
        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    }
    setLoading(false);
  }, []);

  const saveAuth = (out: AuthOutput) => {
    localStorage.setItem('token', out.token);
    localStorage.setItem('user', JSON.stringify(out.user));
    setToken(out.token);
    setUser(out.user);
  };

  const login = async (email: string, password: string) => {
    const res = await api.post<{ data: AuthOutput }>('/api/v1/auth/login', { email, password });
    saveAuth(res.data.data);
  };

  const register = async (name: string, email: string, password: string) => {
    const res = await api.post<{ data: AuthOutput }>('/api/v1/auth/register', { name, email, password });
    saveAuth(res.data.data);
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
