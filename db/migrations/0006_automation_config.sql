-- 0006_automation_config.sql
-- Global automation configuration for issue progress/report pipelines

CREATE TABLE IF NOT EXISTS automation_configs (
    id SERIAL PRIMARY KEY,
    target_repo_url TEXT NOT NULL DEFAULT '',
    target_repo_branch TEXT NOT NULL DEFAULT 'main',
    plane_statuses TEXT[] NOT NULL DEFAULT '{}',
    output_repo_url TEXT NOT NULL DEFAULT '',
    output_branch TEXT NOT NULL DEFAULT 'main',
    output_dir TEXT NOT NULL DEFAULT 'issue-progress',
    report_repo_slug TEXT NOT NULL DEFAULT '1024hub/plane-test',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed a single default row for convenience
INSERT INTO automation_configs (id, target_repo_url, target_repo_branch, plane_statuses, output_repo_url, output_branch, output_dir, report_repo_slug)
VALUES (
    1,
    '',
    'main',
    ARRAY[]::TEXT[],
    'https://cnb.cool/1024hub/plane-test',
    'main',
    'issue-progress',
    '1024hub/plane-test'
)
ON CONFLICT (id) DO NOTHING;
