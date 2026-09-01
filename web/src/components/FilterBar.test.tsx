import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { FilterBar } from './FilterBar';

describe('FilterBar Component', () => {
  const repos = ['brotherlogic/seraphine', 'brotherlogic/home'];

  it('renders search input, repo selector, and container filter pills', () => {
    render(
      <FilterBar
        searchQuery=""
        onSearchChange={vi.fn()}
        selectedRepo="ALL"
        repos={repos}
        onRepoChange={vi.fn()}
        selectedContainerState="ALL"
        onContainerStateChange={vi.fn()}
        totalCount={10}
        filteredCount={10}
      />
    );

    expect(screen.getByPlaceholderText(/search pull requests/i)).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^all$/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^ready$/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^creating$/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^failed$/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^none$/i })).toBeInTheDocument();
  });

  it('handles search input changes and clear button', async () => {
    const onSearchChange = vi.fn();
    const user = userEvent.setup();

    const { rerender } = render(
      <FilterBar
        searchQuery=""
        onSearchChange={onSearchChange}
        selectedRepo="ALL"
        repos={repos}
        onRepoChange={vi.fn()}
        selectedContainerState="ALL"
        onContainerStateChange={vi.fn()}
        totalCount={10}
        filteredCount={10}
      />
    );

    const input = screen.getByPlaceholderText(/search pull requests/i);
    await user.type(input, 'dashboard');
    expect(onSearchChange).toHaveBeenCalled();

    rerender(
      <FilterBar
        searchQuery="dashboard"
        onSearchChange={onSearchChange}
        selectedRepo="ALL"
        repos={repos}
        onRepoChange={vi.fn()}
        selectedContainerState="ALL"
        onContainerStateChange={vi.fn()}
        totalCount={10}
        filteredCount={10}
      />
    );

    const clearBtn = screen.getByLabelText(/clear search/i);
    await user.click(clearBtn);
    expect(onSearchChange).toHaveBeenCalledWith('');
  });

  it('handles repo selection and container status pill toggle', async () => {
    const onRepoChange = vi.fn();
    const onContainerStateChange = vi.fn();
    const user = userEvent.setup();

    render(
      <FilterBar
        searchQuery=""
        onSearchChange={vi.fn()}
        selectedRepo="ALL"
        repos={repos}
        onRepoChange={onRepoChange}
        selectedContainerState="ALL"
        onContainerStateChange={onContainerStateChange}
        totalCount={5}
        filteredCount={2}
      />
    );

    const select = screen.getByRole('combobox');
    await user.selectOptions(select, 'brotherlogic/seraphine');
    expect(onRepoChange).toHaveBeenCalledWith('brotherlogic/seraphine');

    const readyPill = screen.getByRole('button', { name: /^ready$/i });
    await user.click(readyPill);
    expect(onContainerStateChange).toHaveBeenCalledWith('READY');
  });
});
