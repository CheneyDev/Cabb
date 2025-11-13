package ai

import (
    "context"
    "regexp"
    "strings"
)

// BranchNamer suggests a Git branch name for an issue.
type BranchNamer interface {
    SuggestBranchName(ctx context.Context, title, description string) (branch string, reason string, err error)
}

var (
    // Conservative pattern for branch names (lowercase); allow '/','-','_' separators.
    branchPattern = regexp.MustCompile(`^(feat|fix|chore|docs|refactor|test|perf|ci|build|style)(/[a-z0-9][a-z0-9_/-]{0,60})$`)
)

// SanitizeBranch coerces candidate to a safe, lowercased branch and validates against branchPattern.
func SanitizeBranch(candidate string) (string, bool) {
    s := strings.TrimSpace(strings.ToLower(candidate))
    // collapse spaces and disallowed chars
    s = strings.ReplaceAll(s, " ", "-")
    s = strings.ReplaceAll(s, "--", "-")
    if branchPattern.MatchString(s) {
        return s, true
    }
    return s, false
}

// FallbackBranch constructs a branch slug from a title when AI is unavailable.
func FallbackBranch(title string) string {
    t := strings.ToLower(strings.TrimSpace(title))
    // very simple slug: keep alnum and replace others with '-'
    b := make([]rune, 0, len(t))
    for _, r := range t {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
            b = append(b, r)
            continue
        }
        if r == ' ' || r == '-' || r == '_' || r == '/' {
            b = append(b, '-')
            continue
        }
        // skip others
    }
    slug := strings.Trim(betweenDashes(string(b)), "/-")
    if slug == "" {
        slug = "task"
    }
    return "feat/" + slug
}

func betweenDashes(s string) string {
    for strings.Contains(s, "--") {
        s = strings.ReplaceAll(s, "--", "-")
    }
    return s
}

