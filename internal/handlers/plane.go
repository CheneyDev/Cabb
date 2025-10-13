package handlers

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "net/http"
    "strings"

    "github.com/labstack/echo/v4"
)

func (h *Handler) PlaneOAuthStart(c echo.Context) error {
    // Placeholder: redirect to Plane OAuth install page.
    return c.JSON(http.StatusNotImplemented, map[string]string{
        "message": "Plane OAuth start is not implemented in scaffold.",
    })
}

func (h *Handler) PlaneOAuthCallback(c echo.Context) error {
    // Placeholder: exchange code/app_installation_id and store tokens.
    return c.JSON(http.StatusNotImplemented, map[string]string{
        "message": "Plane OAuth callback is not implemented in scaffold.",
    })
}

func (h *Handler) PlaneWebhook(c echo.Context) error {
    // Verify signature if provided
    body, err := io.ReadAll(c.Request().Body)
    if err != nil {
        return c.NoContent(http.StatusBadRequest)
    }
    sig := c.Request().Header.Get("X-Plane-Signature")
    if h.cfg.PlaneWebhookSecret != "" {
        if !verifyHMACSHA256(body, h.cfg.PlaneWebhookSecret, sig) {
            return c.NoContent(http.StatusUnauthorized)
        }
    }

    // Acknowledge; actual routing will be implemented later.
    deliveryID := c.Request().Header.Get("X-Plane-Delivery")
    eventType := c.Request().Header.Get("X-Plane-Event")
    return c.JSON(http.StatusOK, map[string]any{
        "accepted":    true,
        "delivery_id": deliveryID,
        "event_type":  eventType,
    })
}

func verifyHMACSHA256(body []byte, secret, got string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    want := hex.EncodeToString(mac.Sum(nil))
    // Some platforms prefix with "sha256=..."; accept either raw or prefixed.
    got = strings.TrimSpace(got)
    if strings.HasPrefix(strings.ToLower(got), "sha256=") {
        got = got[len("sha256="):]
    }
    // Constant-time compare
    wantB, err1 := hex.DecodeString(want)
    gotB, err2 := hex.DecodeString(got)
    if err1 != nil || err2 != nil {
        return false
    }
    return hmac.Equal(wantB, gotB)
}

