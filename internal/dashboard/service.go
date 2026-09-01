package dashboard

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	manager_pb "github.com/brotherlogic/devcontainer-manager/proto"
	pstore_client "github.com/brotherlogic/pstore/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
)

// DevcontainerClient interface required by dashboard to list containers.
type DevcontainerClient interface {
	List(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error)
}

// Option allows configuring optional parameters for dashboardService.
type Option func(*dashboardService)

// WithWorkers sets the concurrency worker pool size for fetching PR details and checks.
func WithWorkers(num int) Option {
	return func(s *dashboardService) {
		if num > 0 {
			s.workers = num
		}
	}
}

// WithClock sets a custom clock function for timestamp generation.
func WithClock(clock func() time.Time) Option {
	return func(s *dashboardService) {
		if clock != nil {
			s.clock = clock
		}
	}
}

type dashboardService struct {
	mu           sync.RWMutex
	state        *DashboardState
	ghClient     github.Client
	devClient    DevcontainerClient
	pstoreClient pstore_client.PStoreClient
	workers      int
	clock        func() time.Time
}

// NewService creates a new dashboard aggregation Service instance.
func NewService(ghClient github.Client, devClient DevcontainerClient, pstoreClient pstore_client.PStoreClient, opts ...Option) Service {
	s := &dashboardService{
		ghClient:     ghClient,
		devClient:    devClient,
		pstoreClient: pstoreClient,
		workers:      5,
		clock:        time.Now,
		state: &DashboardState{
			PullRequests: []PRSummary{},
			Freshness:    FreshnessMetadata{},
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func mapContainerState(s manager_pb.State) ContainerState {
	switch s {
	case manager_pb.State_DCM_RECEIVED,
		manager_pb.State_DCM_CREATING,
		manager_pb.State_DCM_BRANCHING,
		manager_pb.State_DCM_HARNESS:
		return ContainerStateCreating
	case manager_pb.State_DCM_READY:
		return ContainerStateReady
	case manager_pb.State_DCM_FAILED,
		manager_pb.State_DCM_HARD_FAILED:
		return ContainerStateFailed
	default:
		return ContainerStateUnknown
	}
}

// GetDashboardState returns a thread-safe snapshot of the current dashboard state.
func (s *dashboardService) GetDashboardState(ctx context.Context) (*DashboardState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prsCopy := make([]PRSummary, len(s.state.PullRequests))
	copy(prsCopy, s.state.PullRequests)

	return &DashboardState{
		PullRequests: prsCopy,
		Freshness:    s.state.Freshness,
	}, nil
}

type prJob struct {
	owner    string
	repo     string
	repoFull string
	pr       *github.PullRequest
}

type prJobResult struct {
	summary PRSummary
	err     error
}

func (s *dashboardService) sync(ctx context.Context) error {
	now := s.clock()

	if s.pstoreClient == nil {
		err := fmt.Errorf("pstoreClient is not configured")
		s.recordFailure(now, err)
		return err
	}

	serverState, err := config.ReadServerState(ctx, s.pstoreClient)
	if err != nil {
		s.recordFailure(now, fmt.Errorf("failed to read server state: %w", err))
		return err
	}

	var dcmConfigs []*manager_pb.DevcontainerConfig
	if s.devClient != nil {
		listResp, err := s.devClient.List(ctx, &manager_pb.ListRequest{})
		if err != nil {
			s.recordFailure(now, fmt.Errorf("failed to list devcontainers: %w", err))
			return err
		}
		if listResp != nil {
			dcmConfigs = listResp.GetConfigs()
		}
	}

	dcmByID := make(map[string]*manager_pb.DevcontainerConfig)
	dcmByRepoPR := make(map[string]*manager_pb.DevcontainerConfig)
	for _, cfg := range dcmConfigs {
		if cfg == nil {
			continue
		}
		if cfg.GetId() != "" {
			dcmByID[cfg.GetId()] = cfg
		}
		if req := cfg.GetRequest(); req != nil && req.GetRepo() != "" {
			if id := req.GetIdentifier(); id != nil {
				prNum := id.GetPrNumber()
				if prNum > 0 {
					key := fmt.Sprintf("%s#%d", req.GetRepo(), prNum)
					dcmByRepoPR[key] = cfg
				}
			}
		}
	}

	var jobs []prJob
	for _, repoFullName := range serverState.GetEnrolledRepositories() {
		parts := strings.Split(repoFullName, "/")
		if len(parts) != 2 {
			log.Printf("invalid enrolled repo name: %s", repoFullName)
			continue
		}
		owner, repo := parts[0], parts[1]

		if s.ghClient != nil {
			prs, err := s.ghClient.ListOpenPullRequests(ctx, owner, repo)
			if err != nil {
				s.recordFailure(now, fmt.Errorf("failed to list pull requests for %s: %w", repoFullName, err))
				return err
			}
			for _, pr := range prs {
				if pr != nil {
					jobs = append(jobs, prJob{
						owner:    owner,
						repo:     repo,
						repoFull: repoFullName,
						pr:       pr,
					})
				}
			}
		}
	}

	workerCount := s.workers
	if workerCount <= 0 {
		workerCount = 1
	}
	if workerCount > len(jobs) {
		workerCount = len(jobs)
	}
	if workerCount == 0 {
		workerCount = 1
	}

	jobsChan := make(chan prJob, len(jobs))
	resultsChan := make(chan prJobResult, len(jobs))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsChan {
				summary, err := s.processPRJob(ctx, job, serverState, dcmByID, dcmByRepoPR)
				resultsChan <- prJobResult{
					summary: summary,
					err:     err,
				}
			}
		}()
	}

	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	wg.Wait()
	close(resultsChan)

	var summaries []PRSummary
	for res := range resultsChan {
		if res.err != nil {
			s.recordFailure(now, res.err)
			return res.err
		}
		summaries = append(summaries, res.summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].Repo == summaries[j].Repo {
			return summaries[i].PRNumber < summaries[j].PRNumber
		}
		return summaries[i].Repo < summaries[j].Repo
	})

	s.mu.Lock()
	s.state.PullRequests = summaries
	s.state.Freshness = FreshnessMetadata{
		LastSuccessfulSync: now,
		LastAttemptedSync:  now,
		IsStale:            false,
		ErrorMessage:       "",
	}
	s.mu.Unlock()

	return nil
}

