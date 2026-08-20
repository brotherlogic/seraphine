package dashboard

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	manager_pb "github.com/brotherlogic/devcontainer-manager/proto"
	pstore_client "github.com/brotherlogic/pstore/client"
	"github.com/brotherlogic/seraphine/internal/config"
	"github.com/brotherlogic/seraphine/internal/github"
	pb "github.com/brotherlogic/seraphine/proto"
	"google.golang.org/grpc"
)

type mockGHClient struct {
	listOpenPRsFunc func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error)
	getPRDetailFunc func(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error)
	getCheckFunc    func(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error)
}

func (m *mockGHClient) ListRepositoryInvitations(ctx context.Context) ([]*github.RepositoryInvitation, error) {
	return nil, nil
}
func (m *mockGHClient) AcceptRepositoryInvitation(ctx context.Context, invitationID int64) error {
	return nil
}
func (m *mockGHClient) CreateRuleset(ctx context.Context, owner, repo string, ruleset *github.RulesetRequest) error {
	return nil
}
func (m *mockGHClient) CreateIssue(ctx context.Context, owner, repo string, title, body string, labels []string) (*github.IssueResponse, error) {
	return nil, nil
}
func (m *mockGHClient) IsCollaborator(ctx context.Context, owner, repo, user string) (bool, error) {
	return true, nil
}
func (m *mockGHClient) ListOpenPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	if m.listOpenPRsFunc != nil {
		return m.listOpenPRsFunc(ctx, owner, repo)
	}
	return nil, nil
}
func (m *mockGHClient) GetPullRequestDetails(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
	if m.getPRDetailFunc != nil {
		return m.getPRDetailFunc(ctx, owner, repo, number)
	}
	return nil, nil
}
func (m *mockGHClient) GetCommitCheckStatus(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
	if m.getCheckFunc != nil {
		return m.getCheckFunc(ctx, owner, repo, ref)
	}
	return github.CheckStatusUnknown, nil
}

type mockDevClient struct {
	listFunc func(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error)
}

func (m *mockDevClient) Up(ctx context.Context, in *manager_pb.UpRequest, opts ...grpc.CallOption) (*manager_pb.UpResponse, error) {
	return nil, nil
}
func (m *mockDevClient) Down(ctx context.Context, in *manager_pb.DownRequest, opts ...grpc.CallOption) (*manager_pb.DownResponse, error) {
	return nil, nil
}
func (m *mockDevClient) List(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, in, opts...)
	}
	return &manager_pb.ListResponse{}, nil
}

func TestDashboardService_ZeroReposAndZeroPRs(t *testing.T) {
	fixedTime := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{},
	})

	mockGH := &mockGHClient{}
	mockDev := &mockDevClient{}

	svc := NewService(mockGH, mockDev, pClient, WithClock(func() time.Time { return fixedTime }))
	if s, ok := svc.(*dashboardService); ok {
		if err := s.sync(context.Background()); err != nil {
			t.Fatalf("sync failed: %v", err)
		}
	}

	state, err := svc.GetDashboardState(context.Background())
	if err != nil {
		t.Fatalf("GetDashboardState returned error: %v", err)
	}
	if len(state.PullRequests) != 0 {
		t.Errorf("expected 0 PRs, got %d", len(state.PullRequests))
	}
	if state.Freshness.IsStale {
		t.Errorf("expected IsStale to be false")
	}
	if !state.Freshness.LastSuccessfulSync.Equal(fixedTime) {
		t.Errorf("expected LastSuccessfulSync %v, got %v", fixedTime, state.Freshness.LastSuccessfulSync)
	}
}

