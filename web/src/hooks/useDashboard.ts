import { useState, useEffect, useCallback, useRef } from 'react';
import type { DashboardState } from '../types/dashboard';

export interface UseDashboardOptions {
  pollIntervalMs?: number;
  apiUrl?: string;
}

export interface UseDashboardResult {
  data: DashboardState | null;
  loading: boolean;
  error: Error | null;
  isRefreshing: boolean;
  lastUpdated: Date | null;
  refresh: () => Promise<void>;
}

export function useDashboard(options: UseDashboardOptions = {}): UseDashboardResult {
  const { pollIntervalMs = 30000, apiUrl = '/api/dashboard' } = options;

  const [data, setData] = useState<DashboardState | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  const isMountedRef = useRef<boolean>(true);
  const abortControllerRef = useRef<AbortController | null>(null);

  const fetchDashboard = useCallback(
    async (isManualRefresh = false) => {
      if (isManualRefresh) {
        setIsRefreshing(true);
      }

      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
      const controller = new AbortController();
      abortControllerRef.current = controller;

      try {
        const response = await fetch(apiUrl, {
          signal: controller.signal,
          headers: {
            Accept: 'application/json',
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}: ${response.statusText}`);
        }

        const json: DashboardState = await response.json();

        if (isMountedRef.current) {
          setData(json);
          setError(null);
          setLastUpdated(new Date());
        }
      } catch (err: unknown) {
        if (err instanceof DOMException && err.name === 'AbortError') {
          return;
        }
        if (isMountedRef.current) {
          setError(err instanceof Error ? err : new Error(String(err)));
        }
      } finally {
        if (isMountedRef.current) {
          setLoading(false);
          if (isManualRefresh) {
            setIsRefreshing(false);
          }
        }
      }
    },
    [apiUrl]
  );

  const refresh = useCallback(async () => {
    await fetchDashboard(true);
  }, [fetchDashboard]);

  useEffect(() => {
    isMountedRef.current = true;
    fetchDashboard(false);

    if (pollIntervalMs > 0) {
      const timer = setInterval(() => {
        fetchDashboard(false);
      }, pollIntervalMs);

      return () => {
        isMountedRef.current = false;
        clearInterval(timer);
        if (abortControllerRef.current) {
          abortControllerRef.current.abort();
        }
      };
    }

    return () => {
      isMountedRef.current = false;
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, [fetchDashboard, pollIntervalMs]);

  return {
    data,
    loading,
    error,
    isRefreshing,
    lastUpdated,
    refresh,
  };
}
