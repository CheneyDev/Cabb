package lark

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

// Client is a minimal Feishu (Lark) API client.
// NOTE: API paths/fields are based on current docs but marked 待确认 where necessary.
type Client struct {
	AppID     string
	AppSecret string
	BaseURL   string // default https://open.feishu.cn
	HTTP      *http.Client
}

func NewClient(appID, appSecret string) *Client {
	return &Client{
		AppID:     appID,
		AppSecret: appSecret,
	}
}

func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	token, _, err := c.TenantAccessToken(ctx)
	if err != nil {
		return err
	}
	return c.SendTextToChat(ctx, token, chatID, text)
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark send message status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark reply message status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark send post status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark reply post status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
	}
	return nil
}

// SendCardToChat sends an interactive card message to a chat.
// POST /open-apis/im/v1/messages?receive_id_type=chat_id with msg_type=interactive and content=<card-json-string>
func (c *Client) SendCardToChat(ctx context.Context, tenantToken, chatID string, card map[string]any) error {
    if tenantToken == "" {
        return errors.New("missing tenant token")
    }
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages?receive_id_type=chat_id"
    contentJSON, _ := json.Marshal(card)
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark send card status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
    }
    return nil
}

// ReplyCardInThread replies to a thread with an interactive card message.
// POST /open-apis/im/v1/messages/{message_id}/reply with msg_type=interactive and content=<card-json-string>
func (c *Client) ReplyCardInThread(ctx context.Context, tenantToken, rootMessageID string, card map[string]any) error {
    if tenantToken == "" {
        return errors.New("missing tenant token")
    }
    if rootMessageID == "" {
        return errors.New("missing root message id")
    }
    pathID := url.PathEscape(rootMessageID)
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages/" + pathID + "/reply"
    contentJSON, _ := json.Marshal(card)
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark reply card status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
    }
    return nil
}

