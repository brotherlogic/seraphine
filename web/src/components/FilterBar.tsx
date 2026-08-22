import React from 'react';
import { Search, X, Filter } from 'lucide-react';
import type { ContainerState } from '../types/dashboard';

export interface FilterBarProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  selectedRepo: string;
  repos: string[];
  onRepoChange: (repo: string) => void;
  selectedContainerState: string;
  onContainerStateChange: (state: string) => void;
  totalCount?: number;
  filteredCount?: number;
}

const CONTAINER_FILTER_OPTIONS: Array<{ label: string; value: string }> = [
  { label: 'ALL', value: 'ALL' },
  { label: 'READY', value: 'READY' as ContainerState },
  { label: 'CREATING', value: 'CREATING' as ContainerState },
  { label: 'FAILED', value: 'FAILED' as ContainerState },
  { label: 'NONE', value: 'NONE' as ContainerState },
];

export const FilterBar: React.FC<FilterBarProps> = ({
  searchQuery,
  onSearchChange,
  selectedRepo,
  repos,
  onRepoChange,
  selectedContainerState,
  onContainerStateChange,
  totalCount,
  filteredCount,
}) => {
  return (
    <div className="bg-slate-800/60 border border-slate-700/80 rounded-xl p-4 mb-6 space-y-4 shadow-sm backdrop-blur">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        {/* Search Input */}
        <div className="relative flex-1">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-400">
            <Search className="w-4 h-4" aria-hidden="true" />
          </div>
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="Search pull requests by title, #number, author..."
            className="w-full pl-9 pr-8 py-2 bg-slate-900/80 border border-slate-700 text-slate-100 placeholder-slate-400 text-xs sm:text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
          />
          {searchQuery && (
            <button
              type="button"
              onClick={() => onSearchChange('')}
              aria-label="Clear search"
              className="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400 hover:text-slate-200 transition"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* Repository Dropdown */}
        <div className="w-full md:w-64 shrink-0">
          <label htmlFor="repo-selector" className="sr-only">
            Filter by Repository
          </label>
          <select
            id="repo-selector"
            value={selectedRepo}
            onChange={(e) => onRepoChange(e.target.value)}
            className="w-full py-2 px-3 bg-slate-900/80 border border-slate-700 text-slate-100 text-xs sm:text-sm rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition cursor-pointer"
          >
            <option value="ALL">All Repositories ({repos.length})</option>
            {repos.map((repo) => (
              <option key={repo} value={repo}>
                {repo}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Container State Filter Pills */}
      <div className="flex flex-wrap items-center justify-between gap-3 pt-3 border-t border-slate-700/60">
        <div className="flex flex-wrap items-center gap-2">
          <div className="flex items-center gap-1.5 text-xs text-slate-400 mr-1">
            <Filter className="w-3.5 h-3.5" />
            <span>Container:</span>
          </div>
          {CONTAINER_FILTER_OPTIONS.map((opt) => {
            const isSelected = selectedContainerState === opt.value;
            return (
              <button
                key={opt.value}
                type="button"
                onClick={() => onContainerStateChange(opt.value)}
                className={`px-3 py-1 text-xs font-medium rounded-full border transition ${
                  isSelected
                    ? 'bg-blue-600 border-blue-500 text-white shadow-sm'
                    : 'bg-slate-900/60 border-slate-700/80 text-slate-400 hover:text-slate-200 hover:bg-slate-800'
                }`}
              >
                {opt.label}
              </button>
            );
          })}
        </div>

        {/* Count Summary */}
        {totalCount !== undefined && filteredCount !== undefined && (
          <div className="text-xs text-slate-400 font-medium">
            Showing <span className="text-slate-200 font-semibold">{filteredCount}</span> of{' '}
            <span className="text-slate-200 font-semibold">{totalCount}</span> PRs
          </div>
        )}
      </div>
    </div>
  );
};
