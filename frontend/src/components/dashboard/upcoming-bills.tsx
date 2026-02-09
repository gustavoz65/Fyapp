"use client";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatCurrency, formatDate } from "@/lib/format";
import type { Transaction } from "@/types";

interface UpcomingBillsProps {
  transactions: Transaction[];
}

export function UpcomingBills({ transactions }: UpcomingBillsProps) {
  const items = transactions ?? [];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Proximas Contas</CardTitle>
      </CardHeader>
      <CardContent>
        {items.length === 0 ? (
          <p className="text-sm text-muted-foreground">Nenhuma conta proxima</p>
        ) : (
          <div className="space-y-3">
            {items.map((tx) => (
              <div key={tx.id} className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium">{tx.description}</p>
                  <p className="text-xs text-muted-foreground">
                    {tx.due_date
                      ? formatDate(tx.due_date)
                      : formatDate(tx.transaction_date)}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={tx.is_paid ? "default" : "destructive"}>
                    {tx.is_paid ? "Pago" : "Pendente"}
                  </Badge>
                  <span className="text-sm font-semibold">
                    {formatCurrency(tx.amount)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
