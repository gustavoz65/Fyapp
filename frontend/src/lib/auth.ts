import { api, setTokens, clearTokens } from "./api";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  SocialLoginRequest,
  ListProvidersResponse,
  SetPasswordRequest,
  User,
  UserSettings,
  UpdateUserRequest,
  UpdateUserSettingsRequest,
  ChangePasswordRequest,
} from "@/types";

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/auth/login", data);
  setTokens(response.access_token, response.refresh_token || undefined);
  return response;
}

export async function register(data: RegisterRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/auth/register", data);
  setTokens(response.access_token, response.refresh_token || undefined);
  return response;
}

export async function logout(): Promise<void> {
  try {
    await api.post("/auth/logout", {});
  } catch (error) {
    console.error("Logout error:", error);
  }
  clearTokens();
}

export async function getCurrentUser(): Promise<User> {
  return api.get<User>("/users/me");
}

export async function updateUser(data: UpdateUserRequest): Promise<User> {
  return api.put<User>("/users/me", data);
}

export async function getUserSettings(): Promise<UserSettings> {
  return api.get<UserSettings>("/users/settings");
}

export async function updateUserSettings(data: UpdateUserSettingsRequest): Promise<UserSettings> {
  return api.put<UserSettings>("/users/settings", data);
}

export async function changePassword(data: ChangePasswordRequest): Promise<void> {
  return api.post("/auth/change-password", data);
}

export async function deactivateAccount(): Promise<void> {
  return api.delete("/users/me");
}

export async function socialLogin(data: SocialLoginRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/auth/social/login", data);
  setTokens(response.access_token, response.refresh_token || undefined);
  return response;
}

export async function setPassword(data: SetPasswordRequest): Promise<void> {
  return api.post("/auth/set-password", data);
}

export async function getLinkedProviders(): Promise<ListProvidersResponse> {
  return api.get<ListProvidersResponse>("/auth/social/providers");
}

export async function linkProvider(data: SocialLoginRequest): Promise<void> {
  return api.post("/auth/social/link", data);
}

export async function unlinkProvider(provider: string): Promise<void> {
  return api.delete(`/auth/social/${provider}`);
}
