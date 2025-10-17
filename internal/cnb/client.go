package cnb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 10 * time.Second}
}

func (c *Client) join(path string) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	return u.String(), nil
}

func (c *Client) defaultCreatePath() string {
	if c.IssueCreatePath != "" {
		return c.IssueCreatePath
	}
	return "/{repo}/-/issues"
}
func (c *Client) defaultUpdatePath() string {
	if c.IssueUpdatePath != "" {
		return c.IssueUpdatePath
	}
	return "/{repo}/-/issues/{number}"
}
func (c *Client) defaultCommentPath() string {
	if c.IssueCommentPath != "" {
		return c.IssueCommentPath
	}
	return "/{repo}/-/issues/{number}/comments"
}

func (c *Client) defaultAssigneesPath() string {
    return "/{repo}/-/issues/{number}/assignees"
}

// defaultGetIssuePath returns the path for getting a single issue detail
func (c *Client) defaultGetIssuePath() string {
    return "/{repo}/-/issues/{number}"
}

func expand(tpl string, repo, number string) (string, error) {
	if tpl == "" {
		return "", errors.New("empty path template")
	}
	// Encode repo path by segments so that slashes remain separators
	// e.g., "org/repo" => "org/repo" but each segment escaped
	repoEsc := repo
	if repo != "" {
		parts := strings.Split(repo, "/")
		for i, p := range parts {
			parts[i] = url.PathEscape(p)
		}
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

// Label represents a minimal label entity from CNB.
type Label struct {
    Name  string `json:"name"`
    Color string `json:"color"`
}

// IssueDetail contains a small subset we care about for verification
type IssueDetail struct {
    Number   string `json:"number"`
    Priority string `json:"priority"`
    Title    string `json:"title"`
}

// GetIssue retrieves an issue's details
func (c *Client) GetIssue(ctx context.Context, repo, number string) (*IssueDetail, error) {
    if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
        return nil, errors.New("missing CNB base URL or token")
    }
    p, err := expand(c.defaultGetIssuePath(), repo, number)
    if err != nil {
        return nil, err
    }
    ep, err := c.join(p)
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Accept", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusNotAcceptable {
        _ = resp.Body.Close()
        req2, err2 := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
        if err2 != nil {
            return nil, err2
        }
        req2.Header.Set("Authorization", "Bearer "+c.Token)
        req2.Header.Set("Accept", "application/vnd.cnb.api+json")
        resp2, err2 := c.httpClient().Do(req2)
        if err2 != nil {
            return nil, err2
        }
        defer resp2.Body.Close()
        if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
            var out IssueDetail
            if err := json.NewDecoder(resp2.Body).Decode(&out); err != nil {
                return nil, err
            }
            return &out, nil
        }
        return nil, fmt.Errorf("cnb get issue status=%d", resp2.StatusCode)
    }
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("cnb get issue status=%d", resp.StatusCode)
    }
    var out IssueDetail
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, err
    }
    return &out, nil
}

// EnsureRepoLabels makes sure the given label names exist in the repository.
// It creates missing labels with a default color.
func (c *Client) EnsureRepoLabels(ctx context.Context, repo string, names []string) error {
	if len(names) == 0 {
		return nil
	}
	// Build a quick lookup of existing labels by name (case-insensitive match by exact name)
	existing := map[string]struct{}{}
	// Query by keyword per name to reduce payload; fallback to empty keyword list if needed.
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		ls, _ := c.ListRepoLabels(ctx, repo, n)
		for _, l := range ls {
			if strings.EqualFold(l.Name, n) {
				existing[strings.ToLower(n)] = struct{}{}
				break
			}
		}
	}
	// Create missing labels
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		if _, ok := existing[strings.ToLower(n)]; ok {
			continue
		}
		if err := c.CreateRepoLabel(ctx, repo, n, defaultLabelColor(n)); err != nil {
			return err
		}
	}
	return nil
}

