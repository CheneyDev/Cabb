package lark

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// Client is a minimal Feishu (Lark) API client.
// NOTE: API paths/fields are based on current docs but marked 待确认 where necessary.
type Client struct {
    AppID     string
    AppSecret string
    BaseURL   string // default https://open.feishu.cn
    HTTP      *http.Client
}

func (c *Client) base() string {
    if c.BaseURL != "" { return c.BaseURL }
    return "https://open.feishu.cn"
}

func (c *Client) httpClient() *http.Client {
    if c.HTTP != nil { return c.HTTP }
    return &http.Client{Timeout: 10 * time.Second}
}

// TenantAccessToken fetches an app tenant access token.
// POST /open-apis/auth/v3/tenant_access_token/internal with {app_id, app_secret}
func (c *Client) TenantAccessToken(ctx context.Context) (string, time.Time, error) {
    type resp struct {
        Code int    `json:"code"`
        Msg  string `json:"msg"`
        TenantAccessToken string `json:"tenant_access_token"`
        Expire int `json:"expire"`
    }
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/auth/v3/tenant_access_token/internal"
    payload := map[string]string{"app_id": c.AppID, "app_secret": c.AppSecret}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return "", time.Time{}, err }
    req.Header.Set("Content-Type", "application/json")
    hc := c.httpClient()
    respRaw, err := hc.Do(req)
    if err != nil { return "", time.Time{}, err }
    defer respRaw.Body.Close()
    var out resp
    if err := json.NewDecoder(respRaw.Body).Decode(&out); err != nil { return "", time.Time{}, err }
    if out.Code != 0 || out.TenantAccessToken == "" {
        if out.Msg == "" { out.Msg = "tenant access token failed" }
        return "", time.Time{}, errors.New(out.Msg)
    }
    exp := time.Now().Add(time.Duration(out.Expire) * time.Second)
    return out.TenantAccessToken, exp, nil
}

// SendTextToChat sends a plain text message to a chat.
// POST /open-apis/im/v1/messages?receive_id_type=chat_id with {receive_id, msg_type, content}
func (c *Client) SendTextToChat(ctx context.Context, tenantToken, chatID, text string) error {
    if tenantToken == "" { return errors.New("missing tenant token") }
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages?receive_id_type=chat_id"
    // content itself is JSON string of {"text": "..."}
    content := map[string]string{"text": text}
    contentJSON, _ := json.Marshal(content)
    payload := map[string]any{
        "receive_id": chatID,
        "msg_type":   "text",
        "content":    string(contentJSON),
    }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+tenantToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return errors.New("lark send message status!=2xx") }
    return nil
}

// ReplyTextInThread replies to a message thread (best-effort; API path 待确认)
// POST /open-apis/im/v1/messages/{message_id}/reply with {msg_type, content}
func (c *Client) ReplyTextInThread(ctx context.Context, tenantToken, rootMessageID, text string) error {
    if tenantToken == "" { return errors.New("missing tenant token") }
    if rootMessageID == "" { return errors.New("missing root message id") }
    // Sanitize message id in path
    pathID := url.PathEscape(rootMessageID)
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages/" + pathID + "/reply"
    content := map[string]string{"text": text}
    contentJSON, _ := json.Marshal(content)
    payload := map[string]any{"msg_type": "text", "content": string(contentJSON)}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+tenantToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return errors.New("lark reply message status!=2xx") }
    return nil
}

