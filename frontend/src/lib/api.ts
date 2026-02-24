import type { APIError } from "@/types";

const API_BASE_URL = "http://localhost:3000/api/v1";

let redirectToLogin: (() => void) | null = null;

export function setRedirectCallback(callback: () => void) {
  redirectToLogin = callback;
}

class ApiClient {
  private refreshPromise: Promise<void> | null = null;

  private async refreshAccessToken(): Promise<void> {
    const res = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: "POST",
      credentials: "include",
    });

    if (!res.ok) {
      throw new Error("Token refresh failed");
    }
  }

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    };

    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
      credentials: "include",
    });

    if (res.status === 401) {
      try {
        // Evitar múltiplos refreshes simultâneos
        if (!this.refreshPromise) {
          this.refreshPromise = this.refreshAccessToken().finally(() => {
            this.refreshPromise = null;
          });
        }
        await this.refreshPromise;

        // Retry request original
        const retryRes = await fetch(`${API_BASE_URL}${endpoint}`, {
          ...options,
          headers,
          credentials: "include",
        });

        if (!retryRes.ok) {
          const error = await retryRes.json().catch(() => ({}));
          throw error as APIError;
        }

        if (retryRes.status === 204) return undefined as T;
        return retryRes.json();
      } catch {
        if (redirectToLogin) {
          redirectToLogin();
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
