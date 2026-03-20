import { NextRequest, NextResponse } from "next/server";
import axios from "axios";
import {
  backendAuthClient,
  buildErrorResponse,
  clearAuthCookies,
} from "@/app/api/auth/_lib/backend-auth";
import { REFRESH_COOKIE } from "@/services/auth-cookies";

export async function POST(request: NextRequest) {
  const response = new NextResponse(null, { status: 204 });

  try {
    const refreshToken = request.cookies.get(REFRESH_COOKIE)?.value;

    if (refreshToken) {
      await backendAuthClient.post("/logout", {
        refresh_token: refreshToken,
      });
    }

    clearAuthCookies(response);
    return response;
  } catch (error) {
    clearAuthCookies(response);

    if (axios.isAxiosError(error) && error.response?.status === 401) {
      return response;
    }

    return buildErrorResponse(error, "Logout failed");
  }
}