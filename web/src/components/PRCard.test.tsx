import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { PRCard } from './PRCard';
import type { PRSummary } from '../types/dashboard';

describe('PRCard Component', () => {
  const mockPR: PRSummary = {
    repo: 'brotherlogic/seraphine',
    pr_number: 181,
    title: 'Implement Frontend Dashboard UI Components',
    author: 'brotherlogic-automation',
    commit_count: 3,
    comment_count: 5,
    check_status: 'SUCCESS',
    has_devcontainer: true,
    container_id: 'seraphine-181',
    container_state: 'READY',
  };

  it('renders title, author, pr number, and external github link', () => {
    render(<PRCard pr={mockPR} />);

    expect(screen.getByText('Implement Frontend Dashboard UI Components')).toBeInTheDocument();
    expect(screen.getByText('#181')).toBeInTheDocument();
    expect(screen.getByText('@brotherlogic-automation')).toBeInTheDocument();

    const link = screen.getByRole('link', { name: /view pull request/i });
    expect(link).toHaveAttribute('href', 'https://github.com/brotherlogic/seraphine/pull/181');
    expect(link).toHaveAttribute('target', '_blank');
  });

  it('renders commit count, comment count, and status badges', () => {
    render(<PRCard pr={mockPR} />);

    expect(screen.getByText(/3 commits/i)).toBeInTheDocument();
    expect(screen.getByText(/5 comments/i)).toBeInTheDocument();
    expect(screen.getByText('SUCCESS')).toBeInTheDocument();
    expect(screen.getByText(/READY/i)).toBeInTheDocument();
    expect(screen.getByText(/seraphine-181/i)).toBeInTheDocument();
  });
});
