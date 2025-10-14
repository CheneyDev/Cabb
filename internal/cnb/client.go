package cnb

import (
    "context"
    "errors"
)

// Client is a placeholder for CNB API client.
// NOTE: Outbound CNB API integration is pending confirmation. This client returns Not Implemented errors
// unless feature-flagged and properly configured.
type Client struct {
    BaseURL string
    Token   string
}

var ErrNotImplemented = errors.New("cnb outbound not implemented; 待确认")

func (c *Client) enabled() bool { return c != nil && c.BaseURL != "" && c.Token != "" }

// CreateIssue creates an issue in CNB and returns its IID (string).
func (c *Client) CreateIssue(ctx context.Context, repo string, title, description string) (string, error) {
    if !c.enabled() { return "", ErrNotImplemented }
    // TODO: implement with CNB REST API once spec is confirmed.
    return "", ErrNotImplemented
}

// UpdateIssue updates fields of an issue in CNB.
func (c *Client) UpdateIssue(ctx context.Context, repo, issueIID string, fields map[string]any) error {
    if !c.enabled() { return ErrNotImplemented }
    // TODO: implement
    return ErrNotImplemented
}

// AddComment posts a comment to a CNB issue.
func (c *Client) AddComment(ctx context.Context, repo, issueIID, commentHTML string) error {
    if !c.enabled() { return ErrNotImplemented }
    // TODO: implement
    return ErrNotImplemented
}

// CloseIssue transitions a CNB issue to closed state.
func (c *Client) CloseIssue(ctx context.Context, repo, issueIID string) error {
    if !c.enabled() { return ErrNotImplemented }
    // TODO: implement
    return ErrNotImplemented
}

