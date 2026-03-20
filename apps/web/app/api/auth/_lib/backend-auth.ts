import axios from "axios";
import { NextRequest, NextResponse } from "next/server";
import {
  ACCESS_COOKIE,
  ACCESS_TOKEN_MAX_AGE_SECONDS,
  REFRESH_COOKIE,
  REFRESH_TOKEN_MAX_AGE_SECONDS,
} from "@/services/auth-cookies";

export type BackendAuthResponse = {
  user_id: string;
  name: string;
  email: string;
  tokens: {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
  };
};

export type AuthUser = {
  user_id: string;
  name: string;
  email: string;
};

const backendApiBaseUrl =
  process.env.API_BASE_URL?.replace(/\/$/, "") ??
  process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ??
  "http://localhost:8080";

export const backendAuthClient = axios.create({
  baseURL: `${backendApiBaseUrl}/api/v1/auth`,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 15000,
});

export function toAuthUser(data: AuthUser) {
  return {
    user_id: data.user_id,
    name: data.name,
    email: data.email,
  };
}

export function setAuthCookies(response: NextResponse, data: BackendAuthResponse) {
  response.cookies.set(ACCESS_COOKIE, data.tokens.access_token, {
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: ACCESS_TOKEN_MAX_AGE_SECONDS,
  });

  response.cookies.set(REFRESH_COOKIE, data.tokens.refresh_token, {
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: REFRESH_TOKEN_MAX_AGE_SECONDS,
  });
}

export function clearAuthCookies(response: NextResponse) {
  response.cookies.set(ACCESS_COOKIE, "", {
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 0,
  });

  response.cookies.set(REFRESH_COOKIE, "", {
    httpOnly: false,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 0,
  });
}

export async function refreshSession(request: NextRequest) {
  const refreshToken = request.cookies.get(REFRESH_COOKIE)?.value;

  if (!refreshToken) {
    return null;
  }

  const { data } = await backendAuthClient.post<BackendAuthResponse>("/refresh", {
    refresh_token: refreshToken,
  });

  return data;
}

export function buildErrorResponse(error: unknown, fallback = "Request failed") {
  if (axios.isAxiosError<{ error?: string }>(error)) {
    const status = error.response?.status ?? 500;
    const message = error.response?.data?.error ?? error.message ?? fallback;

    return NextResponse.json({ error: message }, { status });
  }

  const message = error instanceof Error ? error.message : fallback;
  return NextResponse.json({ error: message }, { status: 500 });
}