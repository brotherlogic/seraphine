import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { AlertBanner } from './AlertBanner';

describe('AlertBanner Component', () => {
  it('renders nothing when there is no error and sync is healthy', () => {
    const { container } = render(
      <AlertBanner
        error={null}
        freshness={{
          last_successful_sync: '2026-08-22T09:00:00Z',
          last_attempted_sync: '2026-08-22T09:00:00Z',
          is_stale: false,
        }}
      />
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders error banner when error prop is provided', async () => {
    const handleRetry = vi.fn();
    const user = userEvent.setup();

    render(
      <AlertBanner
        error={new Error('Connection refused to :8080')}
        onRetry={handleRetry}
      />
    );

    expect(screen.getByText(/Unable to reach Seraphine backend/i)).toBeInTheDocument();
    expect(screen.getByText(/Connection refused to :8080/i)).toBeInTheDocument();

    const retryBtn = screen.getByRole('button', { name: /retry/i });
    await user.click(retryBtn);
    expect(handleRetry).toHaveBeenCalledTimes(1);
  });

  it('renders stale warning banner when data is stale', () => {
    render(
      <AlertBanner
        error={null}
        freshness={{
          last_successful_sync: '2026-08-22T08:30:00Z',
          last_attempted_sync: '2026-08-22T09:00:00Z',
          is_stale: true,
          error_message: 'GitHub rate limit exceeded',
        }}
      />
    );

    expect(screen.getByText(/Stale Data Warning/i)).toBeInTheDocument();
    expect(screen.getByText(/GitHub rate limit exceeded/i)).toBeInTheDocument();
  });
});
