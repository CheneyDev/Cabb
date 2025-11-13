package cerebras

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

	ai "cabb/internal/ai"
)

type provider struct {
	model   string
	apiKey  string
	baseURL string
}

// New returns a BranchNamer backed by Cerebras chat completions structured outputs.
func New(model, apiKey, baseURL string) ai.BranchNamer {
	return &provider{model: strings.TrimSpace(model), apiKey: strings.TrimSpace(apiKey), baseURL: strings.TrimSpace(baseURL)}
}

func (p *provider) endpoint(paths ...string) (string, error) {
	base := p.baseURL
	if base == "" {
		base = "https://api.cerebras.ai"
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	path := strings.TrimRight(u.Path, "/")
	for _, s := range paths {
		if !strings.HasPrefix(s, "/") {
			s = "/" + s
		}
		path += s
	}
	u.Path = path
	return u.String(), nil
}

func (p *provider) SuggestBranchName(ctx context.Context, title, description string) (string, string, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return "", "", errors.New("missing CEREBRAS_API_KEY")
	}
    model := p.model
    if model == "" {
        model = "zai-glm-4.6"
    }

    schema := map[string]any{
        "type": "object",
        "properties": map[string]any{
            "branch": map[string]any{
                "type":        "string",
                "description": "lowercase Git branch name in prefix/slug format",
            },
        },
        "required":             []string{"branch"},
        "additionalProperties": false,
    }

	prompt := strings.Join([]string{
		"You generate concise Git branch names for software issues.",
		"Rules:",
		"- Return only a JSON object: {\"branch\":\"...\"}.",
		"- branch must be lowercase.",
		"- Allowed prefixes: feat, fix, chore, docs, refactor, test, perf, ci, build, style.",
		"- Format: <prefix>/<slug> where slug uses [a-z0-9_/-], 2..60 chars after prefix/.",
		"- No punctuation, no emojis, no quotes.",
		"- Keep it short and meaningful.",
		"- If information is insufficient or missing, you must still return a branch.",
		"- When unsure, use prefix 'feat' and derive a short slug from the title; if the title is empty, use 'task' as the slug.",
		"- Do not refuse or add explanations; output JSON only.",
		fmt.Sprintf("Title: %s", strings.TrimSpace(title)),
		fmt.Sprintf("Description (may include HTML): %s", strings.TrimSpace(description)),
	}, "\n")

    body := map[string]any{
		"model":                 model,
		// 参考 Playground，最小只传 user 消息，更容易符合约束
		"messages":              []map[string]any{{"role": "user", "content": prompt}},
        "stream":                false,
        "temperature":           0.6,
        "top_p":                 0.95,
        "max_completion_tokens": 128,
        "tools":                 []any{},
        "response_format": map[string]any{
            "type":        "json_schema",
            // 使用 strict=false，避免 beta 期服务端拒绝；本地进行格式校验
            "json_schema": map[string]any{"name": "branch_schema", "strict": false, "schema": schema},
        },
    }
    // reasoning_effort 仅 gpt-oss-120b 支持；其它模型会 400（param=reasoning_effort）
    if strings.EqualFold(model, "gpt-oss-120b") {
        body["reasoning_effort"] = "medium"
    }
	b, _ := json.Marshal(body)
	ep, err := p.endpoint("/v1/chat/completions")
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	hc := &http.Client{Timeout: 12 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var detail string
		if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
			if len(d) > 400 {
				d = d[:400]
			}
			detail = strings.TrimSpace(string(d))
		}
		return "", "", fmt.Errorf("cerebras chat completions status=%d body=%s", resp.StatusCode, detail)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
    if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
        return "", "", errors.New("empty response content")
    }
    raw := strings.TrimSpace(out.Choices[0].Message.Content)
    // Attempt to strip code fences like ```json ... ```
    if strings.HasPrefix(raw, "```") {
        raw = strings.TrimPrefix(raw, "```")
        // optionally strip language tag
        if i := strings.IndexByte(raw, '\n'); i >= 0 { raw = raw[i+1:] }
        if j := strings.LastIndex(raw, "```"); j >= 0 { raw = raw[:j] }
        raw = strings.TrimSpace(raw)
    }
    // Try direct unmarshal; if fails, try to extract first {...}
    var parsed struct { Branch string `json:"branch"` }
    if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
        // Find first JSON object substring heuristically
        if i := strings.Index(raw, "{"); i >= 0 {
            if j := strings.LastIndex(raw, "}"); j > i {
                frag := raw[i:j+1]
                _ = json.Unmarshal([]byte(frag), &parsed)
            }
        }
    }
    branch, ok := ai.SanitizeBranch(parsed.Branch)
    if !ok || branch == "" {
        branch = ai.FallbackBranch(title)
    }
    return branch, "", nil
}
