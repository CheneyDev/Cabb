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

// Client is CNB API client.
// Endpoints follow CNB swagger (docs/cnb-docs/swagger.json):
//   - POST   /{repo}/-/issues
//   - PATCH  /{repo}/-/issues/{number}
//   - POST   /{repo}/-/issues/{number}/comments
type Client struct {
    BaseURL string
    Token   string
    HTTP    *http.Client

    // Optional overrides (advanced): allow customizing path templates if CNB has a non-standard gateway.
    IssueCreatePath  string // default: "/{repo}/-/issues"
    IssueUpdatePath  string // default: "/{repo}/-/issues/{number}"
    IssueCommentPath string // default: "/{repo}/-/issues/{number}/comments"
}

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

func (c *Client) defaultCreatePath() string {
    if c.IssueCreatePath != "" { return c.IssueCreatePath }
    return "/{repo}/-/issues"
}
func (c *Client) defaultUpdatePath() string {
    if c.IssueUpdatePath != "" { return c.IssueUpdatePath }
    return "/{repo}/-/issues/{number}"
}
func (c *Client) defaultCommentPath() string {
    if c.IssueCommentPath != "" { return c.IssueCommentPath }
    return "/{repo}/-/issues/{number}/comments"
}

func expand(tpl string, repo, number string) (string, error) {
    if tpl == "" { return "", errors.New("empty path template") }
    // Encode repo path by segments so that slashes remain separators
    // e.g., "org/repo" => "org/repo" but each segment escaped
    repoEsc := repo
    if repo != "" {
        parts := strings.Split(repo, "/")
        for i, p := range parts { parts[i] = url.PathEscape(p) }
        repoEsc = strings.Join(parts, "/")
    }
    s := strings.ReplaceAll(tpl, "{repo}", repoEsc)
    // Backward-compat: support {issue_iid} placeholder name
    s = strings.ReplaceAll(s, "{issue_iid}", url.PathEscape(number))
    s = strings.ReplaceAll(s, "{number}", url.PathEscape(number))
    if strings.Contains(s, "{") {
        return "", errors.New("unresolved placeholder in path template")
    }
    return s, nil
}

// CreateIssue creates an issue and returns its number (string).
func (c *Client) CreateIssue(ctx context.Context, repo string, title, descriptionHTML string) (string, error) {
    if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
        return "", errors.New("missing CNB base URL or token")
    }
    p, err := expand(c.defaultCreatePath(), repo, "")
    if err != nil { return "", err }
    ep, err := c.join(p)
    if err != nil { return "", err }
    // Swagger expects PostIssueForm: title/body/labels/assignees...
    payload := map[string]any{"title": title}
    if strings.TrimSpace(descriptionHTML) != "" { payload["body"] = descriptionHTML }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return "", err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
        return "", fmt.Errorf("cnb create issue status=%d", resp.StatusCode)
    }
    var out struct{ Number string `json:"number"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
    if out.Number == "" { return "", errors.New("empty issue number") }
    return out.Number, nil
}

// UpdateIssue updates fields of an issue.
// number: issue number string per swagger.
func (c *Client) UpdateIssue(ctx context.Context, repo, number string, fields map[string]any) error {
    if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
        return errors.New("missing CNB base URL or token")
    }
    p, err := expand(c.defaultUpdatePath(), repo, number)
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

// AddComment posts a comment to an issue.
func (c *Client) AddComment(ctx context.Context, repo, number, commentHTML string) error {
    if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
        return errors.New("missing CNB base URL or token")
    }
    p, err := expand(c.defaultCommentPath(), repo, number)
    if err != nil { return err }
    ep, err := c.join(p)
    if err != nil { return err }
    payload := map[string]any{"body": commentHTML}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
        return fmt.Errorf("cnb add comment status=%d", resp.StatusCode)
    }
    return nil
}

// CloseIssue transitions an issue to closed via UpdateIssue.
func (c *Client) CloseIssue(ctx context.Context, repo, number string) error {
    return c.UpdateIssue(ctx, repo, number, map[string]any{"state": "closed"})
}
