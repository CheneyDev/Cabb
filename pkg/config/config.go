package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    Port              int
    DatabaseURL       string
    Timezone          string

    // Plane
    PlaneBaseURL      string
    PlaneClientID     string
    PlaneClientSecret string
    PlaneRedirectURI  string
    PlaneWebhookSecret string

    // Feishu (Lark)
    LarkAppID         string
    LarkAppSecret     string
    LarkEncryptKey    string
    LarkVerificationToken string

    // CNB
    CNBAppToken       string
    IntegrationToken  string

    // Crypto
    EncryptionKey     string
}

func FromEnv() Config {
    cfg := Config{
        Port:               intFromEnv("PORT", 8080),
        DatabaseURL:        os.Getenv("DATABASE_URL"),
        Timezone:           strFromEnv("TIMEZONE", "Local"),

        PlaneBaseURL:       strFromEnv("PLANE_BASE_URL", "https://api.plane.so"),
        PlaneClientID:      os.Getenv("PLANE_CLIENT_ID"),
        PlaneClientSecret:  os.Getenv("PLANE_CLIENT_SECRET"),
        PlaneRedirectURI:   strFromEnv("PLANE_REDIRECT_URI", "http://localhost:8080/plane/oauth/callback"),
        PlaneWebhookSecret: os.Getenv("PLANE_WEBHOOK_SECRET"),

        LarkAppID:              os.Getenv("LARK_APP_ID"),
        LarkAppSecret:          os.Getenv("LARK_APP_SECRET"),
        LarkEncryptKey:         os.Getenv("LARK_ENCRYPT_KEY"),
        LarkVerificationToken:  os.Getenv("LARK_VERIFICATION_TOKEN"),

        CNBAppToken:        os.Getenv("CNB_APP_TOKEN"),
        IntegrationToken:   os.Getenv("INTEGRATION_TOKEN"),

        EncryptionKey:      os.Getenv("ENCRYPTION_KEY"),
    }
    return cfg
}

func (c Config) Address() string {
    return fmt.Sprintf(":%d", c.Port)
}

func intFromEnv(key string, def int) int {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    if n, err := strconv.Atoi(v); err == nil {
        return n
    }
    return def
}

func strFromEnv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