func TestDashboardService_FullAggregationHappyPath(t *testing.T) {
	fixedTime := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/seraphine", "brotherlogic/pstore"},
		PrContainers: []*pb.PRContainer{
			{Repo: "brotherlogic/seraphine", PrNumber: 101, ContainerId: "c-101"},
		},
	})

	mockGH := &mockGHClient{
		listOpenPRsFunc: func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
			if repo == "seraphine" {
				return []*github.PullRequest{
					{
						Number: 101,
						Title:  "Fix webhook handler",
						User:   github.User{Login: "dev1"},
						Head:   github.PullRequestBranch{SHA: "sha101"},
					},
					{
						Number: 102,
						Title:  "Add metrics",
						User:   github.User{Login: "dev2"},
						Head:   github.PullRequestBranch{SHA: "sha102"},
					},
				}, nil
			}
			if repo == "pstore" {
				return []*github.PullRequest{
					{
						Number: 201,
						Title:  "Upgrade proto",
						User:   github.User{Login: "dev3"},
						Head:   github.PullRequestBranch{SHA: "sha201"},
					},
				}, nil
			}
			return nil, nil
		},
		getPRDetailFunc: func(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
			if repo == "seraphine" && number == 101 {
				return &github.PullRequestDetail{
					Commits:        3,
					Comments:       2,
					ReviewComments: 4,
				}, nil
			}
			if repo == "seraphine" && number == 102 {
				return &github.PullRequestDetail{
					Commits:        1,
					Comments:       0,
					ReviewComments: 0,
				}, nil
			}
			if repo == "pstore" && number == 201 {
				return &github.PullRequestDetail{
					Commits:        5,
					Comments:       1,
					ReviewComments: 1,
				}, nil
			}
			return &github.PullRequestDetail{}, nil
		},
		getCheckFunc: func(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
			switch ref {
			case "sha101":
				return github.CheckStatusSuccess, nil
			case "sha102":
				return github.CheckStatusPending, nil
			case "sha201":
				return github.CheckStatusFailure, nil
			}
			return github.CheckStatusUnknown, nil
		},
	}

	mockDev := &mockDevClient{
		listFunc: func(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error) {
			return &manager_pb.ListResponse{
				Configs: []*manager_pb.DevcontainerConfig{
					{
						Id:    "c-101",
						State: manager_pb.State_DCM_READY,
					},
					{
						Id: "c-201",
						Request: &manager_pb.UpRequest{
							Repo: "brotherlogic/pstore",
							Identifier: &manager_pb.Identifier{
								Id: &manager_pb.Identifier_PrNumber{
									PrNumber: 201,
								},
							},
						},
						State: manager_pb.State_DCM_CREATING,
					},
				},
			}, nil
		},
	}

	svc := NewService(mockGH, mockDev, pClient, WithWorkers(4), WithClock(func() time.Time { return fixedTime }))
	if s, ok := svc.(*dashboardService); ok {
		if err := s.sync(context.Background()); err != nil {
			t.Fatalf("sync failed: %v", err)
		}
	}

	state, err := svc.GetDashboardState(context.Background())
	if err != nil {
		t.Fatalf("GetDashboardState error: %v", err)
	}

	if len(state.PullRequests) != 3 {
		t.Fatalf("expected 3 PRs, got %d", len(state.PullRequests))
	}

	// Verify PR 101
	pr101 := state.PullRequests[1] // seraphine #101
	if pr101.Repo != "brotherlogic/seraphine" || pr101.PRNumber != 101 {
		for _, pr := range state.PullRequests {
			if pr.Repo == "brotherlogic/seraphine" && pr.PRNumber == 101 {
				pr101 = pr
				break
			}
		}
	}
	if pr101.CommitCount != 3 || pr101.CommentCount != 6 {
		t.Errorf("PR 101 counts mismatch: commits=%d comments=%d", pr101.CommitCount, pr101.CommentCount)
	}
	if pr101.CheckStatus != CheckStatusSuccess {
		t.Errorf("PR 101 check status expected SUCCESS, got %s", pr101.CheckStatus)
	}
	if !pr101.HasDevcontainer || pr101.ContainerID != "c-101" || pr101.ContainerState != ContainerStateReady {
		t.Errorf("PR 101 container mismatch: has=%v id=%s state=%s", pr101.HasDevcontainer, pr101.ContainerID, pr101.ContainerState)
	}

	// Verify PR 102 (no container)
	var pr102 PRSummary
	for _, pr := range state.PullRequests {
		if pr.Repo == "brotherlogic/seraphine" && pr.PRNumber == 102 {
			pr102 = pr
			break
		}
	}
	if pr102.HasDevcontainer || pr102.ContainerState != ContainerStateNone {
		t.Errorf("PR 102 container mismatch: has=%v state=%s", pr102.HasDevcontainer, pr102.ContainerState)
	}
	if pr102.CheckStatus != CheckStatusPending {
		t.Errorf("PR 102 check status expected PENDING, got %s", pr102.CheckStatus)
	}

	// Verify PR 201 (creating container)
	var pr201 PRSummary
	for _, pr := range state.PullRequests {
		if pr.Repo == "brotherlogic/pstore" && pr.PRNumber == 201 {
			pr201 = pr
			break
		}
	}
	if !pr201.HasDevcontainer || pr201.ContainerID != "c-201" || pr201.ContainerState != ContainerStateCreating {
		t.Errorf("PR 201 container mismatch: has=%v id=%s state=%s", pr201.HasDevcontainer, pr201.ContainerID, pr201.ContainerState)
	}
	if pr201.CheckStatus != CheckStatusFailure {
		t.Errorf("PR 201 check status expected FAILURE, got %s", pr201.CheckStatus)
	}
}

