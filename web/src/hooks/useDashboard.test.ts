import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useDashboard } from './useDashboard';
import type { DashboardState } from '../types/dashboard';

const mockDashboardData: DashboardState = {
  pull_requests: [
    {
      repo: 'brotherlogic/seraphine',
      pr_number: 180,
      title: 'Scaffold Frontend React Project',
      author: 'brotherlogic-automation',
      commit_count: 2,
      comment_count: 0,
      check_status: 'SUCCESS',
      has_devcontainer: true,
      container_id: 'seraphine-180',
      container_state: 'READY',
    },
  ],
  freshness: {
    last_successful_sync: '2026-08-20T17:00:00Z',
    last_attempted_sync: '2026-08-20T17:00:00Z',
    is_stale: false,
  },
};

describe('useDashboard Hook', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('fetches dashboard data successfully on mount', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: async () => mockDashboardData,
    } as Response);

    const { result } = renderHook(() => useDashboard());

    expect(result.current.loading).toBe(true);
    expect(result.current.data).toBeNull();
    expect(result.current.error).toBeNull();

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.data).toEqual(mockDashboardData);
    expect(result.current.error).toBeNull();
    expect(result.current.lastUpdated).toBeInstanceOf(Date);
    expect(globalThis.fetch).toHaveBeenCalledWith('/api/dashboard', expect.any(Object));
  });

  it('handles fetch failure and updates error state', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
    } as Response);

    const { result } = renderHook(() => useDashboard());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.data).toBeNull();
    expect(result.current.error).toBeTruthy();
    expect(result.current.error?.message).toContain('HTTP error 500');
  });

  it('handles network throw error properly', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network disconnected'));

    const { result } = renderHook(() => useDashboard());

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.data).toBeNull();
    expect(result.current.error?.message).toBe('Network disconnected');
  });

  it('supports manual refresh with isRefreshing state', async () => {
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockDashboardData,
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          ...mockDashboardData,
          pull_requests: [],
        }),
      } as Response);

    const { result } = renderHook(() => useDashboard());

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.data?.pull_requests.length).toBe(1);

    let refreshPromise: Promise<void>;
    act(() => {
      refreshPromise = result.current.refresh();
    });

    expect(result.current.isRefreshing).toBe(true);

    await act(async () => {
      await refreshPromise;
    });

    expect(result.current.isRefreshing).toBe(false);
    expect(result.current.data?.pull_requests.length).toBe(0);
    expect(globalThis.fetch).toHaveBeenCalledTimes(2);
  });

  it('polls automatically based on interval', async () => {
    vi.useFakeTimers();

    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => mockDashboardData,
    } as Response);

    const { result } = renderHook(() => useDashboard({ pollIntervalMs: 10000 }));

    // Flush initial fetch
    await act(async () => {
      await vi.advanceTimersByTimeAsync(1);
    });

    expect(fetchSpy).toHaveBeenCalledTimes(1);
    expect(result.current.data).toEqual(mockDashboardData);

    // Advance by interval
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10000);
    });

    expect(fetchSpy).toHaveBeenCalledTimes(2);

    // Advance another interval
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10000);
    });

    expect(fetchSpy).toHaveBeenCalledTimes(3);
  });

  it('cleans up polling timer on unmount', async () => {
    vi.useFakeTimers();

    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => mockDashboardData,
    } as Response);

    const { unmount } = renderHook(() => useDashboard({ pollIntervalMs: 10000 }));

    await act(async () => {
      await vi.advanceTimersByTimeAsync(1);
    });
    expect(fetchSpy).toHaveBeenCalledTimes(1);

    unmount();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(20000);
    });

    expect(fetchSpy).toHaveBeenCalledTimes(1);
  });
});
