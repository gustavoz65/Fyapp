"use client";

import { useCallback } from "react";
import { useApi } from "./use-api";
import { api } from "@/lib/api";
import type { FinancialHealthSnapshot, HistoricalScoresResponse } from "@/types";

export function useFinancialHealth() {
  const fetcher = useCallback(() => {
    return api.get<FinancialHealthSnapshot>("/health/score");
  }, []);

  const { data: currentScore, error, isLoading, execute } = useApi<FinancialHealthSnapshot>(
    fetcher,
    { immediate: true }
  );

  return {
    currentScore,
    error,
    isLoading,
    refetch: execute,
  };
}

export function useHistoricalHealth(immediate = false) {
  const fetcher = useCallback(() => {
    return api.get<HistoricalScoresResponse>("/health/history");
  }, []);

  const { data, error, isLoading, execute } = useApi<HistoricalScoresResponse>(
    fetcher,
    { immediate }
  );

  return {
    history: data?.snapshots || [],
    error,
    isLoading,
    fetchHistory: execute,
  };
}
