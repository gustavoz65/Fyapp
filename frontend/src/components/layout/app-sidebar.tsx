"use client";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useAuth } from "@/providers/auth-provider";
import {
  ArrowLeftRight,
  Bell,
  Landmark,
  LayoutDashboard,
  LogOut,
  PiggyBank,
  Repeat,
  Settings,
  Tags,
  Target,
} from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  { title: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { title: "Contas", href: "/accounts", icon: Landmark },
  { title: "Transações", href: "/transactions", icon: ArrowLeftRight },
  { title: "Recorrentes", href: "/recurring", icon: Repeat },
  { title: "Orçamentos", href: "/budgets", icon: PiggyBank },
  { title: "Metas", href: "/goals", icon: Target },
  { title: "Categorias", href: "/categories", icon: Tags },
  { title: "Notificações", href: "/notifications", icon: Bell },
  { title: "Configurações", href: "/settings", icon: Settings },
];

export function AppSidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  const initials = user
    ? `${user.first_name[0]}${user.last_name[0]}`.toUpperCase()
    : "??";

  return (
    <Sidebar>
      <SidebarHeader>
        <div className="px-4 py-6 flex items-center justify-center border-b border-sidebar-border/50">
          <Image
            src="/logo_cash_no_dark_mode.png"
            alt="Cashing"
            width={120}
            height={40}
            className="object-contain brightness-0 invert"
            priority
          />
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Menu</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => (
                <SidebarMenuItem key={item.href}>
                  <SidebarMenuButton asChild isActive={pathname === item.href}>
                    <Link href={item.href}>
                      <item.icon className="h-4 w-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton className="w-full">
              <Avatar className="h-6 w-6">
                <AvatarFallback className="text-xs">{initials}</AvatarFallback>
              </Avatar>
              <span className="truncate">
                {user ? `${user.first_name} ${user.last_name}` : "..."}
              </span>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton onClick={logout}>
              <LogOut className="h-4 w-4" />
              <span>Sair</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
