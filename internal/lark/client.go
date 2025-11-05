package lark

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

// Client is a minimal Feishu (Lark) API client.
// NOTE: API paths/fields are based on current docs but marked 待确认 where necessary.
type Client struct {
	AppID     string
	AppSecret string
	BaseURL   string // default https://open.feishu.cn
	HTTP      *http.Client
}

type Chat struct {
	ID        string   `json:"chat_id"`
	Name      string   `json:"name"`
	I18nNames ChatI18n `json:"i18n_names"`
}

type ChatI18n struct {
	ZhCN string `json:"zh_cn"`
	EnUS string `json:"en_us"`
	JaJP string `json:"ja_jp"`
}

func (c *Client) base() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return "https://open.feishu.cn"
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return &http.Client{Timeout: 10 * time.Second}
}

// TenantAccessToken fetches an app tenant access token.
// POST /open-apis/auth/v3/tenant_access_token/internal with {app_id, app_secret}
func (c *Client) TenantAccessToken(ctx context.Context) (string, time.Time, error) {
	type resp struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int    `json:"expire"`
	}
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/auth/v3/tenant_access_token/internal"
	payload := map[string]string{"app_id": c.AppID, "app_secret": c.AppSecret}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	hc := c.httpClient()
	respRaw, err := hc.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer respRaw.Body.Close()
	var out resp
	if err := json.NewDecoder(respRaw.Body).Decode(&out); err != nil {
		return "", time.Time{}, err
	}
	if out.Code != 0 || out.TenantAccessToken == "" {
		if out.Msg == "" {
			out.Msg = "tenant access token failed"
		}
		return "", time.Time{}, errors.New(out.Msg)
	}
	exp := time.Now().Add(time.Duration(out.Expire) * time.Second)
	return out.TenantAccessToken, exp, nil
}

// SendTextToChat sends a plain text message to a chat.
// POST /open-apis/im/v1/messages?receive_id_type=chat_id with {receive_id, msg_type, content}
func (c *Client) SendTextToChat(ctx context.Context, tenantToken, chatID, text string) error {
	if tenantToken == "" {
		return errors.New("missing tenant token")
	}
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
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("lark send message status!=2xx")
	}
	return nil
}

// ReplyTextInThread replies to a message thread (best-effort; API path 待确认)
// POST /open-apis/im/v1/messages/{message_id}/reply with {msg_type, content}
func (c *Client) ReplyTextInThread(ctx context.Context, tenantToken, rootMessageID, text string) error {
	if tenantToken == "" {
		return errors.New("missing tenant token")
	}
	if rootMessageID == "" {
		return errors.New("missing root message id")
	}
	// Sanitize message id in path
	pathID := url.PathEscape(rootMessageID)
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages/" + pathID + "/reply"
	content := map[string]string{"text": text}
	contentJSON, _ := json.Marshal(content)
	payload := map[string]any{"msg_type": "text", "content": string(contentJSON)}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("lark reply message status!=2xx")
	}
	return nil
}

