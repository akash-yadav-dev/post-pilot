import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

const ACCESS_COOKIE = "postpilot_access_token";

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get(ACCESS_COOKIE)?.value;
  const protectedRoutes = [
    "/dashboard",
    "/create",
    "/schedule",
    "/posts",
    "/analytics",
    "/accounts",
    "/team",
    "/notifications",
    "/library",
    "/onboarding",
    "/settings",
  ];
  const isProtectedRoute = protectedRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`)
  );

  if (isProtectedRoute && !token) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/dashboard/:path*",
    "/create/:path*",
    "/schedule/:path*",
    "/posts/:path*",
    "/analytics/:path*",
    "/accounts/:path*",
    "/team/:path*",
    "/notifications/:path*",
    "/library/:path*",
    "/onboarding/:path*",
    "/settings/:path*",
    "/login",
    "/signup",
  ],
};
