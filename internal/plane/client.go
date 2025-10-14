package plane

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
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

func (c *Client) httpClient() *http.Client {
    if c.HTTP != nil { return c.HTTP }
    return http.DefaultClient
}

func (c *Client) join(paths ...string) (string, error) {
    u, err := url.Parse(c.BaseURL)
    if err != nil { return "", err }
    p := strings.TrimRight(u.Path, "/")
    for _, s := range paths {
        if s == "" { continue }
        if !strings.HasPrefix(s, "/") { s = "/" + s }
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
    if err != nil { return "", err }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+bearer)
    req.Header.Set("Content-Type", "application/json")
    hc := c.httpClient()
    hc.Timeout = 10 * time.Second
    resp, err := hc.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return "", fmt.Errorf("plane create issue status=%d", resp.StatusCode) }
    var out struct{ ID string `json:"id"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if out.ID == "" { return "", fmt.Errorf("empty issue id") }
    return out.ID, nil
}

// PatchIssue updates fields of an existing issue
func (c *Client) PatchIssue(ctx context.Context, bearer, workspaceSlug, projectID, issueID string, payload map[string]any) error {
    // PATCH /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/
    path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
    ep, err := c.join(path)
    if err != nil { return err }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+bearer)
    req.Header.Set("Content-Type", "application/json")
    hc := c.httpClient()
    hc.Timeout = 10 * time.Second
    resp, err := hc.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("plane patch issue status=%d", resp.StatusCode) }
    return nil
}

// AddComment posts a comment to an issue
func (c *Client) AddComment(ctx context.Context, bearer, workspaceSlug, projectID, issueID, commentHTML string) error {
    // POST /api/v1/workspaces/{workspace-slug}/projects/{project_id}/issues/{issue_id}/comments/
    path := fmt.Sprintf("/api/v1/workspaces/%s/projects/%s/issues/%s/comments/", url.PathEscape(workspaceSlug), url.PathEscape(projectID), url.PathEscape(issueID))
    ep, err := c.join(path)
    if err != nil { return err }
    payload := map[string]any{"comment_html": commentHTML}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+bearer)
    req.Header.Set("Content-Type", "application/json")
    hc := c.httpClient()
    hc.Timeout = 10 * time.Second
    resp, err := hc.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("plane add comment status=%d", resp.StatusCode) }
    return nil
}

