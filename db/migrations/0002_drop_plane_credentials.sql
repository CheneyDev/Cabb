-- Migration: Drop plane_credentials table
-- Reason: Simplified to use PLANE_SERVICE_TOKEN env variable for global token
-- Date: 2025-10-31

DROP TABLE IF EXISTS plane_credentials;

-- Note: No need to maintain per-workspace tokens in database.
-- Use PLANE_SERVICE_TOKEN environment variable for the global Plane API access token.