func TestDashboardService_UpstreamFailureAndStaleCache(t *testing.T) {
	time1 := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2026, 8, 20, 12, 5, 0, 0, time.UTC)
	currTime := time1

	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/seraphine"},
	})

	failGH := false
	mockGH := &mockGHClient{
		listOpenPRsFunc: func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
			if failGH {
				return nil, errors.New("github rate limit exceeded")
			}
			return []*github.PullRequest{
				{
					Number: 101,
					Title:  "Fix webhook handler",
					User:   github.User{Login: "dev1"},
					Head:   github.PullRequestBranch{SHA: "sha101"},
				},
			}, nil
		},
		getPRDetailFunc: func(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
			return &github.PullRequestDetail{Commits: 1, Comments: 1}, nil
		},
		getCheckFunc: func(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
			return github.CheckStatusSuccess, nil
		},
	}

	mockDev := &mockDevClient{}

	svc := NewService(mockGH, mockDev, pClient, WithClock(func() time.Time { return currTime }))
	s, ok := svc.(*dashboardService)
	if !ok {
		t.Fatalf("expected *dashboardService")
	}

	// 1st sync: Success
	if err := s.sync(context.Background()); err != nil {
		t.Fatalf("first sync failed: %v", err)
	}

	st1, _ := svc.GetDashboardState(context.Background())
	if len(st1.PullRequests) != 1 || st1.Freshness.IsStale {
		t.Fatalf("initial state invalid: %+v", st1)
	}

	// 2nd sync: Upstream failure
	currTime = time2
	failGH = true
	err := s.sync(context.Background())
	if err == nil {
		t.Fatalf("expected sync error on github failure")
	}

	st2, _ := svc.GetDashboardState(context.Background())
	if len(st2.PullRequests) != 1 {
		t.Errorf("expected cached PRs to be retained, got %d PRs", len(st2.PullRequests))
	}
	if !st2.Freshness.IsStale {
		t.Errorf("expected IsStale to be true")
	}
	if !st2.Freshness.LastSuccessfulSync.Equal(time1) {
		t.Errorf("expected LastSuccessfulSync %v, got %v", time1, st2.Freshness.LastSuccessfulSync)
	}
	if !st2.Freshness.LastAttemptedSync.Equal(time2) {
		t.Errorf("expected LastAttemptedSync %v, got %v", time2, st2.Freshness.LastAttemptedSync)
	}
	if st2.Freshness.ErrorMessage == "" {
		t.Errorf("expected non-empty ErrorMessage")
	}
}

