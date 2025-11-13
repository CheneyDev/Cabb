package cerebras

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
    if err != nil { return "", err }
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
        return "", "", errors.New("missing CEREBRAS_API_KEY")
    }
    model := p.model
    if model == "" { model = "gpt-oss-120b" }

    schema := map[string]any{
        "type": "object",
        "properties": map[string]any{
            "branch": map[string]any{
                "type": "string",
                "pattern": "^(feat|fix|chore|docs|refactor|test|perf|ci|build|style)(/[a-z0-9][a-z0-9_/-]{0,60})$",
                "description": "lowercase Git branch name in prefix/slug format",
            },
        },
        "required": []string{"branch"},
        "additionalProperties": false,
    }

    sys := strings.Join([]string{
        "You generate concise Git branch names for software issues.",
        "Rules:",
        "- Respond with JSON matching the provided schema only.",
        "- branch must be lowercase.",
        "- Allowed prefixes: feat, fix, chore, docs, refactor, test, perf, ci, build, style.",
        "- Format: <prefix>/<slug> where slug uses [a-z0-9_/-], 2..60 chars after prefix/.",
        "- No punctuation, no emojis, no quotes.",
        "- Keep it short and meaningful.",
    }, "\n")
    usr := fmt.Sprintf("Title: %s\nDescription (may include HTML): %s\nReturn a fitting branch.", strings.TrimSpace(title), strings.TrimSpace(description))

    body := map[string]any{
        "model": model,
        "messages": []map[string]any{{"role": "system", "content": sys}, {"role": "user", "content": usr}},
        "response_format": map[string]any{
            "type": "json_schema",
            "json_schema": map[string]any{"name": "branch_schema", "strict": true, "schema": schema},
        },
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
        if d, _ := io.ReadAll(resp.Body); len(d) > 0 {
            if len(d) > 400 { d = d[:400] }
            detail = strings.TrimSpace(string(d))
        }
        return "", "", fmt.Errorf("cerebras chat completions status=%d body=%s", resp.StatusCode, detail)
    }
    var out struct{ Choices []struct{ Message struct{ Content string `json:"content"` } `json:"message"` } `json:"choices"` }
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return "", "", err
    }
    if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
        return "", "", errors.New("empty response content")
    }
    var parsed struct{ Branch string `json:"branch"` }
    if err := json.Unmarshal([]byte(out.Choices[0].Message.Content), &parsed); err != nil {
        return "", "", err
    }
    branch, ok := ai.SanitizeBranch(parsed.Branch)
    if !ok || branch == "" {
        branch = ai.FallbackBranch(title)
    }
    return branch, "", nil
}
