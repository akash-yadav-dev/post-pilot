"use client";

import Link from "next/link";
import { ChevronLeft, ChevronRight, PlusCircle, X } from "lucide-react";
import { type ReactNode } from "react";
import { motion, AnimatePresence } from "framer-motion";

export type DashboardModuleItem = {
  id: string;
  label: string;
  description: string;
  href: string;
  icon: ReactNode;
  badge?: string;
  notificationCount?: number;
};

type DashboardSidebarNavProps = {
  modules: DashboardModuleItem[];
  activePath: string;
  isSidebarCollapsed: boolean;
  isMobileNavOpen: boolean;
  onToggleSidebar: () => void;
  onCloseMobileNav: () => void;
};

export function DashboardSidebarNav({
  modules,
  activePath,
  isSidebarCollapsed,
  isMobileNavOpen,
  onToggleSidebar,
  onCloseMobileNav,
}: DashboardSidebarNavProps) {
  return (
    <>
      <aside
        className={`hidden md:flex md:flex-col md:rounded-none md:border md:border-[var(--line)] md:bg-[var(--surface)] md:transition-all md:duration-500 min-h-screen ${
          isSidebarCollapsed ? "md:w-[68px]" : "md:w-[280px]"
        } md:border-bs-0`}
      >
        <div className="flex items-center justify-between border-b border-[var(--line)] p-4">
          <AnimatePresence mode="wait">
            {!isSidebarCollapsed && (
              <motion.div
                initial={{ opacity: 0, width: 0 }}
                animate={{ opacity: 1, width: "auto" }}
                exit={{ opacity: 0, width: 0 }}
                transition={{
                  duration: 0.2,
                  delay: 0.3,
                  ease: "easeOut",
                }}
                className="overflow-hidden whitespace-nowrap"
              >
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-[var(--muted)]">
                    Workspace
                  </p>
                  <p className="mt-1 text-sm font-semibold text-[var(--primary)]">
                    PostPilot Core
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
          <button
            type="button"
            aria-label="Toggle sidebar"
            onClick={onToggleSidebar}
            className="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-[var(--line)] text-[var(--primary)] transition hover:bg-[var(--bg)] cursor-pointer"
          >
            {isSidebarCollapsed ? (
              <ChevronRight size={16} />
            ) : (
              <ChevronLeft size={16} />
            )}
          </button>
        </div>

        <nav className="flex-1 space-y-1 p-3">
          {modules.map((item) => {
            const isActive = activePath === item.href || activePath.startsWith(`${item.href}/`);
            const hasNotification =
              item.notificationCount && item.notificationCount > 0;
            const notificationCount = item.notificationCount || 0;

            return (
              <Link
                key={item.id}
                href={item.href}
                className={`relative flex ${
                  !isSidebarCollapsed ? "w-full" : "w-auto justify-center"
                } items-center rounded-2xl px-3 py-3 text-left transition-all duration-300 cursor-pointer ${
                  isActive
                    ? "bg-[var(--primary)] text-[var(--bg)]"
                    : "text-[var(--primary)] hover:bg-[var(--primary)]/10"
                }`}
              >
                {/* Icon Container with Notification Badge for Collapsed State */}
                <div className="relative shrink-0">
                  <span className="block">{item.icon}</span>

                  {/* Notification Badge - Visible when sidebar is collapsed */}
                  {isSidebarCollapsed && hasNotification && (
                    <span
                      className={`
                      absolute -top-3 -right-3 
                      flex items-center justify-center 
                      min-w-[18px] h-[18px] 
                      px-1 text-[10px] font-bold 
                      rounded-full 
                      ${
                        isActive
                          ? "bg-[var(--primary)] text-[var(--bg)]"
                          : "bg-red-500 text-white"
                      }
                      ring-2 ring-[var(--surface)]
                    `}
                    >
                      {notificationCount > 99 ? "99+" : notificationCount}
                    </span>
                  )}
                </div>

                {/* Expanded State Content with Animation */}
                <AnimatePresence mode="wait">
                  {!isSidebarCollapsed && (
                    <motion.span
                      initial={{ opacity: 0, width: 0 }}
                      animate={{ opacity: 1, width: "auto" }}
                      exit={{ opacity: 0, width: 0 }}
                      transition={{
                        duration: 0.2,
                        delay: 0.25,
                        ease: "easeOut",
                      }}
                      className="flex min-w-0 flex-1 items-center justify-between gap-2 ps-3 overflow-hidden whitespace-nowrap"
                    >
                      <span className="truncate text-sm font-medium">
                        {item.label}
                      </span>
                      <div className="flex items-center gap-1">
                        {/* Notification Badge for Expanded State */}
                        {hasNotification && (
                          <span
                            className={`
                              rounded-full px-2 py-0.5 text-xs font-medium
                              ${
                                isActive
                                  ? "bg-[var(--bg)]/20 text-[var(--bg)]"
                                  : "bg-red-500 text-white"
                              }
                            `}
                          >
                            {notificationCount > 99 ? "99+" : notificationCount}
                          </span>
                        )}
                        {/* Regular Badge */}
                        {item.badge && !hasNotification && (
                          <span
                            className={`rounded-full px-2 py-0.5 text-xs ${
                              isActive
                                ? "bg-[var(--bg)]/20"
                                : "bg-[var(--line)]"
                            }`}
                          >
                            {item.badge}
                          </span>
                        )}
                      </div>
                    </motion.span>
                  )}
                </AnimatePresence>
              </Link>
            );
          })}
        </nav>

        {/* New Campaign Button with Animation */}
        <AnimatePresence mode="wait">
          {!isSidebarCollapsed && (
            <motion.div
              initial={{ opacity: 0, height: 0, y: 20 }}
              animate={{ opacity: 1, height: "auto", y: 0 }}
              exit={{ opacity: 0, height: 0, y: 20 }}
              transition={{
                duration: 0.2,
                delay: 0.25,
                ease: "easeOut",
              }}
              className="border-t border-[var(--line)] p-4 overflow-hidden"
            >
              <button
                type="button"
                className="inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-[var(--primary)] px-4 py-3 text-sm font-semibold text-[var(--bg)] transition hover:opacity-90"
              >
                <PlusCircle size={17} />
                New Campaign
              </button>
            </motion.div>
          )}
        </AnimatePresence>
      </aside>

      {isMobileNavOpen && (
        <div
          className="fixed inset-0 z-50 md:hidden"
          role="dialog"
          aria-modal="true"
        >
          <button
            type="button"
            aria-label="Close navigation overlay"
            className="absolute inset-0 bg-black/40"
            onClick={onCloseMobileNav}
          />

          <div className="absolute left-0 top-0 h-full w-[88%] max-w-sm border-r border-[var(--line)] bg-[var(--surface)] p-4 shadow-xl">
            <div className="mb-4 flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.2em] text-[var(--muted)]">
                  Navigator
                </p>
                <p className="mt-1 text-sm font-semibold">PostPilot Modules</p>
              </div>
              <button
                type="button"
                aria-label="Close mobile navigation"
                onClick={onCloseMobileNav}
                className="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-[var(--line)]"
              >
                <X size={16} />
              </button>
            </div>

            <nav className="space-y-2">
              {modules.map((item) => {
                const isActive = activePath === item.href || activePath.startsWith(`${item.href}/`);
                const hasNotification =
                  item.notificationCount && item.notificationCount > 0;
                const notificationCount = item.notificationCount || 0;

                return (
                  <Link
                    key={item.id}
                    href={item.href}
                    onClick={onCloseMobileNav}
                    className={`flex w-full items-center gap-3 rounded-2xl px-3 py-3 text-left transition cursor-pointer ${
                      isActive
                        ? "bg-[var(--primary)] text-[var(--bg)]"
                        : "text-[var(--primary)] hover:bg-[var(--bg)]"
                    }`}
                  >
                    <div className="relative shrink-0">
                      <span>{item.icon}</span>
                      {/* Mobile notification badge */}
                      {hasNotification && (
                        <span
                          className={`
                          absolute -top-3 -right-3 
                          flex items-center justify-center 
                          min-w-[18px] h-[18px] 
                          px-1 text-[10px] font-bold 
                          rounded-full 
                          ${
                            isActive
                              ? "bg-[var(--bg)]/20 text-[var(--bg)]"
                              : "bg-red-500 text-white"
                          }
                          ring-2 ring-[var(--surface)]
                        `}
                        >
                          {notificationCount > 99 ? "99+" : notificationCount}
                        </span>
                      )}
                    </div>
                    <span className="flex-1">
                      <span className="block text-sm font-medium">
                        {item.label}
                      </span>
                      <span
                        className={`block text-xs ${isActive ? "text-[var(--bg)]/70" : "text-[var(--muted)]"}`}
                      >
                        {item.description}
                      </span>
                    </span>
                    <div className="flex items-center gap-1">
                      {/* Mobile notification count badge */}
                      {hasNotification && (
                        <span
                          className={`rounded-full px-2 py-0.5 text-xs font-medium ${isActive ? "bg-[var(--bg)]/20" : "bg-red-500 text-white"}`}
                        >
                          {notificationCount > 99 ? "99+" : notificationCount}
                        </span>
                      )}
                      {/* Regular badge */}
                      {item.badge && !hasNotification && (
                        <span
                          className={`rounded-full px-2 py-0.5 text-xs ${isActive ? "bg-[var(--bg)]/20" : "bg-[var(--line)]"}`}
                        >
                          {item.badge}
                        </span>
                      )}
                    </div>
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>
      )}
    </>
  );
}