// UpdateInteractiveCard performs delayed update of a previously sent card via callback token.
// POST /open-apis/interactive/v1/card/update with { token: c-xxxx, card: <card-json-object> }
func (c *Client) UpdateInteractiveCard(ctx context.Context, tenantToken, callbackToken string, card map[string]any) error {
    if tenantToken == "" {
        return errors.New("missing tenant token")
    }
    if strings.TrimSpace(callbackToken) == "" {
        return errors.New("missing callback token")
    }
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/interactive/v1/card/update"
    payload := map[string]any{
        "token": callbackToken,
        "card":  card,
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
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return fmt.Errorf("lark update card status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
    }
    return nil
}

// ChatShareLink represents a share link response from Feishu API
type ChatShareLink struct {
    ShareLink   string `json:"share_link"`
    ExpireTime  string `json:"expire_time"`
    IsPermanent bool   `json:"is_permanent"`
}

// GetChatShareLink fetches the share link for a chat group.
// POST /open-apis/im/v1/chats/:chat_id/link
// Per docs/feishu/003-服务端API/群组/群组管理/获取群分享链接.md
func (c *Client) GetChatShareLink(ctx context.Context, tenantToken, chatID, validityPeriod string) (*ChatShareLink, error) {
    if tenantToken == "" {
        return nil, errors.New("missing tenant token")
    }
    if strings.TrimSpace(chatID) == "" {
        return nil, errors.New("missing chat id")
    }
    pathID := url.PathEscape(chatID)
    ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/chats/" + pathID + "/link"
    // Default to week if not specified
    if validityPeriod == "" {
        validityPeriod = "week"
    }
    payload := map[string]any{"validity_period": validityPeriod}
    b, _ := json.Marshal(payload)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+tenantToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient().Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
        _ = json.Unmarshal(body, &er)
        return nil, fmt.Errorf("lark get chat share link status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
    }
    var result struct {
        Code int    `json:"code"`
        Msg  string `json:"msg"`
        Data ChatShareLink `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    if result.Code != 0 {
        if result.Msg == "" {
            result.Msg = "get chat share link failed"
        }
        return nil, errors.New(result.Msg)
    }
    return &result.Data, nil
}

// User represents a Lark user
type User struct {
	UnionID     string `json:"union_id"`
	UserID      string `json:"user_id"`
	OpenID      string `json:"open_id"`
	Name        string `json:"name"`
	EnName      string `json:"en_name"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Avatar      struct {
		Avatar72 string `json:"avatar_72"`
	} `json:"avatar"`
}

// FindByDepartment fetches users in a department.
// GET /open-apis/contact/v3/users/find_by_department
// FindByDepartment fetches users in a department.
// GET /open-apis/contact/v3/users/find_by_department
func (c *Client) FindByDepartment(ctx context.Context, tenantToken string, departmentID string, pageSize int, pageToken string) ([]User, string, bool, error) {
	if tenantToken == "" {
		return nil, "", false, errors.New("missing tenant token")
	}
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/contact/v3/users/find_by_department"
	
	// Build query params
	params := url.Values{}
	params.Set("department_id", departmentID)
	params.Set("department_id_type", "open_department_id") // Default
	params.Set("user_id_type", "open_id") // Default, compatible with person_list component
	if pageSize > 0 {
		params.Set("page_size", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		params.Set("page_token", pageToken)
	}
	
	ep += "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, "", false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
		_ = json.Unmarshal(body, &er)
		return nil, "", false, fmt.Errorf("lark find users status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			HasMore   bool   `json:"has_more"`
			PageToken string `json:"page_token"`
			Items     []User `json:"items"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", false, err
	}
	
	if result.Code != 0 {
		return nil, "", false, fmt.Errorf("lark find users failed: %s", result.Msg)
	}
	
	return result.Data.Items, result.Data.PageToken, result.Data.HasMore, nil
}

type Department struct {
	Name             string `json:"name"`
	DepartmentID     string `json:"department_id"`
	OpenDepartmentID string `json:"open_department_id"`
	ParentDepartmentID string `json:"parent_department_id"`
}

// GetChildrenDepartment fetches sub-departments.
// GET /open-apis/contact/v3/departments/:department_id/children
func (c *Client) GetChildrenDepartment(ctx context.Context, tenantToken string, departmentID string, fetchChild bool, pageSize int, pageToken string) ([]Department, string, bool, error) {
	if tenantToken == "" {
		return nil, "", false, errors.New("missing tenant token")
	}
	// Default to root if empty
	if departmentID == "" {
		departmentID = "0"
	}
	
	pathID := url.PathEscape(departmentID)
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/contact/v3/departments/" + pathID + "/children"
	
	params := url.Values{}
	params.Set("department_id_type", "open_department_id")
	if fetchChild {
		params.Set("fetch_child", "true")
	}
	if pageSize > 0 {
		params.Set("page_size", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		params.Set("page_token", pageToken)
	}
	ep += "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
		_ = json.Unmarshal(body, &er)
		return nil, "", false, fmt.Errorf("lark get children department status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			HasMore   bool         `json:"has_more"`
			PageToken string       `json:"page_token"`
			Items     []Department `json:"items"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", false, err
	}

	if result.Code != 0 {
		return nil, "", false, fmt.Errorf("lark get children department failed: %s", result.Msg)
	}

	return result.Data.Items, result.Data.PageToken, result.Data.HasMore, nil
}

// FindAllUsers recursively fetches all users from all departments.
func (c *Client) FindAllUsers(ctx context.Context, tenantToken string) ([]User, error) {
	// 1. Get all departments (recursive)
	var allDepts []Department
	// Add root department manually as GetChildrenDepartment with fetch_child=true returns children, not root itself
	allDepts = append(allDepts, Department{OpenDepartmentID: "0"})
	
	pageToken := ""
	for {
		depts, nextToken, hasMore, err := c.GetChildrenDepartment(ctx, tenantToken, "0", true, 50, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch departments: %w", err)
		}
		allDepts = append(allDepts, depts...)
		if !hasMore {
			break
		}
		pageToken = nextToken
	}

	// 2. Fetch users for each department
	userMap := make(map[string]User)
	for _, dept := range allDepts {
		deptID := dept.OpenDepartmentID
		if deptID == "" {
			deptID = dept.DepartmentID
		}
		if deptID == "" {
			continue
		}

		userPageToken := ""
		for {
			users, nextToken, hasMore, err := c.FindByDepartment(ctx, tenantToken, deptID, 50, userPageToken)
			if err != nil {
				// Log error but continue? Or fail? For now, fail to be safe.
				// In production, might want to skip inaccessible departments.
				// return nil, fmt.Errorf("failed to fetch users for dept %s: %w", deptID, err)
				// Let's just log and continue to be robust
				fmt.Printf("Warning: failed to fetch users for dept %s: %v\n", deptID, err)
				break 
			}
			for _, u := range users {
				userMap[u.OpenID] = u
			}
			if !hasMore {
				break
			}
			userPageToken = nextToken
		}
	}

	// 3. Convert map to slice
	var allUsers []User
	for _, u := range userMap {
		allUsers = append(allUsers, u)
	}
	return allUsers, nil
}

// BatchSendResult represents the result of a batch send operation
type BatchSendResult struct {
	MessageID            string   `json:"message_id"`
	InvalidDepartmentIDs []string `json:"invalid_department_ids"`
	InvalidOpenIDs       []string `json:"invalid_open_ids"`
	InvalidUserIDs       []string `json:"invalid_user_ids"`
	InvalidUnionIDs      []string `json:"invalid_union_ids"`
}

// BatchSendCard sends an interactive card to multiple users or departments.
// POST /open-apis/message/v4/batch_send/
// At least one of openIDs, userIDs, unionIDs, or departmentIDs must be provided.
func (c *Client) BatchSendCard(ctx context.Context, tenantToken string, openIDs, departmentIDs []string, card map[string]any) (*BatchSendResult, error) {
	if tenantToken == "" {
		return nil, errors.New("missing tenant token")
	}
	if len(openIDs) == 0 && len(departmentIDs) == 0 {
		return nil, errors.New("at least one of openIDs or departmentIDs must be provided")
	}

	ep := strings.TrimRight(c.base(), "/") + "/open-apis/message/v4/batch_send/"

	payload := map[string]any{
		"msg_type": "interactive",
		"card":     card,
	}
	if len(openIDs) > 0 {
		payload["open_ids"] = openIDs
	}
	if len(departmentIDs) > 0 {
		payload["department_ids"] = departmentIDs
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			MessageID            string   `json:"message_id"`
			InvalidDepartmentIDs []string `json:"invalid_department_ids"`
			InvalidOpenIDs       []string `json:"invalid_open_ids"`
			InvalidUserIDs       []string `json:"invalid_user_ids"`
			InvalidUnionIDs      []string `json:"invalid_union_ids"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("batch send failed: code=%d msg=%s", result.Code, result.Msg)
	}

	return &BatchSendResult{
		MessageID:            result.Data.MessageID,
		InvalidDepartmentIDs: result.Data.InvalidDepartmentIDs,
		InvalidOpenIDs:       result.Data.InvalidOpenIDs,
		InvalidUserIDs:       result.Data.InvalidUserIDs,
		InvalidUnionIDs:      result.Data.InvalidUnionIDs,
	}, nil
}

// GetDepartment fetches a single department info.
// GET /open-apis/contact/v3/departments/:department_id
func (c *Client) GetDepartment(ctx context.Context, tenantToken, departmentID string) (*Department, error) {
	if tenantToken == "" {
		return nil, errors.New("missing tenant token")
	}
	if departmentID == "" {
		departmentID = "0"
	}

	pathID := url.PathEscape(departmentID)
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/contact/v3/departments/" + pathID + "?department_id_type=open_department_id"

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

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Department Department `json:"department"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("get department failed: %s", result.Msg)
	}
	return &result.Data.Department, nil
}

// SendCardToUser sends an interactive card message to a single user via open_id.
// POST /open-apis/im/v1/messages?receive_id_type=open_id
func (c *Client) SendCardToUser(ctx context.Context, tenantToken, openID string, card map[string]any) error {
	if tenantToken == "" {
		return errors.New("missing tenant token")
	}
	if strings.TrimSpace(openID) == "" {
		return errors.New("missing open_id")
	}
	ep := strings.TrimRight(c.base(), "/") + "/open-apis/im/v1/messages?receive_id_type=open_id"
	contentJSON, _ := json.Marshal(card)
	payload := map[string]any{
		"receive_id": openID,
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
		body, _ := io.ReadAll(resp.Body)
		var er struct{ Code int `json:"code"`; Msg string `json:"msg"` }
		_ = json.Unmarshal(body, &er)
		return fmt.Errorf("lark send card to user status=%d code=%d msg=%s", resp.StatusCode, er.Code, strings.TrimSpace(er.Msg))
	}
	return nil
}

// ListAllDepartments fetches all departments recursively from root.
func (c *Client) ListAllDepartments(ctx context.Context, tenantToken string) ([]Department, error) {
	var allDepts []Department

	// Add root department
	root, err := c.GetDepartment(ctx, tenantToken, "0")
	if err == nil && root != nil {
		allDepts = append(allDepts, *root)
	}

	// Get all children recursively
	pageToken := ""
	for {
		depts, nextToken, hasMore, err := c.GetChildrenDepartment(ctx, tenantToken, "0", true, 50, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch departments: %w", err)
		}
		allDepts = append(allDepts, depts...)
		if !hasMore {
			break
		}
		pageToken = nextToken
	}

	return allDepts, nil
}
