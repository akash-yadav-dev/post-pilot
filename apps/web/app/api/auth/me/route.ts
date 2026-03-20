import { NextRequest, NextResponse } from "next/server";
import axios from "axios";
import {
  backendAuthClient,
  BackendAuthResponse,
  buildErrorResponse,
  clearAuthCookies,
  refreshSession,
  setAuthCookies,
  toAuthUser,
} from "@/app/api/auth/_lib/backend-auth";
import { ACCESS_COOKIE } from "@/services/auth-cookies";

export async function GET(request: NextRequest) {
  let accessToken = request.cookies.get(ACCESS_COOKIE)?.value;

  if (!accessToken) {
    try {
      const refreshed = await refreshSession(request);

      if (!refreshed) {
        return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
      }

      accessToken = refreshed.tokens.access_token;

      const { data } = await backendAuthClient.get("/me", {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      const response = NextResponse.json(toAuthUser(data));
      setAuthCookies(response, refreshed);
      return response;
    } catch (error) {
      const unauthorized = buildErrorResponse(error, "Not authenticated");
      clearAuthCookies(unauthorized);
      return unauthorized;
    }
  }

  try {
    const { data } = await backendAuthClient.get("/me", {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    return NextResponse.json(toAuthUser(data));
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      try {
        const refreshed = await refreshSession(request);

        if (!refreshed) {
          const unauthorized = NextResponse.json({ error: "Not authenticated" }, { status: 401 });
          clearAuthCookies(unauthorized);
          return unauthorized;
        }

        const { data } = await backendAuthClient.get("/me", {
          headers: {
            Authorization: `Bearer ${refreshed.tokens.access_token}`,
          },
        });

        const response = NextResponse.json(toAuthUser(data));
        setAuthCookies(response, refreshed as BackendAuthResponse);
        return response;
      } catch (refreshError) {
        const unauthorized = buildErrorResponse(refreshError, "Not authenticated");
        clearAuthCookies(unauthorized);
        return unauthorized;
      }
    }

    return buildErrorResponse(error, "Unable to load profile");
  }
}