// EnsureRepoLabelsWithColors ensures labels exist with optional colors.
// If a color is provided (non-empty, may include leading '#'), it will be set when creating the label.
func (c *Client) EnsureRepoLabelsWithColors(ctx context.Context, repo string, want []Label) error {
	if len(want) == 0 {
		return nil
	}
	// Build existing set
	existing := map[string]struct{}{}
	for _, l := range want {
		n := strings.TrimSpace(l.Name)
		if n == "" {
			continue
		}
		ls, _ := c.ListRepoLabels(ctx, repo, n)
		for _, e := range ls {
			if strings.EqualFold(e.Name, n) {
				existing[strings.ToLower(n)] = struct{}{}
				break
			}
		}
	}
	for _, l := range want {
		n := strings.TrimSpace(l.Name)
		if n == "" {
			continue
		}
		if _, ok := existing[strings.ToLower(n)]; ok {
			continue
		}
		color := strings.TrimPrefix(strings.TrimSpace(l.Color), "#")
		if color == "" {
			color = defaultLabelColor(n)
		}
		if err := c.CreateRepoLabel(ctx, repo, n, color); err != nil {
			return err
		}
	}
	return nil
}

// ListRepoLabels lists labels for a repo, optionally filtered by keyword (name contains).
func (c *Client) ListRepoLabels(ctx context.Context, repo, keyword string) ([]Label, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return nil, errors.New("missing CNB base URL or token")
	}
	p, err := expand("/{repo}/-/labels", repo, "")
	if err != nil {
		return nil, err
	}
	ep, err := c.join(p)
	if err != nil {
		return nil, err
	}
	if keyword != "" {
		u, _ := url.Parse(ep)
		q := u.Query()
		q.Set("keyword", keyword)
		u.RawQuery = q.Encode()
		ep = u.String()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cnb list labels status=%d", resp.StatusCode)
	}
	var out []Label
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateRepoLabel creates a label in the repository.
func (c *Client) CreateRepoLabel(ctx context.Context, repo, name, color string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("label name required")
	}
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return errors.New("missing CNB base URL or token")
	}
	p, err := expand("/{repo}/-/labels", repo, "")
	if err != nil {
		return err
	}
	ep, err := c.join(p)
	if err != nil {
		return err
	}
	payload := map[string]any{"name": name}
	if strings.TrimSpace(color) != "" {
		payload["color"] = strings.TrimPrefix(color, "#")
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
		if err2 != nil {
			return err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("cnb create label status=%d", resp2.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cnb create label status=%d", resp.StatusCode)
	}
	return nil
}

// SetIssueLabels sets the full labels list for an issue (replaces existing).
func (c *Client) SetIssueLabels(ctx context.Context, repo, number string, names []string) error {
	if len(names) == 0 {
		return nil
	}
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return errors.New("missing CNB base URL or token")
	}
	p, err := expand("/{repo}/-/issues/{number}/labels", repo, number)
	if err != nil {
		return err
	}
	ep, err := c.join(p)
	if err != nil {
		return err
	}
	payload := map[string]any{"labels": names}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPut, ep, bytes.NewReader(b))
		if err2 != nil {
			return err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
			return nil
		}
		return fmt.Errorf("cnb set issue labels status=%d", resp2.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cnb set issue labels status=%d", resp.StatusCode)
	}
	return nil
}

func defaultLabelColor(name string) string {
	// A neutral default color without '#'; platforms usually accept 6-hex RGB
	return "ededed"
}

