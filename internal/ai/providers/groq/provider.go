package groq

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

// New returns a BranchNamer backed by Groq Chat Completions Structured Outputs.
func New(model, apiKey, baseURL string) ai.BranchNamer {
    return &provider{model: strings.TrimSpace(model), apiKey: strings.TrimSpace(apiKey), baseURL: strings.TrimSpace(baseURL)}
}

func (p *provider) endpoint(paths ...string) (string, error) {
    base := p.baseURL
    if base == "" {
        base = "https://api.groq.com/openai"
    }
    u, err := url.Parse(base)
    if err != nil {
        return "", err
    }
    path := strings.TrimRight(u.Path, "/")
    for _, s := range paths {
        if !strings.HasPrefix(s, "/") { s = "/" + s }
        path += s
    }
    u.Path = path
    return u.String(), nil
}

func (p *provider) SuggestBranchName(ctx context.Context, title, description string) (string, string, error) {
    if strings.TrimSpace(p.apiKey) == "" {
        return "", "", errors.New("missing GROQ_API_KEY")
    }
    model := strings.TrimSpace(p.model)
    if model == "" { model = "moonshotai/kimi-k2-instruct-0905" }

    schema := map[string]any{
        "type": "object",
        "properties": map[string]any{
            "branch": map[string]any{
                "type":        "string",
                "description": "lowercase Git branch name in prefix/slug format",
                // At most 4 hyphen-separated words after the prefix/ (e.g., feat/my-short-slug)
                "pattern":     "^(feat|fix|chore|docs|refactor|test|perf|ci|build|style)(/[a-z0-9_]+(?:-[a-z0-9_]+){0,3})$",
            },
            "reason": map[string]any{
                "anyOf": []any{map[string]any{"type": "string"}, map[string]any{"type": "null"}},
            },
        },
        "required":             []string{"branch", "reason"},
        "additionalProperties": false,
    }

    sys := strings.Join([]string{
        "You generate concise Git branch names for software issues.",
        "Rules:",
        "- Return only JSON with fields: branch, reason.",
        "- branch must be lowercase.",
        "- Allowed prefixes: feat, fix, chore, docs, refactor, test, perf, ci, build, style.",
        "- Format: <prefix>/<slug> where slug uses [a-z0-9_/-], no spaces, 2..60 chars after prefix/.",
        "- No punctuation, no emojis, no quotes.",
        "- Keep it short and meaningful.",
        "- Limit slug to at most 4 hyphen-separated words.",
        "- 中文：分支名单词（用连字符-分隔）不超过4个。",
        "- If information is insufficient or missing, you must still return a branch.",
        "- Output must be ASCII; transliterate or simplify non-ASCII to ASCII.",
        "- When unsure, use prefix 'feat' and derive a short ASCII slug from the title; if the title is empty, use 'task' as the slug.",
        "- Do not refuse or add explanations; output JSON only.",
    }, "\n")
    user := fmt.Sprintf("Title: %s\nDescription (may include HTML): %s\nReturn a fitting branch.", strings.TrimSpace(title), strings.TrimSpace(description))

    body := map[string]any{
        "model": model,
        "messages": []map[string]any{
            {"role": "system", "content": sys},
            {"role": "user", "content": user},
        },
        "response_format": map[string]any{
            "type": "json_schema",
            "json_schema": map[string]any{ "name": "branch_schema", "schema": schema },
        },
        "stream": false,
        "temperature": 0.6,
        "top_p": 0.95,
        "max_completion_tokens": 128,
    }
    b, _ := json.Marshal(body)
    ep, err := p.endpoint("/v1/chat/completions")
    if err != nil { return "", "", err }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(b))
    if err != nil { return "", "", err }
    req.Header.Set("Authorization", "Bearer "+p.apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    hc := &http.Client{Timeout: 12 * time.Second}
    resp, err := hc.Do(req)
    if err != nil { return "", "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        var detail string
        if d, _ := io.ReadAll(resp.Body); len(d) > 0 { if len(d) > 400 { d = d[:400] }; detail = strings.TrimSpace(string(d)) }
        return "", "", fmt.Errorf("groq chat completions status=%d body=%s", resp.StatusCode, detail)
    }
    var out struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", "", err }
    if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
        return "", "", errors.New("empty structured output")
    }
    raw := strings.TrimSpace(out.Choices[0].Message.Content)
    var parsed struct { Branch string `json:"branch"`; Reason *string `json:"reason"` }
    if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
        if i := strings.Index(raw, "{"); i >= 0 {
            if j := strings.LastIndex(raw, "}"); j > i {
                frag := raw[i:j+1]
                _ = json.Unmarshal([]byte(frag), &parsed)
            }
        }
    }
    bname, ok := ai.SanitizeBranch(parsed.Branch)
    if !ok || bname == "" { bname = ai.FallbackBranch(title) }
    reason := ""; if parsed.Reason != nil { reason = strings.TrimSpace(*parsed.Reason) }
    return bname, reason, nil
}
