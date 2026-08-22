import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import App from './App';
import * as dashboardHook from './hooks/useDashboard';
import type { DashboardState } from './types/dashboard';

describe('App Component Integration', () => {
  const mockDashboardData: DashboardState = {
    pull_requests: [
      {
        repo: 'brotherlogic/seraphine',
        pr_number: 180,
        title: 'Scaffold Frontend React Project',
        author: 'brotherlogic-automation',
        commit_count: 1,
        comment_count: 0,
        check_status: 'SUCCESS',
        has_devcontainer: true,
        container_id: 'seraphine-180',
        container_state: 'READY',
      },
      {
        repo: 'brotherlogic/seraphine',
        pr_number: 181,
        title: 'Implement Frontend Dashboard UI Components',
        author: 'developer-one',
        commit_count: 3,
        comment_count: 2,
        check_status: 'PENDING',
        has_devcontainer: false,
        container_state: 'NONE',
      },
      {
        repo: 'brotherlogic/home',
        pr_number: 42,
        title: 'Fix thermostat temperature sync',
        author: 'brotherlogic',
        commit_count: 5,
        comment_count: 1,
        check_status: 'FAILURE',
        has_devcontainer: true,
        container_id: 'home-42',
        container_state: 'FAILED',
      },
    ],
    freshness: {
      last_successful_sync: '2026-08-22T09:00:00Z',
      last_attempted_sync: '2026-08-22T09:00:00Z',
      is_stale: false,
    },
  };

  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders loading state', () => {
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: null,
      loading: true,
      error: null,
      isRefreshing: false,
      lastUpdated: null,
      refresh: vi.fn(),
    });

    render(<App />);
    expect(screen.getByText(/loading dashboard state/i)).toBeInTheDocument();
  });

  it('renders full dashboard layout with repositories and cards', () => {
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: mockDashboardData,
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    expect(screen.getByText('Seraphine Dashboard')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'brotherlogic/seraphine' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'brotherlogic/home' })).toBeInTheDocument();
    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();
    expect(screen.getByText('Fix thermostat temperature sync')).toBeInTheDocument();
  });

  it('filters pull requests by search query', async () => {
    const user = userEvent.setup();
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: mockDashboardData,
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    const searchInput = screen.getByPlaceholderText(/search pull requests/i);
    await user.type(searchInput, 'thermostat');

    expect(screen.getByText('Fix thermostat temperature sync')).toBeInTheDocument();
    expect(screen.queryByText('Scaffold Frontend React Project')).not.toBeInTheDocument();
  });

  it('filters pull requests by repository selector', async () => {
    const user = userEvent.setup();
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: mockDashboardData,
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    const repoSelect = screen.getByRole('combobox');
    await user.selectOptions(repoSelect, 'brotherlogic/home');

    expect(screen.getByText('Fix thermostat temperature sync')).toBeInTheDocument();
    expect(screen.queryByRole('heading', { name: 'brotherlogic/seraphine' })).not.toBeInTheDocument();
  });

  it('filters pull requests by container state filter pills', async () => {
    const user = userEvent.setup();
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: mockDashboardData,
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    const readyPill = screen.getByRole('button', { name: /^ready$/i });
    await user.click(readyPill);

    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();
    expect(screen.queryByText('Fix thermostat temperature sync')).not.toBeInTheDocument();
    expect(screen.queryByText('Implement Frontend Dashboard UI Components')).not.toBeInTheDocument();
  });

  it('shows empty filter state and resets filters on clear button click', async () => {
    const user = userEvent.setup();
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: mockDashboardData,
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    const searchInput = screen.getByPlaceholderText(/search pull requests/i);
    await user.type(searchInput, 'nonexistent-query-xyz');

    expect(screen.getByText(/no matching pull requests/i)).toBeInTheDocument();
    const clearBtn = screen.getByRole('button', { name: /clear filters/i });
    await user.click(clearBtn);

    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();
  });
});