// GetChat fetches metadata for a chat (group) by chat_id.
// GET /open-apis/im/v1/chats/:chat_id per docs (docs/feishu/003-服务端API/群组/群组管理/获取群信息.md).
func (c *Client) GetChat(ctx context.Context, tenantToken, chatID string) (*Chat, error) {
	if tenantToken == "" {
		return nil, errors.New("missing tenant token")
	}
	chatID = strings.TrimSpace(chatID)
	if chatID == "" {
		return nil, errors.New("missing chat id")
	}
	pathID := url.PathEscape(chatID)
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/chats/" + pathID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("lark get chat status=%d", resp.StatusCode)
	}
	var raw struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Chat     Chat     `json:"chat"`
			ChatInfo Chat     `json:"chat_info"`
			ChatID   string   `json:"chat_id"`
			Name     string   `json:"name"`
			I18n     ChatI18n `json:"i18n_names"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if raw.Code != 0 {
		if raw.Msg == "" {
			raw.Msg = "lark get chat failed"
		}
		return nil, errors.New(raw.Msg)
	}
	chat := raw.Data.ChatInfo
	if chat.ID == "" && raw.Data.Chat.ID != "" {
		chat = raw.Data.Chat
	}
	if chat.ID == "" {
		chat.ID = raw.Data.ChatID
	}
	if chat.Name == "" {
		chat.Name = raw.Data.Name
	}
	if chat.I18nNames.ZhCN == "" {
		chat.I18nNames.ZhCN = raw.Data.I18n.ZhCN
	}
	if chat.I18nNames.EnUS == "" {
		chat.I18nNames.EnUS = raw.Data.I18n.EnUS
	}
	if chat.I18nNames.JaJP == "" {
		chat.I18nNames.JaJP = raw.Data.I18n.JaJP
	}
	if chat.Name == "" {
		if chat.I18nNames.ZhCN != "" {
			chat.Name = chat.I18nNames.ZhCN
		} else if chat.I18nNames.EnUS != "" {
			chat.Name = chat.I18nNames.EnUS
		} else if chat.I18nNames.JaJP != "" {
			chat.Name = chat.I18nNames.JaJP
		}
	}
	return &chat, nil
}

// SendPostToChat sends a post (rich text) message to a chat.
// POST /open-apis/im/v1/messages?receive_id_type=chat_id with msg_type=post
func (c *Client) SendPostToChat(ctx context.Context, tenantToken, chatID string, post map[string]any) error {
	if tenantToken == "" {
		return errors.New("missing tenant token")
	}
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages?receive_id_type=chat_id"
	contentJSON, _ := json.Marshal(post)
	payload := map[string]any{
		"receive_id": chatID,
		"msg_type":   "post",
		"content":    string(contentJSON),
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("lark send post status!=2xx")
	}
	return nil
}

// ReplyPostInThread replies to a thread with a post (rich text) message.
// POST /open-apis/im/v1/messages/{message_id}/reply with msg_type=post
func (c *Client) ReplyPostInThread(ctx context.Context, tenantToken, rootMessageID string, post map[string]any) error {
	if tenantToken == "" {
		return errors.New("missing tenant token")
	}
	if rootMessageID == "" {
		return errors.New("missing root message id")
	}
	pathID := url.PathEscape(rootMessageID)
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages/" + pathID + "/reply"
	contentJSON, _ := json.Marshal(post)
	payload := map[string]any{"msg_type": "post", "content": string(contentJSON)}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("lark reply post status!=2xx")
	}
	return nil
}

// SendCardToChat sends an interactive card message to a chat.
// POST /open-apis/im/v1/messages?receive_id_type=chat_id with msg_type=interactive and content={"card": <card>}
func (c *Client) SendCardToChat(ctx context.Context, tenantToken, chatID string, card map[string]any) error {
    if tenantToken == "" {
        return errors.New("missing tenant token")
    }
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages?receive_id_type=chat_id"
    wrapper := map[string]any{"card": card}
    contentJSON, _ := json.Marshal(wrapper)
    payload := map[string]any{
        "receive_id": chatID,
        "msg_type":   "interactive",
        "content":    string(contentJSON),
    }
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+tenantToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return errors.New("lark send card status!=2xx")
    }
    return nil
}

// ReplyCardInThread replies to a thread with an interactive card message.
// POST /open-apis/im/v1/messages/{message_id}/reply with msg_type=interactive
func (c *Client) ReplyCardInThread(ctx context.Context, tenantToken, rootMessageID string, card map[string]any) error {
    if tenantToken == "" {
        return errors.New("missing tenant token")
    }
    if rootMessageID == "" {
        return errors.New("missing root message id")
    }
    pathID := url.PathEscape(rootMessageID)
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages/" + pathID + "/reply"
    wrapper := map[string]any{"card": card}
    contentJSON, _ := json.Marshal(wrapper)
    payload := map[string]any{"msg_type": "interactive", "content": string(contentJSON)}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+tenantToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return errors.New("lark reply card status!=2xx")
    }
    return nil
}
