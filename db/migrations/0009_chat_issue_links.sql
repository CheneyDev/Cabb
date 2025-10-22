-- Map Lark chat (group) to a single bound Plane issue for command routing outside threads
CREATE TABLE IF NOT EXISTS chat_issue_links (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lark_chat_id text NOT NULL,
  lark_thread_id text,
  plane_issue_id uuid NOT NULL,
  plane_project_id uuid,
  workspace_slug text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (lark_chat_id)
);

CREATE INDEX IF NOT EXISTS idx_chat_issue_links_chat ON chat_issue_links(lark_chat_id);
