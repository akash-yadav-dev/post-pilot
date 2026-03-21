"use client";

import { usePathname } from "next/navigation";
import { SiteFooter } from "@/components/layout/site-footer";
import { SiteHeader } from "@/components/layout/site-header";

type SiteShellProps = {
  children: React.ReactNode;
};

export function SiteShell({ children }: SiteShellProps) {
  const pathname = usePathname();
  const dashboardRoutes = [
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
  const isDashboardRoute = dashboardRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`)
  );

  const isAuthRoute = pathname === "/login" || pathname === "/signup";

  if (isDashboardRoute || isAuthRoute) {
    return <>{children}</>;
  }

  return (
    <>
      <SiteHeader />
      <main>{children}</main>
      <SiteFooter />
    </>
  );
}
