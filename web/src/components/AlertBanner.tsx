import React from 'react';
import { AlertCircle, AlertTriangle, RefreshCw } from 'lucide-react';
import type { FreshnessMetadata } from '../types/dashboard';

export interface AlertBannerProps {
  error: Error | null;
  freshness?: FreshnessMetadata | null;
  onRetry?: () => void;
}

export const AlertBanner: React.FC<AlertBannerProps> = ({
  error,
  freshness,
  onRetry,
}) => {
  if (!error && (!freshness || !freshness.is_stale)) {
    return null;
  }

  if (error) {
    return (
      <aside
        aria-label="Error Alert"
        className="mb-6 rounded-xl border border-rose-800/80 bg-rose-950/40 p-4 text-rose-200 shadow-sm"
      >
        <div className="flex items-start justify-between gap-3">
          <div className="flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-rose-400 mt-0.5 shrink-0" aria-hidden="true" />
            <div>
              <h2 className="text-sm font-semibold text-rose-300">
                Unable to reach Seraphine backend
              </h2>
              <p className="mt-1 text-xs text-rose-200/80">
                {error.message} (Displaying last known cached state)
              </p>
            </div>
          </div>
          {onRetry && (
            <button
              type="button"
              onClick={onRetry}
              className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium bg-rose-900/60 hover:bg-rose-900 border border-rose-700/80 rounded-md text-rose-200 transition shrink-0"
            >
              <RefreshCw className="w-3 h-3" />
              <span>Retry</span>
            </button>
          )}
        </div>
      </aside>
    );
  }

  if (freshness?.is_stale) {
    return (
      <aside
        aria-label="Warning Alert"
        className="mb-6 rounded-xl border border-amber-800/80 bg-amber-950/40 p-4 text-amber-200 shadow-sm"
      >
        <div className="flex items-start gap-3">
          <AlertTriangle className="w-5 h-5 text-amber-400 mt-0.5 shrink-0" aria-hidden="true" />
          <div>
            <h2 className="text-sm font-semibold text-amber-300">
              Stale Data Warning
            </h2>
            <p className="mt-1 text-xs text-amber-200/80">
              {freshness.error_message ||
                'Backend sync with upstream GitHub or devcontainer services is delayed. Data may not reflect the latest changes.'}
            </p>
            {freshness.last_successful_sync && (
              <p className="mt-1 text-[11px] text-amber-300/60">
                Last successful sync: {new Date(freshness.last_successful_sync).toLocaleString()}
              </p>
            )}
          </div>
        </div>
      </aside>
    );
  }

  return null;
};
