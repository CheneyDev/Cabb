-- 0007_add_report_repos.sql
-- Extend automation_configs to store multi-repo reporting configuration

ALTER TABLE automation_configs
    ADD COLUMN IF NOT EXISTS report_repos JSONB NOT NULL DEFAULT '[]';

UPDATE automation_configs
SET report_repos = '[]'::jsonb
WHERE report_repos IS NULL;
