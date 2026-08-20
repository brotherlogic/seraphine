package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RepositoryInvitation struct {
	ID         int64      `json:"id"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
}

type User struct {
	Login string `json:"login"`
}

type PullRequestBranch struct {
	SHA string `json:"sha"`
	Ref string `json:"ref"`
}

type PullRequest struct {
	ID      int64             `json:"id"`
	Number  int               `json:"number"`
	Title   string            `json:"title"`
	State   string            `json:"state"`
	User    User              `json:"user"`
	Head    PullRequestBranch `json:"head"`
	Base    PullRequestBranch `json:"base"`
	HTMLURL string            `json:"html_url"`
}

type PullRequestDetail struct {
	PullRequest
	Commits        int `json:"commits"`
	Comments       int `json:"comments"`
	ReviewComments int `json:"review_comments"`
}

type CheckStatus string

const (
	CheckStatusSuccess CheckStatus = "SUCCESS"
	CheckStatusPending CheckStatus = "PENDING"
	CheckStatusFailure CheckStatus = "FAILURE"
	CheckStatusUnknown CheckStatus = "UNKNOWN"
)

type CheckRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	HeadSHA    string `json:"head_sha"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type CheckRunsResponse struct {
	TotalCount int        `json:"total_count"`
	CheckRuns  []CheckRun `json:"check_runs"`
}

type RulesetRequest struct {
	Name        string     `json:"name"`
	Target      string     `json:"target"`
	Enforcement string     `json:"enforcement"`
	Conditions  Conditions `json:"conditions,omitempty"`
	Rules       []Rule     `json:"rules"`
}

type Conditions struct {
	RefName RefName `json:"ref_name"`
}

type RefName struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Rule struct {
	Type       string          `json:"type"`
	Parameters *RuleParameters `json:"parameters,omitempty"`
}

type RuleParameters struct {
	RequiredApprovingReviewCount int  `json:"required_approving_review_count"`
	DismissStaleReviewsOnPush     bool `json:"dismiss_stale_reviews_on_push"`
	RequireCodeOwnerReview        bool `json:"require_code_owner_review"`
	RequireLastPushApproval       bool `json:"require_last_push_approval"`
	RequiredReviewThreadResolution bool `json:"required_review_thread_resolution"`
}

type IssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels,omitempty"`
}

type IssueResponse struct {
	Number int `json:"number"`
}

type Client interface {
	ListRepositoryInvitations(ctx context.Context) ([]*RepositoryInvitation, error)
	AcceptRepositoryInvitation(ctx context.Context, invitationID int64) error
	CreateRuleset(ctx context.Context, owner, repo string, ruleset *RulesetRequest) error
	CreateIssue(ctx context.Context, owner, repo string, title, body string, labels []string) (*IssueResponse, error)
	IsCollaborator(ctx context.Context, owner, repo, user string) (bool, error)
	ListOpenPullRequests(ctx context.Context, owner, repo string) ([]*PullRequest, error)
	GetPullRequestDetails(ctx context.Context, owner, repo string, number int) (*PullRequestDetail, error)
	GetCommitCheckStatus(ctx context.Context, owner, repo, ref string) (CheckStatus, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type githubClient struct {
	httpClient HTTPClient
	token      string
	baseURL    string
}

func NewClient(token string, httpClient HTTPClient) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &githubClient{
		httpClient: httpClient,
		token:      token,
		baseURL:    "https://api.github.com",
	}
}

func (c *githubClient) doRequest(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

func (c *githubClient) ListRepositoryInvitations(ctx context.Context) ([]*RepositoryInvitation, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/user/repository_invitations", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list invitations, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var invitations []*RepositoryInvitation
	if err := json.NewDecoder(resp.Body).Decode(&invitations); err != nil {
		return nil, fmt.Errorf("failed to decode invitations: %w", err)
	}

	return invitations, nil
}

func (c *githubClient) AcceptRepositoryInvitation(ctx context.Context, invitationID int64) error {
	path := fmt.Sprintf("/user/repository_invitations/%d", invitationID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to accept invitation %d, status: %d, body: %s", invitationID, resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (c *githubClient) CreateRuleset(ctx context.Context, owner, repo string, ruleset *RulesetRequest) error {
	path := fmt.Sprintf("/repos/%s/%s/rulesets", owner, repo)
	resp, err := c.doRequest(ctx, http.MethodPost, path, ruleset)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create ruleset, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (c *githubClient) CreateIssue(ctx context.Context, owner, repo string, title, body string, labels []string) (*IssueResponse, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)
	reqBody := &IssueRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create issue, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var issueResp IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		return nil, fmt.Errorf("failed to decode issue response: %w", err)
	}

	return &issueResp, nil
}

func (c *githubClient) IsCollaborator(ctx context.Context, owner, repo, user string) (bool, error) {
	path := fmt.Sprintf("/repos/%s/%s/collaborators/%s", owner, repo, user)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("failed to check collaborator status for user %s in %s/%s, status: %d, body: %s", user, owner, repo, resp.StatusCode, string(bodyBytes))
}

func (c *githubClient) ListOpenPullRequests(ctx context.Context, owner, repo string) ([]*PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls?state=open", owner, repo)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list open pull requests for %s/%s, status: %d, body: %s", owner, repo, resp.StatusCode, string(bodyBytes))
	}

	var prs []*PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("failed to decode pull requests: %w", err)
	}

	return prs, nil
}

func (c *githubClient) GetPullRequestDetails(ctx context.Context, owner, repo string, number int) (*PullRequestDetail, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get pull request details for %s/%s#%d, status: %d, body: %s", owner, repo, number, resp.StatusCode, string(bodyBytes))
	}

	var prDetail PullRequestDetail
	if err := json.NewDecoder(resp.Body).Decode(&prDetail); err != nil {
		return nil, fmt.Errorf("failed to decode pull request details: %w", err)
	}

	return &prDetail, nil
}

func evaluateCheckStatus(runs []CheckRun) CheckStatus {
	if len(runs) == 0 {
		return CheckStatusUnknown
	}

	hasPending := false
	for _, run := range runs {
		status := strings.ToLower(run.Status)
		conclusion := strings.ToLower(run.Conclusion)

		if conclusion == "failure" || conclusion == "timed_out" || conclusion == "cancelled" || conclusion == "action_required" || status == "failure" || status == "timed_out" {
			return CheckStatusFailure
		}

		if status == "in_progress" || status == "queued" || status == "waiting" || status == "pending" || (status != "completed" && conclusion == "") {
			hasPending = true
		}
	}

	if hasPending {
		return CheckStatusPending
	}

	return CheckStatusSuccess
}

func (c *githubClient) GetCommitCheckStatus(ctx context.Context, owner, repo, ref string) (CheckStatus, error) {
	path := fmt.Sprintf("/repos/%s/%s/commits/%s/check-runs", owner, repo, ref)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return CheckStatusUnknown, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return CheckStatusUnknown, fmt.Errorf("failed to get commit check runs for %s/%s ref %s, status: %d, body: %s", owner, repo, ref, resp.StatusCode, string(bodyBytes))
	}

	var checkRunsResp CheckRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkRunsResp); err != nil {
		return CheckStatusUnknown, fmt.Errorf("failed to decode check runs response: %w", err)
	}

	return evaluateCheckStatus(checkRunsResp.CheckRuns), nil
}

