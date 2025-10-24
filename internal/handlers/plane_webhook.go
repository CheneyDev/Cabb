package handlers

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "io"
    "net/http"
    "strings"

    "github.com/labstack/echo/v4"
)

// PlaneWebhook verifies signature, dedupes, and enqueues async processing.
func (h *Handler) PlaneWebhook(c echo.Context) error {
    body, err := io.ReadAll(c.Request().Body)
    if err != nil { return c.NoContent(http.StatusBadRequest) }
    sig := c.Request().Header.Get("X-Plane-Signature")
    if h.cfg.PlaneWebhookSecret != "" {
        if !verifyHMACSHA256(body, h.cfg.PlaneWebhookSecret, sig) {
            return c.NoContent(http.StatusUnauthorized)
        }
    }
    deliveryID := c.Request().Header.Get("X-Plane-Delivery")
    eventType := c.Request().Header.Get("X-Plane-Event")
    hsum := sha256.Sum256(body)
    sum := hex.EncodeToString(hsum[:])
    if h.dedupe != nil && h.dedupe.CheckAndMark("plane."+eventType, deliveryID, sum) {
        return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "duplicate"})
    }
    if hHasDB(h) && deliveryID != "" {
        dup, err := h.db.IsDuplicateDelivery(c.Request().Context(), "plane."+eventType, deliveryID, sum)
        if err == nil && dup {
            return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "duplicate"})
        }
        _ = h.db.UpsertEventDelivery(c.Request().Context(), "plane."+eventType, "incoming", deliveryID, sum, "queued")
    }
    go h.processPlaneWebhook(eventType, body, deliveryID)
    return c.JSON(http.StatusOK, map[string]any{"accepted": true, "delivery_id": deliveryID, "event_type": eventType, "status": "queued"})
}

func verifyHMACSHA256(body []byte, secret, got string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    want := hex.EncodeToString(mac.Sum(nil))
    got = strings.TrimSpace(got)
    if strings.HasPrefix(strings.ToLower(got), "sha256=") {
        got = got[len("sha256="):]
    }
    wantB, err1 := hex.DecodeString(want)
    gotB, err2 := hex.DecodeString(got)
    if err1 != nil || err2 != nil { return false }
    return hmac.Equal(wantB, gotB)
}

type planeWebhookEnvelope struct {
    Event       string         `json:"event"`
    Action      string         `json:"action"`
    WebhookID   string         `json:"webhook_id"`
    WorkspaceID string         `json:"workspace_id"`
    Data        map[string]any `json:"data"`
    Activity    struct {
        Field    string      `json:"field"`
        NewValue any         `json:"new_value"`
        OldValue any         `json:"old_value"`
        Actor    struct {
            ID          string `json:"id"`
            DisplayName string `json:"display_name"`
            Email       string `json:"email"`
        } `json:"actor"`
    } `json:"activity"`
}

func (h *Handler) processPlaneWebhook(event string, body []byte, deliveryID string) {
    var env planeWebhookEnvelope
    if err := json.Unmarshal(body, &env); err != nil { return }
    LogStructured("info", map[string]any{
        "event":        "plane.webhook",
        "delivery_id":  deliveryID,
        "x_event":      strings.ToLower(event),
        "action":       strings.ToLower(env.Action),
        "workspace_id": env.WorkspaceID,
    })
    switch strings.ToLower(event) {
    case "issue":
        h.handlePlaneIssueEvent(env, deliveryID)
    case "issue comment", "issue_comment":
        h.handlePlaneIssueComment(env, deliveryID)
    default:
        // ignore others for CNB integration
    }
}