// CreateIssue creates an issue and returns its number (string).
func (c *Client) CreateIssue(ctx context.Context, repo string, title, descriptionHTML string) (string, error) {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return "", errors.New("missing CNB base URL or token")
	}
	p, err := expand(c.defaultCreatePath(), repo, "")
	if err != nil {
		return "", err
	}
	ep, err := c.join(p)
	if err != nil {
		return "", err
	}
	// Swagger expects PostIssueForm: title/body/labels/assignees...
	payload := map[string]any{"title": title}
	if strings.TrimSpace(descriptionHTML) != "" {
		payload["body"] = descriptionHTML
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	// Some CNB gateways return 406 when Accept lists multiple types; use a single type
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		// Retry with vendor-specific Accept as some gateways require an exact single media type
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
		if err2 != nil {
			return "", err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return "", err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == http.StatusCreated || (resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
			var out2 struct {
				Number string `json:"number"`
			}
			if err := json.NewDecoder(resp2.Body).Decode(&out2); err != nil {
				return "", err
			}
			if out2.Number == "" {
				return "", errors.New("empty issue number")
			}
			return out2.Number, nil
		}
		var msg2 string
		if d, _ := io.ReadAll(resp2.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg2 = string(d)
		}
		u2, _ := url.Parse(ep)
		return "", fmt.Errorf("cnb create issue status=%d path=%s body=%s", resp2.StatusCode, u2.EscapedPath(), msg2)
	}
	if resp.StatusCode != http.StatusCreated && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		// read small error body for diagnostics (truncate)
		var msg string
		if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg = string(d)
		}
		u, _ := url.Parse(ep)
		return "", fmt.Errorf("cnb create issue status=%d path=%s body=%s", resp.StatusCode, u.EscapedPath(), msg)
	}
	var out struct {
		Number string `json:"number"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Number == "" {
		return "", errors.New("empty issue number")
	}
	return out.Number, nil
}

// UpdateIssue updates fields of an issue.
// number: issue number string per swagger.
func (c *Client) UpdateIssue(ctx context.Context, repo, number string, fields map[string]any) error {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return errors.New("missing CNB base URL or token")
	}
	p, err := expand(c.defaultUpdatePath(), repo, number)
	if err != nil {
		return err
	}
	ep, err := c.join(p)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(fields)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
		if err2 != nil {
			return err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
			return nil
		}
		var msg2 string
		if d, _ := io.ReadAll(resp2.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg2 = string(d)
		}
		u2, _ := url.Parse(ep)
		return fmt.Errorf("cnb update issue status=%d path=%s body=%s", resp2.StatusCode, u2.EscapedPath(), msg2)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var msg string
		if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg = string(d)
		}
		u, _ := url.Parse(ep)
		return fmt.Errorf("cnb update issue status=%d path=%s body=%s", resp.StatusCode, u.EscapedPath(), msg)
	}
	return nil
}

// AddComment posts a comment to an issue.
func (c *Client) AddComment(ctx context.Context, repo, number, commentHTML string) error {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return errors.New("missing CNB base URL or token")
	}
	p, err := expand(c.defaultCommentPath(), repo, number)
	if err != nil {
		return err
	}
	ep, err := c.join(p)
	if err != nil {
		return err
	}
	payload := map[string]any{"body": commentHTML}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
		if err2 != nil {
			return err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == http.StatusCreated || (resp2.StatusCode >= 200 && resp2.StatusCode < 300) {
			return nil
		}
		var msg2 string
		if d, _ := io.ReadAll(resp2.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg2 = string(d)
		}
		u2, _ := url.Parse(ep)
		return fmt.Errorf("cnb add comment status=%d path=%s body=%s", resp2.StatusCode, u2.EscapedPath(), msg2)
	}
	if resp.StatusCode != http.StatusCreated && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		var msg string
		if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg = string(d)
		}
		u, _ := url.Parse(ep)
		return fmt.Errorf("cnb add comment status=%d path=%s body=%s", resp.StatusCode, u.EscapedPath(), msg)
	}
	return nil
}

// CloseIssue transitions an issue to closed via UpdateIssue.
func (c *Client) CloseIssue(ctx context.Context, repo, number string) error {
	return c.UpdateIssue(ctx, repo, number, map[string]any{"state": "closed"})
}

// UpdateIssueAssignees replaces the assignees list of an issue via PATCH.
// According to swagger, body schema is api.PatchIssueAssigneesForm: { assignees: [string] }
func (c *Client) UpdateIssueAssignees(ctx context.Context, repo, number string, assignees []string) error {
	if strings.TrimSpace(c.BaseURL) == "" || strings.TrimSpace(c.Token) == "" {
		return errors.New("missing CNB base URL or token")
	}
	p, err := expand(c.defaultAssigneesPath(), repo, number)
	if err != nil {
		return err
	}
	ep, err := c.join(p)
	if err != nil {
		return err
	}
	payload := map[string]any{"assignees": assignees}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotAcceptable {
		_ = resp.Body.Close()
		req2, err2 := http.NewRequestWithContext(ctx, http.MethodPatch, ep, bytes.NewReader(b))
		if err2 != nil {
			return err2
		}
		req2.Header.Set("Authorization", "Bearer "+c.Token)
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Accept", "application/vnd.cnb.api+json")
		resp2, err2 := c.httpClient().Do(req2)
		if err2 != nil {
			return err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode >= 200 && resp2.StatusCode < 300 {
			return nil
		}
		var msg2 string
		if d, _ := io.ReadAll(resp2.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg2 = string(d)
		}
		u2, _ := url.Parse(ep)
		return fmt.Errorf("cnb update assignees status=%d path=%s body=%s", resp2.StatusCode, u2.EscapedPath(), msg2)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var msg string
		if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
			if len(d) > 300 {
				d = d[:300]
			}
			msg = string(d)
		}
		u, _ := url.Parse(ep)
		return fmt.Errorf("cnb update assignees status=%d path=%s body=%s", resp.StatusCode, u.EscapedPath(), msg)
	}
	return nil
}
