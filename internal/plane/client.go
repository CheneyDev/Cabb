package plane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL string // e.g., https://api.plane.so
	// HTTP client; if nil, http.DefaultClient used
	HTTP *http.Client
}

type Workspace struct {
	Name string `json:"name"`
	// Some Plane deployments expose display name via "title"
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type Project struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Slug       string `json:"slug"`
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

func (c *Client) join(paths ...string) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	p := strings.TrimRight(u.Path, "/")
	for _, s := range paths {
		if s == "" {
			continue
		}
		if !strings.HasPrefix(s, "/") {
			s = "/" + s
		}
		p += s
	}
	u.Path = p
	return u.String(), nil
}

// CreateIssue creates a work item in Plane
func (c *Client) CreateIssue(ctx context.Context, bearer, workspaceSlug, projectID string, payload map[string]any) (issueID string, err error) {
	// POST /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/", url.PathEscape(workspaceSlug), url.PathEscape(projectID))
	ep, err := c.join(path)
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-Key", bearer)
	req.Header.Set("Content-Type", "application/json")
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("plane create issue status=%d, body=%s", resp.StatusCode, string(body))
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.ID == "" {
		return "", fmt.Errorf("empty issue id")
	}
	return out.ID, nil
}

// PatchIssue updates fields of an existing issue
func (c *Client) PatchIssue(ctx context.Context, bearer, workspaceSlug, projectID, issueID string, payload map[string]any) error {
	// PATCH /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
	ep, err := c.join(path)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", bearer)
	req.Header.Set("Content-Type", "application/json")
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("plane patch issue status=%d, body=%s", resp.StatusCode, string(body))
	}
	return nil
}

// AddComment posts a comment to an issue
func (c *Client) AddComment(ctx context.Context, bearer, workspaceSlug, projectID, issueID, commentHTML string) error {
	// POST /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/comments/
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/comments/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
	ep, err := c.join(path)
	if err != nil {
		return err
	}
	payload := map[string]any{"comment_html": commentHTML}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", bearer)
	req.Header.Set("Content-Type", "application/json")
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("plane add comment status=%d, body=%s", resp.StatusCode, string(body))
	}
	return nil
}

// GetIssueBySequence fetches an issue by workspace-level sequence id (e.g., KEY-123).
// GET /api/v1/workspaces/{workspace-slug}/issues/{sequence_id}/
// Returns the canonical issue id and project id when available.
func (c *Client) GetIssueBySequence(ctx context.Context, bearer, workspaceSlug, sequenceID string) (issueID, projectID string, err error) {
	path := fmt.Sprintf("/api/v1/workspaces/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(sequenceID))
	ep, err := c.join(path)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("plane get issue by sequence status=%d, body=%s", resp.StatusCode, string(body))
	}
	// Try decode with flexible fields
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", "", err
	}
	getStr := func(key string) string {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}
	issueID = getStr("id")
	projectID = getStr("project")
	if projectID == "" {
		projectID = getStr("project_id")
	}
	if issueID == "" {
		return "", "", fmt.Errorf("empty issue id")
	}
	return issueID, projectID, nil
}

// GetIssueName fetches issue details and returns the issue name (title).
// GET /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/
func (c *Client) GetIssueName(ctx context.Context, bearer, workspaceSlug, projectID, issueID string) (string, error) {
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
	ep, err := c.join(path)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("plane get issue status=%d, body=%s", resp.StatusCode, string(body))
	}
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}
	if v, ok := m["name"].(string); ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("empty name")
}

