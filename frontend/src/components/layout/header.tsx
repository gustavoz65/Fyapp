"use client";

import { Bell } from "lucide-react";
import Link from "next/link";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";
import { ThemeToggle } from "./theme-toggle";
import { api } from "@/lib/api";

export function Header() {
  const [unreadCount, setUnreadCount] = useState(0);

  const fetchUnread = useCallback(async () => {
    try {
      const data = await api.get<{ count: number }>("/notifications/unread/count");
      setUnreadCount(data.count);
    } catch {
      // ignore
    }
  }, []);

  useEffect(() => {
    let mounted = true;

    const loadUnread = async () => {
      if (mounted) {
        await fetchUnread();
      }
    };

    loadUnread();
    const interval = setInterval(() => {
      if (mounted) {
        fetchUnread();
      }
    }, 30000);

    return () => {
      mounted = false;
      clearInterval(interval);
    };
  }, [fetchUnread]);

  return (
    <header className="flex h-14 items-center gap-2 border-b px-4">
      <SidebarTrigger />
      <Separator orientation="vertical" className="h-6" />
      <div className="flex-1" />
      <ThemeToggle />
      <Button variant="ghost" size="icon" asChild className="relative">
        <Link href="/notifications">
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <Badge
              variant="destructive"
              className="absolute -right-1 -top-1 h-5 w-5 rounded-full p-0 text-xs flex items-center justify-center"
            >
              {unreadCount > 9 ? "9+" : unreadCount}
            </Badge>
          )}
        </Link>
      </Button>
    </header>
  );
}
