"use client";

import { setRedirectCallback } from "@/lib/api";
import { useAuthStore } from "@/stores/auth-store";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";

export function AuthInitializer() {
  const initialize = useAuthStore((s) => s.initialize);
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    setRedirectCallback(() => {
      const publicRoutes = ["/login", "/register"];
      const isPublic = publicRoutes.some(
        (r) => pathname === r || pathname?.startsWith(r),
      );
      if (!isPublic) {
        router.push("/login");
      }
    });
  }, [router, pathname]);

  useEffect(() => {
    initialize();
  }, [initialize]);

  return null;
}
