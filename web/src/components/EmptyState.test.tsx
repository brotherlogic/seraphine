import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { EmptyState } from './EmptyState';

describe('EmptyState Component', () => {
  it('renders default empty state without action button', () => {
    render(
      <EmptyState
        title="No repositories enrolled"
        description="There are currently no active pull requests being tracked."
      />
    );

    expect(screen.getByText('No repositories enrolled')).toBeInTheDocument();
    expect(screen.getByText(/no active pull requests/i)).toBeInTheDocument();
  });

  it('renders clear filters button when onReset is provided', async () => {
    const handleReset = vi.fn();
    const user = userEvent.setup();

    render(
      <EmptyState
        title="No matching pull requests"
        description="Try adjusting your search query or filters."
        onReset={handleReset}
      />
    );

    expect(screen.getByText('No matching pull requests')).toBeInTheDocument();
    const resetBtn = screen.getByRole('button', { name: /clear filters/i });
    await user.click(resetBtn);
    expect(handleReset).toHaveBeenCalledTimes(1);
  });
});
