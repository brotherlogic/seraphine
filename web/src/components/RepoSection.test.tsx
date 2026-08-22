import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import { RepoSection } from './RepoSection';
import type { PRSummary } from '../types/dashboard';

describe('RepoSection Component', () => {
  const mockPRs: PRSummary[] = [
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
      author: 'brotherlogic-automation',
      commit_count: 2,
      comment_count: 3,
      check_status: 'PENDING',
      has_devcontainer: false,
      container_state: 'NONE',
    },
  ];

  it('renders repository header and summary badges', () => {
    render(<RepoSection repo="brotherlogic/seraphine" prs={mockPRs} />);

    expect(screen.getByText('brotherlogic/seraphine')).toBeInTheDocument();
    expect(screen.getByText('2 PRs')).toBeInTheDocument();
    expect(screen.getByText('1 Active Container')).toBeInTheDocument();
  });

  it('collapses and expands section on header click', async () => {
    const user = userEvent.setup();
    render(<RepoSection repo="brotherlogic/seraphine" prs={mockPRs} defaultExpanded={true} />);

    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();

    const toggleButton = screen.getByRole('button', { name: /toggle brotherlogic\/seraphine/i });
    await user.click(toggleButton);

    expect(screen.queryByText('Scaffold Frontend React Project')).not.toBeInTheDocument();

    await user.click(toggleButton);
    expect(screen.getByText('Scaffold Frontend React Project')).toBeInTheDocument();
  });
});
