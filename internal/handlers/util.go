package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/labstack/echo/v4"
)

// truncate returns at most n bytes of s (best-effort, not rune-safe for logs only).
func truncate(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n]
}

// writeError writes a unified error response and logs a structured entry.
func writeError(c echo.Context, status int, code, message string, details map[string]any) error {
    if details == nil {
        details = map[string]any{}
    }
    reqID := c.Response().Header().Get("X-Request-ID")
    if reqID == "" {
        reqID = c.Request().Header.Get("X-Request-ID")
    }
    ep := c.Path()
    if ep == "" {
        ep = c.Request().URL.Path
    }
    src := DetectSource(ep)
    errLog := map[string]any{
        "time":       time.Now().UTC().Format(time.RFC3339),
        "level":      "error",
        "request_id": reqID,
        "endpoint":   ep,
        "source":     src,
        "status":     status,
        "result":     "error",
        "error": map[string]any{
            "code":    code,
            "message": message,
        },
    }
    if b, e := json.Marshal(errLog); e == nil {
        log.Println(string(b))
    }
    return c.JSON(status, map[string]any{
        "error": map[string]any{
            "code":       code,
            "message":    message,
            "details":    details,
            "request_id": reqID,
        },
    })
}

// stripTags removes HTML tags (best-effort, sufficient for short previews).
func stripTags(s string) string {
    out := make([]rune, 0, len(s))
    inTag := false
    for _, r := range s {
        switch r {
        case '<':
            inTag = true
        case '>':
            inTag = false
        default:
            if !inTag {
                out = append(out, r)
            }
        }
    }
    return strings.TrimSpace(string(out))
}

