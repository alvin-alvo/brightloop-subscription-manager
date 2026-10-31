// Base API configuration and helper for Brightloop
export const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api';

/**
 * Error thrown by the API client. Carries the HTTP status so callers can tell
 * an auth/permission failure (401/403) — which may warrant a local fallback —
 * from a genuine server error.
 */
export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export async function apiClient<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const token = localStorage.getItem('brightloop_token');
  const headers = new Headers(options?.headers || {});
  headers.set('Content-Type', 'application/json');

  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    const message =
      (errorData as { message?: string })?.message ||
      `Request failed with status ${response.status}`;
    throw new ApiError(message, response.status);
  }

  return response.json();
}
