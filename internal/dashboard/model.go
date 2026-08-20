package dashboard

import (
	"context"
	"time"

	"github.com/brotherlogic/seraphine/internal/github"
)

// CheckStatus aliases GitHub check status for dashboard consumers.
type CheckStatus = github.CheckStatus

const (
	CheckStatusSuccess CheckStatus = github.CheckStatusSuccess
	CheckStatusPending CheckStatus = github.CheckStatusPending
	CheckStatusFailure CheckStatus = github.CheckStatusFailure
	CheckStatusUnknown CheckStatus = github.CheckStatusUnknown
)

// ContainerState represents the lifecycle status of a PR's devcontainer.
type ContainerState string

const (
	ContainerStateCreating ContainerState = "CREATING"
	ContainerStateReady    ContainerState = "READY"
	ContainerStateFailed   ContainerState = "FAILED"
	ContainerStateNone     ContainerState = "NONE"
	ContainerStateUnknown  ContainerState = "UNKNOWN"
)

// PRSummary represents the consolidated dashboard view of a single pull request.
type PRSummary struct {
	Repo            string         `json:"repo"`
	PRNumber        int32          `json:"pr_number"`
	Title           string         `json:"title"`
	Author          string         `json:"author"`
	CommitCount     int            `json:"commit_count"`
	CommentCount    int            `json:"comment_count"`
	CheckStatus     CheckStatus    `json:"check_status"`
	HasDevcontainer bool           `json:"has_devcontainer"`
	ContainerID     string         `json:"container_id,omitempty"`
	ContainerState  ContainerState `json:"container_state"`
}

// FreshnessMetadata tracks sync timestamps and staleness indicators.
type FreshnessMetadata struct {
	LastSuccessfulSync time.Time `json:"last_successful_sync"`
	LastAttemptedSync  time.Time `json:"last_attempted_sync"`
	IsStale            bool      `json:"is_stale"`
	ErrorMessage       string    `json:"error_message,omitempty"`
}

// DashboardState is the root state object containing active PRs and sync metadata.
type DashboardState struct {
	PullRequests []PRSummary       `json:"pull_requests"`
	Freshness    FreshnessMetadata `json:"freshness"`
}

// Service defines the interface for querying dashboard state and running background workers.
type Service interface {
	GetDashboardState(ctx context.Context) (*DashboardState, error)
	RunWorker(ctx context.Context, interval time.Duration)
}
