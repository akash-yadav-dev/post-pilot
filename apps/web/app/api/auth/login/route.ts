import { NextResponse } from "next/server";
import {
  backendAuthClient,
  BackendAuthResponse,
  buildErrorResponse,
  setAuthCookies,
  toAuthUser,
} from "@/app/api/auth/_lib/backend-auth";

export async function POST(request: Request) {
  try {
    const payload = await request.json();
    const { data } = await backendAuthClient.post<BackendAuthResponse>("/login", payload);

    const response = NextResponse.json(toAuthUser(data));
    setAuthCookies(response, data);

    return response;
  } catch (error) {
    return buildErrorResponse(error, "Login failed");
  }
}