func (s *dashboardService) processPRJob(ctx context.Context, job prJob, serverState interface{ GetPrContainers() []*pb.PRContainer }, dcmByID map[string]*manager_pb.DevcontainerConfig, dcmByRepoPR map[string]*manager_pb.DevcontainerConfig) (PRSummary, error) {
	detail, err := s.ghClient.GetPullRequestDetails(ctx, job.owner, job.repo, job.pr.Number)
	if err != nil {
		return PRSummary{}, fmt.Errorf("failed to get PR details for %s#%d: %w", job.repoFull, job.pr.Number, err)
	}

	checkStatus := CheckStatusUnknown
	if job.pr.Head.SHA != "" {
		st, err := s.ghClient.GetCommitCheckStatus(ctx, job.owner, job.repo, job.pr.Head.SHA)
		if err != nil {
			return PRSummary{}, fmt.Errorf("failed to get check status for %s#%d (sha: %s): %w", job.repoFull, job.pr.Number, job.pr.Head.SHA, err)
		}
		checkStatus = st
	}

	hasDevcontainer := false
	containerID := ""
	containerState := ContainerStateNone

	// Correlate with ServerState.PrContainers
	if serverState != nil {
		for _, prc := range serverState.GetPrContainers() {
			if prc != nil && prc.GetRepo() == job.repoFull && int(prc.GetPrNumber()) == job.pr.Number {
				hasDevcontainer = true
				containerID = prc.GetContainerId()
				break
			}
		}
	}

	if hasDevcontainer {
		if cfg, ok := dcmByID[containerID]; ok && cfg != nil {
			containerState = mapContainerState(cfg.GetState())
		} else {
			containerState = ContainerStateUnknown
		}
	} else {
		key := fmt.Sprintf("%s#%d", job.repoFull, job.pr.Number)
		if cfg, ok := dcmByRepoPR[key]; ok && cfg != nil {
			hasDevcontainer = true
			containerID = cfg.GetId()
			containerState = mapContainerState(cfg.GetState())
		}
	}

	commitCount := 0
	commentCount := 0
	if detail != nil {
		commitCount = detail.Commits
		commentCount = detail.Comments + detail.ReviewComments
	}

	return PRSummary{
		Repo:            job.repoFull,
		PRNumber:        int32(job.pr.Number),
		Title:           job.pr.Title,
		Author:          job.pr.User.Login,
		CommitCount:     commitCount,
		CommentCount:    commentCount,
		CheckStatus:     checkStatus,
		HasDevcontainer: hasDevcontainer,
		ContainerID:     containerID,
		ContainerState:  containerState,
	}, nil
}

func (s *dashboardService) recordFailure(now time.Time, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Freshness.LastAttemptedSync = now
	s.state.Freshness.IsStale = true
	if err != nil {
		s.state.Freshness.ErrorMessage = err.Error()
	}
}

// RunWorker runs the background aggregation sync loop at the specified interval.
func (s *dashboardService) RunWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial sync run
	if err := s.sync(ctx); err != nil {
		log.Printf("initial dashboard sync error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.sync(ctx); err != nil {
				log.Printf("dashboard sync error: %v", err)
			}
		}
	}
}
