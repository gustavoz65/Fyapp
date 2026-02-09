"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Transaction } from "@/types";
import { formatCurrency, formatDate } from "@/lib/format";

interface UpcomingBillsProps {
  transactions: Transaction[];
}

export function UpcomingBills({ transactions }: UpcomingBillsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Proximas Contas</CardTitle>
      </CardHeader>
      <CardContent>
        {transactions.length === 0 ? (
          <p className="text-sm text-muted-foreground">Nenhuma conta proxima</p>
        ) : (
          <div className="space-y-3">
            {transactions.map((tx) => (
              <div key={tx.id} className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium">{tx.description}</p>
                  <p className="text-xs text-muted-foreground">
                    {tx.due_date ? formatDate(tx.due_date) : formatDate(tx.transaction_date)}
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
