"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import type { CategoryAmount } from "@/types";
import { formatCurrency, formatPercentage } from "@/lib/format";

interface TopCategoriesProps {
  categories: CategoryAmount[];
}

export function TopCategories({ categories }: TopCategoriesProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Top Categorias de Despesa</CardTitle>
      </CardHeader>
      <CardContent>
        {categories.length === 0 ? (
          <p className="text-sm text-muted-foreground">Sem dados</p>
        ) : (
          <div className="space-y-4">
            {categories.map((cat, i) => (
              <div key={i} className="space-y-2">
                <div className="flex items-center justify-between text-sm">
                  <span className="font-medium">{cat.category_name}</span>
                  <span className="text-muted-foreground">
                    {formatCurrency(cat.amount)} ({formatPercentage(cat.percentage)})
                  </span>
                </div>
                <Progress value={parseFloat(cat.percentage)} className="h-2" />
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
