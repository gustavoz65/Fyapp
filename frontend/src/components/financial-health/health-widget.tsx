"use client";

import { Activity } from "lucide-react";
import { cn } from "@/lib/utils";
import { Progress } from "@/components/ui/progress";
import { useFinancialHealth } from "@/hooks/use-financial-health";
import { useHealthModal } from "@/contexts/health-modal-context";
import { HealthModal } from "./health-modal";

export function HealthWidget() {
  const { openModal } = useHealthModal();
  const { currentScore, isLoading } = useFinancialHealth();

  const scorePercentage = currentScore ? parseFloat(currentScore.score_percentage) : 0;

  // Determine color based on score
  const getScoreColor = (score: number) => {
    if (score < 40) return "text-red-600 dark:text-red-500";
    if (score < 70) return "text-yellow-600 dark:text-yellow-500";
    return "text-green-600 dark:text-green-500";
  };

  const scoreColor = getScoreColor(scorePercentage);

  return (
    <>
      {!isLoading && currentScore && (
        <button
          onClick={openModal}
          className="fixed bottom-6 right-6 z-40 group hover:scale-105 transition-transform focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 rounded-full"
          aria-label="Visualizar saúde financeira"
        >
          <div className="relative w-[70px] h-[70px] flex items-center justify-center bg-background rounded-full border shadow-lg group-hover:shadow-xl transition-shadow">
            {/* Circular Progress using Radix UI */}
            <svg className="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 100 100">
              <circle
                cx="50"
                cy="50"
                r="45"
                fill="none"
                stroke="currentColor"
                strokeWidth="6"
                className="text-muted"
                opacity="0.2"
              />
              <circle
                cx="50"
                cy="50"
                r="45"
                fill="none"
                stroke="currentColor"
                strokeWidth="6"
                strokeDasharray={`${2 * Math.PI * 45}`}
                strokeDashoffset={`${2 * Math.PI * 45 * (1 - scorePercentage / 100)}`}
                strokeLinecap="round"
                className={cn("transition-all duration-500", scoreColor)}
              />
            </svg>

            {/* Center content */}
            <div className="relative flex flex-col items-center justify-center">
              <Activity className={cn("h-5 w-5 mb-0.5", scoreColor)} />
              <span className={cn("text-sm font-bold tabular-nums", scoreColor)}>
                {Math.round(scorePercentage)}%
              </span>
            </div>
          </div>
        </button>
      )}

      <HealthModal currentScore={currentScore} />
    </>
  );
}