// IssueDetail represents the complete issue details from Plane API
type IssueDetail struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	DescriptionHTML     string   `json:"description_html"`
	DescriptionStripped string   `json:"description_stripped"`
	Priority            string   `json:"priority"`
	SequenceID          int      `json:"sequence_id"`
	Project             string   `json:"project"`
	Workspace           string   `json:"workspace"`
	Parent              *string  `json:"parent"`
	State               string   `json:"state"`
	Assignees           []string `json:"assignees"`
	Labels              []string `json:"labels"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
}

// GetIssueDetail fetches complete issue details from Plane API.
// GET /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/
func (c *Client) GetIssueDetail(ctx context.Context, bearer, workspaceSlug, projectID, issueID string) (*IssueDetail, error) {
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
	ep, err := c.join(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("plane get issue detail status=%d, body=%s", resp.StatusCode, string(body))
	}

	var detail IssueDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	if detail.ID == "" {
		return nil, fmt.Errorf("empty issue id")
	}

	return &detail, nil
}

// FindIssueInWorkspace searches for an issue by ID across all projects in a workspace.
// This is useful when we have an issue ID but don't know which project it belongs to.
func (c *Client) FindIssueInWorkspace(ctx context.Context, bearer, workspaceSlug, issueID string) (*IssueDetail, error) {
	// First, get all projects in the workspace
	projects, err := c.ListProjects(ctx, bearer, workspaceSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	// Search for the issue in each project
	for _, project := range projects {
		projectID, ok := project["id"].(string)
		if !ok || projectID == "" {
			continue
		}

		detail, err := c.GetIssueDetail(ctx, bearer, workspaceSlug, projectID, issueID)
		if err != nil {
			// If issue not found in this project, continue to next one
			if strings.Contains(err.Error(), "status=404") {
				continue
			}
			// For other errors, log but continue searching
			continue
		}

		// Found the issue
		return detail, nil
	}

	return nil, fmt.Errorf("issue %s not found in any project of workspace %s", issueID, workspaceSlug)
}

// GetIssueLabels fetches current label IDs of an issue
// GET /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/
func (c *Client) GetIssueLabels(ctx context.Context, bearer, workspaceSlug, projectID, issueID string) ([]string, error) {
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
	ep, err := c.join(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("plane get issue status=%d", resp.StatusCode)
	}
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	var labelIDs []string
	if labels, ok := m["labels"].([]interface{}); ok {
		for _, label := range labels {
			if labelID, ok := label.(string); ok && labelID != "" {
				labelIDs = append(labelIDs, labelID)
			}
		}
	}
	return labelIDs, nil
}

// GetWorkspace fetches metadata of a workspace by slug.
func (c *Client) GetWorkspace(ctx context.Context, bearer, workspaceSlug string) (*Workspace, error) {
	if strings.TrimSpace(workspaceSlug) == "" {
		return nil, fmt.Errorf("workspace slug is empty")
	}
	path := fmt.Sprintf("/api/v1/workspaces/%s/", url.PathEscape(workspaceSlug))
	ep, err := c.join(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("plane get workspace status=%d, body=%s", resp.StatusCode, string(body))
	}
	var out Workspace
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	// Ensure slug fallback is populated from response path when missing
	if out.Slug == "" {
		out.Slug = workspaceSlug
	}
	return &out, nil
}

// GetProject fetches project metadata for a workspace.
func (c *Client) GetProject(ctx context.Context, bearer, workspaceSlug, projectID string) (*Project, error) {
	if strings.TrimSpace(workspaceSlug) == "" {
		return nil, fmt.Errorf("workspace slug is empty")
	}
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("project id is empty")
	}
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID))
	ep, err := c.join(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("plane get project status=%d, body=%s", resp.StatusCode, string(body))
	}
	var out Project
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListProjects fetches all projects within a workspace.
// GET /api/v1/workspaces/{workspace_slug}/projects/
func (c *Client) ListProjects(ctx context.Context, bearer, workspaceSlug string) ([]map[string]any, error) {
	if strings.TrimSpace(workspaceSlug) == "" {
		return nil, fmt.Errorf("workspace slug is empty")
	}
	path := fmt.Sprintf("/api/v1/workspaces/%s/projects/", url.PathEscape(workspaceSlug))
	ep, err := c.join(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", bearer)
	hc := c.httpClient()
	hc.Timeout = 10 * time.Second
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("plane list projects status=%d, body=%s", resp.StatusCode, string(body))
	}
	// Plane API 返回分页响应，解析 results 字段
	var response struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Results, nil
}
