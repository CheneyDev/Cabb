package handlers

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/labstack/echo/v4"
)

// DetectSource returns a coarse-grained source tag from a request path.
func DetectSource(path string) string {
    switch {
    case strings.HasPrefix(path, "/webhooks/plane"):
        return "plane.webhook"
    case strings.HasPrefix(path, "/ingest/cnb"):
        return "cnb.ingest"
    case strings.HasPrefix(path, "/webhooks/lark"):
        return "lark.webhook"
    case strings.HasPrefix(path, "/admin/"):
        return "admin"
    case strings.HasPrefix(path, "/jobs/"):
        return "jobs"
    default:
        return "http"
    }
}

// StructuredLogger is a middleware that logs every request in structured JSON.
// Fields: time, request_id, method, endpoint, path, status, latency_ms, result, source, remote_ip, user_agent
func StructuredLogger() echo.MiddlewareFunc {
    type entry struct {
        Time       string `json:"time"`
        Level      string `json:"level"`
        RequestID  string `json:"request_id"`
        Method     string `json:"method"`
        Endpoint   string `json:"endpoint"`
        Path       string `json:"path"`
        Status     int    `json:"status"`
        LatencyMS  int64  `json:"latency_ms"`
        Result     string `json:"result"`
        Source     string `json:"source"`
        RemoteIP   string `json:"remote_ip"`
        UserAgent  string `json:"user_agent"`
    }

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            err := next(c)

            // Build log entry
            rid := c.Response().Header().Get(echo.HeaderXRequestID)
            if rid == "" { rid = c.Request().Header.Get(echo.HeaderXRequestID) }
            ep := c.Path()
            if ep == "" { ep = c.Request().URL.Path }
            status := c.Response().Status
            // Result classification
            result := "success"
            switch {
            case status >= http.StatusInternalServerError:
                result = "server_error"
            case status >= http.StatusBadRequest:
                result = "client_error"
            }
            e := entry{
                Time:      time.Now().UTC().Format(time.RFC3339),
                Level:     "info",
                RequestID: rid,
                Method:    c.Request().Method,
                Endpoint:  ep,
                Path:      c.Request().URL.Path,
                Status:    status,
                LatencyMS: time.Since(start).Milliseconds(),
                Result:    result,
                Source:    DetectSource(ep),
                RemoteIP:  c.RealIP(),
                UserAgent: c.Request().UserAgent(),
            }
            b, _ := json.Marshal(e)
            log.Println(string(b))
            return err
        }
    }
}

// RequestID attaches a random request id if not present and echoes it back.
func RequestID() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            rid := c.Request().Header.Get(echo.HeaderXRequestID)
            if rid == "" {
                b := make([]byte, 16)
                if _, err := rand.Read(b); err == nil {
                    rid = hex.EncodeToString(b)
                } else {
                    rid = time.Now().UTC().Format("20060102T150405.000Z07:00")
                }
            }
            c.Response().Header().Set(echo.HeaderXRequestID, rid)
            return next(c)
        }
    }
}

// Recover recovers from panics and returns a structured 500 error.
func Recover() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) (err error) {
            defer func() {
                if r := recover(); r != nil {
                    // Log structured panic
                    rid := c.Response().Header().Get(echo.HeaderXRequestID)
                    ep := c.Path(); if ep == "" { ep = c.Request().URL.Path }
                    panicLog := map[string]any{
                        "time":       time.Now().UTC().Format(time.RFC3339),
                        "level":      "error",
                        "request_id": rid,
                        "endpoint":   ep,
                        "source":     DetectSource(ep),
                        "status":     http.StatusInternalServerError,
                        "result":     "error",
                        "panic":      r,
                    }
                    if b, e := json.Marshal(panicLog); e == nil { log.Println(string(b)) }
                    _ = writeError(c, http.StatusInternalServerError, "internal_error", "服务内部错误", nil)
                }
            }()
            return next(c)
        }
    }
}

// LogStructured prints a single structured JSON log entry with given level and fields.
// It adds RFC3339 UTC time and level if not present.
func LogStructured(level string, fields map[string]any) {
    if fields == nil {
        fields = map[string]any{}
    }
    if _, ok := fields["time"]; !ok {
        fields["time"] = time.Now().UTC().Format(time.RFC3339)
    }
    if _, ok := fields["level"]; !ok {
        fields["level"] = level
    }
    b, err := json.Marshal(fields)
    if err != nil {
        log.Printf("{\"level\":\"%s\",\"message\":\"log marshal failed\"}", level)
        return
    }
    log.Println(string(b))
}
