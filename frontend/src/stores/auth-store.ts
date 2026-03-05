import {
  getCurrentUser,
  login as loginFn,
  logout as logoutFn,
  register as registerFn,
  socialLogin as socialLoginFn,
} from "@/lib/auth";
import { getAccessToken } from "@/lib/api";
import type {
  LoginRequest,
  RegisterRequest,
  SocialLoginRequest,
  User,
} from "@/types";
import { create } from "zustand";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  initialize: () => Promise<void>;
  login: (data: LoginRequest) => Promise<User>;
  register: (data: RegisterRequest) => Promise<User>;
  socialLogin: (data: SocialLoginRequest) => Promise<User>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()((set) => ({
  user: null,
  isLoading: true,
  isAuthenticated: false,

  initialize: async () => {
    if (!getAccessToken()) {
      set({ user: null, isAuthenticated: false, isLoading: false });
      return;
    }
    try {
      const user = await getCurrentUser();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch {
      set({ user: null, isAuthenticated: false, isLoading: false });
    }
  },

  login: async (data) => {
    const response = await loginFn(data);
    set({ user: response.user, isAuthenticated: true });
    return response.user;
  },

  register: async (data) => {
    const response = await registerFn(data);
    set({ user: response.user, isAuthenticated: true });
    return response.user;
  },

  socialLogin: async (data) => {
    const response = await socialLoginFn(data);
    set({ user: response.user, isAuthenticated: true });
    return response.user;
  },

  logout: async () => {
    await logoutFn();
    set({ user: null, isAuthenticated: false });
  },

  refreshUser: async () => {
    try {
      const user = await getCurrentUser();
      set({ user, isAuthenticated: true });
    } catch {
      set({ user: null, isAuthenticated: false });
    }
  },
}));
