"use client";

import { createContext, useContext, useEffect, useState, useCallback } from "react";
import { useRouter, usePathname } from "next/navigation";
import type { User, LoginRequest, RegisterRequest, SocialLoginRequest } from "@/types";
import { getCurrentUser, login as loginFn, logout as logoutFn, register as registerFn, socialLogin as socialLoginFn } from "@/lib/auth";
import { setRedirectCallback } from "@/lib/api";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  socialLogin: (data: SocialLoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    setRedirectCallback(() => {
      // Rotas públicas que não devem redirecionar
      const publicRoutes = ["/", "/login", "/register"];
      const isPublicRoute = publicRoutes.some(route => pathname === route || pathname?.startsWith(route));

      if (!isPublicRoute) {
        router.push("/login");
      }
    });
  }, [router, pathname]);

  const refreshUser = useCallback(async () => {
    try {
      const userData = await getCurrentUser();
      setUser(userData);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    let mounted = true;

    const loadUser = async () => {
      try {
        if (mounted) {
          await refreshUser();
        }
      } catch {
        // Se falhar, usuário não está autenticado
      } finally {
        if (mounted) {
          setIsLoading(false);
        }
      }
    };

    loadUser();

    return () => {
      mounted = false;
    };
  }, [refreshUser]);

  const login = async (data: LoginRequest) => {
    const response = await loginFn(data);
    setUser(response.user);
  };

  const register = async (data: RegisterRequest) => {
    const response = await registerFn(data);
    setUser(response.user);
  };

  const socialLogin = async (data: SocialLoginRequest) => {
    const response = await socialLoginFn(data);
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
        socialLogin,
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
