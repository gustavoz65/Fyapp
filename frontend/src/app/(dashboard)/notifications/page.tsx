"use client";

import { useCallback, useEffect, useState } from "react";
import { Bell, CheckCheck, Trash2, Eye } from "lucide-react";
import { api } from "@/lib/api";
import type { Notification } from "@/types";
import { formatRelativeDate, getNotificationTypeLabel } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showUnreadOnly, setShowUnreadOnly] = useState(false);

  const fetchNotifications = useCallback(async () => {
    try {
      const endpoint = showUnreadOnly ? "/notifications/unread" : "/notifications";
      const data = await api.get<Notification[]>(endpoint);
      setNotifications(data || []);
    } catch {} finally { setIsLoading(false); }
  }, [showUnreadOnly]);

  useEffect(() => { fetchNotifications(); }, [fetchNotifications]);

  async function markAsRead(id: string) {
    try {
      await api.patch(`/notifications/${id}/read`);
      fetchNotifications();
    } catch { toast.error("Erro ao marcar como lida"); }
  }

  async function markAllAsRead() {
    try {
      await api.patch("/notifications/read-all");
      toast.success("Todas marcadas como lidas");
      fetchNotifications();
    } catch { toast.error("Erro"); }
  }

  async function deleteNotification(id: string) {
    try {
      await api.delete(`/notifications/${id}`);
      toast.success("Notificacao removida");
      fetchNotifications();
    } catch { toast.error("Erro ao remover"); }
  }

  const notificationIcon: Record<string, string> = {
    budget_alert: "text-yellow-500",
    bill_reminder: "text-blue-500",
    goal_achieved: "text-green-500",
    low_balance: "text-red-500",
    transaction_alert: "text-purple-500",
    system: "text-muted-foreground",
  };

  if (isLoading) {
    return (
      <div className="space-y-8 pb-8">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Notificacoes</h1>
          <p className="text-muted-foreground mt-2">
            Acompanhe suas notificacoes e alertas
          </p>
        </div>
        {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-24" />)}
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Notificacoes</h1>
          <p className="text-muted-foreground mt-2">
            Acompanhe suas notificacoes e alertas
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant={showUnreadOnly ? "default" : "outline"} size="sm" onClick={() => setShowUnreadOnly(!showUnreadOnly)}>
            {showUnreadOnly ? "Ver todas" : "Apenas nao lidas"}
          </Button>
          <Button variant="outline" size="sm" onClick={markAllAsRead}>
            <CheckCheck className="h-4 w-4 mr-2" />Marcar todas como lidas
          </Button>
        </div>
      </div>

      {notifications.length === 0 ? (
        <Card className="border-2">
          <CardContent className="py-16 text-center">
            <Bell className="h-16 w-16 mx-auto text-muted-foreground mb-4 opacity-50" />
            <h3 className="text-lg font-semibold mb-2">Nenhuma notificacao</h3>
            <p className="text-sm text-muted-foreground">
              Voce esta em dia! Nao ha notificacoes no momento
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {notifications.map((notification) => (
            <Card key={notification.id} className={`hover:shadow-md transition-all ${!notification.is_read ? "border-primary/50 bg-primary/5" : ""}`}>
              <CardContent className="flex items-start gap-4 p-5">
                <div className={`mt-1 ${notificationIcon[notification.type] || "text-muted-foreground"}`}>
                  <Bell className="h-5 w-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <p className="text-sm font-medium">{notification.title}</p>
                    <Badge variant="outline" className="text-xs">{getNotificationTypeLabel(notification.type)}</Badge>
                    {!notification.is_read && <Badge variant="default" className="text-xs">Nova</Badge>}
                  </div>
                  <p className="text-sm text-muted-foreground mt-1">{notification.message}</p>
                  <p className="text-xs text-muted-foreground mt-1">{formatRelativeDate(notification.created_at)}</p>
                </div>
                <div className="flex gap-1">
                  {!notification.is_read && (
                    <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => markAsRead(notification.id)}>
                      <Eye className="h-4 w-4" />
                    </Button>
                  )}
                  <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => deleteNotification(notification.id)}>
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
