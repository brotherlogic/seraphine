import React, { useState, useMemo } from 'react';
import { useDashboard } from './hooks/useDashboard';
import { Header } from './components/Header';
import { AlertBanner } from './components/AlertBanner';
import { FilterBar } from './components/FilterBar';
import { RepoSection } from './components/RepoSection';
import { EmptyState } from './components/EmptyState';
import { Loader2, SearchX } from 'lucide-react';
import type { PRSummary } from './types/dashboard';

export const App: React.FC = () => {
  const { data, loading, error, isRefreshing, lastUpdated, refresh } = useDashboard();

  const [searchQuery, setSearchQuery] = useState<string>('');
  const [selectedRepo, setSelectedRepo] = useState<string>('ALL');
  const [selectedContainerState, setSelectedContainerState] = useState<string>('ALL');

  // Derive unique repositories
  const repos = useMemo(() => {
    if (!data?.pull_requests) return [];
    const set = new Set<string>();
    data.pull_requests.forEach((pr) => set.add(pr.repo));
    return Array.from(set).sort();
  }, [data?.pull_requests]);

  // Client-side filtering and search engine
  const filteredPRs = useMemo(() => {
    if (!data?.pull_requests) return [];

    const query = searchQuery.trim().toLowerCase();

    return data.pull_requests.filter((pr: PRSummary) => {
      // 1. Repository Filter
      if (selectedRepo !== 'ALL' && pr.repo !== selectedRepo) {
        return false;
      }

      // 2. Container State Filter
      if (
        selectedContainerState !== 'ALL' &&
        pr.container_state !== selectedContainerState
      ) {
        return false;
      }

      // 3. Search Query Filter (Title, #Number, Author, Repo)
      if (query) {
        const titleMatch = pr.title.toLowerCase().includes(query);
        const numberMatch =
          pr.pr_number.toString().includes(query) ||
          `#${pr.pr_number}`.includes(query);
        const authorMatch = pr.author.toLowerCase().includes(query);
        const repoMatch = pr.repo.toLowerCase().includes(query);

        if (!titleMatch && !numberMatch && !authorMatch && !repoMatch) {
          return false;
        }
      }

      return true;
    });
  }, [data?.pull_requests, selectedRepo, selectedContainerState, searchQuery]);

  // Group filtered PRs by repository
  const groupedPRs = useMemo(() => {
    const map = new Map<string, PRSummary[]>();
    filteredPRs.forEach((pr) => {
      const list = map.get(pr.repo) || [];
      list.push(pr);
      map.set(pr.repo, list);
    });
    return map;
  }, [filteredPRs]);

  const handleResetFilters = () => {
    setSearchQuery('');
    setSelectedRepo('ALL');
    setSelectedContainerState('ALL');
  };

  const isFiltered =
    searchQuery !== '' || selectedRepo !== 'ALL' || selectedContainerState !== 'ALL';
  const totalCount = data?.pull_requests.length ?? 0;
  const filteredCount = filteredPRs.length;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 antialiased selection:bg-blue-600 selection:text-white">
      <div className="px-4 sm:px-6 lg:px-8 py-6">
        <Header
          lastUpdated={lastUpdated}
          loading={loading}
          isRefreshing={isRefreshing}
          onRefresh={refresh}
        />

        <main className="max-w-7xl mx-auto">
          {/* Non-blocking Alert Banner for Backend / Sync Issues */}
          <AlertBanner
            error={error}
            freshness={data?.freshness}
            onRetry={refresh}
          />

          {/* Initial Loading Spinner */}
          {loading && !data && (
            <div className="flex flex-col items-center justify-center py-24 text-slate-400 gap-3">
              <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
              <p className="text-sm font-medium">Loading dashboard state...</p>
            </div>
          )}

          {/* Dashboard Content */}
          {data && (
            <>
              {/* Filter and Search Bar */}
              <FilterBar
                searchQuery={searchQuery}
                onSearchChange={setSearchQuery}
                selectedRepo={selectedRepo}
                repos={repos}
                onRepoChange={setSelectedRepo}
                selectedContainerState={selectedContainerState}
                onContainerStateChange={setSelectedContainerState}
                totalCount={totalCount}
                filteredCount={filteredCount}
              />

              {/* Empty States or Repositories List */}
              {totalCount === 0 ? (
                <EmptyState
                  title="No active pull requests tracked"
                  description="Seraphine is currently not tracking any active pull requests across enrolled repositories."
                />
              ) : filteredCount === 0 ? (
                <EmptyState
                  title="No matching pull requests"
                  description="No pull requests match the active search query or filter combination."
                  icon={<SearchX className="w-8 h-8" aria-hidden="true" />}
                  onReset={isFiltered ? handleResetFilters : undefined}
                />
              ) : (
                <div className="space-y-4">
                  {Array.from(groupedPRs.entries()).map(([repoName, prs]) => (
                    <RepoSection
                      key={repoName}
                      repo={repoName}
                      prs={prs}
                      defaultExpanded={true}
                    />
                  ))}
                </div>
              )}
            </>
          )}
        </main>
      </div>
    </div>
  );
};

export default App;
