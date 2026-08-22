export type CheckStatus = 'SUCCESS' | 'PENDING' | 'FAILURE' | 'UNKNOWN';

export type ContainerState = 'CREATING' | 'READY' | 'FAILED' | 'NONE' | 'UNKNOWN';

export interface PRSummary {
  repo: string;
  pr_number: number;
  title: string;
  author: string;
  commit_count: number;
  comment_count: number;
  check_status: CheckStatus;
  has_devcontainer: boolean;
  container_id?: string;
  container_state: ContainerState;
}

export interface FreshnessMetadata {
  last_successful_sync: string;
  last_attempted_sync: string;
  is_stale: boolean;
  error_message?: string;
}

export interface DashboardState {
  pull_requests: PRSummary[];
  freshness: FreshnessMetadata;
}
