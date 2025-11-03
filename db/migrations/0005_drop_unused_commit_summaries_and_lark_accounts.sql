-- Migration: Drop unused tables commit_summaries and lark_accounts
-- Reason: Not referenced anywhere in the current codebase; remove to simplify schema
-- Date: 2025-11-03

-- idempotent cleanup
DROP TABLE IF EXISTS commit_summaries;
DROP TABLE IF EXISTS lark_accounts;

-- Note: indexes on these tables (if any) are dropped automatically with the table.