func TestDashboardService_DevcontainerListFailure(t *testing.T) {
	time1 := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/seraphine"},
	})

	mockGH := &mockGHClient{
		listOpenPRsFunc: func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
			return []*github.PullRequest{{Number: 1}}, nil
		},
		getPRDetailFunc: func(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
			return &github.PullRequestDetail{}, nil
		},
		getCheckFunc: func(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
			return github.CheckStatusSuccess, nil
		},
	}

	mockDev := &mockDevClient{
		listFunc: func(ctx context.Context, in *manager_pb.ListRequest, opts ...grpc.CallOption) (*manager_pb.ListResponse, error) {
			return nil, errors.New("dcm service unavailable")
		},
	}

	svc := NewService(mockGH, mockDev, pClient, WithClock(func() time.Time { return time1 }))
	s := svc.(*dashboardService)

	err := s.sync(context.Background())
	if err == nil {
		t.Fatalf("expected error on DCM list failure")
	}

	st, _ := svc.GetDashboardState(context.Background())
	if !st.Freshness.IsStale {
		t.Errorf("expected IsStale = true on DCM failure")
	}
}

func TestDashboardService_ConcurrencyAndRace(t *testing.T) {
	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/seraphine"},
	})

	mockGH := &mockGHClient{
		listOpenPRsFunc: func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
			return []*github.PullRequest{
				{Number: 1, Title: "PR 1", User: github.User{Login: "user1"}, Head: github.PullRequestBranch{SHA: "s1"}},
				{Number: 2, Title: "PR 2", User: github.User{Login: "user2"}, Head: github.PullRequestBranch{SHA: "s2"}},
			}, nil
		},
		getPRDetailFunc: func(ctx context.Context, owner, repo string, number int) (*github.PullRequestDetail, error) {
			return &github.PullRequestDetail{Commits: number}, nil
		},
		getCheckFunc: func(ctx context.Context, owner, repo, ref string) (github.CheckStatus, error) {
			return github.CheckStatusSuccess, nil
		},
	}
	mockDev := &mockDevClient{}

	svc := NewService(mockGH, mockDev, pClient, WithWorkers(3))
	s := svc.(*dashboardService)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 5 readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_, _ = svc.GetDashboardState(context.Background())
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	// 2 sync writers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_ = s.sync(context.Background())
					time.Sleep(20 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestDashboardService_RunWorker(t *testing.T) {
	pClient := pstore_client.GetTestClient()
	_ = config.WriteServerState(context.Background(), pClient, &pb.ServerState{
		EnrolledRepositories: []string{"brotherlogic/seraphine"},
	})

	syncCount := 0
	var mu sync.Mutex

	mockGH := &mockGHClient{
		listOpenPRsFunc: func(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
			mu.Lock()
			syncCount++
			mu.Unlock()
			return nil, nil
		},
	}
	mockDev := &mockDevClient{}

	svc := NewService(mockGH, mockDev, pClient)
	ctx, cancel := context.WithCancel(context.Background())

	go svc.RunWorker(ctx, 50*time.Millisecond)

	time.Sleep(150 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	count := syncCount
	mu.Unlock()

	if count < 2 {
		t.Errorf("expected RunWorker to trigger at least 2 syncs, got %d", count)
	}
}

func TestDashboardService_ContainerStateMapping(t *testing.T) {
	tests := []struct {
		dcmState manager_pb.State
		expected ContainerState
	}{
		{manager_pb.State_DCM_RECEIVED, ContainerStateCreating},
		{manager_pb.State_DCM_CREATING, ContainerStateCreating},
		{manager_pb.State_DCM_BRANCHING, ContainerStateCreating},
		{manager_pb.State_DCM_HARNESS, ContainerStateCreating},
		{manager_pb.State_DCM_READY, ContainerStateReady},
		{manager_pb.State_DCM_FAILED, ContainerStateFailed},
		{manager_pb.State_DCM_HARD_FAILED, ContainerStateFailed},
		{manager_pb.State_UNKNOWN_STATE, ContainerStateUnknown},
	}

	for _, tc := range tests {
		got := mapContainerState(tc.dcmState)
		if got != tc.expected {
			t.Errorf("mapContainerState(%v) = %v, expected %v", tc.dcmState, got, tc.expected)
		}
	}
}
