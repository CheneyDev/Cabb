package openai

import (
    "cabb/internal/ai"
)

// New returns nil when not built with 'openai' tag to avoid pulling SDK.
func New(model, apiKey, baseURL string) ai.BranchNamer { return nil }

