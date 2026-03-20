export type AuthUser = {
  user_id: string;
  name: string;
  email: string;
};

type AuthPayload = {
  email: string;
  password: string;
};

type RegisterPayload = AuthPayload & {
  name: string;
};

import { apiClient, getApiErrorMessage } from "@/services/api-client";
import { ACCESS_COOKIE, REFRESH_COOKIE, clearCookie, getCookie } from "@/services/auth-cookies";

export const AuthService = {
  getAccessToken() {
    return getCookie(ACCESS_COOKIE);
  },

  getRefreshToken() {
    return getCookie(REFRESH_COOKIE);
  },

  isAuthenticated() {
    return Boolean(this.getAccessToken());
  },

  clearTokens() {
    clearCookie(ACCESS_COOKIE);
    clearCookie(REFRESH_COOKIE);
  },

  async login(payload: AuthPayload): Promise<AuthUser> {
    try {
      const { data } = await apiClient.post<AuthUser>("/auth/login", payload);
      return data;
    } catch (error) {
      throw new Error(getApiErrorMessage(error, "Unable to log in"));
    }
  },

  async register(payload: RegisterPayload): Promise<AuthUser> {
    try {
      const { data } = await apiClient.post<AuthUser>("/auth/register", payload);
      return data;
    } catch (error) {
      throw new Error(getApiErrorMessage(error, "Unable to create account"));
    }
  },

  async me(): Promise<AuthUser> {
    try {
      const { data } = await apiClient.get<AuthUser>("/auth/me");
      return data;
    } catch (error) {
      throw new Error(getApiErrorMessage(error, "Unable to load profile"));
    }
  },

  async logout(): Promise<void> {
    try {
      await apiClient.post("/auth/logout");
    } catch {
      this.clearTokens();
    }
  },
};
