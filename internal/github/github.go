package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
