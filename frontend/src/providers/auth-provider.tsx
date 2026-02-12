"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import type { User } from "@/types";
import { getCurrentUser, login as loginFn, logout as logoutFn, register as registerFn } from "@/lib/auth";
import type { LoginRequest, RegisterRequest } from "@/types";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    try {
      const userData = await getCurrentUser();
      setUser(userData);
    } catch {
      setUser(null);
      // Cookies são gerenciados pelo backend
    }
  }, []);

  useEffect(() => {
    // Tentar buscar usuário - cookie será enviado automaticamente
    refreshUser()
      .catch(() => {
        // Se falhar, usuário não está autenticado
      })
      .finally(() => setIsLoading(false));
  }, [refreshUser]);

  const login = async (data: LoginRequest) => {
    const response = await loginFn(data);
    setUser(response.user);
  };

  const register = async (data: RegisterRequest) => {
    const response = await registerFn(data);
    setUser(response.user);
  };

  const logout = async () => {
    await logoutFn();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
