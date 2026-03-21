"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import {
  Bell,
  BarChart3,
  CalendarDays,
  FileText,
  Home,
  Layers,
  Library,
  PlusSquare,
  Settings,
  ShieldCheck,
  Users,
} from "lucide-react";
import {
  DashboardSidebarNav,
  type DashboardModuleItem,
} from "@/components/dashboard/dashboard-sidebar-nav";
import { DashboardHeaderNav } from "@/components/dashboard/dashboard-header-nav";
import { AuthService, type AuthUser } from "@/services/auth-service";
import GlobalLoader from "@/components/ui/global-loader";

const modules: DashboardModuleItem[] = [
  {
    id: "overview",
    label: "Overview",
    description: "Pulse and growth metrics",
    href: "/dashboard",
    icon: <Home size={18} />,
    notificationCount: 3,
  },
  {
    id: "create",
    label: "Create",
    description: "Compose and customize posts",
    href: "/create",
    icon: <PlusSquare size={18} />,
  },
  {
    id: "schedule",
    label: "Schedule",
    description: "Calendar, queue, bulk planning",
    href: "/schedule",
    icon: <CalendarDays size={18} />,
  },
  {
    id: "posts",
    label: "Posts",
    description: "Drafts, scheduled, published",
    href: "/posts",
    icon: <FileText size={18} />,
    badge: "12",
    notificationCount: 12,
  },
  {
    id: "analytics",
    label: "Analytics",
    description: "Performance and attribution",
    href: "/analytics",
    icon: <BarChart3 size={18} />,
  },
  {
    id: "accounts",
    label: "Accounts",
    description: "Connected social channels",
    href: "/accounts",
    icon: <Layers size={18} />,
  },
  {
    id: "team",
    label: "Team",
    description: "Roles, approvals, activity logs",
    href: "/team",
    icon: <Users size={18} />,
  },
  {
    id: "notifications",
    label: "Notifications",
    description: "Delivery and workflow alerts",
    href: "/notifications",
    icon: <Bell size={18} />,
  },
  {
    id: "library",
    label: "Library",
    description: "Reusable image and video assets",
    href: "/library",
    icon: <Library size={18} />,
  },
  {
    id: "onboarding",
    label: "Onboarding",
    description: "Connect accounts and set platforms",
    href: "/onboarding",
    icon: <ShieldCheck size={18} />,
  },
  {
    id: "settings",
    label: "Settings",
    description: "Workspace and policies",
    href: "/settings",
    icon: <Settings size={18} />,
  },
];

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [isMobileNavOpen, setIsMobileNavOpen] = useState(false);

  useEffect(() => {
    const boot = async () => {
      if (!AuthService.isAuthenticated()) {
        router.replace("/login");
        return;
      }
      try {
        const fetchedUser = await AuthService.me();
        setUser(fetchedUser);
        setLoading(false);
      } catch {
        AuthService.clearTokens();
        router.replace("/login");
      }
    };
    void boot();
  }, [router]);

  if (loading) return <GlobalLoader />;

  const activeModuleData =
    modules.find((m) => pathname === m.href || pathname.startsWith(`${m.href}/`)) ?? modules[0];

  return (
    <div className="flex min-h-screen flex-col bg-[var(--bg)] text-[var(--primary)]">

      {/* ─── Top header ─────────────────────────────────── */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3, ease: "easeOut" }}
      >
        <DashboardHeaderNav
          activeLabel={activeModuleData.label}
          activeDescription={activeModuleData.description}
          userName={user?.name}
          onOpenMobileNav={() => setIsMobileNavOpen(true)}
        />
      </motion.div>

      {/* ─── Body: sidebar + main ───────────────────────── */}
      <div className="flex flex-1 overflow-hidden">

        {/* Sidebar */}
        <motion.div
          initial={{ opacity: 0, x: -24 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.35, delay: 0.08, ease: "easeOut" }}
          className="shrink-0"
        >
          <DashboardSidebarNav
            modules={modules}
            activePath={pathname}
            isSidebarCollapsed={isSidebarCollapsed}
            isMobileNavOpen={isMobileNavOpen}
            onToggleSidebar={() => setIsSidebarCollapsed((prev) => !prev)}
            onCloseMobileNav={() => setIsMobileNavOpen(false)}
          />
        </motion.div>

        {/* Page content */}
        <motion.main
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, delay: 0.14, ease: "easeOut" }}
          className="flex-1 overflow-auto bg-[var(--surface)]"
        >
          {children}
        </motion.main>

      </div>
    </div>
  );
}
