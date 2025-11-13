//go:build openai

package openai

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "strings"

    ai "cabb/internal/ai"
    openaisdk "github.com/openai/openai-go"
    "github.com/openai/openai-go/option"
)

type provider struct {
    model   string
    apiKey  string
    baseURL string
}

// New returns an OpenAI-backed BranchNamer (requires build tag 'openai').
func New(model, apiKey, baseURL string) ai.BranchNamer {
    return &provider{model: strings.TrimSpace(model), apiKey: strings.TrimSpace(apiKey), baseURL: strings.TrimSpace(baseURL)}
}

func (p *provider) client() *openaisdk.Client {
    if p.baseURL != "" {
        return openaisdk.NewClient(option.WithAPIKey(p.apiKey), option.WithBaseURL(p.baseURL))
    }
    return openaisdk.NewClient(option.WithAPIKey(p.apiKey))
}

func (p *provider) SuggestBranchName(ctx context.Context, title, description string) (string, string, error) {
    if strings.TrimSpace(p.apiKey) == "" {
        return "", "", errors.New("missing OPENAI_API_KEY")
    }
    model := p.model
    if model == "" { model = "gpt-4o-mini" }
    sys := strings.Join([]string{
        "You generate concise Git branch names for software issues.",
        "Rules:",
        "- Return only JSON via structured output schema.",
        "- branch must be lowercase.",
        "- Allowed prefixes: feat, fix, chore, docs, refactor, test, perf, ci, build, style.",
        "- Format: <prefix>/<slug> where slug uses [a-z0-9_/-], no spaces, 2..60 chars after prefix/.",
        "- No punctuation, no emojis, no quotes.",
        "- Keep it short and meaningful.",
        "- If information is insufficient or missing, you must still return a branch.",
        "- Output must be ASCII; transliterate or simplify non-ASCII to ASCII.",
        "- When unsure, use prefix 'feat' and derive a short ASCII slug from the title; if the title is empty, use 'task' as the slug.",
        "- Do not refuse or add explanations; output JSON only.",
    }, "\n")
    user := fmt.Sprintf("Title: %s\nDescription (may include HTML): %s\nReturn a fitting branch.", strings.TrimSpace(title), strings.TrimSpace(description))
    schema := map[string]any{
        "type": "object",
        "properties": map[string]any{
            "branch": map[string]any{"type": "string", "pattern": "^(feat|fix|chore|docs|refactor|test|perf|ci|build|style)(/[a-z0-9][a-z0-9_/-]{0,60})$"},
            "reason": map[string]any{"anyOf": []any{map[string]any{"type": "string"}, map[string]any{"type": "null"}}},
        },
        "required": []string{"branch", "reason"},
        "additionalProperties": false,
    }
    client := p.client()
    resp, err := client.Responses.New(ctx, openaisdk.ResponseNewParams{
        Model: openaisdk.F(model),
        Input: openaisdk.F([]openaisdk.Input{openaisdk.TextInput{Text: sys + "\n\n" + user}}),
        ResponseFormat: openaisdk.F(openaisdk.ResponseFormat{Type: openaisdk.F("json_schema"), JsonSchema: openaisdk.F(openaisdk.ResponseFormatJSONSchema{Name: openaisdk.F("branch_name"), Strict: openaisdk.Bool(true), Schema: openaisdk.F(schema)})}),
    })
    if err != nil { return "", "", err }
    if strings.EqualFold(resp.Status, "incomplete") {
        return "", "", fmt.Errorf("openai incomplete: %v", resp.IncompleteDetails)
    }
    var jsonStr string
    if len(resp.Output) > 0 {
        for _, item := range resp.Output {
            for _, c := range item.Content {
                if c.Type == "refusal" { return "", "", errors.New("openai refusal") }
                if c.Type == "output_text" && c.Text != nil { jsonStr += c.Text.Value }
                if c.Type == "json" && c.JSON != nil { if s, ok := c.JSON.(string); ok { jsonStr = s } }
            }
        }
    }
    if strings.TrimSpace(jsonStr) == "" { return "", "", errors.New("empty structured output") }
    var out struct{ Branch string `json:"branch"`; Reason *string `json:"reason"` }
    if err := json.Unmarshal([]byte(jsonStr), &out); err != nil { return "", "", err }
    b, ok := ai.SanitizeBranch(out.Branch)
    if !ok || b == "" { b = ai.FallbackBranch(title) }
    reason := ""; if out.Reason != nil { reason = strings.TrimSpace(*out.Reason) }
    return b, reason, nil
}
