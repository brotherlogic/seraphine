import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import App from './App';
import * as dashboardHook from './hooks/useDashboard';

describe('App Component', () => {
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

  it('renders dashboard with pull requests', () => {
    vi.spyOn(dashboardHook, 'useDashboard').mockReturnValue({
      data: {
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
        ],
        freshness: {
          last_successful_sync: '2026-08-20T17:00:00Z',
          last_attempted_sync: '2026-08-20T17:00:00Z',
          is_stale: false,
        },
      },
      loading: false,
      error: null,
      isRefreshing: false,
      lastUpdated: new Date(),
      refresh: vi.fn(),
    });

    render(<App />);
    expect(screen.getByText('Seraphine Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();
    expect(screen.getByText('#180')).toBeInTheDocument();
  });
});
