"use client";

import { Activity } from "lucide-react";
import { cn } from "@/lib/utils";
import { useFinancialHealth } from "@/hooks/use-financial-health";
import { useHealthModal } from "@/contexts/health-modal-context";
import { HealthModal } from "./health-modal";

export function HealthWidget() {
  const { openModal } = useHealthModal();
  const { currentScore, isLoading } = useFinancialHealth();

  if (isLoading || !currentScore) {
    return null;
  }

  const scorePercentage = parseFloat(currentScore.score_percentage);

  // Determine color based on score
  const getScoreColor = (score: number) => {
    if (score < 40) return "text-red-600 dark:text-red-500";
    if (score < 70) return "text-yellow-600 dark:text-yellow-500";
    return "text-green-600 dark:text-green-500";
  };

  const getProgressColor = (score: number) => {
    if (score < 40) return "stroke-red-600 dark:stroke-red-500";
    if (score < 70) return "stroke-yellow-600 dark:stroke-yellow-500";
    return "stroke-green-600 dark:stroke-green-500";
  };

  const scoreColor = getScoreColor(scorePercentage);
  const progressColor = getProgressColor(scorePercentage);

  // Calculate circle progress (circumference based on 70mm diameter = 35mm radius)
  const radius = 32;
  const circumference = 2 * Math.PI * radius;
  const progress = ((100 - scorePercentage) / 100) * circumference;

  return (
    <>
      <button
        onClick={openModal}
        className="fixed bottom-6 right-6 z-40 group hover:scale-105 transition-transform focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 rounded-full"
        aria-label="Visualizar saúde financeira"
      >
        <div className="relative w-[70px] h-[70px]">
          {/* Background circle */}
          <svg className="w-full h-full -rotate-90" viewBox="0 0 72 72">
            <circle
              cx="36"
              cy="36"
              r={radius}
              className="stroke-muted fill-none"
              strokeWidth="4"
            />
            {/* Progress circle */}
            <circle
              cx="36"
              cy="36"
              r={radius}
              className={cn("fill-none transition-all duration-500", progressColor)}
              strokeWidth="4"
              strokeDasharray={circumference}
              strokeDashoffset={progress}
              strokeLinecap="round"
            />
          </svg>

          {/* Center content */}
          <div className="absolute inset-0 flex flex-col items-center justify-center bg-background rounded-full border shadow-lg group-hover:shadow-xl transition-shadow">
            <Activity className={cn("h-5 w-5 mb-0.5", scoreColor)} />
            <span className={cn("text-sm font-bold tabular-nums", scoreColor)}>
              {Math.round(scorePercentage)}%
            </span>
          </div>
        </div>
      </button>

      <HealthModal currentScore={currentScore} />
    </>
  );
}
