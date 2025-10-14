package cnb

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// Client is CNB API client. Endpoints are configurable via path templates.
// If any required path template is missing, methods return ErrNotConfigured.
type Client struct {
    BaseURL string
    Token   string

    IssueCreatePath  string // /api/v1/repos/{repo}/issues
    IssueUpdatePath  string // /api/v1/repos/{repo}/issues/{issue_iid}
    IssueCommentPath string // /api/v1/repos/{repo}/issues/{issue_iid}/comments
    IssueClosePath   string // /api/v1/repos/{repo}/issues/{issue_iid}

    HTTP *http.Client
}

var ErrNotConfigured = errors.New("CNB API endpoints not configured; 待确认")

func (c *Client) httpClient() *http.Client {
    if c.HTTP != nil { return c.HTTP }
    return &http.Client{Timeout: 10 * time.Second}
}

func (c *Client) join(path string) (string, error) {
    u, err := url.Parse(c.BaseURL)
    if err != nil { return "", err }
    if !strings.HasPrefix(path, "/") { path = "/" + path }
    u.Path = strings.TrimRight(u.Path, "/") + path
    return u.String(), nil
}

func expand(tpl string, repo, issueIID string) (string, error) {
    if tpl == "" { return "", ErrNotConfigured }
    s := strings.ReplaceAll(tpl, "{repo}", url.PathEscape(repo))
    s = strings.ReplaceAll(s, "{issue_iid}", url.PathEscape(issueIID))
    if strings.Contains(s, "{") { // unresolved placeholder
        return "", ErrNotConfigured
    }
    return s, nil
}

// CreateIssue creates an issue in CNB and returns its IID (string).
// Payload keys are tentative and marked 待确认 per spec.
func (c *Client) CreateIssue(ctx context.Context, repo string, title, descriptionHTML string) (string, error) {
    if c.BaseURL == "" || c.Token == "" { return "", ErrNotConfigured }
    p, err := expand(c.IssueCreatePath, repo, "")
    if err != nil { return "", err }
    ep, err := c.join(p)
    if err != nil { return "", err }
    payload := map[string]any{"title": title}
    if strings.TrimSpace(descriptionHTML) != "" { payload["description_html"] = descriptionHTML }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return "", fmt.Errorf("cnb create issue status=%d", resp.StatusCode) }
    var out struct{ IID string `json:"issue_iid"`; ID string `json:"id"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if out.IID != "" { return out.IID, nil }
    if out.ID != "" { return out.ID, nil }
    return "", errors.New("empty issue iid")
}

// UpdateIssue updates fields of an issue in CNB.
func (c *Client) UpdateIssue(ctx context.Context, repo, issueIID string, fields map[string]any) error {
    if c.BaseURL == "" || c.Token == "" { return ErrNotConfigured }
    p, err := expand(c.IssueUpdatePath, repo, issueIID)
    if err != nil { return err }
    ep, err := c.join(p)
    if err != nil { return err }
    b, _ := json.Marshal(fields)
    req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("cnb update issue status=%d", resp.StatusCode) }
    return nil
}

// AddComment posts a comment to a CNB issue.
func (c *Client) AddComment(ctx context.Context, repo, issueIID, commentHTML string) error {
    if c.BaseURL == "" || c.Token == "" { return ErrNotConfigured }
    p, err := expand(c.IssueCommentPath, repo, issueIID)
    if err != nil { return err }
    ep, err := c.join(p)
    if err != nil { return err }
    payload := map[string]any{"comment_html": commentHTML}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("cnb add comment status=%d", resp.StatusCode) }
    return nil
}

// CloseIssue transitions a CNB issue to closed state (payload tentative).
func (c *Client) CloseIssue(ctx context.Context, repo, issueIID string) error {
    if c.BaseURL == "" || c.Token == "" { return ErrNotConfigured }
    p, err := expand(c.IssueClosePath, repo, issueIID)
    if err != nil { return err }
    ep, err := c.join(p)
    if err != nil { return err }
    payload := map[string]any{"state": "closed"}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("cnb close issue status=%d", resp.StatusCode) }
    return nil
}
