import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { Header } from './Header';

describe('Header Component', () => {
  it('renders title and branding description', () => {
    render(
      <Header
        lastUpdated={new Date('2026-08-22T09:00:00Z')}
        loading={false}
        isRefreshing={false}
        onRefresh={vi.fn()}
      />
    );

    expect(screen.getByText('Seraphine Dashboard')).toBeInTheDocument();
    expect(screen.getByText(/multi-repository/i)).toBeInTheDocument();
  });

  it('renders relative last updated timestamp', () => {
    render(
      <Header
        lastUpdated={new Date(Date.now() - 5000)}
        loading={false}
        isRefreshing={false}
        onRefresh={vi.fn()}
      />
    );

    expect(screen.getByText(/Updated just now/i)).toBeInTheDocument();
  });

  it('calls onRefresh when refresh button is clicked', async () => {
    const handleRefresh = vi.fn();
    const user = userEvent.setup();

    render(
      <Header
        lastUpdated={new Date()}
        loading={false}
        isRefreshing={false}
        onRefresh={handleRefresh}
      />
    );

    const refreshBtn = screen.getByRole('button', { name: /refresh/i });
    await user.click(refreshBtn);
    expect(handleRefresh).toHaveBeenCalledTimes(1);
  });

  it('disables refresh button and shows spinner during refresh', () => {
    render(
      <Header
        lastUpdated={new Date()}
        loading={false}
        isRefreshing={true}
        onRefresh={vi.fn()}
      />
    );

    const refreshBtn = screen.getByRole('button', { name: /refreshing/i });
    expect(refreshBtn).toBeDisabled();
    expect(screen.getByTestId('refresh-spinner')).toBeInTheDocument();
  });
});
