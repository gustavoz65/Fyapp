"use client";

import { useState, useEffect, useCallback } from "react";

interface UseApiOptions {
  immediate?: boolean;
}

export function useApi<T>(
  fetcher: () => Promise<T>,
  options: UseApiOptions = { immediate: true }
) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const execute = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await fetcher();
      setData(result);
      return result;
    } catch (err: unknown) {
      const message = err && typeof err === "object" && "message" in err
        ? (err as { message: string }).message
        : "An error occurred";
      setError(message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [fetcher]);

  useEffect(() => {
    if (options.immediate) {
      execute();
    }
  }, [execute, options.immediate]);

  return { data, error, isLoading, execute, setData };
}

export function useMutation<TData, TArgs = void>(
  mutationFn: (args: TArgs) => Promise<TData>
) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const mutate = useCallback(
    async (args: TArgs) => {
      setIsLoading(true);
      setError(null);
      try {
        const result = await mutationFn(args);
        return result;
      } catch (err: unknown) {
        const message = err && typeof err === "object" && "message" in err
          ? (err as { message: string }).message
          : "An error occurred";
        setError(message);
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [mutationFn]
  );

  return { mutate, isLoading, error };
}
