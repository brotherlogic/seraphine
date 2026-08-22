import React from 'react';
import { Inbox, RotateCcw } from 'lucide-react';

export interface EmptyStateProps {
  title?: string;
  description?: string;
  icon?: React.ReactNode;
  onReset?: () => void;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  title = 'No pull requests tracked',
  description = 'There are currently no active pull requests enrolled in Seraphine.',
  icon,
  onReset,
}) => {
  return (
    <div className="flex flex-col items-center justify-center p-12 text-center bg-slate-800/40 border border-slate-800 rounded-2xl">
      <div className="p-3 mb-4 rounded-full bg-slate-800 border border-slate-700 text-slate-400">
        {icon || <Inbox className="w-8 h-8" aria-hidden="true" />}
      </div>
      <h3 className="text-base font-semibold text-slate-200">{title}</h3>
      <p className="max-w-md mt-1 text-sm text-slate-400">{description}</p>
      {onReset && (
        <button
          type="button"
          onClick={onReset}
          className="inline-flex items-center gap-2 px-4 py-2 mt-5 text-xs font-semibold text-slate-200 bg-slate-800 hover:bg-slate-700 active:bg-slate-600 rounded-lg border border-slate-700 transition"
        >
          <RotateCcw className="w-3.5 h-3.5" />
          <span>Clear Filters</span>
        </button>
      )}
    </div>
  );
};
