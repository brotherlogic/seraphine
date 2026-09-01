import React, { useState } from 'react';
import { ChevronDown, ChevronRight, FolderGit2, Box } from 'lucide-react';
import type { PRSummary } from '../types/dashboard';
import { PRCard } from './PRCard';

export interface RepoSectionProps {
  repo: string;
  prs: PRSummary[];
  defaultExpanded?: boolean;
}

export const RepoSection: React.FC<RepoSectionProps> = ({
  repo,
  prs,
  defaultExpanded = true,
}) => {
  const [isExpanded, setIsExpanded] = useState<boolean>(defaultExpanded);

  const activeContainersCount = prs.filter(
    (pr) => pr.container_state === 'READY' || pr.container_state === 'CREATING'
  ).length;

  return (
    <section className="mb-6 bg-slate-800/40 border border-slate-700/70 rounded-xl overflow-hidden shadow-sm">
      {/* Header / Toggle Button */}
      <button
        type="button"
        onClick={() => setIsExpanded(!isExpanded)}
        aria-label={`Toggle ${repo} section`}
        aria-expanded={isExpanded}
        className="w-full flex items-center justify-between p-4 bg-slate-800/80 hover:bg-slate-800 text-left transition border-b border-slate-700/50"
      >
        <div className="flex items-center gap-3">
          <FolderGit2 className="w-5 h-5 text-blue-400 shrink-0" aria-hidden="true" />
          <h2 className="text-sm sm:text-base font-semibold text-slate-100 font-mono tracking-tight">
            {repo}
          </h2>
        </div>

        <div className="flex items-center gap-3">
          {/* PR Count badge */}
          <span className="px-2.5 py-0.5 text-xs font-semibold rounded-full bg-slate-700/80 text-slate-200 border border-slate-600/60">
            {prs.length} {prs.length === 1 ? 'PR' : 'PRs'}
          </span>

          {/* Active Container Count badge */}
          {activeContainersCount > 0 && (
            <span className="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-0.5 text-xs font-semibold rounded-full bg-emerald-950/60 text-emerald-400 border border-emerald-800/80">
              <Box className="w-3 h-3" />
              <span>
                {activeContainersCount}{' '}
                {activeContainersCount === 1 ? 'Active Container' : 'Active Containers'}
              </span>
            </span>
          )}

          {/* Chevron */}
          <div className="text-slate-400 ml-1">
            {isExpanded ? (
              <ChevronDown className="w-5 h-5" aria-hidden="true" />
            ) : (
              <ChevronRight className="w-5 h-5" aria-hidden="true" />
            )}
          </div>
        </div>
      </button>

      {/* Accordion Content */}
      {isExpanded && (
        <div className="p-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {prs.map((pr) => (
              <PRCard key={`${pr.repo}#${pr.pr_number}`} pr={pr} />
            ))}
          </div>
        </div>
      )}
    </section>
  );
};
