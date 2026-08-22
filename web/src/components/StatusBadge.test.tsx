import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { StatusBadge } from './StatusBadge';

describe('StatusBadge Component', () => {
  it('renders CI check status correctly for SUCCESS', () => {
    render(<StatusBadge type="check" status="SUCCESS" />);
    expect(screen.getByText('SUCCESS')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveClass('text-emerald-400');
  });

  it('renders CI check status correctly for FAILURE', () => {
    render(<StatusBadge type="check" status="FAILURE" />);
    expect(screen.getByText('FAILURE')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveClass('text-rose-400');
  });

  it('renders CI check status correctly for PENDING', () => {
    render(<StatusBadge type="check" status="PENDING" />);
    expect(screen.getByText('PENDING')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveClass('text-amber-400');
  });

  it('renders container state correctly for READY', () => {
    render(<StatusBadge type="container" status="READY" containerId="seraphine-180" />);
    expect(screen.getByText(/READY/i)).toBeInTheDocument();
    expect(screen.getByText(/seraphine-180/i)).toBeInTheDocument();
  });

  it('renders container state correctly for NONE', () => {
    render(<StatusBadge type="container" status="NONE" />);
    expect(screen.getByText(/NONE/i)).toBeInTheDocument();
  });

  it('renders container state correctly for CREATING and FAILED', () => {
    const { rerender } = render(<StatusBadge type="container" status="CREATING" />);
    expect(screen.getByText('CREATING')).toBeInTheDocument();

    rerender(<StatusBadge type="container" status="FAILED" />);
    expect(screen.getByText('FAILED')).toBeInTheDocument();
  });
});
