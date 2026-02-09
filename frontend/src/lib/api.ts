import type { APIError } from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

class ApiClient {
  private accessToken: string | null = null;
  private refreshToken: string | null = null;
  private refreshPromise: Promise<void> | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      this.accessToken = localStorage.getItem("access_token");
      this.refreshToken = localStorage.getItem("refresh_token");
    }
  }

  setTokens(accessToken: string, refreshToken: string) {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    if (typeof window !== "undefined") {
      localStorage.setItem("access_token", accessToken);
      localStorage.setItem("refresh_token", refreshToken);
    }
  }

  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    if (typeof window !== "undefined") {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
    }
  }

  getAccessToken() {
    return this.accessToken;
  }

  private async refreshAccessToken(): Promise<void> {
    if (!this.refreshToken) {
      this.clearTokens();
      throw new Error("No refresh token");
    }

    const res = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: this.refreshToken }),
    });

    if (!res.ok) {
      this.clearTokens();
      throw new Error("Token refresh failed");
    }

    const data = await res.json();
    this.setTokens(data.access_token, data.refresh_token);
  }

  private async ensureValidToken(): Promise<void> {
    if (!this.accessToken) return;

    try {
      const payload = JSON.parse(atob(this.accessToken.split(".")[1]));
      const expiresAt = payload.exp * 1000;
      const now = Date.now();

      if (expiresAt - now < 60000) {
        if (!this.refreshPromise) {
          this.refreshPromise = this.refreshAccessToken().finally(() => {
            this.refreshPromise = null;
          });
        }
        await this.refreshPromise;
      }
    } catch {
      // Token parsing failed, will try the request anyway
    }
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    await this.ensureValidToken();

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    };

    if (this.accessToken) {
      headers["Authorization"] = `Bearer ${this.accessToken}`;
    }

    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (res.status === 401 && this.refreshToken) {
      try {
        await this.refreshAccessToken();
        headers["Authorization"] = `Bearer ${this.accessToken}`;
        const retryRes = await fetch(`${API_BASE_URL}${endpoint}`, {
          ...options,
          headers,
        });
        if (!retryRes.ok) {
          const error = await retryRes.json().catch(() => ({}));
          throw error as APIError;
        }
        if (retryRes.status === 204) return undefined as T;
        return retryRes.json();
      } catch {
        this.clearTokens();
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
        throw new Error("Session expired");
      }
    }

    if (!res.ok) {
      const error = await res.json().catch(() => ({
        code: "UNKNOWN",
        message: "An unexpected error occurred",
      }));
      throw error as APIError;
    }

    if (res.status === 204) return undefined as T;
    return res.json();
  }

  get<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: "GET" });
  }

  post<T>(endpoint: string, body?: unknown) {
    return this.request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  put<T>(endpoint: string, body?: unknown) {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  patch<T>(endpoint: string, body?: unknown) {
    return this.request<T>(endpoint, {
      method: "PATCH",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  delete<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: "DELETE" });
  }
}

export const api = new ApiClient();
