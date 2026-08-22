import React, { useState } from 'react';
import { ExternalLink, GitCommit, MessageSquare, User } from 'lucide-react';
import type { PRSummary } from '../types/dashboard';
import { StatusBadge } from './StatusBadge';

export interface PRCardProps {
  pr: PRSummary;
}

export const PRCard: React.FC<PRCardProps> = ({ pr }) => {
  const [avatarError, setAvatarError] = useState(false);
  const prUrl = `https://github.com/${pr.repo}/pull/${pr.pr_number}`;
  const avatarUrl = `https://github.com/${pr.author}.png?size=40`;

  return (
    <div className="flex flex-col justify-between bg-slate-800/80 border border-slate-700/80 hover:border-slate-600 rounded-xl p-4 transition shadow-sm hover:shadow group">
      <div>
        {/* Top bar: PR number & GitHub Link */}
        <div className="flex items-center justify-between gap-2 mb-2">
          <span className="font-mono text-xs font-semibold text-blue-400 bg-blue-950/40 border border-blue-800/50 px-2 py-0.5 rounded">
            #{pr.pr_number}
          </span>
          <a
            href={prUrl}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={`View Pull Request #${pr.pr_number} on GitHub`}
            className="inline-flex items-center gap-1 text-xs text-slate-400 hover:text-blue-400 transition"
          >
            <span>GitHub</span>
            <ExternalLink className="w-3 h-3" />
          </a>
        </div>

        {/* PR Title */}
        <h3 className="text-sm font-medium text-slate-100 line-clamp-2 mb-3 group-hover:text-blue-200 transition">
          <a href={prUrl} target="_blank" rel="noopener noreferrer">
            {pr.title}
          </a>
        </h3>
      </div>

      <div className="space-y-3 pt-3 border-t border-slate-700/50">
        {/* Author info & Commit/Comment metrics */}
        <div className="flex items-center justify-between text-xs text-slate-400">
          <div className="flex items-center gap-1.5 min-w-0">
            {avatarError ? (
              <div className="w-5 h-5 rounded-full bg-slate-700 flex items-center justify-center text-slate-300 shrink-0">
                <User className="w-3 h-3" />
              </div>
            ) : (
              <img
                src={avatarUrl}
                alt={pr.author}
                onError={() => setAvatarError(true)}
                className="w-5 h-5 rounded-full bg-slate-700 shrink-0"
              />
            )}
            <span className="truncate">@{pr.author}</span>
          </div>

          <div className="flex items-center gap-2.5 shrink-0">
            <span
              className="inline-flex items-center gap-1"
              title={`${pr.commit_count} commits`}
            >
              <GitCommit className="w-3.5 h-3.5 text-slate-400" />
              <span>{pr.commit_count} {pr.commit_count === 1 ? 'commit' : 'commits'}</span>
            </span>
            <span
              className="inline-flex items-center gap-1"
              title={`${pr.comment_count} comments`}
            >
              <MessageSquare className="w-3.5 h-3.5 text-slate-400" />
              <span>{pr.comment_count} {pr.comment_count === 1 ? 'comment' : 'comments'}</span>
            </span>
          </div>
        </div>

        {/* Status badges */}
        <div className="flex flex-wrap items-center justify-between gap-2 pt-1">
          <StatusBadge type="check" status={pr.check_status} />
          <StatusBadge
            type="container"
            status={pr.container_state}
            containerId={pr.container_id}
          />
        </div>
      </div>
    </div>
  );
};
