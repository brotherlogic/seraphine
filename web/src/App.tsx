import React from 'react';
import { useDashboard } from './hooks/useDashboard';

export const App: React.FC = () => {
  const { data, loading, error, isRefreshing, lastUpdated, refresh } = useDashboard();

  return (
    <div className="min-h-screen bg-slate-900 text-slate-100 p-6">
      <header className="max-w-7xl mx-auto flex items-center justify-between border-b border-slate-800 pb-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-slate-100 flex items-center gap-2">
            Seraphine Dashboard
          </h1>
          <p className="text-sm text-slate-400">
            Real-time multi-repository pull requests & devcontainer state
          </p>
        </div>

        <div className="flex items-center gap-4">
          {lastUpdated && (
            <span className="text-xs text-slate-400">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </span>
          )}
          <button
            onClick={() => refresh()}
            disabled={loading || isRefreshing}
            className="px-3 py-1.5 text-sm font-medium bg-slate-800 hover:bg-slate-700 disabled:opacity-50 text-slate-200 rounded-lg border border-slate-700 transition"
          >
            {isRefreshing ? 'Refreshing...' : 'Refresh'}
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto">
        {loading && (
          <div className="flex items-center justify-center py-12 text-slate-400">
            <p>Loading dashboard state...</p>
          </div>
        )}

        {error && (
          <div className="bg-red-950/50 border border-red-800/80 rounded-lg p-4 mb-6 text-red-200 text-sm">
            <p className="font-semibold">Failed to load dashboard data</p>
            <p className="text-red-300 mt-1">{error.message}</p>
          </div>
        )}

        {!loading && data && (
          <div className="space-y-4">
            <div className="flex items-center justify-between text-xs text-slate-400">
              <span>Active Pull Requests: {data.pull_requests.length}</span>
              <span>
                Sync status: {data.freshness.is_stale ? 'Stale' : 'Healthy'}
              </span>
            </div>

            {data.pull_requests.length === 0 ? (
              <div className="bg-slate-800/40 border border-slate-800 rounded-lg p-8 text-center text-slate-400">
                No active pull requests tracked.
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {data.pull_requests.map((pr) => (
                  <div
                    key={`${pr.repo}#${pr.pr_number}`}
                    className="bg-slate-800/60 border border-slate-700/60 rounded-lg p-4 space-y-2"
                  >
                    <div className="flex items-center justify-between text-xs text-slate-400">
                      <span className="font-mono">{pr.repo}</span>
                      <span>#{pr.pr_number}</span>
                    </div>
                    <h3 className="font-medium text-slate-100 line-clamp-2">
                      {pr.title}
                    </h3>
                    <div className="flex items-center justify-between pt-2 text-xs border-t border-slate-700/40">
                      <span className="text-slate-400">@{pr.author}</span>
                      <span className="px-2 py-0.5 rounded bg-slate-700 text-slate-200">
                        {pr.check_status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </main>
    </div>
  );
};

export default App;
