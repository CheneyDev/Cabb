-- 0008_report_notify_config.sql
-- Configuration for AI report notifications (daily/weekly/monthly)

CREATE TABLE IF NOT EXISTS report_notify_configs (
    id SERIAL PRIMARY KEY,
    notify_type VARCHAR(20) NOT NULL DEFAULT 'chat',  -- 'chat', 'users', 'departments'
    chat_id VARCHAR(255) DEFAULT '',
    user_ids JSONB NOT NULL DEFAULT '[]',             -- array of {open_id, name}
    department_ids JSONB NOT NULL DEFAULT '[]',       -- array of {id, name}
    daily_enabled BOOLEAN NOT NULL DEFAULT true,
    weekly_enabled BOOLEAN NOT NULL DEFAULT true,
    monthly_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed a single default row
INSERT INTO report_notify_configs (id, notify_type, chat_id, user_ids, department_ids, daily_enabled, weekly_enabled, monthly_enabled)
VALUES (
    1,
    'chat',
    '',
    '[]'::JSONB,
    '[]'::JSONB,
    true,
    true,
    true
)
ON CONFLICT (id) DO NOTHING;
