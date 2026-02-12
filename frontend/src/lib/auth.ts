import { api } from "./api";
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  User,
  UserSettings,
  UpdateUserRequest,
  UpdateUserSettingsRequest,
  ChangePasswordRequest,
} from "@/types";

export async function login(data: LoginRequest): Promise<LoginResponse> {
  // Cookies são setados automaticamente pelo backend (httpOnly)
  return api.post<LoginResponse>("/auth/login", data);
}

export async function register(data: RegisterRequest): Promise<LoginResponse> {
  // Cookies são setados automaticamente pelo backend (httpOnly)
  return api.post<LoginResponse>("/auth/register", data);
}

export async function logout(): Promise<void> {
  // Backend limpa os cookies automaticamente
  await api.post("/auth/logout");
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
