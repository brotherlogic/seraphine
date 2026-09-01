import React, { useState, useEffect } from 'react';
import { RotateCw, LayoutDashboard } from 'lucide-react';

export interface HeaderProps {
  lastUpdated: Date | null;
  loading: boolean;
  isRefreshing: boolean;
  onRefresh: () => void;
}

function getRelativeTimeString(date: Date): string {
  const seconds = Math.max(0, Math.floor((Date.now() - date.getTime()) / 1000));
  if (seconds < 10) {
    return 'just now';
  }
  if (seconds < 60) {
    return `${seconds}s ago`;
  }
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) {
    return `${minutes}m ago`;
  }
  const hours = Math.floor(minutes / 60);
  return `${hours}h ago`;
}

export const Header: React.FC<HeaderProps> = ({
  lastUpdated,
  loading,
  isRefreshing,
  onRefresh,
}) => {
  const [, setTick] = useState(0);

  // Re-calculate relative time string every 5 seconds
  useEffect(() => {
    const timer = setInterval(() => {
      setTick((t) => t + 1);
    }, 5000);
    return () => clearInterval(timer);
  }, []);

  const relativeTime = lastUpdated ? getRelativeTimeString(lastUpdated) : null;

  return (
    <header className="border-b border-slate-800 bg-slate-900/90 backdrop-blur sticky top-0 z-20 pb-4 mb-6">
      <div className="max-w-7xl mx-auto flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="p-2 rounded-xl bg-blue-600/10 border border-blue-500/20 text-blue-400">
            <LayoutDashboard className="w-6 h-6" aria-hidden="true" />
          </div>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-slate-100 flex items-center gap-2">
              Seraphine Dashboard
            </h1>
            <p className="text-xs text-slate-400">
              Real-time multi-repository pull requests & devcontainer state
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3 self-end sm:self-auto">
          {lastUpdated && (
            <div className="text-right">
              <span
                className="text-xs text-slate-400"
                title={lastUpdated.toLocaleString()}
              >
                Updated {relativeTime}
              </span>
            </div>
          )}
          <button
            type="button"
            onClick={onRefresh}
            disabled={loading || isRefreshing}
            className="inline-flex items-center gap-2 px-3.5 py-1.5 text-xs font-semibold bg-slate-800 hover:bg-slate-700 active:bg-slate-600 disabled:opacity-50 text-slate-200 rounded-lg border border-slate-700 transition shadow-sm"
          >
            <RotateCw
              data-testid="refresh-spinner"
              className={`w-3.5 h-3.5 ${isRefreshing ? 'animate-spin text-blue-400' : 'text-slate-400'}`}
              aria-hidden="true"
            />
            <span>{isRefreshing ? 'Refreshing...' : 'Refresh'}</span>
          </button>
        </div>
      </div>
    </header>
  );
};
