import React from 'react';
import {
  CheckCircle2,
  Clock,
  XCircle,
  HelpCircle,
  Box,
  Loader2,
  AlertTriangle,
  CircleSlash,
} from 'lucide-react';
import type { CheckStatus, ContainerState } from '../types/dashboard';

export interface StatusBadgeProps {
  type: 'check' | 'container';
  status: CheckStatus | ContainerState;
  containerId?: string;
  className?: string;
}

export const StatusBadge: React.FC<StatusBadgeProps> = ({
  type,
  status,
  containerId,
  className = '',
}) => {
  if (type === 'check') {
    const checkStatus = status as CheckStatus;
    switch (checkStatus) {
      case 'SUCCESS':
        return (
          <span
            role="status"
            className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-emerald-950/60 border-emerald-800 text-emerald-400 ${className}`}
          >
            <CheckCircle2 className="w-3.5 h-3.5" aria-hidden="true" />
            <span>SUCCESS</span>
          </span>
        );
      case 'FAILURE':
        return (
          <span
            role="status"
            className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-rose-950/60 border-rose-800 text-rose-400 ${className}`}
          >
            <XCircle className="w-3.5 h-3.5" aria-hidden="true" />
            <span>FAILURE</span>
          </span>
        );
      case 'PENDING':
        return (
          <span
            role="status"
            className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-amber-950/60 border-amber-800 text-amber-400 ${className}`}
          >
            <Clock className="w-3.5 h-3.5 animate-pulse" aria-hidden="true" />
            <span>PENDING</span>
          </span>
        );
      case 'UNKNOWN':
      default:
        return (
          <span
            role="status"
            className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-slate-800 border-slate-700 text-slate-400 ${className}`}
          >
            <HelpCircle className="w-3.5 h-3.5" aria-hidden="true" />
            <span>UNKNOWN</span>
          </span>
        );
    }
  }

  const containerState = status as ContainerState;
  switch (containerState) {
    case 'READY':
      return (
        <span
          role="status"
          className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-emerald-950/60 border-emerald-800 text-emerald-400 ${className}`}
        >
          <Box className="w-3.5 h-3.5" aria-hidden="true" />
          <span>READY</span>
          {containerId && (
            <span className="font-mono text-[11px] opacity-80 pl-1 border-l border-emerald-800/80">
              {containerId}
            </span>
          )}
        </span>
      );
    case 'CREATING':
      return (
        <span
          role="status"
          className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-sky-950/60 border-sky-800 text-sky-400 ${className}`}
        >
          <Loader2 className="w-3.5 h-3.5 animate-spin" aria-hidden="true" />
          <span>CREATING</span>
          {containerId && (
            <span className="font-mono text-[11px] opacity-80 pl-1 border-l border-sky-800/80">
              {containerId}
            </span>
          )}
        </span>
      );
    case 'FAILED':
      return (
        <span
          role="status"
          className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full border bg-rose-950/60 border-rose-800 text-rose-400 ${className}`}
        >
          <AlertTriangle className="w-3.5 h-3.5" aria-hidden="true" />
          <span>FAILED</span>
          {containerId && (
            <span className="font-mono text-[11px] opacity-80 pl-1 border-l border-rose-800/80">
              {containerId}
            </span>
          )}
        </span>
      );
    case 'NONE':
      return (
        <span
          role="status"
          className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-full border bg-slate-800/80 border-slate-700/80 text-slate-400 ${className}`}
        >
          <CircleSlash className="w-3.5 h-3.5" aria-hidden="true" />
          <span>NONE</span>
        </span>
      );
    case 'UNKNOWN':
    default:
      return (
        <span
          role="status"
          className={`inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-full border bg-slate-800 border-slate-700 text-slate-400 ${className}`}
        >
          <HelpCircle className="w-3.5 h-3.5" aria-hidden="true" />
          <span>UNKNOWN</span>
        </span>
      );
  }
};
